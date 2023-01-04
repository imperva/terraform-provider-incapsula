package incapsula

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math"
	"strconv"
	"strings"
)

func dataSourceAccountRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccountRolesRead,

		Description: "Provides the account roles of a given account.",

		// Computed Attributes
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			// Computed Attributes
			"admin_role_id": {
				Description: "Numeric identifier of the default Administrator role.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"reader_role_id": {
				Description: "Numeric identifier of the default Reader role.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"map": {
				Type:        schema.TypeMap,
				Description: "Set of all the account roles",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func dataSourceAccountRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountId := d.Get("account_id").(int)
	responseDTO, err := client.GetAccountRoles(accountId)
	if err != nil {
		return diag.Errorf("Error getting Account Roles: %s", err)
	}

	var accountPermissionsMap = make(map[string]int)

	for _, v := range *responseDTO {
		accountPermissionsMap[v.RoleName] = v.RoleId
		// To avoid user typo error
		accountPermissionsMap[strings.ToLower(v.RoleName)] = v.RoleId
		accountPermissionsMap[strings.ToUpper(v.RoleName)] = v.RoleId

		if v.RoleName == "Administrator" {
			d.Set("admin_role_id", v.RoleId)
		}
		if v.RoleName == "Reader" {
			d.Set("reader_role_id", v.RoleId)
		}
	}

	d.SetId(strconv.Itoa(math.MaxUint8))
	d.Set("map", accountPermissionsMap)

	return nil
}
