---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_ssl_validation"
description: |- 
  Provides an Incapsula SSL validation resource.
---

# incapsula_ssl_validation

Provides an Incapsula SSL validation resource.
This resource is dependent on the incapsula_managed_certificate_settings, the incapsula_ssl_instructions data source, and a DNS terraform provider.
The provider will be blocked until SSL coverage for all domains is configured on the Imperva managed certificate.
This resource can be used only when the SSL validation method for the domains is DNS/CNAME and the DNS records are managed by a DNS terraform provider.
Customers that use this resource will be able to manage fully configured Imperva sites using Terraform.

## Full Imperva Site setup Example Usage

```hcl

# V3 Site
resource "incapsula_site_v3" "example-v3-site" {
  name = "example site"
}

# Manage certificate
resource "incapsula_managed_certificate_settings" "example-site-cert" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  default_validation_method = "CNAME"
}

# Domains
resource "incapsula_domain" "domain1" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  domain="bb.terraform-demo-113311111111.incaptest.co"
}

resource "incapsula_domain" "domain2" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  domain="bb.terraform-demo-1133111dfg11111.incaptest.co"
}

locals {
  domains = toset([incapsula_domain.domain1, incapsula_domain.domain2])
  domain_ids = toset([incapsula_domain.domain1.id, incapsula_domain.domain2.id])
}

# SSL instructions
data "incapsula_ssl_instructions" "example-site-instructions" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  managed_certificate_settings_id = incapsula_managed_certificate_settings.example-site-cert.id
  domain_ids = local.domain_ids
}

# Add the SSL validation records on AWS Route53
resource "aws_route53_record" "ssl-records" {
  for_each = {
    for dom in local.domains : dom.domain =>
    [for ins in data.incapsula_ssl_instructions.example-site-instructions.instructions :  ins if ins.domain_id == tonumber(dom.id)]
  }

  zone_id = "AAAA"
  name    = each.value[0].name
  type    = each.value[0].type
  ttl     = 300
  records = [each.value[0].value]

}

# Block until the certificate is ready
resource "incapsula_ssl_validation" "example-ssl-validation" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  domain_ids = local.domain_ids

  depends_on = [
    aws_route53_record.ssl-records
  ]
}

# Point the traffic to Imperva after the managed certificate is ready
resource "aws_route53_record" "network-records" {

  depends_on = [
    incapsula_ssl_validation.example-ssl-validation
  ]
  for_each = {
    for dom in local.domains : dom.domain => dom
  }

  zone_id = "AAAA"
  name    = each.value.domain
  type    = length(each.value.a_records) > 0 ? "A" : "CNAME"
  ttl     = 300
  records = length(each.value.a_records) > 0 ? each.value.a_records : [incapsula_domain.domain1.cname_redirection_record]

}

# Data centers configuration
resource "incapsula_data_centers_configuration" "example-data-centers-configuration" {
  site_id = incapsula_site_v3.example-v3-site.site_id
  site_topology = "SINGLE_DC"

  data_center {
    name = "DC1"
    ip_mode = "SINGLE_IP"

    origin_server {
      address = "1.2.34"
      is_active = true
    }
  }

}



```


## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `domain_ids` - (Required) List of incapsula_domain ids that .

## Attributes Reference

The following attributes are exported:

* `id` - The id of the SSL validation resource.


## Import/Destroy

SSL validation resource cannot be imported or destroyed


