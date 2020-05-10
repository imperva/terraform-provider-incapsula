---
layout: "incapsula"
page_title: "Incapsula: site"
sidebar_current: "docs-incapsula-resource-site"
description: |-
  Provides a Incapsula Site resource.
---

# incapsula_site

Provides a Incapsula Site resource. 
Sites are the core resource that is required by all other resources.

## Example Usage

```hcl
resource "incapsula_site" "example-site" {
  domain                 = "examplesite.com"
  account_id             = "123"
  ref_id                 = "123"
  send_site_setup_emails = "false"
  site_ip                = "2.2.2.2"
  force_ssl              = "false"
  logs_account_id        = "456"
  data_storage_region    = "US"
  hashing_enabled        = true
  hash_salt              = "foobar"
  log_level              = "full"
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The fully qualified domain name of the site. For example: www.example.com, hello.example.com.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `send_site_setup_emails` - (Optional) If this value is false, end users will not get emails about the add site process such as DNS instructions and SSL setup.
* `site_ip` - (Optional) The web server IP/CNAME.
* `force_ssl` - (Optional) Force SSL. This option is only available for sites with manually configured IP/CNAME and for specific accounts.
* `logs_account_id` - (Optional) Account where logs should be stored. Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.
* `data_storage_region` - (Optional) The data region to use. Options are `APAC`, `AU`, `EU`, and `US`.
* `hashing_enabled` - (Optional) Specify if hashing (masking setting) should be enabled.
* `hash_salt` - (Optional) Specify the hash salt (masking setting), required if hashing is enabled. Maximum length of 64 characters.
* `log_level` - (Optional) The log level. Options are `full`, `security`, and `none`. Defaults to `none`.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the site.
* `site_creation_date` - Numeric representation of the site creation date.
* `dns_cname_record_name` - The CNAME record name.
* `dns_cname_record_value` - The CNAME record value.
* `dns_a_record_name` - The A record name.
* `dns_a_record_value` - The A record value.
* `domain_verification` - The domain verification (e.g. GlobalSign verification).
