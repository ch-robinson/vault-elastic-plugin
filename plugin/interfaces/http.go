package interfaces

import "net/http"

// IHTTP is the interface for http.Client
// https://golang.org/pkg/net/http/
type IHTTP interface {
	// Do wraps the http.Client Do function
	Do(req *http.Request) (*http.Response, error)
}

// IHTTPClient is the interface for functions relating to building http request
type IHTTPClient interface {
	// Do peforms an http request. This is just a wrapper for the http.Client function
	// calls IHTTP.Do for ease of testing
	Do(req *http.Request) (*http.Response, error)
	// BuildPostBasicAuthRequest creates an http.Request with basic authoriztion header.
	// body must be map[string]interface{}
	BuildPostBasicAuthRequest(requestURL, username, password *string, body map[string]interface{}) *http.Request
	// ReadHTTPResponse returns the response body as map[string]interface{}
	ReadHTTPResponse(res *http.Response) (map[string]interface{}, error)
}
