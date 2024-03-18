package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////
// CreateCertificateSigningRequest Tests
////////////////////////////////////////////////////////////////

func TestClientCreateCertificateSigningRequestBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_signing_request_test.TestClientCreateCertificateSigningRequestBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	domain := "sandwich.au"
	email := "test@sandwich.au"
	country := "AU"
	state := "QLD"
	city := "BNE"
	organization := "Tacos Pty Ltd"
	organizationUnit := "Sales"
	certificateSigningRequestCreateResponse, err := client.CreateCertificateSigningRequest(siteID, domain, email, country, state, city, organization, organizationUnit)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when creating certificate signing request for site_id %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if certificateSigningRequestCreateResponse != nil {
		t.Errorf("Should have received a nil certificateSigningRequestCreateResponse instance")
	}
}

func TestClientCreateCertificateSigningRequestBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_certificate_signing_request_test.TestClientCreateCertificateSigningRequestBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateSigningRequestCreate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateSigningRequestCreate, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	domain := "sandwich.au"
	email := "test@sandwich.au"
	country := "AU"
	state := "QLD"
	city := "BNE"
	organization := "Tacos Pty Ltd"
	organizationUnit := "Sales"
	certificateSigningRequestCreateResponse, err := client.CreateCertificateSigningRequest(siteID, domain, email, country, state, city, organization, organizationUnit)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing create certificate signing request JSON response for site_id %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if certificateSigningRequestCreateResponse != nil {
		t.Errorf("Should have received a nil certificateSigningRequestCreateResponse instance")
	}
}
