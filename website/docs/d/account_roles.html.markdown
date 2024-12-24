---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "Incapsula: account-roles"
description: |-
  Provides an Incapsula Account Roles data source.
---

# incapsula_account_roles

Provides the account roles configured in every account. Roles are used to grant a fixed set of permissions to a user.<p>
There are no filters needed for this data source.

The attribute `map` is generated and contains all the account roles, providing a Role Name to Id map.
The mapping should be used for managing users in a subaccount (when the API_KEY and API_ID are associated with a subaccount).
For managing users at the account level (where the API_KEY and API_ID are associated with the parent account), it is recommended to reference the specific role resource instead.

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
    incapsula_account_role.role_1.id,
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