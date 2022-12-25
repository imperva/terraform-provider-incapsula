---
layout: "incapsula"
page_title: "Incapsula: account-permissions"
sidebar_current: "docs-incapsula-data-account-permissions"
description: |-
  Provides an Incapsula Account Permissions data source.
---

# incapsula_account_permissions

A mapping between permission identifiers and "human-readable" permission names.
Provides the ability to use the permission display names when creating and modifying user roles.<p>
To get the current list of permission display names in the account,
use the <b>/v1/abilities/accounts/{accountId}</b> API found in the <b>v1</b> section of the
[Role Management API Definition page.](https://docs.imperva.com/bundle/cloud-application-security/page/roles-api-definition.htm)

To access a subset of the permissions from the data source, use the optional filtering.
The `filter_by_text` argument is case-insensitive. When used, it generates the `keys` attribute.

The attribute `map` is always generated and contains all the account permissions (permission DisplayName to permission Key map).
Using the `map` attribute with an incorrect permission display name will cause the plan step to fail.

## Example Usage

```hcl
data "incapsula_account_permissions" "account_permissions" {
  account_id = data.incapsula_account_data.account_data.current_account
}

resource "incapsula_account_role" "role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
  description = "Sample Role Description 1"
  permissions = ["canAddSite", "canEditSite",
    data.incapsula_account_permissions.account_permissions.map["View Infra Protect settings"],
    data.incapsula_account_permissions.account_permissions.map["Delete exception from policy"],
  ]
}
```

In this example, we are using the generated `keys` attribute filtered by `filter_by_text` argument.

```hcl
data "incapsula_account_permissions" "account_permissions" {
  account_id = data.incapsula_account_data.account_data.current_account
  filter_by_text="site"
}

resource "incapsula_account_role" "role_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 2"
  description = "Sample Role Description 2"
  permissions = data.incapsula_account_permissions.account_permissions.keys
}
```

## Argument Reference

* `filter_by_text` - (Optional) string value - Filter by account permission display names.


## Attributes Reference

The following attributes are exported:

* `map` - Map of all the account permissions where the key is the permission display name and the value is the permission key.

  This attribute is always generated, even if you are using the `filter_by_text` argument

* `keys` - List of account permission keys filtered by `filter_by_text` argument.