package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/ch-robinson/vault-elastic-plugin/plugin/interfaces"
)

// HTTPClient is the wrapper for interacting with http methods
type HTTPClient struct {
	client interfaces.IHTTP
}

// NewHTTPClient instantiates a new HttpClient
func NewHTTPClient(client interfaces.IHTTP) *HTTPClient {
	return &HTTPClient{client}
}

// Do peforms an http request
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

// BuildBasicAuthRequest creates an http.Request with basic authoriztion header.
// body must be map[string]interface{}
func (c *HTTPClient) BuildBasicAuthRequest(requestURL, username, password, httpMethod string, body map[string]interface{}) *http.Request {
	recloser := &ClosingBuffer{}

	if len(body) > 0 {
		reqBody, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		recloser = NewClosingBuffer(bytes.NewBuffer(reqBody)).GetReadCloser().(*ClosingBuffer)

		if err != nil {
			panic(err)
		}
	}

	req, err := http.NewRequest(httpMethod, requestURL, recloser)

	req.SetBasicAuth(username, password)

	if err != nil {
		panic(err)
	}

	c.addHeaders(&req.Header)

	return req
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
