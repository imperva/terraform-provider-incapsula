---
layout: "incapsula"
page_title: "Incapsula: incapsula-siem-log-configuration"
sidebar_current: "docs-incapsula-siem-log-configuration"
description: |-
Provides a Log Configuration resource.
---

# incapsula_siem_log_configuration

Provides a log configuration resource.
This resource is used to manage log configurations which describe the desired log types
and the destination for the log files. The destination is defined by the <connection ID>.
[Learn more](https://docs.imperva.com/bundle/cloud-application-security/page/siem-log-configuration.htm)


## Example Usage

```hcl
resource "incapsula_siem_connection" "example_siem_connection_basic_auth"{
	accountId = "1234567"
	connectionName = "CWAF SIEM-LOGS CONNECTION"
  	storageType = "CUSTOMER_S3_ARN"
  	path = "myBucket/siem/logs"
}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_abp"{
    accountId = 1234567
  	configurationName = "ABP SIEM-LOGS configuration"
  	producer = "ABP"
	datasets = ["ABP"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_netsec"{
    accountId = 1234567
  	configurationName = "NETSEC SIEM-LOGS configuration"
  	producer = "NETSEC"
	datasets = ["CONNECTION", "IP"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_ato"{
    accountId = 1234567
  	configurationName = "ATO SIEM-LOGS configuration"
  	producer = "ATO"
	datasets = ["ATO"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_audit"{
    accountId = 1234567
  	configurationName = "AUDIT Trail SIEM-LOGS configuration"
  	producer = "AUDIT"
	datasets = ["AUDIT_TRAIL"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_csp"{
    accountId = 1234567
  	configurationName = "CSP SIEM-LOGS configuration"
  	producer = "CSP"
	datasets = ["GOOGLE_ANALYTICS_IDS", "SIGNIFICANT_DOMAIN_DISCOVERY"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}
```

## Argument Reference

The following arguments are supported:
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `configurationName` - (Required) Unique configuration name.
* `producer` - (Required) Provider type. Values: `ABP`, `NETSEC`, `ATO`, `AUDIT`
* `datasets` - (Required) An array of strings representing the type of logs. Values:<br /> `ABP` for provider type `ABP`<br /> `CONNECTION`, `NETFLOW`, `IP`, `ATTACK` for provider type `NETSEC`<br /> `ATO` for provider type `ATO`<br /> `AUDIT_TRAIL` for provider type `AUDIT` <br/> `GOOGLE_ANALYTICS_IDS`, `SIGNIFICANT_DOMAIN_DISCOVERY` for provider type `CSP`
* `enabled`  - (Required) Boolean. Values: `true`/ `false`
* `connectionId` - (Required) Connection id associated with this log configuration

**Note**: The connection should be chosen according to conjunction of producer and dataset:

| producer | allowed storage_type | datasets                        |
|----------|----------------------|---------------------------------|
| ABP      | CUSTOMER_S3          | ABP                             |
| NETSEC   | CUSTOMER_S3          | CONNECTION, NETFLOW, IP, ATTACK |
| ATO      | CUSTOMER_S3          | ATO                             |
| AUDIT    | CUSTOMER_S3          | AUDIT_TRAIL                     |
| CSP      | CUSTOMER_S3          | GOOGLE_ANALYTICS_IDS, SIGNIFICANT_DOMAIN_DISCOVERY          |


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the log configuration.

## Import

Customer connection can be imported using `accountId`/`logConfigurationId`:

```
$ terraform import incapsula_siem_log_configuration.example_siem_log_configuration_abp accountId/logConfigurationId
```
