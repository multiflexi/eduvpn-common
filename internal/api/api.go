// Package api implements version 3 of the eduVPN api: https://docs.eduvpn.org/server/v3/api.html
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/jwijenbergh/eduoauth-go"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/eduvpn/eduvpn-common/internal/api/endpoints"
	"github.com/eduvpn/eduvpn-common/internal/api/profiles"
	httpw "github.com/eduvpn/eduvpn-common/internal/http"
	"github.com/eduvpn/eduvpn-common/internal/log"
	"github.com/eduvpn/eduvpn-common/internal/wireguard"
	"github.com/eduvpn/eduvpn-common/types/protocol"
	"github.com/eduvpn/eduvpn-common/types/server"
)

// Callbacks is the API callback interface
// It is used to trigger authorization and forward token updates
type Callbacks interface {
	// TriggerAuth is called when authorization should be triggered
	TriggerAuth(context.Context, string, bool) (string, error)
	// AuthDone is called when authorization has just completed
	AuthDone(string, server.Type)
	// TokensUpdates is called when tokens are updated
	TokensUpdated(string, server.Type, eduoauth.Token)
}

// ServerData is the data for a server that is passed to the API struct
type ServerData struct {
	// ID is the identifier for the server
	ID string
	// Type is the type of server
	Type server.Type
	// BaseWK is the base well-known endpoint
	BaseWK string
	// BaseAuthWK is the base well-known endpoint for authorization. This is only different in case of secure internet
	BaseAuthWK string
	// ProcessAuth processes the OAuth authorization
	ProcessAuth func(context.Context, string) (string, error)
	// DisableAuthorize indicates whether or not new authorization requests should be disabled
	DisableAuthorize bool
	// transport is the HTTP transport, only used for testing currently
	transport http.RoundTripper
}

// Transport returns the transport to be used for the server
// By default it uses the transport from internal/http DefaultTransport
func (s *ServerData) Transport() http.RoundTripper {
	if s.transport == nil {
		return httpw.DefaultTransport
	}
	return s.transport
}

// API is the top-level struct that each method is defined on
type API struct {
	cb Callbacks
	// oauth is the oauth object
	oauth *eduoauth.OAuth
	// Data is the server data
	Data ServerData
}

// NewAPI creates a new API object by creating an OAuth object
func NewAPI(ctx context.Context, clientID string, sd ServerData, cb Callbacks, tokens *eduoauth.Token) (*API, error) {
	cr := customRedirect(clientID)
	// Construct OAuth

	transp := sd.Transport()
	o := eduoauth.OAuth{
		ClientID: clientID,
		EndpointFunc: func(ctx context.Context) (*eduoauth.EndpointResponse, error) {
			ep, err := GetEndpointCache().Get(ctx, sd.BaseAuthWK, transp)
			if err != nil {
				return nil, err
			}
			return &eduoauth.EndpointResponse{
				AuthorizationURL: ep.API.V3.Authorization,
				TokenURL:         ep.API.V3.Token,
			}, nil
		},
		CustomRedirect: cr,
		RedirectPath:   "/callback",
		TokensUpdated: func(tok eduoauth.Token) {
			cb.TokensUpdated(sd.ID, sd.Type, tok)
		},
		Transport: transp,
		UserAgent: httpw.UserAgent,
	}

	if tokens != nil {
		o.UpdateTokens(*tokens)
	}

	api := &API{
		cb:    cb,
		oauth: &o,
		Data:  sd,
	}
	err := api.authorize(ctx)
	if err != nil {
		return nil, err
	}
	return api, nil
}

// ErrAuthorizeDisabled is returned when authorization is disabled but is needed to complete
var ErrAuthorizeDisabled = errors.New("cannot authorize as re-authorization is disabled")

