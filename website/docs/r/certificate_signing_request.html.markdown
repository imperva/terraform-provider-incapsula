---
subcategory: "Cloud WAF - Certificate Management"
layout: "incapsula"
page_title: "incapsula_certificate_signing_request"
description: |-
  Provides a Incapsula Certificate Signing Request resource.
---

# incapsula_certificate_signing_request

Provides a Incapsula Certificate Signing Request resource. 

## Example Usage

```hcl
resource "incapsula_certificate_signing_request" "certificate-signing-request" {
    site_id           = incapsula_site.example-site.id
    domain            = "sandwich.au"
    email             = "test@sandwich.au"
    country           = "AU"
    state             = "QLD"
    city              = "BNE"
    organization      = "Tacos Pty Ltd"
    organization_unit = "Kitchen"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `domain` - (Optional) The common name. For example: `example.com`
* `email` - (Optional) The Email address. For example: `joe@example.com`
* `country` - (Optional) The two-letter ISO code for the country where your organization is located.
* `state` - (Optional) The state/region where your organization is located. This should not be abbreviated.
* `city` - (Optional) The city where your organization is located.
* `organization` - (Optional) The legal name of your organization. This should not be abbreviated or include suffixes such as Inc., Corp., or LLC.
* `organization_unit` - (Optional) The division of your organization handling the certificate. For example, IT Department.

## Attributes Reference

The following attributes are exported:

* `id` - (String) At the moment, only one active certificate can be stored. This exported value is always set to `site_id`.
* `csr_content` - (String) The certificate request data.
