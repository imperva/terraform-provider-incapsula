package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////
// GetPerformanceSettings Tests
////////////////////////////////////////////////////////////////

func TestClientGetPerformanceSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "123"
	performanceSettings, _, err := client.GetPerformanceSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading Incap Performance Settings for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if performanceSettings != nil {
		t.Errorf("Should have received a nil performanceSettings instance")
	}
}

func TestClientGetPerformanceSettingsBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, APIV2BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	performanceSettings, _, err := client.GetPerformanceSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap Performance Settings JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if performanceSettings != nil {
		t.Errorf("Should have received a nil performanceSettings instance")
	}
}

func TestClientGetPerformanceSettingsInvalidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, APIV2BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	performanceSettings, _, err := client.GetPerformanceSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when reading Incap Performance Settings for Site ID %s", 404, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if performanceSettings != nil {
		t.Errorf("Should have received a nil performanceSettings instance")
	}
}

func TestClientGetPerformanceSettingsValidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"mode":{"level":"standard","https":"include_all_resources","time":1000},"key":{"unite_naked_full_cache":true,"comply_vary":true},"response":{"stale_content":{"mode":"custom","time":1000},"cache_shield":true,"cache_response_header":{"mode":"custom","headers":["Access-Control-Allow-Origin","Foo-Bar-Header"]},"tag_response_header":"Example-Tag-Value-Header","cache_empty_responses":true,"cache_300x":true,"cache_http_10_responses":true,"cache_404":{"enabled":true,"time":1000}},"ttl":{"use_shortest_caching":true,"prefer_last_modified":true},"client_side":{"enable_client_side_caching":true,"comply_no_cache":true,"send_age_header":true}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, APIV2BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	performanceSettings, _, err := client.GetPerformanceSettings(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if performanceSettings == nil {
		t.Errorf("Should not have received a nil performanceSettings instance")
	}
	if performanceSettings.Mode.HTTPS != "include_all_resources" {
		t.Errorf("Performance settings mode HTTPS should be include_all_resources")
	}
}

////////////////////////////////////////////////////////////////
// UpdatePerformanceSettings Tests
////////////////////////////////////////////////////////////////

func TestClientUpdatePerformanceSettingsSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "123"
	performanceSettings := PerformanceSettings{}
	performanceSettings.Mode.HTTPS = "include_all_resources"
	_, err := client.UpdatePerformanceSettings(siteID, &performanceSettings)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating Incap Performance Settings for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdatePerformanceSettingsInvalidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	performanceSettings := PerformanceSettings{}
	performanceSettings.Mode.HTTPS = "include_all_resources"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, APIV2BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.UpdatePerformanceSettings(siteID, &performanceSettings)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when updating Incap Performance Settings for Site ID %s", 404, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientUpdatePerformanceSettingsValidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	performanceSettings := PerformanceSettings{}
	performanceSettings.Mode.HTTPS = "include_all_resources"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"mode":{"level":"standard","https":"include_all_resources","time":1000},"key":{"unite_naked_full_cache":true,"comply_vary":true},"response":{"stale_content":{"mode":"custom","time":1000},"cache_shield":true,"cache_response_header":{"mode":"custom","headers":["Access-Control-Allow-Origin","Foo-Bar-Header"]},"tag_response_header":"Example-Tag-Value-Header","cache_empty_responses":true,"cache_300x":true,"cache_http_10_responses":true,"cache_404":{"enabled":true,"time":1000}},"ttl":{"use_shortest_caching":true,"prefer_last_modified":true},"client_side":{"enable_client_side_caching":true,"comply_no_cache":true,"send_age_header":true}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, APIV2BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.UpdatePerformanceSettings(siteID, &performanceSettings)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
