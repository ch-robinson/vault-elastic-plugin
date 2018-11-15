package elastic

import (
	"testing"
	"time"

	"github.com/ch-robinson/vault-elastic-plugin/testdata"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
	"github.com/stretchr/testify/assert"
)

func initializeDatabase(response, connURL *string) *Database {
	mockHTTP := testdata.NewMockHTTP(response)
	httpClient := testdata.NewMockHTTPClient(response, mockHTTP)

	url := "https://mock"

	if connURL != nil {
		url = *connURL
	}

	mockRawconfig := make(map[string]interface{})
	mockRawconfig["password"] = "testpassword"

	return &Database{
		connectionProducer: &connectionProducer{
			Type:          ElasticTypeName,
			HTTPClient:    httpClient,
			ConnectionURL: url,
			Username:      "testuser",
			Password:      "testpassword",
			RawConfig:     mockRawconfig,
		},
		// we can still use this struct despite the name
		CredentialsProducer: &credsutil.SQLCredentialsProducer{
			DisplayNameLen: 15,
			RoleNameLen:    15,
			UsernameLen:    30,
			Separator:      "-",
		},
	}
}

func TestNew(t *testing.T) {
	mockHTTP := testdata.NewMockHTTP(nil)
	httpClient := testdata.NewMockHTTPClient(nil, mockHTTP)

	dbType, err := New(httpClient)

	assert.Nil(t, err)

	dbMiddleware := dbType.(*dbplugin.DatabaseErrorSanitizerMiddleware)

	name, err := dbMiddleware.Type()

	assert.Nil(t, err)
	assert.Equal(t, name, "elasticdb")
}

func TestRun(t *testing.T) {
	mockHTTP := testdata.NewMockHTTP(nil)
	httpClient := testdata.NewMockHTTPClient(nil, mockHTTP)
	err := Run(testdata.MockServePlugin, nil, httpClient)

	assert.Nil(t, err)
	assert.True(t, testdata.MockPluginRan)
}

func TestType(t *testing.T) {
	db := initializeDatabase(nil, nil)

	name, err := db.Type()

	assert.Nil(t, err)
	assert.Equal(t, "elasticdb", name)
}

func TestCreateUserSuccess(t *testing.T) {
	res := `{"user":{"created":true}}`
	db := initializeDatabase(&res, nil)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{
		Creation: []string{"role1"},
	}

	usernameConfig := dbplugin.UsernameConfig{
		RoleName:    "test",
		DisplayName: "test role",
	}

	exp := time.Now().Add(time.Second)

	username, password, err := db.CreateUser(ctx, statements, usernameConfig, exp)

	assert.Nil(t, err)
	assert.True(t, len(username) > 0)
	assert.True(t, len(password) > 0)
}

func TestCreateUserFailCreateUser(t *testing.T) {
	res := `{"user":{"created":false}}`
	db := initializeDatabase(&res, nil)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{
		Creation: []string{"role1"},
	}

	usernameConfig := dbplugin.UsernameConfig{
		RoleName:    "test",
		DisplayName: "test role",
	}

	exp := time.Now().Add(time.Second)

	_, _, err := db.CreateUser(ctx, statements, usernameConfig, exp)

	assert.Equal(t, "User was not created: map[user:map[created:false]]", err.Error())
}

func TestCreateUserFailMissingRoles(t *testing.T) {
	res := `{"user":{"created":false}}`

	db := initializeDatabase(&res, nil)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{
		Creation: []string{},
	}

	usernameConfig := dbplugin.UsernameConfig{
		RoleName:    "test",
		DisplayName: "test role",
	}

	exp := time.Now().Add(time.Second)

	_, _, err := db.CreateUser(ctx, statements, usernameConfig, exp)

	assert.Equal(t, "roles array is required when creating a user", err.Error())
}

func TestCreateUserFailBadRequestURL(t *testing.T) {
	res := `{"user":{"created":false}}`
	url := "bad"

	db := initializeDatabase(&res, &url)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{
		Creation: []string{"role1"},
	}

	usernameConfig := dbplugin.UsernameConfig{
		RoleName:    "test",
		DisplayName: "test role",
	}

	exp := time.Now().Add(time.Second)

	_, _, err := db.CreateUser(ctx, statements, usernameConfig, exp)

	assert.Equal(t, "bad request url", err.Error())
}

func TestCreateUserFailHTTPPost(t *testing.T) {
	res := `{"user":{"created":false}}`
	url := "failedbutcontinue"

	db := initializeDatabase(&res, &url)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{
		Creation: []string{"role1"},
	}

	usernameConfig := dbplugin.UsernameConfig{
		RoleName:    "test",
		DisplayName: "test role",
	}

	exp := time.Now().Add(time.Second)

	_, _, err := db.CreateUser(ctx, statements, usernameConfig, exp)

	assert.Equal(t, "http post test error", err.Error())
}

func TestRenewUserNotImplemented(t *testing.T) {
	db := initializeDatabase(nil, nil)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{
		Creation: []string{"role1"},
	}

	exp := time.Now().Add(time.Second)

	err := db.RenewUser(ctx, statements, "test", exp)

	assert.Nil(t, err)
}

func TestRevokeUserFail(t *testing.T) {
	res := `{}`

	db := initializeDatabase(&res, nil)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{}

	err := db.RevokeUser(ctx, statements, "nouser")

	assert.Equal(t, "user doesn't exist", err.Error())
}

func TestRevokeUserSuccess(t *testing.T) {
	res := `{}`

	db := initializeDatabase(&res, nil)

	ctx := testdata.NewMockVaultContext()

	statements := dbplugin.Statements{}

	err := db.RevokeUser(ctx, statements, "MySuccessTestUser")

	assert.Nil(t, err)
}

func TestRotateRootCredentialsSuccess(t *testing.T) {
	res := `{}`

	db := initializeDatabase(&res, nil)

	response, err := db.RotateRootCredentials(testdata.NewMockVaultContext(), []string{})

	assert.Nil(t, err)
	assert.Equal(t, db.RawConfig["password"], response["password"])
}

func TestRotateRootCredentialsFailWithoutConnectionCredentials(t *testing.T) {
	res := `{}`

	db := initializeDatabase(&res, nil)
	db.Username = ""

	_, err := db.RotateRootCredentials(testdata.NewMockVaultContext(), []string{})

	assert.Equal(t, "both the username and password are required", err.Error())
}

func TestRotateRootCredentialsFailOnBadUsername(t *testing.T) {
	res := `{}`

	db := initializeDatabase(&res, nil)
	db.Username = "nouser"

	_, err := db.RotateRootCredentials(testdata.NewMockVaultContext(), []string{})

	assert.Equal(t, "user doesn't exist", err.Error())
}
