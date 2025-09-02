---
subcategory: "Cloud WAF"
layout: "incapsula"
page_title: "incapsula_waf_security_rule"
description: |-
  Provides a Incapsula WAF Security Rule resource.
---

# incapsula_waf_security_rule

Provides a resource to create a subset of WAF security rules.  See, `incapsula_policy` resource for additional WAF security rule types.

## Example Usage

```hcl
resource "incapsula_waf_security_rule" "example-waf-backdoor-rule" {
  site_id = incapsula_site.example-site.id
  rule_id = "api.threats.backdoor"
  security_rule_action = "api.threats.action.quarantine_url" # (api.threats.action.quarantine_url (default) | api.threats.action.alert | api.threats.action.disabled | api.threats.action.quarantine_url)
}

resource "incapsula_waf_security_rule" "example-waf-bot-access-control-rule" {
  site_id = incapsula_site.example-site.id
  rule_id = "api.threats.bot_access_control"
  block_bad_bots = "true" # true | false (optional, default: true)
  challenge_suspected_bots = "true" # true | false (optional, default: true)
}

resource "incapsula_waf_security_rule" "example-waf-ddos-rule" {
  site_id = incapsula_site.example-site.id
  rule_id = "api.threats.ddos"
  activation_mode = "api.threats.ddos.activation_mode.on" # (api.threats.ddos.activation_mode.auto | api.threats.ddos.activation_mode.off | api.threats.ddos.activation_mode.on | api.threats.ddos.activation_mode.adaptive)
  ddos_traffic_threshold = "5000" # valid values are 10, 20, 50, 100, 200, 500, 750, 1000, 2000, 3000, 4000, 5000
  unknown_clients_challenge = "none" # valid values are none, cookies, javascript, captcha (optional, default: cookies)
  block_non_essential_bots  = "false" # true | false (optional, default: false)
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `rule_id` - (Required) The identifier of the WAF rule, e.g api.threats.cross_site_scripting. Possible values: `api.threats.bot_access_control`, `api.threats.sql_injection`, `api.threats.cross_site_scripting`, `api.threats.backdoor`, `api.threats.ddos`, `api.threats.remote_file_inclusion`.
* `security_rule_action` - (Optional) The action that should be taken when a threat is detected, for example: api.threats.action.block_ip. See above examples for `rule_id` and `action` combinations.
* `activation_mode` - (Optional) The mode of activation for ddos on a site. Possible values: api.threats.ddos.activation_mode.off, api.threats.ddos.activation_mode.auto, api.threats.ddos.activation_mode.on, api.threats.ddos.activation_mode.adaptive.
* `ddos_traffic_threshold` - (Optional) Consider site to be under DDoS if the request rate is above this threshold. The valid values are 10, 20, 50, 100, 200, 500, 750, 1000, 2000, 3000, 4000, 5000.
* `unknown_clients_challenge` - (Optional) Defines a method used for challenging suspicious bots. This argument is valid for the rule_id api.threats.ddos argument value only. If this argument is not provided, then the value stays as it is in the system. Possible values: none, cookies, javascript, captcha
* `block_non_essential_bots` - (Optional) If non-essential bots (bots determined to be legitimate by Imperva's client classification mechanism, such as site helpers and search engines) should be blocked or not. This argument is valid for the rule_id api.threats.ddos argument value only. If this argument is not provided, then the value stays as it is in the system. Possible values: true, false
* `block_bad_bots` - (Optional) Whether or not to block bad bots. Possible values: true, false.
* `challenge_suspected_bots` - (Optional) Whether or not to send a challenge to clients that are suspected to be bad bots (CAPTCHA for example). Possible values: true, false.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the WAF Security Rule.

## Import

Security Rule can be imported using the role `site_id` and `rule_id` separated by /, e.g.:

```
$ terraform import incapsula_waf_security_rule.demo site_id/rule_id
```
