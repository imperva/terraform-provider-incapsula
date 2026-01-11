---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "incapsula_api_client"
description: |-
Provides a Incapsula Account User resource.
---

# incapsula_api_client

Provides an API client resource that manages API credentials (API ID and API key) for a user.

You can also regenerate the API key secret by updating the expiration date, which causes a new secret to be issued.

Note: Before creating API clients through this resource, an account administrator must explicitly enable API Client Creation Consent in the account settings in the Cloud Security Console.

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
* `expiration_date` - (Optional) Expiration date of the API key, in YYYY-MM-DDTHH:mm:ssZ format. Must be a future date. Changing this value causes regeneration of the API key secret.


## Attributes Reference

The following attributes are exported:

* `id` - The API client ID (API ID), which is the unique identifier returned by the API.
  Together with the API key (api_key), this value forms the API credentials used to authenticate API calls.
  You can use this value to import an existing API client into Terraform.



## Import

An API client can be imported using its API ID, which is the value exported as the `id` attribute for this resource.
For example:
```
$ terraform import incapsula_api_client.demo 2222
```

Alternatively, import it using the account ID and the API ID (the value exported as the `id` attribute), separated by a slash (/):
```
$ terraform import incapsula_api_client.demo 1234/2222
```