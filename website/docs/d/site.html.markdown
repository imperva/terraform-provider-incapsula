---
subcategory: "Site Management"
layout: "incapsula"
page_title: "Incapsula: site"
description: |-
    Provides an Incapsula Site data source.
---

# incapsula_site

Provides a way to look up a site by name and retrieve its external ID. The returned site ID can be used in other Terraform resources across different files and states.

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

## Argument Reference

* `name` - (Required) The name of the site to look up.
* `account_id` - (Optional) Sub-account ID to search in.

## Attributes Reference

* `site_id` - The external ID of the site.
* `cname` - The CNAME record of the site.
* `type` - The type of the site.
