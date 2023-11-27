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

func resourceSimplifiedRedirectRulesConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSimplifiedRedirectRulesConfigurationUpdate,
		ReadContext:   resourceSimplifiedRedirectRulesConfigurationRead,
		UpdateContext: resourceSimplifiedRedirectRulesConfigurationUpdate,
		DeleteContext: resourceSimplifiedRedirectRulesConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(data.Id(), "/")
				if len(idSlice) != 1 || idSlice[0] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id", data.Id())
				}

				data.Set("site_id", idSlice[0])
				data.SetId(idSlice[0])

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

			"rule": {
				Description: "List of simplified redirect rules",
				Optional:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_name": {
							Type:        schema.TypeString,
							Description: "The rule name",
							Required:    true,
						},

						"enabled": {
							Type:        schema.TypeBool,
							Description: "Boolean that enables the rule",
							Optional:    true,
							Default:     true,
						},

						"from": {
							Type:        schema.TypeString,
							Description: "From value",
							Required:    true,
						},

						"to": {
							Type:        schema.TypeString,
							Description: "To value",
							Required:    true,
						},

						"response_code": {
							Type:             schema.TypeInt,
							Description:      "Rule's response code",
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IntInSlice([]int{301, 302, 303, 307, 308})),
						},
					},
				},
			},
		},
	}
}

func resourceSimplifiedRedirectRulesConfigurationUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	siteID := data.Get("site_id").(string)

	rulesListDTO := DeliveryRulesListDTO{
		RulesList: createSimplifiedRedirectRulesListFromState(data),
	}

	simplifiedRediectRulesListDTO, diags := client.UpdateDeliveryRuleConfiguration(siteID, "SIMPLIFIED_REDIRECT", &rulesListDTO)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to update delivery rules of category SIMPLIFIED_REDIRECT for Site ID %s", siteID)
		return diags
	} else if simplifiedRediectRulesListDTO.Errors != nil {
		errors, _ := json.Marshal(simplifiedRediectRulesListDTO.Errors)
		log.Printf("[ERROR] Failed to update delivery rules of category SIMPLIFIED_REDIRECT for Site ID %s: %s", siteID, string(errors[:]))
		return []diag.Diagnostic{{
			Severity: diag.Error,
			Summary:  "Failed to update simplified redirect rules",
			Detail:   fmt.Sprintf("Failed to update delivery rules of category SIMPLIFIED_REDIRECT for Site ID %s: %s", siteID, string(errors[:])),
		}}
	}

	diags = append(diags, resourceSimplifiedRedirectRulesConfigurationRead(ctx, data, m)[:]...)
	return diags
}

func resourceSimplifiedRedirectRulesConfigurationRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)
	siteID := data.Get("site_id").(string)

	simplifiedRediectRulesListDTO, diags := client.ReadDeliveryRuleConfiguration(siteID, "SIMPLIFIED_REDIRECT")

	fmt.Println(simplifiedRediectRulesListDTO)

	if simplifiedRediectRulesListDTO != nil && simplifiedRediectRulesListDTO.Errors != nil && simplifiedRediectRulesListDTO.Errors[0].Status == 404 {
		log.Printf("[INFO] Incapsula Site with ID %s has already been deleted\n", data.Get("site_id"))
		data.SetId("")
		return nil
	}
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to read delivery rules in category SIMPLIFIED_REDIRECT for Site ID %s", siteID)
		return diags
	}

	data.Set("rule", serializeSimplifiedRedirectRule(data, *simplifiedRediectRulesListDTO))
	data.SetId(siteID)

	return nil
}

func serializeSimplifiedRedirectRule(data *schema.ResourceData, DeliveryRule DeliveryRulesListDTO) *schema.Set {
	simplifiedRedirectRules := &schema.Set{F: resourceSimplifiedRedirectConfigurationHashFunction}

	for _, rule := range DeliveryRule.RulesList {
		simplifiedRedirectRule := map[string]interface{}{}
		simplifiedRedirectRule["rule_name"] = rule.RuleName
		simplifiedRedirectRule["from"] = rule.From
		simplifiedRedirectRule["to"] = rule.To
		simplifiedRedirectRule["response_code"] = rule.ResponseCode
		simplifiedRedirectRule["enabled"] = rule.Enabled
		simplifiedRedirectRules.Add(simplifiedRedirectRule)
	}

	return simplifiedRedirectRules
}

func resourceSimplifiedRedirectRulesConfigurationDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteID := data.Get("site_id").(string)

	emptyRulesList := DeliveryRulesListDTO{
		RulesList: []DeliveryRuleDto{},
	}

	_, diags = client.UpdateDeliveryRuleConfiguration(siteID, "SIMPLIFIED_REDIRECT", &emptyRulesList)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to delete delivery rules in category SIMPLIFIED_REDIRECT for Site ID %s", siteID)
		return diags
	}

	data.SetId("")
	return diags
}

func createSimplifiedRedirectRulesListFromState(data *schema.ResourceData) []DeliveryRuleDto {
	simplifiedRediectRulesListConf := data.Get("rule").(*schema.Set)
	simplifiedRediectRulesListDTO := make([]DeliveryRuleDto, len(simplifiedRediectRulesListConf.List()))

	for i, deliveryRuleRaw := range simplifiedRediectRulesListConf.List() {
		deliveryRule := deliveryRuleRaw.(map[string]interface{})

		deliveryRuleDTO := DeliveryRuleDto{
			RuleName:     deliveryRule["rule_name"].(string),
			Action:       "RULE_ACTION_SIMPLIFIED_REDIRECT",
			Enabled:      deliveryRule["enabled"].(bool),
			From:         deliveryRule["from"].(string),
			To:           deliveryRule["to"].(string),
			ResponseCode: deliveryRule["response_code"].(int),
		}
		simplifiedRediectRulesListDTO[i] = deliveryRuleDTO
	}

	return simplifiedRediectRulesListDTO
}

func resourceSimplifiedRedirectConfigurationHashFunction(v interface{}) int {
	return schema.HashString(v.(map[string]interface{})["from"])
}
