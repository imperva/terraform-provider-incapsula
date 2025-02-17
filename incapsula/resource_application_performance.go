package incapsula

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceApplicationPerformance() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationPerformanceUpdate,
		Read:   resourceApplicationPerformanceRead,
		Update: resourceApplicationPerformanceUpdate,
		Delete: resourceApplicationPerformanceDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID, err := strconv.Atoi(d.Id())
				if err != nil {
					return nil, fmt.Errorf("failed to convert Site Id from import command for Application Performance resource, actual value: %s, expected numeric id", d.Id())
				}

				d.Set("site_id", siteID)
				log.Printf("[DEBUG] Import Application Performance for Site ID %d", siteID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on. ",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"client_comply_no_cache": {
				Description: "Comply with No-Cache and Max-Age directives in client requests. By default, these cache directives are ignored. Resources are dynamically profiled and re-configured to optimize performance.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"client_enable_client_side_caching": {
				Description: "Cache content on client browsers or applications. When not enabled, content is cached only on the Imperva proxies.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"client_send_age_header": {
				Description: "Send Cache-Control: max-age and Age headers.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"key_comply_vary": {
				Description: "Comply with Vary. Cache resources in accordance with the Vary response header.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"key_unite_naked_full_cache": {
				Description: "Use the Same Cache for Full and Naked Domains. For example, use the same cached resource for www.example.com/a and example.com/a.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"mode_https": {
				Description:  "The resources that are cached over HTTPS, the general level applies. Options are `disabled`, `dont_include_html`, `include_html`, and `include_all_resources`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "dont_include_html", "include_html", "include_all_resources"}, false),
			},
			"mode_level": {
				Description:  "Caching level. Options are `disabled`, `standard`, `smart`, and `all_resources`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "smart",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "standard", "smart", "all_resources"}, false),
			},
			"mode_time": {
				Description: "The time, in seconds, that you set for this option determines how often the cache is refreshed. Relevant for the `include_html` and `include_all_resources` levels only.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"response_cache_300x": {
				Description: "When this option is checked Imperva will cache 301, 302, 303, 307, and 308 redirect response headers containing the target URI.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"response_cache_404_enabled": {
				Description: "Whether or not to cache 404 responses.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"response_cache_404_time": {
				Description:  "The time in seconds to cache 404 responses.",
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntDivisibleBy(60),
			},
			"response_cache_empty_responses": {
				Description: "Cache responses that don’t have a message body.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"response_cache_http_10_responses": {
				Description: "Cache HTTP 1.0 type responses that don’t include the Content-Length header or chunking.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"response_cache_response_header_mode": {
				Description:  "The working mode for caching response headers. Options are `all`, `custom` and `disabled`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "custom", "all"}, false),
			},
			"response_cache_response_headers": {
				Description: "An array of strings representing the response headers to be cached when working in `custom` mode. If empty, no response headers are cached.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				DefaultFunc: func() (interface{}, error) {
					return []interface{}{}, nil
				},
				DiffSuppressFunc: suppressEquivalentStringDiffs,
			},
			"response_cache_shield": {
				Description: "Adds an intermediate cache between other Imperva PoPs and your origin servers to protect your servers from redundant requests.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"response_stale_content_mode": {
				Description:  "The working mode for serving stale content. Options are `disabled`, `adaptive`, and `custom`.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "disabled",
				ValidateFunc: validation.StringInSlice([]string{"disabled", "adaptive", "custom"}, false),
			},
			"response_stale_content_time": {
				Description: "The time, in seconds, to serve stale content for when working in `custom` work mode.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"response_tag_response_header": {
				Description: "Tag the response according to the value of this header. Specify which origin response header contains the cache tags in your resources.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ttl_prefer_last_modified": {
				Description: "Prefer 'Last Modified' over eTag. When this option is checked, Imperva prefers using Last Modified values (if available) over eTag values (recommended on multi-server setups).",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"ttl_use_shortest_caching": {
				Description: "Use shortest caching duration in case of conflicts. By default, the longest duration is used in case of conflict between caching rules or modes. When this option is checked, Imperva uses the shortest duration in case of conflict.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceApplicationPerformanceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	siteIdStr := strconv.Itoa(siteID)

	performanceSettings := PerformanceSettings{}
	performanceSettings.ClientSide.ComplyNoCache = d.Get("client_comply_no_cache").(bool)
	performanceSettings.ClientSide.EnableClientSideCaching = d.Get("client_enable_client_side_caching").(bool)
	performanceSettings.ClientSide.SendAgeHeader = d.Get("client_send_age_header").(bool)
	performanceSettings.Key.ComplyVary = d.Get("key_comply_vary").(bool)
	performanceSettings.Key.UniteNakedFullCache = d.Get("key_unite_naked_full_cache").(bool)
	performanceSettings.Mode.HTTPS = d.Get("mode_https").(string)
	performanceSettings.Mode.Level = d.Get("mode_level").(string)
	performanceSettings.Mode.Time = d.Get("mode_time").(int)
	performanceSettings.Response.Cache300X = d.Get("response_cache_300x").(bool)
	performanceSettings.Response.Cache404.Enabled = d.Get("response_cache_404_enabled").(bool)
	performanceSettings.Response.Cache404.Time = d.Get("response_cache_404_time").(int)
	performanceSettings.Response.CacheEmptyResponses = d.Get("response_cache_empty_responses").(bool)
	performanceSettings.Response.CacheHTTP10Responses = d.Get("response_cache_http_10_responses").(bool)
	performanceSettings.Response.CacheResponseHeader.Mode = d.Get("response_cache_response_header_mode").(string)
	performanceSettings.Response.CacheResponseHeader.Headers = d.Get("response_cache_response_headers").([]interface{})
	performanceSettings.Response.CacheShield = d.Get("response_cache_shield").(bool)
	performanceSettings.Response.StaleContent.Mode = d.Get("response_stale_content_mode").(string)
	performanceSettings.Response.StaleContent.Time = d.Get("response_stale_content_time").(int)
	performanceSettings.Response.TagResponseHeader = d.Get("response_tag_response_header").(string)
	performanceSettings.TTL.PreferLastModified = d.Get("ttl_prefer_last_modified").(bool)
	performanceSettings.TTL.UseShortestCaching = d.Get("ttl_use_shortest_caching").(bool)

	_, err := client.UpdatePerformanceSettings(siteIdStr, &performanceSettings)
	if err != nil {
		log.Printf("[ERROR] Could not update Incapsula performance settings for site_id: %s %s\n", d.Id(), err)
		return err
	}

	return resourceApplicationPerformanceRead(d, m)
}

func resourceApplicationPerformanceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	siteIdStr := strconv.Itoa(siteID)

	performanceSettingsResponse, err := client.GetPerformanceSettings(siteIdStr)
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site peformance settings for site id: %d, %s\n", siteID, err)
		return err
	}

	d.Set("client_comply_no_cache", performanceSettingsResponse.ClientSide.ComplyNoCache)
	d.Set("client_enable_client_side_caching", performanceSettingsResponse.ClientSide.EnableClientSideCaching)
	d.Set("client_send_age_header", performanceSettingsResponse.ClientSide.SendAgeHeader)
	d.Set("key_comply_vary", performanceSettingsResponse.Key.ComplyVary)
	d.Set("key_unite_naked_full_cache", performanceSettingsResponse.Key.UniteNakedFullCache)
	d.Set("mode_https", performanceSettingsResponse.Mode.HTTPS)
	d.Set("mode_level", performanceSettingsResponse.Mode.Level)
	d.Set("mode_time", performanceSettingsResponse.Mode.Time)
	d.Set("response_cache_300x", performanceSettingsResponse.Response.Cache300X)
	d.Set("response_cache_404_enabled", performanceSettingsResponse.Response.Cache404.Enabled)
	d.Set("response_cache_404_time", performanceSettingsResponse.Response.Cache404.Time)
	d.Set("response_cache_empty_responses", performanceSettingsResponse.Response.CacheEmptyResponses)
	d.Set("response_cache_http_10_responses", performanceSettingsResponse.Response.CacheHTTP10Responses)
	d.Set("response_cache_response_header_mode", performanceSettingsResponse.Response.CacheResponseHeader.Mode)
	d.Set("response_cache_response_headers", performanceSettingsResponse.Response.CacheResponseHeader.Headers)
	d.Set("response_cache_shield", performanceSettingsResponse.Response.CacheShield)
	d.Set("response_stale_content_mode", performanceSettingsResponse.Response.StaleContent.Mode)
	d.Set("response_stale_content_time", performanceSettingsResponse.Response.StaleContent.Time)
	d.Set("response_tag_response_header", performanceSettingsResponse.Response.TagResponseHeader)
	d.Set("ttl_prefer_last_modified", performanceSettingsResponse.TTL.PreferLastModified)
	d.Set("ttl_use_shortest_caching", performanceSettingsResponse.TTL.UseShortestCaching)

	d.SetId(siteIdStr)

	return nil
}

func resourceApplicationPerformanceDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
