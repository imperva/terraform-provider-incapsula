package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
)

func resourceApiSecurityEndpointConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiSecurityEndpointConfigCreate,
		Read:   resourceApiSecurityEndpointConfigRead,
		Update: resourceApiSecurityEndpointConfigUpdate,
		Delete: resourceApiSecurityEndpointConfigDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of API Security Endpoint ID (%q), expected api_id/endpoint_id", d.Id())
				}

				apiId, err := strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert API Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				endpointId, err := strconv.Atoi(idSlice[1])
				if err != nil {
					fmt.Errorf("failed to convert Endpoint Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				d.Set("api_id", apiId)
				d.SetId(idSlice[1])
				log.Printf("[DEBUG] Import API Security Endpoint Config JSON for API ID %d, Endpoint ID %d", apiId, endpointId)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"api_id": {
				Description: "The site ID which API security is configured on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"method": {
				Description: "HTTP method that describes a specific endpoint. Possible values: POST, GET, PUT, PATCH, DELETE, HEAD, OPTIONS",
				Type:        schema.TypeString,
				Required:    true,
			},
			"path": {
				Description: "An URL path of specific endpoint ",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"missing_param_violation_action": {
				Description:  "The action taken when an invalid URL Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},

			"invalid_param_value_violation_action": {
				Description:  "The action taken when an invalid parameter value Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},

			"invalid_param_name_violation_action": {
				Description:  "The action taken when an invalid parameter name Violation occurs. Possible values: ALERT_ONLY, BLOCK_REQUEST, BLOCK_USER, BLOCK_IP, IGNORE, DEFAULT. Assigning DEFAULT will inherit the action from parent object.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DEFAULT",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_ONLY", "BLOCK_REQUEST", "BLOCK_USER", "BLOCK_IP", "IGNORE", "DEFAULT"}, false),
			},
		},
	}
}

func resourceApiSecurityEndpointConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Read Incapsula API-security endpoint configuration for ID: %s", d.Id())
	client := m.(*Client)
	endpointGetResponse, err := client.GetApiSecurityEndpointConfig(d.Get("api_id").(int), d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API-security endpoint: %s - %s\n", d.Get("id"), err)
		return err
	}

	d.Set("missing_param_violation_action", endpointGetResponse.Value.ViolationActions.MissingParamViolationAction)
	d.Set("invalid_param_name_violation_action", endpointGetResponse.Value.ViolationActions.InvalidParamNameViolationAction)
	d.Set("invalid_param_value_violation_action", endpointGetResponse.Value.ViolationActions.InvalidParamValueViolationAction)
	d.Set("method", endpointGetResponse.Value.Method)
	d.Set("path", endpointGetResponse.Value.Path)
	return nil
}

func resourceApiSecurityEndpointConfigCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	endpointGetAllResponse, _ := client.GetApiSecurityAllEndpointsConfig(d.Get("api_id").(int))
	var found bool
	var endpointId string
	for _, entry := range endpointGetAllResponse.Value {
		if entry.Path == d.Get("path").(string) && entry.Method == d.Get("method").(string) {
			endpointId = strconv.Itoa(entry.Id)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("[ERROR] API-security endpoint [%s %s] doesn't exist and will not be updated.", d.Get("method").(string), d.Get("path").(string))
	}
	log.Printf("[DEBUG] found endpoint id %s", endpointId)
	d.SetId(endpointId)

	return resourceApiSecurityEndpointConfigUpdate(d, m)
}

func resourceApiSecurityEndpointConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	payload := ApiSecurityEndpointConfigPostPayload{
		ViolationActions: UserViolationActions{
			MissingParamViolationAction:      d.Get("missing_param_violation_action").(string),
			InvalidParamNameViolationAction:  d.Get("invalid_param_name_violation_action").(string),
			InvalidParamValueViolationAction: d.Get("invalid_param_value_violation_action").(string),
		},
	}

	endpointId, err := strconv.Atoi(d.Id())
	if err != nil {
		fmt.Errorf("Endpoint ID should be numeric. Actual value: %s", d.Id())
	}
	_, err = client.PostApiSecurityEndpointConfig(d.Get("api_id").(int), endpointId, &payload)

	if err != nil {
		return err
	}

	return resourceApiSecurityEndpointConfigRead(d, m)
}

func resourceApiSecurityEndpointConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Deleting Incapsula API-security endpoint configuration isn't supported. Request made for endpoint ID: %s \n", d.Id())
	return nil
}
