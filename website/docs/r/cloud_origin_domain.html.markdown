---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_cloud_origin_domain"
description: |-
  Provides a Cloud Origin Domain resource for Imperva for AWS sites.
---

# incapsula_cloud_origin_domain

Provides a Cloud Origin Domain resource for connecting AWS origins to Imperva for AWS (PUBLIC_CLOUD) sites.

This resource registers your AWS origin (ALB, NLB, or custom domain) with Imperva. Imperva then generates an origin domain that you configure in your AWS CloudFront distribution to route traffic through Imperva's security layer.

## Example Usage

### Basic Usage

```hcl
resource "incapsula_site_v3" "aws_site" {
  name       = "my-aws-site"
  type       = "PUBLIC_CLOUD"
  cloud_type = "AWS"
  ref_id     = "cf-dist-E1234567890"
}

resource "incapsula_cloud_origin_domain" "origin" {
  account_id        = incapsula_site_v3.aws_site.account_id
  site_id           = incapsula_site_v3.aws_site.id
  domain            = "internal-alb-1234567890.us-east-1.elb.amazonaws.com"
  region            = "us-east-1"
  origin_tls_policy = "TLS_1_2"
}
```

### With All Parameters

```hcl
resource "incapsula_cloud_origin_domain" "origin" {
  account_id        = incapsula_site_v3.aws_site.account_id
  site_id           = incapsula_site_v3.aws_site.id
  domain            = "internal-alb-1234567890.us-east-1.elb.amazonaws.com"
  region            = "us-east-1"
  port              = 8443
  origin_tls_policy = "TLS_1_2"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `site_id` - (Required) Numeric identifier of the site. The site type must be set to PUBLIC_CLOUD, with cloud_type = "AWS". Cannot be changed after the resource is created.
* `domain` - (Required) The origin domain (FQDN). Must be a valid fully qualified domain name such as an AWS ALB or NLB hostname (e.g., `internal-alb-1234567890.us-east-1.elb.amazonaws.com`). Cannot be changed after the resource is created. Maximum 253 characters.
* `region` - (Required) The AWS region where the origin is located. Supported values: `us-east-1`, `us-east-2`, `us-west-1`, `us-west-2`, `eu-west-1`, `eu-west-2`, `eu-west-3`, `eu-central-1`, `eu-north-1`, `ap-northeast-1`, `ap-northeast-2`, `ap-southeast-1`, `ap-southeast-2`, `ap-south-1`, `sa-east-1`.
* `port` - (Optional) Port number the origin server listens on. Must be 443 or in the range 1024-65535. Default: `443`.
* `origin_tls_policy` - (Required) Minimum TLS version for the connection to the origin. The selected version will be supported along with all higher versions. Supported values: `SSLv3`, `TLS_1_0`, `TLS_1_1`, `TLS_1_2`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in format `account_id/site_id/origin_id`.
* `imperva_origin_domain` - The Imperva-generated routing domain. Use this value as the origin domain in your AWS CloudFront distribution to route traffic through Imperva.
* `created_at` - Timestamp when the cloud origin domain resource was created.
* `updated_at` - Timestamp when the cloud origin domain resource was last updated.

## Import

Cloud Origin Domain can be imported using the format `account_id/site_id/origin_id`:

```
$ terraform import incapsula_cloud_origin_domain.example 55865773/123456/789
```

## Usage Notes

1. The site must be created with `type = "PUBLIC_CLOUD"` and `cloud_type = "AWS"` before adding cloud origin domains.
2. A site can have multiple origin domains spread across different AWS regions.
3. After creation, copy the `imperva_origin_domain` value and configure it as the origin domain in your AWS CloudFront distribution.
4. It is recommended to create one site per CloudFront distribution.