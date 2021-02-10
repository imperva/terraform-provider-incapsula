---
layout: "incapsula"
page_title: "Incapsula: policy"
sidebar_current: "docs-incapsula-resource-policy"
description: |-
  Provides a Incapsula Policy resource.
---

# incapsula_policy

Provides a Incapsula Policy resource. 

## Example Usage

```hcl
# policy_settings internal values:
# policySettingType: IP, GEO, URL
# settingsAction: BLOCK, ALLOW, ALERT, BLOCK_USER, BLOCK_IP, IGNORE
# policySettings.data.url.pattern: CONTAINS, EQUALS, NOT_CONTAINS, NOT_EQUALS, NOT_PREFIX, NOT_SUFFIX, PREFIX, SUFFIX
# exceptionType: GEO, IP, URL, CLIENT_ID, SITE_ID
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The policy name.
* `enabled` - (Required) Enables the policy.
* `policy_type` - (Required) The policy type. Possible values: ACL, WHITELIST.
* `policy_settings` - (Required) The policy settings as JSON string. See Imperva documentation for help with constructing a correct value.
* `account_id` - (Optional) Account ID of the policy.
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