package elastic

import (
	"context"
	"fmt"
	"sync"

	"github.com/ch-robinson/vault-elastic-plugin/httputil"
	"github.com/mitchellh/mapstructure"
)

// connectionProducer implements ConnectionProducer and provides an
// interface for databases to make connections.
type connectionProducer struct {
	ConnectionURL string `json:"connection_url" structs:"connection_url" mapstructure:"connection_url"`
	Username      string `json:"username" structs:"username" mapstructure:"username"`
	Password      string `json:"password" structs:"password" mapstructure:"password"`
	RawConfig     map[string]interface{}
	Type          string
	HTTPClient    httputil.HTTPClienter
	sync.RWMutex
}

// Initialize is deprecated, this is just needed for interface implemenation
func (c *connectionProducer) Initialize(ctx context.Context, conf map[string]interface{}, verifyConnection bool) error {
	_, err := c.Init(ctx, conf, verifyConnection)
	return err
}

// Init parses connection configuration.
func (c *connectionProducer) Init(ctx context.Context, conf map[string]interface{}, verifyConnection bool) (map[string]interface{}, error) {
	c.Lock()
	defer c.Unlock()

	c.RawConfig = conf

	err := mapstructure.WeakDecode(conf, c)
	if err != nil {
		return nil, err
	}

	if len(c.ConnectionURL) == 0 {
		return nil, fmt.Errorf("connection_url cannot be empty")
	}

	if len(c.Username) == 0 {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if len(c.Password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	return conf, nil
}

// Close terminates the database connection.
func (c *connectionProducer) Close() error {
	// NOOP since we are making rest calls instead of
	// direct db connection
	return nil
}

// SecretValues masks the password in case of exeptions so they are not returned in outputs or logs
func (c *connectionProducer) SecretValues() map[string]interface{} {
	return map[string]interface{}{
		c.Password: "[password]",
	}
}
