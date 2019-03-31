package incapsula

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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

				siteID, err := strconv.Atoi(idSlice[0])
				ruleID := idSlice[1]
				if err != nil {
					return nil, err
				}

				d.Set("site_id", siteID)
				d.Set("rule_id", ruleID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"enabled": {
				Description: "todo",
				Type:        schema.TypeString,
				Required:    true,
			},
			"priority": {
				Description: "todo",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "todo",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"action": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"filter": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rule_id": {
				Description: "todo",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"allow_caching": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"dc_id": {
				Description: "todo",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"from": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"to": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"response_code": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"add_missing": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"rewrite_name": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceIncapRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.AddIncapRule(
		d.Get("site_id").(int),
		d.Get("rule_id").(int),
		d.Get("dc_id").(int),
		d.Get("enabled").(string),
		d.Get("priority").(string),
		d.Get("name").(string),
		d.Get("action").(string),
		d.Get("filter").(string),
		d.Get("allow_caching").(string),
		d.Get("response_code").(string),
		d.Get("from").(string),
		d.Get("to").(string),
		d.Get("add_missing").(string),
		d.Get("rewrite_name").(string),
	)

	if err != nil {
		return err
	}

	// Set the rule ID
	d.SetId(d.Get("rule_id").(string))

	return resourceIncapRuleRead(d, m)
}

func resourceIncapRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	listIncapRulesResponse, err := client.ListIncapRules(
		d.Get("include_ad_rules").(string),
		d.Get("include_incap_rules").(string),
	)
	d.Set("todo", listIncapRulesResponse)

	if err != nil {
		return err
	}

	// todo: what is the response

	return nil
}

func resourceIncapRuleUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceIncapRuleCreate(d, m)
}

func resourceIncapRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	// Implement delete by clearing out the rule configuration
	err := client.DeleteIncapRule(
		d.Get("rule_id").(int),
	)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
