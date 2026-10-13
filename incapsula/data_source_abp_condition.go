package incapsula

import (
	"context"
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpCondition() *schema.Resource {
	conditionSchema := map[string]*schema.Schema{
		"account_id": {
			Description:  "ABP account UUID to search within. Managed conditions are visible from any account.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.IsUUID,
		},
		"name": {
			Description:  "Name of the literal condition to look up.",
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},
		"managed": {
			Description: "If `true`, the lookup is restricted to managed (Imperva-owned) " +
				"conditions. If `false` or unset, account-owned conditions are matched too. " +
				"Reflects whether the matched condition is managed.",
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
	maps.Copy(conditionSchema, abpConditionReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpConditionRead,

		Description: "Looks up a literal ABP Condition by name. Names are case-sensitive " +
			"and matched exactly. The lookup fails if more than one condition matches. " +
			"Use `incapsula_abp_conditions` to retrieve every literal Condition of an " +
			"account instead.",

		Schema: conditionSchema,
	}
}

// abpConditionReadOnlyAttributes returns the read-only literal Condition
// attributes shared by the `incapsula_abp_condition` and
// `incapsula_abp_conditions` data sources, keeping the two in lockstep.
func abpConditionReadOnlyAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "Description of the condition.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"code": {
			Description: "Server-side normalized MOI expression.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"template": {
			Description: "Editor template the server uses to validate `code`. Open enumeration; " +
				"`custom` for conditions that do not use a specific editor template.",
			Type:     schema.TypeString,
			Computed: true,
		},
		"last_change_by": {
			Description: "Identifier of the user who last changed this condition.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_at": {
			Description: "RFC3339 timestamp at which the Condition was created.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"modified_at": {
			Description: "RFC3339 timestamp at which the Condition was last modified.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func dataSourceAbpConditionRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	name := data.Get("name").(string)
	managedOnly := data.Get("managed").(bool)

	conditions, err := client.ListAbpConditions(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	var match *AbpCondition
	for i := range conditions {
		c := &conditions[i]
		if c.Kind != AbpConditionKindLiteral {
			continue
		}
		if c.Name != name {
			continue
		}
		isManaged := c.AccountId == ""
		if managedOnly && !isManaged {
			continue
		}
		if match != nil {
			return diag.Errorf("multiple ABP literal Conditions named %q found in account %s", name, accountId)
		}
		match = c
	}
	if match == nil {
		if managedOnly {
			return diag.Errorf("no managed ABP literal Condition named %q found in account %s", name, accountId)
		}
		return diag.Errorf("no ABP literal Condition named %q found in account %s", name, accountId)
	}

	data.SetId(match.Id)

	flat := flattenAbpCondition(match)
	// `id` is the resource id. `account_id` and `name` are user-supplied search
	// arguments; `account_id` in particular must not be written back, since a
	// managed Condition carries an empty AccountId which would clear it.
	delete(flat, "id")
	delete(flat, "account_id")
	delete(flat, "name")

	for _, key := range slices.Sorted(maps.Keys(flat)) {
		if err := data.Set(key, flat[key]); err != nil {
			return diag.Errorf("setting %s: %s", key, err)
		}
	}
	return nil
}
