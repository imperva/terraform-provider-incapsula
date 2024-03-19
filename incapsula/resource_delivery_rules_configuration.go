package incapsula

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDeliveryRulesConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeliveryRulesConfigurationUpdate,
		ReadContext:   resourceDeliveryRulesConfigurationRead,
		UpdateContext: resourceDeliveryRulesConfigurationUpdate,
		DeleteContext: resourceDeliveryRulesConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(data.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/category", data.Id())
				}

				data.Set("site_id", idSlice[0])
				data.Set("category", idSlice[1])

				return []*schema.ResourceData{data}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"category": {
				Description:      "How to load balance between multiple Data Centers.",
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"REDIRECT", "REWRITE", "REWRITE_RESPONSE", "FORWARD"}, false)),
			},

			"rule": {
				Description: "List of delivery rules",
				Optional:    true,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Description: "The rule name",
							Required:    true,
						},

						"action": {
							Type:             schema.TypeString,
							Description:      "Rule action",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"RULE_ACTION_REDIRECT", "RULE_ACTION_REWRITE_URL", "RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_DELETE_HEADER", "RULE_ACTION_DELETE_COOKIE", "RULE_ACTION_FORWARD_TO_DC", "RULE_ACTION_FORWARD_TO_PORT", "RULE_ACTION_RESPONSE_REWRITE_HEADER", "RULE_ACTION_RESPONSE_DELETE_HEADER", "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE", "RULE_ACTION_CUSTOM_ERROR_RESPONSE"}, false)),
						},

						"enabled": {
							Type:        schema.TypeBool,
							Description: "Boolean that enables the rule",
							Optional:    true,
							Default:     true,
						},

						"filter": {
							Type:        schema.TypeString,
							Description: "Defines the conditions that trigger the rule action",
							Optional:    true,
						},

						"from": {
							Type:        schema.TypeString,
							Description: "From value",
							Optional:    true,
						},

						"to": {
							Type:        schema.TypeString,
							Description: "To value",
							Optional:    true,
						},

						"response_code": {
							Type:             schema.TypeInt,
							Description:      "Rule's response code",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(1, 999)),
						},

						"cookie_name": {
							Type:        schema.TypeString,
							Description: "Name of cookie to modify",
							Optional:    true,
						},

						"header_name": {
							Type:        schema.TypeString,
							Description: "Name of header to modify",
							Optional:    true,
						},

						"rewrite_existing": {
							Type:        schema.TypeBool,
							Description: "Apply rewrite rule even if the header/cookie already exists",
							Optional:    true,
							Default:     true,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								//aesthetic only - prevent showing default value for irrelevant action types
								pathParts := strings.Split(k, ".")
								action, _ := d.GetOk(pathParts[0] + "." + pathParts[1] + ".action")
								return !contains(ruleArgsToActionMap["rewrite_existing"], action.(string))
							},
						},

						"add_if_missing": {
							Type:        schema.TypeBool,
							Description: "Rewrite rule would add the header/cookie if it's missing",
							Optional:    true,
						},

						"multiple_headers_deletion": {
							Type:        schema.TypeBool,
							Description: "Delete multiple header occurrences",
							Optional:    true,
						},

						"error_response_format": {
							Type:             schema.TypeString,
							Description:      "The format of the given error response in the error_response_data field",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"JSON", "XML"}, false)),
						},

						"error_response_data": {
							Type:        schema.TypeString,
							Description: "The response returned when the request matches the filter and is blocked",
							Optional:    true,
						},

						"error_type": {
							Type:        schema.TypeString,
							Description: "The error that triggers the rule",
							Optional:    true,
						},

						"dc_id": {
							Type:        schema.TypeInt,
							Description: "Data center ID to forward the request to",
							Optional:    true,
						},

						"port_forwarding_context": {
							Type:             schema.TypeString,
							Description:      "Context for port forwarding",
							Optional:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"port", "header"}, true)),
						},

						"port_forwarding_value": {
							Type:        schema.TypeString,
							Description: "Port number or header name for port forwarding",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceDeliveryRulesConfigurationUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	siteID := data.Get("site_id").(string)
	category := data.Get("category").(string)

	diags := validateConfig(data)
	if diags != nil && diags.HasError() {
		return diags
	}

	rulesListDTO := DeliveryRulesListDTO{
		RulesList: createRulesListFromState(data),
	}

	deliveryRulesListDTO, diags := client.UpdateDeliveryRuleConfiguration(siteID, category, &rulesListDTO)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to update delivery rules of category %s for Site ID %s", category, siteID)
		return diags
	} else if deliveryRulesListDTO.Errors != nil {
		errors, _ := json.Marshal(deliveryRulesListDTO.Errors)
		log.Printf("[ERROR] Failed to update delivery rules of category %s for Site ID %s: %s", category, siteID, string(errors[:]))
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "Failed to update delivery rules",
			Detail:   fmt.Sprintf("Failed to update delivery rules of category %s for Site ID %s: %s", category, siteID, string(errors[:])),
		}}
	}

	diags = append(diags, resourceDeliveryRulesConfigurationRead(ctx, data, m)[:]...)
	return diags
}

func resourceDeliveryRulesConfigurationRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)
	siteID := data.Get("site_id").(string)
	category := data.Get("category").(string)

	deliveryRulesListDTO, diags := client.ReadDeliveryRuleConfiguration(siteID, category)

	fmt.Println(deliveryRulesListDTO)

	if deliveryRulesListDTO != nil && deliveryRulesListDTO.Errors != nil && deliveryRulesListDTO.Errors[0].Status == 404 {
		log.Printf("[INFO] Incapsula Site with ID %s has already been deleted\n", data.Get("site_id"))
		data.SetId("")
		return nil
	}
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to read delivery rules in category %s for Site ID %s", category, siteID)
		return diags
	}

	data.Set("rule", serializeDeliveryRule(data, *deliveryRulesListDTO))
	data.SetId(siteID + "/" + category)

	return nil
}

func serializeDeliveryRule(data *schema.ResourceData, DeliveryRule DeliveryRulesListDTO) []interface{} {
	RulesList := make([]interface{}, len(DeliveryRule.RulesList))
	for i, rule := range DeliveryRule.RulesList {
		RuleSlice := make(map[string]interface{})
		RuleSlice["rule_name"] = rule.RuleName
		RuleSlice["action"] = rule.Action
		RuleSlice["filter"] = rule.Filter
		RuleSlice["add_if_missing"] = rule.AddMissing
		RuleSlice["from"] = rule.From
		RuleSlice["to"] = rule.To
		RuleSlice["response_code"] = rule.ResponseCode
		RuleSlice["cookie_name"] = rule.CookieName
		RuleSlice["header_name"] = rule.HeaderName
		RuleSlice["dc_id"] = rule.DCID
		RuleSlice["port_forwarding_context"] = rule.PortForwardingContext
		RuleSlice["port_forwarding_value"] = rule.PortForwardingValue
		RuleSlice["error_type"] = rule.ErrorType
		RuleSlice["error_response_format"] = rule.ErrorResponseFormat
		RuleSlice["error_response_data"] = rule.ErrorResponseData
		RuleSlice["multiple_headers_deletion"] = rule.MultipleHeaderDeletions
		RuleSlice["enabled"] = rule.Enabled

		if rule.Action == "RULE_ACTION_RESPONSE_REWRITE_HEADER" || rule.Action == "RULE_ACTION_REWRITE_HEADER" || rule.Action == "RULE_ACTION_REWRITE_COOKIE" {
			RuleSlice["rewrite_existing"] = *rule.RewriteExisting
		} else {
			RuleSlice["rewrite_existing"] = false
		}

		RulesList[i] = RuleSlice
	}
	return RulesList
}

func resourceDeliveryRulesConfigurationDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteID := data.Get("site_id").(string)
	category := data.Get("category").(string)

	emptyRulesList := DeliveryRulesListDTO{
		RulesList: []DeliveryRuleDto{},
	}

	_, diags = client.UpdateDeliveryRuleConfiguration(siteID, category, &emptyRulesList)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to delete delivery rules in category %s for Site ID %s", category, siteID)
		return diags
	}

	data.SetId("")
	return diags
}

