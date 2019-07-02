---
layout: "incapsula"
page_title: "Provider: Incapsula"
sidebar_current: "docs-incapsula-index"
description: |-
  The Incapsula provider is used to interact with resources supported by Imperva. The provider needs to be configured with the proper credentials before it can be used.
---

# Incapsula Provider

The Incapsula provider is used to interact with resources supported by Imperva. The provider needs to be configured with the proper credentials before it can be used.

The current API that the Incapsula provider is calling requires sequential execution. You can either use `depends_on` or specify the `parallelism` flag. Imperva recommends the later and setting the value to `1`. Example call: `terraform apply -parallelism=1`.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Incapsula provider
provider "incapsula" {
  api_id = "${var.incapsula_api_id}"
  api_key = "${var.incapsula_api_key}"
}

# Create a site
resource "incapsula_site" "example-site" {
  domain = "examplesite.com"
}

# Create a ACL security rule
resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
  rule_id = "api.acl.blacklisted_ips"
  site_id = "${incapsula_site.example-site.id}"
  ips = "192.168.1.1,192.168.1.2"
}
```

## Argument Reference

The following arguments are supported:

* `api_id` - (Required) The Incapsula API id associated with the account. This can also be
  specified with the `INCAPSULA_API_ID` shell environment variable.
* `api_key` - (Required) The Incapsula API key. This can also be specified with the 
  `INCAPSULA_API_KEY` shell environment variable.
