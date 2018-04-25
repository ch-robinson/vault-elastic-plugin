package util

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/stretchr/testify/assert"
)

func initHTTPClient() *HTTPClient {
	mockHTTP := testdata.NewMockHTTP(nil)
	httpClient := NewHTTPClient(mockHTTP)
	return httpClient
}

func TestNewHTTPClient(t *testing.T) {
	mockHTTP := testdata.NewMockHTTP(nil)

	httpClient := NewHTTPClient(mockHTTP)

	assert.Equal(t, mockHTTP, httpClient.client)
}

func TestReadHTTPResponse(t *testing.T) {
	httpClient := initHTTPClient()

	readCloser := NewClosingBuffer(bytes.NewBufferString("{\"test\":\"body\"}")).GetReadCloser()

	res := &http.Response{
		Status:     "ok",
		StatusCode: 200,
		Body:       readCloser,
	}

	response, err := httpClient.ReadHTTPResponse(res)

	assert.Nil(t, err)
	assert.True(t, response != nil)
}
