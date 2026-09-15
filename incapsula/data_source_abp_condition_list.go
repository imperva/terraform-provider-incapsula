package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpConditionList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpConditionListRead,

		Description: "Looks up an ABP Condition List by name. Names are case-sensitive " +
			"and matched exactly. The lookup fails if more than one condition list matches.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to search within.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Name of the condition list to look up.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"managed": {
				Description: "Whether the matched condition list is managed (Imperva-owned).",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"description": {
				Description: "Description of the condition list.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Condition List was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Condition List was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceAbpConditionListRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	name := data.Get("name").(string)

	conditions, err := client.ListAbpConditions(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	var match *AbpCondition
	for i := range conditions {
		c := &conditions[i]
		if c.Kind != AbpConditionKindList {
			continue
		}
		if c.Name != name {
			continue
		}
		if match != nil {
			return diag.Errorf("multiple ABP Condition Lists named %q found in account %s", name, accountId)
		}
		match = c
	}
	if match == nil {
		return diag.Errorf("no ABP Condition List named %q found in account %s", name, accountId)
	}

	data.SetId(match.Id)
	if err := data.Set("managed", match.AccountId == ""); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("description", match.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("created_at", match.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("modified_at", match.ModifiedAt); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
