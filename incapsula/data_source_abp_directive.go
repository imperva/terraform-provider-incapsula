package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpDirective() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpDirectiveRead,

		Description: "Looks up a single directive within an ABP Policy by its action. " +
			"The lookup fails if more than one directive in the policy has the given action.",

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Description:  "ID of the ABP Policy to search within.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"action": {
				Description:  "Action of the directive to look up.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
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
	}
}

func dataSourceAbpDirectiveRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	policyId := data.Get("policy_id").(string)
	action := data.Get("action").(string)

	policy, diags := client.ReadAbpPolicy(policyId)
	if diags != nil && diags.HasError() {
		return diags
	}
	if policy == nil {
		return diag.Errorf("ABP Policy %s not found", policyId)
	}

	var match *AbpDirective
	for i := range policy.Directives {
		if policy.Directives[i].Action != action {
			continue
		}
		if match != nil {
			return diag.Errorf("multiple directives with action %q found in ABP Policy %s", action, policyId)
		}
		match = &policy.Directives[i]
	}
	if match == nil {
		return diag.Errorf("no directive with action %q found in ABP Policy %s", action, policyId)
	}

	data.SetId(policyId + ":" + action)
	if match.ConditionId != nil {
		if err := data.Set("condition_list_id", *match.ConditionId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := data.Set("condition_list_id", ""); err != nil {
			return diag.FromErr(err)
		}
	}
	if match.SkipConditionId != nil {
		if err := data.Set("skip_condition_list_id", *match.SkipConditionId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := data.Set("skip_condition_list_id", ""); err != nil {
			return diag.FromErr(err)
		}
	}
	if match.ProofOfWorkConfigurationId != nil {
		if err := data.Set("proof_of_work_configuration_id", *match.ProofOfWorkConfigurationId); err != nil {
			return diag.FromErr(err)
		}
	} else {
		if err := data.Set("proof_of_work_configuration_id", ""); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
