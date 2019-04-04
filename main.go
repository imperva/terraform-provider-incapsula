package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-provider-incapsula/incapsula"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: incapsula.Provider})
}