func createRulesListFromState(data *schema.ResourceData) []DeliveryRuleDto {
	deliveryRulesListConf := data.Get("rule").([]interface{})
	deliveryRulesListDTO := make([]DeliveryRuleDto, len(deliveryRulesListConf))

	for i, deliveryRuleRaw := range deliveryRulesListConf {
		deliveryRule := deliveryRuleRaw.(map[string]interface{})

		action := deliveryRule["action"].(string)

		deliveryRuleDTO := DeliveryRuleDto{
			RuleName: deliveryRule["rule_name"].(string),
			Action:   action,
			Enabled:  deliveryRule["enabled"].(bool),
			Filter:   deliveryRule["filter"].(string),
		}

		rewriteExisting := new(bool)
		//rewrite_existing mustn't be set for other rule types
		if action == "RULE_ACTION_RESPONSE_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_HEADER" || action == "RULE_ACTION_REWRITE_COOKIE" {
			*rewriteExisting = deliveryRule["rewrite_existing"].(bool)
		} else {
			rewriteExisting = nil
		}

		switch deliveryRule["action"] {
		case "RULE_ACTION_REDIRECT":
			deliveryRuleDTO.From = deliveryRule["from"].(string)
			deliveryRuleDTO.To = deliveryRule["to"].(string)
			deliveryRuleDTO.ResponseCode = deliveryRule["response_code"].(int)
		case "RULE_ACTION_REWRITE_URL":
			deliveryRuleDTO.From = deliveryRule["from"].(string)
			deliveryRuleDTO.To = deliveryRule["to"].(string)
		case "RULE_ACTION_REWRITE_HEADER":
			deliveryRuleDTO.HeaderName = deliveryRule["header_name"].(string)
			deliveryRuleDTO.From = deliveryRule["from"].(string)
			deliveryRuleDTO.To = deliveryRule["to"].(string)
			deliveryRuleDTO.RewriteExisting = rewriteExisting
			deliveryRuleDTO.AddMissing = deliveryRule["add_if_missing"].(bool)
		case "RULE_ACTION_RESPONSE_REWRITE_HEADER":
			deliveryRuleDTO.HeaderName = deliveryRule["header_name"].(string)
			deliveryRuleDTO.From = deliveryRule["from"].(string)
			deliveryRuleDTO.To = deliveryRule["to"].(string)
			deliveryRuleDTO.RewriteExisting = rewriteExisting
			deliveryRuleDTO.AddMissing = deliveryRule["add_if_missing"].(bool)
		case "RULE_ACTION_REWRITE_COOKIE":
			deliveryRuleDTO.CookieName = deliveryRule["cookie_name"].(string)
			deliveryRuleDTO.From = deliveryRule["from"].(string)
			deliveryRuleDTO.To = deliveryRule["to"].(string)
			deliveryRuleDTO.RewriteExisting = rewriteExisting
			deliveryRuleDTO.AddMissing = deliveryRule["add_if_missing"].(bool)
		case "RULE_ACTION_DELETE_HEADER":
			deliveryRuleDTO.HeaderName = deliveryRule["header_name"].(string)
			deliveryRuleDTO.MultipleHeaderDeletions = deliveryRule["multiple_headers_deletion"].(bool)
		case "RULE_ACTION_RESPONSE_DELETE_HEADER":
			deliveryRuleDTO.HeaderName = deliveryRule["header_name"].(string)
			deliveryRuleDTO.MultipleHeaderDeletions = deliveryRule["multiple_headers_deletion"].(bool)
		case "RULE_ACTION_DELETE_COOKIE":
			deliveryRuleDTO.CookieName = deliveryRule["cookie_name"].(string)
		case "RULE_ACTION_FORWARD_TO_DC":
			deliveryRuleDTO.DCID = deliveryRule["dc_id"].(int)
		case "RULE_ACTION_FORWARD_TO_PORT":
			deliveryRuleDTO.PortForwardingContext = deliveryRule["port_forwarding_context"].(string)
			deliveryRuleDTO.PortForwardingValue = deliveryRule["port_forwarding_value"].(string)
		case "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE":
			deliveryRuleDTO.ResponseCode = deliveryRule["response_code"].(int)
		case "RULE_ACTION_CUSTOM_ERROR_RESPONSE":
			deliveryRuleDTO.ErrorType = deliveryRule["error_type"].(string)
			deliveryRuleDTO.ErrorResponseFormat = deliveryRule["error_response_format"].(string)
			deliveryRuleDTO.ErrorResponseData = deliveryRule["error_response_data"].(string)
			deliveryRuleDTO.ResponseCode = deliveryRule["response_code"].(int)
		}
		deliveryRulesListDTO[i] = deliveryRuleDTO
	}

	return deliveryRulesListDTO
}

