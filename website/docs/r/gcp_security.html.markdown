---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "Security Policies for Imperva for Google Cloud Sites"
description: |-
  Describes the Terraform security resources supported for Imperva for Google Cloud sites.
---

# Security Policies for Imperva for Google Cloud Sites

Imperva for Google Cloud load balancer websites and domain websites use the V3 site model for security
policies. Policies and rules operate on the site or asset ID regardless of the website type.
You can use the following security resources with Imperva for Google Cloud sites:

* `incapsula_policy` - Defines an account-level policy.
* `incapsula_policy_asset_association` - Associates a policy with a site asset.
* `incapsula_incap_rule` - Defines rules for a site.
* `incapsula_security_rule_exception` - Defines security rule exceptions for a site.

Policies and rules apply at the site level. For a load balancer website, the same policies and rules apply
to all domains protected by that website. For a domain website, you can apply dedicated security policies
and rules to the domain managed separately by that website.

These resources use the same schema as Cloud WAF and Imperva for AWS sites. See the Cloud WAF and
Imperva for AWS security guides for the full set of arguments.

## Resources that do not apply

The following resources do not apply to Imperva for Google Cloud sites (`type = "PUBLIC_CLOUD"` and `cloud_type = "GCP"`):

* `incapsula_data_center`, `incapsula_data_center_server`, and `incapsula_data_centers_configuration` - Imperva for Google Cloud does not use Cloud WAF data center configuration to route traffic to your origin.
* CDN and caching resources (`incapsula_site_cache_configuration`, `incapsula_cache_rule`) - Imperva for Google Cloud does not use Cloud WAF content delivery or caching.
* Certificate resources (`incapsula_custom_certificate`, `incapsula_managed_certificate_settings`, `incapsula_site_ssl_settings`) - Imperva for Google Cloud does not use Cloud WAF certificate management.

## Example

```hcl
resource "incapsula_site_v3" "gcp_lb" {
  name                  = "my-gcp-lb-site"
  type                  = "PUBLIC_CLOUD"
  cloud_type            = "GCP"
  is_load_balancer_site = true
}

resource "incapsula_incap_rule" "alert_rule" {
  name    = "gcpAlert"
  site_id = incapsula_site_v3.gcp_lb.id
  action  = "RULE_ACTION_ALERT"
  filter  = "Full-URL == \"/login\""
  enabled = true
}
```

