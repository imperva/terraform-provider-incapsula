package incapsula

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"math"
	"strconv"
	"strings"
)

type ClientApps struct {
	Res            *int              `json:"res"`
	ResMessage     string            `json:"res_message"`
	DebugInfo      map[string]string `json:"debug_info"`
	ClientApps     map[string]string `json:"clientApps"`
	ClientAppTypes map[string]string `json:"clientAppTypes"`
}

func dataSourceClientApps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClientAppsRead,

		Description: "Provides the properties of all the client applications.",

		// Computed Attributes
		Schema: map[string]*schema.Schema{
			"filter": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by Client Application name",
				Optional:    true,
			},
			"map": {
				Type:        schema.TypeMap,
				Description: "Map of all the client applications",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"ids": {
				Type:        schema.TypeSet,
				Description: "Set of client applications ids filtered by name",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func UniqueStrings(stringSlices ...[]string) []string {
	uniqueMap := map[string]bool{}

	for _, stringSlice := range stringSlices {
		for _, string := range stringSlice {
			uniqueMap[string] = true
		}
	}

	// Create a slice with the capacity of unique items
	// This capacity make appending flow much more efficient
	result := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		result = append(result, key)
	}
	return result
}

func dataSourceClientAppsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	responseDTO, err := client.GetClientApplicationsMetadata()
	if err != nil {
		return diag.Errorf("Error getting Client Applications: %s", err)
	}

	if *responseDTO.Res != 0 || responseDTO.ResMessage != "OK" {
		if err != nil {
			panic(err)
		}
		return diag.Errorf("Error getting Client Applications Metadata: %v", responseDTO.DebugInfo)
	}

	// Preparing 2 maps for filter client name not found error
	// ClientAppsGroup contains the client names group by the first 2 letters to propose other suggestions
	ClientAppsLower := make(map[string]struct{}, len(responseDTO.ClientApps))
	ClientAppsGroup := make(map[string][]string, len(responseDTO.ClientApps))
	for _, clientName := range responseDTO.ClientApps {
		clientNameLower := strings.ToLower(clientName)
		ClientAppsLower[clientNameLower] = struct{}{}

		// We initiate supposing Client Name have 2 or more chars
		clientNameLower2Letters := clientNameLower[0:2]
		ClientAppsGroup[clientNameLower2Letters] = UniqueStrings(append(ClientAppsGroup[clientNameLower2Letters], clientName))
	}

	v, _ := d.GetOk("filter")
	filteredValues := v.(*schema.Set)
	filteredValuesMap := make(map[string]struct{}, len(filteredValues.List()))
	for _, v := range filteredValues.List() {
		botNameLowerCase := strings.ToLower(v.(string))
		filteredValuesMap[botNameLowerCase] = struct{}{}

		if _, ok := ClientAppsLower[botNameLowerCase]; !ok {
			proposalMessage := ""
			if len(botNameLowerCase) >= 2 && len(ClientAppsGroup[botNameLowerCase[0:2]]) > 0 {
				proposalClientNames := "'" + strings.Join(ClientAppsGroup[botNameLowerCase[0:2]], "','") + "'"
				proposalMessage = fmt.Sprintf("- Do you mean: %+v:", proposalClientNames)
			}
			return diag.Errorf("Client '%s' not found %s", v.(string), proposalMessage)
		}

	}

	var clientApps = make(map[string]int, len(responseDTO.ClientApps))
	var clientAppsIds = make([]int, 0)

	for clientId, clientName := range responseDTO.ClientApps {
		clientIdInt, _ := strconv.Atoi(clientId)
		clientNameLower := strings.ToLower(clientName)
		clientNameUpper := strings.ToUpper(clientName)
		if len(filteredValuesMap) > 0 {
			if _, ok := filteredValuesMap[clientNameLower]; ok {
				clientAppsIds = append(clientAppsIds, clientIdInt)
			}
		}
		clientApps[clientName] = clientIdInt
		// To avoid user typo error
		clientApps[clientNameLower] = clientIdInt
		clientApps[clientNameUpper] = clientIdInt
	}

	d.SetId(strconv.Itoa(math.MaxUint8))
	d.Set("map", clientApps)
	d.Set("ids", clientAppsIds)

	return nil
}
