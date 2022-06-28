---
layout: "incapsula"
page_title: "Incapsula: account-data"
sidebar_current: "docs-incapsula-data-account-data"
description: |-
Provides an Incapsula Account Data data source.
---

# incapsula_account_data

There are no filters needed for this data source

## Example Usage


```hcl
data "incapsula_account_data" "account_data" {
}

# Policy: Use reference to current_account field of incapsula_account_data data source in "default_website_accounts"
and "available_for_accounts"
resource "incapsula_policy" "example-whitelist-ip-policy" {
    name        = "Example WHITELIST IP Policy"
    enabled     = true 
    policy_type = "WHITELIST"
    description = "Example WHITELIST IP Policy description"
    default_website_accounts = [
        data.incapsula_account_data.account_data.current_account,
        "1234"
    ]
    available_for_accounts = [
        data.incapsula_account_data.account_data.current_account,
        "1234",
        "5678"
    ]
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

There are no filters in this resource.

## Attributes Reference

The following attributes are exported:

* `current_account` - Current account ID.
* `plan_name` - Plan name.