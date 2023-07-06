package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	ATOSiteMitigationConfigurationPath = "/mitigation"
)

func TestATOSiteMitigationConfigurationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteId := 42
	accountId := 55
	endpointId := "123"

	ret, err := client.GetAtoEndpointMitigationConfiguration(accountId, siteId, endpointId)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[Error] Error executing get ATO mitigation configuration request") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteMitigationConfiguration(&ATOEndpointMitigationConfigurationDTO{})

	if err == nil {
		t.Errorf("Should have received an error")
	}

	// Site ID is not present and we should produce this error
	if !strings.Contains(err.Error(), "site_id is not specified in updating ATO Mitigation configuration") {
		t.Errorf("Should have received an client error, got: %s", err)
	}

	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestATOSiteMitigationConfigurationErrorResponse(t *testing.T) {
	apiId := "foo"
	apiKey := "bar"
	accountId := 55
	siteId := 42
	endpointId := "123"
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d&endpointIds=123", ATOSitePath, siteId, ATOSiteMitigationConfigurationPath, accountId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`Server error`))
	}))

	defer server.Close()

	config := &Config{APIID: apiId, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	ret, err := client.GetAtoEndpointMitigationConfiguration(accountId, siteId, endpointId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "Error response from server for fetching ATO mitigation configuration for site") {
		t.Errorf("Error response from server for fetching ATO mitigation configuration for site, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteMitigationConfiguration(&ATOEndpointMitigationConfigurationDTO{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "site_id is not specified in updating ATO Mitigation configuration") {
		t.Errorf("Should have received 'Error parsing ATO mitigation configuration response for site', got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestATOSiteAMitigationConfigurationInvalidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteId := 42
	accountId := 55
	endpointId := "123"
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d&endpointIds=123", ATOSitePath, siteId, ATOSiteMitigationConfigurationPath, accountId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`site-config=wrong-value`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	ret, err := client.GetAtoEndpointMitigationConfiguration(accountId, siteId, endpointId)
	if err == nil {
		t.Errorf("Should have received an error")
		return
	}
	if !strings.Contains(err.Error(), "Error in parsing JSON response for ATO mitigation configuration") {
		t.Errorf("Should have received 'Error in parsing JSON response for ATO mitigation configuration', got: %s", err)
		return
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
		return
	}

	err = client.UpdateATOSiteMitigationConfiguration(&ATOEndpointMitigationConfigurationDTO{})
	if err == nil {
		t.Errorf("Should have received an error")
		return
	}
	if !strings.Contains(err.Error(), "site_id is not specified in updating ATO Mitigation configuration") {
		t.Errorf("Should have received 'site_id is not specified in updating ATO Mitigation configuration', got: %s", err)
		return
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
		return
	}
}

func TestATOSiteMitigationConfigurationResponse(t *testing.T) {
	apiId := "foo"
	apiKey := "bar"
	accountId := 55
	siteId := 42
	endpointId := "5000"
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d&endpointIds=5000", ATOSitePath, siteId, ATOSiteMitigationConfigurationPath, accountId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`[

    {
        "endpointId": "5000",
        "lowAction": "NONE",
        "mediumAction": "CAPTCHA",
        "highAction": "BLOCK"
    }

]`))
	}))

	defer server.Close()

	// Initialize config and client
	config := &Config{APIID: apiId, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	// Fetch the mitigation configuration for the site
	mitigationConfigurationItem, err := client.GetAtoEndpointMitigationConfigurationWithRetries(accountId, siteId, endpointId)

	// Check for no value edge cases
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if mitigationConfigurationItem == nil {
		t.Errorf("Should have received a response for GetAtoEndpointMitigationConfigurationWithRetries")
	}

	if mitigationConfigurationItem.EndpointId != "5000" {
		t.Errorf("Expected mitigation configuration endpointId : 5000, received : %s", mitigationConfigurationItem.EndpointId)
	}

	if mitigationConfigurationItem.LowAction != "NONE" {
		t.Errorf("Expected LowAction 'NONE', received : %s", mitigationConfigurationItem.LowAction)
	}

	if mitigationConfigurationItem.MediumAction != "CAPTCHA" {
		t.Errorf("Expected MediumAction 'NONE', received : %s", mitigationConfigurationItem.MediumAction)
	}

	if mitigationConfigurationItem.HighAction != "BLOCK" {
		t.Errorf("Expected HighAction 'BLOCK', received : %s", mitigationConfigurationItem.HighAction)
	}

}
