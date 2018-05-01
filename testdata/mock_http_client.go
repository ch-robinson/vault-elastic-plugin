package testdata

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ch-robinson/vault-elastic-plugin/plugin/interfaces"
)

// MockHTTPClient mocks HTTPClient
type MockHTTPClient struct {
	responseBody *string
	client       interfaces.IHTTP
}

// NewMockHTTPClient instantiates a new mock http client
func NewMockHTTPClient(responseBody *string, client interfaces.IHTTP) interfaces.IHTTPClient {
	return &MockHTTPClient{responseBody, client}
}

// Get mocks Get
func (m *MockHTTPClient) Get(url string, accessToken, authType *string) (*http.Response, error) {
	if url == "bad" {
		return nil, errors.New("get test error")
	}
	fmt.Printf("\n\n%s\n\n", *m.responseBody)
	readCloser := newClosingBuffer(bytes.NewBufferString(*m.responseBody)).GetReadCloser()

	return &http.Response{
		Status:     "success",
		StatusCode: 200,
		Body:       readCloser,
	}, nil
}

// BuildPostBasicAuthRequest mocks building a post request with basic auth
func (m *MockHTTPClient) BuildBasicAuthRequest(requestURL, username, password, httpMethod string, body map[string]interface{}) (*http.Request, error) {
	if requestURL == "bad" {
		return nil, nil
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	err := enc.Encode(body)

	if err != nil {
		panic(err)
	}

	readCloser := newClosingBuffer(bytes.NewBufferString(*m.responseBody)).GetReadCloser()

	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "", readCloser)

	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

// Do mocks a request
func (m *MockHTTPClient) Do(req *http.Request) (resp *http.Response, err error) {
	if req == nil {
		return nil, errors.New("http post test error")
	}

	readCloser := newClosingBuffer(bytes.NewBufferString(*m.responseBody)).GetReadCloser()

	return &http.Response{
		Status:     "success",
		StatusCode: 200,
		Body:       readCloser,
	}, nil
}

// ReadHTTPResponse mocks ReadHTTPResponse
func (m *MockHTTPClient) ReadHTTPResponse(res *http.Response) (map[string]interface{}, error) {
	var mockRes map[string]interface{}

	err := json.Unmarshal([]byte(*m.responseBody), &mockRes)

	if err != nil {
		panic(err)
	}

	return mockRes, nil
}
