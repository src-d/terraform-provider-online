package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/src-d/terraform-provider-online-net/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
