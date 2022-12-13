---
layout: "incapsula"
page_title: "Incapsula: account-default-roles"
sidebar_current: "docs-incapsula-data-account-default-roles"
description: |-
  Provides an Incapsula Account Default Roles data source.
---

# incapsula_account_default_roles

Provides the account default roles configured in every account.<p>
There are no filters needed for this data source.

## Example Usage

```hcl
data "incapsula_account_default_roles" "default_roles" {
  account_id = data.incapsula_account_data.account_data.current_account
}

resource "incapsula_account_user" "example_basic_user_1" {
  account_id = data.incapsula_account_data.account_data.current_account
  email = "example@terraform.com"
  first_name = "First"
  last_name = "Last"
  role_ids = [
    data.incapsula_account_default_roles.default_roles.admin_role_id,
    data.incapsula_account_default_roles.default_roles.reader_role_id,
  ]
}

```

## Argument Reference

There are no filters in this resource.

## Attributes Reference

The following attributes are exported:

* `admin_role_id` - Default Administrator Role Id.
* `reader_role_id` - Default Reader Role Id.