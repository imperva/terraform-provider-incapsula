---
layout: "incapsula"
page_title: "Incapsula: incap-site-tls-settings"
sidebar_current: "docs-incapsula-resource-site-tls-settings"
description: |- Provides an Incapsula Site TLS Settings resource.
---

# incapsula_site_tls_settings

Provides an Incapsula Site TLS Settings resource.
If your site needs to support client certificates, you can upload your CA certificate to Imperva and configure your websites to use it.

## Example Usage
Associate 2 different certificates with site.

```hcl
resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_1" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate        = filebase64("./cert.der")
    certificate_name   = "der certificate example 1"
}

resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_2" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate        = filebase64("./cert.pem")
    certificate_name   = "pem certificate example 2"
}

resource "incapsula_mtls_site_client_to_imperva_certificate_association" "site_client_ca_certificate_association_1" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_1.id
    site_id            = incapsula_site.example-site.id
}

resource "incapsula_mtls_site_client_to_imperva_certificate_association" "site_client_ca_certificate_association_2" {
    account_id         = data.incapsula_account_data.account_data.current_account
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_2.id
    site_id            = incapsula_site.example-site.id
}
```

Use of `depends_on` parameter to ensure proper order of editing related resources

```
resource "incapsula_site_tls_settings" "demo_site_tls_configuration" {
    site_id                          = incapsula_site.example-site.id
    mandatory                        = true
    ports                            = [100,120,292]
    is_ports_exception               = false
    hosts                            = ["host.com", "site.ca"]
    is_hosts_exception               = true
    fingerprints                     = ["fingerprint1", "fingerprint2"]
    forward_to_origin                = false
    header_name                      = "header"
    header_value                     = "COMMON_NAME"
    is_disable_session_resumption    = true
    depends_on                       = [
        incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_1,
        incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_2,
        incapsula_mtls_site_client_to_imperva_certificate_association.site_client_ca_certificate_association_2,
        incapsula_mtls_site_client_to_imperva_certificate_association.site_client_ca_certificate_association_2
    ]
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `mandatory` - (Optional) When set to true, the end user is required to present the client certificate in order to access the site. By default, set to false.
* `ports` - (Optional) The ports on which client certificate authentication is supported. If left empty, client certificates are supported on all ports.
* `is_ports_exception` - (Optional) When set to true, client certificates are not supported on the ports listed in the Ports field ('blacklisted'). By default, set to false.
* `hosts` - (Optional) The hosts on which client certificate authentication is supported. If left empty, client certificates are supported on all hosts.
* `is_hosts_exception` - (Optional) When set to true, client certificates are not supported on the hosts listed in the Hosts field ('blacklisted'). By default, set to false.
* `fingerprints` - (Optional) Permitted client certificate fingerprints. If left empty, all fingerprints are permitted.
* `forward_to_origin` - (Optional) When set to true, the contents specified in headerValue are sent to the origin server in the header specified by headerName. By default, set to false.
* `header_name` - (Optional) The name of the header to send header content in. By default, the header name is 'clientCertificateInfo'.
* `header_value` - (Optional) The content to send in the header specified by headerName. One of the following: FULL_CERT (for full certificate in Base64) COMMON_NAME (for certificate's common name (CN)) FINGERPRINT (for the certificate fingerprints in SHA1) SERIAL_NUMBER (for the certificate's serial number).
* `is_disable_session_resumption` - (Optional)

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Site TLS Configuration.

## Import

Site TLS Settings can be imported using Site ID :

```
$ terraform import incapsula_site_tls_settings.demo_site_tls_configuration 1234

```

