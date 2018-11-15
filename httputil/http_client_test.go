package httputil

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/stretchr/testify/assert"
)

func initClient(resBody *string) *ClientWrapper {
	mockHTTP := testdata.NewMockHTTP(resBody)
	clientWrapper := New(mockHTTP)
	return clientWrapper
}

func TestNewClient(t *testing.T) {
	mockHTTP := testdata.NewMockHTTP(nil)

	clientWrapper := New(mockHTTP)

	assert.Equal(t, mockHTTP, clientWrapper.client)
}

func TestReadHTTPResponse(t *testing.T) {
	clientWrapper := initClient(nil)

	readCloser := NewClosingBuffer(bytes.NewBufferString("{\"test\":\"body\"}")).GetReadCloser()

	res := &http.Response{
		Status:     "ok",
		StatusCode: 200,
		Body:       readCloser,
	}

	response, err := clientWrapper.ReadHTTPResponse(res)

	assert.Nil(t, err)
	assert.True(t, response != nil)
}

func TestAddHeaders(t *testing.T) {
	clientWrapper := initClient(nil)
	var h http.Header = make(map[string][]string)

	clientWrapper.addHeaders(&h)

	assert.Equal(t, "application/json", h.Get("Content-Type"))
	assert.Equal(t, "application/json", h.Get("Accept"))
}

func TestDo(t *testing.T) {
	b := "good"
	clientWrapper := initClient(&b)

	req, _ := http.NewRequest("GET", "mocked", nil)

	_, err := clientWrapper.Do(req)

	assert.Nil(t, err)
}

func TestBuildBasicAuthRequestWithBody(t *testing.T) {
	clientWrapper := initClient(nil)
	body := make(map[string]interface{})
	body["test"] = true

	req, err := clientWrapper.BuildBasicAuthRequest("http://test", "testuser", "testpassword", "POST", body)

	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk", req.Header.Get("Authorization"))
}

func TestBuildBasicAuthRequestNilBody(t *testing.T) {
	clientWrapper := initClient(nil)

	req, err := clientWrapper.BuildBasicAuthRequest("http://test", "testuser", "testpassword", "DELETE", nil)

	assert.Nil(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "Basic dGVzdHVzZXI6dGVzdHBhc3N3b3Jk", req.Header.Get("Authorization"))
}
