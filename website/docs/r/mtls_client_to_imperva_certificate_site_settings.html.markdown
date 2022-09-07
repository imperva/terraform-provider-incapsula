---
layout: "incapsula"
page_title: "Incapsula: incap-mtls_client_to_imperva_certificate_site_settings"
sidebar_current: "docs-incapsula-resource-mtls-client-to-imperva-certificate-site-settings"
description: |- Provides an Incapsula MTLS Client to Imperva Certificate Site Settings resource.
---

# incapsula_mtls_client_to_imperva_certificate_site_settings

Provides an Incapsula mTLS Client to Imperva Certificate Site Settings resource.
If your site needs to support client certificates, you can upload your CA certificate to Imperva and configure your websites to use it.

## Example Usage
Associate 2 different certificates to the same site.
Use of `depends_on` parameter to ensure proper order of editing related resources

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

resource "incapsula_mtls_client_to_imperva_certificate_site_association" "client_ca_certificate__site_association_1" {
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_1.id
    site_id            = incapsula_site.example-site.id
}

resource "incapsula_mtls_client_to_imperva_certificate_site_association" "site_client_ca_certificate_site_association_2" {
    certificate_id     = incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_2.id
    site_id            = incapsula_site.example-site.id
}
```

Use of `depends_on` parameter to ensure proper order of editing related resources

```hcl
resource "incapsula_mtls_client_to_imperva_certificate_site_settings" "demo_client_to_imperva_certificate_site_settings" {
    site_id                          = incapsula_site.example-site.id
    require_client_certificate       = true
    ports                            = [100,120,292]
    is_ports_exception               = false
    hosts                            = ["host.com", "site.ca"]
    is_hosts_exception               = true
    fingerprints                     = ["fingerprint1", "fingerprint2"]
    forward_to_origin                = false
    header_name                      = "header"
    header_value                     = "SERIAL_NUMBER"
    is_disable_session_resumption    = true
    depends_on                       = [
        incapsula_mtls_client_to_imperva_certificate_site_association.client_ca_certificate_site_association_1,
        incapsula_mtls_client_to_imperva_certificate_site_association.client_ca_certificate_site_association_2
    ]
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `require_client_certificate` - (Optional) When set to true, the end user is required to present the client certificate in order to access the site. Default - false.
* `ports` - (Optional) The ports on which client certificate authentication is supported. If left empty, client certificates are supported on all ports. Default: empty list.
* `is_ports_exception` - (Optional) When set to true, client certificates are not supported on the ports listed in the Ports field ('blacklisted'). Default - false.
* `hosts` - (Optional) The hosts on which client certificate authentication is supported. If left empty, client certificates are supported on all hosts. Default - empty list.
* `is_hosts_exception` - (Optional) When set to true, client certificates are not supported on the hosts listed in the Hosts field ('blacklisted'). Default - false.
* `fingerprints` - (Optional) Permitted client certificate fingerprints. If left empty, all fingerprints are permitted. Default - empty list.
* `forward_to_origin` - (Optional) When set to true, the contents specified in headerValue are sent to the origin server in the header specified by headerName. Default - false. If parameter is set to true, specify of `header_name`, `header_value` are required.
* `header_name` - (Optional) The name of the header to send header content in. By default, the header name is 'clientCertificateInfo'. Specifying this parameter is relevant only if `forward_to_origin` is set to true. Default - `clientCertificateInfo`.
* `header_value` - (Optional) The content to send in the header specified by headerName. One of the following: FULL_CERT (for full certificate in Base64) COMMON_NAME (for certificate's common name (CN)) FINGERPRINT (for the certificate fingerprints in SHA1) SERIAL_NUMBER (for the certificate's serial number). Parameter is required if `forward_to_origin` is set to true. Default: `FULL_CERT`.
* `is_disable_session_resumption` - (Optional) Disables SSL session resumption for site. Needed when Incapsula Client Certificate is needed only for specific hosts/ports and site have clients that reuse TLS session across different hosts/ports. Default - false.

## Attributes Reference

The following attributes are exported:

* `id` - Value of Site ID in the Incapsula MTLS Client to Imperva Certificate Site Settings.

## Import

Incapsula mTLS Client to Imperva Certificate Site Settings can be imported using Site ID :

```
$ terraform import incapsula_mtls_client_to_imperva_certificate_site_settings.demo_client_to_imperva_certificate_site_settings 1234

```

