---
layout: "incapsula"
page_title: "Incapsula: incapsula-siem-connection-s3arn"
sidebar_current: "docs-incapsula-siem-connection-s3arn"
description: |-
Provides a customer S3 ARN connection configuration resource.
---

# incapsula_siem_connection_s3arn

Provides a customer S3 ARN connection configuration resource.
This resource is used to manage connection to customer AWS S3 bucket that used ARN role for access.
The connection doesn't contain connection credentials, but instead you should grant putObject permission of Imperva Logs AWS role :
[Link to documentation](https://docs.imperva.com/bundle/cloud-application-security/page/siem-log-configuration.htm)


## Example Usage

```hcl
resource "incapsula_siem_connection_s3arn" "example_siem_connection_arn"{
	accountId = "1234567"
	connectionName = "CWAF SIEM-LOGS CONNECTION"
  	storageType = "CUSTOMER_S3_ARN"
  	path = "myBucket/siem/logs"
}
```

## Argument Reference

The following arguments are supported:

* `connectionName` - (Required) Unique connection name.
* `storageType` - (Required) Storage type CUSTOMER_S3_ARN.
* `path` - (Required) Path to the files inside bucket including bucket name: `bucketName/folder/subfolder`.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the Customer S3 connection.

## Import

Customer connection can be imported using `connectionId`:

```
$ terraform import incapsula_siem_connection_s3arn.example_siem_connection_arn connectionId
```