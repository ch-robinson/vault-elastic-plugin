package elastic

import (
	"testing"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/stretchr/testify/assert"
)

func initializeConnectionProducer(response, connURL *string) *connectionProducer {
	mockHTTP := testdata.NewMockHTTP(response)
	httpClient := testdata.NewMockHTTPClient(response, mockHTTP)

	url := "https://mock"

	if connURL != nil {
		url = *connURL
	}

	mockRawconfig := make(map[string]interface{})
	mockRawconfig["password"] = "testpassword"

	return &connectionProducer{
		Type:          ElasticTypeName,
		HTTPClient:    httpClient,
		ConnectionURL: url,
		Username:      "testuser",
		Password:      "testpassword",
		RawConfig:     mockRawconfig,
	}
}

func TestInitialize(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	urlString := "http://testurl"
	var empty string

	cp := initializeConnectionProducer(&empty, &urlString)

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Nil(t, err)
}

func TestInitializeFailOnMissingConnectionUrl(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	var empty string

	cp := initializeConnectionProducer(&empty, &empty)

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Equal(t, "connection_url cannot be empty", err.Error())
}

func TestInitializeFailOnMissingUsername(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	urlString := "http://testurl"
	var empty string

	cp := initializeConnectionProducer(&empty, &urlString)

	cp.Username = empty

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Equal(t, "username cannot be empty", err.Error())
}

func TestInitializeFailOnMissingPassword(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	urlString := "http://testurl"
	var empty string

	cp := initializeConnectionProducer(&empty, &urlString)

	cp.RawConfig["password"] = empty

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Equal(t, "password cannot be empty", err.Error())
}

func TestClose(t *testing.T) {
	urlString := "http://testurl"
	var empty string

	cp := initializeConnectionProducer(&empty, &urlString)

	err := cp.Close()

	assert.Nil(t, err)
}
