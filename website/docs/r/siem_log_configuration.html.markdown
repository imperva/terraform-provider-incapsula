---
subcategory: "SIEM"
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

```

## Argument Reference

The following arguments are supported:
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `configurationName` - (Required) Unique configuration name.
* `producer` - (Required) Provider type, see the table below. The provider does not check it; the API validates it, and an unknown producer fails `terraform apply` with an error naming the `provider` field.
* `datasets` - (Required) An array of strings representing the type of logs. The datasets must match the producer, see the table below. The provider does not check them; the API validates them and `terraform apply` fails with the API error: an unknown dataset names the field (e.g. `datasets[0]`), and a dataset that does not match the producer names both.
* `enabled`  - (Required) Boolean. Values: `true`/ `false`
* `connectionId` - (Required) Connection id associated with this log configuration
* `logs_level` - (Optional) Security log level - compatible only with CLOUD_WAF producer. Values: `NONE`, `FULL`, `SECURITY`
* `compress_logs` - (Optional) Boolean - compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers. Values: `true`/ `false`
* `format` - (Optional) Log format - compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers. Values: `CEF`, `W3C` , `LEEF`
* `public_key` - (Optional) Public key for encryption - compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers.
* `public_key_file_name` - (Optional) The name of the public key file corresponding to the public_key field. This is compatible only with CLOUD_WAF and ATTACK_ANALYTICS producers.
* `batch_packaging_type` - (Optional) How log batches are packaged in the payload sent to the connection. Values: `NONE`, `IMPERVA_METADATA_PREFIX`. When unset, no payload settings are sent.

**Note**: Producers and datasets at the time of writing. New producers and datasets can be added to the API without a provider release.

| producer         | datasets                                                                                                                                                              |
|------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| ABP              | ABP                                                                                                                                                                   |
| NETSEC           | CONNECTION, NETFLOW, IP, ATTACK, NOTIFICATIONS                                                                                                                        |
| ATO              | ATO                                                                                                                                                                   |
| AUDIT            | AUDIT_TRAIL                                                                                                                                                           |
| CSP              | GOOGLE_ANALYTICS_IDS, SIGNIFICANT_DOMAIN_DISCOVERY, SIGNIFICANT_SCRIPT_DISCOVERY, SIGNIFICANT_DATA_TRANSFER_DISCOVERY,DOMAIN_DISCOVERY_ENFORCE_MODE,CSP_HEADER_HEALTH |
| CLOUD_WAF        | WAF_RAW_LOGS, CLOUD_WAF_ACCESS                                                                                                                                        |
| ATTACK_ANALYTICS | WAF_ANALYTICS_LOGS                                                                                                                                                    |


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier of the log configuration.

## Import

Customer connection can be imported using `accountId`/`logConfigurationId`:

```
$ terraform import incapsula_siem_log_configuration.example_siem_log_configuration_abp accountId/logConfigurationId
```
