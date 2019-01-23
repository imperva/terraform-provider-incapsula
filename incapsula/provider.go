package incapsula

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var baseURL string
var descriptions map[string]string

func init() {
	// Storing this in the provider rather than making it configurable
	// This endpoint is unlikely to change in the near future and will be augmented for testing
	baseURL = "https://my.incapsula.com/api/prov/v1"

	descriptions = map[string]string{
		"api_id": "The API identifier for API operations. You can retrieve this\n" +
			"from the Incapsula management console.",

		"api_key": "The API key for API operations. You can retrieve this\n" +
			"from the Incapsula management console.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiID := d.Get("api_id").(string)
	apiKey := d.Get("api_key").(string)

	config := Config{
		APIID:   apiID,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	return config.Client()
}

// Provider returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["api_id"],
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: descriptions["api_key"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"incapsula_site": resourceSite(),
		},

		ConfigureFunc: providerConfigure,
	}
}
