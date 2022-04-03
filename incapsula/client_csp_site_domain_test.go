package incapsula

import (
	"encoding/base64"
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

	updatedDom, err := client.updateCSPPreApprovedDomain(siteID, &CSPPreApprovedDomain{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from CSP API while updating pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updatedDom != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.deleteCSPPreApprovedDomains(siteID, "ref")
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

	updatedDom, err := client.updateCSPPreApprovedDomain(siteID, &CSPPreApprovedDomain{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error status code 500 from CSP API when updating pre-approved domain") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updatedDom != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.deleteCSPPreApprovedDomains(siteID, "ref")
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

	updatedDom, err := client.updateCSPPreApprovedDomain(siteID, &CSPPreApprovedDomain{})
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

func TestCSPSiteDomainPreApprovedUpdateResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s/%d/preapprovedlist", CSPSiteApiPath, siteID)

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

	domain, err := client.updateCSPPreApprovedDomain(siteID, &CSPPreApprovedDomain{})
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if domain == nil {
		t.Errorf("Should have received a response")
	}
	if domain.Domain != "domain.com" {
		t.Errorf("Incorrect value inresponse from updateCSPPreApprovedDomain")
	}
	if domain.Subdomains != true {
		t.Errorf("Incorrect value inresponse from updateCSPPreApprovedDomain")
	}
}

func TestCSPSiteDomainNotesResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	domain := "google.com"
	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))
	endpoint := fmt.Sprintf("%s/%d/domains/%s/notes", CSPSiteApiPath, siteID, domainRef)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`[
			{
				"text": "its google",
				"author": "Amiran Chachashvili (Amiran.Chachashvili@imperva.com)",
				"date": 1646804283517
			}
		]`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	notes, err := client.getCSPDomainNotes(siteID, domain)
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if notes == nil {
		t.Errorf("Should have received a response")
	}
	if len(notes) != 1 {
		t.Errorf("Incorrect value inresponse from getCSPDomainNotes")
	}
	if notes[0].Text != "its google" {
		t.Errorf("Incorrect value inresponse from getCSPDomainNotes")
	}
}

func TestCSPSiteDomainStatusResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	domain := "google.com"
	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))
	endpoint := fmt.Sprintf("%s/%d/domains/%s/status", CSPSiteApiPath, siteID, domainRef)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"blocked": false,
			"reviewed": true,
			"reviewedAt": 1646810654947
		}`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	notes, err := client.getCSPDomainStatus(siteID, domain)
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if notes == nil {
		t.Errorf("Should have received a response")
	}
}

func TestCSPSiteDomainStatusEmptyResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	domain := "google.com"
	domainRef := base64.RawURLEncoding.EncodeToString([]byte(domain))
	endpoint := fmt.Sprintf("%s/%d/domains/%s/status", CSPSiteApiPath, siteID, domainRef)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{}`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	notes, err := client.getCSPDomainStatus(siteID, domain)
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if notes == nil {
		t.Errorf("Should have received a response")
	}
}
