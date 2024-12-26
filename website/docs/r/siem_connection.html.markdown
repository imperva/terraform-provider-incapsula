---
subcategory: "SIEM"
layout: "incapsula"
page_title: "incapsula_siem_connection"
description: |- 
    Provides a customer S3 connection configuration resource.
---

# incapsula_siem_connection

Provides a customer S3 connection configuration resource.

The follow-on action is to use the `incapsula_siem_log_configuration` resource, to configure the connection.

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
> **NOTE:**
For security reasons, when a resource is exported, the `access_key` and `secret_key` fields will be replaced with the value `Sensitive data placeholder`.
The actual key values are still used in the communication with the s3 bucket.
Note - This resource cannot be updated unless you specify a real value for the `access_key` and `secret_key` fields instead of `Sensitive data placeholder`.
To clarify, none of the fields in exported resources can be updated unless real `access_key` and `secret_key` values are set.

Example of exported resource:

```hcl
resource "incapsula_siem_connection" "example_siem_connection_s3_basic"{
	account_id = "1234567"
	connection_name = "APB siem-logs connection basic auth"
	storage_type = "CUSTOMER_S3"
  	access_key = "Sensitive data placeholder"
  	secret_key = "Sensitive data placeholder"
  	path = "myBucket/siem/logs"
}
```
## Argument Reference

The following arguments are supported:

* `connection_name` - (Required) Unique connection name.
* `path` - (Required) Path to the files inside bucket including bucket name: `bucketName/folder/subfolder`.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `access_key` - (Required when storage_type="CUSTOMER_S3" ) AWS access key.
* `secret_key` - (Required when storage_type="CUSTOMER_S3") AWS access secret.
* `storage_type` - (Required) Storage type. Possible values: `CUSTOMER_S3`, `CUSTOMER_S3_ARN` 

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the customer’s S3 connection.

## Import

Customer connection  can be imported using `accountId`/`connectionId`:

```
$ terraform import incapsula_siem_connection.example_siem_connection accountId/connectionId
```
