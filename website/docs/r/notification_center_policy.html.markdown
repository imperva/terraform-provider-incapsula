---
subcategory: "Account and User Management"
layout: "incapsula"
page_title: "incapsula_notification_center_policy"
description: |-
  Provides a Incapsula Notification Center Policy resource.
---

# incapsula_notification_center_policy

Provides an Incapsula Notification Center Policy resource.

## Example Usage
Notification center policy that applies on subaccounts

```hcl
resource "incapsula_notification_center_policy" "notification-policy-subaccount" {
  account_id = 12345
  policy_name = "Terraform policy sub account"
  status = "ENABLE"
  sub_category = "SITE_NOTIFICATIONS"
  emailchannel_user_recipient_list = [1111, 2222]
  emailchannel_external_recipient_list=["john.doe@company.com", "another.email@company.com"]      
  policy_type = "SUB_ACCOUNT"
  sub_account_list = [123456, incapsula_subaccount.tmp-subaccount.id]
}
```

Notification policy that applies to assets of type "incapsula_site"
```hcl
resource "incapsula_notification_center_policy" "notification-policy-account-with-assets" {
  account_id = 12345
  policy_name = "Terraform policy account with assets"
  asset {
    asset_type = "SITE"
    asset_id = incapsula_site.tmp-site.id
  }
  asset {
    asset_type = "SITE"
    asset_id = 7999203
  }   
  status = "ENABLE"
  sub_category = "SITE_NOTIFICATIONS"
  emailchannel_user_recipient_list = [1111, 2222]
  emailchannel_external_recipient_list=["john.doe@company.com", "another.exernal.email@company.com"] 
  policy_type = "ACCOUNT"
  apply_to_new_assets = "FALSE"
}
```
Notification policy on sub-category with no relevance to assets
```hcl
resource "incapsula_notification_center_policy" "notification-policy-account-without-assets" { 
  account_id = 12345
  policy_name = "Terraform policy account without assets"
  status = "ENABLE"
  sub_category = "ACCOUNT_NOTIFICATIONS"
  emailchannel_user_recipient_list = [1111, 2222]
  emailchannel_external_recipient_list=["john.doe@company.com", "another.exernal.email@company.com"]      
  policy_type = "ACCOUNT"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to work on.
* `policy_name` - (Required) The name of the policy. Cannot contain special characters
* `status` - (Optional) Indicates whether the policy is enabled or disabled. Possible
  values: `ENABLE` (default value), `DISABLE`.
* `sub_category` - (Required) The sub category of notifications that the policy will apply to. The possible
  values are available via the API at https://api.imperva.com/notification-settings/v3/subtypes.
* `emailchannel_user_recipient_list` - (Optional) List of numeric identifiers of the users from the Imperva account 
  to receive emails notifications. There must be at least one value in this list or in the `emailchannel_external_recipient_list` list.
* `emailchannel_external_recipient_list` - (Optional) List of email addresses (for recipients who are not Imperva users) to receive email notifications.
  There must be at least one value in this list or in the `emailchannel_user_recipient_list` list.
* `apply_to_new_assets` - (Optional) If value is `TRUE`, all newly onboarded assets are automatically added to the
  notification policy's assets list. Possible values: `TRUE`, `FALSE` (default value).\
  We recommend always setting this field's value to `FALSE`, to disable automatic updates of assets on the policy, so you
  have full control over your resources.
* `policy_type` - (Optional) If the value is `ACCOUNT`, the policy will apply only to the current account that is 
  specified by the account_id. If the value is `SUB_ACCOUNT` the policy applies to the sub accounts only.
  The parent account will receive notifications for activity in the sub accounts that are specified in the 
  `sub_account_list` parameter. This `sub_account_list` is available only in accounts that can contain sub accounts.
  Possible values: `ACCOUNT` (default value), `SUB_ACCOUNT`.
* `sub_account_list` - (Optional) List of numeric identifiers of sub accounts of this account for which the parent account will
  receive notifications. Should be set if the `policy_type` is `SUB_ACCOUNT`.
* `apply_to_new_sub_accounts` - (Optional) If value is `TRUE`, all newly onboarded sub accounts are automatically added
  to the notification policy's sub account list. Possible values: `TRUE`, `FALSE` (default value)\
  Relevant if the `policy_type` is `SUB_ACCOUNT`.\
  We recommend always setting this field's value to `FALSE`, to disable automatic updates of sub-accounts on the policy, 
  so you have full control over your resources.

  
Under the following conditions, you need to define at least 1 asset:
If the `policy_type` argument is `ACCOUNT`, and the chosen `sub_category` requires configuration of assets, and the
argument `apply_to_new_assets` is `FALSE`, then at least 1 asset must be defined.\
For example, when configuring a policy for the `SITE_NOTIFICATIONS` `sub_category`, if the argument `apply_to_new_assets` is FALSE, at least one SITE asset must be specified.
The arguments that are supported in `asset` sub resource are:
* `asset_type` - Indicates the Imperva-protected entity that triggers the notification. Possible values: `SITE`, `IP_RANGE`, `EDGE_IP`, `ORIGIN_CONNECTIVITY`,
  `NETFLOW_EXPORTER`, `DOMAIN`.
* `asset_id` - Numeric identifier of the asset.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier for the Notification Policy.

## Import

Notification Policy can be imported using the account_id/policy_id

```
$ terraform import incapsula_notification_center_policy.notification-policy-account-without-assets 12345/9999
```
