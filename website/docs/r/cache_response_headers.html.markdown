---
layout: "incapsula"
page_title: "Incapsula: cache-response-headers"
sidebar_current: "docs-incapsula-resource-cache-response-headers"
description: |-
  Provides a Incapsula WAF cache response headers resource.
---

# incapsula_cache_response_headers

Provides a Incapsula WAF cache response headers resource. 

## Example Usage

```hcl
resource "incapsula_cache_response_headers" "example-cache-reponse-headers" {
  site_id = "${incapsula_site.example-site.id}"
  cache_headers = "server, x-min-v"
}

resource "incapsula_cache_response_headers" "example-cache-all-response-headers" {
  site_id = "${incapsula_site.example-site.id}"
  cache_all_headers = "true"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `cache_headers` - (Optional) List of header names to be cached in comma separated values.
* `cache_all_headers` - (Optional) Cache all response headers. Pass 'true' or 'false' in the value parameter. Cannot be selected together with cache_headers. Default:false"
