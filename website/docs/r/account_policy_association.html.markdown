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

### Basic Usage - Account Policy Association for account before migration - `default_waf_policy_id` is not set
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
}
```

### Basic Usage - Account Policy Association for account after migration - `default_waf_policy_id` is set
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
* `account_id` - (Required) The account to operate on.
* `default_waf_policy_id` - (Optional)  The WAF policy which is set as default for the account. The account can only have 1 such ID.
  The Default policy will be applied automatically to sites that are created after setting it to default.
  This default setting can be set for the current account, or if used by users with credentials of the parent account can also be set for sub-accounts.
  This parameter is MANDATORY For customers that have account level WAF RULES policies enabled. This means that a default WAF RULES policy resource must be created.
  For customers who have not migrated yet, this parameter should not be set. HOWEVER, once migration occurs, the above is true, a WAF RULES policy must be added and set as default.
  Default setting - Non
* `default_non_mandatory_policy_ids` - (Optional)  This list is currently relevant to Allow lists and ACL policies. More than one policy can be set as default.
  The default policies can be set for the current account, or if used by users with credentials of the parent account can also be set for sub-accounts.
  Default setting – empty list. No default policy.

## Import

Account Policy Association can be imported using the `id` (account ID), e.g.:
```
$ terraform import incapsula_account_policy_association.example-account-policy-association 1234
```