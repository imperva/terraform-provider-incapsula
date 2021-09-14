package incapsula

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiSecuritySiteConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiSecuritySiteConfigCreate,
		Read:   resourceApiSecuritySiteConfigRead,
		Update: resourceApiSecuritySiteConfigUpdate,
		Delete: resourceApiSecuritySiteConfigDelete,
		Importer: &schema.ResourceImporter{ //todo - check
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID, err := strconv.Atoi(d.Id())
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", d.Id())
				}

				d.Set("site_id", siteID)
				log.Printf("[DEBUG] Import  Site Config JSON for Site ID %d", siteID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "The Site ID of the the site the API security is configured on.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			//optional
			"account_id": {
				Description: "The Account ID of the the site the API security is configured on.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"site_name": {
				Description: "The site name",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"api_only_site": {
				Description: "", //todo api_only_site set description
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"discovery_enabled": {
				Description: "Parameter shows whether automatic API discovery is enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"non_api_request_violation_action": {
				Description: "", //todo add description
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_url_violation_action": {
				Description: "The action taken when an invalid URL Violation occurs. Actions available: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_method_violation_action": {
				Description: "The action taken when an invalid method Violation occurs. Actions available: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"missing_param_violation_action": {
				Description: "The action taken when a missing parameter Violation occurs. Actions available: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_param_value_violation_action": {
				Description: "The action taken when an invalid parameter value Violation occurs. Actions available: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_param_name_violation_action": {
				Description: "The action taken when an invalid parameter value Violation occurs. Actions available: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"last_modified": {
				Description: "The latest date when the resource was updated",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"is_automatic_discovery_api_integration_enabled": {
				Description: "Parameter shows whether automatic API discovery integration is enabled",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

func resourceApiSecuritySiteConfigCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Create Incapsula API-security site configuration for site ID: %d redircted to update function\n", d.Get("site_id"))
	return resourceApiSecuritySiteConfigUpdate(d, m)
}

func resourceApiSecuritySiteConfigUpdate(
	d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	payload := ApiSecuritySiteConfigPostPayload{
		ApiOnlySite:                  d.Get("api_only_site").(bool),
		NonApiRequestViolationAction: d.Get("non_api_request_violation_action").(string),
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        d.Get("invalid_url_violation_action").(string),
			InvalidMethodViolationAction:     d.Get("invalid_method_violation_action").(string),
			MissingParamViolationAction:      d.Get("missing_param_violation_action").(string),
			InvalidParamNameViolationAction:  d.Get("invalid_param_name_violation_action").(string),
			InvalidParamValueViolationAction: d.Get("invalid_param_value_violation_action").(string),
		},
	}

	apiSecuritySiteConfigPostResponse, err := client.UpdateApiSecuritySiteConfig(
		d.Get("site_id").(int),
		&payload)

	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula API-security Site Configuration on site id: %d - %s\n", d.Get("site_id"), err)
		return err
	}

	siteID := strconv.Itoa(apiSecuritySiteConfigPostResponse.Value.SiteId)
	d.SetId(siteID)
	log.Printf("[INFO] Updated Incapsula API-security site configuration with ID: %s\n", siteID)

	return resourceApiSecuritySiteConfigRead(d, m)
}

func resourceApiSecuritySiteConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id")

	apiSecuritySiteConfigGetResponse, err := client.ReadApiSecuritySiteConfig(siteId.(int))
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API-security site configuration for site ID: %d - %s\n", siteId, err)
		return err
	}

	// Set computed values
	d.Set("site_name", apiSecuritySiteConfigGetResponse.Value.SiteName)
	d.Set("account_id", apiSecuritySiteConfigGetResponse.Value.AccountId)
	d.Set("api_only_site", apiSecuritySiteConfigGetResponse.Value.ApiOnlySite)
	d.Set("invalid_method_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidMethodViolationAction)
	d.Set("invalid_param_name_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidParamNameViolationAction)
	d.Set("invalid_param_value_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidParamValueViolationAction)
	d.Set("invalid_url_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidUrlViolationAction)
	d.Set("missing_param_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.MissingParamViolationAction)
	d.Set("is_automatic_discovery_api_integration_enabled", apiSecuritySiteConfigGetResponse.Value.IsAutomaticDiscoveryApiIntegrationEnabled)
	d.Set("non_api_request_violation_action", apiSecuritySiteConfigGetResponse.Value.NonApiRequestViolationAction)
	d.Set("discovery_enabled", apiSecuritySiteConfigGetResponse.Value.DiscoveryEnabled)
	return nil
}

func resourceApiSecuritySiteConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[ERROR] Deleting Incapsula API-security site configuration isn't supported. request made for site ID: %d \n", d.Get("site_id"))
	return nil
}
