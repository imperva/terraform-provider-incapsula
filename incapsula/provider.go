package incapsula

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TfResource struct {
	Type string
	Id   string
}

var baseURL string
var baseURLRev2 string
var baseURLRev3 string
var baseURLAPI string
var descriptions map[string]string

func init() {
	baseURL = "https://my.incapsula.com/api/prov/v1"
	baseURLRev2 = "https://my.imperva.com/api/prov/v2"
	baseURLRev3 = "https://my.imperva.com/api/prov/v3"
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

		"base_url_rev_3": "The base URL (revision 3) for API operations. Used for provider development.",

		"base_url_api": "The base URL (same as v2 but with different subdomain) for API operations. Used for provider development.",
	}
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, error) {
	config := Config{
		APIID:       d.Get("api_id").(string),
		APIKey:      d.Get("api_key").(string),
		BaseURL:     d.Get("base_url").(string),
		BaseURLRev2: d.Get("base_url_rev_2").(string),
		BaseURLRev3: d.Get("base_url_rev_3").(string),
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
			"base_url_rev_3": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_BASE_URL_REV_3", baseURLRev3),
				Description: descriptions["base_url_rev_3"],
			},
			"base_url_api": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("INCAPSULA_BASE_URL_API", baseURLAPI),
				Description: descriptions["base_url_api"],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"incapsula_role_abilities":      dataSourceRoleAbilities(),
			"incapsula_data_center":         dataSourceDataCenter(),
			"incapsula_account_data":        dataSourceAccount(),
			"incapsula_client_apps_data":    dataSourceClientApps(),
			"incapsula_account_permissions": dataSourceAccountPermissions(),
			"incapsula_account_roles":       dataSourceAccountRoles(),
			"incapsula_ssl_instructions":    dataSourceSSLInstructions(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"incapsula_cache_rule":                                             resourceCacheRule(),
			"incapsula_certificate_signing_request":                            resourceCertificateSigningRequest(),
			"incapsula_custom_certificate":                                     resourceCertificate(),
			"incapsula_custom_hsm_certificate":                                 resourceCustomCertificateHsm(),
			"incapsula_data_center":                                            resourceDataCenter(),
			"incapsula_data_center_server":                                     resourceDataCenterServer(),
			"incapsula_incap_rule":                                             resourceIncapRule(),
			"incapsula_origin_pop":                                             resourceOriginPOP(),
			"incapsula_policy":                                                 resourcePolicy(),
			"incapsula_account_policy_association":                             resourceAccountPolicyAssociation(),
			"incapsula_policy_asset_association":                               resourcePolicyAssetAssociation(),
			"incapsula_security_rule_exception":                                resourceSecurityRuleException(),
			"incapsula_site":                                                   resourceSite(),
			"incapsula_managed_certificate_settings":                           resourceManagedCertificate(),
			"incapsula_site_v3":                                                resourceSiteV3(),
			"incapsula_waf_security_rule":                                      resourceWAFSecurityRule(),
			"incapsula_account":                                                resourceAccount(),
			"incapsula_subaccount":                                             resourceSubAccount(),
			"incapsula_waf_log_setup":                                          resourceWAFLogSetup(),
			"incapsula_txt_record":                                             resourceTXTRecord(),
			"incapsula_data_centers_configuration":                             resourceDataCentersConfiguration(),
			"incapsula_api_security_site_config":                               resourceApiSecuritySiteConfig(),
			"incapsula_api_security_api_config":                                resourceApiSecurityApiConfig(),
			"incapsula_api_security_endpoint_config":                           resourceApiSecurityEndpointConfig(),
			"incapsula_notification_center_policy":                             resourceNotificationCenterPolicy(),
			"incapsula_site_ssl_settings":                                      resourceSiteSSLSettings(),
			"incapsula_site_log_configuration":                                 resourceSiteLogConfiguration(),
			"incapsula_ssl_validation":                                         resourceDomainsValidation(),
			"incapsula_csp_site_configuration":                                 resourceCSPSiteConfiguration(),
			"incapsula_csp_site_domain":                                        resourceCSPSiteDomain(),
			"incapsula_ato_site_allowlist":                                     resourceATOSiteAllowlist(),
			"incapsula_ato_endpoint_mitigation_configuration":                  ATOEndpointMitigationConfiguration(),
			"incapsula_application_delivery":                                   resourceApplicationDelivery(),
			"incapsula_site_monitoring":                                        resourceSiteMonitoring(),
			"incapsula_account_ssl_settings":                                   resourceAccountSSLSettings(),
			"incapsula_mtls_imperva_to_origin_certificate":                     resourceMtlsImpervaToOriginCertificate(),
			"incapsula_mtls_imperva_to_origin_certificate_site_association":    resourceMtlsImpervaToOriginCertificateSiteAssociation(),
			"incapsula_mtls_client_to_imperva_ca_certificate":                  resourceMtlsClientToImpervaCertificate(),
			"incapsula_mtls_client_to_imperva_ca_certificate_site_association": resourceMtlsClientToImpervaCertificateSiteAssociation(),
			"incapsula_mtls_client_to_imperva_ca_certificate_site_settings":    resourceMtlsClientToImpervaCertificateSetings(),
			"incapsula_site_domain_configuration":                              resourceSiteDomainConfiguration(),
			"incapsula_domain":                                                 resourceSiteSingleDomainConfiguration(),
			"incapsula_bots_configuration":                                     resourceBotsConfiguration(),
			"incapsula_account_role":                                           resourceAccountRole(),
			"incapsula_account_user":                                           resourceAccountUser(),
			"incapsula_siem_connection":                                        resourceSiemConnection(),
			"incapsula_siem_splunk_connection":                                 resourceSiemSplunkConnection(),
			"incapsula_siem_sftp_connection":                                   resourceSiemSftpConnection(),
			"incapsula_siem_log_configuration":                                 resourceSiemLogConfiguration(),
			"incapsula_waiting_room":                                           resourceWaitingRoom(),
			"incapsula_abp_websites":                                           resourceAbpWebsites(),
			"incapsula_delivery_rules_configuration":                           resourceDeliveryRulesConfiguration(),
			"incapsula_simplified_redirect_rules_configuration":                resourceSimplifiedRedirectRulesConfiguration(),
			"incapsula_site_cache_configuration":                               resourceSiteCacheConfiguration(),
			"incapsula_short_renewal_cycle":                                    resourceShortRenewalCycle(),
		},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			terraformVersion = "0.11+compatible"
		}
		diags := getLLMSuggestions(d)
		client, _ := providerConfigure(d, terraformVersion)
		return client, diags
	}

	return provider
}

