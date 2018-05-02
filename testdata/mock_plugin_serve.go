package testdata

import "github.com/hashicorp/vault/api"

var MockPluginRan = false

var MockServePlugin = func(plugin interface{}, tlsConfig *api.TLSConfig) {
	MockPluginRan = true
}
