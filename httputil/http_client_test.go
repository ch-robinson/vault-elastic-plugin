package httputil

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/stretchr/testify/assert"
)

func initHTTPClient(resBody *string) *HTTPClient {
	mockHTTP := testdata.NewMockHTTP(resBody)
	httpClient := New(mockHTTP)
	return httpClient
}

func TestNewHTTPClient(t *testing.T) {
	mockHTTP := testdata.NewMockHTTP(nil)

	httpClient := New(mockHTTP)

	assert.Equal(t, mockHTTP, httpClient.client)
}

func TestReadHTTPResponse(t *testing.T) {
	httpClient := initHTTPClient(nil)

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

func TestAddHeaders(t *testing.T) {
	httpClient := initHTTPClient(nil)
	var h http.Header = make(map[string][]string)

	httpClient.addHeaders(&h)

	assert.Equal(t, "application/json", h.Get("Content-Type"))
	assert.Equal(t, "application/json", h.Get("Accept"))
}

func TestDo(t *testing.T) {
	b := "good"
	httpClient := initHTTPClient(&b)

	req, _ := http.NewRequest("GET", "mocked", nil)

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
}

func TestBuildBasicAuthRequestWithBody(t *testing.T) {
	httpClient := initHTTPClient(nil)
	body := make(map[string]interface{})
	body["test"] = true

	req, err := httpClient.BuildBasicAuthRequest("http://test", "testuser", "testpassword", "POST", body)

	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk", req.Header.Get("Authorization"))
}

func TestBuildBasicAuthRequestNilBody(t *testing.T) {
	httpClient := initHTTPClient(nil)

	req, err := httpClient.BuildBasicAuthRequest("http://test", "testuser", "testpassword", "DELETE", nil)

	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk", req.Header.Get("Authorization"))
}
