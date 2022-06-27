---
layout: "incapsula"
page_title: "Incapsula: site_monitoring"
sidebar_current: "docs-incapsula-resource-site_monitoring"
description: |- Provides a Incapsula Site Monitoring resource.
---

# incapsula_site_monitoring

Configure settings to determine when origin servers should be considered “up” or “down” (active or inactive) by the Imperva Load Balancer. 
Select which failure scenarios you want to produce alarm messages, and how to send them.

Note that destroy action doesn't do any change in Imperva system.

## Example Usage

### Basic Usage - Site Monitoring

```hcl
resource "incapsula_site_monitoring" "example_site_monitoring" {
    id      = "1234"
    site_id = 1234

    failed_request_criteria {
        http_request_timeout       = 1
        http_request_timeout_units = "MINUTES"
        http_response_error        = "501,503"
    }

    monitoring {
        failed_requests_duration       = 2
        failed_requests_duration_units = "MINUTES"
        failed_requests_min_number     = 10
        failed_requests_percentage     = 10
    }

    notifications {
        alarm_on_dc_failover        = false
        alarm_on_server_failover    = true
        alarm_on_stands_by_failover = true
        required_monitors           = "MANY"
    }

    up_down_verification {
        monitoring_url            = "/users"
        up_check_retries          = 5
        up_checks_interval        = 1
        up_checks_interval_units  = "MINUTES"
        use_verification_for_down = false
        expected_received_string  = "example string"
    }
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `monitoring` - (Optional) Monitoring Parameters
* `failed_request_criteria` - (Optional) TCP connection errors always count as failed requests. In addition, the properties below can be configured to define failed requests.
* `up_down_verification` - (Optional) Use verification checks to mark server as "Down"
* `notifications` - (Optional) Set up email alerts for different scenarious

The following parameters are supported for the `monitoring` resource:

* `failed_requests_percentage` - (Optional) The percentage of failed requests to the origin server. Default: 40
* `failed_requests_min_number` - (Optional) The minimum number of failed requests to be considered as failure. Default: 3
* `failed_requests_duration` - (Optional) The minimum duration of failures above the threshold to consider server as down. 20-180 SECONDS or 1-2 MINUTES. Default: 40.
* `failed_requests_duration_units` - (Optional) Time unit. Possible values: SECONDS, MINUTES. Default: SECONDS.


The following parameters are supported for the `failed_request_criteria` resource:
* `http_request_timeout` - (Optional) The maximum time to wait for an HTTP response. 1-200 SECONDS or 1-2 MINUTES. Default: 35
* `http_request_timeout_units` - (Optional) Time unit. Default: SECONDS.
* `http_response_error` - (Optional) The HTTP response error codes or patterns that will be counted as request failures. Default: "501-599".


The following parameters are supported for the `up_down_verification` resource:
* `use_verification_for_down` - (Optional) If Imperva determines that an origin server is down according to failed request criteria, it will initiate another request to verify that the origin server is down. Default: true
* `monitoring_url` - (Optional) The URL to use for monitoring your website. Default: "/"
* `expected_received_string` - (Optional) The expected string. If left empty, any response, except for the codes defined in the HTTP response error codes to be treated as Down parameter, will be considered successful. If the value is non-empty, then the defined value must appear within the response string for the response to be considered successful.
* `up_checks_interval` - (Optional) After an origin server was identified as down, Imperva will periodically test it to see whether it has recovered, according to the frequency defined in this parameter. 10-120 SECONDS or 1-2 MINUTES. Default: 20
* `up_checks_interval_units` - (Optional) Time unit. Default: SECONDS.
* `up_check_retries` - (Optional) Every time an origin server is tested to see whether it’s back up, the test will be retried this number of times. Default: 3

The following parameters are supported for the `notifications` resource:
* `alarm_on_stands_by_failover` - (Optional) Indicates whether or not an email will be sent upon failover to a standby data center. Default: true
* `alarm_on_dc_failover`- (Optional) Indicates whether or not an email will be sent upon data center failover. Default: true
* `alarm_on_server_failover` - (Optional) Indicates whether or not an email will be sent upon server failover. Default: false
* `required_monitors` - (Optional) Monitors required to report server / data center as down. Default: "MOST"

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the data centers configuration. The id is identical to Site id.

## Import

Site Monitoring configuration can be imported using the `id`, e.g.:

```
$ terraform import incapsula_site_monitoring.example_site_monitoring 1234
```