func getLLMSuggestions(d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	folderPath := "terraform_client/nir-terraform-tf-testing/"
	allResourcesFromState := getAllResourcesFromState(folderPath + "terraform.tfstate")
	for _, res := range allResourcesFromState {
		log.Printf("Resource: %s\n", res)
	}

	resources := getAllResourcesTypeAndId(folderPath + "terraform.tfstate")
	for _, res := range resources {
		log.Printf("Resource Type: %s, ID: %s\n", res.Type, res.Id)
	}
	allResourcesFromFiles, _ := getAllResourcesFromTfFiles(folderPath)
	log.Printf("Resource from file: %s\n", allResourcesFromFiles)

	diags = getMissingResources(d, resources, diags)
	diags = getBestPractices(d, allResourcesFromFiles, diags)
	diags = getResourceReplaceSuggestions(d, allResourcesFromFiles, diags)
	diags = getResourceSuggestions(d, allResourcesFromFiles, diags)

	return diags
}

func getResourceSuggestions(d *schema.ResourceData, resources string, diags diag.Diagnostics) diag.Diagnostics {
	question := "you are slim shady, whats your name? out put should be a string only."

	answer, _ := answerWithTools(question, d.Get("api_id").(string), d.Get("api_key").(string))
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "New Resources Suggestion",
		Detail:   answer,
	})
	return diags
}

