---
subcategory: "Advanced Bot Protection"
layout: "incapsula"
page_title: "incapsula_abp_websites"
description: |-
  Provides an ABP (Advanced Bot Protection) website resource.

---

# incapsula_abp_websites

Provides an ABP (Advanced Bot Protection) website resource. Allows you to enable and configure ABP for given websites.

This represents the canonical configuration of the ABP website groups and websites, if there are existing website groups or website, or they have different settings those will be removed or changed and not added as additional items.

NOTE: Due to limitations in ABP, the API key/id used to deploy this resource must match the `account_id` used in the resource (API key/id for a parent account do not work). All Incapsula sites associated with the resource must also be defined in that account.

## Example Usage

```terraform
resource "incapsula_abp_websites" "abp_websites" {
    account_id = data.incapsula_account_data.account_data.current_account
    auto_publish = true
    website_group {
        name = "sites-1"
        website {
            site_id = incapsula_site.sites-1.id
            enable_mitigation = false
        }
    }
    website_group {
        name = "sites-2"
        website {
            site_id = incapsula_site.sites-2.id
            enable_mitigation = true
        }
    }
}


resource "incapsula_abp_websites" "abp_websites" {
    account_id = data.incapsula_account_data.account_data.current_account
    auto_publish = true
    website_group {
        name = "sites"
        website {
            site_id = incapsula_site.sites-1.id
            enable_mitigation = false
        }
        website {
            site_id = incapsula_site.sites-2.id
            enable_mitigation = true
        }
    }
    website_group {
        name = "sites" # Duplicate name
        name_id = "sites-2" # name_id can be used to disambiguate names in case of duplicates
        website {
            site_id = incapsula_site.sites-3.id
            enable_mitigation = true
        }
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `account_id` (Number) The account these websites belongs to.

### Optional

- `auto_publish` (Boolean) Whether to publish the changes automatically. Changes don't take take effect until they have been published.
- `website_group` (Block List) List of website groups which are associated to ABP. Website groups are matched in a top-down fashion. If a more specific website group should take precedence over over a wild card entry then that should be higher in the list. (see [below for nested schema](#nestedblock--website_group))

### Read-Only

- `id` (String) The ID of this resource.
- `last_publish` (String) When the last publish was done for this terraform resource. Changes are published when `auto_publish` is true and the terraform config is applied.

<a id="nestedblock--website_group"></a>
### Nested Schema for `website_group`

Required:

- `name` (String) Name for the website group. Must be unique unless `name_id` is specified.

Optional:

- `name_id` (String) Unique user-defined identifier used to differentiate website groups whose `name` is identical
- `website` (Block List) List of websites within the website group. Websites are matched in a top-down fashion. If a more specific website should take precedence over over a wild card entry then that should be higher in the list (see [below for nested schema](#nestedblock--website_group--website))

Read-Only:

- `id` (String) The ID of this resource.

<a id="nestedblock--website_group--website"></a>
### Nested Schema for `website_group.website`

Required:

- `incapsula_site_id` (Number) Which `incapsula_site` this website refers to

Optional:

- `enable_mitigation` (Boolean) Enables the ABP conditions for this website. Defaults to true.

Read-Only:

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax:

```shell
terraform import incapsula_abp_websites.websites 1234
```
