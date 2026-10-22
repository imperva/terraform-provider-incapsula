---
subcategory: "SIEM"
layout: "incapsula"
page_title: "incapsula_siem_generic_http_connection"
description: |-
    Provides a customer Generic HTTP connection configuration resource.
---

# incapsula_siem_generic_http_connection

Provides a customer Generic HTTP connection configuration resource. Logs are delivered to an HTTPS endpoint, authenticated either with secret headers (`HEADERS`) or with basic authentication (`BASIC_AUTH`).

The follow-on action is to use the `incapsula_siem_log_configuration` resource to configure the connection.

## Example usage

### Headers authentication

```hcl
resource "incapsula_siem_generic_http_connection" "example_headers" {
  account_id      = "1234567"
  connection_name = "Google SecOps"
  url             = "https://malachiteingestion-pa.googleapis.com/v2/unstructuredlogentries:batchCreate"
  auth_type       = "HEADERS"

  auth_header {
    header_name   = "X-goog-api-key"
    secret_format = "{secret}"
    secret        = var.google_api_key
  }

  auth_header {
    header_name   = "X-Webhook-Access-Key"
    secret_format = "{secret}"
    secret        = var.webhook_access_key
  }

  headers = {
    "Content-Type" = "application/json"
  }

  disable_cert_verification = false
}
```

### Basic authentication

```hcl
resource "incapsula_siem_generic_http_connection" "example_basic_auth" {
  account_id      = "1234567"
  connection_name = "My HTTP receiver"
  url             = "https://logs.example.com/ingest"
  auth_type       = "BASIC_AUTH"
  username        = "admin"
  password        = var.receiver_password
}
```

### Using the connection in a log configuration

Referencing the connection's `id` lets Terraform create the connection first and destroy the log configuration first.
Log formats `CEF`, `LEEF`, and `W3C` cannot be sent with `Content-Type: application/json`, so this connection sends `text/plain`.

```hcl
resource "incapsula_siem_generic_http_connection" "example_webhook" {
  connection_name = "My webhook"
  url             = "https://logs.example.com/ingest"
  auth_type       = "HEADERS"

  auth_header {
    header_name   = "Authorization"
    secret_format = "Bearer {secret}"
    secret        = var.webhook_token
  }

  headers = {
    "Content-Type" = "text/plain"
  }
}

resource "incapsula_siem_log_configuration" "example" {
  configuration_name   = "Cloud WAF to HTTP"
  producer             = "CLOUD_WAF"
  datasets             = ["WAF_RAW_LOGS"]
  enabled              = true
  connection_id        = incapsula_siem_generic_http_connection.example_webhook.id
  format               = "CEF"
  batch_packaging_type = "NONE"
}
```

> **NOTE:**
Create and update perform a live connectivity test against `url`. If the test fails, `terraform apply` fails with code `1000` and a `detail` explaining why, for example:
```
status code: 400
code: 1000
detail: The destination couldn't be reached. Check the host and port, and that the destination accepts connections from Imperva's IP ranges.
```

> **NOTE:**
A connection cannot be deleted while a log configuration references it. Reference the connection's `id` from the log configuration (as above), or use `depends_on`, so Terraform destroys the log configuration first.

> **NOTE:**
`secret` and `password` are write-only: they are never read back from the API. Changing them outside Terraform is not detected as drift. After import, set them in the configuration and run `terraform apply`.

## Argument reference

The following arguments are supported:

* `connection_name` - (Required) Unique connection name.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `url` - (Required) HTTPS endpoint the logs are sent to.
* `auth_type` - (Required) Authentication type. Values: `HEADERS`, `BASIC_AUTH`.
* `auth_header` - (Optional) Secret authentication header, 1 to 5 blocks, sent in declaration order. Required when `auth_type` is `HEADERS`; not allowed with `BASIC_AUTH`. Each block supports:
    * `header_name` - (Required) Header name. Must be unique (case-insensitive), including against `headers`.
    * `secret_format` - (Required) Header value template containing the `{secret}` placeholder, for example `Bearer {secret}`.
    * `secret` - (Required, Sensitive) Secret substituted into `secret_format`.
* `username` - (Optional) Username. Required when `auth_type` is `BASIC_AUTH`; not allowed with `HEADERS`.
* `password` - (Optional, Sensitive) Password. Required when `auth_type` is `BASIC_AUTH`; not allowed with `HEADERS`.
* `headers` - (Optional) Map of additional non-secret HTTP headers, for example `Content-Type`. Up to 20 entries.
* `disable_cert_verification` - (Optional) Disable server certificate verification, for endpoints using a self-signed certificate. Default: `false`.

## Attributes reference

The following attributes are exported:

* `id` - Unique identifier of the customer's Generic HTTP connection.

## Import

You can import a customer connection using `accountId`/`connectionId`:

```
$ terraform import incapsula_siem_generic_http_connection.example_headers accountId/connectionId
```
