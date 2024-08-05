---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_site_ssl_settings"
description: |- 
  Provides an Incapsula Site SSL Settings resource.
---
# incapsula_site_ssl_settings

Provides an Incapsula Site SSL Settings resource.

In this resource you can configure:
- HSTS: A security mechanism enabling websites to announce themselves as accessible only via HTTPS. 
For more information about HSTS, click [here](https://www.imperva.com/blog/hsts-strict-transport-security/).
- TLS settings: Define the supported TLS version and cipher suites used for encryption of the TLS handshake between client and Imperva. 
For more information about supported TLS versions and ciphers, click [here](https://docs.imperva.com/bundle/cloud-application-security/page/cipher-suites.htm).

If you run the SSL settings resource from a site for which SSL is not yet enabled and the SSL certificate is not approved, it will result in the following error response:
- `status:` 406 
- `message:` Site does not have SSL configured
- To enable this feature for your site, you must first configure its SSL settings including a valid certificate.

## Example Usage

```hcl
resource "incapsula_site_ssl_settings" "example"  {
  site_id       = incapsula_site.mysite.id
  account_id    = 4321
  
  hsts { 
    is_enabled               = true
    max_age                  = 31536000
    sub_domains_included     = false
    pre_loaded               = false
  }

  inbound_tls_settings {
    configuration_profile = "CUSTOM"

    tls_configuration {
      tls_version     = "TLS_1_2"
      ciphers_support = [
          "TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
          "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
        ]
    }
    tls_configuration {
      tls_version     = "TLS_1_3"
      ciphers_support = [
        "TLS_AES_128_GCM_SHA256",
        "TLS_CHACHA20_POLY1305_SHA256",
        "TLS_AES_256_GCM_SHA384",
      ]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `account_id` - (Optional) Numeric identifier of the account in which the site is located.
* `hsts` - (Optional): HTTP Strict Transport Security (HSTS) configuration settings for the site.
    - Type: `set` of `hsts_config` resource (defined below)
* `inbound_tls_settings` - (Optional): Transport Layer Security (TLS) configuration settings for the site.
  - Type: `set` of `inbound_tls_settings` resource (defined below)

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

## Schema of `inbound_tls_settings` resource

The `inbound_tls_settings` resource represents the configuration settings for Transport Layer Security (TLS).

* `configuration_profile` - (Required): Where to use a pre-defined or custom configuration for TLS settings. Possible values: DEFAULT, ENHANCED_SECURITY, CUSTOM.
  - Type: `string`
* `tls_configuration` - (Optional): List supported TLS versions and ciphers.
  - Type: `List`

### Nested Schema for `tls_configuration`

* `tls_version` - (Required): TLS supported versions.
  - Type: `string`
* `ciphers_support` - (Required): List of ciphers to use for this TLS version.
  - Type: `List`


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the Site SSL settings. The id is identical to Site id.

## Import

Site SSL settings can be imported using the `siteId` or `siteId`/`accountId` for sub-accounts:
```
terraform import incapsula_site_ssl_settings.example 1234
terraform import incapsula_site_ssl_settings.example 1234/4321
```



