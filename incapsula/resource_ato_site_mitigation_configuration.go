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
		Read:   resourceATOEndpointMitigationConfigurationRead,
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

func resourceATOEndpointMitigationConfigurationRead(d *schema.ResourceData, m interface{}) error {

	// Fetch our http client
	client := m.(*Client)

	// Extract the required identifiers siteId, accountId and endpointId
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

	// Extract the Account ID
	accountId, ok := d.Get("account_id").(int)

	if !ok {
		return fmt.Errorf("account_id should be of type string. Received : %s", d.Get("account_id"))
	}

	// Extract the EndpointId
	endpointId, ok := d.Get("endpoint_id").(string)

	if !ok {
		return fmt.Errorf("endpoint_id should be of type string. Received : %s", d.Get("endpoint_id"))
	}

	atoEndpointMitigationConfigurationDTO, err := client.GetAtoEndpointMitigationConfigurationWithRetries(accountId, siteId, endpointId)

	// Handle fetch error
	if err != nil {
		return fmt.Errorf("[Error] getting ATO site mitigation configuration for site : %d Error : : %s", siteId, err)
	}

	// Assign ID
	d.SetId(fmt.Sprintf("%d/%d/%s", accountId, siteId, endpointId))

	// Assign the mitigation configuration if present to the terraform compatible map
	if atoEndpointMitigationConfigurationDTO != nil {
		d.Set("site_id", atoEndpointMitigationConfigurationDTO.SiteId)
		d.Set("endpoint_id", atoEndpointMitigationConfigurationDTO.EndpointId)
		d.Set("low_action", atoEndpointMitigationConfigurationDTO.LowAction)
		d.Set("medium_action", atoEndpointMitigationConfigurationDTO.MediumAction)
		d.Set("high_action", atoEndpointMitigationConfigurationDTO.HighAction)
	}

	return nil
}

func resourceATOSiteMitigationConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	// Extract the required identifiers siteId, accountId and endpointId
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

	// Extract the Account ID
	accountId, ok := d.Get("account_id").(int)

	if !ok {
		return fmt.Errorf("account_id should be of type string. Received : %s", d.Get("account_id"))
	}

	// Extract the EndpointId
	endpointId, ok := d.Get("endpoint_id").(string)

	if !ok {
		return fmt.Errorf("endpoint_id should be of type string. Received : %s", d.Get("endpoint_id"))
	}

	// Extract low action
	lowAction, ok := d.Get("low_action").(string)

	if !ok {
		return fmt.Errorf("low_action should be of type string. Received : %s", d.Get("low_action"))
	}

	// Extract medium action
	mediumAction, ok := d.Get("medium_action").(string)

	if !ok {
		return fmt.Errorf("medium_action should be of type string. Received : %s", d.Get("medium_action"))
	}

	// Extract high action
	highAction, ok := d.Get("high_action").(string)

	if !ok {
		return fmt.Errorf("high_action should be of type string. Received : %s", d.Get("high_action"))
	}

	// convert terraform compatible map to ATOEndpointMitigationConfigurationDTO
	var atoMitigationConfigurationDTO ATOEndpointMitigationConfigurationDTO

	atoMitigationConfigurationDTO.SiteId = siteId
	atoMitigationConfigurationDTO.AccountId = accountId
	atoMitigationConfigurationDTO.EndpointId = endpointId
	atoMitigationConfigurationDTO.LowAction = lowAction
	atoMitigationConfigurationDTO.MediumAction = mediumAction
	atoMitigationConfigurationDTO.HighAction = highAction

	// Fetch our http client
	client := m.(*Client)

	err := client.UpdateATOEndpointMitigationConfigurationWithRetries(&atoMitigationConfigurationDTO)
	if err != nil {
		// If endpoints do not exist return the appropriate error
		if strings.Contains(err.Error(), `Endpoints`) && strings.Contains(err.Error(), `do not exist`) {
			return fmt.Errorf("[ERROR] Endpoints do not exist for updating mitigation configuration")
		}
		// Return the error from the api call
		e := fmt.Errorf("[ERROR] Could not update ATO site mitigation configuration for site ID : %d Error : %s \n", atoMitigationConfigurationDTO.SiteId, err)
		return e
	}

	return resourceATOEndpointMitigationConfigurationRead(d, m)
}

func resourceATOSiteMitigationConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	accountId := d.Get("account_id").(int)
	endpointId := d.Get("endpoint_id").(string)

	log.Printf("[DEBUG] Disabling ATO site mitigation for site ID %d \n", siteId)

	err := client.DisableATOEndpointMitigationConfiguration(accountId, siteId, endpointId)
	if err != nil {
		e := fmt.Errorf("[ERROR] Could not disable ATO site mitigation configuration for site ID : %d Error : %s \n", siteId, err)
		return e
	}

	return nil
}
