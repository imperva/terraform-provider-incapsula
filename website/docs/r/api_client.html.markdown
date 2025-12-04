---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "incapsula_api_client"
description: |-
Provides a Incapsula Account User resource.
---

# incapsula_api_client

Provides an api_client resource.
This resource enables you to create api-client (credentials for API calls) for a user.

This resource also provides the option to regenerate the API Key (the secret part of the credentials), by extending the expiration date.

## Example Usage

API-client creation for the current user.

```hcl
resource "incapsula_api_client" "api_client_1" {
  name = "First"
  description = "Last"
  expiration_date = "2026-09-12"
  enabled = true
}
```

API-client creation for another user.

```hcl
resource "incapsula_api_client" "api_client_1" {
  user_email = "example@terraform.com"
  name = "First"
  description = "Last"
  expiration_date = "2026-09-12"
  enabled = true
}
```

API-client creation for a user on another account.

```hcl
resource "incapsula_api_client" "api_client_1" {
  account_id = 1234
  user_email = "example@terraform.com"
  name = "First"
  description = "Last"
  expiration_date = "2026-09-12"
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on.
* `user_email` - (Optional) Identifier of the user on whom the api-client will be created.
* `name` - (Required) The name of the API client.
* `description` - (Optional) The description of the API client.
* `api_key` - (Computed) The API key. This attribute is a secret and will not be exposed on terraform show. to be used as a reference only.
* `enabled` - (Optional) whether the API-client is enabled or not. The default is false. 
* `expiration_date` - Expiration date of the API key (YYYY-MM-DD format only). Must be a future date. **Changing this value will cause regeneration of the key**. 
* `last_used_at` - (Computed) The last time when the api client was used. 



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
$ terraform import incapsula_api_client.demo 1234/2222
```