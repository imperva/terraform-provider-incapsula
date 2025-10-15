---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_domain"
description: |- 
  Provides an Incapsula Domain resource for V3 Sites.
---

# incapsula_domain

Provides an Incapsula Domain resource for V3 Sites.
The provider will add/delete domains to/from an Imperva site, based on this resource.
These domains are protected by Imperva and share the website settings and configuration of the onboarded website. Legitimate traffic for all verified domains is allowed.

Note:
This resource applies only to sites managed by the incapsula_site_v3 resource. For sites managed by the incapsula_site resource, please use the incapsula_site_domain_configuration resource instead.
Adding an apex domain without its corresponding www subdomain is not supported.

## Example Usage

```hcl
resource "incapsula_site_v3" "example_site" {
  name              = "example-site.com"
}

resource "incapsula_domain" "example_domain" {
    site_id = incapsula_site_v3.example_site.site_id
    domain = "example-domain.com"
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site.
* `domain` - (Required) The fully qualified domain name of the site.

## Attributes Reference

The following attributes are exported:

* `cname_redirection_record` - CNAME validation code for traffic redirection.  Point your domain's DNS to this record in order to forward the traffic to Imperva.
* `status` - Status of the domain. Indicates if domain DNS is pointed to Imperva's CNAME. Options: BYPASSED, VERIFIED, PROTECTED, MISCONFIGURED.
* `a_records` - Will appear for apex domains only. A Records for traffic redirection. Point your apex domain's DNS to this record in order to forward the traffic to Imperva.

## Import

Domains can be imported using the site_id and domain_id separated by /, e.g.:

```
$ terraform import incapsula_domain.example_domain site_id/domain_id
```

## Limitations
### Maximum domains per Imperva site: 
As per the official Website Domain Management feature, the number of domains permitted for each website is
up to 1,000.<br />
Note: This includes auto-discovered domains.<br />
The official docs for Website Domain Management API are located here: https://docs.imperva.com/bundle/cloud-application-security/page/website-domain-api-definition.htm


