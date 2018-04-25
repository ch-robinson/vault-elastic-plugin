package testdata

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/ch-robinson/vault-elastic-plugin/plugin/interfaces"
)

// MockHTTP Mock object for http methods
type MockHTTP struct {
	responseBody *string
}

// NewMockHTTP instantiates a new MockHTTP
func NewMockHTTP(responseBody *string) interfaces.IHTTP {
	return &MockHTTP{responseBody}
}

// Get Mock Get request
func (m *MockHTTP) Get(url string) (resp *http.Response, err error) {
	if url == "bad" {
		return nil, errors.New("http get test error")
	}

	readCloser := newClosingBuffer(bytes.NewBufferString(*m.responseBody)).GetReadCloser()

	return &http.Response{
		Status:     "success",
		StatusCode: 200,
		Body:       readCloser,
	}, nil
}

// Do Mock Do request
func (m *MockHTTP) Do(req *http.Request) (*http.Response, error) {
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
