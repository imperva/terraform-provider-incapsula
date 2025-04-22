package incapsula

const VerifyAccount = "verify_account"

const CreateSite = "create_site"
const ReadSite = "read_site"
const UpdateSite = "update_site"
const DeleteSite = "delete_site"

const CreatePolicy = "create_policy"
const ReadPolicy = "read_policy"
const UpdatePolicy = "update_policy"
const DeletePolicy = "delete_policy"

const ReadPoliciesAll = "read_policies_all"

const ReadPolicyAccountAssociatiation = "read_policy_account_association"
const UpdatePolicyAccountAssociatiation = "update_policy_account_association"

const CreatePolicyAssetAssociation = "create_policy_asset_association"
const ReadPolicyAssetAssociation = "read_policy_asset_association"
const DeletePolicyAssetAssociation = "delete_policy_asset_association"

const CreateAccount = "create_account"
const ReadAccount = "read_account"
const UpdateAccount = "update_account"
const DeleteAccount = "delete_account"

const CreateSubAccount = "create_sub_account"
const ReadSubAccount = "read_sub_account"
const DeleteSubAccount = "delete_sub_account"

const CreateWAFLogSetup = "create_waf_log_setup"
const DeleteWAFLogSetup = "delete_waf_log_setup"
const ActivateWAFLogSetup = "activate_waf_log_setup"
const UpdateStatusWAFLogSetup = "update_status_waf_log_setup"

const ReadAccountDataStorageRegion = "read_account_data_storage_region"
const UpdateAccountDataStorageRegion = "update_account_data_storage_region"

const ReadSiteMasking = "read_site_masking"
const UpdateSiteMasking = "update_site_masking"

const UpdateLogLevel = "update_log_level"

const ReadSitePerformance = "read_site_performance"
const UpdateSitePerformance = "update_site_performance"

const CreateApiSecApiConfig = "create_api_sec_api_config"
const ReadApiSecApiConfig = "read_api_sec_api_config"
const UpdateApiSecApiConfig = "update_api_sec_api_config"
const DeleteApiSecApiConfig = "delete_api_sec_api_config"

const ReadApiSecEndpointConfig = "read_api_sec_endpoint_config"
const UpdateApiSecEndpointConfig = "update_api_sec_endpoint_config"

const ReadApiSecSiteConfig = "read_api_sec_site_config"
const UpdateApiSecSiteConfig = "update_api_sec_site_config"

const CreateCacheRule = "create_cache_rule"
const ReadCacheRule = "read_cache_rule"
const UpdateCacheRule = "update_cache_rule"
const DeleteCacheRule = "delete_cache_rule"

const CreateIncapRule = "create_incap_rule"
const ReadIncapRule = "read_incap_rule"
const UpdateIncapRule = "update_incap_rule"
const DeleteIncapRule = "delete_incap_rule"

const UpdateSecurityRule = "update_security_rule"

const CreateSecurityRuleException = "create_security_rule_exception"
const ReadSecurityRuleException = "read_security_rule_exception"
const UpdateSecurityRuleException = "update_security_rule_exception"
const DeleteSecurityRuleException = "delete_security_rule_exception"

const CreateCertificateSigningRequest = "create_certificate_signing_request"

const CreateCustomCertificate = "create_custom_certificate"
const ReadCustomCertificate = "read_custom_certificate"
const UpdateCustomCertificate = "update_custom_certificate"
const DeleteCustomCertificate = "delete_custom_certificate"

const CreateHSMCustomCertificate = "create_hsm_custom_certificate"
const ReadHSMCustomCertificate = "read_hsm_custom_certificate"
const DeleteHsmCustomCertificate = "delete_hsm_custom_certificate"

const CreateDataCenter = "create_data_center"
const ReadDataCenter = "read_data_center"
const UpdateDataCenter = "update_data_center"
const DeleteDataCenter = "delete_data_center"

const CreateDataCenterConfiguration = "create_data_configuration"
const ReadDataCenterConfiguration = "read_data_configuration"

const CreateBotConfiguration = "create_bot_configuration"
const ReadBotConfiguration = "read_bot_configuration"
const ReadClientApplications = "read_client_applications"

const ReadDataStorageRegion = "read_data_storage_region"
const UpdateDataStorageRegion = "update_data_storage_region"

const UpdateOriginPop = "update_origin_pop"

const CreateDataCenterServer = "create_data_center_server"
const UpdateDataCenterServer = "update_data_center_server"
const DeleteDataCenterServer = "delete_data_center_server"

const CreateTxtRecord = "create_txt_record"
const ReadTxtRecord = "read_txt_record"
const UpdateTxtRecord = "update_txt_record"
const DeleteTxtRecord = "delete_txt_record"

const ReadCspSiteConfiguration = "read_csp_site_configuration"
const UpdateCspSiteConfiguration = "update_csp_site_configuration"

const CreateCspSiteDomain = "create_csp_site_domain"
const ReadCspSiteDomain = "read_csp_site_domain"
const UpdateCspSiteDomain = "update_csp_site_domain"
const DeleteCspSiteDomain = "delete_csp_site_domain"

const CreateATOSiteAllowlistOperation = "create_ato_site_allowlist"
const ReadATOSiteAllowlistOperation = "read_ato_site_allowlist"
const UpdateATOSiteAllowlistOperation = "read_ato_site_allowlist"
const DeleteATOSiteAllowlistOperation = "read_ato_site_allowlist"

