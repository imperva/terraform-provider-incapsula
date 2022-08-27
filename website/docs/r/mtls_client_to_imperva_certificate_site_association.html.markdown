layout: "incapsula"
page_title: "Incapsula: site-mtls-client-to-imperva-certificate-association"
sidebar_current: "docs-incapsula-site-mtls-client-to-imperva-certificate-association"
description: |-
Provides an Incapsula Site to TLS Client to Imperva CA Certificate Association resource.
---

# incapsula_site_client_to_imperva_certificate_association

Provides an Incapsula Site to TLS Client to Imperva CA Certificate Association resource.

## Example Usage

```hcl
data "incapsula_account_data" "account_data" {
}

resource "incapsula_site" "example-site-1" {
  domain                 = "examplesite.com"
}

resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_1" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate        = filebase64("./cert1.der")
    certificate_name   = "der certificate example 1"
}

resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_2" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate        = filebase64("./cert2.pfx")
    certificate_name   = "pfx certificate example 2"
}

resource "incapsula_mtls_client_to_imperva_certificate_site_association" "site_certificate_association_1" {
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_1.id
    site_id            = incapsula_site.example-site-1.id
}

resource "incapsula_mtls_client_to_imperva_certificate_site_association" "site_certificate_association_2" {
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_2.id
    site_id            = incapsula_site.example-site-1.id
}
```

## Argument Reference

The following arguments are supported:

* `certificate_id` - (Required) The Mutual TLS Client to Imperva CA Certificate ID.
* `site_id` - (Required) Numeric identifier of the site to operate on.

## Attributes Reference

The following attributes are exported:

* `id` - Incapsula Site to Client to Imperva CA Association ID. The ID composed of 2 parts: `site_id` and `certificate_id` separated by slash.

## Import

Incapsula Mutual TLS Imperva to Origin Certificate can be imported using `site_id` and `certificate_id` separated by slash:

```
$ terraform import incapsula_mtls_client_to_imperva_certificate_site_association.site_certificate_association_1 site_id/certificate_id
```

