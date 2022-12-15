---
layout: "incapsula"
page_title: "Incapsula: incapsula-siem-connection-s3"
sidebar_current: "docs-incapsula-siem-connection-s3"
description: |-
Provides a customer S3 connection configuration resource.
---

# incapsula_siem_connection_s3

Provides a customer S3 connection configuration resource.
This resource is used to manage the connection to the customer’s AWS S3 bucket.
The connection contains connection details: AWS AccessKey, AWS AccessSecret 
and path to the destination repository (including bucket name) for storing the log files.
[Learn more](https://docs.imperva.com/bundle/cloud-application-security/page/siem-log-configuration.htm)

## Example Usage

```hcl
resource "incapsula_siem_connection_s3" "example_siem_connection"{
	accountId = "1234567"
	connectionName = "CWAF SIEM-LOGS CONNECTION"
  	accessKey = "AKIAIOSFODNN7EXAMPLE"
  	secretKey = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  	path = "myBucket/siem/logs"
}
```

## Argument Reference

The following arguments are supported:

* `connectionName` - (Required) Unique connection name.
* `accessKey` - (Required) AWS Access key.
* `secretKey` - (Required) AWS access secret.
* `path` - (Required) Path to the files inside bucket including bucket name: `bucketName/folder/subfolder`.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the customer’s S3 connection.

## Import

Customer connection  can be imported using `connectionId`:

```
$ terraform import incapsula_siem_connection_s3.example_siem_connection connectionId
```