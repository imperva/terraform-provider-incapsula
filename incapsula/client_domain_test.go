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
	log.Printf("ENDPOINT: " + endpoint)
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
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domain_configuration_test.TestClientGetDomainsForSiteValidCase")
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
	log.Printf("ENDPOINT: " + endpoint)
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
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domain_configuration_test.TestClientGetDomainsForSiteValidCase")
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
	log.Printf("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.AddDomainToSite(siteID, "a.co")

	if err.Error() != "add domain request failed: The provided domain is not valid." {
		t.Errorf("Should have received an error")
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
