---
layout: "incapsula"
page_title: "Incapsula: site"
sidebar_current: "docs-incapsula-resource-site"
description: |-
  Provides a Incapsula Site resource.
---

# incapsula_site

Provides a Incapsula Site resource. 
Sites are the core resource that is required by all other resources.

## Example Usage

```hcl
resource "incapsula_site" "example-site" {
  domain                 = "examplesite.com"
  account_id             = 123
  ref_id                 = "123"
  send_site_setup_emails = "false"
  site_ip                = "2.2.2.2"
  force_ssl              = "false"
  logs_account_id        = "456"
  data_storage_region    = "US"
  hashing_enabled        = true
  hash_salt              = "foobar"
  log_level              = "full"
  # Performance
  perf_client_comply_no_cache              = true
  perf_client_enable_client_side_caching   = true
  perf_client_send_age_header              = true
  perf_key_comply_vary                     = true
  perf_key_unite_naked_full_cache          = true
  perf_mode_https                          = "include_all_resources"
  perf_mode_level                          = "standard"
  perf_mode_time                           = 1000
  perf_response_cache_300x                 = true
  perf_response_cache_404_enabled          = true
  perf_response_cache_404_time             = 60
  perf_response_cache_empty_responses      = true
  perf_response_cache_http_10_responses    = true
  perf_response_cache_response_header_mode = "custom"
  perf_response_cache_response_headers     = ["Access-Control-Allow-Origin", "Foo-Bar-Header"]
  perf_response_cache_shield               = true
  perf_response_stale_content_mode         = "custom"
  perf_response_stale_content_time         = 1000
  perf_response_tag_response_header        = "Example-Tag-Value-Header"
  perf_ttl_prefer_last_modified            = true
  perf_ttl_use_shortest_caching            = true
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The fully qualified domain name of the site. For example: www.example.com, hello.example.com.
* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `send_site_setup_emails` - (Optional) If this value is false, end users will not get emails about the add site process such as DNS instructions and SSL setup.
* `site_ip` - (Optional) The web server IP/CNAME. This field should be specified when creating a site and the domain does not yet exist or the domain already points to Imperva Cloud. When specified, its value will be used for adding site only. After site is already created this field will be ignored. To modify site ip, please use resource incapsula_data_centers_configuration instead.
* `force_ssl` - (Optional) Force SSL. This option is only available for sites with manually configured IP/CNAME and for specific accounts.
* `logs_account_id` - (Optional) Account where logs should be stored. Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.
* `active` - (Optional) Whether the site is active or bypassed by the Imperva network. Options are `active` and `bypass`.
* `restricted_cname_reuse` - (Optional) Use this option to allow Imperva to detect and add domains that are using the Imperva-provided CNAME (not recommended). One of: true | false.
* `domain_validation` - (Optional) Sets the domain validation method that will be used to generate an SSL certificate. Options are `email`, `html`, and `dns`.
* `approver` - (Optional) Sets the approver e-mail address that will be used to perform SSL domain validation.
* `ignore_ssl` - (Optional) Sets the ignore SSL flag (if the site is in pending-select-approver state). Pass "true" or empty string in the value parameter.
* `acceleration_level` - (Optional) Sets the acceleration level of the site. Options are `none`, `standard`, and `aggressive`.
* `seal_location` - (Optional) Sets the seal location. Options are `api.seal_location.none`, `api.seal_location.bottom_left`, `api.seal_location.right_bottom`, `api.seal_location.left`, and `api.seal_location.right`.
* `domain_redirect_to_full` - (Optional) Sets the redirect naked to full flag. Pass "true" or empty string in the value parameter.
* `remove_ssl` - (Optional) Sets the remove SSL from site flag. Pass "true" or empty string in the value parameter.
* `data_storage_region` - (Optional) The data region to use. Options are `APAC`, `AU`, `EU`, and `US`.
* `hashing_enabled` - (Optional) Specify if hashing (masking setting) should be enabled.
* `hash_salt` - (Optional) Specify the hash salt (masking setting), required if hashing is enabled. Maximum length of 64 characters.
* `log_level` - (Optional) The log level. Options are `full`, `security`, and `none`.
* `naked_domain_san` - (Optional) Use `true` to add the naked domain SAN to a www site’s SSL certificate. Default value: `true`
* `wildcard_san` - (Optional) Use `true` to add the wildcard SAN or `false` to add the full domain SAN to the site’s SSL certificate. Default value: `true`
* `perf_client_comply_no_cache` - (Optional) Comply with No-Cache and Max-Age directives in client requests. By default, these cache directives are ignored. Resources are dynamically profiled and re-configured to optimize performance.
* `perf_client_enable_client_side_caching` - (Optional) Cache content on client browsers or applications. When not enabled, content is cached only on the Imperva proxies.
* `perf_client_send_age_header` - (Optional) Send Cache-Control: max-age and Age headers.
* `perf_key_comply_vary` - (Optional) Comply with Vary. Cache resources in accordance with the Vary response header.
* `perf_key_unite_naked_full_cache` - (Optional) Use the Same Cache for Full and Naked Domains. For example, use the same cached resource for www.example.com/a and example.com/a.
* `perf_mode_https` - (Optional) The resources that are cached over HTTPS, the general level applies. Options are `disabled`, `dont_include_html`, `include_html`, and `include_all_resources`.
* `perf_mode_level` - (Optional) Caching level. Options are `disabled`, `standard`, `smart`, and `all_resources`.
* `perf_mode_time` - (Optional) The time, in seconds, that you set for this option determines how often the cache is refreshed. Relevant for the `include_html` and `include_all_resources` levels only.
* `perf_response_cache_300x` - (Optional) When this option is checked Imperva will cache 301, 302, 303, 307, and 308 redirect response headers containing the target URI.
* `perf_response_cache_404_enabled` - (Optional) Whether or not to cache 404 responses.
* `perf_response_cache_404_time` - (Optional) The time in seconds to cache 404 responses. Value should be divisible by
  60.
* `perf_response_cache_empty_responses` - (Optional) Cache responses that don’t have a message body.
* `perf_response_cache_http_10_responses` - (Optional) Cache HTTP 1.0 type responses that don’t include the Content-Length header or chunking.
* `perf_response_cache_response_header_mode` - (Optional) The working mode for caching response headers. Options are `all` and `custom`.
* `perf_response_cache_response_headers` - (Optional) An array of strings representing the response headers to be cached when working in `custom` mode. If empty, no response headers are cached.
For example: `["Access-Control-Allow-Origin","Access-Control-Allow-Methods"]`.
* `perf_response_cache_shield` - (Optional) Adds an intermediate cache between other Imperva PoPs and your origin servers to protect your servers from redundant requests.
* `perf_response_stale_content_mode` - (Optional) The working mode for serving stale content. Options are `disabled`, `adaptive`, and `custom`.
* `perf_response_stale_content_time` - (Optional) The time, in seconds, to serve stale content for when working in `custom` work mode.
* `perf_response_tag_response_header` - (Optional) Tag the response according to the value of this header. Specify which origin response header contains the cache tags in your resources.
* `perf_ttl_prefer_last_modified` - (Optional) Prefer 'Last Modified' over eTag. When this option is checked, Imperva prefers using Last Modified values (if available) over eTag values (recommended on multi-server setups).
* `perf_ttl_use_shortest_caching` - (Optional) Use shortest caching duration in case of conflicts. By default, the longest duration is used in case of conflict between caching rules or modes. When this option is checked, Imperva uses the shortest duration in case of conflict.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the site.
* `site_creation_date` - Numeric representation of the site creation date.
* `dns_cname_record_name` - The CNAME record name.
* `dns_cname_record_value` - The CNAME record value.
* `dns_a_record_name` - The A record name.
* `dns_a_record_value` - The A record value.
* `domain_verification` - The domain verification (e.g. GlobalSign verification, HTML meta tag).
* `dns_record_name` - The TXT record that needs to be updated with the `domain_verification` value.
* `original_data_center_id` - Numeric representation of the data center created with the site.

## Import

Site can be imported using the `id`, e.g.:

```
$ terraform import incapsula_site.demo 1234
```
