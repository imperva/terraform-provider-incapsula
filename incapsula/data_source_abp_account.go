package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAbpAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpAccountRead,

		Description: "Looks up the ABP account belonging to the configured API credentials. " +
			"`account_id` and the data source's own `id` are both the account UUID, so either " +
			"can be referenced.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "ABP account UUID. Same value as `id`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAbpAccountRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)

	accountId, err := client.ReadAbpAccountId()
	if err != nil {
		return diag.FromErr(err)
	}

	data.SetId(accountId)
	if err := data.Set("account_id", accountId); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
