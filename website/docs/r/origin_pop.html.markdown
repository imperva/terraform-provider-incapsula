---
subcategory: "Deprecated"
layout: "incapsula"
page_title: "incapsula_origin_pop"
description: |-
  Provides a Incapsula Data Center Origin POP association resource.
---

-> DEPRECATED: incapsula_origin_pop

This resource has been DEPRECATED. It will be removed in a future version. 
Please use the current `incapsula_data_centers_configuration` resource instead.

## Example Usage

```hcl
resource "incapsula_data_center" "example-data-center" {
  site_id = incapsula_site.example-site.id
  name = "Example data center"
  server_address = "8.8.4.4"
  is_content = "true"
}

resource "incapsula_origin_pop" "aws-east" {
  dc_id = incapsula_data_center.example-data-center.id
  site_id = incapsula_site.example-site.id
  origin_pop = "iad"
}
```

## Argument Reference

The following arguments are supported:

* `dc_id` - (Required) Numeric identifier of the data center.
* `origin_pop` - (Required) The Origin POP code (must be lowercase), e.g: `iad`. Note, this field is create/update only. Reads are not supported as the API doesn't exist yet. Note that drift may happen.
* `site_id` - (Required) Numeric identifier of the site to operate on.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier for the Origin POP association.

## Import

Origin Pop can be imported using the `site_id` and `dc_id` separated by /, e.g.:

```
$ terraform import incapsula_origin_pop.aws-east site_id/dc_id
```