package elastic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ch-robinson/vault-elastic-plugin/plugin/interfaces"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/logical/database/dbplugin"
	"github.com/hashicorp/vault/plugins/helper/database/credsutil"
)

// ElasticTypeName is the name of the plugin type
const ElasticTypeName = "elasticdb"

// Database is an implementation of github.com/hashicorp/vault/builtin/logical/database/dbplugin Database interface
type Database struct {
	*connectionProducer
	credsutil.CredentialsProducer
}

var _ dbplugin.Database = &Database{}

// New returns a new Elastic instance with provided implementation of http.Client
func New(httpClient interfaces.IHTTPClient) (interface{}, error) {
	// setup struct
	db := &Database{
		connectionProducer: &connectionProducer{
			Type:       ElasticTypeName,
			HTTPClient: httpClient,
		},
		// we can still use this struct despite the name
		CredentialsProducer: &credsutil.SQLCredentialsProducer{
			DisplayNameLen: 15,
			RoleNameLen:    15,
			UsernameLen:    30,
			Separator:      "-",
		},
	}

	// This just set's struct fields so plugin implementations are used
	dbType := dbplugin.NewDatabaseErrorSanitizerMiddleware(db, db.SecretValues)

	return dbType, nil
}

// Run instantiates the Database struct, and runs the RPC server for the plugin
func Run(serve func(plugin interface{}, tlsConfig *api.TLSConfig), apiTLSConfig *api.TLSConfig, httpClient interfaces.IHTTPClient) error {
	dbType, err := New(httpClient)

	if err != nil {
		return err
	}

	serve(dbType, apiTLSConfig)

	return nil
}

// Type returns the TypeName for this backend
func (m *Database) Type() (string, error) {
	return ElasticTypeName, nil
}

// CreateUser creates a new user in Elastic DB
func (m *Database) CreateUser(ctx context.Context, statements dbplugin.Statements, usernameConfig dbplugin.UsernameConfig, expiration time.Time) (string, string, error) {
	// Generates the new password
	newPassword, err := m.GeneratePassword()

	if err != nil {
		return "", "", err
	}

	// Generates the new username
	newUsername, err := m.GenerateUsername(usernameConfig)

	if err != nil {
		return "", "", err
	}

	var body = make(map[string]interface{})
	body["password"] = newPassword
	body["roles"] = statements.Creation

	if len(statements.Creation) == 0 {
		return "", "", fmt.Errorf("roles array is required when creating a user")
	}

	var url = fmt.Sprintf("%s/_xpack/security/user/%s", m.ConnectionURL, newUsername)

	request, err := m.HTTPClient.BuildBasicAuthRequest(url, m.Username, m.Password, "POST", body)

	if err != nil {
		return "", "", err
	}

	res, err := m.HTTPClient.Do(request)

	if err != nil {
		return "", "", err
	}

	response, err := m.HTTPClient.ReadHTTPResponse(res)

	if err != nil {
		return "", "", err
	}

	// elastic doesn't throw an exeception if the user was not created.
	// instead they return {"user":{"created":true}} so we need to check
	// and return an error if the user was not created.
	user := response["user"].(map[string]interface{})
	created := user["created"].(bool)

	if !created {
		return "", "", fmt.Errorf("User was not created: %+v", response)
	}

	return newUsername, newPassword, nil
}

// RenewUser is not currently used
func (m *Database) RenewUser(ctx context.Context, statements dbplugin.Statements, username string, expiration time.Time) error {
	// NOOP
	return nil
}

// RevokeUser drops the specified user from the authentication database.
func (m *Database) RevokeUser(ctx context.Context, statements dbplugin.Statements, username string) error {
	var url = fmt.Sprintf("%s/_xpack/security/user/%s", m.ConnectionURL, username)

	request, err := m.HTTPClient.BuildBasicAuthRequest(url, m.Username, m.Password, "DELETE", nil)

	if err != nil {
		return err
	}

	res, err := m.HTTPClient.Do(request)

	if err != nil {
		return err
	}

	_, err = m.HTTPClient.ReadHTTPResponse(res)

	if err != nil {
		return err
	}

	return nil
}

// RotateRootCredentials rotates the root superuser credentials stored for the database connection
func (m *Database) RotateRootCredentials(ctx context.Context, statements []string) (map[string]interface{}, error) {
	if len(m.Username) == 0 || len(m.Password) == 0 {
		return nil, errors.New("Both the username and password are required.")
	}

	password, err := m.GeneratePassword()

	if err != nil {
		return nil, err
	}

	var body = make(map[string]interface{})
	body["password"] = password

	url := fmt.Sprintf("%s/_xpack/security/user/%s/_password", m.ConnectionURL, m.Username)

	request, err := m.HTTPClient.BuildBasicAuthRequest(url, m.Username, m.Password, "PUT", body)

	if err != nil {
		return nil, err
	}

	res, err := m.HTTPClient.Do(request)

	if err == nil || res.StatusCode == 200 {
		// Need to return the password back to Vault so it knows what the new password is
		m.RawConfig["password"] = password
		return m.RawConfig, nil
	}

	return m.RawConfig, err
}
