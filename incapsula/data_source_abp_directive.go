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
				Description:  "ID of the ABP Policy to search within. Mutually exclusive with `account_global_policy`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"policy_id", "account_global_policy"},
			},
			"account_global_policy": {
				Description:  "If true, look up a directive within the account global policy. Requires `account_id` and is mutually exclusive with `policy_id`.",
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ExactlyOneOf: []string{"policy_id", "account_global_policy"},
				RequiredWith: []string{"account_id"},
			},
			"account_id": {
				Description:  "ABP account UUID. Required when `account_global_policy` is true.",
				Type:         schema.TypeString,
				Optional:     true,
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
	action := data.Get("action").(string)

	var policy *AbpPolicy
	var diags diag.Diagnostics
	var policyDescription, idPrefix string
	if data.Get("account_global_policy").(bool) {
		accountId := data.Get("account_id").(string)
		policy, diags = client.ReadAbpAccountGlobalPolicy(accountId)
		if diags != nil && diags.HasError() {
			return diags
		}
		if policy == nil {
			return diag.Errorf("ABP account global policy not found for account %s", accountId)
		}
		policyDescription = "ABP account global policy for account " + accountId
		idPrefix = "global:" + accountId
	} else {
		policyId := data.Get("policy_id").(string)
		policy, diags = client.ReadAbpPolicy(policyId)
		if diags != nil && diags.HasError() {
			return diags
		}
		if policy == nil {
			return diag.Errorf("ABP Policy %s not found", policyId)
		}
		policyDescription = "ABP Policy " + policyId
		idPrefix = policyId
	}

	var match *AbpDirective
	for i := range policy.Directives {
		if policy.Directives[i].Action != action {
			continue
		}
		if match != nil {
			return diag.Errorf("multiple directives with action %q found in %s", action, policyDescription)
		}
		match = &policy.Directives[i]
	}
	if match == nil {
		return diag.Errorf("no directive with action %q found in %s", action, policyDescription)
	}

	data.SetId(idPrefix + ":" + action)
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
