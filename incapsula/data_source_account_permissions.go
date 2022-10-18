package incapsula

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math"
	"strconv"
	"strings"
)

// AbilitiesDetailsDTO - Same DTO for: GET response, POST request, and POST response
type AbilitiesDetailsDTO struct {
	RoleAbilities []RoleAbility `json:"roleAbilities"`
	ErrorCode     int           `json:"errorCode"`
	Description   string        `json:"description"`
}

func dataSourceAccountPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClientAppsRead,

		Description: "Provides the properties of all the client applications.",

		// Computed Attributes
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Numeric identifier of the account to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			// Optional Arguments
			"filter_by_text": {
				Type:        schema.TypeString,
				Description: "Filter by text representing the permission to fetch",
				Optional:    true,
			},

			// Computed Attributes
			"map": {
				Type:        schema.TypeMap,
				Description: "Set of all the ability keys",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"keys": {
				Type:        schema.TypeSet,
				Description: "Set of ability keys ids filtered by name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceClientAppsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountId := d.Get("account_id").(int)
	responseDTO, err := client.GetAccountAbilities(accountId)
	if err != nil {
		return diag.Errorf("Error getting Account Permissions: %s", err)
	}

	var accountPermissionsMap = make(map[string]string)
	var accountPermissionsKeys = make([]string, 0)

	textToFilter, filterExists := d.GetOk("filter_by_text")

	for _, v := range *responseDTO {
		if filterExists {
			displayNameLowerCase := strings.ToLower(v.AbilityDisplayName)
			textToFilterLowerCase := strings.ToLower(textToFilter.(string))
			if strings.Contains(displayNameLowerCase, textToFilterLowerCase) {
				accountPermissionsKeys = append(accountPermissionsKeys, v.AbilityKey)
			}
		}
		accountPermissionsMap[v.AbilityDisplayName] = v.AbilityKey
		// To avoid user typo error
		accountPermissionsMap[strings.ToLower(v.AbilityDisplayName)] = v.AbilityKey
		accountPermissionsMap[strings.ToUpper(v.AbilityDisplayName)] = v.AbilityKey
	}

	d.SetId(strconv.Itoa(math.MaxUint8))
	d.Set("map", accountPermissionsMap)
	d.Set("keys", accountPermissionsKeys)

	return nil
}
