package incapsula

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// abpDirectiveActionRegexp matches a valid directive action: snake_case only.
var abpDirectiveActionRegexp = regexp.MustCompile(`^[_a-z][_a-z0-9]*$`)

// standardDirectiveActions mirrors the actions the ABP frontend creates
// when a user picks "Standard Directives" in the create-policy modal.
var standardDirectiveActions = []string{
	"allow",
	"block",
	"captcha_cleared",
	"captcha",
	"identify",
	"tarpit",
	"delay",
}

func standardDirectives() []AbpDirective {
	directives := make([]AbpDirective, len(standardDirectiveActions))
	for i, action := range standardDirectiveActions {
		directives[i] = AbpDirective{Action: action}
	}
	return directives
}

type AbpPolicy struct {
	Id          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Directives  []AbpDirective `json:"directives"`
	ModifiedAt  string         `json:"modified_at,omitempty"`
}

type AbpDirective struct {
	Action                     string  `json:"action"`
	ConditionId                *string `json:"condition_id,omitempty"`
	SkipConditionId            *string `json:"skip_condition_id,omitempty"`
	ProofOfWorkConfigurationId *string `json:"proof_of_work_configuration_id,omitempty"`
}

func resourceAbpPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPolicyCreate,
		ReadContext:   resourceAbpPolicyRead,
		UpdateContext: resourceAbpPolicyUpdate,
		DeleteContext: resourceAbpPolicyDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			if !d.Get("use_standard_directives").(bool) {
				return nil
			}
			directiveCfg := d.GetRawConfig().GetAttr("directive")
			if !directiveCfg.IsNull() && directiveCfg.LengthInt() > 0 {
				return fmt.Errorf("`directive` blocks must not be set when `use_standard_directives` is true")
			}
			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Description: `Provides an ABP Policy resource. A Policy is an ordered collection of
directives, each pairing an action with the conditions that trigger it, applied
to traffic that a Site's selector maps to this Policy.

Set ` + "`use_standard_directives`" + ` to create the Policy with the standard set of
directives (matching the ABP UI's "Standard Directives" choice), or provide
explicit ` + "`directive`" + ` blocks. Each directive exposes a ` + "`condition_list_id`" + ` to
which conditions are attached via ` + "`incapsula_abp_condition_list_entry`" + `.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account this policy belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description:  "Policy name. 1..100 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description: "Optional policy description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"use_standard_directives": {
				Description: "If true, the policy is created with the standard set of directives (matching the ABP UI's \"Standard Directives\" choice). When set, custom `directive` blocks must not be specified.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"directive": {
				Description: "Ordered list of directives evaluated top-down for this policy. A policy must have at least one directive. Computed when `use_standard_directives` is true; otherwise required.",
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "Action to take when this directive matches. Must be snake_case, 1..100 characters.",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 100),
								validation.StringMatch(
									abpDirectiveActionRegexp,
									"action must be snake_case matching ^[_a-z][_a-z0-9]*$",
								),
							),
						},
						"condition_list_id": {
							Description: "Condition list containing conditions for this directive.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"skip_condition_list_id": {
							Description: "Condition list whose matches skip this directive. Only meaningful when `action` is `proof_of_work`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"proof_of_work_configuration_id": {
							Description: "ID of the proof-of-work configuration to apply. Only valid when `action` is `proof_of_work`.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func extractDescription(data *schema.ResourceData) *string {
	if v, ok := data.GetOk("description"); ok {
		s := v.(string)
		return &s
	}
	return nil
}

func extractDirectives(data *schema.ResourceData) []AbpDirective {
	raw := data.Get("directive").([]interface{})
	directives := make([]AbpDirective, len(raw))
	for i, item := range raw {
		m := item.(map[string]interface{})
		d := AbpDirective{Action: m["action"].(string)}
		if cid, ok := m["condition_list_id"].(string); ok && cid != "" {
			d.ConditionId = &cid
		}
		if powId, ok := m["proof_of_work_configuration_id"].(string); ok && powId != "" {
			d.ProofOfWorkConfigurationId = &powId
		}
		directives[i] = d
	}
	return directives
}

func flattenDirectives(directives []AbpDirective) []interface{} {
	out := make([]interface{}, len(directives))
	for i, d := range directives {
		m := map[string]interface{}{"action": d.Action}
		if d.ConditionId != nil {
			m["condition_list_id"] = *d.ConditionId
		}
		if d.SkipConditionId != nil {
			m["skip_condition_list_id"] = *d.SkipConditionId
		}
		if d.ProofOfWorkConfigurationId != nil {
			m["proof_of_work_configuration_id"] = *d.ProofOfWorkConfigurationId
		}
		out[i] = m
	}
	return out
}

func resourceAbpPolicyCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountId := data.Get("account_id").(string)

	var directives []AbpDirective
	if data.Get("use_standard_directives").(bool) {
		directives = standardDirectives()
	} else {
		directives = extractDirectives(data)
	}

	policy := AbpPolicy{
		Name:        data.Get("name").(string),
		Description: extractDescription(data),
		Directives:  directives,
	}

	created, diags := client.CreateAbpPolicy(accountId, policy)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to create ABP policy for Account ID %s", accountId)
		return diags
	}

	data.SetId(created.Id)

	err := data.Set("directive", flattenDirectives(created.Directives))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAbpPolicyRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	policyId := data.Id()

	policy, diags := client.ReadAbpPolicy(policyId)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to read ABP policy ID %s", policyId)
		return diags
	}

	if policy == nil {
		log.Printf("[INFO] ABP policy ID %s no longer exists upstream, clearing state", policyId)
		data.SetId("")
		return diags
	}

	err := data.Set("name", policy.Name)
	if err != nil {
		return diag.FromErr(err)
	}
	if policy.Description != nil {
		err = data.Set("description", *policy.Description)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = data.Set("description", "")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	err = data.Set("directive", flattenDirectives(policy.Directives))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAbpPolicyUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	policyId := data.Id()

	var directives []AbpDirective
	if data.Get("use_standard_directives").(bool) {
		directives = standardDirectives()
	} else {
		directives = extractDirectives(data)
	}

	policy := AbpPolicy{
		Name:        data.Get("name").(string),
		Description: extractDescription(data),
		Directives:  directives,
	}

	updated, diags := client.UpdateAbpPolicy(policyId, policy)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to update ABP policy ID %s", policyId)
		return diags
	}

	err := data.Set("directive", flattenDirectives(updated.Directives))
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAbpPolicyDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	policyId := data.Id()

	diags := client.DeleteAbpPolicy(policyId)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to delete ABP policy ID %s", policyId)
		return diags
	}

	data.SetId("")
	return diags
}
