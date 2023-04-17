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

	ret, err := client.GetAtoSiteMitigationConfiguration(accountId, siteId)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[Error] Error executing get ATO mitigation configuration request") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteMitigationConfiguration(&ATOSiteMitigationConfigurationDTO{})

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
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d", ATOSitePath, siteId, ATOSiteMitigationConfigurationPath, accountId)

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

	ret, err := client.GetAtoSiteMitigationConfiguration(accountId, siteId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "Error response from server for fetching ATO mitigation configuration for site") {
		t.Errorf("Error response from server for fetching ATO mitigation configuration for site, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteMitigationConfiguration(&ATOSiteMitigationConfigurationDTO{})
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
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d", ATOSitePath, siteId, ATOSiteMitigationConfigurationPath, accountId)

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

	ret, err := client.GetAtoSiteMitigationConfiguration(accountId, siteId)
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

	err = client.UpdateATOSiteMitigationConfiguration(&ATOSiteMitigationConfigurationDTO{})
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
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d", ATOSitePath, siteId, ATOSiteMitigationConfigurationPath, accountId)

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
    }, 
	{
        "endpointId": "5001",
        "lowAction": "NONE",
        "mediumAction": "CAPTCHA",
        "highAction": "TARPIT"
    }

]`))
	}))

	defer server.Close()

	// Initialize config and client
	config := &Config{APIID: apiId, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	// Fetch the mitigation configuration for the site
	response, err := client.GetAtoSiteMitigationConfigurationWithRetries(accountId, siteId)

	// Check for no value edge cases
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if response == nil {
		t.Errorf("Should have received a response")
	}

	if response.MitigationConfiguration == nil {
		t.Errorf("ATO mitigation configuration should not be nil")
	}

	// Verify that there are 2 items in the ATO mitigation configuration
	if len(response.MitigationConfiguration) != 2 {
		t.Errorf("Size of  mitigation configuration should be 2, received : %d", len(response.MitigationConfiguration))
	}

	// Use the first item for testing values
	mitigationConfigurationItem := response.MitigationConfiguration[0]

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

	// Verify that both the mitigation configuration items are not the same
	if mitigationConfigurationItem.EndpointId == response.MitigationConfiguration[1].EndpointId {
		t.Errorf("Mitigation configuration endpoint are not expected to be identical with a value of %s", mitigationConfigurationItem.EndpointId)
	}

}
