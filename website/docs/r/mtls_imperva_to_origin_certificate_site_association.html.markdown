---
subcategory: "Cloud WAF - Certificate Management"
layout: "incapsula"
page_title: "incapsula_mtls_imperva_to_origin_certificate"
description: |-
  Provides a Mutual TLS Imperva to Origin certificate resource.
---

# incapsula_mtls_imperva_to_origin_certificate_site_association

Provides a Mutual TLS Imperva to Origin certificate Association resource.
This resource is used to associate between mTLS client certificates and site.

## Example Usage

```hcl
resource "incapsula_mtls_imperva_to_origin_certificate_site_association" "mtls_imperva_to_origin_certificate_site_association"{
  site_id       = incapsula_site.example-site.id
  certificate_id  = incapsula_certificate.example-certificate.id
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Site id to assign to a given mTLS client certificate.
* `certificate_id` - (Required) The mTLS certificate id you want to assign to your site.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Mutual TLS Imperva to Origin Certificate.