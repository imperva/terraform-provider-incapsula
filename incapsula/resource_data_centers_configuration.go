package incapsula

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDataCentersConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataCentersConfigurationCreate,
		Read:   resourceDataCentersConfigurationRead,
		Update: resourceDataCentersConfigurationCreate,
		Delete: resourceDataCentersConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("site_id", d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"site_lb_algorithm": {
				Description: "How to load balance between multiple Data Centers.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "BEST_CONNECTION_TIME",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					strVal := val.(string)
					allowedVals := []string{"BEST_CONNECTION_TIME", "GEO_PREFERRED", "GEO_REQUIRED", "WEIGHTED_LB"}
					if !isValidEnum(strVal, key, allowedVals) {
						errs = append(errs, fmt.Errorf("%q must be one of: [%s]. Got: %s",
							key, strings.Join(allowedVals, ","), strVal))
					}
					return
				},
			},
			"fail_over_required_monitors": {
				Description: "How many Imperva PoPs should assess Data Center as down before failover is performed.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "MOST",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					strVal := val.(string)
					allowedVals := []string{"ONE", "MANY", "MOST", "ALL"}
					if !isValidEnum(strVal, key, allowedVals) {
						errs = append(errs, fmt.Errorf("%q must be one of: [%s]. Got: %s",
							key, strings.Join(allowedVals, ","), strVal))
					}
					return
				},
			},
			"site_topology": {
				Description: "One of: 'SINGLE_SERVER' - No LB supported, 'SINGLE_DC' - Multiple servers on single Data Center, or 'MULTIPLE_DC' - Multiple Data Centers having multiple origin servers.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "SINGLE_DC",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					strVal := val.(string)
					allowedVals := []string{"SINGLE_SERVER", "SINGLE_DC", "MULTIPLE_DC"}
					if !isValidEnum(strVal, key, allowedVals) {
						errs = append(errs, fmt.Errorf("%q must be one of: [%s]. Got: %s",
							key, strings.Join(allowedVals, ","), strVal))
					}
					return
				},
			},
			"min_available_servers_for_dc_up": {
				Description: "The minimal number of available data center's servers to consider that data center as UP",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
			},
			"kickstart_url": {
				Description: "The URL that will be sent to the standby server when Imperva performs failover based on our monitoring. E.g. https://www.example.com/kickStart",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"kickstart_user": {
				Description: "User name, if required by the kickstart URL.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"kickstart_password": {
				Description: "Password, if required by the kickstart URL.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_persistent": {
				Description: "When true our proxy servers will maintain session stickiness to origin servers by a cookie.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"data_center": {
				Description: "A set of Data Centers and their Origin Servers",
				Required:    true,
				MinItems:    1,
				Type:        schema.TypeSet,
				Set:         resourceDataCentersConfigurationDataCenterHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Data Center name",
							Required:    true,
						},

						"dc_id": {
							Type:        schema.TypeInt,
							Description: "Internal Data Center id.",
							Computed:    true,
						},

						"ip_mode": {
							Type:        schema.TypeString,
							Description: "SINGLE_IP supports multiple processes on same origin server each listening to a different port, MULTIPLE_IP support multiple origin servers all listening to same port.",
							Optional:    true,
							Default:     "MULTIPLE_IP",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								strVal := val.(string)
								allowedVals := []string{"SINGLE_IP", "MULTIPLE_IP"}
								if !isValidEnum(strVal, key, allowedVals) {
									errs = append(errs, fmt.Errorf("%q must be one of: [%s]. Got: %s",
										key, strings.Join(allowedVals, ","), strVal))
								}
								return
							},
						},

						"web_servers_per_server": {
							Type:        schema.TypeInt,
							Description: "When IP mode = SINGLE_IP, number of web servers (processes) per server. Each web server listens to different port. E.g. when web_servers_per_server = 5, HTTP traffic will use ports 80-84 while HTTPS traffic will use ports 443-447",
							Optional:    true,
							Default:     1,
						},

						"dc_lb_algorithm": {
							Type:        schema.TypeString,
							Description: "How to load balance between the servers of this data center.",
							Optional:    true,
							Default:     "LB_LEAST_PENDING_REQUESTS",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								strVal := val.(string)
								allowedVals := []string{"LB_LEAST_PENDING_REQUESTS", "LB_LEAST_OPEN_CONNECTIONS", "LB_SOURCE_IP_HASH", "RANDOM", "WEIGHTED"}
								if !isValidEnum(strVal, key, allowedVals) {
									errs = append(errs, fmt.Errorf("%q must be one of: [%s]. Got: %s",
										key, strings.Join(allowedVals, ","), strVal))
								}
								return
							},
						},

						"weight": {
							Type:        schema.TypeInt,
							Description: "When site_lb_algorithm = WEIGHTED_LB, the weight in pecentage of this Data Center. Then, total weight of all Data Centers must be 100.",
							Optional:    true,
							Default:     0,
						},

						"is_enabled": {
							Type:        schema.TypeBool,
							Description: "Specifies if the Data Center is enabled.",
							Optional:    true,
							Default:     true,
						},

						"is_active": {
							Type:        schema.TypeBool,
							Description: "Specifies if the Data Center is an active or a standby Data Center. No more than one standby Data Center can be defined.",
							Optional:    true,
							Default:     true,
						},

						"is_content": {
							Type:        schema.TypeBool,
							Description: "When true, this data center will only handle requests that were routed to it using application delivery forward rules. If true, must be an enabled data center.",
							Optional:    true,
							Default:     false,
						},
						"is_rest_of_the_world": {
							Type:        schema.TypeBool,
							Description: "When true and site_lb_algorithm = GEO_PREFERRED or GEO_REQUIRED, exactly one data center must have is_rest_of_the_world = true. This data center will handle traffic from any region that is not assigned to a specific data center.",
							Optional:    true,
							Default:     false,
						},

						"geo_locations": {
							Type:        schema.TypeString,
							Description: "List of geo regions that this data center will serve. Mandatory if site_lb_algorithm = GEO_PREFERRED or GEO_REQUIRED. E.g. \"ASIA,AFRICA\"",
							Optional:    true,
							Default:     "",
							ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
								strVals := strings.Split(val.(string), ",")
								allowedVals := []string{"EUROPE", "AUSTRALIA", "US_EAST", "US_WEST", "AFRICA", "ASIA", "SOUTH_AMERICA", "NORTH_AMERICA"}
								for _, strVal := range strVals {
									if strVal != "" && !isValidEnum(strVal, key, allowedVals) {
										errs = append(errs, fmt.Errorf("%q must be an empty string or any of: [%s]. Got: %s",
											key, strings.Join(allowedVals, ","), strVal))
									}
								}
								return
							},
						},

						"origin_pop": {
							Type:        schema.TypeString,
							Description: "The ID of the PoP that serves as an access point between Imperva and the customer’s origin server. E.g. \"lax\", for Los Angeles. When not specified, all Imperva PoPs can send traffic to this data center. The list of available PoPs is documented at: https://docs.imperva.com/bundle/cloud-application-security/page/more/pops.htm",
							Optional:    true,
							Default:     "",
						},

						"origin_server": {
							Description: "A set of Origin Servers",
							Required:    true,
							MinItems:    1,
							Type:        schema.TypeSet,
							Set:         resourceDataCentersConfigurationOriginServerHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": {
										Type:        schema.TypeString,
										Description: "Origin server address. can be specified as IPv4, IPv6, or DNS host name.",
										Required:    true,
									},

									"weight": {
										Type:        schema.TypeInt,
										Description: "When dc_lb_algorithm = WEIGHTED, the weight in pecentage of this origin server. Then, total weight of all origin server on each Data Center must be 100.",
										Optional:    true,
										Default:     0,
									},

									"is_enabled": {
										Type:        schema.TypeBool,
										Description: "Specifies if the origin server is enabled.",
										Optional:    true,
										Default:     true,
									},

									"is_active": {
										Type:        schema.TypeBool,
										Description: "Specifies if the origin server is an active or a standby origin server.",
										Optional:    true,
										Default:     true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func isValidEnum(val string, key string, allowedValues []string) bool {
	for _, allowedVal := range allowedValues {
		if allowedVal == val {
			return true
		}
	}

	return false
}

func populateFromConfOriginServers(dc map[string]interface{}) []OriginServerStruct {
	originServerConf := dc["origin_server"].(*schema.Set)
	var originServerStructs = make([]OriginServerStruct, len(originServerConf.List()))
	var osInd int = 0
	for _, originServer := range originServerConf.List() {
		os := originServer.(map[string]interface{})
		originServerStructs[osInd] = OriginServerStruct{}
		if attr, ok := os["address"]; ok && attr != "" {
			originServerStructs[osInd].Address = attr.(string)
		}
		if attr, ok := os["is_enabled"]; ok && attr != "" {
			originServerStructs[osInd].IsEnabled = attr.(bool)
		}
		if attr, ok := os["is_active"]; ok && !attr.(bool) {
			originServerStructs[osInd].ServerMode = "STANDBY"
		} else {
			originServerStructs[osInd].ServerMode = "ACTIVE"
		}
		if attr, ok := os["weight"]; ok && attr != "" && dc["dc_lb_algorithm"] != nil && dc["dc_lb_algorithm"] == "WEIGHTED" {
			var weight int = attr.(int)
			originServerStructs[osInd].Weight = &weight
		} else {
			originServerStructs[osInd].Weight = nil
		}

		osInd++
	}

	return originServerStructs
}

func populateFromConfDataCenters(d *schema.ResourceData) []DataCenterStruct {
	dataCentersConf := d.Get("data_center").(*schema.Set)
	var dataCentersStructs = make([]DataCenterStruct, len(dataCentersConf.List()))
	var dcInd int = 0
	for _, dataCenter := range dataCentersConf.List() {
		dc := dataCenter.(map[string]interface{})
		dataCentersStructs[dcInd] = DataCenterStruct{}
		if attr, ok := dc["name"]; ok && attr != "" {
			dataCentersStructs[dcInd].Name = attr.(string)
		}
		if attr, ok := dc["dc_id"]; ok && attr != "" && (attr.(int) > 0) {
			dcId := attr.(int)
			dataCentersStructs[dcInd].ID = &dcId
		}
		if attr, ok := dc["ip_mode"]; ok && attr != "" {
			dataCentersStructs[dcInd].IpMode = attr.(string)
		}
		if attr, ok := dc["web_servers_per_server"]; ok && attr != "" && dataCentersStructs[dcInd].IpMode == "SINGLE_IP" {
			var webServersPerServer int = attr.(int)
			dataCentersStructs[dcInd].WebServersPerServer = &webServersPerServer
		}
		if attr, ok := dc["dc_lb_algorithm"]; ok && attr != "" {
			dataCentersStructs[dcInd].DcLbAlgorithm = attr.(string)
		}
		if attr, ok := dc["weight"]; ok && attr != "" && d.Get("site_lb_algorithm") == "WEIGHTED_LB" {
			var weight int = attr.(int)
			dataCentersStructs[dcInd].Weight = &weight
		}
		if attr, ok := dc["is_enabled"]; ok && attr != "" {
			dataCentersStructs[dcInd].IsEnabled = attr.(bool)
		}
		if attr, ok := dc["is_active"]; ok && attr != "" {
			dataCentersStructs[dcInd].IsActive = attr.(bool)
		}
		if attr, ok := dc["is_content"]; ok && attr != "" {
			dataCentersStructs[dcInd].IsContent = attr.(bool)
		}
		if attr, ok := dc["is_rest_of_the_world"]; ok && attr != "" {
			dataCentersStructs[dcInd].IsRestOfTheWorld = attr.(bool)
		}
		if attr, ok := dc["geo_locations"]; ok && attr != "" {
			dataCentersStructs[dcInd].GeoLocations = strings.Split(attr.(string), ",")
		}
		if attr, ok := dc["origin_pop"]; ok && attr != "" {
			dataCentersStructs[dcInd].OriginPoP = attr.(string)
		}

		dataCentersStructs[dcInd].OriginServers = populateFromConfOriginServers(dc)

		dcInd++
	}

	return dataCentersStructs
}

func populateFromConfDataCentersConfigurationDTO(d *schema.ResourceData) DataCentersConfigurationDTO {
	requestDTO := DataCentersConfigurationDTO{}
	requestDTO.Data = make([]DataCentersStruct, 1)
	requestDTO.Data[0].DataCenterMode = d.Get("site_topology").(string)
	requestDTO.Data[0].FailOverRequiredMonitors = d.Get("fail_over_required_monitors").(string)
	requestDTO.Data[0].IsPersistent = d.Get("is_persistent").(bool)
	requestDTO.Data[0].KickStartURL = d.Get("kickstart_url").(string)
	requestDTO.Data[0].KickStartUser = d.Get("kickstart_user").(string)
	requestDTO.Data[0].KickStartPass = d.Get("kickstart_password").(string)
	requestDTO.Data[0].MinAvailableServersForDataCenterUp = d.Get("min_available_servers_for_dc_up").(int)
	requestDTO.Data[0].SiteLbAlgorithm = d.Get("site_lb_algorithm").(string)
	requestDTO.Data[0].DataCenters = populateFromConfDataCenters(d)

	return requestDTO
}

func resourceDataCentersConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	requestDTO := populateFromConfDataCentersConfigurationDTO(d)
	responseDTO, err := client.PutDataCentersConfiguration(d.Get("site_id").(string), requestDTO)
	if err != nil {
		return fmt.Errorf("Error updating Data Centers configuration for site (%s): %s",
			d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 {
		return fmt.Errorf("Error updating Data Centers configuration for site (%s): %s",
			d.Get("site_id"), responseDTO.Errors)
	}

	// Set the dc ID
	d.SetId(d.Get("site_id").(string))

	return resourceDataCentersConfigurationRead(d, m)
}

func resourceDataCentersConfigurationRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListDataCentersResponse for the data center
	client := m.(*Client)

	responseDTO, err := client.GetDataCentersConfiguration(d.Get("site_id").(string))
	if err != nil {
		return fmt.Errorf("Error getting Data Centers configuration for site (%s): %s", d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 {
		if responseDTO.Errors[0].Status == "404" {
			log.Printf("[INFO] Incapsula Site ID %s has already been deleted: %s\n", d.Get("site_id"), responseDTO.Errors)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error getting Data Centers configuration for site (%s): %s", d.Get("site_id"), responseDTO.Errors)
	}

	d.Set("site_lb_algorithm", responseDTO.Data[0].SiteLbAlgorithm)
	d.Set("fail_over_required_monitors", responseDTO.Data[0].FailOverRequiredMonitors)
	d.Set("site_topology", responseDTO.Data[0].DataCenterMode)
	d.Set("min_available_servers_for_dc_up", responseDTO.Data[0].MinAvailableServersForDataCenterUp)
	d.Set("kickstart_url", responseDTO.Data[0].KickStartURL)
	d.Set("kickstart_user", responseDTO.Data[0].KickStartUser)
	d.Set("kickstart_password", responseDTO.Data[0].KickStartPass)
	d.Set("is_persistent", responseDTO.Data[0].IsPersistent)

	dataCenters := &schema.Set{F: resourceDataCentersConfigurationDataCenterHash}
	for _, v := range responseDTO.Data[0].DataCenters {
		dataCenter := map[string]interface{}{}
		dataCenter["name"] = v.Name
		if v.ID != nil {
			dataCenter["dc_id"] = *v.ID
		} else {
			dataCenter["dc_id"] = 0
		}
		dataCenter["ip_mode"] = v.IpMode
		if v.WebServersPerServer != nil {
			dataCenter["web_servers_per_server"] = *v.WebServersPerServer
		} else {
			dataCenter["web_servers_per_server"] = 1
		}
		dataCenter["dc_lb_algorithm"] = v.DcLbAlgorithm
		if v.Weight != nil {
			dataCenter["weight"] = *v.Weight
		} else {
			dataCenter["weight"] = 0
		}
		dataCenter["is_enabled"] = v.IsEnabled
		dataCenter["is_active"] = v.IsActive
		dataCenter["is_content"] = v.IsContent
		dataCenter["is_rest_of_the_world"] = v.IsRestOfTheWorld
		dataCenter["geo_locations"] = strings.Join(v.GeoLocations, ",")
		dataCenter["origin_pop"] = v.OriginPoP

		originServers := &schema.Set{F: resourceDataCentersConfigurationOriginServerHash}
		for _, v := range v.OriginServers {
			originServer := map[string]interface{}{}
			originServer["address"] = v.Address
			if v.Weight != nil {
				originServer["weight"] = *v.Weight
			} else {
				originServer["weight"] = 0
			}
			originServer["is_enabled"] = v.IsEnabled
			originServer["is_active"] = v.ServerMode != "STANDBY"
			originServers.Add(originServer)
		}
		dataCenter["origin_server"] = originServers
		dataCenters.Add(dataCenter)
	}

	d.Set("data_center", dataCenters)

	return nil
}

func resourceDataCentersConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	responseDTO, err := client.GetDataCentersConfiguration(d.Get("site_id").(string))
	if err != nil {
		return fmt.Errorf("Error deleting Data Centers configuration for site (%s): %s", d.Get("site_id"), err)
	}

	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 && responseDTO.Errors[0].Status != "404" {
		return fmt.Errorf("Error deleting Data Centers configuration for site (%s): %s", d.Get("site_id"), responseDTO.Errors)
	}

	d.SetId("")
	return nil
}

func resourceDataCentersConfigurationOriginServerHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["address"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["weight"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	if v, ok := m["is_enabled"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m["is_active"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	return PositiveHash(buf.String())
}

func resourceDataCentersConfigurationDataCenterHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["name"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["ip_mode"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["web_servers_per_server"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	if v, ok := m["dc_lb_algorithm"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["weight"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	if v, ok := m["is_enabled"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m["is_active"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m["is_content"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m["is_rest_of_the_world"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m["geo_locations"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["origin_pop"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	originServerConf := m["origin_server"].(*schema.Set)
	for _, originServer := range originServerConf.List() {
		os := originServer.(map[string]interface{})
		if v, ok := os["address"]; ok {
			buf.WriteString(fmt.Sprintf("%s-", v.(string)))
		}

		if v, ok := os["weight"]; ok {
			buf.WriteString(fmt.Sprintf("%d-", v.(int)))
		}

		if v, ok := os["is_enabled"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}

		if v, ok := os["is_active"]; ok {
			buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
		}
	}

	return PositiveHash(buf.String())
}

func PositiveHash(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