func validateConfig(data *schema.ResourceData) diag.Diagnostics {
	diags := []diag.Diagnostic{}
	rulesList := data.GetRawConfig().GetAttr("rule").AsValueSlice()

	for i := 0; i < len(data.Get("rule").([]interface{})); i++ {
		ruleAction := rulesList[i].GetAttr("action").AsString()

		for attr, allowedActions := range ruleArgsToActionMap {
			if !rulesList[i].GetAttr(attr).IsNull() && !contains(allowedActions, ruleAction) {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Configuration argument '%s' is not applicable to action %s", attr, ruleAction),
					Detail:   fmt.Sprintf("rule[%d].%s", i, attr),
				})
			}
		}
	}

	return diags
}

var ruleArgsToActionMap = map[string][]string{
	"from":                      {"RULE_ACTION_REDIRECT", "RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_RESPONSE_REWRITE_HEADER", "RULE_ACTION_REWRITE_URL"},
	"to":                        {"RULE_ACTION_REDIRECT", "RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_RESPONSE_REWRITE_HEADER", "RULE_ACTION_REWRITE_URL"},
	"response_code":             {"RULE_ACTION_REDIRECT", "RULE_ACTION_REWRITE_URL", "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE", "RULE_ACTION_CUSTOM_ERROR_RESPONSE"},
	"header_name":               {"RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_DELETE_HEADER", "RULE_ACTION_RESPONSE_REWRITE_HEADER", "RULE_ACTION_RESPONSE_DELETE_HEADER"},
	"cookie_name":               {"RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_DELETE_COOKIE"},
	"rewrite_existing":          {"RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_RESPONSE_REWRITE_HEADER"},
	"add_if_missing":            {"RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_RESPONSE_REWRITE_HEADER"},
	"multiple_headers_deletion": {"RULE_ACTION_DELETE_HEADER", "RULE_ACTION_RESPONSE_DELETE_HEADER"},
	"error_response_format":     {"RULE_ACTION_CUSTOM_ERROR_RESPONSE"},
	"error_response_data":       {"RULE_ACTION_CUSTOM_ERROR_RESPONSE"},
	"error_type":                {"RULE_ACTION_CUSTOM_ERROR_RESPONSE"},
	"port_forwarding_context":   {"RULE_ACTION_FORWARD_TO_PORT"},
	"port_forwarding_value":     {"RULE_ACTION_FORWARD_TO_PORT"},
	"dc_id":                     {"RULE_ACTION_FORWARD_TO_DC"},
	"filter":                    {"RULE_ACTION_REDIRECT", "RULE_ACTION_REWRITE_HEADER", "RULE_ACTION_REWRITE_COOKIE", "RULE_ACTION_REWRITE_URL", "RULE_ACTION_DELETE_HEADER", "RULE_ACTION_DELETE_COOKIE", "RULE_ACTION_RESPONSE_REWRITE_HEADER", "RULE_ACTION_RESPONSE_DELETE_HEADER", "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE", "RULE_ACTION_CUSTOM_ERROR_RESPONSE", "RULE_ACTION_FORWARD_TO_DC", "RULE_ACTION_FORWARD_TO_PORT"},
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
