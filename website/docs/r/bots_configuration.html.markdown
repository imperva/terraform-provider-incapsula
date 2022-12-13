---
layout: "incapsula"
page_title: "Incapsula: bots-configuration"
sidebar_current: "docs-incapsula-resource-bots-configuration"
description: |-
  Provides a Incapsula BOT Access Control Configuration resource.
---

# incapsula_bots_configuration

Provides an Incapsula BOT Access Control Configuration resource.
Each Site has a Good Bots list and a Bad Bots list already configured. This resource allows you to customize them.
<br/>
<strong>canceled_good_bots</strong> list is used to cancel (uncheck in UI) the default Good Bots.
<br/>
<strong>bad_bots</strong> list is used to customize additional Bad Bots
from Imperva’s predefined list of bad bots.

In order to get the latest list, use the <b>/api/integration/v1/clapps</b> found in the <b>Integration</b> section of the 
[Cloud Application Security v1/v3 API Definition page.](https://docs.imperva.com/bundle/cloud-application-security/page/cloud-v1-api-definition.htm)


## Example Usage

### Basic Usage - Lists

The basic usage is to use lists of client application ids.

```hcl
resource "incapsula_bots_configuration" "example-basic-bots-configuration" {
  site_id = incapsula_site.example-basic-site.id
  canceled_good_bots = [6, 17]
  bad_bots = [1, 62, 245, 18]
}
```

### Data Sources Usage

Using `incapsula_client_apps_data` data sources we can use Client Application names that is more "human-readable".

Both lists (canceled_good_bots and bad_bots) can access the data sources in 2 ways:
* `ids` - Contains the Ids of each Client Application name set in `filter` argument (if set in the data source)
* `map` - Contains all the Client Application (names to ids map)

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

### Combination Usage

We always can combine both usages; list and datasource in the same resource or map and list in the same block.

```hcl
data "incapsula_client_apps_data" "client_apps_bad_bots" {
}

resource "incapsula_bots_configuration" "example-basic-bots-configuration" {
  site_id = incapsula_site.example-basic-site.id
  
  canceled_good_bots = [6, 17]  

  bad_bots = [
        data.incapsula_client_apps_data.client_apps_bad_bots.map["Google Translate"],
        data.incapsula_client_apps_data.client_apps_bad_bots.map["Googlebot"],
        530, 531, 537
  ]
}
```



## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.

* `canceled_good_bots` - (Optional) List of Bot IDs taken from Imperva’s predefined list of bad bots

  Default value is an empty list - Will restore the default Canceled Good Bots list

* `bad_bots` - (Optional) List of Bot IDs taken from Imperva’s predefined list of bad bots

  Default value is an empty list - Will restore the default Bad Bots list (empty list)

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the bots configuration.

## Import

Bots Configuration can be imported using the `id` -
the id is identical to Site id, e.g.:
```
$ terraform import incapsula_bots_configuration.demo 1234
```