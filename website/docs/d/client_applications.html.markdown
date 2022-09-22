---
layout: "incapsula"
page_title: "Incapsula: client-applications"
sidebar_current: "docs-incapsula_client_apps_data"
description: |-
Provides an Incapsula Client Applications data source.
---

# incapsula_client_apps_data

Provides the list of all the client applications as a map.
There are no filters needed for this data source

## Example Usage


```hcl
data "incapsula_client_apps_data" "client_apps" {
}

resource "incapsula_bots_configuration" "example-basic-bots-configuration" {
  site_id = incapsula_site.example-basic-site.id
  canceled_good_bots = [
        6,
        data.incapsula_client_apps_data.client_apps.map["Firefox"],
        16, 17
  ]
  bad_bots = [
        data.incapsula_client_apps_data.client_apps.map["Google Translate"],
        data.incapsula_client_apps_data.client_apps.map["Googlebot"],
        530, 531, 537
  ]
}
```

## Argument Reference

There are no filters in this resource.

## Attributes Reference

The following attributes are exported:

* `map` - Map of all the client applications where the key is the client name and value is the client ID.