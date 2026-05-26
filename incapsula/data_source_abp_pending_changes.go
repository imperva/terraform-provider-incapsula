package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAbpPendingChanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpPendingChangesRead,

		Description: "Exposes a hash to facilitate change detection caused by e.g. a depends_on = [module.x]\n",

		Schema: map[string]*schema.Schema{
			"hash": {
				Description: "Stable (but computed) marker facilitating change-detection if any resource that this data depends on changes",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAbpPendingChangesRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	data.SetId("incapsula_abp_pending_changes")
	err := data.Set("hash", "incapsula_abp_pending_changes")
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
