package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCspSiteDomainBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	domain, err := client.getCspDomainData(siteID, "ref")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API for when getting domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if domain != nil {
		t.Errorf("Should have received a nil response")
	}

	preApproved, err := client.getCspPreApprovedDomains(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API for when getting pre-approved domains list") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if preApproved != nil {
		t.Errorf("Should have received a nil response")
	}

	updatedDom, err := client.updateCspPreApprovedDomain(siteID, &CspPreApprovedDomain{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API while updating pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updatedDom != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.deleteCspPreApprovedDomains(siteID, "ref")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API for when deleting pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestCspSiteDomainErrorResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		rw.Write([]byte(`Server error`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	domain, err := client.getCspDomainData(siteID, "ref")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error status code 500 from CSP API when getting domain data") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if domain != nil {
		t.Errorf("Should have received a nil response")
	}

	preApproved, err := client.getCspPreApprovedDomains(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error status code 500 from CSP API when getting pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if preApproved != nil {
		t.Errorf("Should have received a nil response")
	}

	updatedDom, err := client.updateCspPreApprovedDomain(siteID, &CspPreApprovedDomain{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error status code 500 from CSP API when updating pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updatedDom != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.deleteCspPreApprovedDomains(siteID, "ref")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error status code 500 from CSP API when deleting pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestCspSiteDomainInvalidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			rw.WriteHeader(200)
		} else if req.Method == http.MethodPost {
			rw.WriteHeader(201)
		}
		rw.Write([]byte(`[not a json`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	domain, err := client.getCspDomainData(siteID, "ref")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error parsing JSON response") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if domain != nil {
		t.Errorf("Should have received a nil response")
	}

	preApproved, err := client.getCspPreApprovedDomains(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error parsing JSON response") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if preApproved != nil {
		t.Errorf("Should have received a nil response")
	}

	updatedDom, err := client.updateCspPreApprovedDomain(siteID, &CspPreApprovedDomain{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error parsing JSON response") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updatedDom != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestCspSiteDomainDataResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d/domains/%s", CspSiteApiPath, siteID, "Z29vZ2xlLmNvbQ")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"id": "Z29vZ2xlLmNvbQ",
			"domain": "google.com",
			"status": {
				"blocked": false,
				"reviewed": true
			},
			"domainRisk": "Low",
			"notes": [],
			"timeBucket": 1644624000001,
			"significance": 0,
			"resourceTypes": [
				"SCRIPT"
			],
			"browserStats": {
				"Chrome": 1
			},
			"countryStats": {
				"US": 1
			},
			"ipsSample": [
				"1.1.1.1",
				"2.2.2.2"
			],
			"sources": 1,
			"discoveredAt": 1645054932433,
			"lastSeenMs": 1645054932433,
			"domainInfo": {
				"baseDomain": "google.com",
				"companyName": null,
				"domainCategory": "Unclassified",
				"countries": [],
				"sslCertificateInfo": null,
				"registrationTime": null,
				"registrar": null,
				"orgOwner": null,
				"dynamicDnsBased": false,
				"domainQuality": {
					"score": 99.0,
					"scoreFromMl": 0.0,
					"scoreFromHeuristics": 99.0,
					"scoreOverride": null
				},
				"additionalInsights": [],
				"domainCategorySemrush": "{\"Internet & Telecom/Web Apps & Online Tools\":\"0.892291\",\"Internet & Telecom/Search Engines\":\"0.883972\",\"Internet & Telecom\":\"0.880098\",\"Internet & Telecom/Web Services\":\"0.791557\",\"Internet & Telecom/Web Services/Search Engine Optimization & Marketing\":\"0.772449\",\"Computers & Electronics\":\"0.694975\",\"Computers & Electronics/Software\":\"0.670513\",\"Internet & Telecom/Web Services/Web Stats & Analytics\":\"0.619495\",\"Online Communities\":\"0.466804\",\"Reference\":\"0.449689\"}"
			},
			"domainReports": [
				{
					"documentUri": "https://test.mage.com/test/source/document",
					"sourceFile": "https://test.mage.com/testSourceFile.js",
					"blockedUri": "https://google.com",
					"lineNumber": 0,
					"sourceType": "SCRIPT"
				}
			],
			"partOfProfile": true,
			"frequent": false
		}`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	domain, err := client.getCspDomainData(siteID, "Z29vZ2xlLmNvbQ")
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if domain == nil {
		t.Errorf("Should have received a response")
	}
	if domain.Domain != "google.com" || domain.Frequent != false || domain.PartOfProfile != true || len(domain.IPSamples) != 2 ||
		domain.Status.Blocked != false || domain.Status.Reviewed != true || domain.DomainRisk != "Low" {
		t.Errorf("Incorrect value inresponse from getCspDomainData")
	}
}

func TestCspSiteDomainPreApprovedResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d/preapprovedlist", CspSiteApiPath, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`[
			{
				"domain": "domain.com",
				"subdomains": true,
				"referenceId": "ZG9tYWluLmNvbQ"
			},
			{
				"domain": "google.com",
				"subdomains": false,
				"referenceId": "Z29vZ2xlLmNvbQ"
			}
		]`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	domains, err := client.getCspPreApprovedDomains(siteID)
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if domains == nil {
		t.Errorf("Should have received a response")
	}
	if len(domains) != 2 {
		t.Errorf("Incorrect value inresponse from getCspPreApprovedDomains")
	}
	if _, ok := domains["ZG9tYWluLmNvbQ"]; ok != true {
		t.Errorf("Incorrect value inresponse from getCspPreApprovedDomains")
	}
	if _, ok := domains["Z29vZ2xlLmNvbQ"]; ok != true {
		t.Errorf("Incorrect value inresponse from getCspPreApprovedDomains")
	}
	if domains["ZG9tYWluLmNvbQ"].Domain != "domain.com" || domains["ZG9tYWluLmNvbQ"].Subdomains != true {
		t.Errorf("Incorrect value inresponse from getCspPreApprovedDomains")
	}

	if domains["Z29vZ2xlLmNvbQ"].Domain != "google.com" || domains["Z29vZ2xlLmNvbQ"].Subdomains != false {
		t.Errorf("Incorrect value inresponse from getCspPreApprovedDomains")
	}
}

func TestCspSiteDomainPreApprovedUpdateResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d/preapprovedlist", CspSiteApiPath, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(201)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"domain": "domain.com",
			"subdomains": true,
			"referenceId": "ZG9tYWluLmNvbQ"
		}`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	domain, err := client.updateCspPreApprovedDomain(siteID, &CspPreApprovedDomain{})
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if domain == nil {
		t.Errorf("Should have received a response")
	}
	if domain.Domain != "domain.com" {
		t.Errorf("Incorrect value inresponse from updateCspPreApprovedDomain")
	}
	if domain.Subdomains != true {
		t.Errorf("Incorrect value inresponse from updateCspPreApprovedDomain")
	}
}
