---
layout: "incapsula"
page_title: "Incapsula: data-center"
sidebar_current: "docs-incapsula-resource-data-center"
description: |-
  Provides a Incapsula Data Center resource.
---

# incapsula_data_center

Provides a Incapsula Data Center resource. 

## Example Usage

```hcl
resource "incapsula_data_center" "example-data-center" {
  site_id = "${incapsula_site.example-site.id}"
  name = "Example data center"
  server_address = "8.8.4.4"
  is_content = "yes"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `name` - (Required) The new data center's name.
* `server_address` - (Required) The server's address. Possible values: IP, CNAME.
* `is_enabled` - (Optional) Enables the data center.
* `is_standby` - (Optional) Defines the data center as standby for failover.
* `is_content` - (Optional) The data center will be available for specific resources (Forward Delivery Rules).

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the data center.
