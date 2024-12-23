---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_policy"
description: |-
  Provides a Incapsula Policy resource.
---

# incapsula_policy

Provides a resource to define WAF security, Whitelist, and ACL policies. All policies are created at the parent account level. 

The follow-on action is to use the `incapsula_account_policy_association` resource, to assign the policy to a sub account.

To simplify the use of policies, you can utilize this [cloud-waf Module](https://registry.terraform.io/modules/imperva/cloud-waf/incapsula/latest) along with its submodules.

For full feature documentation, see [Create and Manage Policies](https://docs.imperva.com/bundle/cloud-application-security/page/policies.htm).

## Example Usage

```hcl

resource "incapsula_policy" "example-whitelist-ip-policy" {
    name        = "Example WHITELIST IP Policy"
    enabled     = true 
    policy_type = "WHITELIST"
    description = "Example WHITELIST IP Policy description"
    policy_settings = <<POLICY
    [
      {
        "settingsAction": "ALLOW",
        "policySettingType": "IP",
        "data": {
          "ips": [
            "1.2.3.4"
          ]
        }
      }
    ]
    POLICY
}

resource "incapsula_policy" "example-acl-country-block-policy" {
    description     = "EXAMPLE ACL Block Countries based on attack."
    enabled         = true
    policy_type = "ACL"
    name            = var.dynamic_country_block_policy_name
    policy_settings = jsonencode(
        [
            {
                data = {
                    geo = {
                        countries = var.countries
                    }
                }
                policySettingType = "GEO"
                settingsAction    = "BLOCK"
            },
        ]
    )
}

resource "incapsula_policy" "example-waf-rule-illegal-resource-access-policy" {
    name        = "Example WAF-RULE ILLEGAL RESOURCE ACCESS Policy"
    enabled     = true 
    policy_type = "WAF_RULES"
    policy_settings = <<POLICY
    [
    {
      "settingsAction": "BLOCK",
      "policySettingType": "REMOTE_FILE_INCLUSION"

    },
    {
      "settingsAction": "BLOCK",
      "policySettingType": "ILLEGAL_RESOURCE_ACCESS",
      "policyDataExceptions": [
        {
          "data": [
            {
              "exceptionType": "URL",
              "values": [
                "/cmd.exe"
              ]
            },
            {
              "exceptionType": "SITE_ID",
              "values": [
              "132456789"
              ]
            }
          ]
        }
      ]
    },
    {
      "settingsAction": "BLOCK",
      "policySettingType": "CROSS_SITE_SCRIPTING"
      
    },
    {
      "settingsAction": "BLOCK",
      "policySettingType": "SQL_INJECTION"
      
    }
    ]
    POLICY
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The policy name.
* `enabled` - (Required) Enables the policy.
* `policy_type` - (Required) The policy type. Possible values: ACL, WHITELIST, WAF_RULES.  Note: For (policy_type=WAF_RULES), all 4 setting types (policySettingType) are mandatory (REMOTE_FILE_INCLUSION, ILLEGAL_RESOURCE_ACCESS, CROSS_SITE_SCRIPTING, SQL_INJECTION).
* `policy_settings` - (Required) The policy settings as JSON string. See Imperva documentation for help with constructing a correct value.
Policy_settings internal values:
policySettingType: IP, GEO, URL
settingsAction: BLOCK, ALLOW, ALERT, BLOCK_USER, BLOCK_IP, IGNORE
policySettings.data.url.pattern: CONTAINS, EQUALS, NOT_CONTAINS, NOT_EQUALS, NOT_PREFIX, NOT_SUFFIX, PREFIX, SUFFIX 
exceptionType: GEO, IP, URL, CLIENT_ID, SITE_ID
* `description` - (Optional) The policy description.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the policy.
* `account_id` - Account ID of the policy.

## Import

Policy can be imported using the `id`, e.g.:

```
$ terraform import incapsula_policy.demo 1234
```