package incapsula

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"api_id": "The API identifier for API operations. You can retrieve this\n" +
			"from the Incapsula management console.",

		"api_key": "The API key for API operations. You can retrieve this\n" +
			"from the Incapsula management console.",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIID:  d.Get("api_id").(string),
		APIKey: d.Get("api_key").(string),
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
