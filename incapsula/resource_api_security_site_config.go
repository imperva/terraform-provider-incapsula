package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
)

func resourceApiSecuritySiteConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiSecuritySiteConfigUpdate,
		Read:   resourceApiSecuritySiteConfigRead,
		Update: resourceApiSecuritySiteConfigUpdate,
		Delete: resourceApiSecuritySiteConfigDelete,
		Importer: &schema.ResourceImporter{
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
			"is_automatic_discovery_api_integration_enabled": {
				Description: "Parameter shows whether automatic API discovery is enabled",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			//Optional
			"invalid_url_violation_action": {
				Description:  "The action taken when an invalid URL Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE"}, false),
				Default:      "ALERT_ONLY",
			},
			"invalid_method_violation_action": {
				Description:  "The action taken when an invalid method Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE"}, false),
				Default:      "ALERT_ONLY",
			},
			"missing_param_violation_action": {
				Description:  "The action taken when a missing parameter Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE"}, false),
				Default:      "ALERT_ONLY",
			},
			"invalid_param_value_violation_action": {
				Description:  "The action taken when an invalid parameter value Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE"}, false),
				Default:      "ALERT_ONLY",
			},
			"invalid_param_name_violation_action": {
				Description:  "The action taken when an invalid parameter value Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE"}, false),
				Default:      "ALERT_ONLY",
				Deprecated:   "invalid_param_name_violation_action field is deprecated",
			},
			"is_api_only_site": {
				Description: "Apply positive security model for all traffic on the site. Applying the positive security model for all traffic on the site may lead to undesired request blocking.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"non_api_request_violation_action": {
				Description:  "Action to be taken for traffic on the site that does not target the uploaded APIs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE"}, false),
			},
			"last_modified": {
				Description: "The last modified timestamp",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceApiSecuritySiteConfigUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Update Incapsula API-security site configuration for site ID: %d", d.Get("site_id"))

	client := m.(*Client)
	payload := ApiSecuritySiteConfigPostPayload{
		ApiOnlySite:                               d.Get("is_api_only_site").(bool),
		NonApiRequestViolationAction:              d.Get("non_api_request_violation_action").(string),
		IsAutomaticDiscoveryApiIntegrationEnabled: d.Get("is_automatic_discovery_api_integration_enabled").(bool),
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
	d.Set("invalid_method_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidMethodViolationAction)
	d.Set("invalid_param_name_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidParamNameViolationAction)
	d.Set("invalid_param_value_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidParamValueViolationAction)
	d.Set("invalid_url_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.InvalidUrlViolationAction)
	d.Set("missing_param_violation_action", apiSecuritySiteConfigGetResponse.Value.ViolationActions.MissingParamViolationAction)
	d.Set("non_api_request_violation_action", apiSecuritySiteConfigGetResponse.Value.NonApiRequestViolationAction)
	d.Set("is_automatic_discovery_api_integration_enabled", apiSecuritySiteConfigGetResponse.Value.IsAutomaticDiscoveryApiIntegrationEnabled)
	d.Set("is_api_only_site", apiSecuritySiteConfigGetResponse.Value.ApiOnlySite)
	return nil
}

func resourceApiSecuritySiteConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[ERROR] Deleting Incapsula API-security site configuration isn't supported. request made for site ID: %d \n", d.Get("site_id"))
	d.SetId("")
	return nil
}
