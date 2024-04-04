---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_custom_certificate"
description: |-
  Provides a Incapsula Custom Certificate resource.
---

# incapsula_custom_certificate

Provides a Incapsula Custom Certificate resource. 
Custom certificates must be one of the following formats: PFX, PEM, or CER.

## Example Usage

```hcl
resource "incapsula_custom_certificate" "custom-certificate" {
    site_id = incapsula_site.example-site.id
    certificate = filebase64("${"path/to/your/cert.crt"}")
    private_key = filebase64("${"path/to/your/private_key.key"}")
    auth_type   = "RSA/ECC"
    passphrase = "yourpassphrase"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `certificate` - (Required) The certificate file in base64 format. You can use the Terraform HCL `file` directive to pull in the contents from a file. You can also inline the certificate in the configuration.
* `private_key` - (Optional) The private key of the certificate in base64 format. Optional in case of PFX certificate file format.
* `passphrase` - (Optional) The passphrase used to protect your SSL certificate.
* `auth_type` - (Optional) The authentication type of the certificate (RSA/ECC). If not provided then RSA will be taken as a default.
* `input_hash` - (Optional) Currently ignored. If terraform plan flags this field as changed, it means that any of: `certificate`, `private_key`, or `passphrase` has changed.

## Attributes Reference

The following attributes are exported:

* `id` - At the moment, only one active certificate can be stored. This exported value is always set as `12345`. This will be augmented in future versions of the API.

## Import

Custom Certificate cannot be imported.