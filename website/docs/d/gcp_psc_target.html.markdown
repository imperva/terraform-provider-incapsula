---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_gcp_psc_target"
description: |-
  Retrieves the PSC target service address for a Google Cloud region.
---

# incapsula_gcp_psc_target (Data Source)

Retrieves the PSC target service address for a specified Google Cloud region. Use the returned `address` to configure the PSC Network Endpoint Group (NEG) that provides the private connection between your Google Cloud project and Imperva. Use the [`incapsula_gcp_regions`](gcp_regions.html) data source to retrieve the supported regions.

~> **Note:** The account must have a load balancer website (`incapsula_site_v3` with `is_load_balancer_site = true`) or a legacy Imperva for Google Cloud default site that was created automatically before multi-site support.

## Example Usage

```hcl
data "incapsula_gcp_psc_target" "target" {
  region     = "us-east1"
  depends_on = [incapsula_site_v3.gcp-lb-site]
}

output "psc_target" {
  value = data.incapsula_gcp_psc_target.target.address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The Google Cloud region for which you want to retrieve the PSC target service address.
* `account_id` - (Optional) The account ID (`caid`) used to scope the request.

## Attributes Reference

The following attributes are exported:

* `id` - Set to `<account_id>/<region>` when `account_id` is set, otherwise to the requested region.
* `address` - The PSC target service address for the specified Google Cloud region.
