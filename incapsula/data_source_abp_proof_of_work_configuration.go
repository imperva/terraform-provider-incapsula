package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpProofOfWorkConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpProofOfWorkConfigurationRead,

		Description: "Looks up an ABP Proof Of Work Configuration in an account by name.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to search within.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Name of the Proof Of Work Configuration to look up.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"difficulty": {
				Description: "Number of credits a client gets for completing a proof-of-work challenge.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"algorithm": {
				Description: "Proof Of Work algorithm.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Proof Of Work Configuration was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Proof Of Work Configuration was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAbpProofOfWorkConfigurationRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	name := data.Get("name").(string)

	configs, err := client.ListAbpProofOfWorkConfigurations(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	var match *AbpProofOfWorkConfiguration
	for i := range configs {
		if configs[i].Name == name {
			if match != nil {
				return diag.Errorf("multiple ABP Proof Of Work Configurations named %q found in account %s", name, accountId)
			}
			match = &configs[i]
		}
	}
	if match == nil {
		return diag.Errorf("no ABP Proof Of Work Configuration named %q found in account %s", name, accountId)
	}

	data.SetId(match.Id)
	if err := serializeAbpProofOfWorkConfiguration(data, match); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