func (a *API) authorize(ctx context.Context) (err error) {
	_, err = a.oauth.AccessToken(ctx)
	// already authorized
	if err == nil {
		return nil
	}

	// otherwise check if invalid tokens,
	// if not then something else is wrong with the API
	// return an error
	tErr := &eduoauth.TokensInvalidError{}
	if !errors.As(err, &tErr) {
		return err
	}

	if a.Data.DisableAuthorize {
		return ErrAuthorizeDisabled
	}

	defer func() {
		if err == nil {
			a.cb.AuthDone(a.Data.ID, a.Data.Type)
		}
	}()

	scope := "config"
	url, err := a.oauth.AuthURL(ctx, scope)
	if err != nil {
		return err
	}
	if a.Data.ProcessAuth != nil {
		url, err = a.Data.ProcessAuth(ctx, url)
		if err != nil {
			return err
		}
	}
	// We expect an uri if custom redirect is non empty
	uri, err := a.cb.TriggerAuth(ctx, url, a.oauth.CustomRedirect != "")
	if err != nil {
		return err
	}
	// The uri is only given here if a custom redirect is done
	err = a.oauth.Exchange(ctx, uri)
	if err != nil {
		return err
	}
	return nil
}

func (a *API) authorized(ctx context.Context, method string, endpoint string, opts *httpw.OptionalParams) (http.Header, []byte, error) {
	ep, err := GetEndpointCache().Get(ctx, a.Data.BaseWK, a.Data.Transport())
	if err != nil {
		return nil, nil, err
	}
	u := ep.API.V3.API + endpoint

	// TODO: Cache HTTP client?
	httpC := httpw.NewClient(a.oauth.NewHTTPClient())
	return httpC.Do(ctx, method, u, opts)
}

func (a *API) authorizedRetry(ctx context.Context, method string, endpoint string, opts *httpw.OptionalParams) (http.Header, []byte, error) {
	h, body, err := a.authorized(ctx, method, endpoint, opts)
	if err == nil {
		return h, body, nil
	}

	statErr := &httpw.StatusError{}
	// Only retry authorized if we get an HTTP 401
	// TODO: Can the OAuth client handle this instead?
	if errors.As(err, &statErr) && statErr.Status == 401 {
		log.Logger.Debugf("Got a 401 error after HTTP method: %s, endpoint: %s. Marking token as expired...", method, endpoint)
		// Mark the token as expired and retry, so we trigger the refresh flow
		a.oauth.SetTokenExpired()
		h, body, err = a.authorized(ctx, method, endpoint, opts)
	}
	// Tokens is invalid we need to renew and authorize again
	tErr := &eduoauth.TokensInvalidError{}
	if err != nil && errors.As(err, &tErr) {
		// Mark the token as invalid and retry, so we trigger the authorization flow
		a.oauth.SetTokenRenew()
		log.Logger.Debugf("the tokens were invalid, trying again...")
		if autherr := a.authorize(ctx); autherr != nil {
			return nil, nil, autherr
		}
		return a.authorized(ctx, method, endpoint, opts)
	}
	return h, body, err
}

// Disconnect disconnects a client from the server by sending a /disconnect API call
// This cleans up resources such as WireGuard IP allocation
func (a *API) Disconnect(ctx context.Context) error {
	_, _, err := a.authorized(ctx, http.MethodPost, "/disconnect", &httpw.OptionalParams{Timeout: 5 * time.Second})
	return err
}

// Info does the /info API call
func (a *API) Info(ctx context.Context) (*profiles.Info, error) {
	_, body, err := a.authorizedRetry(ctx, http.MethodGet, "/info", nil)
	if err != nil {
		return nil, fmt.Errorf("failed API /info: %w", err)
	}
	p := profiles.Info{}
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed API /info: %w", err)
	}
	return &p, nil
}

// ConnectData is the data that is returned when the /connect call completes without error
type ConnectData struct {
	// Configuration is the VPN configuration
	Configuration string
	// Protocol tells us what protocol it is, OpenVPN or WireGuard (proxied or not)
	Protocol protocol.Protocol
	// Expires tells us when this configuration expires
	Expires time.Time
	// Proxy is filled when WireGuard is proxied
	Proxy *wireguard.Proxy
}

// see https://github.com/eduvpn/documentation/blob/v3/API.md#request-1
func boolToYesNo(preferTCP bool) string {
	if preferTCP {
		return "yes"
	}
	return "no"
}

func protocolFromCT(ct string) (protocol.Protocol, error) {
	switch ct {
	case "application/x-wireguard-profile":
		return protocol.WireGuard, nil
	case "application/x-wireguard+tcp-profile":
		return protocol.WireGuardProxy, nil
	case "application/x-openvpn-profile":
		return protocol.OpenVPN, nil
	}
	return protocol.Unknown, fmt.Errorf("invalid content type: %s", ct)
}

