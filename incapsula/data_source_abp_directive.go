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

		Description: "Looks up a single Directive within an ABP Policy by its action. " +
			"The lookup fails if more than one Directive in the Policy has the given action.",

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Artificial identifier for this data source. Directives have no " +
					"stable ID in the API, so this is synthesized from the Policy and action " +
					"(`<policy_id>:<action>`, or `global:<account_id>:<action>` for the account " +
					"global Policy) and is not a real API identifier.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy_id": {
				Description: "ID of the ABP Policy to search within. Mutually exclusive with " +
					"`account_global_policy`. This is a search argument, not an output: it stays " +
					"empty when the Directive is looked up in the account global Policy.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"policy_id", "account_global_policy"},
			},
			"account_global_policy": {
				Description: "If true, look up a Directive within the account global Policy. " +
					"Requires `account_id` and is mutually exclusive with `policy_id`. Only ever " +
					"write `true`: writing `false` explicitly — including from a variable — is " +
					"rejected with an \"exactly one of ...\" error, because the attribute counts as " +
					"set even when its value is the default. Omit it instead.",
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				ExactlyOneOf: []string{"policy_id", "account_global_policy"},
				RequiredWith: []string{"account_id"},
			},
			"account_id": {
				Description:  "ABP account UUID. Required when `account_global_policy` is true, and ignored otherwise.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsUUID,
			},
			"action": {
				Description:  "Action of the Directive to look up.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"condition_list_id": {
				Description: "Condition List containing Conditions for this Directive.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"skip_condition_list_id": {
				Description: "Condition List whose matches skip this Directive. Only set when `action` is `proof_of_work`; an empty string otherwise, never null.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"proof_of_work_configuration_id": {
				Description: "ID of the proof-of-work configuration applied. Only set when `action` is `proof_of_work`; an empty string otherwise, never null.",
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
	var err error
	var policyDescription, idPrefix string
	if data.Get("account_global_policy").(bool) {
		accountId := data.Get("account_id").(string)
		policy, err = client.ReadAbpAccountGlobalPolicy(accountId)
		if err != nil {
			return diag.FromErr(err)
		}
		if policy == nil {
			return diag.Errorf("ABP account global policy not found for account %s", accountId)
		}
		policyDescription = "ABP account global policy for account " + accountId
		idPrefix = "global:" + accountId
	} else {
		policyId := data.Get("policy_id").(string)
		policy, err = client.ReadAbpPolicy(policyId)
		if err != nil {
			return diag.FromErr(err)
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
