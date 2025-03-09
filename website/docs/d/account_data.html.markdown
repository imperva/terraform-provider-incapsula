---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "Incapsula: account-data"
description: |-
    Provides an Incapsula Account Data data source.
---

# incapsula_account_data

There are no filters needed for this data source

## Example Usage


```hcl
data "incapsula_account_data" "account_data" {
}

resource "incapsula_account_policy_association" "example-account-policy-association-parent" {
    account_id                       = data.incapsula_account_data.account_data.current_account
    default_non_mandatory_policy_ids = [
        "123456"
    ]
    default_waf_policy_id            = "789012"
}
```

## Argument Reference

There are no filters in this resource.

## Attributes Reference

The following attributes are exported:

* `current_account` - Current account ID.
* `plan_name` - Plan name.