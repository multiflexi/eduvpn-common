package eduvpn

import (
	"fmt"

	"github.com/jwijenbergh/eduvpn-common/internal/config"
	"github.com/jwijenbergh/eduvpn-common/internal/discovery"
	"github.com/jwijenbergh/eduvpn-common/internal/fsm"
	"github.com/jwijenbergh/eduvpn-common/internal/log"
	"github.com/jwijenbergh/eduvpn-common/internal/oauth"
	"github.com/jwijenbergh/eduvpn-common/internal/server"
	"github.com/jwijenbergh/eduvpn-common/internal/types"
	"github.com/jwijenbergh/eduvpn-common/internal/util"
)

type VPNState struct {
	// The chosen server
	Servers server.Servers `json:"servers"`

	// The list of servers and organizations from disco
	Discovery discovery.Discovery `json:"-"`

	// The fsm
	FSM fsm.FSM `json:"-"`

	// The logger
	Logger log.FileLogger `json:"-"`

	// The config
	Config config.Config `json:"-"`

	// Whether to enable debugging
	Debug bool `json:"-"`
}

func (state *VPNState) GetSavedServers() string {
	serversJSON, serversJSONErr := state.Servers.GetServersConfiguredJSON()

	if serversJSONErr != nil {
		return "{}"
	}

	return serversJSON
}

func (state *VPNState) Register(name string, directory string, stateCallback func(string, string, string), debug bool) error {
	errorMessage := "failed to register with the GO library"
	if !state.FSM.InState(fsm.DEREGISTERED) {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: fsm.DeregisteredError{}.CustomError()}
	}
	// Initialize the logger
	logLevel := log.LOG_WARNING

	if debug {
		logLevel = log.LOG_INFO
	}

	loggerErr := state.Logger.Init(logLevel, name, directory)
	if loggerErr != nil {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: loggerErr}
	}

	// Initialize the FSM
	state.FSM.Init(name, stateCallback, directory, debug)
	state.Debug = debug

	// Initialize the Config
	state.Config.Init(name, directory)

	// Try to load the previous configuration
	if state.Config.Load(&state) != nil {
		// This error can be safely ignored, as when the config does not load, the struct will not be filled
		state.Logger.Log(log.LOG_INFO, "Previous configuration not found")
	}

	// Go to the No Server state with the saved servers
	state.FSM.GoTransitionWithData(fsm.NO_SERVER, state.GetSavedServers(), false)

	state.GetDiscoServers()
	state.GetDiscoOrganizations()
	return nil
}

func (state *VPNState) Deregister() error {
	// Close the log file
	state.Logger.Close()

	// Save the config
	state.Config.Save(&state)

	// Empty out the state
	*state = VPNState{}
	return nil
}

func (state *VPNState) GoBack() error {
	errorMessage := "failed to go back"
	if state.FSM.InState(fsm.DEREGISTERED) {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: fsm.DeregisteredError{}.CustomError()}
	}

	// FIXME: Abitrary back transitions don't work because we need the approriate data
	state.FSM.GoTransitionWithData(fsm.NO_SERVER, state.GetSavedServers(), false)
	// state.FSM.GoBack()
	return nil
}

func (state *VPNState) getConfig(chosenServer server.Server, forceTCP bool) (string, string, error) {
	errorMessage := "failed to get a configuration for OpenVPN/Wireguard"
	if state.FSM.InState(fsm.DEREGISTERED) {
		return "", "", &types.WrappedErrorMessage{Message: errorMessage, Err: fsm.DeregisteredError{}.CustomError()}
	}

	// Relogin with oauth
	// This moves the state to authorized
	if server.NeedsRelogin(chosenServer) {
		loginErr := server.Login(chosenServer)

		if loginErr != nil {
			// We are possibly in oauth started
			// Go back
			state.GoBack()
			return "", "", &types.WrappedErrorMessage{Message: errorMessage, Err: loginErr}
		}
	} else { // OAuth was valid, ensure we are in the authorized state
		state.FSM.GoTransition(fsm.AUTHORIZED)
	}

	state.FSM.GoTransition(fsm.REQUEST_CONFIG)

	config, configType, configErr := server.GetConfig(chosenServer, forceTCP)

	if configErr != nil {
		// Go back
		state.GoBack()
		return "", "", &types.WrappedErrorMessage{Message: errorMessage, Err: configErr}
	}

	// Signal the server display info
	state.FSM.GoTransitionWithData(fsm.HAS_CONFIG, state.getServerInfoData(), false)

	// Save the config
	state.Config.Save(&state)

	return config, configType, nil
}

func (state *VPNState) SetSecureLocation(countryCode string) error {
	server, serverErr := state.Discovery.GetServerByCountryCode(countryCode, "secure_internet")

	if serverErr != nil {
		return &types.WrappedErrorMessage{Message: "failed asking secure location", Err: serverErr}
	}

	state.Servers.SetSecureLocation(server, &state.FSM)
	return nil
}

