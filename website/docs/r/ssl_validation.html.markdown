---
subcategory: "Cloud WAF - Certificate Management"
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
<br/>
**Note: This resource is designed to work with sites represented by the "incapsula_site_v3" resource only.**

## Full Imperva Site setup Example Usage

```hcl

# V3 Site
resource "incapsula_site_v3" "example-v3-site" {
  name = "example site"
}

# Manage certificate
resource "incapsula_managed_certificate_settings" "example-site-cert" {
  site_id = incapsula_site_v3.example-v3-site.id
  default_validation_method = "CNAME"
}

# Domains
resource "incapsula_domain" "domain1" {
  site_id = incapsula_site_v3.example-v3-site.id
  domain="bb.terraform-demo-113311111111.incaptest.co"
}

resource "incapsula_domain" "domain2" {
  site_id = incapsula_site_v3.example-v3-site.id
  domain="bb.terraform-demo-1133111dfg11111.incaptest.co"
}

locals {
  domains = toset([incapsula_domain.domain1, incapsula_domain.domain2])
  domain_ids = toset([incapsula_domain.domain1.id, incapsula_domain.domain2.id])
}

# SSL instructions
data "incapsula_ssl_instructions" "example-site-instructions" {
  site_id = incapsula_site_v3.example-v3-site.id
  managed_certificate_settings_id = incapsula_managed_certificate_settings.example-site-cert.id
  domain_ids = local.domain_ids
}

# Add the SSL validation records on your DNS provider
# Use the response from data.incapsula_ssl_instructions.example-site-instructions.instructions and review the instructions for each for each domain.
# Note: In some cases the incapsula_ssl_instructions data source does not return instructions for all the domains.
# For more details see the documentation of incapsula_ssl_instructions.

# Block until the certificate is ready
resource "incapsula_ssl_validation" "example-ssl-validation" {
  site_id = incapsula_site_v3.example-v3-site.id
  domain_ids = local.domain_ids

  depends_on = [
    # Your DNS provider resource that creates the records
  ]
}

# Point the traffic to Imperva after the managed certificate is ready
# Note: the resource depends on incapsula_ssl_validation.example-ssl-validation
  

# Data centers configuration
resource "incapsula_data_centers_configuration" "example-data-centers-configuration" {
  site_id = incapsula_site_v3.example-v3-site.id
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


