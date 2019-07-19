package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

// Security Rule Enumerations
// NOTE: no whitelists for rule: api.acl.whitelisted_ips
const blacklistedCountriesExceptionRuleId = "api.acl.blacklisted_countries"
const blacklistedURLsExceptionRuleId = "api.acl.blacklisted_urls"
const blacklistedIPsExceptionRuleId = "api.acl.blacklisted_ips"
const backdoorExceptionRuleId = "api.threats.backdoor"
const crossSiteScriptingExceptionRuleId = "api.threats.cross_site_scripting"
const illegalResourceAccessExceptionRuleId = "api.threats.illegal_resource_access"
const remoteFileInclusionExceptionRuleId = "api.threats.remote_file_inclusion"
const sqlInjectionExceptionRuleId = "api.threats.sql_injection"
const ddosExceptionRuleId = "api.threats.ddos"
const botAccessControlExceptionRuleId = "api.threats.bot_access_control"

type DeleteSecurityRuleExceptionResponse struct {
	Res int `json:"res"`
}

func resourceSecurityRuleException() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityRuleExceptionCreate,
		Read:   resourceSecurityRuleExceptionRead,
		Update: resourceSecurityRuleExceptionUpdate,
		Delete: resourceSecurityRuleExceptionDelete,
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
				Description: "The identifier of the security rule, e.g api.threats.cross_site_scripting.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_app_types": &schema.Schema{
				Description: "A comma separated list of client application types,",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"client_apps": &schema.Schema{
				Description: "A comma separated list of client application IDs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"countries": &schema.Schema{
				Description: "A comma separated list of country codes.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"continents": &schema.Schema{
				Description: "A comma separated list of continent codes.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ips": &schema.Schema{
				Description: "A comma separated list of IPs or IP ranges, e.g: 192.168.1.1, 192.168.1.1-192.168.1.100 or 192.168.1.1/24",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"url_patterns": &schema.Schema{
				Description: "A comma separated list of url patterns. One of: contains | equals | prefix | suffix | not_equals | not_contain | not_prefix | not_suffix. The patterns should be in accordance with the matching urls sent by the urls parameter.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"urls": &schema.Schema{
				Description: "A comma separated list of resource paths. For example, /home and /admin/index.html are resource paths, while http://www.example.com/home is not. Each URL should be encoded separately using percent encoding as specified by RFC 3986 (http://tools.ietf.org/html/rfc3986#section-2.1). An empty URL list will remove all URLs.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"user_agents": &schema.Schema{
				Description: "A comma separated list of encoded user agents.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"parameters": &schema.Schema{
				Description: "A comma separated list of encoded parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"whitelist_id": &schema.Schema{
				Description: "The id (an integer) of the whitelist to be set. This field is optional - in case no id is supplied, a new whitelist will be created.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"exception_id_only": &schema.Schema{
				Description: "The id (an integer) of the whitelist to be set. This field is optional - in case no id is supplied, a new whitelist will be created.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceSecurityRuleExceptionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Configuring Incapsula Security Rule Exception for rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))
	siteStatusResponse, err := client.AddSecurityRuleException(
		d.Get("site_id").(int),
		ruleID,
		d.Get("client_app_types").(string),
		d.Get("client_apps").(string),
		d.Get("countries").(string),
		d.Get("continents").(string),
		d.Get("ips").(string),
		d.Get("url_patterns").(string),
		d.Get("urls").(string),
		d.Get("user_agents").(string),
		d.Get("parameters").(string),
	)
	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
		return err
	}

	// Set the rule exception ID
	d.SetId(siteStatusResponse.ExceptionID)

	log.Printf("[INFO] Created Incapsula security rule exception for rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	return resourceSecurityRuleExceptionRead(d, m)
}

func resourceSecurityRuleExceptionRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	siteID := strconv.Itoa(d.Get("site_id").(int))
	ruleID := d.Get("rule_id").(string)
	whitelistID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Reading Incapsula security rule exception whitelist_id (%d) on rule_id (%s) \n", whitelistID, ruleID)

	siteStatusResponse, err := client.ListSecurityRuleExceptions(siteID, ruleID)

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula security rule exception whitelist_id (%d) on rule_id (%s) %s\n", whitelistID, ruleID, err)
		return err
	}

	// Now with the site status, iterate through the rules and find our ID
	exceptionFound := false
	if ruleID == blacklistedCountriesExceptionRuleId || ruleID == blacklistedURLsExceptionRuleId || ruleID == blacklistedIPsExceptionRuleId {
		for _, entry := range siteStatusResponse.Security.Acls.Rules {
			if entry.ID == d.Get("rule_id").(string) {
				for _, exception := range entry.Exceptions {
					if exception.ID == whitelistID {
						for _, value := range exception.Values {
							d.Set(value.ID, value.Name)
							exceptionFound = true
							break
						}
					}
				}
			}
		}
	} else {
		for _, entry := range siteStatusResponse.Security.Waf.Rules {
			if entry.ID == d.Get("rule_id").(string) {
				for _, exception := range entry.Exceptions {
					if exception.ID == whitelistID {
						for _, value := range exception.Values {
							d.Set(value.ID, value.Name)
							exceptionFound = true
							break
						}
					}
				}
			}
		}
	}
	if exceptionFound == false {
		log.Printf("[ERROR] Read Incapsula security rule exception failed, exception not found: whitelist_id (%d) and rule_id (%s) on site_id (%d)\n", whitelistID, ruleID, d.Get("site_id").(int))
	} else {
		log.Printf("[INFO] Read Incapsula security rule exception whitelist_id (%d) and rule_id (%s) on site_id (%d)\n", whitelistID, ruleID, d.Get("site_id").(int))
	}

	return nil
}

func resourceSecurityRuleExceptionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)
	whitelistID := d.Id()

	log.Printf("[INFO] Updating Incapsula security rule exception for rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	// Add the appropriate exception params based on ruleID, set exception_id_only to return the whitelist_id for newly created rule
	switch ruleID {
	// ACL RuleIDs
	case blacklistedCountriesExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			d.Get("client_app_types").(string),
			"",
			"",
			"",
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			"",
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case blacklistedIPsExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			"",
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case blacklistedURLsExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			"",
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case backdoorExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			d.Get("user_agents").(string),
			d.Get("parameters").(string),
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case botAccessControlExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			"",
			"",
			"",
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			d.Get("user_agents").(string),
			"",
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case crossSiteScriptingExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			"",
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			d.Get("parameters").(string),
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case ddosExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			"",
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case illegalResourceAccessExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			d.Get("parameters").(string),
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case remoteFileInclusionExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			d.Get("user_agents").(string),
			d.Get("parameters").(string),
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	case sqlInjectionExceptionRuleId:
		_, err := client.EditSecurityRuleException(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("client_apps").(string),
			d.Get("countries").(string),
			d.Get("continents").(string),
			d.Get("ips").(string),
			d.Get("url_patterns").(string),
			d.Get("urls").(string),
			"",
			"",
			whitelistID,
		)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula security rule exception for rule_id (%s) on site_id (%d), %s\n", ruleID, d.Get("site_id").(int), err)
			return err
		}
	}

	// Set the rule ID as whitelistID
	d.SetId(whitelistID)

	log.Printf("[INFO] Updated Incapsula security rule exception for rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	return resourceWAFSecurityRuleRead(d, m)
}

func resourceSecurityRuleExceptionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)
	whitelistID := d.Id()

	log.Printf("[INFO] Deleting Incapsula security rule exception whitelist_id (%s) for rule_id (%s) on site_id (%d)\n", whitelistID, ruleID, d.Get("site_id").(int))

	err := client.DeleteSecurityRuleException(
		d.Get("site_id").(int),
		ruleID,
		whitelistID,
	)
	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula security rule exception whitelist_id (%s) for rule_id (%s) on site_id (%d), %s\n", whitelistID, ruleID, d.Get("site_id").(int), err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
