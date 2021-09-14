package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				if len(idSlice) != 3 || idSlice[0] == "" || idSlice[1] == "" || idSlice[2] == "" {
					return nil, fmt.Errorf("unexpected format of API Security Endpoint ID (%q), expected api_id/method/path", d.Id())
				}

				apiId, err := strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert API Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				method := idSlice[1]
				path := idSlice[2]

				d.Set("api_id", apiId)
				d.Set("method", method)
				d.Set("path", strings.ReplaceAll(path, "_", "/"))
				log.Printf("[DEBUG] Import API Security Endpoint Config JSON for API ID %d, Path %s, Method  %s", apiId, path, method)

				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"api_id": {
				Description: "The site ID which API security is configured on.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			"missing_param_violation_action": {
				Description:  "The action taken when an invalid URL Violation occurs. Assigning DEFAULT will inherit the action from parent object.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"invalid_param_value_violation_action", "invalid_param_name_violation_action"},
				Default:      "DEFAULT",
			},

			"invalid_param_value_violation_action": {
				Description:  "The action taken when an invalid parameter value Violation occurs. Assigning DEFAULT will inherit the action from parent object. Possible values: Alert Only, Block Request, Block User, Block IP, Ignore.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"missing_param_violation_action", "invalid_param_name_violation_action"},
				Default:      "DEFAULT",
			},

			"invalid_param_name_violation_action": {
				Description:  "The action taken when an invalid parameter name Violation occurs. Assigning DEFAULT will inherit the action from parent object. Possible values: Alert Only, Block Request, Block User, Block IP, Ignore.",
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"missing_param_violation_action", "invalid_param_value_violation_action"},
				Default:      "DEFAULT",
			},

			"method": {
				Description: "HTTP method that describes a specific endpoint",
				Type:        schema.TypeString,
				Required:    true,
			},

			"path": {
				Description: "An URL path of specific endpoint ",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceApiSecurityEndpointConfigCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Create Incapsula API-security endpoint configuration for API ID: %d redircted to update function\n", d.Get("api_id"))
	return resourceApiSecurityEndpointConfigUpdate(d, m)
}

func resourceApiSecurityEndpointConfigRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Read Incapsula API-security endpoint configuration for ID: %s", d.Id())
	client := m.(*Client)
	endpointGetAllResponse, err := client.GetApiSecurityAllEndpointsConfig(d.Get("api_id").(int))
	var currentId int
	var found bool

	for _, entry := range endpointGetAllResponse.Value {
		if entry.Path == d.Get("path").(string) && entry.Method == d.Get("method").(string) {
			currentId = entry.Id
			found = true
			break
		}
	}

	if !found {
		log.Printf("[INFO] API-security endpoint [%s %s] doesn't exist and will not be updated.", d.Get("method").(string), d.Get("path").(string))
		d.SetId("")
		return nil
	}
	log.Printf("found endpoint id %d", currentId)

	endpointID := strconv.Itoa(currentId)

	endpointGetResponse, err := client.GetApiSecurityEndpointConfig(d.Get("api_id").(int), endpointID)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API-security endpoint: %s - %s\n", endpointID, err)
		return err
	}

	// Set computed values
	d.Set("missing_param_violation_action", endpointGetResponse.Value.ViolationActions.MissingParamViolationAction)
	d.Set("invalid_param_name_violation_action", endpointGetResponse.Value.ViolationActions.InvalidParamNameViolationAction)
	d.Set("invalid_param_value_violation_action", endpointGetResponse.Value.ViolationActions.InvalidParamValueViolationAction)
	d.Set("method", endpointGetResponse.Value.Method)
	d.Set("path", endpointGetResponse.Value.Path)
	d.SetId(strconv.Itoa(currentId))

	return nil
}

func resourceApiSecurityEndpointConfigUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	//get all Endpoints
	endpointGetAllResponse, err := client.GetApiSecurityAllEndpointsConfig(d.Get("api_id").(int))

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula API-security all endpoints: %s - %s\n", d.Get("api_id"), err)
		return err
	}

	//find current endpoint by path+method
	var currentId int
	var found bool
	for _, entry := range endpointGetAllResponse.Value {
		if entry.Path == d.Get("path").(string) && entry.Method == d.Get("method").(string) {
			currentId = entry.Id
			found = true
			break
		}
	}

	if !found {
		log.Printf("[INFO] API-security endpoint [%s %s] doesn't exist and will not be updated.", d.Get("method").(string), d.Get("path").(string))
		d.SetId("")
		return fmt.Errorf("Endpoint with method %s, path %s, wasn't found", d.Get("method"), d.Get("path"))
	}

	payload := ApiSecurityEndpointConfigPostPayload{
		ViolationActions: UserViolationActions{
			MissingParamViolationAction:      d.Get("missing_param_violation_action").(string),
			InvalidParamNameViolationAction:  d.Get("invalid_param_name_violation_action").(string),
			InvalidParamValueViolationAction: d.Get("invalid_param_value_violation_action").(string),
		},
	}

	_, err = client.PostApiSecurityEndpointConfig(d.Get("api_id").(int), currentId, &payload)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(currentId))

	return nil
}

func resourceApiSecurityEndpointConfigDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[INFO] Deleting Incapsula API-security endpoint configuration isn't supported. Request made for endpoint ID: %s \n", d.Id())
	return nil
}
