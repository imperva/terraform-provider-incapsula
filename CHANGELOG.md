## 3.25.5 (Jul 15, 2024)

IMPROVEMENTS:
- fix site ssl settings documentation ([#454](https://github.com/imperva/terraform-provider-incapsula/pull/454))


## 3.25.4 (Jul 03, 2024)

IMPROVEMENTS:
- Change outputs for SSL verification records so they always exist and have values when using CNAME verification. ([#439](https://github.com/imperva/terraform-provider-incapsula/pull/439))
- Ato - support site id in import command ([#449](https://github.com/imperva/terraform-provider-incapsula/pull/449))

## 3.25.3 (June 13, 2024)

IMPROVEMENTS:
- Account export capability for ATO ([#430](https://github.com/imperva/terraform-provider-incapsula/pull/430))


## 3.25.2 (June 3, 2024)

IMPROVEMENTS:
- Add optional account id parameter to the site ssl settings resource ([#424](https://github.com/imperva/terraform-provider-incapsula/pull/424))


## 3.25.0 (May 20, 2024)

FEATURES:
* **New Resource:** `incapsula_mtls_imperva_to_origin_certificate`

IMPROVEMENTS:
- SiemLogConfiguration: Adding support for SIGNIFICANT_SCRIPT_DISCOVERY and SIGNIFICANT_DATA_TRANSFER_DISCOVERY ([#420](https://github.com/imperva/terraform-provider-incapsula/pull/420))


## 3.24.0 (May 12, 2024)

FEATURES:
* **New Resource:** `incapsula_certificate_signing_request`

IMPROVEMENTS:
- Change error pages to full-update endpoint ([#414](https://github.com/imperva/terraform-provider-incapsula/pull/414))
- add docs for mtls imperva to origin resource ([#418](https://github.com/imperva/terraform-provider-incapsula/pull/418))


## 3.23.1 (Apr 1, 2024)

- removing deprecated account id from policy dto upon update ([#411](https://github.com/imperva/terraform-provider-incapsula/pull/411))
- fix client_certificate_test ([#412](https://github.com/imperva/terraform-provider-incapsula/pull/412))


## 3.23.0 (Mar 31, 2024)

IMPROVEMENTS:
- Add auth type to custom certificate resource ([#402](https://github.com/imperva/terraform-provider-incapsula/pull/402))
- Extending the site ssl settings resource and adding inbound TLS settings that will allow clients to configure which tls versions and which ciphers to use at the site level.
  And extending the already existing ssl settings resource tests to cover inbound TLS settings. ([#406](https://github.com/imperva/terraform-provider-incapsula/pull/406))
- Use different operation name for the new delivery rules resource ([#407](https://github.com/imperva/terraform-provider-incapsula/pull/407))


## 3.22.1 (Mar 25, 2024)

IMPROVEMENTS:

- Delivery rules configuration improvements ([#398](https://github.com/imperva/terraform-provider-incapsula/pull/398))
- deprecating account id for incapsula_policy use incapsula_account_policy_association instead ([#400](https://github.com/imperva/terraform-provider-incapsula/pull/398))
- siem_log_configuration: fix documentation ([#401](https://github.com/imperva/terraform-provider-incapsula/pull/401))



## 3.22.0 (Mar 10, 2024)

FEATURES:

* **New Resource:** `incapsula_siem_connnection_splunk`
* **New Resource:** `incapsula_simplified_redirect_rules_configuration.go`
* **New Resource:** `incapsula_delivery_rules_configuration.go`


## 3.21.5 (Feb 14, 2024)

- Change account_user URLs (internal change) ([#391](https://github.com/imperva/terraform-provider-incapsula/pull/391))


## 3.21.4 (Feb 4, 2024)

IMPROVEMENTS:

- Add HTTP2 attributes in account and sub-account resources ([#388](https://github.com/imperva/terraform-provider-incapsula/pull/388))

## 3.21.3 (Jan 25, 2024)

- incapsula_site_ssl_settings: Revert last changes ([#386](https://github.com/imperva/terraform-provider-incapsula/pull/386))


## 3.21.2 (Jan 10, 2024)

- Documentation fixes


## 3.21.1 (Jan 8, 2024)

- Documentation fixes
- Incapsula_account: fix bug of roles name ([#378](https://github.com/imperva/terraform-provider-incapsula/pull/378))
- incapsula_mtls_client_to_imperva_ca_certificate: fix parameter name ([#379](https://github.com/imperva/terraform-provider-incapsula/pull/379))
- incapsula_security_rule_exception: fix bug ([#380](https://github.com/imperva/terraform-provider-incapsula/pull/380))


## 3.21.0 (Nov 20, 2023)

IMPROVEMENTS:

* Add ABP identification failed error page to application delivery resource ([#370](https://github.com/imperva/terraform-provider-incapsula/pull/370))


## 3.20.6 (Nov 7, 2023)

* Incapsula_Api_Security_API_Config bug fix - missing base path in update method ([#368](https://github.com/imperva/terraform-provider-incapsula/pull/368))


## 3.20.6 (Oct 29, 2023)

* Vulnerabilities fixes ([#361](https://github.com/imperva/terraform-provider-incapsula/pull/361)([#365](https://github.com/imperva/terraform-provider-incapsula/pull/365))
* Incapsula_site documentation fixes ([#363](https://github.com/imperva/terraform-provider-incapsula/pull/363))
* Incapsula_waf_security_rule documentation fixes ([#362](https://github.com/imperva/terraform-provider-incapsula/pull/362))

## 3.20.5 (Oct 18, 2023)

* Incapsula_application_delivery - fix http2 flag([#359](https://github.com/imperva/terraform-provider-incapsula/pull/359))

## 3.20.4 (Sep 27, 2023)

* Incapsula_notification_policy documentation fixes ([#355](https://github.com/imperva/terraform-provider-incapsula/pull/355))
* Incapsula_abp_websites documentation fixes ([#356](https://github.com/imperva/terraform-provider-incapsula/pull/356q))


## 3.20.3 (Sep 18, 2023)

* Incapsula_account_policy_association and client resources bug fixes ([#351](https://github.com/imperva/terraform-provider-incapsula/pull/351))


## 3.20.2 (Sep 10, 2023)

* incapsula_bots_configuration resource documentation fix ([#345](https://github.com/imperva/terraform-provider-incapsula/pull/345))
* Incapsula_account_policy_association resource documentation fix ([#346](https://github.com/imperva/terraform-provider-incapsula/pull/346))


## 3.20.1 (Aug 28, 2023)

* Incapsula_data_center_configuration resource: fix documentation ([#341](https://github.com/imperva/terraform-provider-incapsula/pull/341))
* Incapsula_account_ssl_settings resource: fix documentation ([#342](https://github.com/imperva/terraform-provider-incapsula/pull/342))


## 3.20.0 (Jul 26, 2023)

FEATURES:

* **New Resource:** `incapsula_ato_site_allowlist`
* **New Resource:** `incapsula_ato_endpoint_mitigation_configuration`


## 3.19.0 (Jul 12, 2023)

FEATURES:

* **New Resource:** `incapsula_abp_websites`

IMPROVEMENTS:

* Adding CSP for SIEM Log configuration producer and the accompanying datasets ([#334](https://github.com/imperva/terraform-provider-incapsula/pull/334)))

## 3.18.3 (Jul 05, 2023)

BUG FIXES:

* Incapsula_account resource changes: allow to edit account_name + map naked_domain_san_for_new_www_sites value on Read ([#325](https://github.com/imperva/terraform-provider-incapsula/pull/325))


## 3.18.2 (Jun 26, 2023)

BUG FIXES:

* Policy resource - fix bug with empty policyDataExceptions array in local resource always shows diff ([#322](https://github.com/imperva/terraform-provider-incapsula/pull/322))


## 3.18.1 (Jun 11, 2023)

BUG FIXES:

* Fix documentation site_ssl_settings resource ([#317](https://github.com/imperva/terraform-provider-incapsula/pull/317))
* Fix import for site_ssl_settings resource ([#318](https://github.com/imperva/terraform-provider-incapsula/pull/318))
* Fix rewrite_existing cannot be set to false in incap_rule resource ([#319](https://github.com/imperva/terraform-provider-incapsula/pull/319))


## 3.18.0 (May 29, 2023)

FEATURES:

* **New Resource:** `incapsula_site_ssl_settings`

## 3.17.0 (May 17, 2023)

IMPROVEMENTS:

* incapsula_account - support managing consent ([#307](https://github.com/imperva/terraform-provider-incapsula/pull/307))
* incapsula_siem_log_configuration - Support ATO and AUDIT_TRAIL  ([#308](https://github.com/imperva/terraform-provider-incapsula/pull/308))

## 3.16.1 (Mar 22, 2023)

IMPROVEMENTS:

* incapsula_application_delivery - support compression_type ([#301](https://github.com/imperva/terraform-provider-incapsula/pull/301))

## 3.15.2 (Feb 26, 2023)

BUG FIXES:

* Fix a bug of '+' character in a user's email ([#292](https://github.com/imperva/terraform-provider-incapsula/pull/292))

## 3.15.1 (Feb 6, 2023)

BUG FIXES:  

* Fix unchangeable attributes bug in account resource ([#284](https://github.com/imperva/terraform-provider-incapsula/pull/284))
* Fix bug in import command of siem_log_configuration and incapsula_siem_connection resources ([#285](https://github.com/imperva/terraform-provider-incapsula/pull/285))


## 3.15.0 (Feb 6, 2023)

FEATURES:

* **New Resource:** `incapsula_waiting_room`


* ## 3.14.0 (Jan 15, 2023)

FEATURES:

* **New Resource:** `incapsula_site_domain_configuration`
* **New Resource:** `incapsula_siem_log_configuration`
* **New Resource:** `incapsula_siem_connection`


## 3.13.0 (Jan 8, 2023)

FEATURES:

* **New Resource:** `incapsula_account_role`
* **New Resource:** `incapsula_account_user`

* **New DataSource:** `incapsula_account_permissions`


## 3.12.0 (Dec 11, 2022)

IMPROVEMENTS:

* incapsula_incap_rule - Incap rules enable flag ([#259](https://github.com/imperva/terraform-provider-incapsula/pull/259))


## 3.11.0 (Dec 4, 2022)

FEATURES:

* **New Resource:** incapsula_bots_configuration

IMPROVEMENTS:

* incapsula_incap_rule - Support overrideExisting Flag ([#244](https://github.com/imperva/terraform-provider-incapsula/pull/244))
* incapsula_account_policy_association - added available_policy_ids optional argument + move to v3 apis to improve performance ([#250](https://github.com/imperva/terraform-provider-incapsula/pull/250))


BUG FIXES:

* Fix issue #234 - remove omitempty for boolean fields ([#247](https://github.com/imperva/terraform-provider-incapsula/pull/247))


## 3.10.3 (Nov 20, 2022)

BUG FIXES:

* adding current account id support to incapsula_policy_asset_association ([#243](https://github.com/imperva/terraform-provider-incapsula/pull/243))


## 3.10.2 (Oct 31, 2022)

BUG FIXES:

* policy resource fails to read when account_id param is not provided ([#240](https://github.com/imperva/terraform-provider-incapsula/pull/240))

## 3.10.1 (Oct 31, 2022)

BUG FIXES:

* Fix account ssl settings resource documentation ([#238](https://github.com/imperva/terraform-provider-incapsula/pull/238))


## 3.10.0 (Oct 30, 2022)

FEATURES:

* **New Resource:** incapsula_account_ssl_settings

Deprecations: wildcard_san_for_new_sites, naked_domain_san_for_new_www_sites and support_all_tls_versions in account resource are now deprecated, matched arguments in the account SSL settings resource should be used instead

BUG FIXES:

* Adding account status response to the client object. This allows to have the account context on any client request. ([#232](https://github.com/imperva/terraform-provider-incapsula/pull/232))
* Adding account type to the account status response. ([#232](https://github.com/imperva/terraform-provider-incapsula/pull/232))
* Adding current account to the policy actions. This allows a reseller to manage its accounts' policies ([#232](https://github.com/imperva/terraform-provider-incapsula/pull/232))

## 3.9.1 (Oct 20, 2022)

BUG FIXES:

* documentation corrections ([#229](https://github.com/imperva/terraform-provider-incapsula/pull/229))


## 3.9.0 (Oct 20, 2022)

FEATURES:

* **New Resource:** `incapsula_mtls_client_to_imperva_ca_certificate`
* **New Resource:** `incapsula_mtls_client_to_imperva_ca_certificate_site_settings`
* **New Resource:** `incapsula_mtls_client_to_imperva_ca_certificate_site_association`

BUG FIXES:

* fix documentation of api_security_api_config ([#224](https://github.com/imperva/terraform-provider-incapsula/pull/224))


## 3.8.7 (Oct 3, 2022)

BUG FIXES:

* remove future resource from the documentation ([#219](https://github.com/imperva/terraform-provider-incapsula/pull/219))


## 3.8.6 (Oct 2, 2022)

IMPROVEMENTS:

* incapsula_subaccount: Support for setting default data region for subaccounts ([#207](https://github.com/imperva/terraform-provider-incapsula/pull/207))

BUG FIXES:

* incapsula_policy: fixing bug that clears policy account's defaults when updating policy resource. ([#211](https://github.com/imperva/terraform-provider-incapsula/pull/211))
* The parameters `incapsula_site.restricted_cname_reuse` and `invalid_param_name_violation_action` in all `incapsula_api_security` resources should not be used as they are currently not supported (will be in the future) ([#215](https://github.com/imperva/terraform-provider-incapsula/pull/215))
   


## 3.8.5 (Aug 31, 2022)

BUG FIXES:

* Add retries of read operations when fail ([#200](https://github.com/imperva/terraform-provider-incapsula/pull/200))
* incapsula_api_security_site_config: make is_automatic_discovery_api_integration_enabled optional to align with BE API ([#205](https://github.com/imperva/terraform-provider-incapsula/pull/205))


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
