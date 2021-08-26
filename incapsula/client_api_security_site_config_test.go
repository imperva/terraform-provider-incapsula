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
// UpdateApiSecuritySiteConfig Tests
////////////////////////////////////////////////////////////////

func TestUpdateApiSecuritySiteConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	payload := ApiSecuritySiteConfigPostPayload{
		ApiOnlySite: true,
		IsAutomaticDiscoveryApiIntegrationEnabled: false,
		NonApiRequestViolationAction:              "IGNORE",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	apiSecuritySiteConfigPostResponse, err := client.UpdateApiSecuritySiteConfig(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from Incapsula service while updating API security site configuration") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecuritySiteConfigPostResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestUpdateApiSecuritySiteConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s%d", siteConfigUrl, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	payload := ApiSecuritySiteConfigPostPayload{
		ApiOnlySite: true,
		IsAutomaticDiscoveryApiIntegrationEnabled: false,
		NonApiRequestViolationAction:              "IGNORE",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	apiSecuritySiteConfigPostResponse, err := client.UpdateApiSecuritySiteConfig(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing API security JSON response")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiSecuritySiteConfigPostResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigPostResponse instance")
	}

}

func TestUpdateApiSecuritySiteConfigInvalidSiteConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s%d", siteConfigUrl, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"value": "An internal error occurred. Contact support specifying your account ID and site ID.",
							"isError": true}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	//invalid payload
	payload := ApiSecuritySiteConfigPostPayload{
		ApiOnlySite: true,
	}

	apiSecuritySiteConfigPostResponse, err := client.UpdateApiSecuritySiteConfig(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from Incapsula service when updating api-security site configuration")) {
		t.Errorf("Should have received a bad api securiy site config rule error, got: %s", err)
	}
	if apiSecuritySiteConfigPostResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigPostResponse instance")
	}
}

func TestUpdateApiSecuritySiteConfigValidSiteConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s%d", siteConfigUrl, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"value": {"siteId": 42},"isError": false}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	payload := ApiSecuritySiteConfigPostPayload{
		ApiOnlySite: true,
		IsAutomaticDiscoveryApiIntegrationEnabled: false,
		NonApiRequestViolationAction:              "IGNORE",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	apiSecuritySiteConfigPostResponse, err := client.UpdateApiSecuritySiteConfig(
		siteID,
		&payload)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if apiSecuritySiteConfigPostResponse == nil {
		t.Errorf("Should not have received a nil apiSecuritySiteConfigPostResponse instance")
	}
	if apiSecuritySiteConfigPostResponse.Value.SiteId == 0 {
		t.Errorf("Should not have received an empty site config ID")
	}
}

////////////////////////////////////////////////////////////////
// ReadApiSecuritySiteConfig Tests
////////////////////////////////////////////////////////////////
func TestClientReadApiSecuritySiteConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	apiSecuritySiteConfigGetResponse, err := client.ReadApiSecuritySiteConfig(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR]Error from Incapsula service while reading Api-Security Site Config for site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecuritySiteConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestClientReadApiSecuritySiteConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s%d", siteConfigUrl, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiSecuritySiteConfigGetResponse, err := client.ReadApiSecuritySiteConfig(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing GET Api-Security Site Config JSON response for site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecuritySiteConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}

}

func TestClientReadApiSecuritySiteConfigInvalidSiteConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s%d", siteConfigUrl, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"value": "Specified site or account do not exist. Make sure you have a valid account and a valid site before attempting to execute an operation on an API",
				"isError": true}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiSecuritySiteConfigGetResponse, err := client.ReadApiSecuritySiteConfig(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from Incapsula service when reading Api-Security Site Config for site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiSecuritySiteConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestClientReadApiSecuritySiteConfigValidSiteConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s%d", siteConfigUrl, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"value": {"siteId": 123},"isError": false}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiSecuritySiteConfigGetResponse, err := client.ReadApiSecuritySiteConfig(siteID)
	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if apiSecuritySiteConfigGetResponse == nil {
		t.Errorf("Should not have received a nil apiSecuritySiteConfigGetResponse instance")
	}
	if apiSecuritySiteConfigGetResponse.Value.SiteId != 123 {
		t.Errorf("Site ID doesn't match")
	}
}
