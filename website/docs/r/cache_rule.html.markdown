---
layout: "incapsula"
page_title: "Incapsula: cache-rule"
sidebar_current: "docs-incapsula-resource-cache-rule"
description: |-
  Provides a Incapsula Cache Rule resource.
---

# incapsula_cache_rule

Provides a Incapsula Cache Rule resource. 

## Example Usage

```hcl
resource "incapsula_cache_rule" "example-incap-cache-rule" {
  name = "Example cache rule"
  site_id = "${incapsula_site.example-site.id}"
  action = "HTTP_CACHE_MAKE_STATIC"
  filter = "isMobile == Yes"
  enabled = true
  ttl = 3600
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `name` - (Required) Rule name.
* `action` - (Required) Rule action. See the detailed descriptions in the API documentation. Possible values: `HTTP_CACHE_MAKE_STATIC`, `HTTP_CACHE_CLIENT_CACHE_CTL`, `HTTP_CACHE_FORCE_UNCACHEABLE`, `HTTP_CACHE_ADD_TAG`, `HTTP_CACHE_DIFFERENTIATE_SSL`, `HTTP_CACHE_DIFFERENTIATE_BY_HEADER`, `HTTP_CACHE_DIFFERENTIATE_BY_COOKIE`, `HTTP_CACHE_DIFFERENTIATE_BY_GEO`, `HTTP_CACHE_IGNORE_PARAMS`, `HTTP_CACHE_ENRICH_CACHE_KEY`, `HTTP_CACHE_FORCE_VALIDATION`, `HTTP_CACHE_IGNORE_AUTH_HEADER`.
* `filter` - (Required) The filter defines the conditions that trigger the rule action, if left empty, the rule is always run.
* `enabled` - (Required) Boolean that specifies if the rule should be enabled.
* `ttl` - (Optional) TTL in seconds. Relevant for `HTTP_CACHE_MAKE_STATIC` and `HTTP_CACHE_CLIENT_CACHE_CTL` actions.
* `ignored_params` - (Optional) Parameters to ignore. Relevant for `HTTP_CACHE_IGNORE_PARAMS` action. An array containing `'*'` means all parameters are ignored.
* `text` - (Optional) Tag name if action is `HTTP_CACHE_ADD_TAG` action, text to be added to the cache key as suffix if action is `HTTP_CACHE_ENRICH_CACHE_KEY`.
* `differentiate_by_value` - (Optional) Value to differentiate resources by. Relevant for `HTTP_CACHE_DIFFERENTIATE_BY_HEADER`, `HTTP_CACHE_DIFFERENTIATE_BY_COOKIE` and `HTTP_CACHE_DIFFERENTIATE_BY_GEO` actions.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Cache Rule.
