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
	certId := 100

	_, err := client.GetSiteMtlsCertificateAssociation(certId, siteID, "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error getting Site to mutual TLS Imperva to Origin Certificate association for Site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

}

func TestClientGetSiteMtlsCertificateAssociationBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	certifiateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls/origin/%d/associated-sites/%d", certifiateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	certificateExists, err := client.GetSiteMtlsCertificateAssociation(certifiateID, siteID, "")

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 406 from Incapsula service on fetching Incapsula Site to mutual TLS Imperva to Origin certificate association for Site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if certificateExists == true {
		t.Errorf("Should have received a nil mTLSCertificate instance")
	}
}

func TestClientGetSiteMtlsCertificateAssociationInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	certifiateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls/origin/%d/associated-sites/%d", certifiateID, siteID)

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
               "pointer": "/v3/mtls/origin"
           },
           "title": "Bad Request",
           "detail": "handleRequest - Got response headers:org.springframework.web.reactive.function.client.DefaultClientResponse$DefaultHeaders@20c80d50, status: 400 BAD_REQUEST, body: {\"errors\":[{\"status\":400,\"id\":\"de31602becdf6d4b\",\"source\":{\"pointer\":\"/mtls/origin\"},\"title\":\"Bad Request\",\"detail\":\"Certificate already exists\"}]}"
       }
   ]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	certificateExists, err := client.GetSiteMtlsCertificateAssociation(certifiateID, siteID, "")

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 400 from Incapsula service on fetching Incapsula Site to mutual TLS Imperva to Origin certificate association for Site ID %d", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if certificateExists == true {
		t.Errorf("Should have received a nil mTLSCertificate instance")
	}
}

func TestClientGetSiteMtlsCertificateAssociationValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	certifiateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls/origin/%d/associated-sites/%d", certifiateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	certificateExists, err := client.GetSiteMtlsCertificateAssociation(certifiateID, siteID, "")

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if certificateExists == false {
		t.Errorf("Should not have received a false in GetSiteMtlsCertificateAssociation")
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

	err := client.CreateSiteMtlsCertificateAssociation(certificateID, siteID, "")
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
               "pointer": "/v3/mtls/origin"
           },
           "title": "Bad Request",
           "detail": "handleRequest - Got response headers:org.springframework.web.reactive.function.client.DefaultClientResponse$DefaultHeaders@20c80d50, status: 400 BAD_REQUEST, body: {\"errors\":[{\"status\":400,\"id\":\"de31602becdf6d4b\",\"source\":{\"pointer\":\"/mtls/origin\"},\"title\":\"Bad Request\",\"detail\":\"Certificate already exists\"}]}"
       }
   ]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.CreateSiteMtlsCertificateAssociation(certificateID, siteID, "")

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

	err := client.DeleteSiteMtlsCertificateAssociation(certificateID, siteID, "")
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
		"pointer": "/v3/mtls/origin/111"
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

	err := client.CreateSiteMtlsCertificateAssociation(certificateID, siteID, "")

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on creating Incapsula Site to mutual TLS Imperva to Origin certificate Association for Site ID %d, Certificate ID %d", siteID, certificateID)) {
		t.Errorf("Should have received a status error for Incapsula Site to mutual TLS Imperva to Origin certificate association, got: %s", err)
	}
}
