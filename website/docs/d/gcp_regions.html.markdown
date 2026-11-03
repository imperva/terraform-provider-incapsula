---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_gcp_regions"
description: |-
  Retrieves the supported Google Cloud regions for Imperva for Google Cloud.
---

# incapsula_gcp_regions (Data Source)

Retrieves the supported Google Cloud regions for Imperva for Google Cloud. Use a region from this list when retrieving a PSC target service address with the [`incapsula_gcp_psc_target`](gcp_psc_target.html) data source.

## Example Usage

```hcl
data "incapsula_gcp_regions" "all" {}

data "incapsula_gcp_psc_target" "target" {
  region = "us-east1"
}

output "region_supported" {
  value = contains(data.incapsula_gcp_regions.all.regions, "us-east1")
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account ID (`caid`) used to scope the request.

## Attributes Reference

The following attributes are exported:

* `id` - A stable identifier for the data source.
* `regions` - The supported Google Cloud region identifiers, such as `us-east1` and `europe-west1`, sorted alphabetically. The supported regions can change over time.
