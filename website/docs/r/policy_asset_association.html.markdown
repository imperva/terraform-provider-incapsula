---
layout: "incapsula"
page_title: "Incapsula: policy-asset-association"
sidebar_current: "docs-incapsula-resource-policy-asset-association"
description: |-
  Provides a Incapsula Policy Asset Association resource.
---

# incapsula_policy_asset_association

Provides a Incapsula Policy Asset Association resource. 

## Example Usage

```hcl
resource "incapsula_policy_asset_association" "example-policy-asset-association" {
  policy_id  = incapsula_policy.example-policy.id
  asset_id   = incapsula_site.example-site-dns.id 
  asset_type = "WEBSITE"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The Policy ID for the asset association.
* `asset_id` - (Required) The Asset ID for the asset association. Only type of asset supported at the moment is site.
* `asset_type` - (Required) The Policy type for the asset association. Only value at the moment is `WEBSITE`.
* `account_id` - (Optional) The Asset's Account ID. Set this field in case the asset's account is different from the account used in the credentials. e.g setting a sub account asset association set from the parent account

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the policy asset association.

## Import

Policy can be imported using the `policy_id`, `asset_id` and `asset_type` e.g.:

```
$ terraform import incapsula_policy_asset_association.example-policy-asset-association policy_id/asset_id/asset_type
```

