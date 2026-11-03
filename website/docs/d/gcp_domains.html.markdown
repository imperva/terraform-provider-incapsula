---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_gcp_domains"
description: |-
  Retrieves the domains associated with a Google Cloud website in Imperva.
---

# incapsula_gcp_domains (Data Source)

Retrieves the domains associated with a Google Cloud website in Imperva. For a load balancer website, the data source returns the domains protected through that website. For a domain website, it returns the domain managed separately by that website.

## Example Usage

```hcl
data "incapsula_gcp_domains" "site_domains" {
  site_id = incapsula_site_v3.gcp-lb-site.id
}

output "domain_names" {
  value = [for d in data.incapsula_gcp_domains.site_domains.domains : d.domain]
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) The numeric identifier of the site whose domains should be listed.
* `account_id` - (Optional) The account ID (`caid`) used to scope the request.

## Attributes Reference

The following attributes are exported:

* `id` - Set to `<account_id>/<site_id>` when `account_id` is set, otherwise to the queried `site_id`.
* `domains` - The domains associated with the specified Google Cloud website. Each element contains:
  - `id` - The numeric identifier of the domain.
  - `domain` - The domain name.
  - `status` - The status of the domain.
  - `cname` - The CNAME assigned to the domain.
  - `creation_date` - The creation time of the domain, in epoch milliseconds.
  - `a_records` - The A records of the domain, when reported.
