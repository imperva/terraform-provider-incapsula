---
layout: "incapsula"
page_title: "Incapsula: custom-certifcate"
sidebar_current: "docs-incapsula-resource-custom-certificate"
description: |-
  Provides a Incapsula Custom Certificate resource.
---

# incapsula_custom_certificate

Provides a Incapsula Custom Certificate resource. 
Custom certificates must be one of the following formats: PFX, PEM, or CER.

## Example Usage

```hcl
resource "incapsula_custom_certificate" "custom-certificate" {
    site_id = "${incapsula_site.example-site.id}"
    certificate = "${file("path/to/your/cert.crt")}"
    private_key = "${file("path/to/your/private_key.key")}"
    passphrase = "yourpassphrase"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `certificate` - (Required) The certificate file in base64 format. You can use the Terraform HCL `file` directive to pull in the contents from a file. You can also inline the certificate in the configuration.
* `private_key` - (Optional) The private key of the certificate in base64 format. Optional in case of PFX certificate file format.
* `passphrase` - (Optional) The passphrase used to protect your SSL certificate.

## Attributes Reference

The following attributes are exported:

* `id` - At the moment, only one active certificate can be stored. This exported value is always set as `12345`. This will be augmented in future versions of the API.
