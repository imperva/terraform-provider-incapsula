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
// GetSiteMtlsClientToImpervaCertificateAssociation Tests
////////////////////////////////////////////////////////////////

func TestClientGetSiteMtlsClientToImpervaCertificateAssociationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	accountID := 88
	certificateID := 100

	_, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error getting Site to mutual TLS Imperva to Origin Certificate association for Site ID %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

}

func TestClientGetSiteMtlsClientToImpervaCertificateAssociationBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 88
	certificateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls-origin/certificates/%d/associated-sites/%d", certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	clientCaCertificateWithSites, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("ERROR] Error parsing Incapsula Site to mutual TLS Client to Imperva Certificate association JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if clientCaCertificateWithSites != nil {
		t.Errorf("Should have received a nil clientCaCertificateWithSites instance")
	}
}

func TestClientGetSiteMtlsClientToImpervaCertificateAssociationInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 88
	certificateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls-origin/certificates/%d/associated-sites/%d", certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
  "errors": [
	{
		"status": 500,
		"id": "cca667c1371c31ff",
		"source": {
		"pointer": "/v3/mtls-origin/certificates/111"
	},
		"title": "Internal Server Error",
		"detail": "Internal Server Error"
	}":[123,222,42]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	clientCaCertificateWithSites, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on fetching Incapsula Site to mutual TLS Client to Imperva Certificate association for Site ID %s", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if clientCaCertificateWithSites != nil {
		t.Errorf("Should have received a nil clientCaCertificateWithSites instance")
	}
}

func TestClientGetSiteMtlsClientToImpervaCertificateAssociationValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 88
	certificateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls-origin/certificates/%d/associated-sites/%d", certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
   {
  "id": 1,
  "name": "some name",
  "serialNumber": "string",
  "issuer": "string",
  "creationDate": "string",
  "assignedSites":[123,222,42]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	clientCaCertificateWithSites, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if clientCaCertificateWithSites.Name == "some name" {
		t.Errorf("Should not have received a false in GetSiteMtlsCertificateAssociation")
	}
}


////////////////////////////////////////////////////////////////
// CreateSiteMtlsClientToImpervaCertificateAssociation Tests
////////////////////////////////////////////////////////////////

func TestClientCreateSiteMtlsClientToImpervaCertificateAssociationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	accountID := 88
	certificateID := 100

	_, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error creating Incapsula Site to mutual TLS Client to Imperva Certificate Association for certificate ID %d, Site ID %d", certificateID,siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

}

func TestClientCreateSiteMtlsClientToImpervaCertificateAssociationBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 88
	certificateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls-origin/certificates/%d/associated-sites/%d", certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	clientCaCertificateWithSites, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("ERROR] Error parsing Incapsula Site to mutual TLS Client to Imperva Certificate association JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if clientCaCertificateWithSites != nil {
		t.Errorf("Should have received a nil clientCaCertificateWithSites instance")
	}
}

func TestClientCreateSiteMtlsClientToImpervaCertificateAssociationInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 88
	certificateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls-origin/certificates/%d/associated-sites/%d", certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
  "errors": [
	{
		"status": 500,
		"id": "cca667c1371c31ff",
		"source": {
		"pointer": "/v3/mtls-origin/certificates/111"
	},
		"title": "Internal Server Error",
		"detail": "Internal Server Error"
	}":[123,222,42]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	clientCaCertificateWithSites, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on fetching Incapsula Site to mutual TLS Client to Imperva Certificate association for Site ID %s", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if clientCaCertificateWithSites != nil {
		t.Errorf("Should have received a nil clientCaCertificateWithSites instance")
	}
}

func TestClientCreateSiteMtlsClientToImpervaCertificateAssociationnValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 88
	certificateID := 100

	endpoint := fmt.Sprintf("/certificates-ui/v3/mtls-origin/certificates/%d/associated-sites/%d", certificateID, siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
   {
  "id": 1,
  "name": "some name",
  "serialNumber": "string",
  "issuer": "string",
  "creationDate": "string",
  "assignedSites":[123,222,42]
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	clientCaCertificateWithSites, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if clientCaCertificateWithSites.Name == "some name" {
		t.Errorf("Should not have received a false in GetSiteMtlsCertificateAssociation")
	}
}