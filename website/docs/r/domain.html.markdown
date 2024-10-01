---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_domain"
description: |- 
  Provides an Incapsula Domain resource for V3 Sites.
---

# incapsula_domain

Provides an Incapsula Domain resource for V3 Sites.
The provider will add/delete domains to/from an Imperva site, based on this resource.
These domains are protected by Imperva and share the website settings and configuration of the onboarded website. Legitimate traffic for all verified domains is allowed.

## Example Usage

```hcl
resource "incapsula_domain" "incapsula_domain-1111_example-domain.com" {
    site_id = 1111
    domain = "example-domain.com"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site.
* `domain` - (Required) The fully qualified domain name of the site.

## Attributes Reference

The following attributes are exported:

* `site_id` - Numeric identifier of the site.
* `domain` - The fully qualified domain name of the site.
* `cname_redirection_record` - CNAME validation code for traffic redirection.  Point your domain's DNS to this record in order to forward the traffic to Imperva.
* `status` - Status of the domain. Indicates if domain DNS is pointed to Imperva's CNAME. Options: BYPASSED, VERIFIED, PROTECTED, MISCONFIGURED.
* `a_records` - A Records for traffic redirection. Point your apex domain's DNS to this record in order to forward the traffic to Imperva.

## Import

Domains can be imported using the site_id and domain_id, e.g. when site_id is 1111 and domain_id is 2222:

```
$ terraform import incapsula_domain.incapsula_domain_1111_example-domain.com 1111/2222
```

## Limitations
### Maximum domains per Imperva site: 
As per the official Website Domain Management feature, the number of domains permitted for each website is
up to 1,000.<br />
Note: This includes auto-discovered domains.<br />
The official docs for Website Domain Management API are located here: https://docs.imperva.com/bundle/cloud-application-security/page/website-domain-api-definition.htm


