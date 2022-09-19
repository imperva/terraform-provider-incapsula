---
layout: "incapsula"
page_title: "Incapsula: incapsula-mtls-client-to-imperva-ca-certificate"
sidebar_current: "docs-incapsula-mtls-client-to-imperva-ca-certificate"
description: |-
Provides a Mutual TLS Client to Imperva CA Certificate resource.
---

# incapsula_mtls_client_to_imperva_ca_certificate

Provides a Mutual TLS Client to Imperva CA Certificate resource.
This resource is used to upload the client CA certificate used by Imperva to validate the client certificate
Mutual TLS Client to Imperva Certificates must be one of the following formats: PEM, CRT, CER or CA.
Update action is not supported for current resource. Please create a new Mutual TLS Client to Imperva CA Certificate resource and only then, delete the old one.

## Example Usage
Reference to account data source in `account_id` field

```hcl
data "incapsula_account_data" "account_data" {
}

resource "incapsula_mtls_client_to_imperva_ca_certificate" "client_ca_certificate_1"{
  certificate_name = "pem certificate example"
  certificate      = filebase64("./ca_certificate.pem")
  account_id       =  data.incapsula_account_data.account_data.current_account
}
```

Reference to subaccount resource in `account_id` field

```hcl
resource "incapsula_subaccount" "example-subaccount" {
  sub_account_name  = "Example SubAccount"
  logs_account_id   = "789"
  log_level         = "full"
}

resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_1"{
  certificate_name = "pem certificate example"
  certificate      = filebase64("./ca_certificate.pem")
  account_id       = incapsula_subaccount.example-subaccount.id
}
```

`Exported Certificate - data placeholder` value will be set in the `certificate` field in the result of account export.
The reason for using these values is that this certificate currently exists in the account configuration and this resource allows using it with new sites configured via Terraform.
Note - This resource cannot be updated unless you specify a real values to the `certificate` field instead of `Exported Certificate - data placeholder`.
In exported resources, the certificate_name cannot be changed as well until real values are set.

Example of exported resource:

```hcl
resource "incapsula_mtls_client_to_imperva_ca_certificate" "client_ca_certificate_1"{
  certificate_name = "pem certificate example"
  certificate      = "Exported Certificate - data placeholder"
  account_id           = incapsula_subaccount.example-subaccount.id
}
```

## Argument Reference

The following arguments are supported:

* `certificate` - (Required) Your mTLS client certificate file. Supported formats: PEM, CRT, CER and CA.
  You can use the Terraform HCL `filebase64` directive to pull in the contents from a file. You can also inline the certificate in the configuration.
* `account_id` - (Required) Numeric identifier of the account to operate on.
* `certificate_name` - (Optional) A descriptive name for your mTLS Client Certificate.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Mutual TLS Imperva to Origin Certificate.

## Import

Incapsula Mutual TLS Imperva to Origin Certificate can be imported using `account_id` and `certificate_id`:

```
$ terraform import incapsula_mtls_client_to_imperva_ca_certificate.client_ca_certificate_1 account_id/certificate_id
```