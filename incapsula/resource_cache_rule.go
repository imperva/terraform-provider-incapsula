package incapsula

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCacheRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceCacheRuleCreate,
		Read:   resourceCacheRuleRead,
		Update: resourceCacheRuleUpdate,
		Delete: resourceCacheRuleDelete,
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
				ForceNew:    true,
			},
			"name": {
				Description: "Rule name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description: "Rule action. See the detailed descriptions in the API documentation. Possible values: `HTTP_CACHE_MAKE_STATIC`, `HTTP_CACHE_CLIENT_CACHE_CTL`, `HTTP_CACHE_FORCE_UNCACHEABLE`, `HTTP_CACHE_ADD_TAG`, `HTTP_CACHE_DIFFERENTIATE_SSL`, `HTTP_CACHE_DIFFERENTIATE_BY_HEADER`, `HTTP_CACHE_DIFFERENTIATE_BY_COOKIE`, `HTTP_CACHE_DIFFERENTIATE_BY_GEO`, `HTTP_CACHE_IGNORE_PARAMS`, `HTTP_CACHE_ENRICH_CACHE_KEY`, `HTTP_CACHE_FORCE_VALIDATION`, `HTTP_CACHE_IGNORE_AUTH_HEADER`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"filter": {
				Description: "The filter defines the conditions that trigger the rule action, if left empty, the rule is always run.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "Boolean that specifies if the rule should be enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			// Optional Arguments
			"ttl": {
				Description: "TTL in seconds. Relevant for `HTTP_CACHE_MAKE_STATIC` and `HTTP_CACHE_CLIENT_CACHE_CTL` actions.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"ignored_params": {
				Description: "Parameters to ignore. Relevant for `HTTP_CACHE_IGNORE_PARAMS` action. An array containing `'*'` means all parameters are ignored.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"text": {
				Description: "Tag name if action is `HTTP_CACHE_ADD_TAG` action, text to be added to the cache key as suffix if action is `HTTP_CACHE_ENRICH_CACHE_KEY`.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"differentiate_by_value": {
				Description: "Value to differentiate resources by. Relevant for `HTTP_CACHE_DIFFERENTIATE_BY_HEADER`, `HTTP_CACHE_DIFFERENTIATE_BY_COOKIE` and `HTTP_CACHE_DIFFERENTIATE_BY_GEO` actions.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCacheRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	rule := CacheRule{
		Name:                 d.Get("name").(string),
		Action:               d.Get("action").(string),
		Filter:               d.Get("filter").(string),
		Enabled:              d.Get("enabled").(bool),
		TTL:                  d.Get("ttl").(int),
		IgnoredParams:        d.Get("ignored_params").(string),
		Text:                 d.Get("text").(string),
		DifferentiateByValue: d.Get("differentiate_by_value").(string),
	}

	ruleWithID, err := client.AddCacheRule(d.Get("site_id").(string), &rule)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(ruleWithID.RuleID))

	return resourceCacheRuleRead(d, m)
}

func resourceCacheRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	rule, statusCode, err := client.ReadCacheRule(d.Get("site_id").(string), ruleID)

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
	d.Set("enabled", rule.Enabled)
	d.Set("ttl", rule.TTL)
	d.Set("ignored_params", rule.IgnoredParams)
	d.Set("text", rule.Text)
	d.Set("differentiate_by_value", rule.DifferentiateByValue)

	return nil
}

func resourceCacheRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	rule := CacheRule{
		Name:                 d.Get("name").(string),
		Action:               d.Get("action").(string),
		Filter:               d.Get("filter").(string),
		Enabled:              d.Get("enabled").(bool),
		TTL:                  d.Get("ttl").(int),
		IgnoredParams:        d.Get("ignored_params").(string),
		Text:                 d.Get("text").(string),
		DifferentiateByValue: d.Get("differentiate_by_value").(string),
	}

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = client.UpdateCacheRule(d.Get("site_id").(string), ruleID, &rule)

	if err != nil {
		return err
	}

	return resourceCacheRuleRead(d, m)
}

func resourceCacheRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	err = client.DeleteCacheRule(d.Get("site_id").(string), ruleID)
	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