func (state *VPNState) askSecureLocation() error {
	errorMessage := "failed asking Secure Internet location"
	locations, locationsErr := state.Discovery.GetSecureLocationList()
	if locationsErr != nil {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: locationsErr}
	}

	// Ask for the location in the callback
	state.FSM.GoTransitionWithData(fsm.ASK_LOCATION, locations, false)
	return nil
}

func (state *VPNState) addSecureInternetHomeServer(orgID string) (server.Server, error) {
	errorMessage := fmt.Sprintf("failed adding Secure Internet home server with organization ID %s", orgID)
	// Get the secure internet URL from discovery
	secureOrg, secureServer, discoErr := state.Discovery.GetSecureHomeArgs(orgID)
	if discoErr != nil {
		return nil, &types.WrappedErrorMessage{Message: errorMessage, Err: discoErr}
	}

	// Add the secure internet server
	server, serverErr := state.Servers.AddSecureInternet(secureOrg, secureServer, &state.FSM)

	if serverErr != nil {
		return nil, &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	var locationErr error

	if !state.Servers.HasSecureLocation() {
		locationErr = state.askSecureLocation()
	} else {
		// reinitialize
		locationErr = state.SetSecureLocation(state.Servers.GetSecureLocation())
	}

	if locationErr != nil {
		return nil, &types.WrappedErrorMessage{Message: errorMessage, Err: locationErr}
	}

	return server, nil
}

func (state *VPNState) GetConfigSecureInternet(orgID string, forceTCP bool) (string, string, error) {
	errorMessage := fmt.Sprintf("failed getting a configuration for Secure Internet organization %s", orgID)
	state.FSM.GoTransition(fsm.LOADING_SERVER)
	server, serverErr := state.addSecureInternetHomeServer(orgID)

	if serverErr != nil {
		return "", "", &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	state.FSM.GoTransition(fsm.CHOSEN_SERVER)

	return state.getConfig(server, forceTCP)
}

func (state *VPNState) addInstituteServer(url string) (server.Server, error) {
	errorMessage := fmt.Sprintf("failed adding Institute Access server with url %s", url)
	instituteServer, discoErr := state.Discovery.GetServerByURL(url, "institute_access")
	if discoErr != nil {
		return nil, &types.WrappedErrorMessage{Message: errorMessage, Err: discoErr}
	}
	// Add the secure internet server
	server, serverErr := state.Servers.AddInstituteAccessServer(instituteServer, &state.FSM)

	if serverErr != nil {
		return nil, &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	state.FSM.GoTransition(fsm.CHOSEN_SERVER)

	return server, nil
}

func (state *VPNState) addCustomServer(url string) (server.Server, error) {
	errorMessage := fmt.Sprintf("failed adding Custom server with url %s", url)

	customServer := &types.DiscoveryServer{BaseURL: url, DisplayName: map[string]string{"en": url}, Type: "custom_server"}

	// A custom server is just an institute access server under the hood
	server, serverErr := state.Servers.AddCustomServer(customServer, &state.FSM)

	if serverErr != nil {
		return nil, &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	state.FSM.GoTransition(fsm.CHOSEN_SERVER)

	return server, nil
}

func (state *VPNState) GetConfigInstituteAccess(url string, forceTCP bool) (string, string, error) {
	errorMessage := fmt.Sprintf("failed getting a configuration for Institute Access %s", url)
	state.FSM.GoTransition(fsm.LOADING_SERVER)
	server, serverErr := state.addInstituteServer(url)

	if serverErr != nil {
		return "", "", &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	return state.getConfig(server, forceTCP)
}

func (state *VPNState) GetConfigCustomServer(url string, forceTCP bool) (string, string, error) {
	errorMessage := fmt.Sprintf("failed getting a configuration for custom server %s", url)
	state.FSM.GoTransition(fsm.LOADING_SERVER)
	server, serverErr := state.addCustomServer(url)

	if serverErr != nil {
		return "", "", &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	return state.getConfig(server, forceTCP)
}

func (state *VPNState) CancelOAuth() error {
	errorMessage := "failed to cancel OAuth"
	if !state.FSM.InState(fsm.OAUTH_STARTED) {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: fsm.WrongStateError{Got: state.FSM.Current, Want: fsm.OAUTH_STARTED}.CustomError()}
	}

	currentServer, serverErr := state.Servers.GetCurrentServer()

	if serverErr != nil {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}
	server.CancelOAuth(currentServer)
	return nil
}

func (state *VPNState) ChangeSecureLocation() error {
	errorMessage := "failed to change location from the main screen"

	if !state.FSM.InState(fsm.NO_SERVER) {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: fsm.WrongStateError{Got: state.FSM.Current, Want: fsm.NO_SERVER}.CustomError()}
	}

	askLocationErr := state.askSecureLocation()

	if askLocationErr != nil {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: askLocationErr}
	}

	// Go back to the main screen
	state.FSM.GoTransitionWithData(fsm.NO_SERVER, state.GetSavedServers(), false)

	return nil
}

func (state *VPNState) GetDiscoOrganizations() (string, error) {
	if state.FSM.InState(fsm.DEREGISTERED) {
		return "", &types.WrappedErrorMessage{Message: "failed to get the organizations with Discovery", Err: fsm.DeregisteredError{}.CustomError()}
	}
	return state.Discovery.GetOrganizationsList()
}

func (state *VPNState) GetDiscoServers() (string, error) {
	if state.FSM.InState(fsm.DEREGISTERED) {
		return "", &types.WrappedErrorMessage{Message: "failed to get the servers with Discovery", Err: fsm.DeregisteredError{}.CustomError()}
	}
	return state.Discovery.GetServersList()
}

func (state *VPNState) SetProfileID(profileID string) error {
	errorMessage := "failed to set the profile ID for the current server"
	server, serverErr := state.Servers.GetCurrentServer()
	if serverErr != nil {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: serverErr}
	}

	base, baseErr := server.GetBase()

	if baseErr != nil {
		return &types.WrappedErrorMessage{Message: errorMessage, Err: baseErr}
	}
	base.Profiles.Current = profileID
	return nil
}

func (state *VPNState) SetSearchServer() error {
	if !state.FSM.HasTransition(fsm.SEARCH_SERVER) {
		return &types.WrappedErrorMessage{Message: "failed to set search server", Err: fsm.WrongStateTransitionError{Got: state.FSM.Current, Want: fsm.CONNECTED}.CustomError()}
	}

	state.FSM.GoTransition(fsm.SEARCH_SERVER)
	return nil
}

func (state *VPNState) getServerInfoData() string {
	jsonString, _ := state.Servers.GetCurrentServerInfoJSON()
	return jsonString
}

func (state *VPNState) SetConnected() error {
	if state.FSM.InState(fsm.CONNECTED) {
		// already connected, show no error
		return nil
	}
	if !state.FSM.HasTransition(fsm.CONNECTED) {
		return &types.WrappedErrorMessage{Message: "failed to set connected", Err: fsm.WrongStateTransitionError{Got: state.FSM.Current, Want: fsm.CONNECTED}.CustomError()}
	}

	state.FSM.GoTransitionWithData(fsm.CONNECTED, state.getServerInfoData(), false)
	return nil
}

func (state *VPNState) SetConnecting() error {
	if state.FSM.InState(fsm.CONNECTING) {
		// already loading connection, show no error
		return nil
	}
	if !state.FSM.HasTransition(fsm.CONNECTING) {
		return &types.WrappedErrorMessage{Message: "failed to set connecting", Err: fsm.WrongStateTransitionError{Got: state.FSM.Current, Want: fsm.CONNECTING}.CustomError()}
	}

	state.FSM.GoTransition(fsm.CONNECTING)
	return nil
}

func (state *VPNState) SetDisconnected() error {
	if state.FSM.InState(fsm.HAS_CONFIG) {
		// already disconnected, show no error
		return nil
	}
	if !state.FSM.HasTransition(fsm.HAS_CONFIG) {
		return &types.WrappedErrorMessage{Message: "failed to set disconnected", Err: fsm.WrongStateTransitionError{Got: state.FSM.Current, Want: fsm.HAS_CONFIG}.CustomError()}
	}

	state.FSM.GoTransitionWithData(fsm.HAS_CONFIG, state.getServerInfoData(), false)

	return nil
}

func (state *VPNState) RenewSession() error {
	errorMessage := "failed to renew session"

	currentServer, currentServerErr := state.Servers.GetCurrentServer()

	if currentServerErr != nil {
		return &types.WrappedErrorMessage{Message: "failed to renew session", Err: currentServerErr}
	}

	oauthStructure := currentServer.GetOAuth()
	oauthStructure.Token = oauth.OAuthToken{Access: "",Refresh: "",Type: "",Expires: 0,ExpiredTimestamp: util.GetCurrentTime()}
	// Make sure the FSM is initialized
	oauthStructure.FSM = &state.FSM

	loginErr := server.Login(currentServer)

	if loginErr != nil {
		// We are possibly in oauth started
		// Go back
		state.GoBack()
		return &types.WrappedErrorMessage{Message: errorMessage, Err: loginErr}
	}

	return nil
}


func (state *VPNState) ShouldRenewButton() bool {
	if !state.FSM.InState(fsm.CONNECTED) {
		return false
	}

	currentServer, currentServerErr := state.Servers.GetCurrentServer()

	if currentServerErr != nil {
		state.Logger.Log(log.LOG_INFO, fmt.Sprintf("No server found to renew with err: %s", GetErrorTraceback(currentServerErr)))
		return false
	}

	return server.ShouldRenewButton(currentServer)
}

func GetErrorCause(err error) error {
	return types.GetErrorCause(err)
}

func GetErrorLevel(err error) types.ErrorLevel {
	return types.GetErrorLevel(err)
}

func GetErrorTraceback(err error) string {
	return types.GetErrorTraceback(err)
}

func GetErrorJSONString(err error) string {
	return types.GetErrorJSONString(err)
}
