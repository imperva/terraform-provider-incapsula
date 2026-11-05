package incapsula

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpPolicies() *schema.Resource {
	policyAttributes := map[string]*schema.Schema{
		"id": {
			Description: "ID of the Policy.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Name of the Policy.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
	maps.Copy(policyAttributes, abpPolicyReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpPoliciesRead,

		Description: "Lists the ABP Policies of an account, ordered by name. Use " +
			"`incapsula_abp_policy` to look up a single Policy by name or id instead.\n\n" +
			"Note what the list contains: the account global Policy is **not** included (use " +
			"`incapsula_abp_policy` with `account_global = true` for it), while the `\"<site> " +
			"Default\"` Policies the API creates for every Site **are** — on a busy account " +
			"those outnumber hand-written Policies, and nothing in the returned data " +
			"distinguishes them beyond their name.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to list the Policies of.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"policies": {
				Description: "The Policies belonging to the account, except the account global Policy.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: policyAttributes,
				},
			},
		},
	}
}

func dataSourceAbpPoliciesRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	policies, err := client.ListAbpPolicies(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API defines no order for Policies, so sort to keep the list stable
	// across reads.
	slices.SortFunc(policies, func(a, b AbpPolicy) int {
		if byName := strings.Compare(a.Name, b.Name); byName != 0 {
			return byName
		}
		return strings.Compare(a.Id, b.Id)
	})

	flat := make([]any, 0, len(policies))
	for i := range policies {
		flat = append(flat, flattenAbpPolicy(&policies[i]))
	}

	data.SetId(accountId)
	if err := data.Set("policies", flat); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
