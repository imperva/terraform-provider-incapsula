---
layout: "incapsula"
page_title: "Incapsula: account-roles"
sidebar_current: "docs-incapsula-data-account-roles"
description: |-
  Provides an Incapsula Account Roles data source.
---

# incapsula_account_roles

Provides the account roles configured in every account.<p>
There are no filters needed for this data source.

The attribute `map` is generated and contain all the Account Roles (Role Name to Id map).
The map should be use for user creation managed by SubAccount (API_KEY and API_ID associated to a SubAccount).
In case of account management (API_KEY and API_ID associated to an Account), it's recommended to use role references.

## Example Usage

```hcl
data "incapsula_account_roles" "roles" {
  account_id = data.incapsula_account_data.account_data.current_account
}

resource "incapsula_account_user" "user_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  role_ids = [
    data.incapsula_account_roles.roles.admin_role_id,
    data.incapsula_account_roles.roles.reader_role_id,
    data.incapsula_account_roles.roles.map["Sample Role 1"],
  ]
}

```

## Argument Reference

There are no filters in this resource.

## Attributes Reference

The following attributes are exported:

* `admin_role_id` - Default Administrator Role Id.
* `reader_role_id` - Default Reader Role Id.
* `map` - Map of all the Account Roles where the key is the Role Name and value is the Role Id.