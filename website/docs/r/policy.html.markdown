---
layout: "incapsula"
page_title: "Incapsula: policy"
sidebar_current: "docs-incapsula-resource-policy"
description: |-
  Provides a Incapsula Policy resource.
---

# incapsula_policy

Provides a Incapsula Policy resource. 

## Example Usage

```hcl
resource "incapsula_policy" "example-policy" {
  name        = "Example policy"
  enabled     = true 
  policy_type = "ACL"
  description = "Example policy description" 
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The policy name.
* `enabled` - (Required) Enables the policy.
* `policy_type` - (Required) The policy type. Possible values: ACL, WHITELIST.
* `account_id` - (Optional) Account ID of the policy.
* `description` - (Optional) The policy description.

## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the policy.
* `account_id` - Account ID of the policy.
