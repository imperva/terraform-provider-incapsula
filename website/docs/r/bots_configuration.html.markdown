---
layout: "incapsula"
page_title: "Incapsula: bots-configuration"
sidebar_current: "docs-incapsula-resource-bots-configuration"
description: |-
  Provides a Incapsula BOT Access Control Configuration resource.
---

# incapsula_bots_configuration

Provides an Incapsula BOT Access Control Configuration resource.
Each Site has Good and Bad Bots list already configured. This resource allows you to customize them.
<br/>
<strong>canceled_good_bots</strong> list is used to cancel (uncheck in UI) the default Good Bots.
<br/>
<strong>bad_bots</strong> list is used to customize additional bad bots
<br/>Imperva’s predefined list of bad bots.

In order to get the latest list, use the <b>/api/integration/v1/clapps</b> found in the <b>Integration</b> section of the 
[Cloud Application Security v1/v3 API Definition page.](https://docs.imperva.com/bundle/cloud-application-security/page/cloud-v1-api-definition.htm)


## Example Usage

```hcl
data "incapsula_client_apps_data" "client_apps_canceled_good_bots" {
  filter=["Googlebot","SiteUptime"]
}

data "incapsula_client_apps_data" "client_apps_bad_bots" {
  filter=["Setoozbot","Chinese Bot","Firefox"]
}

resource "incapsula_bots_configuration" "example-basic-bots-configuration" {
  site_id = incapsula_site.example-basic-site.id
  
  canceled_good_bots=data.incapsula_client_apps_data.client_apps_canceled_good_bots.ids

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
* `canceled_good_bots` - (Optional) List of Bot IDs taken from Imperva’s predefined list of bad bots **
* `bad_bots` - (Optional) List of Bot IDs taken from Imperva’s predefined list of bad bots **


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the bots configuration. The id is identical to Site id.

## Import

Bots Configuration can be imported using the `id`, e.g.:

```
$ terraform import incapsula_bots_configuration.demo 1234
```