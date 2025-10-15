---
subcategory: "Cloud WAF - Certificate Management"
layout: "incapsula"
page_title: "incapsula_managed_certificate_settings"
description: |- 
  Provides an Incapsula Site Managed Certificate Settings resource.
---

# incapsula_managed_certificate_settings

Provides an Incapsula Site's managed certificate settings resource.
The provider will configure or remove a managed certificate for the sites' domains, based on this resource.
<br/>

Note: This resource applies only to sites managed by the incapsula_site_v3 resource.

## Example Usage

```hcl
resource "incapsula_managed_certificate_settings" "example-managed_certificate_settings" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  default_validation_method = "CNAME"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `default_validation_method` - (Optional) The default SSL validation method that will be used for new domains. Options are `CNAME`, `DNS` and `EMAIL`. Defaults to `CNAME`.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the managed certificate settings resource.

## Import

Managed certificate settings can be imported using the site_id, e.g.:

```
$ terraform import incapsula_managed_certificate_settings.example-managed_certificate_settings site_id
```

Or by using the account_id and site_id separated by /, e.g.:

```
$ terraform import incapsula_managed_certificate_settings.example-managed_certificate_settings account_id/site_id
```

The official docs for Manage Certificate settings API are located here: https://docs.imperva.com/bundle/cloud-application-security/page/certificatesUI-api-definition.htm


