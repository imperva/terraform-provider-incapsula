---
layout: "incapsula"
page_title: "Incapsula: policy-asset-association"
sidebar_current: "docs-incapsula-resource-policy-asset-association"
description: |-
  Provides a Incapsula Policy Asset Association resource.
---

# incapsula_policy_asset_association

Provides an Incapsula Policy Asset Association resource. This resource enables you to apply policies to assets in your account.

 

## Example Usage

```hcl
resource "incapsula_policy_asset_association" "example-policy-asset-association" {
  policy_id  = "123456"
  asset_id   = "456789"
  asset_type = "WEBSITE"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required) The Policy ID for the asset association.
* `asset_id` - (Required) The Asset ID for the asset association. Only type of asset supported at the moment is site.
* `asset_type` - (Required) The Policy type for the asset association. Only value at the moment is `WEBSITE`.
* `account_id` - (Optional) The account ID of the asset. Set this field if the asset's account is different than the account used in the credentials. For example, when setting a sub account’s asset association from the parent account.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the policy asset association.

## Import

Policy can be imported using the `policy_id`, `asset_id` and `asset_type` e.g.:

```
$ terraform import incapsula_policy_asset_association.example-policy-asset-association policy_id/asset_id/asset_type
```

