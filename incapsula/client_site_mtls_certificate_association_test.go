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
// GetSiteMtlsCertificateAssociation Tests
////////////////////////////////////////////////////////////////

func TestClientGetSiteMtlsCertificateAssociationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	mTLSCertificate, err := client.GetSiteMtlsCertificateAssociation(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprint("[ERROR] Error getting Incapsula Site to mTLS certificate association for Site ID 42")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSCertificate != nil {
		t.Errorf("Should have received a nil mTLSCertificate instance")
	}
}

func TestClientGetSiteMtlsCertificateAssociationBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s?siteId=%d", endpointMTLSCertificate, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificate, err := client.GetSiteMtlsCertificateAssociation(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing Incapsula Site to mutual TLS Imperva to Origin Certificate association JSON response for Site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if mTLSCertificate != nil {
		t.Errorf("Should have received a nil mTLSCertificate instance")
	}
}

func TestClientGetSiteMtlsCertificateAssociationInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s?siteId=%d", endpointMTLSCertificate, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
    "errors": [
        {
            "status": 400,
            "id": "16d37a3dfb2b3aff",
            "source": {
                "pointer": "/v3/mtls-origin/certificates"
            },
            "title": "Bad Request",
            "detail": "handleRequest - Got response headers:org.springframework.web.reactive.function.client.DefaultClientResponse$DefaultHeaders@20c80d50, status: 400 BAD_REQUEST, body: {\"errors\":[{\"status\":400,\"id\":\"de31602becdf6d4b\",\"source\":{\"pointer\":\"/mtls-origin/certificates\"},\"title\":\"Bad Request\",\"detail\":\"Certificate already exists\"}]}"
        }
    ]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificate, err := client.GetSiteMtlsCertificateAssociation(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 400 from Incapsula service on fetching Incapsula Site to mutual TLS Imperva to Origin certificate association for Site ID %d", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if mTLSCertificate != nil {
		t.Errorf("Should have received a nil mTLSCertificate instance")
	}
}

func TestClientGetSiteMtlsCertificateAssociationValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("%s?siteId=%d", endpointMTLSCertificate, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "data": [
        {
            "certificateId": 1,
            "issuedTo": "company",
            "issuedBy": "ca company",
            "validFrom": 1655011098999,
            "validUntil": 168456887638629,
            "chain": "ca cert",
            "name": "mtls-cert",
            "creationDate": 1685036052000,
            "lastUpdate": 1685036052000,
            "serialNumber": "4345667",
            "appliedSitesDetails": [
                {
                    "externalId": 122222,
                    "name": "site.incapsula.co",
                    "accountId": 2121221
                }
            ]
        }
    ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificate, err := client.GetSiteMtlsCertificateAssociation(siteID)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if mTLSCertificate == nil {
		t.Errorf("Should not have received a nil addIncapRuleResponse instance")
	}
	if mTLSCertificate.Id != 1 {
		t.Errorf("Certifcate ID doesn't match. Actual : %v", mTLSCertificate.Id)
	}
}

////////////////////////////////////////////////////////////////
// UpdateSiteMtlsCertificateAssociation Tests
////////////////////////////////////////////////////////////////

func TestClientCreateSiteMtlsCertificateAssociationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	certificateID := 100

	err := client.CreateSiteMtlsCertificateAssociation(certificateID, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error creating Incapsula Site to Imperva to Origin mutual TLS Certificate Association for certificate ID %d, Site ID %d", certificateID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdateSiteMtlsCertificateAssociationInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	certificateID := 100

	endpoint := fmt.Sprintf("%s/%d/associated-sites/%d", endpointMTLSCertificate, certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
    "errors": [
        {
            "status": 400,
            "id": "16d37a3dfb2b3aff",
            "source": {
                "pointer": "/v3/mtls-origin/certificates"
            },
            "title": "Bad Request",
            "detail": "handleRequest - Got response headers:org.springframework.web.reactive.function.client.DefaultClientResponse$DefaultHeaders@20c80d50, status: 400 BAD_REQUEST, body: {\"errors\":[{\"status\":400,\"id\":\"de31602becdf6d4b\",\"source\":{\"pointer\":\"/mtls-origin/certificates\"},\"title\":\"Bad Request\",\"detail\":\"Certificate already exists\"}]}"
        }
    ]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.CreateSiteMtlsCertificateAssociation(certificateID, siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 400 from Incapsula service on creating Incapsula Site to mutual TLS Imperva to Origin certificate Association for Site ID %d, Certificate ID %d", siteID, certificateID)) {
		t.Errorf("Should have received a bad Incapsula Site to mutual TLS Imperva to Origin certificate association error, got: %s", err)
	}
}

////////////////////////////////////////////////////////////////
// DeleteSiteMtlsCertificateAssociation Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteSiteMtlsCertificateAssociationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	certificateID := 100

	err := client.DeleteSiteMtlsCertificateAssociation(certificateID, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error deleting Incapsula Site to Imperva to Origin mutual TLS Certificate Association for certificate ID %d for Site ID %d", certificateID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
}

func TestClientDeleteSiteMtlsCertificateAssociationInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	certificateID := 100

	endpoint := fmt.Sprintf("%s/%d/associated-sites/%d", endpointMTLSCertificate, certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
    {
			"errors": [
	{
		"status": 500,
		"id": "cca667c1371c31ff",
		"source": {
		"pointer": "/v3/mtls-origin/certificates/111"
	},
		"title": "Internal Server Error",
		"detail": "Internal Server Error"
	}
]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.CreateSiteMtlsCertificateAssociation(certificateID, siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on creating Incapsula Site to mutual TLS Imperva to Origin certificate Association for Site ID %d, Certificate ID %d", siteID, certificateID)) {
		t.Errorf("Should have received a status error for Incapsula Site to mutual TLS Imperva to Origin certificate association, got: %s", err)
	}
}
