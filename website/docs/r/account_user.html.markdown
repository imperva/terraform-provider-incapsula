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
In addition, the default account roles may be taken from `incapsula_account_default_roles` datasource.
This data source contains the Administrator and Reader default role ids.

## Example Usage

### Basic Usage

Sample user creation with no roles.

```hcl
resource "incapsula_account_user" "example_basic_user_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
}

```

### Role References Usage

The usage is with role reference.

```hcl
resource "incapsula_account_role" "example_basic_role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
}
resource "incapsula_account_role" "example_basic_role_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 2"
}

resource "incapsula_account_user" "example_basic_user_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  role_ids = [
    incapsula_account_role.example_basic_role_1.id,
    incapsula_account_role.example_basic_role_2.id,
  ]
}
```

### Role References & Data Sources Usage

Using `incapsula_account_default_roles` data sources we can use admin_role_id/reader_role_id exported attributes.

```hcl
data "incapsula_account_default_roles" "default_roles" {
  account_id = data.incapsula_account_data.account_data.current_account
}
resource "incapsula_account_role" "example_basic_role_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 1"
}
resource "incapsula_account_role" "example_basic_role_2" {
  account_id = data.incapsula_account_data.account_data.current_account
  name = "Sample Role 2"
}

resource "incapsula_account_user" "example_basic_user_3" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  role_ids = [
    incapsula_account_role.example_basic_role_1.id,
    incapsula_account_role.example_basic_role_2.id,
    data.incapsula_account_default_roles.default_roles.reader_role_id,
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