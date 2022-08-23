
## 3.8.4 (Aug 17, 2022)

BUG FIXES:

* incapsula_subaccount: fix 'read' method to use another API to the backend ([#197](https://github.com/imperva/terraform-provider-incapsula/pull/197))

## 3.8.3 (Aug 11, 2022)

BUG FIXES:

* remove future resource from the documentation ([#195](https://github.com/imperva/terraform-provider-incapsula/pull/195))

## 3.8.2 (Aug 3, 2022)

BUG FIXES:

* incapsula_application_delivery: return ports to default upon deleting the resource ([#189](https://github.com/imperva/terraform-provider-incapsula/pull/189))

## 3.8.1 (Jul 5, 2022)

IMPROVEMENTS:

* incapsula_site: Add cname option to domain validation options

BUG FIXES:

* incapsula_account_policy_association: change default_waf_policy_id to be optional - for customers who have not migrated yet to waf policy ([#185](https://github.com/imperva/terraform-provider-incapsula/issues/185))
* incapsula_site: change default values from string to bool ([#186](https://github.com/imperva/terraform-provider-incapsula/issues/186))


## 3.8.0 (Jul 4, 2022)

FEATURES:

* **New Resource:** `incapsula_account-policy-association`
* **New Resource:** `incapsula_application_delivery`
* **New Resource:** `incapsula_site_monitoring`

* **New DataSource:** `account-data`


## 3.7.0 (June 30, 2022)

FEATURES:

* **New Resource:** `incapsula_waf_log_setup`

## 3.6.0 (June 8, 2022)

FEATURES:

* **New Resource:** `incapsula_custom_hsm_certificate`

## 3.5.2 (May 16, 2022)

IMPROVEMENTS:

* Add deprecation message to already deprecated resources (old data_center resources)

BUG FIXES:

* incapsula_site: formatting parameters with %t fails if the values are strings, not bool ([#158](https://github.com/imperva/terraform-provider-incapsula/issues/158))
* incapsula_site: add retries when configuring site after creating it - to allow the site creation to fully finish ([#165](https://github.com/imperva/terraform-provider-incapsula/issues/165))
* incapsula_notification_center_policy: Fix redundant slash in path issue ([#162](https://github.com/imperva/terraform-provider-incapsula/issues/162))
* incapsula_origin_pop: avoid crashing when upgrading from version 2* to 3* without changing the resource format in the state file([#167](https://github.com/imperva/terraform-provider-incapsula/issues/167))

## 3.5.1 (Released)

BUG FIXES:

* Fix a bug where naked_domain_san and wildcard_san attributes on site resource weren't handled by 'modify' method

## 3.5.0 (Released)

FEATURES:

* **New Resource:** `CSP_Site_configuration`, `CSP_Site_domain`

IMPROVEMENTS:

* Add operation type to HTTP client calls
* Fix acceptance test for Custom Certificate resource


## 3.4.0 (Released) 

* Add support for notification center

## 3.3.4 (Released) 

* Fix bug in Custom Certificate resource

## 3.3.3 (Released)

* Fix pagination bug in sub-account resource

## 3.3.2 (Released)

* Edit business logic, add acceptance test for incapsula_txt_record resource

## 3.3.1 (Released)

* No Changes were detected

## 3.3.0 (Released)

* SubAccount resource addition (incapsula_subaccount)

## 3.2.2 (Released)

* Support 'force-risky-operation' header for cache settings

## 3.2.1 (Released)

* Fix custom_certificate resource: remove 'ForceNew' and unecessary base64 encoding

## 3.1.1 (Released)

* Fix `perf_response_cache_404_time` in the `site` resource. Validate if its value is divisible by 60
* Make `original_data_center_id` deprecated in the `site` resource
* Add "ForceNew" for `method` and `path` arguments in the `api_security_endpoint_config` resource

## 3.1.0 (Released)

* Support API-Security with new resources: api_security_api_config, api_security_endpoint_config,
  api_security_site_config

## 3.0.1 (Released)

* Fix TTL attribute of incap_rule resource: enable zero as value
* New attribute on site resource: strict_cname_reuse

## 3.0.0 (Released)

* New resource: data_centers_configuration
* Remove unnecessary 'ForceNew' from attributes of site resource

## 2.9.0 (Released)

* Fix `acceleration_level` in the `site` resource
* Fix support for `continents` in the `security_rule_exception` resource
* Fix the default value for `seal_location` in the `site` resource
* Fix documentation for the `security_rule_exception` resource
* Add TF provider version to HTTP client calls

## 2.8.0 (Released)
* Add new `incap_rule` properties and `RULE_ACTION_FORWARD_TO_PORT` action
* Fix `security_rule_exception` import
* Fix redundant `data_center_server` when `enabled=false`
* Fix `origin_pop` import
* Add `policy_asset_association` import
* Fix `site` resource to use the `logs_account_id` for various methods (read/update)
* Fix `site` resource to read/update `seal_location`

## 2.7.5 (Released)
* Fixed `naked_domain_san` and `wildcard_san` on the `site` resource.
* Adding edit `server_address` ability to `data_center` resource.
* Updating several resources parameters to include `ForceNew`.

## 2.7.4 (Released)

* Add retry logic to `site` and `data_center` resources
* Set the `dns_record_name` in the `site` resource
* Add the `naked_domain_san` in the `site` resource
* Add the `wildcard_san` in the `site` resource
* Update type `PerformanceSettings` struct to omit when empty except the mode_level
* Add `txt_record_value_*` resource
* Update the `site` resource value `site_ip` to be computed
* Fix the `site_ip` to store in state file
* Fix `data_storage_region` issue in `account` resource, defaults to `US`

## 2.7.3 (Released)

* Set the `log_level` in the `site` resource to be optional
* Fix `omitempty` issues in `client_performance`

## 2.7.2 (Released)

* Fix documentation bug for `account` resource

## 2.7.1 (Released)

* Fix documentation bug for origin POP codes

## 2.7.0 (Released)

* Add support for setting the data center's origin POP with a new synthetic resource
* Fix a bug on setting `is_content` on the `data_center` resource
* Remove old `acl_security_rule` resource as it has been migrated to the `policy` resource
* Update all documentation to include all import operations

## 2.6.2 (Released)

* Provider fixes `parameters` on SQL Injection Security Exception resource
* Merged support for account creation (used by re-sellers) w/ fixes

## 2.6.1 (Released)

* Provider has landed in the Terraform Registry
* Cleared out old GitHub workflows

## 2.6.0 (Released)

* Add support for policy management

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

## 2.1.0 (Released)

Add checks for resource destruction during reads. The following resources have been updated:

* ACL Security Rule
* Certificate
* Data Center
* Data Center Server
* Incap Rule
* Security Rule Exception
* Site
* WAF Security Rule

## 2.0.0 (Released)

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

## 1.0.0 (Released)

Initial release of the Incapsula Terraform Provider.