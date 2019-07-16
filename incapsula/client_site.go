package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

const endpointSiteAdd = "sites/add"
const endpointSiteStatus = "sites/status"
const endpointSiteUpdate = "sites/configure"
const endpointSiteDelete = "sites/delete"

// SiteAddResponse contains the relevant site information when adding an Incapsula managed site
type SiteAddResponse struct {
	SiteID int `json:"site_id"`
	Res    int `json:"res"`
}

// SiteUpdateResponse contains the relevant site information when updating an Incapsula managed site
type SiteUpdateResponse struct {
	SiteID int `json:"site_id"`
	Res    int `json:"res"`
}

// SiteStatusResponse contains managed site information
type SiteStatusResponse struct {
	SiteID            int      `json:"site_id"`
	Status            string   `json:"status"`
	Domain            string   `json:"domain"`
	RefID             string   `json:"ref_id,omitempty"`
	AccountID         int      `json:"account_id"`
	AccelerationLevel string   `json:"acceleration_level"`
	SiteCreationDate  int64    `json:"site_creation_date"`
	Ips               []string `json:"ips"`
	DNS               []struct {
		DNSRecordName string   `json:"dns_record_name"`
		SetTypeTo     string   `json:"set_type_to"`
		SetDataTo     []string `json:"set_data_to"`
	} `json:"dns"`
	OriginalDNS []struct {
		DNSRecordName string   `json:"dns_record_name"`
		SetTypeTo     string   `json:"set_type_to"`
		SetDataTo     []string `json:"set_data_to"`
	} `json:"original_dns"`
	Warnings                     []interface{} `json:"warnings"`
	Active                       string        `json:"active"`
	SupportAllTLSVersions        bool          `json:"support_all_tls_versions"`
	WildcardSanForNewSites       bool          `json:"wildcard_san_for_new_sites"`
	NakedDomainSanForNewWwwSites bool          `json:"naked_domain_san_for_new_www_sites"`
	AdditionalErrors             []interface{} `json:"additionalErrors"`
	DisplayName                  string        `json:"display_name"`
	Security                     struct {
		Waf struct {
			Rules []struct {
				Action                 string `json:"action,omitempty"`
				ActionText             string `json:"action_text,omitempty"`
				ID                     string `json:"id"`
				Name                   string `json:"name"`
				BlockBadBots           bool   `json:"block_bad_bots,omitempty"`
				ChallengeSuspectedBots bool   `json:"challenge_suspected_bots,omitempty"`
				ActivationMode         string `json:"activation_mode,omitempty"`
				ActivationModeText     string `json:"activation_mode_text,omitempty"`
				DdosTrafficThreshold   int    `json:"ddos_traffic_threshold,omitempty"`
				Exceptions             []struct {
					Values []struct {
						ID   string   `json:"id,omitempty"`
						Name string   `json:"name,omitempty"`
						Ips  []string `json:"ips,omitempty"`
						Urls []struct {
							Value   string `json:"value,omitempty"`
							Pattern string `json:"pattern,omitempty"`
						} `json:"urls,omitempty"`
						Geo struct {
							Countries  []string `json:"countries,omitempty"`
							Continents []string `json:"continents,omitempty"`
						} `json:"geo,omitempty"`
						ClientApps     []string `json:"client_apps,omitempty"`
						ClientAppTypes []string `json:"client_app_types,omitempty"`
						Parameters     []string `json:"parameters,omitempty"`
						UserAgents     []string `json:"user_agents,omitempty"`
					} `json:"values,omitempty"`
					ID int `json:"id,omitempty"`
				} `json:"exceptions,omitempty"`
			} `json:"rules"`
		} `json:"waf"`
		Acls struct {
			Rules []struct {
				Ips  []string `json:"ips,omitempty"`
				ID   string   `json:"id"`
				Name string   `json:"name"`
				Geo  struct {
					Countries  []string `json:"countries"`
					Continents []string `json:"continents"`
				} `json:"geo,omitempty"`
				Urls []struct {
					Value   string `json:"value"`
					Pattern string `json:"pattern"`
				} `json:"urls,omitempty"`
				Exceptions []struct {
					Values []struct {
						ID   string   `json:"id"`
						Name string   `json:"name"`
						Ips  []string `json:"ips,omitempty"`
						Urls []struct {
							Value   string `json:"value"`
							Pattern string `json:"pattern"`
						} `json:"urls,omitempty"`
						Geo struct {
							Countries  []string `json:"countries"`
							Continents []string `json:"continents"`
						} `json:"geo,omitempty"`
						ClientApps     []string `json:"client_apps,omitempty"`
						ClientAppTypes []string `json:"client_app_types,omitempty"`
						Parameters     []string `json:"parameters,omitempty"`
						UserAgents     []string `json:"user_agents,omitempty"`
					} `json:"values"`
					ID int `json:"id"`
				} `json:"exceptions"`
			} `json:"rules"`
		} `json:"acls"`
	} `json:"security"`
	SealLocation struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"sealLocation"`
	Ssl struct {
		OriginServer struct {
			Detected        bool   `json:"detected"`
			DetectionStatus string `json:"detectionStatus"`
		} `json:"origin_server"`
		GeneratedCertificate struct {
			San []interface{} `json:"san"`
		} `json:"generated_certificate"`
	} `json:"ssl"`
	SiteDualFactorSettings struct {
		SpecificUsers                []interface{} `json:"specificUsers"`
		Enabled                      bool          `json:"enabled"`
		CustomAreas                  []interface{} `json:"customAreas"`
		CustomAreasExceptions        []interface{} `json:"customAreasExceptions"`
		AllowAllUsers                bool          `json:"allowAllUsers"`
		ShouldSuggestApplicatons     bool          `json:"shouldSuggestApplicatons"`
		AllowedMedia                 []string      `json:"allowedMedia"`
		ShouldSendLoginNotifications bool          `json:"shouldSendLoginNotifications"`
		Version                      int           `json:"version"`
	} `json:"siteDualFactorSettings"`
	LoginProtect struct {
		Enabled               bool          `json:"enabled"`
		SpecificUsersList     []interface{} `json:"specific_users_list"`
		SendLpNotifications   bool          `json:"send_lp_notifications"`
		AllowAllUsers         bool          `json:"allow_all_users"`
		AuthenticationMethods []string      `json:"authentication_methods"`
		Urls                  []interface{} `json:"urls"`
		URLPatterns           []interface{} `json:"url_patterns"`
	} `json:"login_protect"`
	PerformanceConfiguration struct {
		AdvancedCachingRules struct {
			NeverCacheResources  []interface{} `json:"never_cache_resources"`
			AlwaysCacheResources []interface{} `json:"always_cache_resources"`
		} `json:"advanced_caching_rules"`
		AccelerationLevel         string        `json:"acceleration_level"`
		AsyncValidation           bool          `json:"async_validation"`
		MinifyJavascript          bool          `json:"minify_javascript"`
		MinifyCSS                 bool          `json:"minify_css"`
		MinifyStaticHTML          bool          `json:"minify_static_html"`
		CompressJpeg              bool          `json:"compress_jpeg"`
		CompressJepg              bool          `json:"compress_jepg"`
		ProgressiveImageRendering bool          `json:"progressive_image_rendering"`
		AggressiveCompression     bool          `json:"aggressive_compression"`
		CompressPng               bool          `json:"compress_png"`
		OnTheFlyCompression       bool          `json:"on_the_fly_compression"`
		TCPPrePooling             bool          `json:"tcp_pre_pooling"`
		ComplyNoCache             bool          `json:"comply_no_cache"`
		ComplyVary                bool          `json:"comply_vary"`
		UseShortestCaching        bool          `json:"use_shortest_caching"`
		PerferLastModified        bool          `json:"perfer_last_modified"`
		PreferLastModified        bool          `json:"prefer_last_modified"`
		DisableClientSideCaching  bool          `json:"disable_client_side_caching"`
		Cache300X                 bool          `json:"cache300x"`
		CacheHeaders              []interface{} `json:"cache_headers"`
	} `json:"performance_configuration"`
	ExtendedDdos int    `json:"extended_ddos"`
	ExceptionID  string `json:"exception_id,omitempty"`
	Res          int    `json:"res"`
	ResMessage   string `json:"res_message"`
	DebugInfo    struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// AddSite adds a site to be managed by Incapsula
func (c *Client) AddSite(domain, accountID, refID, sendSiteSetupEmails, siteIP, forceSSL string) (*SiteAddResponse, error) {
	log.Printf("[INFO] Adding Incapsula site for domain: %s\n", domain)

	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteAdd), url.Values{
		"api_id":                 {c.config.APIID},
		"api_key":                {c.config.APIKey},
		"domain":                 {domain},
		"account_id":             {accountID},
		"ref_id":                 {refID},
		"send_site_setup_emails": {sendSiteSetupEmails},
		"site_ip":                {siteIP},
		"force_ssl":              {forceSSL},
	})
	if err != nil {
		return nil, fmt.Errorf("Error adding site for domain %s: %s", domain, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula add site JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteAddResponse SiteAddResponse
	err = json.Unmarshal([]byte(responseBody), &siteAddResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add site JSON response for domain %s: %s", domain, err)
	}

	// Look at the response status code from Incapsula
	if siteAddResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding site for domain %s: %s", domain, string(responseBody))
	}

	return &siteAddResponse, nil
}

// SiteStatus gets the Incapsula managed site's status
func (c *Client) SiteStatus(domain string, siteID int) (*SiteStatusResponse, error) {
	log.Printf("[INFO] Getting Incapsula site status for domain: %s (site id: %d)\n", domain, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteStatus), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
	})
	if err != nil {
		return nil, fmt.Errorf("Error getting site status for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula site status JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteStatusResponse SiteStatusResponse
	err = json.Unmarshal([]byte(responseBody), &siteStatusResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing site status JSON response for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Look at the response status code from Incapsula
	if siteStatusResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when getting site status for domain %s (site id: %d): %s", domain, siteID, string(responseBody))
	}

	return &siteStatusResponse, nil
}

// UpdateSite will update the specific param/value on the site resource
func (c *Client) UpdateSite(siteID, param, value string) (*SiteUpdateResponse, error) {
	log.Printf("[INFO] Updating Incapsula site for siteID: %s\n", siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteUpdate), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {siteID},
		"param":   {param},
		"value":   {value},
	})
	if err != nil {
		return nil, fmt.Errorf("Error updating param (%s) with value (%s) on site_id: %s: %s", param, value, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update site JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteUpdateResponse SiteUpdateResponse
	err = json.Unmarshal([]byte(responseBody), &siteUpdateResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update site JSON response for siteID %s: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if siteUpdateResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when updating site for siteID %s: %s", siteID, string(responseBody))
	}

	return &siteUpdateResponse, nil
}

// DeleteSite deletes a site currently managed by Incapsula
func (c *Client) DeleteSite(domain string, siteID int) error {
	// Specifically shaded this struct, no need to share across funcs or export
	// We only care about the response code and possibly the message
	type SiteDeleteResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
	}

	log.Printf("[INFO] Deleting Incapsula site for domain: %s (site id: %d)\n", domain, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteDelete), url.Values{
		"api_id":  {c.config.APIID},
		"api_key": {c.config.APIKey},
		"site_id": {strconv.Itoa(siteID)},
	})
	if err != nil {
		return fmt.Errorf("Error deleting site for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula delete site JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var siteDeleteResponse SiteDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &siteDeleteResponse)
	if err != nil {
		return fmt.Errorf("Error parsing delete site JSON response for domain %s (site id: %d): %s", domain, siteID, err)
	}

	// Look at the response status code from Incapsula
	if siteDeleteResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when deleting site for domain %s (site id: %d): %s", domain, siteID, string(responseBody))
	}

	return nil
}
