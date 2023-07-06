---
layout: "incapsula"
page_title: "Incapsula: incapsula-ato-site-mitigation-configuration"
sidebar_current: "docs-incapsula-resource-ato-site-mitigation-configuration"
description: |- Provides an Incapsula ATO site mitigation onfiguration resource.
---

# incapsula_ato_endpoint_mitigation_configuration

Provides an Incapsula ATO site mitigation configuration resource.

## Example Usage

```hcl
resource "incapsula_ato_endpoint_mitigation_configuration" "demo-terraform-ato-site-mitigation-configuration" {
  account_id      = incapsula_site.example-site.account_id
  site_id         = incapsula_site.example-site.id
  allowlist       = [ { "ip": "192.10.20.0", "mask": "24", "desc": "Test IP 1" }, { "ip": "192.10.20.1", "mask": "8", "desc": "Test IP 2" } ]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on. This is required only if the site belongs to the sub account associated with the api key and the api ID 
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `allowlist` - (Required) Array of AllowlistItem which consists of
  - ip : string. example: 192.10.20.0 IP address to exclude. This will be either an IPv4 (e.g. 50.3.183.2) or normalized IPv6 representation (e.g. 2001:db8:0:0:1:0:0:1).
  - mask :  string. example: 24 [Optional] IP subnet mask to use for excluding a range of IPs. This is the number of bits to use from the IP address as a subnet mask to apply on the source IP of incoming traffic.
  - desc :  string. example: My own IP to always allow Description of the IP/subnet.```


## Import

ATO Site allowlist configuration can be imported using the site_id 

```
$ terraform import incapsula_ato_site_allowlist.demo-terraform-ato-site-allowlist 1234
```
