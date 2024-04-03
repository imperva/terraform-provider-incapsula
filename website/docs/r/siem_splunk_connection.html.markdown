---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_siem_splunk_connection"
description: |- 
    Provides a customer Splunk connection configuration resource.
---

# incapsula_siem_splunk_connection

Provides a customer Splunk connection configuration resource.

The follow-on action is to use the `incapsula_siem_log_configuration` resource to configure the connection.

## Example Usage

```hcl
resource "incapsula_siem_splunk_connection" "example_siem_splunk_connection"{
	account_id = "1234567"
	connection_name = "APB siem-logs Splunk connection"
  	host = "my.splunk.com"
  	port = 8080
  	token = "Sensitive data placeholder"
  	disable_cert_verification = false
}
```
> **NOTE:**
For security reasons, when a resource is exported, the `token` field will be replaced with the value `Sensitive data placeholder`.
The actual values are still used in the communication with the Splunk server.
Note - This resource cannot be updated unless you specify a real value for the `token` field instead of `Sensitive data placeholder`.
To clarify, none of the fields in exported resources can be updated unless a real `token` value is set.

Example of exported resource:

```hcl
resource "incapsula_siem_splunk_connection" "example_siem_splunk_connection"{
	account_id = "1234567"
	connection_name = "APB siem-logs Splunk connection"
  	host = "my.splunk.com"
  	port = 8080
  	token = "Sensitive data placeholder"
  	disable_cert_verification = false
}
```
## Argument Reference

The following arguments are supported:

* `connection_name` - (Required) Unique connection name.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `host` - (Required) Splunk server host.
* `port` - (Required) Splunk server port.
* `token` - (Required) Splunk access token - Version 4 UUID format. 
* `disable_cert_verification` - (Optional) Flag to disable/enable server certificate. Used when a self-signed certificate applied on the server side.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the customer’s Splunk connection.

## Import

Customer connection  can be imported using `accountId`/`connectionId`:

```
$ terraform import incapsula_siem_splunk_connection.example_siem_splunk_connection accountId/connectionId
```
