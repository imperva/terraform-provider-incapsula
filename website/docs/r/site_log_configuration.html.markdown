---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_site_log_configuration"
description: |-
    Provides an Incapsula Site Log Configuration resource.
---
# incapsula_site_log_configuration

Provides an Incapsula Site Log Configuration resource.

Note: This resource applies only to sites managed by the incapsula_site_v3 resource. For sites managed by the incapsula_site resource, please use the relevant fields from the incapsula_site resource instead.

In this resource, you can configure:
- Log Level: The log level for the site. Options include security logs only, both security and access logs, or no logging.
- Logs Account ID: The account ID that collects the logs.
- Data Storage Region: The region where the data is stored.

Note: This resource is designed to work with sites represented by the "incapsula_site_v3" resource only.

## Example Usage

```hcl

resource "incapsula_site_log_configuration" "example" {
  site_id             = "1234341"
  logs_account_id     = "67890"
  log_level           = "full"
  data_storage_region = "US"
  hashing_enabled     = true
  hash_salt           = "EJKHRT48375N4TKE7956NG"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `logs_account_id` - (Optional) Numeric identifier of the account that collects the logs.
* `log_level` - (Optional) The log level options are `full`, `security`, and `none`. Full logging includes both security and access logs.
* `data_storage_region` - (Optional) The data region to use. Options are `APAC`, `AU`, `EU`, and `US`.
* `hashing_enabled` - (Optional) Use the hashing method for masking fields in your logs and in the Security Events page, instead of the default (XXX) data masking.
* `hash_salt` - (Optional) Hashing salt to use for the hashing process. Required if hashing is enabled. Maximum length of 64 characters.
## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Site Log Configuration. The id is identical to Site id.

## Dependencies
The log_level attribute is specific to logs for the Cloud WAF service.

To configure the log level, you need the `incapsula_siem_log_configuration` resource with logs configured for the Cloud WAF service.

## Import

Site Log Configuration can be imported using the `site_id`:
```
$ terraform import incapsula_site_log_configuraton.demo 1234
```