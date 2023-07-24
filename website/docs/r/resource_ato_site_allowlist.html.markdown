---
layout: "incapsula"
page_title: "Incapsula: incapsula-ato-site-allowlist"
sidebar_current: "docs-incapsula-resource-ato-site-allowlist"
description: "Provides an Incapsula ATO site allowlist configuration resource"
---

# incapsula_ato_site_allowlist

Provides an Incapsula ATO site allowlist configuration resource.

## Example Usage

```hcl
resource "incapsula_ato_site_allowlist_configuration" "demo-terraform-ato-site-allowlist-configuration" {
  account_id      = incapsula_site.example-site.account_id
  site_id         = incapsula_site.example-site.id
  allowlist       = [ { "ip": "192.10.20.0", "mask": "24", "desc": "Test IP 1" }, { "ip": "192.10.20.1", "mask": "8", "desc": "Test IP 2" } ]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) Numeric identifier of the account to operate on. This is required only if the site belongs to the sub account associated with the api key and the api ID 
* `site_id` - (Required) Numeric identifier of the site to operate on.
* `allowlist` - (Required) Array of [AllowlistItem](#allowlistitem) objects

## Object definitions 

#### AllowlistItem

* `ip`   :  (required) string. IP address to exclude. You can use either IPv4 (e.g. 50.3.183.2) or normalized IPv6 representation (e.g. 2001:db8:0:0:1:0:0:1).
  - example: "192.10.20.0"  
* `mask` :  (optional) string. IP subnet mask to use for excluding a range of IPs. This is the number of bits to use from the IP address as a subnet mask to apply on the source IP of incoming traffic.
  - example: "24" 
* `desc` :  (optional) string. Reason for adding this entry to the allowlist  
  - example: "My own IP to always allow Description of the IP/subnet." 

## Import

ATO Site allowlist configuration can be imported using the site_id 

```
$ terraform import incapsula_ato_site_allowlist.demo-terraform-ato-site-allowlist 1234
```
