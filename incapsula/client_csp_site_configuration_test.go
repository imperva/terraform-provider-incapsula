package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCspSiteConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	ret, err := client.GetCSPSite(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API for when reading site") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	ret, err = client.UpdateCSPSite(siteID, &CSPSiteConfig{})

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API while updating site configuration") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestCSPSiteConfigErrorResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d", CSPSiteApiPath, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`Server error`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	ret, err := client.GetCSPSite(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from CSP API when reading site config for ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	ret, err = client.UpdateCSPSite(siteID, &CSPSiteConfig{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from CSP API when reading site config for ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestCSPSiteConfigInvalidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d", CSPSiteApiPath, siteID)

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

	ret, err := client.GetCSPSite(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing JSON response for site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	ret, err = client.UpdateCSPSite(siteID, &CSPSiteConfig{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing JSON response for site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestCSPSiteConfigResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d", CSPSiteApiPath, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"name": "mage-website.abp-monsters.com",
			"mode": "monitor",
			"discovery": "start",
			"settings": {
				"emails": [
					{
						"email": "email@imperva.com"
					}
				]
			}}`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	ret, err := client.GetCSPSite(siteID)
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if ret == nil {
		t.Errorf("Should have received a response")
	}
	if ret.Name != "mage-website.abp-monsters.com" {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
	if ret.Mode != "monitor" {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
	if len(ret.Settings.Emails) != 1 {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
	if ret.Settings.Emails[0].Email != "email@imperva.com" {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}

	ret, err = client.UpdateCSPSite(siteID, &CSPSiteConfig{})
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if ret == nil {
		t.Errorf("Should have received a response")
	}
	if ret.Name != "mage-website.abp-monsters.com" {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
	if ret.Mode != "monitor" {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
	if len(ret.Settings.Emails) != 1 {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
	if ret.Settings.Emails[0].Email != "email@imperva.com" {
		t.Errorf("Incorrect value inresponse from GetCSPSite")
	}
}
