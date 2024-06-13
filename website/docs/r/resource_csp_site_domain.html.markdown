---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_csp_site_domain"
description: |- 
  Provides an Incapsula CSP domain resource.
---

# incapsula_csp_site_domain

Provides an Incapsula CSP domain resource.

## Example Usage

```hcl
resource "incapsula_csp_site_domain" "demo-terraform-csp-site-domain" {
  account_id          = incapsula_csp_site_configuration.example-site.account_id
  site_id             = incapsula_csp_site_configuration.example-site.site_id
  domain              = "www.imperva.com"
  status              = "allowed"
  include_subdomains  = false
  notes               = ["first note", "second note"]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to operate on.
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `domain` - (Required) The fully qualified domain name of the site. For example: `www.imperva.com`.
* `include_subdomains` - (Required) Defines Whether subdomains will inherit the allowance of the parent domain.
  Possible values: `true`, `false`
* `status` - (Optional) Defines whether the domain should be allowed or blocked once the site's mode changes to the Enforcement.
  Possible values: `allowed` (default value), `blocked`.
* `notes` -  (Optional) An array of notes for the domain, to help in future analysis and investigation. You can add as many notes as you like.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the CSP Domain.

## Import

CSP Site Configuration can be imported using the account_id, site_id and base4 encoded string of the domain separated by /, e.g.

```
$ terraform import incapsula_csp_site_domain.demo-terraform-csp-site-domain 555/1234/d3d3LmltcGVydmEuY29t
```
