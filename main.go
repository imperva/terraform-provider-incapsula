package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-incapsula/incapsula"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: incapsula.Provider})
}
