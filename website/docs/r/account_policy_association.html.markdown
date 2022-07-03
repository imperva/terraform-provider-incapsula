---
layout: "incapsula"
page_title: "Incapsula: account-policy-association"
sidebar_current: "docs-incapsula-resource-account-policy-association"
description: |-
Provides an Incapsula Account Policy Association resource.
---
# incapsula_account_policy_association
Provides an Incapsula Account Policy Association resource.
## Example Usage
```hcl
resource "incapsula_account_policy_association" "example-account-policy-association" {
    account_id                       = "51414296"
    default_non_mandatory_policy_ids = [
        "1477685",
        "1479747",
        "1480155",
        "1480322",
        "1480323",
        "402617",
    ]
    default_waf_policy_id            = "1480033"
}
```
## Argument Reference
The following arguments are supported:
* `account_id` - (Mandatory) The account to operate on.
* `default_waf_policy_id` - (Mandatory) The WAF policy which is set as default to the account. The account can only have 1 such id. The Default policy will be applied automatically to sites that were created after setting it to default.
* `default_non_mandatory_policy_ids` - (Optional) This list is currently relevant to whitelist and acl policies. More than one policy can be set as default.
  Account Policy Association can be imported using the `id`, e.g.:
```
$ terraform import incapsula_account_policy_association.example-account-policy-association 1234
```