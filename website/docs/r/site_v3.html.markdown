---
subcategory: "Cloud WAF - Site Management"
layout: "incapsula"
page_title: "incapsula_site_v3"
description: |- 
  Provides a Incapsula Site resource.
---

# incapsula_site_v3

Provides an Incapsula V3 site resource. A V3 site resource is the core resource that is required by all other resources. incapsula_site_v3 is a newer version of incapsula_site. Site can be managed by incapsula_site_v3 or incapsula_site, but not both simultaneously.
Full site onboarding example with TF can be found [here](https://docs.imperva.com/bundle/cloud-application-security/page/website-certificate-terraform.htm).

## Example Usage

```hcl
resource "incapsula_site_v3" "example-site-v3" {
  name = "example.com"
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional) The account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.
* `name` - (Required) The site name.
* `type` - (Optional) The website type. Indicates which kind of website is created. The allowed options is CLOUD_WAF for a website onboarded to Imperva Cloud WAF.
* `ref_id` - (Optional) Sets the Reference ID. A free-text field that enables you to add a unique identifier to correlate a website in our service with an object on the customer side.
* `active` - (Optional) Whether the site is active or bypassed by the Imperva network.
## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the site.
* `creation_time` - Creation time of the site.
* `cname` - The CNAME provided by Imperva that is used for pointing your website traffic to the Imperva network.


## Import

Site can be imported using the `account Id`/`id`, e.g.:

```
$ terraform import incapsula_site_v3.example 543/1234
```
