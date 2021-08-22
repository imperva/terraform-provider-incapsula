package incapsula

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
	"strings"
)

func dataSourceDataCenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataCenterRead,
		Description: "Provides the properties of a single Data Center. All 'filter_by_' arguments are optional. When specified, a logical AND operator is assumed.",

		Schema: map[string]*schema.Schema{
			// Computed Attributes
			"site_id": {
				Description: "Site ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"filter_by_id": {
				Type:        schema.TypeInt,
				Description: "Filter by Data Center internal ID",
				Optional:    true,
			},
			"filter_by_name": {
				Type:        schema.TypeString,
				Description: "Filter by Data Center name",
				Optional:    true,
			},
			"filter_by_is_enabled": {
				Description: "Filter by whether Data Center is enabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_by_is_active": {
				Description: "Filter by whether Data Center is active",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_by_is_standby": {
				Description: "Filter by whether Data Center is standby",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_by_is_rest_of_the_world": {
				Description: "Filter by whether this Data Center handles only traffic from geo regions that are not assigned to any other Data Center",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_by_is_content": {
				Description: "Filter by whether this Data Center will only handle traffic routed by Application Delivery Forward-to-DC rule.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_by_geo_location": {
				Description: "Filter by assigned geo location or the Data Center, which serves the rest of the world, if the geo location is not assigned to any other Data Center.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Data Center name",
				Computed:    true,
			},
			"is_enabled": {
				Description: "When true, Data Center is enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_active": {
				Description: "When false, Data Center in standby mode",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_rest_of_the_world": {
				Description: "When true and lb_algorithm is one of: GEO_PREFERED, GEO_REQUIRED, Then this Data Center will handle the traffic of all geo locations that are not assigned to any other Data Center",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_content": {
				Description: "When true, Data Center will only accept traffic routed by Application Delivery Forward-to-DC rule.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"ip_mode": {
				Type:        schema.TypeString,
				Description: "SINGLE_IP supports multiple processes on same origin server each listening to a different port, MULTIPLE_IP support multiple origin servers all listening to same port.",
				Computed:    true,
			},
			"lb_algorithm": {
				Type:        schema.TypeString,
				Description: "How to load balance between the servers of this data center.",
				Computed:    true,
			},
			"weight": {
				Type:        schema.TypeInt,
				Description: "The weight in pecentage of this Data Center. Populated only when Site's LB algorithem is WEIGHTED_LB.",
				Computed:    true,
			},
			"origin_pop": {
				Type:        schema.TypeString,
				Description: "The ID of the PoP that serves as an access point between Imperva and the customer’s origin server. E.g. \"lax\", for Los Angeles. When not specified, all Imperva PoPs can send traffic to this data center.",
				Computed:    true,
			},
			"geo_locations": {
				Type:        schema.TypeString,
				Description: "List of geo regions that this data center will serve. Populated if Site's LB algorithm = GEO_PREFERRED or GEO_REQUIRED. E.g. \"ASIA,AFRICA\"",
				Computed:    true,
			},
		},
	}
}

func dataSourceDataCenterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	responseDTO, err := client.GetDataCentersConfiguration(d.Get("site_id").(string))
	if err != nil {
		return diag.Errorf("Error getting Data Centers configuration for site (%s): %s", d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 {
		out, err := json.Marshal(responseDTO.Errors)
		if err != nil {
			panic(err)
		}
		return diag.Errorf("Error getting Data Centers configuration for site (%s): %s", d.Get("site_id"), string(out))
	}

	var matchedDC DataCenterStruct
	for _, dc := range responseDTO.Data[0].DataCenters {
		if v, ok := d.GetOk("filter_by_geo_location"); ok {
			found := false
			for _, geoLocation := range dc.GeoLocations {
				if geoLocation == v {
					found = true
				}
			}
			if found && len(matchedDC.OriginServers) > 0 {
				matchedDC.OriginServers = nil
			}
			if (!found) && ((!dc.IsRestOfTheWorld) || len(matchedDC.OriginServers) > 0) {
				continue
			}
		}
		if v, ok := d.GetOk("filter_by_id"); ok && (v != *dc.ID) {
			continue
		}
		if v, ok := d.GetOk("filter_by_name"); ok && (v != dc.Name) {
			continue
		}
		if v, ok := d.GetOk("filter_by_is_enabled"); ok && (v != dc.IsEnabled) {
			continue
		}
		if v, ok := d.GetOk("filter_by_is_active"); ok && (v != dc.IsActive) {
			continue
		}
		if v, ok := d.GetOk("filter_by_is_standby"); ok && (v == dc.IsActive) {
			continue
		}
		if v, ok := d.GetOk("filter_by_is_rest_of_the_world"); ok && (v != dc.IsRestOfTheWorld) {
			continue
		}
		if v, ok := d.GetOk("filter_by_is_content"); ok && (v != dc.IsContent) {
			continue
		}

		if len(matchedDC.OriginServers) == 0 {
			matchedDC = dc
		} else {
			return diag.Errorf("More than one Data Center matched specified filters for site (%s). First two matches are for DC names: %s and %s", d.Get("site_id"), matchedDC.Name, dc.Name)
		}
	}

	if len(matchedDC.OriginServers) == 0 {
		return diag.Errorf("No Data Center matched specified filters for site (%s)", d.Get("site_id"))
	}

	d.SetId(strconv.FormatInt(int64(*matchedDC.ID), 10))

	d.Set("name", matchedDC.Name)
	d.Set("is_enabled", matchedDC.IsEnabled)
	d.Set("is_active", matchedDC.IsActive)
	d.Set("is_rest_of_the_world", matchedDC.IsRestOfTheWorld)
	d.Set("is_content", matchedDC.IsContent)
	d.Set("ip_mode", matchedDC.IpMode)
	d.Set("lb_algorithm", matchedDC.DcLbAlgorithm)
	d.Set("weight", matchedDC.Weight)
	d.Set("origin_pop", matchedDC.OriginPoP)
	d.Set("geo_locations", strings.Join(matchedDC.GeoLocations, ","))

	return nil
}
