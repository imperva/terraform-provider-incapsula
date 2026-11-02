package incapsula

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// abpDirectiveActionRegexp matches a valid directive action: snake_case only.
var abpDirectiveActionRegexp = regexp.MustCompile(`^[_a-z][_a-z0-9]*$`)

type AbpPolicy struct {
	Id          string         `json:"id,omitempty"`
	AccountId   string         `json:"account_id,omitempty"`
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
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpPolicyImport,
		},

		Description: `Provides an ABP Policy resource. A Policy is an ordered collection of
directives, each pairing an action with the conditions that trigger it, applied
to traffic that a Site's selector maps to this Policy.

All ` + "`directive`" + ` blocks must be specified explicitly. Each directive exposes a
` + "`condition_list_id`" + ` to which conditions are attached via
` + "`incapsula_abp_condition_list_entry`" + `.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this Policy belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Policy name. 1..100 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description: "Description of the policy. Set to empty string if omitted",
				Type:        schema.TypeString,
				Required:    true,
			},
			"directive": {
				Description: "Ordered list of directives evaluated top-down for this policy. A policy must have at least one directive.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
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

func extractAbpPolicy(data *schema.ResourceData) AbpPolicy {
	description := data.Get("description").(string)
	return AbpPolicy{
		Name:        data.Get("name").(string),
		Description: &description,
		Directives:  extractAbpDirectives(data),
	}
}

func extractAbpDirectives(data *schema.ResourceData) []AbpDirective {
	raw := data.Get("directive").([]any)
	directives := make([]AbpDirective, len(raw))
	for i, item := range raw {
		m := item.(map[string]any)
		d := AbpDirective{Action: m["action"].(string)}
		if cid, ok := m["condition_list_id"].(string); ok && cid != "" {
			d.ConditionId = &cid
		}
		if skipCid, ok := m["skip_condition_list_id"].(string); ok && skipCid != "" {
			d.SkipConditionId = &skipCid
		}
		if powId, ok := m["proof_of_work_configuration_id"].(string); ok && powId != "" {
			d.ProofOfWorkConfigurationId = &powId
		}
		directives[i] = d
	}
	return directives
}

func flattenAbpDirectives(directives []AbpDirective) []any {
	out := make([]any, len(directives))
	for i, d := range directives {
		m := map[string]any{"action": d.Action}
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

// serializeAbpPolicy writes the server's view of a policy back into state.
func serializeAbpPolicy(data *schema.ResourceData, policy *AbpPolicy) error {
	if policy.AccountId != "" {
		if err := data.Set("account_id", policy.AccountId); err != nil {
			return err
		}
	}
	if err := data.Set("name", policy.Name); err != nil {
		return err
	}
	if policy.Description != nil {
		if err := data.Set("description", policy.Description); err != nil {
			return err
		}
	} else {
		if err := data.Set("description", ""); err != nil {
			return err
		}
	}
	if err := data.Set("directive", flattenAbpDirectives(policy.Directives)); err != nil {
		return err
	}
	return nil
}

func resourceAbpPolicyCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	created, err := client.CreateAbpPolicy(accountId, extractAbpPolicy(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Policy create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpPolicy(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Policy %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpPolicyRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	policy, err := client.ReadAbpPolicy(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if policy == nil {
		log.Printf("[INFO] ABP Policy %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpPolicy(data, policy); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpPolicyUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	updated, err := client.UpdateAbpPolicy(id, extractAbpPolicy(data))
	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpPolicy(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpPolicyDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpPolicy(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpPolicyImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<policy_id>'")
	}

	client := m.(*Client)
	policy, err := client.ReadAbpPolicy(id)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("ABP Policy %s not found", id)
	}

	data.SetId(id)
	if err := serializeAbpPolicy(data, policy); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
