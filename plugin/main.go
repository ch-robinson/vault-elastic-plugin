package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ch-robinson/vault-elastic-plugin/plugin/elastic"
	"github.com/ch-robinson/vault-elastic-plugin/plugin/util"
	"github.com/hashicorp/vault/helper/pluginutil"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	clientWrapper := util.NewHTTPClient(&http.Client{})

	err := elastic.Run(apiClientMeta.GetTLSConfig(), clientWrapper)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
