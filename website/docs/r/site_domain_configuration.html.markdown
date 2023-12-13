---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_site_domain_configuration"
description: |- 
  Provides an Incapsula Site Domains Configuration resource.
---

# incapsula_site_domain_configuration

Provides an Incapsula Site Domain Configuration resource.
The provider will add/delete domains to/from an Imperva site, based on this resource.
Note: The provider is using a single update request, hence domains that exists in the account, but
are missing from the TF file will be deleted.

## Example Usage

```hcl
resource "incapsula_site_domain_configuration" "site-domain-configuration" {
    site_id = incapsula_site.example-site.id
    domain {
      name="example-a.my-web-site.com"
    }
    domain {
      name="example-b.my-web-site.com"
    }
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `domain` - (Required) The fully qualified domain name of the site.

## Attributes Reference

The following attributes are exported:

* `domain` - A list of added domains.
* `site_id` - The id of the site.
* `cname_redirection_record` - CNAME validation code for traffic redirection.  Point your domain's DNS to this record in order to forward the traffic to Imperva

For Each domain the following data will be stored:
  * `id` - The id of the domain.
  * `name` - The address of the domain.
  * `status` - PROTECTED, VERIFIED, BYPASSED, MISCONFIGURED.

## Import

Site Domains Configurations cannot be imported.

## Limitations
### Auto-discovered domains: 
The provider ignores Auto-Discovered domains, hence it will not delete such domains, and it will
not manage them on the TF state.

### Maximum domains per Imperva site: 
As per the official Website Domain Management feature, the number of domains permitted for each website is
up to 1,000.<br />
Note: This includes auto-discovered domains.<br />
The official docs for Website Domain Management API are located here: https://docs.imperva.com/bundle/cloud-application-security/page/website-domain-api-definition.htm


