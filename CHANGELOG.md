## 2.5.0 (Unreleased)

* Add support for performance settings in the `site` resource

## 2.4.0 (Unreleased)

* Add support for site masking settings
* Add support for specifying the log level on a site
* Re-factor internal site resource (lots of copy/pasta in create and update)
* Fix an issue with computed and optional attributes for `data_storage_region`
* Configure `site` resource values during an update for: `active`, `acceleration_level`, `seal_location`, `domain_redirect_to_full`, `remove_ssl`, `ignore_ssl`

## 2.3.0 (Unreleased)

* Add support for setting the data storage region on a site
* Remove the deprecated setting for `is_standby` on the `data_center` resource (`is_enabled` replaces this functionality); should resolve flapping integration tests + potential production issues
* Properly configuring `site` resource values for: `active`, `acceleration_level`, `seal_location`, `domain_redirect_to_full`, `remove_ssl`, `ignore_ssl`
* Added `domain_verification` as an exported variable for the `site` resource

## 2.2.0 (Unreleased)

* Add support for cache rules
* Improve documentation
* Fix Incap Rule example bugs

## 2.1.0 (Unreleased)

Add checks for resource destruction during reads. The following resources have been updated:

* ACL Security Rule
* Certificate
* Data Center
* Data Center Server
* Incap Rule
* Security Rule Exception
* Site
* WAF Security Rule

## 2.0.0 (Unreleased)

As we near certification readiness, we've made lots of changes to the provider. Backwards incompatible changes have been made to the Incap Rule resources. Please review the documentation. Changes below:

* All acceptance and unit tests now pass. There was a race condition issue with dependencies - see this Hashicorp issues: https://github.com/hashicorp/terraform/issues/23169 and https://github.com/hashicorp/terraform/issues/23635
* Migrate `resource_incap_rule` to use APIv2. See updated documentation and example files for the latest resource spec.
* Add fixes for data center and data center server result codes (oscillation between strings and ints)
* Fix importing of various resources: data center, data center server, Incap Rule, etc.
* Fix ceriticate argument requirements (thanks @areifert)
* Added GitHub workflow integration for side builds prior to certification (thanks @pklime2)
* Upgrade to Terraform v0.12
* Migrate to standalone Terraform SDK
* Started to improve consistency of error log messages (Site IDs, Rule IDs, etc.) 

## 1.0.0 (Unreleased)

Initial release of the Incapsula Terraform Provider.