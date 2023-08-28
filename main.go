package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-incapsula/incapsula"
)

func main() {
	schema.DescriptionKind = schema.StringMarkdown

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: incapsula.Provider})
}
