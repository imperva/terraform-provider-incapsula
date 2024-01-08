---
subcategory: "Deprecated"
layout: "incapsula"
page_title: "incapsula_data_center"
description: |-
  Provides a Incapsula Data Center resource.
---

-> DEPRECATED: incapsula_data_center

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
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `name` - (Required) The new data center's name.
* `server_address` - (Required) The server's address. Possible values: IP, CNAME.
* `is_enabled` - (Optional) Enables the data center.
* `is_content` - (Optional) The data center will be available for specific resources (Forward Delivery Rules).

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the data center.

## Import

Data Center can be imported using the `id`, e.g.:

```
$ terraform import incapsula_data_center.demo 1234
```