---
layout: "incapsula"
page_title: "Incapsula: data-centers-configuration"
sidebar_current: "docs-incapsula-data-data-center"
description: |-
  Provides an Incapsula Data Center data source.
---

# incapsula_data_center

Provides filtering of a single Data Center in specific site. 
Then its id can be referenced by other resources such as incapsula_incap_rule.

Filtering by site id is mandatory. All other filters are optional. A logical AND is applied on all specified filters.

Exactly one Data Center must be qualified or an error will be raised.

## Example Usage


```hcl
resource "incapsula_data_centers_configuration" "example-two-data-centers-configuration" {
  site_id = incapsula_site.example-site.id
  site_topology = "MULTIPLE_DC"

  data_center {
    name = "AD Forward Rules DC"
    is_content = true

    origin_server {
      address = "55.66.77.123"
    }

  }

  data_center {
    name = "Main DC"

    origin_server {
      address = "54.74.193.120"
    }

  }

}

data "incapsula_data_center" "content_dc" {
  site_id              = incapsula_data_centers_configuration.example-two-data-centers-configuration.id
  filter_by_is_content = true
}

# Incap Rule: Forward to Data Center (ADR)
resource "incapsula_incap_rule" "example-incap-rule-fwd-to-data-center" {
  name    = "Example incap rule forward to data center"
  site_id = incapsula_site.example-site.id
  action  = "RULE_ACTION_FORWARD_TO_DC"
  filter  = "Full-URL == \"/someurl\""
  dc_id   = data.incapsula_data_center.content_dc.id
}

```

If more than one content Data Center is defined for the site, you may add filters. E.g.

```hcl
data "incapsula_data_center" "content_dc" {
  site_id              = incapsula_data_centers_configuration.example-two-data-centers-configuration.id
  filter_by_is_content = true
  filter_by_name       = "AD Forward Rules DC"
}

```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Reference to resource incapsula_data_centers_configuration.*resource_name*.id 
* `filter_by_id` - (Optional) Filters by the numeric unique internal Data Center id.
* `filter_by_name` - (Optional) string value - Filters by Data Center name. DC Name is unique per Site.
* `filter_by_is_enabled` - (Optional) boolean value - Filters by is_enabled.
* `filter_by_is_active` - (Optional) boolean value - Filters by is_active == true. 
* `filter_by_is_standby` - (Optional) boolean value - Filters by is_active == false. 
* `filter_by_is_rest_of_the_world` - (Optional) boolean value - Filters by is_rest_of_the_world. True value will provide the single rest-of-the-world DC.
* `filter_by_is_content` - (Optional) boolean value - Filters by is_content. True value will provide the single content DC (handles only traffic routed by AD forward-to-dc rule).
* `filter_by_geo_location` - (Optional) string value - One of: EUROPE, AUSTRALIA, US_EAST, US_WEST, AFRICA, ASIA, SOUTH_AMERICA, NORTH_AMERICA. If no DC is assigned to handle traffic from the specified region, then qualify the the single rest-of-the-world DC, 

The five boolean filters accept true value only.
To get the standby Data Center: use `filter_by_is_standby = true` (instead of `filter_by_is_active = false` which will be ignored).

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the qualified data center.
* `name` - Data Center Name.
* `ip_mode` - One of: SINGLE_IP, MULTIPLE_IP.
* `lb_algorithm` - How to load balance between the servers of this data center. One of: LB_LEAST_PENDING_REQUESTS, LB_LEAST_OPEN_CONNECTIONS, LB_SOURCE_IP_HASH, RANDOM, WEIGHTED.
* `weight` - The weight in percentage of this Data Center. Populated only when Site's LB algorithm is WEIGHTED_LB.
* `is_enabled` - When true, this Data Center is enabled. I.e. can serve requests.
* `is_active` - When true, this Data Center is active. When false, this Data center will Standby.
* `is_content` - When true, this Data Center will only serve requests that were routed using AD Forward rules.
* `is_rest_of_the_world` - When true and Site's LB algorithm = GEO_PREFERRED or GEO_REQUIRED, this data center will handle traffic from any region that is not assigned to a specific data center.
* `geo_locations` - Comma separated list of geo regions that this data center will serve. Populated only when Site's LB algorithm is GEO_PREFERRED or GEO_REQUIRED. E.g. "ASIA,AFRICA". Allowed regions: EUROPE, AUSTRALIA, US_EAST, US_WEST, AFRICA, ASIA, SOUTH_AMERICA, NORTH_AMERICA.
* `origin_pop` - (Optional) The ID of the PoP that serves as an access point between Imperva and the customer’s origin server. E.g. "lax", for Los Angeles. When not specified, all Imperva PoPs can send traffic to this data center. The list of available PoPs is documented at: <https://docs.imperva.com/bundle/cloud-application-security/page/more/pops.htm>.
