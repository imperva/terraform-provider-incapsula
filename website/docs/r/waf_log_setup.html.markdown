---
subcategory: "Deprecated"
layout: "incapsula"
page_title: "incapsula_waf_log_setup"
description: |-
    Provides an Incapsula WAF Log Setup resource.
---
-> DEPRECATED: incapsula_data_center

This resource has been DEPRECATED. It will be removed in a future version.
Please use the current `incapsula_siem_log_configuration` for CWAF log configuration resource instead.
For SFTP Connection please use the current `incapsula_siem_sftp_connection` resource, and for S3 Connection please use the `incapsula_siem_s3_connection` resource.

# incapsula_waf_log_setup

Provides an Incapsula WAF Log Setup resource.

## Example Usage

Example #1: Setup Activated SFTP Connection
```hcl
resource "incapsula_waf_log_setup" "ex-sftp-waf_log_setup" {
    account_id = 102030
    sftp_host = "samplehost"
    sftp_user_name = "sampleuser"
    sftp_password = "**********"
    sftp_destination_folder = "/home/user_name/log_folder"
}
```

Example #2: Setup Activated S3 Connection
```hcl
resource "incapsula_waf_log_setup" "ex-s3-waf_log_setup1" {
    account_id = 102030
    s3_bucket_name = "bucket_name/log_folder"
    s3_access_key = "AKIAIOSFODNN7EXAMPLE"
    s3_secret_key = "****************************************"
}
```

Example #3: Setup Suspended S3 Connection
```hcl
resource "incapsula_waf_log_setup" "ex-s3-waf_log_setup2" {
    account_id = 102040
    enabled = false
    s3_bucket_name = "bucket_name/log_folder"
    s3_access_key = "AKIAIOSFODNN7EXAMPLE"
    s3_secret_key = "****************************************"
}
```

Example #4: Setup Suspended Default (API) Connection
```hcl
resource "incapsula_waf_log_setup" "ex-s3-waf_log_setup3" {
    account_id = 102050
    enabled = false
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Mandatory) The Numeric identifier of the account to operate on.
* `enabled` - (Optional) A boolean flag to enable or disable WAF Logs. Default value is true.
* `sftp_host` - (Optional) The IP address of your SFTP server.
* `sftp_user_name` - (Optional) A username that will be used to log in to the SFTP server.
* `sftp_password` - (Optional, Sensitive) A corresponding password for the user account used to log in to the SFTP server.
* `sftp_destination_folder` - (Optional) The path to the directory on the SFTP server.
* `s3_bucket_name` - (Optional) S3 bucket name.
* `s3_access_key` - (Optional) S3 access key.
* `s3_secret_key` - (Optional, Sensitive) S3 secret key.

Please note, either sftp_* or s3_* arguments are required group. If neither groups specified default (API) will be set up

Import is not supported yet