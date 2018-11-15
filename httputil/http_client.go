package httputil

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// HTTP is the interface for http.Client
// https://golang.org/pkg/net/http/
type HTTPer interface {
	// Do wraps the http.Client Do function
	Do(req *http.Request) (*http.Response, error)
}

// HTTPClient is the interface for functions relating to building http request
type HTTPClienter interface {
	// Do peforms an http request. This is just a wrapper for the http.Client function
	// calls HTTP.Do for ease of testing
	Do(req *http.Request) (*http.Response, error)
	// BuildBasicAuthRequest creates an http.Request with basic authoriztion header.
	// body must be map[string]interface{}
	BuildBasicAuthRequest(requestURL, username, password, httpMethod string, body map[string]interface{}) (*http.Request, error)
	// ReadHTTPResponse returns the response body as map[string]interface{}
	ReadHTTPResponse(res *http.Response) (map[string]interface{}, error)
}

// HTTPClient is the wrapper for interacting with http methods
type HTTPClient struct {
	client HTTPer
}

// NewHTTPClient instantiates a new HttpClient
func New(client HTTPer) *HTTPClient {
	return &HTTPClient{client}
}

// Do peforms an http request
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// BuildBasicAuthRequest creates an http.Request with basic authoriztion header.
// body must be map[string]interface{}
func (c *HTTPClient) BuildBasicAuthRequest(requestURL, username, password, httpMethod string, body map[string]interface{}) (*http.Request, error) {
	var req *http.Request
	var err error

	if body != nil && len(body) > 0 {
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		recloser := NewClosingBuffer(bytes.NewBuffer(reqBody)).GetReadCloser().(*ClosingBuffer)

		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(httpMethod, requestURL, recloser)
	} else {
		req, err = http.NewRequest(httpMethod, requestURL, nil)
	}

	req.SetBasicAuth(username, password)

	if err != nil {
		return nil, err
	}

	c.addHeaders(&req.Header)

	return req, nil
}

// ReadHTTPResponse returns the response body as map[string]interface{}
func (c *HTTPClient) ReadHTTPResponse(res *http.Response) (map[string]interface{}, error) {
	resBody, err := ioutil.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	var body map[string]interface{}

	err = json.Unmarshal(resBody, &body)

	if err != nil {
		return nil, err
	}

	// Throw error if not ok. Might need to watch out for other success codes, but this should be ok.
	if res.StatusCode != 200 {
		return nil, err
	}

	return body, nil
}

// addHeaders adds http.Headers. If accessToken is provided, an Authorization header will be added with given authType (Bearer, token, etc.)
func (c *HTTPClient) addHeaders(header *http.Header) {
	header.Add("Content-Type", "application/json")
	header.Add("Accept", "application/json")
}
