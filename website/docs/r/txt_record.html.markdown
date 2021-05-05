---
layout: "incapsula"
page_title: "Incapsula: txt-record"
sidebar_current: "docs-incapsula-resource-txt-record"
description: |-
  Provides a Incapsula TXT Record(s) association resource.
---

# incapsula_txt_record

Provides a TXT Record(s) association resource. 

## Example Usage

```hcl
resource "incapsula_site" "example-site" {
  domain = "www.example.com"
}

resource "incapsula_txt_record" "test" {
  site_id = incapsula_site.example-site.id
  txt_record_value_one = "test1"
  txt_record_value_two = "test2"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site.
* `txt_record_value_one` - (Optional) New value for txt record number one.
* `txt_record_value_two` - (Optional) New value for txt record number two.
* `txt_record_value_three` - (Optional) New value for txt record number three.
* `txt_record_value_four` - (Optional) New value for txt record number four.
* `txt_record_value_five` - (Optional) New value for txt record number five.


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier for the TXT Records association.

## Import

TXT Record(s) cannot be imported.