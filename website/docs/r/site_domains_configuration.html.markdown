---
layout: "incapsula"
page_title: "Incapsula: site-domains-configuration"
sidebar_current: "docs-incapsula-resource-site-domains-configuration"
description:
  Provides an Incapsula Site Domains Configuration resource.
---

# incapsula_site_domain_configuration

Provides an Incapsula Site Domains Configuration resource.
The provider will add/delete domains to/from an Imperva site, based on the resource.
Note: The provider is using a single update request, hence domains that exists on the account, but are missing from the TF file will be deleted.

## Example Usage

```hcl
resource "incapsula_site_domains_configuration" "site-domains-configuration" {
    site_id = incapsula_site.example-site.id
    domain {
      domain_name="example-a.my-web-site.com"
    }
    domain {
      domain_name="example-b.my-web-site.com"
    }
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `domain` - (Required) the address of the domain to add.

## Attributes Reference

The following attributes are exported:

* `domain` - a list of added domains.

For Each domain the following data will be stored:
  * `id` - the id of the domain.
  * `domain_name` - the address of the domain.
  * `status` - PROTECTED, VERIFIED, BYPASSED, MISCONFIGURED.

## Import

Site Domains Configurations cannot be imported.

## Limitations
### Auto-discovered domains: 
The provider ignores Auto-Discovered domains, hence it will not delete such domains, and it will not manage them on the TF state.

### Maximum domains per Imperva site: 
as per the official Website Domain Management, the amount of domains permitted to each website is up to 1,000.
Please note: this includes auto-discovered domains. 

The official docs for Website Domain Management API can be found here: https://docs.imperva.com/bundle/cloud-application-security/page/website-domain-api-definition.htm

