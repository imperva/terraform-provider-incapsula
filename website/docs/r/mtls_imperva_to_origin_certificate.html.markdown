layout: "incapsula"
page_title: "Incapsula: mtls-imperva-to-origin-certificate"
sidebar_current: "docs-incapsula-resource-mtls-imperva-to-origin-certificate"
description: |-
Provides an Incapsula Mutual TLS Imperva to Origin Certificate resource.
---

# incapsula_mtls_imperva_to_origin_certificate

Provides an Incapsula Mutual TLS Imperva to Origin Certificate resource.
Upload and manage your mTLS client certificates to enable mutual authentication between Imperva and your origin servers. 
Mutual TLS Imperva to Origin Certificates must be one of the following formats: PFX, DER, or PEM.
Replace an existing mTLS client certificate that is uploaded to your account. The Imperva certificate ID remains the same after replacement.

## Example Usage
Reference to subaccount resource in `account_id` field

```hcl
resource "incapsula_subaccount" "example-subaccount" {
  sub_account_name  = "Example SubAccount"
  logs_account_id   = "789"
  log_level         = "full"
}

resource "incapsula_mtls_imperva_to_origin_certificate" "mtls_certificate"{
  certificate       = filebase64("./cert.der")
  private_key       = filebase64("./key.der")
  passphrase        = "12345"
  certificate_name  = "pem certificate example"
  account_id        = incapsula_subaccount.mtls_certificate.id
}
```

Account ID is not specified. In this case operation will be performed on the account identified by the authentication parameters.

```hcl
resource "incapsula_mtls_imperva_to_origin_certificate" "mtls_certificate"{
  certificate       = filebase64("./cert.der")
  private_key       = filebase64("./key.der")
  passphrase        = "12345"
  certificate_name  = "pem certificate example"
}

## Argument Reference

The following arguments are supported:

* `certificate` - (Required)Your mTLS client certificate file in base64 format. Supported formats: PEM, DER and PFX. Only RSA certificates are currently supported. The certificate RSA key size must be 2048 bit or less. The certificate must be issued by a certificate authority (CA) and cannot be self-signed. 
You can use the Terraform HCL `filebase64` directive to pull in the contents from a file. You can also inline the certificate in the configuration.
* `private_key` - (Optional) Your private key file in base64 format. Supported formats: PEM, DER. If PFX certificate is used, then this field can remain empty.
* `passphrase` - (Optional) Your private key passphrase. Leave empty if the private key is not password protected.
* `certificate_name` - (Optional) A descriptive name for your mTLS client certificate.
* `account_id` - (Optional) Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.

* `input_hash` - (Optional) Currently ignored. If terraform plan flags this field as changed, it means that any of: `certificate`, `private_key`, or `passphrase` has changed.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Mutual TLS Imperva to Origin Certificate.

## Import

Incapsula Mutual TLS Imperva to Origin Certificate can be imported using `certificate_id`:

```
$ terraform import incapsula_mtls_imperva_to_origin_certificate.mtls_certificate certificate_id
```