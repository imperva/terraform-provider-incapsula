package incapsula

import (
	"context"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// abpConditionsManagedAll lists account-owned and managed conditions.
	abpConditionsManagedAll = "all"
	// abpConditionsManagedOnly lists managed conditions only.
	abpConditionsManagedOnly = "managed_only"
	// abpConditionsManagedExcluded lists account-owned conditions only. This is
	// the default.
	abpConditionsManagedExcluded = "managed_excluded"
)

func dataSourceAbpConditions() *schema.Resource {
	conditionAttributes := map[string]*schema.Schema{
		"id": {
			Description: "ID of the Condition.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"account_id": {
			Description: "ABP account UUID the Condition belongs to. Empty for managed Conditions.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "Name of the Condition.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"managed": {
			Description: "Whether this Condition is managed (Imperva-owned) rather than account-owned.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
	maps.Copy(conditionAttributes, abpConditionReadOnlyAttributes())

	return &schema.Resource{
		ReadContext: dataSourceAbpConditionsRead,

		Description: "Lists the literal ABP Conditions of an account. Use this to " +
			"reference conditions that are not managed by this Terraform configuration when " +
			"all of them are of interest. Use `incapsula_abp_condition` to look up a single " +
			"Condition by `name` instead.\n\n" +
			"Only literal Conditions are returned. Condition lists and condition list " +
			"entries are not included here; see `incapsula_abp_condition_list`.\n\n" +
			"Note: unless `managed` says otherwise, only account-owned conditions are " +
			"listed. Managed conditions are excluded by default\n\n" +
			"Note: The time of reading this datasource is not defined by the terraform " +
			"module. That means that the read is subject to race-conditions. If a condition " +
			"is created or deleted in the same plan as this datasource is read, it may yield " +
			"unexpected results, like an inexhaustive listing. If deterministic operation is " +
			"required, specify resource dependencies explicitly to force terraform plan order.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID to list the Conditions of. Managed conditions are visible from any account.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.IsUUID,
			},
			"managed": {
				Description: "Which Conditions to list with respect to ownership. One of:\n\n" +
					"  - `managed_excluded` (default): account-owned conditions only\n" +
					"  - `managed_only`: managed conditions only\n" +
					"  - `all`: both",
				Type:     schema.TypeString,
				Optional: true,
				Default:  abpConditionsManagedExcluded,
				ValidateFunc: validation.StringInSlice([]string{
					abpConditionsManagedAll,
					abpConditionsManagedOnly,
					abpConditionsManagedExcluded,
				}, false),
			},
			"conditions": {
				Description: "All matching literal conditions, ordered by `name` and then by `id`.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: conditionAttributes,
				},
			},
		},
	}
}

// abpConditionManagedFilterMatches reports whether a condition with the given
// ownership passes the `managed` filter of the `incapsula_abp_conditions` data
// source. An unrecognized filter cannot occur: the schema validates it.
func abpConditionManagedFilterMatches(filter string, isManaged bool) bool {
	switch filter {
	case abpConditionsManagedOnly:
		return isManaged
	case abpConditionsManagedExcluded:
		return !isManaged
	default: // abpConditionsManagedAll
		return true
	}
}

func dataSourceAbpConditionsRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)
	managed := data.Get("managed").(string)

	conditions, err := client.ListAbpConditions(accountId)
	if err != nil {
		return diag.FromErr(err)
	}

	// The account endpoint returns all three Condition variants mixed together,
	// including managed Conditions which carry an empty AccountId.
	matches := make([]*AbpCondition, 0, len(conditions))
	for i := range conditions {
		c := &conditions[i]
		if c.Kind != AbpConditionKindLiteral {
			continue
		}
		if !abpConditionManagedFilterMatches(managed, c.AccountId == "") {
			continue
		}
		matches = append(matches, c)
	}

	// Unlike Sites, where the order returned by the API is the priority and thus
	// significant, the API defines no order for Conditions. Sort so that the list
	// is stable across reads and does not produce spurious diffs in consumers.
	slices.SortFunc(matches, func(a, b *AbpCondition) int {
		if byName := strings.Compare(a.Name, b.Name); byName != 0 {
			return byName
		}
		return strings.Compare(a.Id, b.Id)
	})

	flat := make([]any, 0, len(matches))
	for _, condition := range matches {
		flat = append(flat, flattenAbpCondition(condition))
	}

	data.SetId(accountId)
	if err := data.Set("conditions", flat); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
