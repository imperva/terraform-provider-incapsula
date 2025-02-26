---
subcategory: "Application Performance and Delivery"
layout: "incapsula"
page_title: "incapsula_site_cache_configuration"
description: |-
  Provides an Incapsula Site Cache Configuration resource.
---

# incapsula_site_cache_configuration

Configure content caching for your website.

## Example Usage

### Basic Usage - Site Cache Configuration

```hcl
resource "incapsula_site_cache_configuration" "example_site_cache_configuration" {
	site_id = incapsula_site_v3.example-site-v3.id
	client_comply_no_cache              = true
	client_enable_client_side_caching   = true
	client_send_age_header              = true
	key_comply_vary                     = true
	key_unite_naked_full_cache          = true
	mode_https                          = "include_all_resources"
	mode_level                          = "standard"
	mode_time                           = 1000
	response_cache_300x                 = true
	response_cache_404_enabled          = true
	response_cache_404_time             = 60
	response_cache_empty_responses      = true
	response_cache_http_10_responses    = true
	response_cache_response_header_mode = "custom"
	response_cache_response_headers     = ["Access-Control-Allow-Origin", "Foo-Bar-Header"]
	response_cache_shield               = true
	response_stale_content_mode         = "custom"
	response_stale_content_time         = 1000
	response_tag_response_header        = "Example-Tag-Value-Header"
	ttl_prefer_last_modified            = true
	ttl_use_shortest_caching            = true		
}
```

## Argument Reference

The following arguments are supported:

* `client_comply_no_cache` - (Optional) Comply with No-Cache and Max-Age directives in client requests. By default, these cache directives are ignored. Resources are dynamically profiled and re-configured to optimize performance. **Default:** `false`
* `client_enable_client_side_caching` - (Optional) Cache content on client browsers or applications. When not enabled, content is cached only on the Imperva proxies. **Default:** `false`
* `client_send_age_header` - (Optional) Send Cache-Control: max-age and Age headers. **Default:** `false`
* `key_comply_vary` - (Optional) Comply with Vary. Cache resources in accordance with the Vary response header. **Default:** `false`
* `key_unite_naked_full_cache` - (Optional) Use the Same Cache for Full and Naked Domains. For example, use the same cached resource for _www.example.com/a_ and _example.com/a_. **Default:** `false`
* `mode_https` - (Optional) The resources that are cached over HTTPS, the general level applies. **Default:** `disabled`
Possible values:
  * **disabled** - No resources are cached over HTTPS.
  * **dont_include_html** - Only images, css files, js files, and resources defined with the 'Cache-Control: public' header may be cached over HTTPS.
  * **include_html** - Same as dont_include_html, also HTML files may be cached over HTTPS.
  * **include_all_resources** - All resources may be cached over HTTPS.
* `mode_level` - (Optional) Caching level. **Default:** `smart` 
Possible values:
  * **disabled** - Turn off site caching entirely, including user-defined custom cache rules. **Note:** Since all other attributes require caching, it is recommended to clear all resource attributes from configuration when disabling the cache.
  * **custom_cache_rules_only** - Disable caching, unless specified otherwise by a custom cache rule.
  * **standard** - Cache according to standard HTTP headers.
  * **smart** - Cache according to standard HTTP headers as well as profile dynamic pages to identify and cache static content that was not marked as static. In addition to content that was marked by the site's developer / web server as static using standard HTTP headers, Imperva also profiles other resources to identify and cache static content that was not marked as such.
  * **all_resources** - Cache every resource on the web server for the specified amount of time. All site content is cached.
* `mode_time` - (Optional) The time, in seconds, that you set for this option determines how often the cache is refreshed. Relevant for the `include_html` and `include_all_resources` levels only.
* `response_cache_300x` - (Optional) When this option is checked Imperva will cache 301, 302, 303, 307, and 308 redirect response headers containing the target URI. **Default:** `false`
* `response_cache_404_enabled` - (Optional) Whether or not to cache 404 responses. **Default:** `false`
* `response_cache_404_time` - (Optional) The time in seconds to cache 404 responses. Value should be divisible by 60.
* `response_cache_empty_responses` - (Optional) Cache responses that don’t have a message body. **Default:** `false`
* `response_cache_http_10_responses` - (Optional) Cache HTTP 1.0 type responses that don’t include the Content-Length header or chunking. **Default:** `false`
* `response_cache_response_header_mode` - (Optional) The working mode for caching response headers. **Default:** `disabled`
Possible values:
  * **disabled**
  * **all** - Cache all headers in responses.
  * **custom** - Specify which response headers should be cached along with the resource.
* `response_cache_response_headers` - (Optional) An array of strings representing the response headers to be cached when working in `custom` mode. If empty, no response headers are cached.
For example: `["Access-Control-Allow-Origin","Access-Control-Allow-Methods"]`.
* `response_cache_shield` - (Optional) Adds an intermediate cache between other Imperva PoPs and your origin servers to protect your servers from redundant requests. **Default:** `false`
* `response_stale_content_mode` - (Optional) The working mode for serving stale content. **Default:** `disabled`
Possible values:
  * **disabled** - Disable serving of stale content.
  * **adaptive** - Stale content is served for a duration of 2 to 24 hours based on the time passed since the resource was last updated.
  * **custom** - Serve stale content for the specified amount of time.
* `response_stale_content_time` - (Optional) The time, in seconds, to serve stale content for when working in `custom` work mode.
* `response_tag_response_header` - (Optional) Tag the response according to the value of this header. Specify which origin response header contains the cache tags in your resources.
* `ttl_prefer_last_modified` - (Optional) Prefer 'Last Modified' over eTag. When this option is checked, Imperva prefers using Last Modified values (if available) over eTag values (recommended on multi-server setups). **Default:** `false`
* `ttl_use_shortest_caching` - (Optional) Use shortest caching duration in case of conflicts. By default, the longest duration is used in case of conflict between caching rules or modes. When this option is checked, Imperva uses the shortest duration in case of conflict. **Default:** `false`

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Site Cache configuration. The id is identical to Site id.

## Import

Site Cache Configuration can be imported using the `id`, e.g.:

```
$ terraform import incapsula_site_cache_configuration.example_site_cache_configuration 1234
```
