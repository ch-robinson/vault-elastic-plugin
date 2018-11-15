package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ch-robinson/vault-elastic-plugin/elastic"
	"github.com/ch-robinson/vault-elastic-plugin/httputil"
	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/plugins"
)

var Version = "1.0"

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	clientWrapper := httputil.New(&http.Client{})

	if err := elastic.Run(plugins.Serve, apiClientMeta.GetTLSConfig(), clientWrapper); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
