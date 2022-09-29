package incapsula

import (
	"context"
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

	v, _ := d.GetOk("filter")
	filteredValues := v.(*schema.Set)
	filteredValuesMap := make(map[string]struct{}, len(filteredValues.List()))
	for _, v := range filteredValues.List() {
		filteredValuesMap[strings.ToLower(v.(string))] = struct{}{}
	}

	var clientApps = make(map[string]int, len(responseDTO.ClientApps))
	var clientAppsIds = make([]int, 0)

	for clientId, clientName := range responseDTO.ClientApps {
		clientIdInt, _ := strconv.Atoi(clientId)
		if len(filteredValuesMap) > 0 {
			if _, ok := filteredValuesMap[strings.ToLower(clientName)]; ok {
				clientAppsIds = append(clientAppsIds, clientIdInt)
			}
		}
		clientApps[clientName] = clientIdInt
		// To avoid user typo error
		clientApps[strings.ToLower(clientName)] = clientIdInt
		clientApps[strings.ToUpper(clientName)] = clientIdInt
	}

	d.SetId(strconv.Itoa(math.MaxUint8))
	d.Set("map", clientApps)
	d.Set("ids", clientAppsIds)

	return nil
}
