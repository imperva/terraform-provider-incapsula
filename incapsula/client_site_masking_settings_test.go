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
// GetMaskingSettings Tests
////////////////////////////////////////////////////////////////

func TestClientGetMaskingSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "123"
	maskingSettings, err := client.GetMaskingSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading masking settings for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if maskingSettings != nil {
		t.Errorf("Should have received a nil maskingSettings instance")
	}
}

func TestClientGetMaskingSettingsBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/masking?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	maskingSettings, err := client.GetMaskingSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap masking settings JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if maskingSettings != nil {
		t.Errorf("Should have received a nil maskingSettings instance")
	}
}

func TestClientGetMaskingSettingsInvalidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/masking?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	maskingSettings, err := client.GetMaskingSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when reading masking settings for Site ID %s", 404, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if maskingSettings != nil {
		t.Errorf("Should have received a nil maskingSettings instance")
	}
}

func TestClientGetMaskingSettingsValidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/masking?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"hashing_enabled":true,"hash_salt":"abc"}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	maskingSettings, err := client.GetMaskingSettings(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if maskingSettings == nil {
		t.Errorf("Should not have received a nil maskingSettings instance")
	}
	if !maskingSettings.HashingEnabled {
		t.Errorf("Hashing should be enabled")
	}
	if maskingSettings.HashSalt != "abc" {
		t.Errorf("Hash salt doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// UpdateMaskingSettings Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateMaskingSettingsSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "123"
	maskingSettings := MaskingSettings{HashingEnabled: true, HashSalt: "salt"}
	err := client.UpdateMaskingSettings(siteID, &maskingSettings)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating masking settings for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdateMaskingSettingsInvalidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	maskingSettings := MaskingSettings{HashingEnabled: true, HashSalt: "salt"}

	endpoint := fmt.Sprintf("/sites/%s/settings/masking?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.UpdateMaskingSettings(siteID, &maskingSettings)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when updating masking settings for Site ID %s", 404, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientUpdateMaskingSettingsValidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	maskingSettings := MaskingSettings{HashingEnabled: true, HashSalt: "salt"}

	endpoint := fmt.Sprintf("/sites/%s/settings/masking?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"hashing_enabled":true,"hash_salt":"abc"}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.UpdateMaskingSettings(siteID, &maskingSettings)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
