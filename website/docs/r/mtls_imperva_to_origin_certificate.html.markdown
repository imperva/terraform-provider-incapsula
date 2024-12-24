---
subcategory: "Cloud WAF - Certificate Management"
layout: "incapsula"
page_title: "incapsula_mtls_imperva_to_origin_certificate"
description: |-
  Provides a Mutual TLS Imperva to Origin certificate resource.
---

# incapsula_mtls_imperva_to_origin_certificate

Provides a Mutual TLS Imperva to Origin certificate resource.
This resource is used to upload mTLS client certificates to enable mutual authentication between Imperva and origin servers.
Mutual TLS Imperva to Origin Certificates must be in one of the following formats: pem, der, pfx, cert, crt, p7b, cer, p12, key, ca-bundle, bundle, priv, cert.

## Example Usage

```hcl
resource "incapsula_mtls_imperva_to_origin_certificate" "mtls_imperva_to_origin_certificate"{
  certificate       = filebase64("${"path/to/your/cert.pem"}")
  private_key       = filebase64("${"path/to/your/private_key.pem"}")
  passphrase        = "my_passphrase"
  certificate_name  = "pem certificate example"
  account_id        = "incapsula_account.example-account.id"  
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required) Your mTLS client certificate file. Supported formats: pem, der, pfx, cert, crt, p7b, cer, p12, ca-bundle, bundle, cert.
  You can use the Terraform HCL `filebase64` directive to pull in the contents from a file. You can also embed the certificate in the configuration.
* `private_key` - Your private key file. supported formats: pem, der, priv, key. If pfx or p12 certificate is used, then this field can remain empty.
* `passphrase` - Your private key passphrase. Leave empty if the private key is not password protected.
* `certificate_name` - (Optional) A descriptive name for your mTLS Certificate.
* `account_id` - (Required) Numeric identifier of the account to operate on.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Mutual TLS Imperva to Origin Certificate.