package incapsula

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceACLSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceACLSecurityRuleCreate,
		Read:   resourceACLSecurityRuleRead,
		Update: resourceACLSecurityRuleUpdate,
		Delete: resourceACLSecurityRuleDelete,

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
			"countries": &schema.Schema{
				Description: "A comma separated list of country codes.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ips": &schema.Schema{
				Description: "A comma separated list of IPs or IP ranges, e.g: 192.168.1.1, 192.168.1.1-192.168.1.100 or 192.168.1.1/24.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"urls": &schema.Schema{
				Description: "A comma separated list of resource paths.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"url_patterns": &schema.Schema{
				Description: "The patterns should be in accordance with the matching urls sent by the urls parameter. Options: CONTAINS | EQUALS | PREFIX | SUFFIX | NOT_EQUALS | NOT_CONTAIN | NOT_PREFIX | NOT_SUFFIX",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceACLSecurityRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.ConfigureACLSecurityRule(
		d.Get("site_id").(int),
		d.Get("rule_id").(string),
		d.Get("countries").(string),
		d.Get("ips").(string),
		d.Get("urls").(string),
		d.Get("url_patterns").(string),
	)

	if err != nil {
		return err
	}

	// Set the rule ID
	d.SetId(d.Get("rule_id").(string))

	return resourceACLSecurityRuleRead(d, m)
}

func resourceACLSecurityRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	siteStatusResponse, err := client.SiteStatus("acl-rule-read", d.Get("site_id").(int))

	if err != nil {
		return err
	}

	// Now with the site status, iterate through the rules and find our ID
	for _, entry := range siteStatusResponse.Security.Acls.Rules {
		if entry.ID == d.Get("rule_id").(string) {
			// Set different attributes based on the rule id
			switch entry.ID {
			case blacklistedCountries:
				d.Set("countries", strings.Join(entry.Geo.Countries, ","))
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

	return nil
}

func resourceACLSecurityRuleUpdate(d *schema.ResourceData, m interface{}) error {
	// This is the same as create
	return resourceACLSecurityRuleCreate(d, m)
}

func resourceACLSecurityRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	// Implement delete by clearing out the rule configuration
	_, err := client.ConfigureACLSecurityRule(
		d.Get("site_id").(int),
		d.Get("rule_id").(string),
		"", // countries
		"", // ips
		"", // urls
		"", // url_patterns
	)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
