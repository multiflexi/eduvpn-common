package http

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type URLParameters map[string]string

type HTTPOptionalParams struct {
	Headers       http.Header
	URLParameters URLParameters
	Body          url.Values
}

// Construct an URL including on parameters
func HTTPConstructURL(baseURL string, parameters URLParameters) (string, error) {
	url, parseErr := url.Parse(baseURL)
	if parseErr != nil {
		return "", &HTTPConstructURLError{URL: baseURL, Parameters: parameters, Err: parseErr}
	}

	q := url.Query()

	for parameter, value := range parameters {
		q.Set(parameter, value)
	}
	url.RawQuery = q.Encode()
	return url.String(), nil
}

// Convenience functions
func HTTPGet(url string) (http.Header, []byte, error) {
	return HTTPMethodWithOpts(http.MethodGet, url, nil)
}

func HTTPPost(url string, body url.Values) (http.Header, []byte, error) {
	return HTTPMethodWithOpts(http.MethodGet, url, &HTTPOptionalParams{Body: body})
}

func HTTPGetWithOpts(url string, opts *HTTPOptionalParams) (http.Header, []byte, error) {
	return HTTPMethodWithOpts(http.MethodGet, url, opts)
}

func HTTPPostWithOpts(url string, opts *HTTPOptionalParams) (http.Header, []byte, error) {
	return HTTPMethodWithOpts(http.MethodPost, url, opts)
}

func httpOptionalURL(url string, opts *HTTPOptionalParams) (string, error) {
	if opts != nil {
		url, urlErr := HTTPConstructURL(url, opts.URLParameters)

		if urlErr != nil {
			return url, &HTTPRequestCreateError{URL: url, Err: urlErr}
		}
		return url, nil
	}
	return url, nil
}

func httpOptionalHeaders(req *http.Request, opts *HTTPOptionalParams) {
	// Add headers
	if opts != nil && req != nil {
		for k, v := range opts.Headers {
			req.Header.Add(k, v[0])
		}
	}
}

func httpOptionalBodyReader(opts *HTTPOptionalParams) io.Reader {
	if opts != nil && opts.Body != nil {
		return strings.NewReader(opts.Body.Encode())
	}
	return nil
}

func HTTPMethodWithOpts(method string, url string, opts *HTTPOptionalParams) (http.Header, []byte, error) {
	// Make sure the url contains all the parameters
	// This can return an error,
	// it already has the right error so so we don't wrap it further
	url, urlErr := httpOptionalURL(url, opts)
	if urlErr != nil {
		// No further type wrapping is needed here
		return nil, nil, urlErr
	}

	// Create a client
	client := &http.Client{}

	// Create request object with the body reader generated from the optional arguments
	req, reqErr := http.NewRequest(method, url, httpOptionalBodyReader(opts))
	if reqErr != nil {
		return nil, nil, &HTTPRequestCreateError{URL: url, Err: reqErr}
	}

	// See https://stackoverflow.com/questions/17714494/golang-http-request-results-in-eof-errors-when-making-multiple-requests-successi
	req.Close = true

	// Make sure the headers contain all the parameters
	httpOptionalHeaders(req, opts)

	// Do request
	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, nil, &HTTPResourceError{URL: url, Err: respErr}
	}

	// Request successful, make sure body is closed at the end
	defer resp.Body.Close()

	// Return a string
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return resp.Header, nil, &HTTPReadError{URL: url, Err: readErr}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp.Header, body, &HTTPStatusError{URL: url, Status: resp.StatusCode}
	}

	// Return the body in bytes and signal the status error if there was one
	return resp.Header, body, nil
}

type HTTPResourceError struct {
	URL string
	Err error
}

func (e *HTTPResourceError) Error() string {
	return fmt.Sprintf("failed obtaining HTTP resource: %s with error: %v", e.URL, e.Err)
}

type HTTPStatusError struct {
	URL    string
	Status int
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("failed obtaining HTTP resource: %s as it gave an unsuccesful status code: %d", e.URL, e.Status)
}

type HTTPReadError struct {
	URL string
	Err error
}

func (e *HTTPReadError) Error() string {
	return fmt.Sprintf("failed reading HTTP resource: %s with error: %v", e.URL, e.Err)
}

type HTTPParseJsonError struct {
	URL  string
	Body string
	Err  error
}

func (e *HTTPParseJsonError) Error() string {
	return fmt.Sprintf("failed parsing json %s for HTTP resource: %s with error: %v", e.Body, e.URL, e.Err)
}

type HTTPRequestCreateError struct {
	URL string
	Err error
}

func (e *HTTPRequestCreateError) Error() string {
	return fmt.Sprintf("failed to create HTTP request with url: %s and error: %v", e.URL, e.Err)
}

type HTTPConstructURLError struct {
	URL        string
	Parameters URLParameters
	Err        error
}

func (e *HTTPConstructURLError) Error() string {
	return fmt.Sprintf("failed to construct url: %s including parameters: %v with error: %v", e.URL, e.Parameters, e.Err)
}