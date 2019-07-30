package incapsula

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCacheResponseHeaders() *schema.Resource {
	return &schema.Resource{
		Create: resourceCacheResponseHeadersAdd,
		Read:   resourceCacheResponseHeadersRead,
		Update: resourceCacheResponseHeadersUpdate,
		Delete: resourceCacheResponseHeadersDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"site_id": &schema.Schema{
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"cache_headers": &schema.Schema{
				Description: "List of header names to be cached.",
				Type:        schema.TypeString, // TypeList after PR Fix/string to array #13
				//				Elem:        &schema.Schema{Type: schema.TypeString}, PR Fix/string to array #13
				Optional: true,
			},
			"cache_all_headers": &schema.Schema{
				Description: "Cache all response headers. Pass 'true' or 'false' in the value parameter. Cannot be selected together with cache_headers. Default:false",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCacheResponseHeadersAdd(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteID := d.Get("site_id")
	cacheHeaders := d.Get("cache_headers")

	log.Printf("[INFO] Creating Incapsula cache response headers: %s, %s\n", siteID, cacheHeaders)

	_, err := client.ConfigureAdvanceCache(
		d.Get("site_id").(int),
		d.Get("cache_headers").(string),
		//		strings.Join(convertStringArr(d.Get("cache_headers").([]interface{})), ","), PR Fix/string to array #13
		d.Get("cache_all_headers").(string),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula cache response headers: %s, %s\n", siteID, err)
		return err
	}

	d.SetId(fmt.Sprint(siteID))

	log.Printf("[INFO] Created Incapsula cache response headers: %s, %s\n", siteID, cacheHeaders)

	return nil
}

func resourceCacheResponseHeadersRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteID := d.Get("site_id")

	log.Printf("[INFO] Reading Incapsula cache response for site id: %s\n", siteID)

	listCacheHeaderResponse, err := client.SiteStatus("sec-rule-read", d.Get("site_id").(int))

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula cache response for id: %s, %s\n", siteID, err)
		return err
	}

	// now loop through values in status response
	for _, entry := range listCacheHeaderResponse.PerformanceConfiguration.CacheHeaders {
		cacheHeaders := d.Get("cache_headers").(string)

		// cache headers could be multiple values separated with comma
		// handle this by splitting the values and interate through the values
		s := strings.Split(cacheHeaders, ",")

		for _, v := range s {
			if v == entry {
				break
			}
		}
	}

	return nil
}

func resourceCacheResponseHeadersUpdate(d *schema.ResourceData, m interface{}) error {

	// no modify handler, we will use the same call virtually just updating whatever values in configuration
	return resourceCacheResponseHeadersAdd(d, m)
}
func resourceCacheResponseHeadersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteID := d.Get("site_id")

	log.Printf("[INFO] Resetting Incapsula cache response headers to default: %s\n", siteID)

	// there is no delete handler, instead we send empty cache headers keys to the same api
	_, err := client.ConfigureAdvanceCache(
		d.Get("site_id").(int),
		"",
		"",
	)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula site for domain: %s, %s\n", siteID, err)
		return err
	}

	log.Printf("[INFO] Reset Incapsula cache response headers to default: %s\n", siteID)

	return nil
}
