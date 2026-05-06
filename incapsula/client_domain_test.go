package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGetDomainForSiteValidCase(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domain_configuration_test.TestClientGetDomainsForSiteValidCase")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	domainID := "12"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains/%s", siteID, domainID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
            "id": 12,
            "siteId": 111,
            "domain": "a.co",
            "mainDomain": true,
            "managed": true,
            "status": "BYPASSED",
            "creationDate": 1665485465000
}`))
	}))
	defer server.Close()
	log.Print("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	addDomainResponse, err := client.GetDomain(siteID, domainID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDomainResponse == nil {
		t.Errorf("Should not have received a nil addDomainResponse instance")
	}

	verifyResponse(t, addDomainResponse)
}
func TestClientAddDomainForSiteValidCase(t *testing.T) {
	log.Print("======================== BEGIN TEST ========================")
	log.Print("[DEBUG] Running test client_site_domain_configuration_test.TestClientGetDomainsForSiteValidCase")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
            "id": 12,
            "siteId": 111,
            "domain": "a.co",
            "mainDomain": true,
            "managed": true,
            "status": "BYPASSED",
            "creationDate": 1665485465000
}`))
	}))
	defer server.Close()
	log.Print("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	addDomainResponse, err := client.AddDomainToSite(siteID, "a.co")

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDomainResponse == nil {
		t.Errorf("Should not have received a nil addDomainResponse instance")
	}

	verifyResponse(t, addDomainResponse)
}

func TestClientAddDomainForSiteInvalidCase(t *testing.T) {
	log.Print("======================== BEGIN TEST ========================")
	log.Print("[DEBUG] Running test client_site_domain_configuration_test.TestClientGetDomainsForSiteValidCase")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  		"domain": "example.com",
		"errors": [
			{
			  "id": "error-1",
			  "status": 400,
			  "title": "Bad Request",
			  "detail": "The provided domain is not valid.",
			  "source": {
				"pointer": "/data/attributes/domain",
				"parameter": "domain"
			  }
			}
		  ]
		}`))
	}))
	defer server.Close()
	log.Print("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.AddDomainToSite(siteID, "a.co")

	if err.Error() != "add domain request failed (status 400): The provided domain is not valid." {
		t.Errorf("Should have received an error")
	}
}

func TestClientAddWildcardDomainAsPrimaryReturnsDetailedError(t *testing.T) {
	log.Print("======================== BEGIN TEST ========================")
	log.Print("[DEBUG] Running test TestClientAddWildcardDomainAsPrimaryReturnsDetailedError")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)

		if req.URL.String() != endpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"errors": [
				{
					"id": "error-1",
					"status": 400,
					"title": "create.wildcard.domains.for.not.allowed",
					"detail": "Creating wildcard domains for this site is not allowed",
					"source": {
						"pointer": "/data/attributes/domain",
						"parameter": "domain"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.AddDomainToSite(siteID, "*.example.com")

	if err == nil {
		t.Errorf("Should have received an error when adding wildcard as primary domain")
	}
	expectedErr := "add domain request failed (status 400): Creating wildcard domains for this site is not allowed"
	if err.Error() != expectedErr {
		t.Errorf("Expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestClientAddWildcardDomainAfterPrimarySucceeds(t *testing.T) {
	log.Print("======================== BEGIN TEST ========================")
	log.Print("[DEBUG] Running test TestClientAddWildcardDomainAfterPrimarySucceeds")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"id": 13,
			"siteId": 111,
			"domain": "*.example.com",
			"mainDomain": false,
			"managed": true,
			"status": "BYPASSED",
			"creationDate": 1665485465000
		}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.AddDomainToSite(siteID, "*.example.com")

	if err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
	if resp == nil {
		t.Errorf("Should not have received a nil response")
	}
	if resp.Domain != "*.example.com" {
		t.Errorf("Expected domain *.example.com, got: %s", resp.Domain)
	}
	if resp.MainDomain != false {
		t.Errorf("Expected MainDomain false for wildcard, got: %v", resp.MainDomain)
	}
}

func verifyResponse(t *testing.T, addDomainResponse *SiteDomainDetails) {
	if addDomainResponse.Domain != "a.co" {
		t.Errorf("Should have received a.co domain. Got: %s", addDomainResponse.Domain)
	}

	if addDomainResponse.MainDomain != true {
		t.Errorf("Should have received true for MainDomain. Got: %v", addDomainResponse.MainDomain)
	}

	if addDomainResponse.Managed != true {
		t.Errorf("Should have received true for Managed. Got: %v", addDomainResponse.Managed)
	}

	if addDomainResponse.Status != "BYPASSED" {
		t.Errorf("Should have received BYPASSED for Status. Got: %s", addDomainResponse.Status)
	}

	if addDomainResponse.CreationDate != 1665485465000 {
		t.Errorf("Should have received 1665485465000 for CreationDate. Got: %d", addDomainResponse.CreationDate)
	}

	if addDomainResponse.Id != 12 {
		t.Errorf("Should have received 12 for Id. Got: %d", addDomainResponse.Id)
	}
}
