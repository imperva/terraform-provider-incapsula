package incapsula

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
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
			"map": {
				Type:        schema.TypeMap,
				Description: "Map of all the client applications",
				StateFunc: func(val any) string {
					log.Printf("StateFunc IN")
					return strings.ToLower(val.(string))
				},
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
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

	var clientApps = make(map[string]string, len(responseDTO.ClientApps))
	for clientId, clientName := range responseDTO.ClientApps {
		clientApps[clientName] = clientId
		// To avoid user typo error
		clientApps[strings.ToLower(clientName)] = clientId
		clientApps[strings.ToUpper(clientName)] = clientId
	}

	d.SetId(strconv.Itoa(math.MaxUint8))
	d.Set("map", clientApps)
	log.Printf("[DEBUG] map: %v\n", d.Get("map"))

	return nil
}
