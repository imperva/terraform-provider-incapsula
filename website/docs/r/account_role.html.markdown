---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_account_role"
description: |-
  Provides a Incapsula Account Role resource.
---

# incapsula_account_role

Provides an account role resource.
Each account has the option to create roles, to grant a fixed set of permissions to users. This resource enables you to create roles.

The role permissions should be added as keys (strings) and may be taken from the `incapsula_account_permissions` data source.
The `incapsula_account_permissions` data source contains the account permissions list.
To get the current list of permission in the account, use the <b>/v1/abilities/accounts/{accountId}</b> API found in the <b>v1</b> section of the
[Role Management API Definition page.](https://docs.imperva.com/bundle/cloud-application-security/page/roles-api-definition.htm)


## Example Usage

### Basic Usage - List

The basic usage is to use lists of account permissions keys.

```hcl
resource "incapsula_account_role" "role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
  description = "Sample Role Description 1"
  permissions = ["canAddSite", "canEditSite"]
}
```

### Data Sources Usage

The incapsula_account_permissions data sources provide the Account Permissions display names that are more "human-readable".

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

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to operate on - a reference to the account datasource may be used

* `name` - (Required) The role name

* `description` - (Optional) The role description

* `permissions` - (Optional) List of account permission keys

  Default value is an empty list (role with no permissions).
  `incapsula_account_permissions` data source can be used in different ways (see examples above)


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the account role.

## Import

Account Role can be imported using the `id`
```
$ terraform import incapsula_account_role.demo 1234
```