---
subcategory: "SIEM"
layout: "incapsula"
page_title: "incapsula_siem_sftp_connection"
description: |- 
    Provides a customer Sftp connection configuration resource.
---

# incapsula_siem_sftp_connection

Provides a customer SFTP connection configuration resource.

The follow-on action is to use the `incapsula_siem_log_configuration` resource to configure the connection.

## Example Usage

```hcl
resource "incapsula_siem_sftp_connection" "example_siem_sftp_connection"{
	account_id = "1234567"
	connection_name = "ABP siem-logs SFTP connection"
  	host = "ec2.eu-west-2.compute.amazonaws.com"
  	path = "/example/accounts/1234567"
  	username = "example_sftp_user"
  	password = "Sensitive data placeholder"
}
```
> **NOTE:**
For security reasons, when a resource is exported, the `password` field will be replaced with the value `Sensitive data placeholder`.
The actual values are still used in the communication with the SFTP server.
Note - This resource cannot be updated unless you specify a real value for the `password` field instead of `Sensitive data placeholder`.
To clarify, none of the fields in exported resources can be updated unless a real `password` value is set.

Example of exported resource:

```hcl
resource "incapsula_siem_sftp_connection" "example_siem_sftp_connection"{
	account_id = "1234567"
	connection_name = "APB siem-logs SFTP connection"
  	host = "ec2.eu-west-2.compute.amazonaws.com"
  	path = "/example/accounts/1234567"
  	username = "example_sftp_user"
  	password = "Sensitive data placeholder"
}
```
## Argument Reference

The following arguments are supported:

* `connection_name` - (Required) Unique connection name.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `host` - (Required) SFTP server host.
* `path` - (Required) SFTP server path.
* `username` - (Required) SFTP access username.
* `password` - (Required) SFTP access password. 

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the customer’s SFTP connection.

## Import

Customer connection can be imported using `accountId`/`connectionId`:

```
$ terraform import incapsula_siem_sftp_connection.example_siem_sftp_connection accountId/connectionId
```
