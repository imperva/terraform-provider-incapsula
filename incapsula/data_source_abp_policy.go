package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpPolicyRead,

		Description: "Looks up an ABP Policy by name or id. Exactly one of `name` or `id` " +
			"must be set. Names are case-sensitive and matched exactly; the lookup fails " +
			"if more than one policy matches.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to search within. Required when looking up by `name`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				RequiredWith: []string{"name"},
			},
			"name": {
				Description:  "Name of the policy to look up. Mutually exclusive with `id`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
				ExactlyOneOf: []string{"name", "id"},
			},
			"id": {
				Description:  "ID of the policy to look up. Mutually exclusive with `name`.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"name", "id"},
			},
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
	if match.Description != nil {
		if err := data.Set("description", *match.Description); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := data.Set("description", ""); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := data.Set("directive", flattenDirectives(match.Directives)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("modified_at", match.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func lookupAbpPolicy(client *Client, data *schema.ResourceData) (*AbpPolicy, diag.Diagnostics) {
	if id, ok := data.GetOk("id"); ok {
		policy, diags := client.ReadAbpPolicy(id.(string))
		if diags.HasError() {
			return nil, diags
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
