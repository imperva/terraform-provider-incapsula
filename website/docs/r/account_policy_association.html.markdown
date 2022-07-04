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
data "incapsula_account_data" "account_data" {
}

resource "incapsula_account_policy_association" "example-account-policy-association" {
    account_id                       = data.incapsula_account_data.account_data.current_account
    default_non_mandatory_policy_ids = [
        "123777",
        "123888",
        "123999",
        "123444",
        incapsula.policy.default_acl_policy.id,
        incapsula.policy.default_allowlist_policy.id,
    ]
    default_waf_policy_id            = "1480033"
}
```
## Argument Reference
The following arguments are supported:
* `account_id` - (Mandatory) The account to operate on.
* `default_waf_policy_id` - (Mandatory) The WAF policy which is set as default to the account. The account can only have 1 such id. The Default policy will be applied automatically to sites that were created after setting it to default.
* `default_non_mandatory_policy_ids` - (Optional) This list is currently relevant to whitelist and acl policies. More than one policy can be set as default.

## Import

Account Policy Association can be imported using the `id`, e.g.:
```
$ terraform import incapsula_account_policy_association.example-account-policy-association 1234
```