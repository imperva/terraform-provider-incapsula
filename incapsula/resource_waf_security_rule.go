package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Default actions to reset polies to upon delete/destroy
const backdoorRuleIDDefaultAction = "api.threats.action.quarantine_url"
const crossSiteScriptingRuleIDDefaultAction = "api.threats.action.block_request"
const illegalResourceAccessRuleIDDefaultAction = "api.threats.action.block_request"
const remoteFileInclusionRuleIDDefaultAction = "api.threats.action.block_request"
const sqlInjectionRuleIDDefaultAction = "api.threats.action.block_request"
const ddosRuleIDDefaultActivationMode = "api.threats.ddos.activation_mode.auto"
const ddosRuleIDDefaultDDOSTrafficThreshold = "1000"
const botAccessControlBlockBadBotsDefaultAction = "true"
const botAccessControlChallengeSuspectedBotsDefaultAction = "false"

func resourceWAFSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceWAFSecurityRuleCreate,
		Read:   resourceWAFSecurityRuleRead,
		Update: resourceWAFSecurityRuleUpdate,
		Delete: resourceWAFSecurityRuleDelete,
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
				Description: "The identifier of the WAF rule, e.g api.threats.cross_site_scripting.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Required for rule_id: api.threats.backdoor, api.threats.cross_site_scripting, api.threats.illegal_resource_access, api.threats.remote_file_inclusion, api.threats.sql_injection
			"security_rule_action": &schema.Schema{
				Description: "The action that should be taken when a threat is detected, for example: api.threats.action.block_ip.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Required for rule_id: api.threats.ddos
			"activation_mode": &schema.Schema{
				Description: "The mode of activation for ddos on a site. Possible values: off, auto, on.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ddos_traffic_threshold": &schema.Schema{
				Description: "Consider site to be under DDoS if the request rate is above this threshold. The valid values are 10, 20, 50, 100, 200, 500, 750, 1000, 2000, 3000, 4000, 5000.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Required for rule_id: api.threats.bot_access_control
			"block_bad_bots": &schema.Schema{
				Description: "Whether or not to block bad bots. Possible values: true, false.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"challenge_suspected_bots": &schema.Schema{
				Description: "Whether or not to send a challenge to clients that are suspected to be bad bots (CAPTCHA for example). Possible values: true, false.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceWAFSecurityRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Creating Incapsula WAF Rule rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	if ruleID == backdoorRuleId || ruleID == crossSiteScriptingRuleId || ruleID == illegalResourceAccessRuleId || ruleID == remoteFileInclusionRuleId || ruleID == sqlInjectionRuleId {
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			d.Get("security_rule_action").(string),
			"",
			"",
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not create Incapsula WAF Rule rule_id (%s) and security_rule_action (%s) on site_id (%d), %s\n", ruleID, d.Get("security_rule_action").(string), d.Get("site_id").(int), err)
			return err
		}
	} else if ruleID == ddosRuleId {
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			"",
			d.Get("activation_mode").(string),
			d.Get("ddos_traffic_threshold").(string),
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not create Incapsula WAF Rule rule_id (%s) with activation_mode (%s) and ddos_traffic_threshold (%s) on site_id (%d), %s\n", ruleID, d.Get("activation_mode").(string), d.Get("ddos_traffic_threshold").(string), d.Get("site_id").(int), err)
			return err
		}
	} else if ruleID == botAccessControlRuleId {
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			"",
			"",
			"",
			d.Get("block_bad_bots").(string),
			d.Get("challenge_suspected_bots").(string),
		)
		if err != nil {
			log.Printf("[ERROR] Could not create Incapsula WAF Rule rule_id (%s) with block_bad_bots (%s) and challenge_suspected_bots (%s) on site_id (%d), %s\n", ruleID, d.Get("block_bad_bots").(string), d.Get("challenge_suspected_bots").(string), d.Get("site_id").(int), err)
			return err
		}
	}

	// Set the rule ID
	d.SetId(d.Get("rule_id").(string))

	log.Printf("[INFO] Created Incapsula WAF Rule rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	return resourceWAFSecurityRuleRead(d, m)
}

func resourceWAFSecurityRuleRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Reading Incapsula WAF Rule for id: %s\n", ruleID)

	siteStatusResponse, err := client.SiteStatus("waf-rule-read", d.Get("site_id").(int))

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula WAF Rule for id: %s, %s\n", ruleID, err)
		return err
	}

	// Now with the site status, iterate through the rules and find our ID
	for _, entry := range siteStatusResponse.Security.Waf.Rules {
		if entry.ID == d.Get("rule_id").(string) {
			// Set different attributes based on the rule id
			switch entry.ID {
			case backdoorRuleId:
				d.Set("action", entry.Action)
			case crossSiteScriptingRuleId:
				d.Set("action", entry.Action)
			case illegalResourceAccessRuleId:
				d.Set("action", entry.Action)
			case remoteFileInclusionRuleId:
				d.Set("action", entry.Action)
			case sqlInjectionRuleId:
				d.Set("action", entry.Action)
			case ddosRuleId:
				d.Set("activation_mode", entry.ActivationMode)
				d.Set("ddos_traffic_threshold", entry.DdosTrafficThreshold)
			case botAccessControlRuleId:
				d.Set("block_bad_bots", entry.BlockBadBots)
				d.Set("block_bad_bots", entry.ChallengeSuspectedBots)
			}
			break
		}
	}

	log.Printf("[INFO] Read Incapsula WAF Rule rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	return nil
}

func resourceWAFSecurityRuleUpdate(d *schema.ResourceData, m interface{}) error {
	// This is the same as create
	return resourceWAFSecurityRuleCreate(d, m)
}

func testAccStateWAFSecurityRuleID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_waf_security_rule" {
			continue
		}

		ruleID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, ruleID), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func resourceWAFSecurityRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	ruleID := d.Get("rule_id").(string)

	log.Printf("[INFO] Resetting Incapsula WAF Rule rule_id (%s) on site_id (%d)\n", ruleID, d.Get("site_id").(int))

	// Set WAF rule type defaults based on specific rule id
	switch ruleID {
	case backdoorRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			backdoorRuleIDDefaultAction,
			"",
			"",
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with security_rule_action (%s) on site_id (%d) %s\n", ruleID, backdoorRuleIDDefaultAction, d.Get("site_id").(int), err)
			return err
		}
	case crossSiteScriptingRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			crossSiteScriptingRuleIDDefaultAction,
			"",
			"",
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with security_rule_action (%s) on site_id (%d) %s\n", ruleID, crossSiteScriptingRuleIDDefaultAction, d.Get("site_id").(int), err)
			return err
		}
	case illegalResourceAccessRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			illegalResourceAccessRuleIDDefaultAction,
			"",
			"",
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with security_rule_action (%s) on site_id (%d) %s\n", ruleID, illegalResourceAccessRuleIDDefaultAction, d.Get("site_id").(int), err)
			return err
		}
	case remoteFileInclusionRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			remoteFileInclusionRuleIDDefaultAction,
			"",
			"",
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with security_rule_action (%s) on site_id (%d) %s\n", ruleID, remoteFileInclusionRuleIDDefaultAction, d.Get("site_id").(int), err)
			return err
		}
	case sqlInjectionRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			sqlInjectionRuleIDDefaultAction,
			"",
			"",
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with security_rule_action (%s) on site_id (%d) %s\n", ruleID, sqlInjectionRuleIDDefaultAction, d.Get("site_id").(int), err)
			return err
		}
	case ddosRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			"",
			ddosRuleIDDefaultActivationMode,
			ddosRuleIDDefaultDDOSTrafficThreshold,
			"",
			"",
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with default_activation_mode (%s) and ddos_traffic_threshold (%s) on site_id (%d) %s\n", ruleID, ddosRuleIDDefaultActivationMode, ddosRuleIDDefaultDDOSTrafficThreshold, d.Get("site_id").(int), err)
			return err
		}
	case botAccessControlRuleId:
		_, err := client.ConfigureWAFSecurityRule(
			d.Get("site_id").(int),
			ruleID,
			"",
			"",
			"",
			botAccessControlBlockBadBotsDefaultAction,
			botAccessControlChallengeSuspectedBotsDefaultAction,
		)
		if err != nil {
			log.Printf("[ERROR] Could not reset Incapsula WAF Rule rule_id (%s) with block_bad_bots (%s) and challenge_suspected_bots (%s) on site_id (%d) %s\n", ruleID, botAccessControlBlockBadBotsDefaultAction, botAccessControlChallengeSuspectedBotsDefaultAction, d.Get("site_id").(int), err)
			return err
		}
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] RESET Incapsula WAF Rule for id: %s\n", ruleID)

	return nil
}
