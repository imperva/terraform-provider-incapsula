---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "Incapsula: account"
description: |-
    Provides an Incapsula Account data source that finds an account by its name.
---

# incapsula_account

Finds an account below the current account by its name, so there is no need to hardcode account IDs.

* When using reseller credentials, the reseller's direct accounts are searched.
* When using parent account credentials, its sub accounts are searched.

The credentials must have one of the abilities `canAccessSubAccounts`, `canManageAccountSubAccounts` or `canExportConfigurations`.

## Example Usage

```hcl
data "incapsula_account" "example" {
    account_name = "My Sub Account"
}

resource "incapsula_site" "site" {
    account_id  = data.incapsula_account.example.id
    domain      = "www.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) The name of the account to find. Matching is case-insensitive; the value is not trimmed. The data source fails if no account matches, or if more than one account matches (the matching account IDs are listed in the error).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the account.
* `account_id` - The ID of the account.
* `account_name` - The name of the account, as returned by the API.
* `parent_id` - The ID of the parent account.
