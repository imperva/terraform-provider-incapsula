---
layout: "incapsula"
page_title: "Incapsula: data-center-server"
sidebar_current: "docs-incapsula-resource-data-center-server"
description: |-
  Provides a Incapsula Data Center Server resource.
---

# incapsula_data_center_server

Provides a Incapsula Data Center Server resource. 

## Example Usage

```hcl
resource "incapsula_data_center_server" "example-data-center-server" {
  dc_id = "${incapsula_data_center.example-data-center.id}"
  site_id = "${incapsula_site.example-site.id}"
  server_address = "4.4.4.4"
  is_standby = "no"
}
```

## Argument Reference

The following arguments are supported:

* `dc_id` - (Required) Numeric identifier of the data center server to operate on.
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `server_address` - (Optional) The server's address.
* `is_standby` - (Optional) Set the server as Active (P0) or Standby (P1).
* `is_enabled` - (Optional) Enables the data center server.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the data center server.
