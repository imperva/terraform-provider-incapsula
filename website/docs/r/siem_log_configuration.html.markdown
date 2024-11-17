---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_siem_log_configuration"
description: |- 
    Provides a Log Configuration resource.
---

# incapsula_siem_log_configuration

Provides a resource to configure a log connection.

Dependency is on existing connection, created using the `incapsula_siem_connection` or `incapsula_siem_splunk_connection` resource.

## Example Usage

```hcl
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
	datasets = ["GOOGLE_ANALYTICS_IDS", "SIGNIFICANT_DOMAIN_DISCOVERY", "SIGNIFICANT_SCRIPT_DISCOVERY", "SIGNIFICANT_DATA_TRANSFER_DISCOVERY"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_csp"{
    accountId = 1234567
  	configurationName = "CSP SIEM-LOGS configuration"
  	producer = "CSP"
	datasets = ["GOOGLE_ANALYTICS_IDS", "SIGNIFICANT_DOMAIN_DISCOVERY", "SIGNIFICANT_SCRIPT_DISCOVERY", "SIGNIFICANT_DATA_TRANSFER_DISCOVERY"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_cloudwaf"{
    accountId = 1234567
  	configurationName = "CLOUD-WAF SIEM-LOGS configuration"
  	producer = "CLOUD_WAF"
	datasets = [WAF_RAW_LOGS", "CLOUD_WAF_ACCESS"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id
  	logs_level = "NONE"
  	compress_logs = false
  	format = "CEF"
  	publickey = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1IMTnnpYzq0EyjSY6py0
9yI83NfyNI7JFDimfTDbJ/1lpyMnutDeSjrtcq+M0zNj1aBkdxkPXMrutQEOE7F7
6TaYsEITDfE3VGYCNCCvKOWiZ3z8W/FiMAb/3U7akS+j8f5Hrn8nftexfotxq30V
W2ng4pt/1xCyqE8FzQ8Y5DsoemI0CMIauTJpP7E4XRfSqHAeqRxBa27yVcadLnp8
zr41yCpGIiOqD2ubWALQSYOX8gp+Rde0zKBOAlIfct7k4UZQnOxxj8ugAN6zGA5T
8Cn4cGBJffzQCN72Dy4pMPmPNpFQuWAmmW1B5mWWptTa4sV/KUzLAJzP+wyPufKD
NwIDAQAB
-----END PUBLIC KEY-----
  EOT
  	publicKeyFileName = "examplePublicKey.pem"
}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_csp"{
    accountId = 1234567
  	configurationName = "ATTACK-ANALYTICS SIEM-LOGS configuration"
  	producer = "ATTACK_ANALYTICS"
	datasets = ["WAF_ANALYTICS_LOGS"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id
   	format = "CEF"
   	compress_logs = true
    publickey = <<-EOT
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1IMTnnpYzq0EyjSY6py0
9yI83NfyNI7JFDimfTDbJ/1lpyMnutDeSjrtcq+M0zNj1aBkdxkPXMrutQEOE7F7
6TaYsEITDfE3VGYCNCCvKOWiZ3z8W/FiMAb/3U7akS+j8f5Hrn8nftexfotxq30V
W2ng4pt/1xCyqE8FzQ8Y5DsoemI0CMIauTJpP7E4XRfSqHAeqRxBa27yVcadLnp8
zr41yCpGIiOqD2ubWALQSYOX8gp+Rde0zKBOAlIfct7k4UZQnOxxj8ugAN6zGA5T
8Cn4cGBJffzQCN72Dy4pMPmPNpFQuWAmmW1B5mWWptTa4sV/KUzLAJzP+wyPufKD
NwIDAQAB
-----END PUBLIC KEY-----
  EOT
  	publicKeyFileName = "examplePublicKey.pem"

}

resource "incapsula_siem_log_configuration" "example_siem_log_configuration_csp"{
    accountId = 1234567
  	configurationName = "DNSMS SIEM-LOGS configuration"
  	producer = "DNSMS"
	datasets = ["DNSMS_SECURITY_LOGS"]
  	enabled = true
  	connectionId = incapsula_siem_connection.example_siem_connection_basic_auth.id
  	
}
```

## Argument Reference

The following arguments are supported:
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `configurationName` - (Required) Unique configuration name.
* `producer` - (Required) Provider type. Values: `ABP`, `NETSEC`, `ATO`, `AUDIT`, `CLOUD_WAF`, `ATTACK_ANALYTICS`, `DNSMS`
* `datasets` - (Required) An array of strings representing the type of logs. Values:<br /> `ABP` for provider type `ABP`<br /> `CONNECTION`, `NETFLOW`, `IP`, `ATTACK`,`NOTIFICATIONS` for provider type `NETSEC`<br /> `ATO` for provider type `ATO`<br /> `AUDIT_TRAIL` for provider type `AUDIT` <br /> `GOOGLE_ANALYTICS_IDS`, `SIGNIFICANT_DOMAIN_DISCOVERY`, `SIGNIFICANT_SCRIPT_DISCOVERY`, `SIGNIFICANT_DATA_TRANSFER_DISCOVERY`, `DOMAIN_DISCOVERY_ENFORCE_MODE`, `CSP_HEADER_HEALTH` for provider type `CSP`<br /> `WAF_RAW_LOGS`, `CLOUD_WAF_ACCESS` for provider type `CLOUD_WAF` <br /> `WAF_ANALYTICS_LOGS` for provider type `ATTACK_ANALYTICS`<br /> `DNSMS_SECURITY_LOGS` for provider type `DNSMS`
* `enabled`  - (Required) Boolean. Values: `true`/ `false`
* `connectionId` - (Required) Connection id associated with this log configuration
* `logs_level` - (Optional) Security log level - compatible only with CLOUD_WAF producer. Values: `NONE`, `FULL`, `SECURITY`
* `compress_logs` - (Optional) Boolean - compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers. Values: `true`/ `false`
* `format` - (Optional) Log format - compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers. Values: `CEF`, `W3C` , `LEEF`
* `public_key` - (Optional) Public key for encryption - compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers.
* `public_key_file_name` - (Optional) The name of the public key file corresponding to the public_key field. This is compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers.

**Note**: The connection should be chosen according to conjunction of producer and dataset:

| producer         | datasets                                                                                                                                                              |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ABP              | ABP                                                                                                                                                                   |
| NETSEC           | CONNECTION, NETFLOW, IP, ATTACK, NOTIFICATIONS                                                                                                                        |
| ATO              | ATO                                                                                                                                                                   |
| AUDIT            | AUDIT_TRAIL                                                                                                                                                           |
| CSP              | GOOGLE_ANALYTICS_IDS, SIGNIFICANT_DOMAIN_DISCOVERY, SIGNIFICANT_SCRIPT_DISCOVERY, SIGNIFICANT_DATA_TRANSFER_DISCOVERY,DOMAIN_DISCOVERY_ENFORCE_MODE,CSP_HEADER_HEALTH |
| CLOUD_WAF        | WAF_RAW_LOGS, CLOUD_WAF_ACCESS                                                                                                                                        |
| ATTACK_ANALYTICS | WAF_ANALYTICS_LOGS                                                                                                                                                    |
| DNSMS            | DNSMS_SECURITY_LOGS                                                                                                                                                   |


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the log configuration.

## Import

Customer connection can be imported using `accountId`/`logConfigurationId`:

```
$ terraform import incapsula_siem_log_configuration.example_siem_log_configuration_abp accountId/logConfigurationId
```
