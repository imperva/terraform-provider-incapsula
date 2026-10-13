package incapsula

import (
	"context"
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpPolicy() *schema.Resource {
	policySchema := map[string]*schema.Schema{
		"account_id": {
			Description:  "ABP account UUID to search within. Required when looking up by `name` or `account_global`.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},
		"name": {
			Description:  "Name of the policy to look up. Mutually exclusive with `id` and `account_global`.",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
			ExactlyOneOf: []string{"name", "id", "account_global"},
			RequiredWith: []string{"account_id"},
		},
		"id": {
			Description:  "ID of the policy to look up. Mutually exclusive with `name` and `account_global`.",
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IsUUID,
			ExactlyOneOf: []string{"name", "id", "account_global"},
		},
		"account_global": {
			Description:  "If true, look up the account global policy. Requires `account_id` and is mutually exclusive with `name` and `id`.",
			Type:         schema.TypeBool,
			Optional:     true,
			Default:      false,
			ExactlyOneOf: []string{"name", "id", "account_global"},
			RequiredWith: []string{"account_id"},
		},
	}
	maps.Copy(policySchema, abpPolicyReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpPolicyRead,

		Description: "Looks up an ABP Policy by name, id, or by selecting the account global " +
			"policy. Exactly one of `name`, `id`, or `account_global` must be set. Names are " +
			"case-sensitive and matched exactly; the lookup fails if more than one policy matches. " +
			"Use `incapsula_abp_policies` to retrieve every Policy of an account instead.",

		Schema: policySchema,
	}
}

// abpPolicyReadOnlyAttributes returns the read-only Policy attributes shared by
// the `incapsula_abp_policy` and `incapsula_abp_policies` data sources, keeping
// the two in lockstep.
func abpPolicyReadOnlyAttributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Description: "Description of the policy.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"directive": {
			Description: "Ordered list of directives evaluated top-down for this policy.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"action": {
						Description: "Action taken when this directive matches.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"condition_list_id": {
						Description: "Condition list containing conditions for this directive.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"skip_condition_list_id": {
						Description: "Condition list whose matches skip this directive. Only set when `action` is `proof_of_work`.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"proof_of_work_configuration_id": {
						Description: "ID of the proof-of-work configuration applied. Only set when `action` is `proof_of_work`.",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
		"modified_at": {
			Description: "RFC3339 timestamp at which the Policy was last modified.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func dataSourceAbpPolicyRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)

	match, diags := lookupAbpPolicy(client, data)
	if diags.HasError() {
		return diags
	}

	data.SetId(match.Id)

	flat := flattenAbpPolicy(match)
	// `id` is the resource id and `name` is a user-supplied search argument.
	delete(flat, "id")
	delete(flat, "name")

	for _, key := range slices.Sorted(maps.Keys(flat)) {
		if err := data.Set(key, flat[key]); err != nil {
			return diag.Errorf("setting %s: %s", key, err)
		}
	}
	return nil
}

func lookupAbpPolicy(client *Client, data *schema.ResourceData) (*AbpPolicy, diag.Diagnostics) {
	if data.Get("account_global").(bool) {
		accountId := data.Get("account_id").(string)
		policy, err := client.ReadAbpAccountGlobalPolicy(accountId)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if policy == nil {
			return nil, diag.Errorf("no global ABP Policy found for account %s", accountId)
		}
		return policy, nil
	}

	if id, ok := data.GetOk("id"); ok {
		policy, err := client.ReadAbpPolicy(id.(string))
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if policy == nil {
			return nil, diag.Errorf("no ABP Policy with id %q found", id.(string))
		}
		return policy, nil
	}

	accountId := data.Get("account_id").(string)
	name := data.Get("name").(string)

	policies, err := client.ListAbpPolicies(accountId)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	var match *AbpPolicy
	for i := range policies {
		if policies[i].Name != name {
			continue
		}
		if match != nil {
			return nil, diag.Errorf("multiple ABP Policies named %q found in account %s", name, accountId)
		}
		match = &policies[i]
	}
	if match == nil {
		return nil, diag.Errorf("no ABP Policy named %q found in account %s", name, accountId)
	}
	return match, nil
}
