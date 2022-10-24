---
layout: "incapsula"
page_title: "Incapsula: incap-account-ssl-settings"
sidebar_current: "docs-incapsula-resource-account-ssl-settings"
description: |- Provides an Incapsula Account SSL Settings resource.
---

# incapsula_account_ssl_settings

Provides an Incapsula Account SSL Settings resource.

## Example Usage

```hcl
resource "incapsula_account_ssl_settings" "ssl-52546413" {
    account_id = 123
    allow_support_old_tls_versions = false
    enable_hsts_for_new_sites = true
    use_wild_card_san_instead_of_fqdn = true
    add_naked_domain_san_for_www_sites = false
    allow_cname_validation = true
    allowed_domain_for_cname_validation = {
       name = "example.com"
    }
    allowed_domain_for_cname_validation = {
       name = "example2.com"
    }
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to operate on, Enterprise accounts can use the current_account from the data_source_account.
* `allow_support_old_tls_versions` - (Optional) When true, sites under the account or sub-accounts can allow support of old TLS versions traffic. This can be configured only on the parent account level, default false.
* `enable_hsts_for_new_sites` - (Optional) When true, enables HSTS support for newly created websites, default true.
* `use_wild_card_san_instead_of_fqdn` - (Optional) Newly created sites will use wildcard SAN instead of FQDN SAN on the Imperva-generated certificate, default true.
* `add_naked_domain_san_for_www_sites` - (Optional) Newly created WWW sites will have also naked domain SAN on the Imperva-generated certificate, default true.
* `allow_cname_validation` - (Optional) Delegate Imperva the ability to automatically prove ownership of SANs under the domains in the allowed_domains_for_cname_validation list, default false.

Optional `allowed_domain_for_cname_validation` sub resources can be defined.
The following allowed_domain_for_cname_validation arguments are supported:  
* `name` - (Required) The domain name.

## Attributes Reference

The following attributes are exported:  

* `id` - The resource id.
* `value_for_cname_validation` - The CNAME record value to allow CA validation delegation .

The following attributes are exported from allowed_domains_for_cname_validation sub resource:  

* `id` - The domain id.
* `status` - The domain status.
* `cname_record_value` - The CNAME record value to use to configure this domain for delegation.
* `cname_record_host` - The CNAME record host to use.
* `creation_date` - The domain creation date.
* `status_since` - The domain status since date.
* `last_status_check` - The domain last status check date.

## Import

Account SSL Settings Configuration can be imported using the account_id (id), e.g.:

```
$ terraform import incapsula_account_ssl_settings.example-terraform-account-ssl-settings-config 123

```
