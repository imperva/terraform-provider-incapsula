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

func dataSourceAbpConditions() *schema.Resource {
	conditionAttributes := map[string]*schema.Schema{
		"id": {
			Description: "ID of the Condition.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"account_id": {
			Description: "ABP account UUID the Condition belongs to. Empty for managed Conditions.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Name of the Condition.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"managed": {
			Description: "Whether the condition is managed rather than account-owned.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
	maps.Copy(conditionAttributes, abpConditionReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpConditionsRead,

		Description: "Lists the literal ABP Conditions of an account, ordered by name. " +
			"Condition lists and condition list entries are not included. Use " +
			"`incapsula_abp_condition` to look up a single Condition by name instead.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to list the Conditions of.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"managed": {
				Description: "Only list managed conditions. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"account_owned": {
				Description: "Only list account-owned conditions. Defaults to the opposite of `managed`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"conditions": {
				Description: "The matching Conditions.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: conditionAttributes,
				},
			},
		},
	}
}

func dataSourceAbpConditionsRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	onlyManaged := data.Get("managed").(bool)

	// account_owned defaults to the opposite of managed, which a schema Default
	// cannot express, so tell unset apart from false via the raw config.
	onlyAccountOwned := !onlyManaged
	if raw := data.GetRawConfig(); !raw.IsNull() {
		if v := raw.GetAttr("account_owned"); v.IsKnown() && !v.IsNull() {
			onlyAccountOwned = v.True()
		}
	}

	if onlyManaged == onlyAccountOwned {
		return diag.Errorf("exactly one of managed or account_owned must be true")
	}

	conditions, err := client.ListAbpConditions(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	// The account endpoint returns all three Condition variants mixed together,
	// including managed Conditions, which carry an empty AccountId.
	matches := make([]*AbpCondition, 0, len(conditions))
	for i := range conditions {
		c := &conditions[i]
		if c.Kind != AbpConditionKindLiteral {
			continue
		}
		if isManaged := c.AccountId == ""; isManaged != onlyManaged {
			continue
		}
		matches = append(matches, c)
	}

	// The API defines no order for Conditions, so sort to keep the list stable
	// across reads.
	slices.SortFunc(matches, func(a, b *AbpCondition) int {
		if byName := strings.Compare(a.Name, b.Name); byName != 0 {
			return byName
		}
		return strings.Compare(a.Id, b.Id)
	})

	flat := make([]any, 0, len(matches))
	for _, condition := range matches {
		flat = append(flat, flattenAbpCondition(condition))
	}

	data.SetId(accountId)
	if err := data.Set("conditions", flat); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
