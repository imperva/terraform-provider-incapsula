---
layout: "incapsula"
page_title: "Incapsula: incapsula-ato-site-mitigation-configuration"
sidebar_current: "docs-incapsula-resource-ato-site-mitigation-configuration"
description: |- Provides an Incapsula ATO site allowlist resource.
---

# incapsula_ato_site_mitigation_configuration

Provides an Incapsula ATO site allowlist configuration resource.

## Example Usage

```hcl
resource "incapsula_ato_site_mitigation_configuration" "demo-terraform-ato-site-mitigation-configuration" {
  account_id                     = incapsula_site.example-site.account_id
  site_id                        = incapsula_site.example-site.id
  mitigation_configuration       = [ { "endpointId": "5000", "lowAction": "NONE", "mediumAction": "CAPTCHA", "highAction": "BLOCK" }, { "endpointId": "5001", "lowAction": "NONE", "mediumAction": "CAPTCHA", "highAction": "TARPIT" } ] 
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on. This is required only if the site belongs to the sub account associated with the api key and the api ID 
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `mitigation_configuration` - (Required) Array of MitigationConfigurationItem which consists of :
  - endpointId:   string, readOnly: true, Endpoint ID associated with this request.
  - lowAction	  string, readOnly: true, Mitigation action configured for low risk requests - in UPPER CASE.
  - mediumAction  string, readOnly: true, Mitigation action configured for low risk requests - in UPPER CASE.
  - highAction	  string, readOnly: true, Mitigation action configured for low risk requests - in UPPER CASE.

##### Mitigation action can be one of : 
  - NONE 
  - CAPTCHA
  - BLOCK
  - TARPIT


## Import

ATO Site mitigation configuration can be imported using the site_id 

```
$ terraform import incapsula_ato_site_mitigation_configuration.demo-terraform-ato-site-mitigation-configuration 1234
```
