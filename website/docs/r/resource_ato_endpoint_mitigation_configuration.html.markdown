---
layout: "incapsula"
page_title: "Incapsula: incapsula-ato-endpoint-mitigation-configuration"
sidebar_current: "docs-incapsula-resource-ato-endpoint-mitigation-configuration"
description: "Provides an Incapsula ATO endpoint mitigation resource"
---

# incapsula_ato_endpoint_mitigation_configuration

Provides an Incapsula ATO mitigation configuration resource for an endpoint.

## Example Usage

```hcl
resource "incapsula_ato_endpoint_mitigation_configuration" "demo-terraform-ato-endpoint-mitigation-configuration" {
  account_id                     = incapsula_site.example-site.account_id
  site_id                        = incapsula_site.example-site.id
  endpoint_id                    = "5001"
  mitigation_action_for_high_risk                    = "BLOCK"
  mitigation_action_for_medium_risk                  = "BLOCK"
  mitigation_action_for_low_risk                     = "NONE" 
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on. This is required only if the site belongs to the sub account associated with the api key and the api ID 
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `endpoint_id` - (Required) string, Endpoint ID associated with this request.
* `mitigation_action_for_low_risk` - (Required) string, Mitigation action configured for low risk requests - in UPPER CASE.
* `mitigation_action_for_medium_risk` - (Required) string, Mitigation action configured for low risk requests - in UPPER CASE.
* `mitigation_action_for_high_risk` - (Required) string, Mitigation action configured for low risk requests - in UPPER CASE.

##### Mitigation action can be one of : 
  - NONE 
  - CAPTCHA
  - BLOCK
  - TARPIT

## Import

ATO endpoint mitigation configuration can be imported using the account_id/site_id/endpoint_id 

```
$ terraform import incapsula_ato_endpoint_mitigation_configuration.demo-terraform-ato-endpoint-mitigation-configuration 1234/567/89012
```
