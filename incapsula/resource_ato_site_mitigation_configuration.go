package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceATOSiteMitigationConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceATOSiteMitigationConfigurationUpdate,
		Read:   resourceATOSiteMitigationConfigurationRead,
		Update: resourceATOSiteMitigationConfigurationUpdate,
		Delete: resourceATOSiteMitigationConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteId, err := strconv.Atoi(d.Id())
				err = d.Set("site_id", siteId)
				if err != nil {
					return nil, fmt.Errorf("failed to extract site ID from import command, actual value: %s, error : %s", d.Id(), err)
				}
				log.Printf("[DEBUG] Import ATO allowlist for site ID %d", siteId)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Account ID that the site belongs to.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"site_id": {
				Description: "Site ID to get the allowlist for.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"mitigation_configuration": {
				Description: "The mitigation configuration of IPs and IP ranges for the given site ID",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
		},
	}
}

func resourceATOSiteMitigationConfigurationRead(d *schema.ResourceData, m interface{}) error {

	// Fetch our http client
	client := m.(*Client)

	// Fetch the ATO allowlist of IPs and subnets
	siteId, ok := d.Get("site_id").(int)

	if !ok {
		return fmt.Errorf("site_id should be of type int. Received : %s", d.Get("site_id"))
	}

	if siteId == 0 {
		siteIdFromResourceId, conversionError := strconv.Atoi(d.Id())
		if conversionError != nil {
			return fmt.Errorf("atleast one of id or site_id should be set for incapsula_ato_site_allowlist")
		}
		siteId = siteIdFromResourceId
	}

	accountId := d.Get("account_id").(int)
	atoSiteMitigationConfigurationDTO, err := client.GetAtoSiteMitigationConfigurationWithRetries(accountId, siteId)

	// Handle fetch error
	if err != nil {
		return fmt.Errorf("[Error] getting ATO site mitigation configuration for site : %d Error : : %s", siteId, err)
	}

	// Initialize the data map that terraform uses
	var mitigationConfigurationMap = make(map[string]interface{})

	// Assign values from our received DTO to the map that terraform understands.
	mitigationConfigurationMap["site_id"] = atoSiteMitigationConfigurationDTO.SiteId

	// Assign the allowlist if present to the terraform compatible map
	if atoSiteMitigationConfigurationDTO.MitigationConfiguration != nil {

		mitigationConfigurationMap["mitigation_configuration"] = make([]map[string]interface{}, len(atoSiteMitigationConfigurationDTO.MitigationConfiguration))

		for i, mitigationConfigurationForEndpoint := range atoSiteMitigationConfigurationDTO.MitigationConfiguration {

			// Assign the properties to the map
			mitigationConfigurationForEndpointMap := make(map[string]interface{})
			mitigationConfigurationForEndpointMap["endpoint_id"] = mitigationConfigurationForEndpoint.EndpointId
			mitigationConfigurationForEndpointMap["low_action"] = mitigationConfigurationForEndpoint.LowAction
			mitigationConfigurationForEndpointMap["medium_action"] = mitigationConfigurationForEndpoint.MediumAction
			mitigationConfigurationForEndpointMap["high_action"] = mitigationConfigurationForEndpoint.HighAction

			// Assign the mitigation for endpoint to the site mitigation configuration
			mitigationConfigurationMap["mitigation_configuration"].([]map[string]interface{})[i] = mitigationConfigurationForEndpointMap
		}

	} else {
		mitigationConfigurationMap["mitigation_configuration"] = make([]interface{}, 0)
	}

	d.SetId(strconv.Itoa(siteId))
	err = d.Set("mitigation_configuration", mitigationConfigurationMap["mitigation_configuration"])
	if err != nil {
		e := fmt.Errorf("[Error] Error in reading mitigation configuration values : %s", err)
		return e
	}

	return nil
}

func resourceATOSiteMitigationConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	accountId := d.Get("account_id").(int)

	atoMitigationConfigurationMap := make(map[string]interface{})
	atoMitigationConfigurationMap["account_id"] = accountId
	atoMitigationConfigurationMap["site_id"] = siteId
	atoMitigationConfigurationMap["mitigation_configuration"] = d.Get("mitigation_configuration").([]interface{})

	log.Printf("[DEBUG] Updating ATO site mitigation configuration site ID %d \n", siteId)

	// convert terraform compatible map to ATOSiteMitigationConfigurationDTO
	atoMitigationConfigurationDTO, err := formAtoMitigationConfigurationDTOFromMap(atoMitigationConfigurationMap)
	if err != nil {
		e := fmt.Errorf("[Error] Error forming ATO mitigation configuration object for API call : %s", err)
		log.Printf(e.Error())
		return err
	}

	err = client.UpdateATOSiteMitigationConfigurationWithRetries(atoMitigationConfigurationDTO)
	if err != nil {
		// If endpoints do not exist return the appropriate error
		if strings.Contains(err.Error(), `Endpoints`) && strings.Contains(err.Error(), `do not exist`) {
			return fmt.Errorf("[ERROR] Endpoints do not exist for updating mitigation configuration")
		}
		// Return the error from the api call
		e := fmt.Errorf("[ERROR] Could not update ATO site mitigation configuration for site ID : %d Error : %s \n", atoMitigationConfigurationDTO.SiteId, err)
		return e
	}

	return resourceATOSiteMitigationConfigurationRead(d, m)
}

func resourceATOSiteMitigationConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	accountId := d.Get("account_id").(int)

	log.Printf("[DEBUG] Disabling ATO site mitigation for site ID %d \n", siteId)

	err := client.DisableATOSiteMitigationConfiguration(accountId, siteId)
	if err != nil {
		e := fmt.Errorf("[ERROR] Could not disable ATO site mitigation configuration for site ID : %d Error : %s \n", siteId, err)
		return e
	}

	return nil
}
