package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceACLSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceACLSecurityRuleCreate,
		Read:   resourceACLSecurityRuleRead,
		Update: resourceACLSecurityRuleUpdate,
		Delete: resourceACLSecurityRuleDelete,
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
			"site_id": &schema.Schema{
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"rule_id": &schema.Schema{
				Description: "The id of the acl, e.g api.acl.blacklisted_ips.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"continents": &schema.Schema{
				Description: "A comma separated list of continents codes.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"countries": &schema.Schema{
				Description: "A list of country codes.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"ips": &schema.Schema{
				Description: "A list of IPs or IP ranges, e.g: [192.168.1.1] or [192.168.1.1-192.168.1.100] or [192.168.1.1/24].",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"urls": &schema.Schema{
				Description: "A comma separated list of resource paths. NOTE: this is a 1:1 list with url_patterns e.q:  urls = \"Test,/Values\" url_patterns = \"CONTAINS,PREFIX\"",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"url_patterns": &schema.Schema{
				Description: "The patterns should be in accordance with the matching urls sent by the urls parameter. Options: CONTAINS | EQUALS | PREFIX | SUFFIX | NOT_EQUALS | NOT_CONTAIN | NOT_PREFIX | NOT_SUFFIX",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"client_apps": &schema.Schema{
				Description: "The client apps",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
		},
	}
}

func resourceACLSecurityRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Creating Incapsula ACL Rule for id: %s\n", ruleID)

	_, err := client.ConfigureACLSecurityRule(
		d.Get("site_id").(int),
		ruleID,
		strings.Join(convertStringArr(d.Get("continents").([]interface{})), ","),
		strings.Join(convertStringArr(d.Get("countries").([]interface{})), ","),
		strings.Join(convertStringArr(d.Get("ips").([]interface{})), ","),
		strings.Join(convertStringArr(d.Get("urls").([]interface{})), ","),
		strings.Join(convertStringArr(d.Get("url_patterns").([]interface{})), ","),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula ACL Rule for id: %s, %s\n", ruleID, err)
		return err
	}

	// Set the rule ID
	d.SetId(d.Get("rule_id").(string))

	log.Printf("[INFO] Created Incapsula ACL Rule for id: %s\n", ruleID)

	return resourceACLSecurityRuleRead(d, m)
}

func resourceACLSecurityRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Reading Incapsula ACL Rule for id: %s\n", ruleID)

	siteStatusResponse, err := client.SiteStatus("acl-rule-read", d.Get("site_id").(int))

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula ACL Rule for id: %s, %s\n", ruleID, err)
		return err
	}

	// Now with the site status, iterate through the rules and find our ID
	for _, entry := range siteStatusResponse.Security.Acls.Rules {
		if entry.ID == d.Get("rule_id").(string) {
			// Set different attributes based on the rule id
			switch entry.ID {
			case blacklistedCountries:
				d.Set("countries", strings.Join(entry.Geo.Countries, ","))
				d.Set("continents", strings.Join(entry.Geo.Continents, ","))
			case blacklistedURLs:
				urls := make([]string, 0)
				urlPatterns := make([]string, 0)
				for _, url := range entry.Urls {
					urls = append(urls, url.Value)
					urlPatterns = append(urlPatterns, url.Pattern)
				}
				d.Set("urls", strings.Join(urls, ","))
				d.Set("url_patterns", strings.Join(urlPatterns, ","))
			case blacklistedIPs:
				d.Set("ips", strings.Join(entry.Ips, ","))
			case whitelistedIPs:
				d.Set("ips", strings.Join(entry.Ips, ","))
			}
			break
		}
	}

	log.Printf("[INFO] Read Incapsula ACL Rule for id: %s\n", ruleID)

	return nil
}

func resourceACLSecurityRuleUpdate(d *schema.ResourceData, m interface{}) error {
	// This is the same as create
	return resourceACLSecurityRuleCreate(d, m)
}

func resourceACLSecurityRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Deleting Incapsula ACL Rule for id: %s\n", ruleID)

	// Implement delete by clearing out the rule configuration
	_, err := client.ConfigureACLSecurityRule(
		d.Get("site_id").(int),
		ruleID,
		"", // countries
		"", // contenents
		"", // ips
		"", // urls
		"", // urls
	)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula ACL Rule for id: %s, %s\n", ruleID, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula ACL Rule for id: %s\n", ruleID)

	return nil
}
