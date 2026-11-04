---
subcategory: "Site Management"
layout: "incapsula"
page_title: "Incapsula: site"
description: |-
    Provides an Incapsula Site data source.
---

# incapsula_site

Provides a way to look up a site by name and retrieve its external ID. The returned site ID can be used in other Terraform resources across different files and states. For Imperva for Google Cloud load balancer websites, the data source also returns the onboarding parameters, including the `authority_header`.

## Example Usage

```hcl
data "incapsula_site" "example" {
  name = "mysite.example.com"
}

```

### Search in Sub-Account

```hcl
data "incapsula_site" "subaccount" {
  name       = "mysite.example.com"
  account_id = 1234568
}
```

### Imperva for Google Cloud Load Balancer Website

```hcl
data "incapsula_site" "lb" {
  name = "my-gcp-lb-site"
}

output "authority_header" {
  value = data.incapsula_site.lb.authority_header
}
```

## Argument Reference

* `name` - (Required) The name of the site to look up.
* `account_id` - (Optional) ID of the subaccount to search. If omitted, searches the account associated with your Incapsula provider credentials. This value is required only when searching a subaccount.

## Attributes Reference

* `site_id` - The external ID of the site.
* `cname` - The CNAME record of the site.
* `type` - The type of the site.
* `cloud_type` - The cloud provider type of the site, `AWS` or `GCP`. Empty for sites that are not `PUBLIC_CLOUD`.
* `is_load_balancer_site` - Indicates whether the site is an Imperva for Google Cloud load balancer website. `false` for domain websites, legacy Imperva for Google Cloud default sites, and sites that are not part of Imperva for Google Cloud.
* `authority_header` - The authority header assigned by Imperva for an Imperva for Google Cloud load balancer website. Available for load balancer websites and legacy Imperva for Google Cloud default sites. Empty for domain websites and sites that are not part of Imperva for Google Cloud.
