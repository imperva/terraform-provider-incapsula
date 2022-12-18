---
layout: "incapsula"
page_title: "Incapsula: account-user"
sidebar_current: "docs-incapsula-resource-account-user"
description: |-
Provides a Incapsula Account User resource.
---

# incapsula_account_user

Provides an Account User resource.
Each Account has the option to create users with roles assigned to it. This resource allows you to add them.

The user roles should be added as ids and may be taken from `incapsula_account_role` resources as reference.
In addition, the default account roles may be taken from `incapsula_account_roles` datasource.
This data source contains the Administrator and Reader default role ids.

This resource give also the option to assign users to SubAccounts, the usage is the same but the behavior differ a bit.
Please look at the 'SubAccount User Assignment Usage' example for more details


## Example Usage

### Basic Usage

Sample user creation with no roles.

```hcl
resource "incapsula_account_user" "user_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
}

```

### Role References Usage

The usage is with role reference.

```hcl
resource "incapsula_account_role" "role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
}
resource "incapsula_account_role" "role_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 2"
}

resource "incapsula_account_user" "user_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  role_ids = [
    incapsula_account_role.role_1.id,
    incapsula_account_role.role_2.id,
  ]
}
```

### Role References & Data Sources Usage

Using `incapsula_account_roles` data sources we can use admin_role_id/reader_role_id exported attributes.

```hcl
data "incapsula_account_roles" "roles" {
  account_id = data.incapsula_account_data.account_data.current_account
}
resource "incapsula_account_role" "role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
}
resource "incapsula_account_role" "role_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 2"
}

resource "incapsula_account_user" "user_3" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  role_ids = [
    incapsula_account_role.role_1.id,
    incapsula_account_role.role_2.id,
    data.incapsula_account_roles.roles.reader_role_id,
  ]
}
```

### SubAccount User Assignment Usage Manage by Account

For SubAccounts we are talking about assignments so the user should exist in the parent account.</p>
In terms of resource, it means the email attribute must be taken from an existing user, hardcoded or by reference (preferred option).
The first and last name are redundant and then, ignored and taken from the existing chosen account.
The roles have to be chosen independently, they are not coming from the existing user.

```hcl
resource "incapsula_account_role" "role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
}

resource "incapsula_account_user" "user_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
}

resource "incapsula_account_user" "user_2" {
  account_id = incapsula_subaccount.example-subaccount.id
  email = incapsula_account_user.user_1.email
  role_ids = [
    incapsula_account_role.role_1.id,
  ]
}
```

### SubAccount User Assignment Usage Manage by SubAccount

if the API_KEY and API_ID are associated to a SubAccount, we won't have the option to manage roles, so we should use `incapsula_account_roles` datasource.
This datasource have a `map` attribute generated that contains all the Account Roles (Role Name to Id map).

```hcl
data "incapsula_account_roles" "roles" {
  account_id = data.incapsula_account_data.account_data.current_account
}

resource "incapsula_account_user" "user_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  role_ids = [
    data.incapsula_account_roles.roles.map["Sample Role 1"],
  ]
}
```


## Argument Reference

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to operate on. <p/>
  Using reference to account datasource
* `email` - (Required) The user email
* `first_name` - (Optional) The user first name
* `last_name` - (Optional) The user last name
* `role_ids` - (Optional) List of role ids to be associated to the user. <p/>
  Default value is an empty list (user with no roles).


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the account user.

## Import

Account User can be imported using the `id`
```
$ terraform import incapsula_account_user.demo 1234
```