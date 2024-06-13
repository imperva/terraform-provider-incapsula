---
subcategory: "Provider Reference"
layout: "incapsula"
page_title: "incapsula_account_policy_association"
description: |-
  Provides an Incapsula Account Policy Association resource.
---

# incapsula_account_policy_association

Provides an Incapsula Account Policy Association resource.

Dependency is on existing policies, created using the `incapsula_policy` resource.

## Example Usage 

### Basic Usage - Account Policy Association for account before WAF settings are migrated to the WAF RULES policy - `default_waf_policy_id` is not set
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
        incapsula_policy.default_acl_policy.id,
        incapsula_policy.default_allowlist_policy.id,
    ]
}
```

### Basic Usage - Account Policy Association for account after WAF settings are migrated to the WAF RULES policy - `default_waf_policy_id` is set
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
        incapsula_policy.default_acl_policy.id,
        incapsula_policy.default_allowlist_policy.id,
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
  This parameter is MANDATORY for customers that have account level WAF RULES policies enabled. This means that a default WAF RULES policy resource must be created.
  For customers who were not migrated yet, this parameter should not be set. However, after migration, a WAF RULES policy must be added and set as default.
  Default setting - None.
* `default_non_mandatory_policy_ids` - (Optional)  This list is currently relevant to Allowlist and ACL policies. More than one policy can be set as default.
  The default policies can be set for the current account, or if used by users with credentials of the parent account can also be set for sub-accounts.
  Default setting – empty list. No default policy. Providing an empty list or omitting this argument will clear all the non-mandatory default policies.
* `available_policy_ids` - (Optional) Comma separated list of the account’s available policies. These policies can be applied to the websites in the account.
  e.g. available_policy_ids = format(\"%s,%s\", incapsula_policy.acl1-policy.id, incapsula_policy.waf3-policy.id)
  Specify this argument only for a parent account trying to update policy availability for its subaccounts. To remove availability for all policies, specify "NO_AVAILABLE_POLICIES".
  
## Destroy
Destroying this resource will cause the following behavior:
* Default WAF policy will remain unchanged
* Default non-mandatory policies will be unset as default 
* Availability will remain unchanged unless the resource is pointing to a sub account and managed by the parent account. In that case, availability to all policies except for the WAF policy will be removed.
## Import

Account Policy Association can be imported using the `id` (account ID), e.g.:
```
$ terraform import incapsula_account_policy_association.example-account-policy-association 1234
```
