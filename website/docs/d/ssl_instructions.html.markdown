---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_ssl_instructions"
description: |- 
  Provides an Incapsula DNS and SSL instructions.
---
# incapsula_ssl_instructions

Provides an Incapsula Site SSL instruction resource.

This data resource enables you to retrieve instructions for configuring your DNS and SSL for completing the domain validation process.



## Example Usage

```hcl
data "incapsula_ssl_instructions" "example"  {
  site_id       = incapsula_site.mysite.id
  domain_ids    = [incapsula_domain.my_domain1.id,incapsula_domain.my_domain2.id]
  managed_certificate_settings_id = incapsula_managed_certificate_settings.my_managed_certificate_settings.id
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `domain_ids` - (Required) Numeric identifiers of the domains of the site to which the settings will be applied.
  - Type: `list` of `int`
* `managed_certificate_settings_id` - (Required): Numeric identifier of the managed certificate settings related to the domains.
  - Type: `int`


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Site SSL settings. The id is identical to Site id.
* `instructions` - The SSL settings instructions for the domain. It will return a set of instructions for the domains. Each instruction will contain the following fields:
  - `domain_id` - The domain id.
  - `san_id` - The SAN id.
  - `name` - the domain name to add to the DNS server
  - `type` - The record type used for the instructions. e.g TXT.
  - `value` - The certificate verification code