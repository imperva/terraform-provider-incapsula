package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAbpPendingChanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpPendingChangesRead,

		Description: "Produces a marker value used to trigger a new `incapsula_abp_preflight` " +
			"(and thereby a new publish) when the ABP configuration changes.\n\n" +
			"**How it actually works.** This data source performs no API call and does not " +
			"inspect the account: `hash` is a fixed value. Change detection comes entirely " +
			"from `depends_on`. When something this data source depends on has a planned " +
			"change, Terraform postpones reading it, `hash` becomes unknown for the duration " +
			"of the plan, and the preflight that consumes it is therefore replaced. When " +
			"nothing in `depends_on` changes, `hash` stays the same and no new preflight is " +
			"created.\n\n" +
			"Two consequences worth knowing. **Every ABP resource must be reachable through " +
			"`depends_on`** (for example by keeping them all in one module and depending on " +
			"that module): a change to a resource outside it does not trigger a preflight or " +
			"a publish, so it stays pending — `terraform plan` reports no changes while the " +
			"account still has unpublished changes — until something else publishes the " +
			"account, which then activates it along with everything else. And **changes made " +
			"outside Terraform are never noticed**, so they do not trigger a publish either.",

		Schema: map[string]*schema.Schema{
			"hash": {
				Description: "Change-detection marker. Its value is constant and carries no " +
					"information about the account; it is only unknown — and therefore a trigger " +
					"— while something this data source depends on has a planned change.",
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAbpPendingChangesRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	data.SetId("incapsula_abp_pending_changes")
	err := data.Set("hash", "incapsula_abp_pending_changes")
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
