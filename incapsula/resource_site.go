package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"domain": {
				Description: "The fully qualified domain name of the site. For example: www.example.com, hello.example.com.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					d := val.(string)
					parts := strings.Split(d, ".")
					if len(parts) <= 2 {
						errs = append(errs, fmt.Errorf("%q must be a fully qualified domain name (www.example.com, not example.com), got: %s", key, d))
					}
					return
				},
			},

			// Optional Arguments
			"account_id": {
				Description: "Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ref_id": {
				Description: "Customer specific identifier for this operation.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"send_site_setup_emails": {
				Description: "If this value is false, end users will not get emails about the add site process such as DNS instructions and SSL setup.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"site_ip": {
				Description: "Manually set the web server IP/CNAME.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"force_ssl": {
				Description: "If this value is true, manually set the site to support SSL. This option is only available for sites with manually configured IP/CNAME and for specific accounts.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"logs_account_id": {
				Description: "Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"active": {
				Description: "active or bypass.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_validation": {
				Description: "email or html or dns.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"approver": {
				Description: "my.approver@email.com (some approver email address).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"ignore_ssl": {
				Description: "true or empty string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"acceleration_level": {
				Description: "none | standard | aggressive.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"seal_location": {
				Description: "api.seal_location.bottom_left | api.seal_location.none | api.seal_location.right_bottom | api.seal_location.right | api.seal_location.left | api.seal_location.bottom_right | api.seal_location.bottom.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domain_redirect_to_full": {
				Description: "true or empty string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"remove_ssl": {
				Description: "true or empty string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"data_storage_region": {
				Description: "The data region to use. Options are `APAC`, `AU`, `EU`, and `US`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"hashing_enabled": {
				Description: "Specify if hashing (masking setting) should be enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"hash_salt": {
				Description: "Specify the hash salt (masking setting), required if hashing is enabled. Maximum length of 64 characters.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					salt := val.(string)
					if len(salt) > 64 {
						errs = append(errs, fmt.Errorf("%q must be a max of 64 characters, got: %s", key, salt))
					}
					return
				},
			},
			"log_level": {
				Description: "The log level. Options are `full`, `security`, and `none`. Defaults to `none`.",
				Type:        schema.TypeString,
				Default:     "none",
				Optional:    true,
			},
			"perf_client_comply_no_cache": {
				Description: "Comply with No-Cache and Max-Age directives in client requests. By default, these cache directives are ignored. Resources are dynamically profiled and re-configured to optimize performance.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_client_enable_client_side_caching": {
				Description: "Cache content on client browsers or applications. When not enabled, content is cached only on the Imperva proxies.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_client_send_age_header": {
				Description: "Send Cache-Control: max-age and Age headers.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_key_comply_vary": {
				Description: "Comply with Vary. Cache resources in accordance with the Vary response header.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_key_unite_naked_full_cache": {
				Description: "Use the Same Cache for Full and Naked Domains. For example, use the same cached resource for www.example.com/a and example.com/a.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_mode_https": {
				Description: "The resources that are cached over HTTPS, the general level applies. Options are `disabled`, `dont_include_html`, `include_html`, and `include_all_resources`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"perf_mode_level": {
				Description: "Caching level. Options are `disable`, `standard`, `smart`, and `all_resources`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"perf_mode_time": {
				Description: "The time, in seconds, that you set for this option determines how often the cache is refreshed. Relevant for the `include_html` and `include_all_resources` levels only.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_300x": {
				Description: "When this option is checked Imperva will cache 301, 302, 303, 307, and 308 redirect response headers containing the target URI.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_404_enabled": {
				Description: "Whether or not to cache 404 responses.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_404_time": {
				Description: "The time in seconds to cache 404 responses.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_empty_responses": {
				Description: "Cache responses that don’t have a message body.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_http_10_responses": {
				Description: "Cache HTTP 1.0 type responses that don’t include the Content-Length header or chunking.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_response_header_mode": {
				Description: "The working mode for caching response headers. Options are `all` and `custom`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_cache_response_headers": {
				Description: "An array of strings representing the response headers to be cached when working in `custom` mode. If empty, no response headers are cached.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppressEquivalentStringDiffs,
			},
			"perf_response_cache_shield": {
				Description: "Adds an intermediate cache between other Imperva PoPs and your origin servers to protect your servers from redundant requests.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_stale_content_mode": {
				Description: "The working mode for serving stale content. Options are `disabled`, `adaptive`, and `custom`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_stale_content_time": {
				Description: "The time, in seconds, to serve stale content for when working in `custom` work mode.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
			"perf_response_tag_response_header": {
				Description: "Tag the response according to the value of this header. Specify which origin response header contains the cache tags in your resources.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"perf_ttl_prefer_last_modified": {
				Description: "Prefer 'Last Modified' over eTag. When this option is checked, Imperva prefers using Last Modified values (if available) over eTag values (recommended on multi-server setups).",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"perf_ttl_use_shortest_caching": {
				Description: "Use shortest caching duration in case of conflicts. By default, the longest duration is used in case of conflict between caching rules or modes. When this option is checked, Imperva uses the shortest duration in case of conflict.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},

			// Computed Attributes
			"site_creation_date": {
				Description: "Numeric representation of the site creation date.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"dns_cname_record_name": {
				Description: "CNAME record name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_cname_record_value": {
				Description: "CNAME record value.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_a_record_name": {
				Description: "A record name.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"dns_a_record_value": {
				Description: "A record value.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"domain_verification": {
				Description: "Domain verification (e.g. GlobalSign verification).",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceSiteCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	domain := d.Get("domain").(string)

	log.Printf("[INFO] Creating Incapsula site for domain: %s\n", domain)

	siteAddResponse, err := client.AddSite(
		domain,
		d.Get("account_id").(string),
		d.Get("ref_id").(string),
		d.Get("send_site_setup_emails").(string),
		d.Get("site_ip").(string),
		d.Get("force_ssl").(string),
	)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula site for domain: %s, %s\n", domain, err)
		return err
	}

	// Set the Site ID
	d.SetId(strconv.Itoa(siteAddResponse.SiteID))
	log.Printf("[INFO] Created Incapsula site for domain: %s\n", domain)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	err = updateAdditionalSiteProperties(client, d)
	if err != nil {
		return err
	}

	err = updateDataStorageRegion(client, d)
	if err != nil {
		return err
	}

	err = updateMaskingSettings(client, d)
	if err != nil {
		return err
	}

	err = updateLogLevel(client, d)
	if err != nil {
		return err
	}

	err = updatePerformanceSettings(client, d)
	if err != nil {
		return err
	}

	// Set the rest of the state from the resource read
	return resourceSiteRead(d, m)
}

func resourceSiteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	domain := d.Get("domain").(string)
	siteID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Reading Incapsula site for domain: %s\n", domain)

	siteStatusResponse, err := client.SiteStatus(domain, siteID)

	// Site object may have been deleted
	if siteStatusResponse != nil && siteStatusResponse.Res.(float64) == 9413 {
		log.Printf("[INFO] Incapsula Site ID %d has already been deleted: %s\n", siteID, err)
		d.SetId("")
		return nil
	}

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site for domain: %s, %s\n", domain, err)
		return err
	}

	d.Set("site_creation_date", siteStatusResponse.SiteCreationDate)
	d.Set("domain", siteStatusResponse.Domain)
	d.Set("account_id", siteStatusResponse.AccountID)

	// Set the DNS information
	dnsARecordValues := make([]string, 0)
	for _, entry := range siteStatusResponse.DNS {
		if entry.SetTypeTo == "CNAME" && len(entry.SetDataTo) > 0 {
			d.Set("dns_cname_record_name", entry.DNSRecordName)
			d.Set("dns_cname_record_value", entry.SetDataTo[0])
		}
		if entry.SetTypeTo == "A" {
			d.Set("dns_a_record_name", entry.DNSRecordName)
			dnsARecordValues = append(dnsARecordValues, entry.SetDataTo...)
		}
	}
	d.Set("dns_a_record_value", dnsARecordValues)

	// Set the GlobalSign verification (may not exist)
	if siteStatusResponse.Ssl.GeneratedCertificate.ValidationMethod == "dns" {
		d.Set("domain_verification", siteStatusResponse.Ssl.GeneratedCertificate.ValidationData[0].SetDataTo[0])
	}

	// Get the log level for the site
	if siteStatusResponse.LogLevel != "" {
		d.Set("log_level", siteStatusResponse.LogLevel)
	}

	// Get the data storage region for the site
	dataStorageRegionResponse, err := client.GetDataStorageRegion(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site data storage region for domain: %s and site id: %d, %s\n", domain, siteID, err)
		return err
	}
	d.Set("data_storage_region", dataStorageRegionResponse.Region)

	// Get the masking settings for the site
	maskingResponse, err := client.GetMaskingSettings(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site masking settings for domain: %s and site id: %d, %s\n", domain, siteID, err)
		return err
	}
	d.Set("hashing_enabled", maskingResponse.HashingEnabled)
	d.Set("hash_salt", maskingResponse.HashSalt)

	// Get the performance settings for the site
	performanceSettingsResponse, _, err := client.GetPerformanceSettings(d.Id())
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site peformance settings for domain: %s and site id: %d, %s\n", domain, siteID, err)
		return err
	}
	d.Set("perf_client_comply_no_cache", performanceSettingsResponse.ClientSide.ComplyNoCache)
	d.Set("perf_client_enable_client_side_caching", performanceSettingsResponse.ClientSide.EnableClientSideCaching)
	d.Set("perf_client_send_age_header", performanceSettingsResponse.ClientSide.SendAgeHeader)
	d.Set("perf_key_comply_vary", performanceSettingsResponse.Key.ComplyVary)
	d.Set("perf_key_unite_naked_full_cache", performanceSettingsResponse.Key.UniteNakedFullCache)
	d.Set("perf_mode_https", performanceSettingsResponse.Mode.HTTPS)
	d.Set("perf_mode_level", performanceSettingsResponse.Mode.Level)
	d.Set("perf_mode_time", performanceSettingsResponse.Mode.Time)
	d.Set("perf_response_cache_300x", performanceSettingsResponse.Response.Cache300X)
	d.Set("perf_response_cache_404_enabled", performanceSettingsResponse.Response.Cache404.Enabled)
	d.Set("perf_response_cache_404_time", performanceSettingsResponse.Response.Cache404.Time)
	d.Set("perf_response_cache_empty_responses", performanceSettingsResponse.Response.CacheEmptyResponses)
	d.Set("perf_response_cache_http_10_responses", performanceSettingsResponse.Response.CacheHTTP10Responses)
	d.Set("perf_response_cache_response_header_mode", performanceSettingsResponse.Response.CacheResponseHeader.Mode)
	d.Set("perf_response_cache_response_headers", performanceSettingsResponse.Response.CacheResponseHeader.Headers)
	d.Set("perf_response_cache_shield", performanceSettingsResponse.Response.CacheShield)
	d.Set("perf_response_stale_content_mode", performanceSettingsResponse.Response.StaleContent.Mode)
	d.Set("perf_response_stale_content_time", performanceSettingsResponse.Response.StaleContent.Time)
	d.Set("perf_response_tag_response_header", performanceSettingsResponse.Response.TagResponseHeader)
	d.Set("perf_ttl_prefer_last_modified", performanceSettingsResponse.TTL.PreferLastModified)
	d.Set("perf_ttl_use_shortest_caching", performanceSettingsResponse.TTL.UseShortestCaching)

	log.Printf("[INFO] Finished reading Incapsula site for domain: %s\n", domain)

	return nil
}

func resourceSiteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := updateAdditionalSiteProperties(client, d)
	if err != nil {
		return err
	}

	err = updateDataStorageRegion(client, d)
	if err != nil {
		return err
	}

	err = updateMaskingSettings(client, d)
	if err != nil {
		return err
	}

	err = updateLogLevel(client, d)
	if err != nil {
		return err
	}

	err = updatePerformanceSettings(client, d)
	if err != nil {
		return err
	}

	// Set the rest of the state from the resource read
	return resourceSiteRead(d, m)
}

func resourceSiteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	domain := d.Get("domain").(string)
	siteID, _ := strconv.Atoi(d.Id())

	log.Printf("[INFO] Deleting Incapsula site for domain: %s\n", domain)

	err := client.DeleteSite(domain, siteID)

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula site for domain: %s, %s\n", domain, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	log.Printf("[INFO] Deleted Incapsula site for domain: %s\n", domain)

	return nil
}

func updateAdditionalSiteProperties(client *Client, d *schema.ResourceData) error {
	updateParams := [7]string{"acceleration_level", "active", "approver", "domain_redirect_to_full", "domain_validation", "ignore_ssl", "remove_ssl"}
	for i := 0; i < len(updateParams); i++ {
		param := updateParams[i]
		if d.HasChange(param) && d.Get(param) != "" {
			log.Printf("[INFO] Updating Incapsula site param (%s) with value (%s) for site_id: %s\n", param, d.Get(param).(string), d.Id())
			_, err := client.UpdateSite(d.Id(), param, d.Get(param).(string))
			if err != nil {
				log.Printf("[ERROR] Could not update Incapsula site param (%s) with value (%s) for site_id: %s %s\n", param, d.Get(param).(string), d.Id(), err)
				return err
			}
		}
	}
	return nil
}

func updateDataStorageRegion(client *Client, d *schema.ResourceData) error {
	if d.HasChange("data_storage_region") {
		dataStorageRegion := d.Get("data_storage_region").(string)
		_, err := client.UpdateDataStorageRegion(d.Id(), dataStorageRegion)
		if err != nil {
			log.Printf("[ERROR] Could not set Incapsula site data storage region with value (%s) for site_id: %s %s\n", dataStorageRegion, d.Id(), err)
			return err
		}
	}
	return nil
}

func updateMaskingSettings(client *Client, d *schema.ResourceData) error {
	if d.HasChange("hashing_enabled") || d.HasChange("hash_salt") {
		hashingEnabled := d.Get("hashing_enabled").(bool)
		hashSalt := d.Get("hash_salt").(string)
		maskingSettings := MaskingSettings{HashingEnabled: hashingEnabled, HashSalt: hashSalt}
		err := client.UpdateMaskingSettings(d.Id(), &maskingSettings)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula site masking settings for site_id: %s %s\n", d.Id(), err)
			return err
		}
	}
	return nil
}

func updateLogLevel(client *Client, d *schema.ResourceData) error {
	if d.HasChange("log_level") {
		logLevel := d.Get("log_level").(string)
		err := client.UpdateLogLevel(d.Id(), logLevel)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula site log level: %s for site_id: %s %s\n", logLevel, d.Id(), err)
			return err
		}
	}
	return nil
}

func updatePerformanceSettings(client *Client, d *schema.ResourceData) error {
	if d.HasChange("perf_client_comply_no_cache") ||
		d.HasChange("perf_client_enable_client_side_caching") ||
		d.HasChange("perf_client_send_age_header") ||
		d.HasChange("perf_key_comply_vary") ||
		d.HasChange("perf_key_unite_naked_full_cache") ||
		d.HasChange("perf_mode_https") ||
		d.HasChange("perf_mode_level") ||
		d.HasChange("perf_mode_time") ||
		d.HasChange("perf_response_cache_300x") ||
		d.HasChange("perf_response_cache_404_enabled") ||
		d.HasChange("perf_response_cache_404_time") ||
		d.HasChange("perf_response_cache_empty_responses") ||
		d.HasChange("perf_response_cache_http_10_responses") ||
		d.HasChange("perf_response_cache_response_header_mode") ||
		d.HasChange("perf_response_cache_response_headers") ||
		d.HasChange("perf_response_cache_shield") ||
		d.HasChange("perf_response_stale_content_mode") ||
		d.HasChange("perf_response_stale_content_time") ||
		d.HasChange("perf_response_tag_response_header") ||
		d.HasChange("perf_ttl_prefer_last_modified") ||
		d.HasChange("perf_ttl_use_shortest_caching") {
		performanceSettings := PerformanceSettings{}
		performanceSettings.ClientSide.ComplyNoCache = d.Get("perf_client_comply_no_cache").(bool)
		performanceSettings.ClientSide.EnableClientSideCaching = d.Get("perf_client_enable_client_side_caching").(bool)
		performanceSettings.ClientSide.SendAgeHeader = d.Get("perf_client_send_age_header").(bool)
		performanceSettings.Key.ComplyVary = d.Get("perf_key_comply_vary").(bool)
		performanceSettings.Key.UniteNakedFullCache = d.Get("perf_key_unite_naked_full_cache").(bool)
		performanceSettings.Mode.HTTPS = d.Get("perf_mode_https").(string)
		performanceSettings.Mode.Level = d.Get("perf_mode_level").(string)
		performanceSettings.Mode.Time = d.Get("perf_mode_time").(int)
		performanceSettings.Response.Cache300X = d.Get("perf_response_cache_300x").(bool)
		performanceSettings.Response.Cache404.Enabled = d.Get("perf_response_cache_404_enabled").(bool)
		performanceSettings.Response.Cache404.Time = d.Get("perf_response_cache_404_time").(int)
		performanceSettings.Response.CacheEmptyResponses = d.Get("perf_response_cache_empty_responses").(bool)
		performanceSettings.Response.CacheHTTP10Responses = d.Get("perf_response_cache_http_10_responses").(bool)
		performanceSettings.Response.CacheResponseHeader.Mode = d.Get("perf_response_cache_response_header_mode").(string)
		performanceSettings.Response.CacheResponseHeader.Headers = d.Get("perf_response_cache_response_headers").([]interface{})
		performanceSettings.Response.CacheShield = d.Get("perf_response_cache_shield").(bool)
		performanceSettings.Response.StaleContent.Mode = d.Get("perf_response_stale_content_mode").(string)
		performanceSettings.Response.StaleContent.Time = d.Get("perf_response_stale_content_time").(int)
		performanceSettings.Response.TagResponseHeader = d.Get("perf_response_tag_response_header").(string)
		performanceSettings.TTL.PreferLastModified = d.Get("perf_ttl_prefer_last_modified").(bool)
		performanceSettings.TTL.UseShortestCaching = d.Get("perf_ttl_use_shortest_caching").(bool)

		_, err := client.UpdatePerformanceSettings(d.Id(), &performanceSettings)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula performance settings for site_id: %s %s\n", d.Id(), err)
			return err
		}
	}
	return nil
}
