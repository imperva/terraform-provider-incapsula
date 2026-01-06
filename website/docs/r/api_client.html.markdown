---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "incapsula_api_client"
description: |-
Provides a Incapsula Account User resource.
---

# incapsula_api_client

Provides an api_client resource. 
This resource lets you create an API client (credentials for API calls) for a user.

You can also regenerate the API key secret by updating the expiration date, which causes a new secret to be issued.

## Example Usage

API client creation for the current user.

```hcl
resource "incapsula_api_client" "api_client_1" {
  name = "First"
  description = "Last"
  expiration_date = "2026-09-12T10:33:29Z"
  enabled = true
}
```

API client creation for another user.

```hcl
resource "incapsula_api_client" "api_client_1" {
  user_email = "example@terraform.com"
  name = "First"
  description = "Last"
  expiration_date = "2026-09-12T10:33:29Z"
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
  expiration_date = "2026-09-12T10:33:29Z"
  enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on.
* `user_email` - (Optional) Email address of the user for whom the API client will be created.
* `name` - (Required) The name of the API client.
* `description` - (Optional) The description of the API client.
* `api_key` - (Computed) The API key secret. This attribute is sensitive and is not exposed in terraform show. Use it for reference only.
* `enabled` - (Optional) Whether the API client is enabled. The default is false. 
* `expiration_date` - (Optional) Expiration date of the API key, in YYYY-MM-DD format. Must be a future date. Changing this value causes regeneration of the API key secret. 
* `last_used_at` - (Computed) The last time the API client was used. 



## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the account user.



## Import

An API client can be imported using the API ID, for example:
```
$ terraform import incapsula_api_client.demo 2222
```

Alternatively, import it using the account_id and API ID, separated by a slash (/):
```
$ terraform import incapsula_api_client.demo 1234/2222
```