layout: "incapsula" 
page_title: "Incapsula: mtls-imperva-to-origin-certificate-site-association"
sidebar_current: "docs-incapsula-mtls-imperva-to-origin-certificate-site-association"
description: |- 
  Provides an Incapsula Site to Mutual TLS Imperva to Origin Certificate Association resource.
---

# incapsula_mtls_imperva_to_origin_certificate_site_association

Provides an Incapsula Site to Mutual TLS Imperva to Origin Certificate Association resource.

## Example Usage

```hcl
resource "incapsula_mtls_imperva_to_origin_certificate" "mtls_certificate"{
  certificate       = filebase64("./cert.der")
  private_key       = filebase64("./key.der")
  passphrase        = "12345"
  certificate_name  = "pem certificate example"
}

resource "incapsula_mtls_imperva_to_origin_certificate_site_association" "site_mtls_association-site1" {
  certificate_id  =  incapsula_mtls_imperva_to_origin_certificate.mtls_certificate.id
  site_id         =  incapsula_site.example-site.id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) The account to operate on.
* `certificate_id` - (Required) The Mutual TLS Imperva to Origin Certificate ID.
* `site_id` - (Required) Numeric identifier of the site to operate on.

## Attributes Reference

The following attributes are exported:

* `id` - Incapsula Site to Mutual TLS Imperva to Origin Certificate Association. The ID composed of 2 parts: `site_id` and `certificate_id` separated by slash.

## Import

Incapsula Mutual TLS Imperva to Origin Certificate can be imported using `site_id` and `certificate_id` separated by slash:

```
$ terraform import incapsula_mtls_imperva_to_origin_certificate_site_association.site_mtls_association-site1 site_id/certificate_id
```
