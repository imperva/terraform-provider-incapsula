---
layout: "incapsula"
page_title: "Incapsula: incapsula-siem-connection"
sidebar_current: "docs-incapsula-siem-connection"
description: |-
Provides a customer S3 connection configuration resource.
---

# incapsula_siem_connection_s3

Provides a customer S3 connection configuration resource.
This resource is used to manage the connection to the customer’s AWS S3 bucket.
[Learn more](https://docs.imperva.com/bundle/cloud-application-security/page/siem-log-configuration.htm)

## Example Usage

```hcl
resource "incapsula_siem_connection" "example_siem_connection_s3_basic"{
	account_id = "1234567"
	connection_name = "APB siem-logs connection basic auth"
	storage_type = "CUSTOMER_S3"
  	access_key = "AKIAIOSFODNN7EXAMPLE"
  	secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  	path = "myBucket/siem/logs"
}

resource "incapsula_siem_connection" "example_siem_connection_s3_arn"{
	account_id = "1234567"
	connection_name = "ABP siem-logs connection arn auth"
	storage_type = "CUSTOMER_S3_ARN"
  	path = "myBucket/siem/logs"
}
```

## Argument Reference

The following arguments are supported:

* `connection_name` - (Required) Unique connection name.
* `path` - (Required) Path to the files inside bucket including bucket name: `bucketName/folder/subfolder`.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `access_key` - (Required when storage_type="CUSTOMER_S3" ) AWS Access key.
* `secret_key` - (Required when storage_type="CUSTOMER_S3") AWS access secret.
* `storage_type` - (Required) Storage type. Possible values: `CUSTOMER_S3`, `CUSTOMER_S3_ARN` 

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the customer’s S3 connection.

## Import

Customer connection  can be imported using `connectionId`:

```
$ terraform import incapsula_siem_connection_s3.example_siem_connection connectionId
```