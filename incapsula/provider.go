package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var baseURL string
var incapRuleBaseURL string
var descriptions map[string]string

func init() {
	// Storing this in the provider rather than making it configurable
	// This endpoint is unlikely to change in the near future
	baseURL = "https://my.incapsula.com/api/prov/v1"

	// Efforts on APIv2 are underway and IncapRules are supported
	// The other endpoints will eventually move over but we'll need this for now
	incapRuleBaseURL = "https://my.imperva.com/api/prov/v2"

	descriptions = map[string]string{
		"api_id": "The API identifier for API operations. You can retrieve this\n" +
			"from the Incapsula management console. Can be set via INCAPSULA_API_ID " +
			"environment variable.",

		"api_key": "The API key for API operations. You can retrieve this\n" +
			"from the Incapsula management console. Can be set via INCAPSULA_API_KEY " +
			"environment variable.",
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	apiID := d.Get("api_id").(string)
	apiKey := d.Get("api_key").(string)

	config := Config{
		APIID:            apiID,
		APIKey:           apiKey,
		BaseURL:          baseURL,
		IncapRuleBaseURL: incapRuleBaseURL,
	}

	return config.Client()
}

// Provider returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
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
		},

		ResourcesMap: map[string]*schema.Resource{
			"incapsula_site":                    resourceSite(),
			"incapsula_acl_security_rule":       resourceACLSecurityRule(),
			"incapsula_waf_security_rule":       resourceWAFSecurityRule(),
			"incapsula_security_rule_exception": resourceSecurityRuleException(),
			"incapsula_incap_rule":              resourceIncapRule(),
			"incapsula_data_center":             resourceDataCenter(),
			"incapsula_data_center_server":      resourceDataCenterServer(),
			"incapsula_custom_certificate":      resourceCertificate(),
		},

		ConfigureFunc: configureProvider,
	}
}
