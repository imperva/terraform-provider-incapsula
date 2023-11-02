package incapsula

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
)

func resourceIncapRulePriority() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIncapRulePriorityRead,
		//Update: resourceIncapRulePriorityUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/rule_id", d.Id())
				}

				siteID := idSlice[0]
				d.Set("site_id", siteID)

				ruleID := idSlice[1]
				d.SetId(ruleID)

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"rule_name": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "Rule action. See the detailed descriptions in the API documentation",
				Type:        schema.TypeString,
				Required:    true,
			},
			"category": {
				Description:  "Rule action. See the detailed descriptions in the API documentation",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"REDIRECT", "SIMPLIFIED_ REDIRECT", "REWRITE", "FORWARD", "REWRITE_RESPONSE"}, false),
			},
			// Optional Arguments
			"filter": {
				Description: "Defines the conditions that trigger the rule action",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"response_code": {
				Description: "For `RULE_ACTION_REDIRECT` or `RULE_ACTION_SIMPLIFIED_REDIRECT` rule's response code, valid values are `302`, `301`, `303`, `307`, `308`. For `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` rule's response code, valid values are all 3-digits numbers. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"add_if_missing": {
				Description: "Add cookie or header if it doesn't exist (Rewrite cookie rule only).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"rewrite_existing": {
				Description: "Rewrite cookie or header if it exists.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"from": {
				Description: "Pattern to rewrite. For `RULE_ACTION_REWRITE_URL` - Url to rewrite. For `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to rewrite. For `RULE_ACTION_REWRITE_COOKIE` - Cookie value to rewrite.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"to": {
				Description: "Pattern to change to. `RULE_ACTION_REWRITE_URL` - Url to change to. `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to change to. `RULE_ACTION_REWRITE_COOKIE` - Cookie value to change to.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"rewrite_name": {
				Description: "Name of cookie or header to rewrite. Applies only for `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cookie_name": {
				Description: "Name of cookie to rewrite,delete,",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"header_name": {
				Description: "Name of header to REWRITE,REWRITE_RESPONSE",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dc_id": {
				Description: "Data center ID to forward the request to",
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
			},
			"port_forwarding_context": {
				Description: "Context for port forwarding. Use Port Value or Use Header Name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"port_forwarding_value": {
				Description: "Port number or header name for port forwarding",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_type": {
				Description: "The error that triggers the rule. <code>error.type.all</code> triggers the rule regardless of the error type.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_response_format": {
				Description:  "The format of the given error response in the error_response_data field.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"json", "xml"}, false),
			},
			"error_response_data": {
				Description: "",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"multiple_headers_deletion": {
				Description: "Delete multiple header occurrences. Applies only to rules using `RULE_ACTION_DELETE_HEADER` and `RULE_ACTION_RESPONSE_DELETE_HEADER`.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"enabled": {
				Description: "Enable or disable rule.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceIncapRulePriorityRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	_, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] The resource ID should be numeric")
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Detail:   fmt.Sprintf("The ID should be numeric Error : %s", err),
		}}
	}

	rules, statusCode, diags := client.ReadIncapRulePriorities(d.Get("site_id").(string), d.Get("category").(string))

	fmt.Println(rules)
	// If the rule is deleted on the server, blow it out locally and run through the normal TF cycle
	if statusCode != 200 {
		return diags
	}

	if diags != nil {
		return diags
	}

	// Update all of the properties
	//d.Set("name", rule.Name)

	//action := d.Get("action").(string)

	//if action == "RULE_ACTION_RESPONSE_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_COOKIE" {
	//	if rule.RewriteExisting != nil {
	//		d.Set("rewrite_existing", *rule.RewriteExisting)
	//	}
	//} else {
	//	//align with schema default to avoid diff when importing resources
	//	d.Set("rewrite_existing", true)
	//}

	return nil
}
