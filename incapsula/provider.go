package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var baseURL string
var baseURLRev2 string
var baseURLAPI string
var descriptions map[string]string

func init() {
	baseURL = "https://my.incapsula.com/api/prov/v1"
	baseURLRev2 = "https://my.imperva.com/api/prov/v2"
	baseURLAPI = "https://api.imperva.com"

	descriptions = map[string]string{
		"api_id": "The API identifier for API operations. You can retrieve this\n" +
			"from the Incapsula management console. Can be set via INCAPSULA_API_ID " +
			"environment variable.",

		"api_key": "The API key for API operations. You can retrieve this\n" +
			"from the Incapsula management console. Can be set via INCAPSULA_API_KEY " +
			"environment variable.",

		"base_url": "The base URL for API operations. Used for provider development.",

		"base_url_rev_2": "The base URL (revision 2) for API operations. Used for provider development.",

		"base_url_api": "The base URL (same as v2 but with different subdomain) for API operations. Used for provider development.",
	}
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		APIID:       d.Get("api_id").(string),
		APIKey:      d.Get("api_key").(string),
		BaseURL:     d.Get("base_url").(string),
		BaseURLRev2: d.Get("base_url_rev_2").(string),
		BaseURLAPI:  d.Get("base_url_api").(string),
	}

	return config.Client()
}

// Provider returns a *schema.Provider.
func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_API_ID", ""),
				Description: descriptions["api_id"],
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_API_KEY", ""),
				Description: descriptions["api_key"],
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_BASE_URL", baseURL),
				Description: descriptions["base_url"],
			},
			"base_url_rev_2": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_BASE_URL_REV_2", baseURLRev2),
				Description: descriptions["base_url_rev_2"],
			},
			"base_url_api": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_BASE_URL_API", baseURLAPI),
				Description: descriptions["base_url_api"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"incapsula_role_abilities": dataSourceRoleAbilities(),
			"incapsula_data_center":    dataSourceDataCenter(),
			"incapsula_account_data":   dataSourceAccount(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"incapsula_cache_rule":                   resourceCacheRule(),
			"incapsula_custom_certificate":           resourceCertificate(),
			"incapsula_custom_hsm_certificate":       resourceCustomCertificateHsm(),
			"incapsula_data_center":                  resourceDataCenter(),
			"incapsula_data_center_server":           resourceDataCenterServer(),
			"incapsula_incap_rule":                   resourceIncapRule(),
			"incapsula_origin_pop":                   resourceOriginPOP(),
			"incapsula_policy":                       resourcePolicy(),
			"incapsula_account_policy_association":   resourceAccountPolicyAssociation(),
			"incapsula_policy_asset_association":     resourcePolicyAssetAssociation(),
			"incapsula_security_rule_exception":      resourceSecurityRuleException(),
			"incapsula_site":                         resourceSite(),
			"incapsula_waf_security_rule":            resourceWAFSecurityRule(),
			"incapsula_account":                      resourceAccount(),
			"incapsula_subaccount":                   resourceSubAccount(),
			"incapsula_waf_log_setup":                resourceWAFLogSetup(),
			"incapsula_txt_record":                   resourceTXTRecord(),
			"incapsula_data_centers_configuration":   resourceDataCentersConfiguration(),
			"incapsula_api_security_site_config":     resourceApiSecuritySiteConfig(),
			"incapsula_api_security_api_config":      resourceApiSecurityApiConfig(),
			"incapsula_api_security_endpoint_config": resourceApiSecurityEndpointConfig(),
			"incapsula_notification_center_policy":   resourceNotificationCenterPolicy(),
			"incapsula_csp_site_configuration":       resourceCSPSiteConfiguration(),
			"incapsula_csp_site_domain":              resourceCSPSiteDomain(),
			"incapsula_application_delivery":         resourceApplicationDelivery(),
			"incapsula_site_monitoring":              resourceSiteMonitoring(),
			"incapsula_account_ssl_settings":         resourceAccountSSLSettings(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}