const CreateATOSiteMitigationConfigurationOperation = "create_ato_site_mitigation_configuration"
const ReadATOSiteMitigationConfigurationOperation = "read_ato_site_mitigation_configuration"
const UpdateATOSiteMitigationConfigurationOperation = "read_ato_site_mitigation_configuration"
const DeleteATOSiteMitigationConfigurationOperation = "read_ato_site_mitigation_configuration"

const CreateNotificationCenterPolicy = "create_notification_center_policy"
const ReadNotificationCenterPolicy = "read_notification_center_policy"
const UpdateNotificationCenterPolicy = "update_notification_center_policy"
const DeleteNotificationCenterPolicy = "delete_notification_center_policy"

const UpdateAccountSSLSettings = "update_account_ssl_settings"
const GetAccountSSLSettings = "get_account_ssl_settings"
const DeleteAccountSSLSettings = "delete_account_ssl_settings"

const CreateMtlsImpervaToOriginCertifiate = "create_mtls_imperva_to_origin_certificate"
const ReadMtlsImpervaToOriginCertifiate = "read_mtls_imperva_to_origin_certificate"
const UpdateMtlsImpervaToOriginCertifiate = "update_mtls_imperva_to_origin_certificate"
const DeleteMtlsImpervaToOriginCertifiate = "delete_mtls_imperva_to_origin_certificate"

const CreateSiteMtlsImpervaToOriginCertifiateAssociation = "create_site_mtls_imperva_to_origin_certificate_association"
const ReadSiteMtlsImpervaToOriginCertifiateAssociation = "read_site_mtls_imperva_to_origin_certificate_association"
const DeleteSiteMtlsImpervaToOriginCertifiateAssociation = "delete_site_mtls_imperva_to_origin_certificate_association"

const CreateMtlsClientToImpervaCertifiate = "create_mtls_client_to_imperva_certificate"
const ReadMtlsClientToImpervaCertifiate = "read_mtls_client_to_imperva_certificate"
const UpdateMtlsClientToImpervaCertifiate = "update_mtls_client_to_imperva_certificate"
const DeleteMtlsClientToImpervaCertifiate = "delete_mtls_client_to_imperva_certificate"

const CreateMtlsClientToImpervaCertifiateSiteAssociation = "create_mtls_client_to_imperva_certificate_site_accociation"
const ReadMtlsClientToImpervaCertifiateSiteAssociation = "read_mtls_client_to_imperva_certificate_site_accociation"
const DeleteMtlsClientToImpervaCertifiateSiteAssociation = "delete_mtls_client_to_imperva_certificate_site_accociation"

const CreateSiteTlsSettings = "create_site_tls_settings"
const ReadSiteTlsSettings = "read_site_tls_settings"
const UpdateSiteSSLSettings = "update_site_ssl_settings"
const ReadSiteSSLSettings = "read_site_ssl_settings"

const CreateAccountRole = "create_account_role"
const ReadAccountRole = "read_account_role"
const UpdateAccountRole = "update_account_role"
const DeleteAccountRole = "delete_account_role"

const ReadAccountAbilities = "read_account_abilities"
const ReadAccountRoles = "read_account_roles"

const CreateAccountUser = "create_account_user"
const CreateSubAccountUser = "create_sub_account_user"
const ReadAccountUser = "read_account_user"
const UpdateAccountUser = "update_account_user"
const DeleteAccountUser = "delete_account_user"

const UpdateDomain = "update_domain"
const CreateDomain = "create_domain"
const DeleteDomain = "delete_domain"

const ReadDomain = "read_domain"
const ReadDomainExtraDetails = "read_domain_extra_details"

const CreateSiemConnection = "create_siem_connection"
const ReadSiemConnection = "read_siem_connection"
const UpdateSiemConnection = "update_siem_connection"
const DeleteSiemConnection = "delete_siem_connection"

const CreateSiemLogConfiguration = "create_siem_log_configuration"
const ReadSiemLogConfiguration = "read_siem_log_configuration"
const UpdateSiemLogConfiguration = "update_siem_log_configuration"
const DeleteSiemLogConfiguration = "delete_siem_log_configuration"

const CreateWaitingRoom = "create_waiting_room"
const ReadWaitingRoom = "read_waiting_room"
const UpdateWaitingRoom = "update_waiting_room"
const DeleteWaitingRoom = "delete_waiting_room"

const CreateAbpWebsites = "create_abp_websites"
const ReadAbpWebsites = "read_abp_websites"
const UpdateAbpWebsites = "update_abp_websites"
const DeleteAbpWebsites = "delete_abp_websites"

const RequestSiteCert = "request_site_cert"

const AddV3Site = "add_v3_site"
const UpdateV3Site = "update_v3_site"

const ReadDeliveryRuleConfiguration = "read_delivery_rules_configuration"
const UpdateDeliveryRuleConfiguration = "update_delivery_rules_configuration"

const ReadFastRenewalConfiguration = "read_fast_renewal_configuration"
const CreateFastRenewalConfiguration = "create_fast_renewal_configuration"
const DeleteFastRenewalConfiguration = "delete_fast_renewal_configuration"
