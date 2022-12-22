---
layout: "incapsula"
subcategory: Roles & User Management
page_title: "Incapsula: account-permissions"
sidebar_current: "docs-incapsula-data-account-permissions"
description: |-
  Provides an Incapsula Account Permissions data source.
---

# incapsula_account_permissions

Provides the ability to use "human-readable" strings identifying the account permissions for roles.<p>
In order to get the latest list, use the <b>/v1/abilities/accounts/{accountId}</b> API found in the <b>v1</b> section of the
[Role Management API Definition page.](https://docs.imperva.com/bundle/cloud-application-security/page/roles-api-definition.htm)

Filtering is optional. When used, it will generate the `keys` attribute.
`filter_by_text` argument is case-insensitive.

The attribute `map` is always generated and contain all the Account Permissions (DisplayName to Key map).
Using map access with a wrong Account Permissions display name, it will fail during plan process.

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

* `filter_by_text` - (Optional) string value - Filter by Account Permissions display names.


## Attributes Reference

The following attributes are exported:

* `map` - Map of all the Account Permissions where the key is the display names and value is the key.

  This attribute is always generated, even if you are using `filter_by_text` argument.

* `keys` - List of Account Permissions keys filtered by `filter_by_text` argument.