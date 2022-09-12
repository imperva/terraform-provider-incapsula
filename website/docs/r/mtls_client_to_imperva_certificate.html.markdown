---
layout: "incapsula"
page_title: "Incapsula: site-mtls-certificate-association"
sidebar_current: "docs-incapsula-site-mtls-certificate-association"
description: |-
Provides an Incapsula Site to Mutual TLS Imperva to Origin Certificate Association resource.
---

# incapsula_mtls_client_to_imperva_certificate

Provides a Mutual TLS Client to Imperva CA Certificate resource.
Mutual TLS Imperva to Origin Certificates must be one of the following formats: PEM, CRT, CER or CA.
Update action is not supported for current resource. Please create a new Imperva CA Certificate resource and only then, delete the old one.

## Example Usage
Reference to account data source in `account_id` field

```hcl
data "incapsula_account_data" "account_data" {
}

resource "incapsula_mtls_client_to_imperva_certificate" "client_ca_certificate_1"{
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
$ terraform import incapsula_mtls_client_to_imperva_certificate.client_ca_certificate_1 account_id/certificate_id
```