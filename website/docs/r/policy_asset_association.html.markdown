---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_policy_asset_association"
description: |-
  Provides a Incapsula Policy Asset Association resource.
---

# incapsula_policy_asset_association

Provides an Incapsula Policy Asset Association resource. This resource enables you to apply existing policies to assets in your account.

Dependency is on existing policies, created using the `incapsula_policy` resource.

To simplify the use of policies, you can utilize this [cloud-waf Module](https://registry.terraform.io/modules/imperva/cloud-waf/incapsula/latest) along with its submodules.

For full feature documentation, see [Create and Manage Policies](https://docs.imperva.com/bundle/cloud-application-security/page/policies.htm).

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
* `account_id` - (Optional) The account ID of the asset. Set this field if the asset's account is different than the account used in the credentials. For example, when setting a sub account’s asset association from the parent account.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the policy asset association.

## Import

Policy can be imported using the `policy_id`, `asset_id` and `asset_type` e.g.:

```
$ terraform import incapsula_policy_asset_association.example-policy-asset-association policy_id/asset_id/asset_type
```