// ErrNoProtocols is returned when a connect call is given with an empty protocol slice
var ErrNoProtocols = errors.New("no protocols supplied")

// ErrUnknownProtocol is returned when the client in a connect gives an unknown protocol
var ErrUnknownProtocol = errors.New("unknown protocol supplied")

// Connect sends a /connect to an eduVPN server
// `ctx` is the context used for cancellation
// protos is the list of protocols supported and wanted by the client
func (a *API) Connect(ctx context.Context, prof profiles.Profile, protos []protocol.Protocol, pTCP bool) (*ConnectData, error) {
	hdrs := http.Header{
		"content-type": {"application/x-www-form-urlencoded"},
	}
	uv := url.Values{
		"profile_id": {prof.ID},
	}

	if len(protos) == 0 {
		return nil, ErrNoProtocols
	}

	var wgKey *wgtypes.Key

	// Loop over the protocols and set the correct headers and values
	for _, p := range protos {
		switch p {
		case protocol.WireGuard:
			gk, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				return nil, err
			}
			wgKey = &gk
			// Set the public key
			pubkey := wgKey.PublicKey()
			uv.Set("public_key", pubkey.String())
			hdrs.Add("accept", "application/x-wireguard-profile")
			hdrs.Add("accept", "application/x-wireguard+tcp-profile")
		case protocol.OpenVPN:
			hdrs.Add("accept", "application/x-openvpn-profile")
		default:
			return nil, ErrUnknownProtocol
		}
	}
	// set prefer TCP
	uv.Set("prefer_tcp", boolToYesNo(pTCP))

	// Construct the parameters
	params := &httpw.OptionalParams{Headers: hdrs, Body: uv}
	h, body, err := a.authorizedRetry(ctx, http.MethodPost, "/connect", params)
	if err != nil {
		return nil, fmt.Errorf("failed API /connect call: %w", err)
	}

	// Parse expiry
	expH := h.Get("expires")
	if expH == "" {
		return nil, errors.New("the server did not give an expires header")
	}
	expT, err := http.ParseTime(expH)
	if err != nil {
		return nil, fmt.Errorf("failed parsing expiry time: %w", err)
	}

	vpnCfg := string(body)
	// Parse content type
	contentH := h.Get("content-type")
	proto, err := protocolFromCT(contentH)
	if err != nil {
		return nil, err
	}

	if proto == protocol.OpenVPN {
		// ensure scripts are not ran by default by append script-security 0 to the config
		vpnCfg += "\nscript-security 0"
		return &ConnectData{
			Configuration: vpnCfg,
			Protocol:      proto,
			Expires:       expT,
		}, nil
	}

	vpnCfg, proxy, err := wireguard.Config(vpnCfg, wgKey, proto == protocol.WireGuardProxy)
	if err != nil {
		return nil, err
	}
	return &ConnectData{
		Configuration: vpnCfg,
		Protocol:      proto,
		Expires:       expT,
		Proxy:         proxy,
	}, nil
}

func getEndpoints(ctx context.Context, url string, tp http.RoundTripper) (*endpoints.Endpoints, error) {
	uStr, err := httpw.JoinURLPath(url, "/.well-known/vpn-user-portal")
	if err != nil {
		return nil, err
	}
	httpC := httpw.NewClient(nil)
	httpC.Client.Transport = tp
	_, body, err := httpC.Get(ctx, uStr)
	if err != nil {
		return nil, fmt.Errorf("failed getting server endpoints with error: %w", err)
	}

	ep := endpoints.Endpoints{}
	if err = json.Unmarshal(body, &ep); err != nil {
		return nil, fmt.Errorf("failed getting server endpoints with error: %w", err)
	}
	err = ep.Validate()
	if err != nil {
		return nil, err
	}
	return &ep, nil
}

// OAuthLogger is defined here to update the internal logger
// for the eduoauth library
type OAuthLogger struct{}

// Logf logs a message with parameters
func (ol *OAuthLogger) Logf(msg string, params ...interface{}) {
	log.Logger.Debugf(msg, params...)
}

// Log logs a message
func (ol *OAuthLogger) Log(msg string) {
	log.Logger.Debugf("%s", msg)
}

func init() {
	eduoauth.UpdateLogger(&OAuthLogger{})
}
