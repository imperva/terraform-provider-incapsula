---
layout: "incapsula"
page_title: "Incapsula: bots-configuration"
sidebar_current: "docs-incapsula-resource-bots-configuration"
description: |-
  Provides a Incapsula BOT Access Control Configuration resource.
---

# incapsula_bots_configuration

Provides an Incapsula BOT Access Control Configuration resource. 
Each Site have a Good Bots list already configured and may be changed by calling the API/via Terraform.
<br/>
<strong>canceled_good_bots</strong> list is used to cancel (uncheck in UI) the default Good Bots.
<br/>
<strong>bad_bots</strong> list is used to customize additional bad bots
<br/>
Imperva’s predefined list of bad bots:
https://docs.imperva.com/bundle/cloud-application-security/page/settings/client-classification.htm

## Example Usage

```hcl
resource "incapsula_bots_configuration" "example-basic-bots-configuration" {
  site_id = incapsula_site.example-basic-site.id
  canceled_good_bots = [6,16,17,20]
  bad_bots = [537,530,531,66]
}
```



## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `canceled_good_bots` - (Optional) List of Bot IDs taken from Imperva’s predefined list of bad bots **
* `bad_bots` - (Optional) List of Bot IDs taken from Imperva’s predefined list of bad bots **

** Imperva’s predefined list of bad bots:
https://docs.imperva.com/bundle/cloud-application-security/page/settings/client-classification.htm

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the nots configuration. The id is identical to Site id.

## Import

Bots Configuration can be imported using the `id`, e.g.:

```
$ terraform import incapsula_bots_configuration.demo 1234
```