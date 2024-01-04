---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_csp_site_configuration"
description: |- 
  Provides an Incapsula CSP site configuration resource.
---

# incapsula_csp_site_configuration

Provides an Incapsula CSP site configuration resource.

## Example Usage

```hcl
resource "incapsula_csp_site_configuration" "demo-terraform-csp-site-configuration" {
  account_id      = incapsula_site.example-site.account_id
  site_id         = incapsula_site.example-site.id
  mode            = "monitor"
  email_addresses = [ "test@imperva.com" ]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to operate on.
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `mode` - (Optional) Website Protection Mode. When in "enforce" mode, blocked resources will not be available in the application and new resources will be automatically blocked. When in "monitor" mode, all resources are available in the application and the system keeps track of all new domains that are discovered.
  Possible values: `monitor` (default value), `enforce`, `off`.
* `email_addresses` -  (Optional) An array of email address for the event notification recipient list of a specific website. Notifications are reasonably small and limited in frequency.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the CSP Site Configuration.

## Import

CSP Site Configuration can be imported using the account_id and site_id separated by /, e.g.

```
$ terraform import incapsula_csp_site_configuration.demo-terraform-csp-site-configuration 555/1234
```
