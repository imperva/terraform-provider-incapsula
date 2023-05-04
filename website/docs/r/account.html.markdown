---
layout: "incapsula"
page_title: "Incapsula: account"
sidebar_current: "docs-incapsula-resource-account"
description: |-
  Provides a Incapsula Account resource.
---

# incapsula_account

Provides a Incapsula Account resource. 

## Example Usage

```hcl
resource "incapsula_account" "example-account" {
  email                              = "example@example.com"
  parent_id                          = 123
  ref_id                             = "123"
  user_name                          = "John Doe"
  plan_id                            = "ent100"
  account_name                       = "Example Account"
  logs_account_id                    = "456"
  log_level                          = "full"

  data_storage_region                = "US"

  # Base64 Encoded HTML
  error_page_template                = "RlP5QhsBHAECGUVDFxYZVCQFBwkDBggLBA0MFB0cGhsYFTgCIgUgJx3EG8LuM6ZpqwR8ScEztVwTqbxuB8..."
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) Email address of the account admin. For example: joe@example.com.
* `parent_id` - (Optional) The newly created account's parent id. If not specified, the invoking account will be assigned as the parent.
* `ref_id` - (Optional) Customer specific identifier for this operation.
* `user_name` - (Optional) The account owner's name. For example: John Doe.
* `plan_id` - (Optional) An identifier of the plan to assign to the new account. For example, ent100 for the Enterprise 100 plan (values can be provided by your account manager).
* `account_name` - (Optional) Account name.
* `logs_account_id` - (Optional) Account where logs should be stored. Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.
* `log_level` - (Optional) The log level. Options are `full`, `security`, and `none`.
* `data_storage_region` - (Optional) Default data region of the account for newly created sites. Options are `APAC`, `EU`, `US` and `AU`. Defaults to `US`.
* `support_all_tls_versions` - (Optional) Allow sites in the account to support all TLS versions for connectivity between clients (visitors) and the Imperva service.  
                               Note: This argument is deprecated. Use add_naked_domain_san_for_www_sites in the account_ssl_settings resource instead.  
* `naked_domain_san_for_new_www_sites` - (Optional) Add naked domain SAN to Incapsula SSL certificates for new www sites. Options are `true` and `false`. Defaults to `true`.  
                                         Note: This argument is deprecated. Use add_naked_domain_san_for_www_sites in the account_ssl_settings resource instead.
* `wildcard_san_for_new_sites` - (Optional) Add wildcard SAN to Incapsula SSL certificates for new sites. Options are `true`, `false` and `default`. Defaults to `default`.  
                               Note: This argument is deprecated. Use use_wild_card_san_instead_of_fqdn in the account_ssl_settings resource instead.
* `error_page_template` - (Optional) Base64 encoded template for an error page.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the account.
* `trial_end_date` - Numeric representation of the site creation date.
* `support_level` - The CNAME record name.
* `plan_name` - The CNAME record value.

## Import

Account can be imported using the `id`, e.g.:

```
$ terraform import incapsula_account.demo 1234
```