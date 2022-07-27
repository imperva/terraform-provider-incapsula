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
    imperva_certificate_use_wild_card_san_instead_of_fqdn = true
    imperva_certificate_add_naked_domain_san_for_www_sites = false
    imperva_certificate_delegation_allow_cname_validation = true
    imperva_certificate_delegation_allowed_domains_for_cname_validation = ["example.com"]
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Required) Numeric identifier of the account to operate on, Enterprise accounts can use the current_account from the data_source_account.
* `use_wild_card_san_instead_of_fqdn` - (Optional) Newly created sites will use wildcard SAN instead of FQDN SAN on the Imperva-generated certificate, default true.
* `add_naked_domain_san_for_www_sites` - (Optional) Newly created WWW sites will have also naked domain SAN on the Imperva-generated certificate, default true.
* `allow_cname_validation` - (Optional) Delegate Imperva the ability to automatically prove ownership of SANs under the domains in the allowed_domains_for_cname_validation list, default false.
* `allowed_domains_for_cname_validation` - (Optional) The list of domains that Imperva can automatically prove ownership to the CA on behalf of the customer.

## Attributes Reference

The following attributes are exported:

* `id` - The resource id
* `value_for_cname_validation` - The CNAME record value to allow CA validation delegation 

## Import

Account SSL Settings Configuration can be imported using the account_id (id), e.g.:

```
$ terraform import incapsula_account_ssl_settings.example-terraform-account-ssl-settings-config 123

```
