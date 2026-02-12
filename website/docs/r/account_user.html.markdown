---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "incapsula_account_user"
description: |-
  Provides a Incapsula Account User resource.
---

# incapsula_account_user

Provides an account user resource.
This resource enables you to create users in an account and assign roles to them.

The user roles should be added as ids which can be taken from `incapsula_account_role` resources as reference.
In addition, the default account roles may be taken from the `incapsula_account_roles` data source.

This resource also provides the option to assign users to subaccounts.
The usage is the same but the behavior differs slightly.
See the 'SubAccount User Assignment Usage' example below for more details.


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

### Usage with Approved IPs

User creation with approved IP addresses for login restrictions.

```hcl
resource "incapsula_account_user" "user_with_ips" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  approved_ips = ["192.168.1.1", "10.0.0.5", "172.16.0.10"]
}

```

### Role References Usage

Usage with role reference.

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
  approved_ips = ["192.168.1.1", "10.0.0.0/24"]
}
```

### SubAccount User Assignment Usage - Manage by Account

For subaccounts we are not creating a new user but assigning an existing user from the parent account.
In terms of the TF resource, it means the email attribute must be taken from an existing user, by reference (preferred option) or hardcoded.
The first and last name are redundant and will be taken from the existing selected account.
The roles are not taken from the existing user and must be assigned independently.

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

### SubAccount User Assignment Usage - Manage by SubAccount

If the API_KEY and API_ID are associated with a subaccount, we don't have the option to manage roles and need to use the `incapsula_account_roles` data source.
Using the `map` attribute, this data source generates a mapping for all the account roles (Role Name to Id map).

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
* `email` - (Required) The user email. This attribute cannot be updated.
* `first_name` - (Optional) The user's first name. This attribute cannot be updated.
* `last_name` - (Optional) The user's last name. This attribute cannot be updated.
* `role_ids` - (Optional) List of role ids to be associated with the user. <p/>
  Default value is an empty list (user with no roles).
* `approved_ips` - (Optional) List of approved IP addresses from which the user is permitted to log in. <p/>
  Supports individual IPs, IP ranges, and CIDR notation. Default value is an empty list (no IP restrictions).


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the account user.

## Import

Account User can be imported using the `account_id` and `email` separated by `/`, e.g.:
```
$ terraform import incapsula_account_user.demo 1234/example@terraform.com
```