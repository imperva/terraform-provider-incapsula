package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceApiSecurityApiConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiSecurityAPIConfigCreate,
		Read:   resourceApiSecurityAPIConfigRead,
		Update: resourceApiSecurityAPIConfigUpdate,
		Delete: resourceApiSecurityAPIConfigDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/api_id", d.Id())
				}

				siteID, err := strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				apiID := idSlice[1]
				_, err = strconv.Atoi(apiID)
				if err != nil {
					fmt.Errorf("failed to convert API Id from import command, actual value: %s, expected numeric id", apiID)
				}

				d.Set("site_id", siteID)
				d.SetId(apiID)
				log.Printf("[DEBUG] Import  Api Config JSON for Site ID %d, API Id %s", siteID, apiID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "The site ID which API security is configured on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
			},
			"id": {
				Description: "The API ID which API security is configured.",
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
			},
			"api_specification": {
				Description: "The API specification document content. The supported format is OAS2 or OAS3",
				Type:        schema.TypeString,
				Required:    true,
			},
			//Optional

			"site_name": {
				Description: "The site name",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},

			"validate_host": {
				Description: "When set to true, verifies that the host name and site name match. Set to false in cases such as CNAME reuse or API management integrations where the host name and site name do not match. Default value : true.",
				Type:        schema.TypeBool,
				Optional:    true,
			},

			"invalid_method_violation_action": {
				Description: "The action taken when an invalid method Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALERT_ONLY",
			},

			"missing_param_violation_action": {
				Description: "The action taken when a missing parameter Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALERT_ONLY",
			},

			"invalid_param_value_violation_action": {
				Description: "The action taken when an invalid parameter value Violation occurs. Assigning DEFAULT will inherit the action from parent object, DEFAULT is not applicable for site-level configuration APIs.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALERT_ONLY",
			},

			"invalid_param_name_violation_action": {
				Description: "The violation action taken when invalid request parameter name was sent. Possible values: Alert Only, Block Request, Block User, Block IP, Ignore.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALERT_ONLY",
			},

			"invalid_url_violation_action": {
				Description: "The violation action taken when invalid URL was used. Possible values: Alert Only, Block Request, Block User, Block IP, Ignore.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "ALERT_ONLY",
			},

			"description": {
				Description: "A description that will help recognize the API in the dashboard",
				Type:        schema.TypeString,
				Optional:    true,
			},

			"base_path": {
				Description: "Override the spec basePath / server base path with this value",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},

			"last_modified": {
				Description: "The latest date when the resource was updated",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},

			"host_name": {
				Description: "Host name from the swagger file",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},

			"api_source": {
				Description: "Parameter shows the way API was added. Possible values: USER, AUTO, MIXED",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},

			"discovery_enabled": {
				Description: "Parameter indicates whether automatic API discovery is enabled",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceApiSecurityAPIConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	payload := ApiSecurityApiConfigPostPayload{

		ValidateHost:     false,
		Description:      d.Get("description").(string),
		ApiSpecification: d.Get("api_specification").(string),
		BasePath:         d.Get("base_path").(string),
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        d.Get("invalid_url_violation_action").(string),
			InvalidMethodViolationAction:     d.Get("invalid_method_violation_action").(string),
			MissingParamViolationAction:      d.Get("missing_param_violation_action").(string),
			InvalidParamNameViolationAction:  d.Get("invalid_param_name_violation_action").(string),
			InvalidParamValueViolationAction: d.Get("invalid_param_value_violation_action").(string),
		},
	}

	apiSecurityApiConfigPostResponse, err := client.CreateApiSecurityApiConfig(
		d.Get("site_id").(int),
		&payload)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula API-security site configuration on site id: %d - %s\n", d.Get("site_id"), err)
		return err
	}

	apiID := strconv.Itoa(apiSecurityApiConfigPostResponse.Value.ApiId)
	d.SetId(apiID)
	log.Printf("[INFO] Updated Incapsula API-security api configuration with ID: %s\n", apiID)

	return resourceApiSecurityAPIConfigRead(d, m)
}

func resourceApiSecurityAPIConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     false,
		Description:      d.Get("description").(string),
		ApiSpecification: d.Get("api_specification").(string),
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        d.Get("invalid_url_violation_action").(string),
			InvalidMethodViolationAction:     d.Get("invalid_method_violation_action").(string),
			MissingParamViolationAction:      d.Get("missing_param_violation_action").(string),
			InvalidParamNameViolationAction:  d.Get("invalid_param_name_violation_action").(string),
			InvalidParamValueViolationAction: d.Get("invalid_param_value_violation_action").(string),
		},
	}

	apiSecurityApiConfigPostResponse, err := client.UpdateApiSecurityApiConfig(
		d.Get("site_id").(int),
		d.Id(),
		&payload)

	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula API-security API configuration on site id: %d - %s\n", d.Get("site_id"), err)
		return err
	}

	apiID := strconv.Itoa(apiSecurityApiConfigPostResponse.Value.ApiId)
	d.SetId(apiID)
	log.Printf("[INFO] Updated Incapsula API-security api configuration with ID: %s\n", apiID)

	return resourceApiSecurityAPIConfigRead(d, m)
}

func resourceApiSecurityAPIConfigRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	apiID, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not read API Security API Config ID: %s - %s\n", d.Id(), err)
		return err
	}

	apiSecurityApiConfigGetResponse, err := client.GetApiSecurityApiConfig(siteID, apiID)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API Security API: %d - %s\n", apiID, err)
		return err
	}
	// Set computed values
	d.Set("site_id", apiSecurityApiConfigGetResponse.Value.SiteId)
	d.Set("id", strconv.Itoa(apiSecurityApiConfigGetResponse.Value.Id))
	d.Set("site_name", apiSecurityApiConfigGetResponse.Value.SiteName)
	d.Set("host_name", apiSecurityApiConfigGetResponse.Value.HostName)
	d.Set("base_path", apiSecurityApiConfigGetResponse.Value.BasePath)
	d.Set("description", apiSecurityApiConfigGetResponse.Value.Description)
	d.Set("last_modified", apiSecurityApiConfigGetResponse.Value.LastModified)
	d.Set("api_source", apiSecurityApiConfigGetResponse.Value.ApiSource)
	d.Set("invalid_method_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidMethodViolationAction)
	d.Set("invalid_url_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidUrlViolationAction)
	d.Set("missing_param_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.MissingParamViolationAction)
	d.Set("invalid_param_name_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidParamNameViolationAction)
	d.Set("invalid_param_value_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidParamValueViolationAction)
	//In current implementation validateHost value is always set as "false". Will be changed in next releases
	d.Set("validate_host", false)

	apiSecurityApiConfigGetFileResponse, err := client.GetApiSecurityApiSwaggerConfig(siteID, apiID)
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API Security API swagger file: %d - %s\n", apiID, err)
		return err
	}
	d.Set("api_specification", apiSecurityApiConfigGetFileResponse.Value)

	return nil
}

func resourceApiSecurityAPIConfigDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	err := client.DeleteApiSecurityApiConfig(siteID, d.Id())

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
