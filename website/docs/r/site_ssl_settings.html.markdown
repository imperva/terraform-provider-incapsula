---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_site_ssl_settings"
description: |- 
  Provides an Incapsula Site SSL Settings resource.
---
# incapsula_site_ssl_settings

Provides an Incapsula Site SSL Settings resource.

If you run the same resource from a site for which SSL is not yet enabled and **approved** will result in the following error response:
- `status:` 406 
- `message:` Site does not have SSL configured
- To enable this feature for your site, you must first configure its SSL settings including a valid certificate.

For more information what HSTS is click [here](https://www.imperva.com/blog/hsts-strict-transport-security/).

## Example Usage

```hcl
resource "incapsula_site_ssl_settings" "example"  {
  site_id = incapsula_site.mysite.id
  
  hsts {
    is_enabled = true
    max_age = 86400
    sub_domains_included = true
    pre_loaded = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `hsts` - (Optional): HTTP Strict Transport Security (HSTS) configuration settings for the site.
    - Type: `set` of `hsts_config` resource (defined below)

## Schema of `hsts_config` resource

The `hsts_config` resource represents the configuration settings for HTTP Strict Transport Security (HSTS).

* `is_enabled` - (Optional): Whether HSTS is enabled for the site.
    - Type: `bool`
    - Default: `false`
* `max_age` - (Optional): The maximum age, in seconds, that the HSTS policy should be enforced for the site.
    - Type: `int`
    - Default: `31536000` (1 year)
* `sub_domains_included` - (Optional): Whether sub-domains should be included in the HSTS policy.
    - Type: `bool`
    - Default: `false`
* `pre_loaded` - (Optional): Whether the site is preloaded in the HSTS preload list maintained by browsers.
    - Type: `bool`
    - Default: `false`

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Site SSL settings. The id is identical to Site id.

## Import

Site SSL settings can be imported using the `id`:
```
terraform import incapsula_site_ssl_settings.example 1234
```



