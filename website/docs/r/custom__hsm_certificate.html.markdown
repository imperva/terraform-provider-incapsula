---
layout: "incapsula"
page_title: "Incapsula: custom-hsm-certifcate"
sidebar_current: "docs-incapsula-resource-custom-hsm-certificate"
description: |-
  Provides a Incapsula Custom HSM Certificate resource.
---

# incapsula_custom_hsm_certificate

Provides an Incapsula Custom HSM Certificate resource.
The certificate content must be in base64 format.

## Example Usage

```hcl
resource "incapsula_custom_certificate" "custom-certificate" {
    site_id = incapsula_site.example-site.id
    certificate = filebase64("${"path/to/your/cert.crt"}")
    api_detail {
      api_id = "345345-dfg44534-d34534tdfg-dsf4435rg" 
      api_key = "Mdrghg56G5dfHER445hjy5Ghhfg5rth5435hkj3hgd8r7ty948rjslkfhiu4how3hrioeuhtiuer"
      hostname = "api.amer.smartkey.io"
    }
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `certificate` - (Required) The certificate file in base64 format. You can use the Terraform HCL `file` directive to pull in the contents from a file. You can also inline the certificate in the configuration.
* `input_hash` - (Optional) Currently ignored. If terraform plan flags this field as changed, it means that any of: `certificate`, `site_id`, or `api_detail` has changed.
* `api_id` - The key ID. This is the UUID of the Fortanix security object.
* `api_key` - The API key. This is the REST API authentication key from the Fortanix application you created.
* `hostname` - The hostname. This is the location of your assets in the HSM service. In this case, it's the URI (host name) of the Fortanix region as it appears in the security object. For example, api.amer.smartkey.io.

## Attributes Reference

Custom HSM Certificate cannot be exported.

## Import

Custom HSM Certificate cannot be imported.