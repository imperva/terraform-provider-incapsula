package incapsula

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceIncapRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIncapRuleCreate,
		Read:   resourceIncapRuleRead,
		Update: resourceIncapRuleUpdate,
		Delete: resourceIncapRuleDelete,
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
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Rule name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "Rule action. See the detailed descriptions in the API documentation. Possible values: `RULE_ACTION_REDIRECT`, `RULE_ACTION_SIMPLIFIED_REDIRECT`, `RULE_ACTION_REWRITE_URL`, `RULE_ACTION_REWRITE_HEADER`, `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_DELETE_HEADER`, `RULE_ACTION_DELETE_COOKIE`, `RULE_ACTION_RESPONSE_REWRITE_HEADER`, `RULE_ACTION_RESPONSE_DELETE_HEADER`, `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE`, `RULE_ACTION_FORWARD_TO_DC`, `RULE_ACTION_ALERT`, `RULE_ACTION_BLOCK`, `RULE_ACTION_BLOCK_USER`, `RULE_ACTION_BLOCK_IP`, `RULE_ACTION_RETRY`, `RULE_ACTION_INTRUSIVE_HTML`, `RULE_ACTION_CAPTCHA`, `RULE_ACTION_RATE`, `RULE_ACTION_CUSTOM_ERROR_RESPONSE`",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Optional Arguments
			"filter": {
				Description: "The filter defines the conditions that trigger the rule action. For action `RULE_ACTION_SIMPLIFIED_REDIRECT` filter is not relevant. For other actions, if left empty, the rule is always run.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"response_code": {
				Description: "For `RULE_ACTION_REDIRECT` or `RULE_ACTION_SIMPLIFIED_REDIRECT` rule's response code, valid values are `302`, `301`, `303`, `307`, `308`. For `RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE` rule's response code, valid values are all 3-digits numbers. For `RULE_ACTION_CUSTOM_ERROR_RESPONSE`, valid values are `400`, `401`, `402`, `403`, `404`, `405`, `406`, `407`, `408`, `409`, `410`, `411`, `412`, `413`, `414`, `415`, `416`, `417`, `419`, `420`, `422`, `423`, `424`, `500`, `501`, `502`, `503`, `504`, `505`, `507`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"add_missing": {
				Description: "Add cookie or header if it doesn't exist (Rewrite cookie rule only).",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"from": {
				Description: "Pattern to rewrite. For `RULE_ACTION_REWRITE_URL` - Url to rewrite. For `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to rewrite. For `RULE_ACTION_REWRITE_COOKIE` - Cookie value to rewrite.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"to": {
				Description: "Pattern to change to. `RULE_ACTION_REWRITE_URL` - Url to change to. `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER` - Header value to change to. `RULE_ACTION_REWRITE_COOKIE` - Cookie value to change to.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rewrite_name": {
				Description: "Name of cookie or header to rewrite. Applies only for `RULE_ACTION_REWRITE_COOKIE`, `RULE_ACTION_REWRITE_HEADER` and `RULE_ACTION_RESPONSE_REWRITE_HEADER`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dc_id": {
				Description: "Data center to forward request to. Applies only for `RULE_ACTION_FORWARD_TO_DC`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"rate_context": {
				Description: "The context of the rate counter. Possible values `IP` or `Session`. Applies only to rules using `RULE_ACTION_RATE`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rate_interval": {
				Description: "The interval in seconds of the rate counter. Possible values is a multiple of `10`; minimum `10` and maximum `300`. Applies only to rules using `RULE_ACTION_RATE`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"error_type": {
				Description: "The error that triggers the rule. `error.type.all` triggers the rule regardless of the error type. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `error.type.all`, `error.type.connection_timeout`, `error.type.access_denied`, `error.type.parse_req_error`, `error.type.parse_resp_error`, `error.type.connection_failed`, `error.type.deny_and_retry`, `error.type.ssl_failed`, `error.type.deny_and_captcha`, `error.type.2fa_required`, `error.type.no_ssl_config`, `error.type.no_ipv6_config`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_response_format": {
				Description: "The format of the given error response in the error_response_data field. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`. Possible values: `json`, `xml`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"error_response_data": {
				Description: "The response returned when the request matches the filter and is blocked. Applies only for `RULE_ACTION_CUSTOM_ERROR_RESPONSE`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceIncapRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	rule := IncapRule{
		Name:                d.Get("name").(string),
		Action:              d.Get("action").(string),
		Filter:              d.Get("filter").(string),
		ResponseCode:        d.Get("response_code").(int),
		AddMissing:          d.Get("add_missing").(bool),
		From:                d.Get("from").(string),
		To:                  d.Get("to").(string),
		RewriteName:         d.Get("rewrite_name").(string),
		DCID:                d.Get("dc_id").(int),
		RateContext:         d.Get("rate_context").(string),
		RateInterval:        d.Get("rate_interval").(int),
		ErrorType:           d.Get("error_type").(string),
		ErrorResponseFormat: d.Get("error_response_format").(string),
		ErrorResponseData:   d.Get("error_response_data").(string),
	}

	ruleWithID, err := client.AddIncapRule(d.Get("site_id").(string), &rule)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ruleWithID.RuleID))

	return resourceIncapRuleRead(d, m)
}

func resourceIncapRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	rule, statusCode, err := client.ReadIncapRule(d.Get("site_id").(string), ruleID)

	// If the rule is deleted on the server, blow it out locally and run through the normal TF cycle
	if statusCode == 404 {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	// Update all of the properties
	d.Set("name", rule.Name)
	d.Set("action", rule.Action)
	d.Set("filter", rule.Filter)
	d.Set("response_code", rule.ResponseCode)
	d.Set("add_missing", rule.AddMissing)
	d.Set("from", rule.From)
	d.Set("to", rule.To)
	d.Set("rewrite_name", rule.RewriteName)
	d.Set("dc_id", rule.DCID)
	d.Set("rate_context", rule.RateContext)
	d.Set("rate_interval", rule.RateInterval)
	d.Set("error_type", rule.ErrorType)
	d.Set("error_response_format", rule.ErrorResponseFormat)
	d.Set("error_response_data", rule.ErrorResponseData)

	return nil
}

func resourceIncapRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	rule := IncapRule{
		Name:                d.Get("name").(string),
		Action:              d.Get("action").(string),
		Filter:              d.Get("filter").(string),
		ResponseCode:        d.Get("response_code").(int),
		AddMissing:          d.Get("add_missing").(bool),
		From:                d.Get("from").(string),
		To:                  d.Get("to").(string),
		RewriteName:         d.Get("rewrite_name").(string),
		DCID:                d.Get("dc_id").(int),
		RateContext:         d.Get("rate_context").(string),
		RateInterval:        d.Get("rate_interval").(int),
		ErrorType:           d.Get("error_type").(string),
		ErrorResponseFormat: d.Get("error_response_format").(string),
		ErrorResponseData:   d.Get("error_response_data").(string),
	}

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	_, err = client.UpdateIncapRule(d.Get("site_id").(string), ruleID, &rule)

	if err != nil {
		return err
	}

	return nil
}

func resourceIncapRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteIncapRule(d.Get("site_id").(string), ruleID)
	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
