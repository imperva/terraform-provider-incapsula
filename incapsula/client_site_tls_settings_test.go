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
// GetSiteTlsSettings Tests
////////////////////////////////////////////////////////////////
func TestGetSiteTlsSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	siteTlsSettings, err := client.GetSiteTlsSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error getting Site TLS Settings for Site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteTlsSettings != nil {
		t.Errorf("Should have received a nil siteTlsSettings instance")
	}

}

func TestGetSiteTlsSettingsBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/certificate-manager/v2/sites/%d%s", siteID, endpointSiteTlsSettings)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteTlsSettings, err := client.GetSiteTlsSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprint("[ERROR] Error parsing Incapsula Site TLS Settings")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteTlsSettings != nil {
		t.Errorf("Should have received a nil siteTlsSettings instance")
	}
}

func TestGetSiteTlsSettingsInvalidApiConfigInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/certificate-manager/v2/sites/%d%s", siteID, endpointSiteTlsSettings)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(401)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteTlsSettings, err := client.GetSiteTlsSettings(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 401 from Incapsula service on fetching Incapsula Site TLS Settings for Site ID %d", siteID)) {
		t.Errorf("Should have received an invalid status code error, got: %s", err)
	}
	if siteTlsSettings != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestGetSiteTlsSettingsValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/certificate-manager/v2/sites/%d%s", siteID, endpointSiteTlsSettings)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "mandatory": true,
  "ports": [
    80,
    9000
  ],
  "isPortsException": false,
  "hosts": [
    "imperva.com",
    "imprevaservices.com"
  ],
  "isHostsException": false,
  "fingerprints": ["F009B2EABECCBE9BFBE23B8C57A684650B8564A9"],
  "forwardToOrigin": true,
  "headerName": "Full-Cert",
  "headerValue": "FULL_CERT",
  "isDisableSessionResumption": true
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteTlsSettings, err := client.GetSiteTlsSettings(siteID)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), siteTlsSettings)
	}
	if siteTlsSettings == nil {
		t.Errorf("Should not have received a nil siteTlsSettings instance")
	}
	if siteTlsSettings.Mandatory != true {
		t.Errorf("Parameter Mandatory doesn't match. Actual : %v", siteTlsSettings.Mandatory)
	}
}

////////////////////////////////////////////////////////////////
// UpdateSiteTlsSetings Tests
////////////////////////////////////////////////////////////////
func TestUpdateSiteTlsSetingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	siteTlsSettingsPayload := SiteTlsSettings{}

	siteTlsSettings, err := client.UpdateSiteTlsSetings(siteID, siteTlsSettingsPayload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error updating Site TLS Settings for Site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteTlsSettings != nil {
		t.Errorf("Should have received a nil siteTlsSettings instance")
	}

}

func TestUpdateSiteTlsSetingsBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	siteTlsSettingsPayload := SiteTlsSettings{}

	endpoint := fmt.Sprintf("/certificate-manager/v2/sites/%d/configuration/client-certificates", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteTlsSettings, err := client.UpdateSiteTlsSetings(siteID, siteTlsSettingsPayload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing Incapsula Site TLS Settings JSON response for Site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteTlsSettings != nil {
		t.Errorf("Should have received a nil siteTlsSettings instance")
	}
}

func TestUpdateSiteTlsSetingsInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	siteTlsSettingsPayload := SiteTlsSettings{}

	endpoint := fmt.Sprintf("/certificate-manager/v2/sites/%d/configuration/client-certificates", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteTlsSettings, err := client.UpdateSiteTlsSetings(siteID, siteTlsSettingsPayload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 406 from Incapsula service on update Incapsula Site TLS Settings for Site ID %d", siteID)) {
		t.Errorf("Should have received an invalid status code error, got: %s", err)
	}
	if siteTlsSettings != nil {
		t.Errorf("Should have received a nil siteTlsSettings instance")
	}
}

func TestUpdateSiteTlsSetingsValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	siteTlsSettingsPayload := SiteTlsSettings{}

	endpoint := fmt.Sprintf("/certificate-manager/v2/sites/%d%s", siteID, endpointSiteTlsSettings)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "mandatory": true,
  "ports": [
    80,
    9000
  ],
  "isPortsException": false,
  "hosts": [
    "imperva.com",
    "imprevaservices.com"
  ],
  "isHostsException": false,
  "fingerprints": ["F009B2EABECCBE9BFBE23B8C57A684650B8564A9"],
  "forwardToOrigin": true,
  "headerName": "Full-Cert",
  "headerValue": "FULL_CERT",
  "isDisableSessionResumption": true
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteTlsSettings, err := client.UpdateSiteTlsSetings(siteID, siteTlsSettingsPayload)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), siteTlsSettings)
	}
	if siteTlsSettings == nil {
		t.Errorf("Should not have received a nil siteTlsSettings instance")
	}
	if siteTlsSettings.Mandatory != true {
		t.Errorf("Parameter Mandatory doesn't match. Actual : %v", siteTlsSettings.Mandatory)
	}
}
