package incapsula

import (
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: "The site name.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"api_only_site": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"discovery_enabled": {
				Description: "",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"non_api_request_violation_action": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
				//RequiredWith:    []string{"api_only_site"},
			},
			"invalid_url_violation_action": {
				Description: "The action taken when an invalid URL Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_method_violation_action": {
				Description: "The action taken when an invalid method Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"missing_param_violation_action": {
				Description: "The action taken when a missing parameter Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_param_value_violation_action": {
				Description: "The action taken when an invalid parameter value Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"invalid_param_name_violation_action": {
				Description: "The action taken when an invalid parameter value Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"last_modified": {
				Description: "", //todo add description
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"is_automatic_discovery_api_integration_enabled": {
				Description: "",
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
		ApiOnlySite: d.Get("api_only_site").(bool),
		IsAutomaticDiscoveryApiIntegrationEnabled: d.Get("is_automatic_discovery_api_integration_enabled").(bool),
		NonApiRequestViolationAction:              d.Get("non_api_request_violation_action").(string),
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
		log.Printf("[ERROR] Could not update Incapsula API-security site configuration on site id: %d - %s\n", d.Get("site_id"), err)
		return err
	}

	siteID := strconv.Itoa(apiSecuritySiteConfigPostResponse.Value.SiteId)
	d.SetId(siteID)
	log.Printf("[INFO] Updated Incapsula API-security site configuration with ID: %s\n", siteID)

	return resourceApiSecuritySiteConfigRead(d, m)
}

func resourceApiSecuritySiteConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id")

	apiSecuritySiteConfigGetResponse, err := client.ReadApiSecuritySiteConfig(siteID.(int))
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API-security site configuration for site ID: %s - %s\n", siteID, err)
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
	return nil
}

func resourceApiSecuritySiteConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[ERROR] Deleting Incapsula API-security site configuration isn't supported. request made for site ID: %d \n", d.Get("site_id"))
	return nil
}
