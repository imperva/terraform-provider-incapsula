---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "Imperva for Google Cloud Site Onboarding with Terraform"
description: |-
  End-to-end guide for onboarding Imperva for Google Cloud sites using Terraform.
---

# Imperva for Google Cloud Site Onboarding with Terraform

This guide walks you through onboarding Imperva for Google Cloud sites using the
Incapsula Terraform provider. It covers creating a **load balancer website** and a **domain website**,
retrieving the site's `authority_header`, discovering supported Google Cloud regions, and retrieving
a PSC target service address.

## 1. Create the load balancer website

`is_load_balancer_site = true` is valid only when `type` is `PUBLIC_CLOUD` and `cloud_type` is `GCP`.
You can set `is_load_balancer_site` only when you create the site. Changing the value for an existing site causes the apply operation to fail, and Terraform does not replace the site.
To intentionally recreate the site, run `terraform apply -replace=<resource address>`.

```hcl
resource "incapsula_site_v3" "gcp_lb" {
  name                  = "my-gcp-lb-site"
  type                  = "PUBLIC_CLOUD"
  cloud_type            = "GCP"
  is_load_balancer_site = true
}

output "authority_header" {
  value = incapsula_site_v3.gcp_lb.authority_header
}
```

The `authority_header` (format `<key>-<accountId>.<domain>`) is populated when the site is read and is used to authenticate the integration with the Imperva SaaS platform.

## 2. Retrieve supported Google Cloud regions

```hcl
data "incapsula_gcp_regions" "all" {}
```

The supported regions are sorted alphabetically and can change over time.

## 3. Retrieve the PSC target service address

The account must have a load balancer website or a legacy Imperva for Google Cloud default site before you retrieve a PSC target service address.
Use `depends_on` to run the lookup after Terraform creates the load balancer website in step 1.

```hcl
data "incapsula_gcp_psc_target" "target" {
  region     = "us-east1"
  depends_on = [incapsula_site_v3.gcp_lb]
}
```

Use `data.incapsula_gcp_psc_target.target.address` as the target service when creating the
PSC Network Endpoint Group (NEG) that provides the private connection between your Google Cloud project and Imperva.

## 4. Create a domain website

A domain website lets you manage a protected domain separately with its own security policies and rules. Set `is_load_balancer_site = false`, or omit the field. The `authority_header` is empty for domain websites.
Add the domain you want to manage separately using the `incapsula_domain` resource.

```hcl
resource "incapsula_site_v3" "gcp_domain" {
  name                  = "my-domain-site"
  type                  = "PUBLIC_CLOUD"
  cloud_type            = "GCP"
  is_load_balancer_site = false
}

resource "incapsula_domain" "gcp_domain" {
  site_id = incapsula_site_v3.gcp_domain.id
  domain  = "www.example.com"
}
```

## 5. Retrieve the site's domains

```hcl
data "incapsula_gcp_domains" "lb_domains" {
  site_id = incapsula_site_v3.gcp_lb.id
}
```

For a load balancer website, the data source returns the domains protected by that website, including auto-discovered domains. For a domain website, it returns the domain managed separately by that website.

## 6. Apply security policies

You can apply supported security-policy resources to Imperva for Google Cloud sites. See the
[Imperva for Google Cloud security guide](gcp_security.html) for details.
