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

* `account_id` - (Optional) Numeric identifier of the account to operate on.
* `name` - (Required) The name of the API client.
* `description` - (Optional) The description of the API client.
* `apiKey` - (Computed) The API key. This attribute is a secret and will not be exposed on terraform show. to be used as a reference only.
* `enabled` - (Optional) wether the api client is enabled or not. default is false. 
* `regenerate` - (Optional) when true and expiration period was set to a future date - the api key will be regenerated. 
* `lastActionTime` - (Computed) The last time when the api client was used. 
* `expiration_period` - (required when regenerate is true) a future date when the api client will be expired. 
* `gracePeriodInSeconds` - (Optional) The period of time when the old key will be still enabled after a new key was regenerated. 


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the account user.

## Import

Api Client can be imported using the `API ID`, e.g.:
```
$ terraform import incapsula_api_client.demo 2222
```

OR, it can be imported using the `account_id` and `API ID` separated by `/`, e.g.:
```
$ terraform import incapsula_account_user.demo 1234/2222
```