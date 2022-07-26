package incapsula

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Description: "Provides data about user Account",

		Schema: map[string]*schema.Schema{
			// Computed Attributes
			"current_account": {
				Type:        schema.TypeString,
				Description: "Current account ID",
				Computed:    true,
			},
			"plan_name": {
				Type:        schema.TypeString,
				Description: "Plan name",
				Computed:    true,
			},
		},
	}
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountStatusResponse, err := client.Verify()
	if err != nil {
		return diag.Errorf("Error checking account details: %v", err)
	}
	d.SetId(strconv.Itoa(accountStatusResponse.AccountID))
	d.Set("current_account", strconv.Itoa(accountStatusResponse.AccountID))
	d.Set("plan_name", accountStatusResponse.Account.PlanName)

	return nil
}
