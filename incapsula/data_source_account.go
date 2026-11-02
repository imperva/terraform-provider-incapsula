package incapsula

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAccount() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRead,
		Description: "Finds an account below the current account (a reseller's accounts, or a parent account's sub accounts) by its name.",

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_name": {
				Description: "The name of the account to find. Matching is case-insensitive.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Computed Attributes
			"account_id": {
				Description: "The ID of the account.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"parent_id": {
				Description: "The ID of the parent account.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

// findAccountsByName returns all the accounts whose name equals the given name, ignoring case.
func findAccountsByName(accounts []SubAccountItem, name string) []SubAccountItem {
	var matches []SubAccountItem
	for _, account := range accounts {
		if strings.EqualFold(account.SubAccountName, name) {
			matches = append(matches, account)
		}
	}
	return matches
}

func dataSourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	name := d.Get("account_name").(string)

	accounts, err := client.ListSubAccounts()
	if err != nil {
		return diag.Errorf("Error listing accounts: %v", err)
	}

	matches := findAccountsByName(accounts, name)
	switch len(matches) {
	case 0:
		return diag.Errorf("No account named %q was found below the current account", name)
	case 1:
	default:
		ids := make([]string, len(matches))
		for i, match := range matches {
			ids[i] = strconv.Itoa(match.SubAccountID)
		}
		return diag.Errorf("Found %d accounts named %q (case-insensitive), IDs: %s", len(matches), name, strings.Join(ids, ", "))
	}

	account := matches[0]
	id := strconv.Itoa(account.SubAccountID)
	d.SetId(id)
	d.Set("account_id", id)
	d.Set("account_name", account.SubAccountName)
	d.Set("parent_id", fmt.Sprint(account.ParentID))
	return nil
}
