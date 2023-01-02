package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Description: "Numeric identifier of the site to operate on. ",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"api_specification": {
				Description: "The API specification document content. The supported format is OAS2 or OAS3",
				Type:        schema.TypeString,
				Required:    true,
			},

			//Optional
			"invalid_url_violation_action": {
				Description:  "The violation action taken when invalid URL was used. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},
			"invalid_method_violation_action": {
				Description:  "The action taken when an invalid method Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},
			"missing_param_violation_action": {
				Description:  "The action taken when a missing parameter Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},
			"invalid_param_value_violation_action": {
				Description:  "The action taken when an invalid parameter value Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},
			"invalid_param_name_violation_action": {
				Description:  "The violation action taken when invalid request parameter name was sent. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object",
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
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
			"host_name": {
				Description: "The host name from the swagger file",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified": {
				Description: "The last modified timestamp",
				Type:        schema.TypeInt,
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

	_, err := client.UpdateApiSecurityApiConfig(
		d.Get("site_id").(int),
		d.Id(),
		&payload)

	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula API-security API configuration on site id: %d - %s\n", d.Get("site_id"), err)
		return err
	}

	log.Printf("[INFO] Updated Incapsula API-security api configuration with ID: %s\n", d.Id())
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
	d.SetId(strconv.Itoa(apiSecurityApiConfigGetResponse.Value.Id))
	d.Set("site_id", apiSecurityApiConfigGetResponse.Value.SiteId)
	d.Set("host_name", apiSecurityApiConfigGetResponse.Value.HostName)
	d.Set("base_path", apiSecurityApiConfigGetResponse.Value.BasePath)
	d.Set("description", apiSecurityApiConfigGetResponse.Value.Description)
	d.Set("last_modified", apiSecurityApiConfigGetResponse.Value.LastModified)
	d.Set("invalid_method_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidMethodViolationAction)
	d.Set("invalid_url_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidUrlViolationAction)
	d.Set("missing_param_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.MissingParamViolationAction)
	d.Set("invalid_param_name_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidParamNameViolationAction)
	d.Set("invalid_param_value_violation_action", apiSecurityApiConfigGetResponse.Value.ViolationActions.InvalidParamValueViolationAction)

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
	apiID, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Error converting Api Security API configuration ID for site (%d). Expected numeric value, got %s", siteID, d.Id())
	}

	err = client.DeleteApiSecurityApiConfig(siteID, d.Id())

	if err != nil {
		return fmt.Errorf("Error deleting Api Security API configuration for site (%d), API Id (%d): %s", siteID, apiID, err)
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
