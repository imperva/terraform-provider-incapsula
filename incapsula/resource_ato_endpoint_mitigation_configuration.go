package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func ATOEndpointMitigationConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: ATOEndpointMitigationConfigurationUpdate,
		Read:   resourceATOEndpointMitigationConfigurationRead,
		Update: ATOEndpointMitigationConfigurationUpdate,
		Delete: ATOEndpointMitigationConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				keyParts := strings.Split(d.Id(), "/")
				print("id is %s", d.Id())
				if len(keyParts) != 3 {
					return nil, fmt.Errorf("Error parsing ID, actual value: %s, expected two numeric IDs and string seperated by '/'\n", d.Id())
				}
				accountId, err := strconv.Atoi(keyParts[0])
				if err != nil {
					return nil, fmt.Errorf("failed to convert account ID from import command, actual value: %s, expected numeric id", keyParts[0])
				}
				siteId, err := strconv.Atoi(keyParts[1])
				if err != nil {
					return nil, fmt.Errorf("failed to convert site ID from import command, actual value: %s, expected numeric id", keyParts[1])
				}
				endpointId := keyParts[1]

				d.Set("account_id", accountId)
				d.Set("site_id", siteId)
				d.Set("endpoint_id", endpointId)
				log.Printf("[DEBUG] To Import ATO endpoint mitigation for account ID %d , site ID %d , endpoint ID %s", accountId, siteId, endpointId)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Account ID that the site belongs to.",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"site_id": {
				Description: "Site ID to get the allowlist for.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"endpoint_id": {
				Description: "Endpoint ID associated with this request",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"mitigation_action_for_low_risk": {
				Description: "Mitigation action configured for low risk requests - in UPPER CASE.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"mitigation_action_for_medium_risk": {
				Description: "Mitigation action configured for medium risk requests - in UPPER CASE.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"mitigation_action_for_high_risk": {
				Description: "Mitigation action configured for high risk requests - in UPPER CASE.",
				Type:        schema.TypeString,
				Required:    true,
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
		siteIdFromResourceId, conversionError := strconv.Atoi(strings.Split(d.Id(), "/")[1])
		if conversionError != nil {
			return fmt.Errorf("either id : %s or (account_id, site_id and endpoint_id) should be set for incapsula_ato_endpoint_mitigation_configuration for reading mitigation configuration", d.Id())
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

	atoEndpointMitigationConfigurationDTO, status, err := client.GetAtoEndpointMitigationConfigurationWithRetries(accountId, siteId, endpointId)

	// Handle fetch error
	if err != nil {
		return fmt.Errorf("[Error] getting ATO site mitigation configuration for site : %d Error : : %s", siteId, err)
	}

	// Handle site does not exist in ATO. If this is a permissions issue, then let this be addressed in the update phase
	if status == http.StatusUnauthorized {
		// Remove this resource from the state file by setting empty ID as it does not exist. Terraform will remove it
		d.Set("id", "")
	} else {
		// Assign ID
		d.SetId(fmt.Sprintf("%d/%d/%s", accountId, siteId, endpointId))
	}

	// Assign the mitigation configuration if present to the terraform compatible map
	if atoEndpointMitigationConfigurationDTO != nil {
		d.Set("site_id", atoEndpointMitigationConfigurationDTO.SiteId)
		d.Set("endpoint_id", atoEndpointMitigationConfigurationDTO.EndpointId)
		d.Set("mitigation_action_for_low_risk", atoEndpointMitigationConfigurationDTO.LowAction)
		d.Set("mitigation_action_for_medium_risk", atoEndpointMitigationConfigurationDTO.MediumAction)
		d.Set("mitigation_action_for_high_risk", atoEndpointMitigationConfigurationDTO.HighAction)
	} else {
		// Setting empty ID as it does not exist. Terraform will remove it
		d.Set("id", "")
	}

	return nil
}

func ATOEndpointMitigationConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	// Extract the required identifiers siteId, accountId and endpointId
	siteId, ok := d.Get("site_id").(int)

	if !ok {
		return fmt.Errorf("site_id should be of type int. Received : %s", d.Get("site_id"))
	}

	if siteId == 0 {
		siteIdFromResourceId, conversionError := strconv.Atoi(d.Id())
		if conversionError != nil {
			return fmt.Errorf("either id or (account_id, site_id and endpoint_id) should be set for incapsula_ato_endpoint_mitigation_configuration for updating mitigation configuration")
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
	lowAction, ok := d.Get("mitigation_action_for_low_risk").(string)

	if !ok {
		return fmt.Errorf("mitigation_action_for_low_risk should be of type string. Received : %s", d.Get("mitigation_action_for_low_risk"))
	}

	// Extract medium action
	mediumAction, ok := d.Get("mitigation_action_for_medium_risk").(string)

	if !ok {
		return fmt.Errorf("mitigation_action_for_medium_risk should be of type string. Received : %s", d.Get("mitigation_action_for_medium_risk"))
	}

	// Extract high action
	highAction, ok := d.Get("mitigation_action_for_high_risk").(string)

	if !ok {
		return fmt.Errorf("mitigation_action_for_high_risk should be of type string. Received : %s", d.Get("mitigation_action_for_high_risk"))
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
		// Return the error from the api call
		e := fmt.Errorf("[ERROR] Could not update ATO site mitigation configuration for site ID : %d Error : %s \n", atoMitigationConfigurationDTO.SiteId, err)
		return e
	}

	return resourceATOEndpointMitigationConfigurationRead(d, m)
}

func ATOEndpointMitigationConfigurationDelete(d *schema.ResourceData, m interface{}) error {
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
