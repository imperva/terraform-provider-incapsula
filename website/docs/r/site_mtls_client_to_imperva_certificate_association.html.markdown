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

resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_1" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate        = filebase64("./cert.der")
    certificate_name   = "der certificate example 1"
}

resource "incapsula_mtls_site_client_to_imperva_certificate_association" "site_certificate_association_1" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_1.id
    site_id            = incapsula_site.example-site.id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) The account to operate on.
* `certificate_id` - (Required) The Mutual TLS Client to Imperva CA Certificate ID.
* `site_id` - (Required) Numeric identifier of the site to operate on.

## Attributes Reference

The following attributes are exported:

* `id` - Incapsula Site to Client to Imperva CA Association ID. The ID composed of 3 parts: `account_id`, `site_id` and `certificate_id` separated by slash.

## Import

Incapsula Mutual TLS Imperva to Origin Certificate can be imported using `account_id`, `site_id` and `certificate_id` separated by slash:

```
$ terraform import incapsula_mtls_site_client_to_imperva_certificate_association.site_certificate_association_1 account_id/site_id/certificate_id
```

