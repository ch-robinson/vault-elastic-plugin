package elastic

import (
	"testing"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	urlString := "http://testurl"
	var empty string

	cp := initializeDatabase(&empty, &urlString)

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Nil(t, err)
}

func TestInitializeFailOnMissingConnectionUrl(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	var empty string

	cp := initializeDatabase(&empty, &empty)

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Equal(t, "connection_url cannot be empty", err.Error())
}

func TestInitializeFailOnMissingUsername(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	urlString := "http://testurl"
	var empty string

	cp := initializeDatabase(&empty, &urlString)

	cp.Username = empty

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Equal(t, "username cannot be empty", err.Error())
}

func TestInitializeFailOnMissingPassword(t *testing.T) {
	ctx := testdata.NewMockVaultContext()

	urlString := "http://testurl"
	var empty string

	cp := initializeDatabase(&empty, &urlString)

	cp.RawConfig["password"] = empty

	err := cp.Initialize(ctx, cp.RawConfig, false)

	assert.Equal(t, "password cannot be empty", err.Error())
}
