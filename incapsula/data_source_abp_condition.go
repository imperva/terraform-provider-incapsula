package incapsula

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAbpCondition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAbpConditionRead,

		Description: "Looks up a literal ABP Condition by name. Names are case-sensitive " +
			"and matched exactly. The lookup fails if more than one condition matches.",

		Schema: map[string]*schema.Schema{
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
	if err := data.Set("managed", match.AccountId == ""); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("description", match.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("code", match.Code); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("last_change_by", match.LastChangeBy); err != nil {
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