func getResourceReplaceSuggestions(d *schema.ResourceData, resources string, diags diag.Diagnostics) diag.Diagnostics {
	question := "you are slim shady, whats your name? out put should be a string only."

	answer, _ := answerWithTools(question, d.Get("api_id").(string), d.Get("api_key").(string))
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource Replacement Suggestion",
		Detail:   answer,
	})
	return diags
}

func getBestPractices(d *schema.ResourceData, resources string, diags diag.Diagnostics) diag.Diagnostics {
	question := "you are slim shady, whats your name? out put should be a string only."

	answer, _ := answerWithTools(question, d.Get("api_id").(string), d.Get("api_key").(string))
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Best Practice Suggestion",
		Detail:   answer,
	})
	return diags
}

func getMissingResources(d *schema.ResourceData, resources []TfResource, diags diag.Diagnostics) diag.Diagnostics {
	//question := "Based on the giving resources, which comes in the following structure [{{resource name resource id}}]" +
	//	" fetch all the sites from the backend and compare them with the given sites resources. " +
	//	" check which resources are missing and output the missing resources only" +
	//	" output should be in the following json format: " +
	//	"[{{ \"resource_type\": \"<resource_type>\", \"resource_id\": \"<resource_id>\", \"site name\": \"<site_name>\" }}]" +
	//	" given resources: " + fmt.Sprintf("%v", resources)

	question := "Fetch all the sites from the backend, use pagination 50, or fetch all pages. then output all the sites in the following json format: " +
		"[{{ \"resource_type\": \"<resource_type>\", \"resource_id\": \"<resource_id>\", \"site name\": \"<site_name>\" }}]"

	answer, _ := answerWithTools(question, d.Get("api_id").(string), d.Get("api_key").(string))
	log.Printf("[Info] LLM Missing Resources Answer: %s\n", answer)
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Missing Resource Suggestion",
		Detail:   answer,
	})
	return diags
}

func getAllResourcesTypeAndId(statePath string) []TfResource {
	var resources []TfResource
	file, err := os.Open(statePath)
	if err != nil {
		log.Printf("[Error] Unable to open state file: %v", err)
		return resources
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("[Error] Unable to close state file: %v", err)
		}
	}(file)

	var state struct {
		Resources []struct {
			Type      string `json:"type"`
			Instances []struct {
				Attributes map[string]interface{} `json:"attributes"`
			} `json:"instances"`
		} `json:"resources"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		log.Printf("[Error] Unable to decode state file: %v", err)
		return resources
	}

	for _, resource := range state.Resources {
		for _, instance := range resource.Instances {
			id, ok := instance.Attributes["id"]
			if ok {
				if idStr, isStr := id.(string); isStr {
					resources = append(resources, TfResource{Type: resource.Type, Id: idStr})
				}
			}
		}
	}
	return resources
}

func getAllResourcesFromState(statePath string) []map[string]interface{} {
	var resources []map[string]interface{}
	file, err := os.Open(statePath)
	if err != nil {
		log.Printf("[Error] Unable to open state file: %v", err)
		return resources
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("[Error] Unable to close state file: %v", err)
		}
	}(file)

	var state struct {
		Resources []struct {
			Type      string `json:"type"`
			Instances []struct {
				Attributes map[string]interface{} `json:"attributes"`
			} `json:"instances"`
		} `json:"resources"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&state); err != nil {
		log.Printf("[Error] Unable to decode state file: %v", err)
		return resources
	}

	for _, resource := range state.Resources {
		for _, instance := range resource.Instances {
			resources = append(resources, instance.Attributes)
		}
	}
	return resources
}

func getAllResourcesFromTfFiles(dir string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.tf"))
	if err != nil {
		return "", err
	}
	var content string
	for _, file := range files {
		src, err := ioutil.ReadFile(file)
		if err != nil {
			continue
		}
		content += string(src) + "\n"
	}
	return content, nil
}
