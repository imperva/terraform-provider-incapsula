---
layout: "incapsula"
page_title: "Incapsula: client-applications"
sidebar_current: "docs-incapsula_client_apps_data"
description: |-
Provides an Incapsula Client Applications data source.
---

# incapsula_client_apps_data

Provides the ability to use "human-readable" strings identifying the different applications.<p>
In order to get the latest list, use the <b>/api/integration/v1/clapps</b> found in the <b>Integration</b> section of the
[Cloud Application Security v1/v3 API Definition page.](https://docs.imperva.com/bundle/cloud-application-security/page/cloud-v1-api-definition.htm)

Filtering is optional. When used, it will generate the `ids` attribute.
`filter` argument is case-insensitive. If we provide a wrong client application name, it won't fail and just ignore it.

The attribute `map` is always generated and contain all the Client Application (names to ids map).
Using map access with a wrong client application name, it will fail during plan process.

## Example Usage

```hcl
data "incapsula_client_apps_data" "client_apps_canceled_good_bots" {
  filter=["Googlebot","SiteUptime"]
}

data "incapsula_client_apps_data" "client_apps_bad_bots" {
}

resource "incapsula_bots_configuration" "example-basic-bots-configuration" {
  site_id = incapsula_site.example-basic-site.id
  
  canceled_good_bots = data.incapsula_client_apps_data.client_apps_canceled_good_bots.ids

  bad_bots = [
        data.incapsula_client_apps_data.client_apps_bad_bots.map["Google Translate"],
        data.incapsula_client_apps_data.client_apps_bad_bots.map["Googlebot"]
  ]
}
```

## Argument Reference

* `filter` - (Optional) string value - Filter by Client Application names.

  <b>Client Applications names are not unique, so it may return multiple Client Applications ids.</b>


## Attributes Reference

The following attributes are exported:

* `map` - Map of all the client applications where the key is the client name and value is the client ID.
 
  This attribute is always generated, even if you are using `filter` argument.

* `ids` - List of client applications ids filtered by `filter` argument.