package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestClientGetDomainsForSiteValidCase(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_configuration_test.TestClientGetDomainsForSiteValidCase")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/%s", siteID, "domains?pageSize=-1")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "data": [
        {
            "id": 12,
            "siteId": 111,
            "domain": "a.co",
            "mainDomain": true,
            "managed": true,
            "status": "BYPASSED",
            "creationDate": 1665485465000
        },
        {
            "id": 13,
            "siteId": 111,
            "domain": "b.a.co",
            "managed": true,
            "validationMethod": "CNAME",
            "status": "BYPASSED",
            "creationDate": 1668699134000
        },
        {
            "id": 14,
            "siteId": 111,
            "domain": "c.a.co",
            "managed": true,
            "validationMethod": "CNAME",
            "status": "BYPASSED",
            "creationDate": 1668700703096
        }
    ],
    "meta": {
        "totalPages": 0
    }
}`))
	}))
	defer server.Close()
	log.Printf("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	addCertificateResponse, err := client.GetWebsiteDomains(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addCertificateResponse == nil {
		t.Errorf("Should not have received a nil addCertificateResponse instance")
	}
}

func TestClientGetDomainsForSiteBadJsonResponse(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_configuration_test.TestClientGetDomainsForSiteBadJsonResponse")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/%s", siteID, "domains?pageSize=-1")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()
	log.Printf("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	addCertificateResponse, err := client.GetWebsiteDomains(siteID)

	if err == nil {
		t.Errorf("expected to get error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula get domain details %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientGetDomainsForSiteUnAuthorizedResponse(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_configuration_test.TestClientGetDomainsForSiteUnAuthorizedResponse")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	endpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/%s", siteID, "domains")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(4001)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{}`))
	}))
	defer server.Close()
	log.Printf("ENDPOINT: " + endpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	addCertificateResponse, err := client.GetWebsiteDomains(siteID)

	if err == nil {
		t.Errorf("expected to get error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when reading domain configuration details %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddDomainsForSiteValidCase(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_configuration_test.TestClientAddDomainsForSiteValidCase")
	apiID := "foo"
	apiKey := "bar"
	siteID := "111"
	bulkUpdateEndpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains", siteID)
	siteExtraDetailsEndpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains/extraDetails", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if (req.URL.String() != bulkUpdateEndpoint) && (req.URL.String() != siteExtraDetailsEndpoint) {
			t.Errorf("Should have have hit %s or %s endpoints. Got: %s", bulkUpdateEndpoint, siteExtraDetailsEndpoint, req.URL.String())
		}

		if req.URL.String() == bulkUpdateEndpoint {
			rw.Write([]byte(`{
			"data": [
				{
					"id": 12,
					"siteId": 111,
					"domain": "a.co",
					"mainDomain": true,
					"managed": true,
					"status": "BYPASSED",
					"creationDate": 1665485465000
				},
				{
					"id": 13,
					"siteId": 111,
					"domain": "b.a.co",
					"managed": true,
					"validationMethod": "CNAME",
					"status": "BYPASSED",
					"creationDate": 1668699134000
				},
				{
					"id": 14,
					"siteId": 111,
					"domain": "c.a.co",
					"managed": true,
					"validationMethod": "CNAME",
					"status": "BYPASSED",
					"creationDate": 1668700703096
				}
			],
			"meta": {
				"totalPages": 0
			}
			}`))
		}

		if req.URL.String() == siteExtraDetailsEndpoint {
			rw.Write([]byte(`{
			"data": [
				{
					"numberOfAutoDiscoveredDomains": 1,
					"maxAllowedDomains": 1000
				}
			]
			}`))
		}
	}))

	defer server.Close()
	log.Printf("ENDPOINT: " + bulkUpdateEndpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteIdStr, err := strconv.Atoi(siteID)
	addCertificateResponse, err := client.BulkUpdateDomainsToSite(siteID, getFakeSiteDomainDetailsArray(siteIdStr, "a.com", "b.com"))
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addCertificateResponse == nil {
		t.Errorf("Should not have received a nil addCertificateResponse instance")
	}
}

func TestClientAddDomainsForSiteErrorResponse(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_configuration_test.TestClientAddDomainsForSiteErrorResponse")
	apiID := "foo"
	apiKey := "bar"
	siteId := "111"
	bulkUpdateEndpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains", siteId)
	siteExtraDetailsEndpoint := fmt.Sprintf("/site-domain-manager/v2/sites/%s/domains/extraDetails", siteId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)

		if (req.URL.String() != bulkUpdateEndpoint) && (req.URL.String() != siteExtraDetailsEndpoint) {
			t.Errorf("Should have have hit %s or %s endpoints. Got: %s", bulkUpdateEndpoint, siteExtraDetailsEndpoint, req.URL.String())
		}

		rw.Write([]byte(`{
    "errors": [
        {
            "status": 400,
            "id": "f83b824fc6d9505e",
            "source": {
                "pointer": "/v2/sites/65887086/domains"
            },
            "title": "BAD_REQUEST",
            "detail": "upload failed due to the following errors: Invalid domain name:?a.com,?b.com"
        }
    ]
}`))

		if req.URL.String() == siteExtraDetailsEndpoint {
			rw.Write([]byte(`{
			"data": [
				{
					"numberOfAutoDiscoveredDomains": 1,
					"maxAllowedDomains": 1000
				}
			]
			}`))
		}
	}))

	defer server.Close()
	log.Printf("ENDPOINT: " + bulkUpdateEndpoint)
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	var domainA = "?a.com"
	var domainB = "?b.com"
	siteIdStr, err := strconv.Atoi(siteId)
	addCertificateResponse, err := client.BulkUpdateDomainsToSite(siteId, getFakeSiteDomainDetailsArray(siteIdStr, domainA, domainB))

	if err == nil {
		t.Errorf("expected to get error")
	}

	errorString := err.Error()
	if !strings.Contains(errorString, fmt.Sprintf("[ERROR] Error status code 400 from Incapsula get site domains extra details for siteId %s", siteId)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}

	if !strings.Contains(err.Error(), fmt.Sprintf("Invalid domain name:%s,%s", domainA, domainB)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}

	if addCertificateResponse != nil {
		t.Errorf("expexcted to get nil on addCertificateResponse instance")
	}
}
func getFakeSiteDomainDetailsArray(siteId int, domainA string, domainB string) []SiteDomainDetails {
	siteDomainDetails := make([]SiteDomainDetails, 2)
	siteDomainDetails[0] = SiteDomainDetails{SiteId: siteId, Domain: domainA}
	siteDomainDetails[1] = SiteDomainDetails{SiteId: siteId, Domain: domainB}
	return siteDomainDetails
}
