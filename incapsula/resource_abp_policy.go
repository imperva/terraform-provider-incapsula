package incapsula

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AbpPolicy struct {
	Id          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Directives  []AbpDirective `json:"directives"`
	ModifiedAt  string         `json:"modified_at,omitempty"`
}

type AbpDirective struct {
	Action      string  `json:"action"`
	ConditionId *string `json:"condition_id,omitempty"`
}

func resourceAbpPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPolicyCreate,
		ReadContext:   resourceAbpPolicyRead,
		UpdateContext: resourceAbpPolicyUpdate,
		DeleteContext: resourceAbpPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
			},
		},

		Description: "Incapsula ABP policy resource\n",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account this policy belongs to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Policy name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Optional policy description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"directive": {
				Description: "Ordered list of directives evaluated top-down for this policy.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "Action to take when this directive matches.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"condition_id": {
							Description: "Optional condition this directive applies to. If omitted, the backend generates one.",
							Type:        schema.TypeString,
							Computed:    true,
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
		if cid, ok := m["condition_id"].(string); ok && cid != "" {
			d.ConditionId = &cid
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
			m["condition_id"] = *d.ConditionId
		}
		out[i] = m
	}
	return out
}

func resourceAbpPolicyCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountId := data.Get("account_id").(string)

	policy := AbpPolicy{
		Name:        data.Get("name").(string),
		Description: extractDescription(data),
		Directives:  extractDirectives(data),
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

	policy := AbpPolicy{
		Name:        data.Get("name").(string),
		Description: extractDescription(data),
		Directives:  extractDirectives(data),
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
