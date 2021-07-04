---
layout: "incapsula"
page_title: "Incapsula: data-centers-configuration"
sidebar_current: "docs-incapsula-resource-data-centers-configuration"
description: |-
  Provides a Incapsula Data Centers Configuration resource.
---

# incapsula_data_centers_configuration

Provides a Incapsula Data Centers Configuration resource. 
Each Site must have one or more Data Centers.
Each Data Center must have one or more Origin Servers
Both Load Balancing and Failover can be configured for both Data Centers and Origin Servers.
Each Data Center can be assigned to serve a specific list of Geo locations.
Or be dedicated to serve requests that were routed using AD Forward Rules.

## Example Usage

### Basic Usage - Single Data Center with one Active and one Standby Origin Servers

```hcl
resource "incapsula_data_centers_configuration" "example-basic-data-centers-configuration" {
  site_id = incapsula_site.example-basic-site.id
  site_topology = "SINGLE_DC"

  data_center {
    name = "New DC"
    ip_mode = "MULTIPLE_IP"

    origin_server {
      address = "54.74.193.120"
      is_active = true
    }

    origin_server {
      address = "44.72.103.175"
      is_active = false
    }

  }

}
```

### Multiple Data Centers across Geo Locations

```hcl
resource "incapsula_data_centers_configuration" "example-geo-assigned-data-centers-configuration" {
  site_id = incapsula_site.example-geo-assigned-site.id
  is_persistent = true
  site_lb_algorithm = "GEO_PREFERRED"
  site_topology = "MULTIPLE_DC"

  data_center {
    name = "Rest of the world DC"
    ip_mode = "MULTIPLE_IP"
    is_rest_of_the_world = true
    origin_pop = "hkg"

    origin_server {
      address = "55.66.77.123"
      is_active = true
      is_enabled = true
    }

  }

  data_center {
    name = "EMEA DC"
    geo_locations = "AFRICA,EUROPE,ASIA"
    origin_pop = "lon"

    origin_server {
      address = "54.74.193.120"
      is_active = true
    }

  }

  data_center {
    name = "Americas DC"
    geo_locations = "US_EAST,US_WEST"
    origin_pop = "iad"

    origin_server {
      address = "54.90.145.67"
      is_active = true
      is_enabled = true
    }

  }

}
```

### Multiple Data Centers with different capacities plus a dedicated Data Center for handling AD forward rules' traffic 

```hcl
resource "incapsula_data_centers_configuration" "example-geo-assigned-data-centers-configuration" {
  site_id = incapsula_site.example-geo-assigned-site.id
  is_persistent = true
  site_lb_algorithm = "WEIGHTED_LB"
  site_topology = "MULTIPLE_DC"

  data_center {
    name = "AD Forward Rules DC"
    is_content = true

    origin_server {
      address = "55.66.77.123"
    }

  }

  data_center {
    name = "Powerful DC"
    weight = 67

    origin_server {
      address = "54.74.193.120"
      is_active = true
    }

  }

  data_center {
    name = "Slagish DC"
    weight = 33
    dc_lb_algorithm = "WEIGHTED"

    origin_server {
      address = "54.90.145.67"
      weight = 50
    }

    origin_server {
      address = "54.90.145.68"
      weight = 30
    }

    origin_server {
      address = "54.90.145.69"
      weight = 20
    }

  }

}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `site_topology` - (Optional) One of: SINGLE_SERVER (no failover), SINGLE_DC (allows failover and LB), or MULTIPLE_DC (allows also Geo and/or AD Forward rules assignment)
* `site_lb_algorithm` - (Optional) How to load balance between multiple Data Centers. One of: BEST_CONNECTION_TIME, GEO_PREFERRED, GEO_REQUIRED, WEIGHTED_LB.
* `fail_over_required_monitors` - (Optional) How many Imperva PoPs should assess Data Center as down before failover is performed. One of: ONE, MANY, MOST, ALL.
* `min_available_servers_for_dc_up` - (Optional) The minimal number of available data center's servers to consider that data center as UP. Default: 1.
* `kickstart_url` - (Optional) The URL that will be sent to the standby server when Imperva performs failover based on our monitoring. E.g. "https://www.example.com/kickStart".
* `kickstart_user` - (Optional) User name, if required by the kickstart URL.
* `kickstart_password` - (Optional) User name, if required by the kickstart URL.
* `is_persistent` - (Optional) When true (the default) our proxy servers will maintain session stickiness to origin servers by a cookie.

At least one `data_center` sub resource must be defined.
The following Data Center arguments are supported: 

* `name` - (Required) Data Center Name. Must be unique within a Site. 
* `ip_mode` - (Optional) SINGLE_IP supports multiple processes on same origin server each listening to a different port, MULTIPLE_IP (the default) support multiple origin servers all listening to same port.
* `web_servers_per_server` - (Optional) When IP mode = SINGLE_IP, number of web servers (processes) per server. Each web server listens to different port. E.g. when web_servers_per_server = 5, HTTP traffic will use ports 80-84 while HTTPS traffic will use ports 443-447. Default: 1.
* `dc_lb_algorithm` - (Optional) How to load balance between the servers of this data center. One of: LB_LEAST_PENDING_REQUESTS (the default), LB_LEAST_OPEN_CONNECTIONS, LB_SOURCE_IP_HASH, RANDOM, WEIGHTED.
* `weight` - (Optional) When site_lb_algorithm = WEIGHTED_LB, the weight in pecentage of this Data Center. Then, total weight of all Data Centers must be 100.
* `is_enabled` - (Optional) When true (the default), this Data Center is enabled. I.e. can serve requests.
* `is_active` - (Optional) When true (the default), this Data Center is active. When false, this Data center will Standby. Automatic failover will happen only if all active Data Centers are not available.
* `is_content` - (Optional) When true, this Data Center will only serve requests that were routed using AD Forward rules. If true, it must also be enabled.
* `is_rest_of_the_world` - (Optional) When true and site_lb_algorithm = GEO_PREFERRED or GEO_REQUIRED, exactly one data center must have is_rest_of_the_world = true. This data center will handle traffic from any region that is not assigned to a specific data center.
* `geo_locations` - (Optional) Commma separated list of geo regions that this data center will serve. Mandatory if site_lb_algorithm = GEO_PREFERRED or GEO_REQUIRED. E.g. "ASIA,AFRICA". Allowed regions: EUROPE, AUSTRALIA, US_EAST, US_WEST, AFRICA, ASIA, SOUTH_AMERICA, NORTH_AMERICA.
* `origin_pop` - (Optional) The ID of the PoP that serves as an access point between Imperva and the customer’s origin server. E.g. "lax", for Los Angeles. When not specified, all Imperva PoPs can send traffic to this data center. The list of available PoPs is documented at: <https://docs.imperva.com/bundle/cloud-application-security/page/more/pops.htm>.

For each `data_center` sub resource, at least one `origin_server` sub resource must be defined.
The following Origin Server arguments are supported: 

* `address` - (Required) Server Address speciied as: ipv4, ipv6, or DNS server name. 
* `weight` - (Optional) When dc_lb_algorithm = WEIGHTED, the weight in pecentage of this Origin Server. Then, total weight of all Data Center's Origin Servers must be 100.
* `is_enabled` - (Optional) When true (the default), this Origin Server is enabled. I.e. can serve requests.
* `is_active` - (Optional) When true (the default), this Origin Server is active. When false, this Origin Server will Standby until failover is performed.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the data centers configuration. The id is identical to Site id.

## Import

Data Centers Configuration can be imported using the `id`, e.g.:

```
$ terraform import incapsula_data_centers_configuration.demo 1234
```