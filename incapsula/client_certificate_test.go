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
// AddCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientAddCertificateBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAddCertificateBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "abc", "def", "efg", "RSA", "hij")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding custom certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_certificate_test.TestClientAddCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "", "RSA", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add custom certificate JSON response for site_id %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateInvalidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAddCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateAdd, req.URL.String())
		}
		rw.Write([]byte(`{"res":3015,"res_message":"Internal error","debug_info":{"id-info":"13008","Error":"Unexpected error occurred"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "", "RSA", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding custom certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a bad certificate error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAddCertificateValidRule")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateAdd, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "", "RSA", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addCertificateResponse == nil {
		t.Errorf("Should not have received a nil addCertificateResponse instance")
	}
	if addCertificateResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// ListCertificates Tests
////////////////////////////////////////////////////////////////

func TestClientListCertificatesBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientListCertificatesBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID, ReadHSMCustomCertificate)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting custom certificates for site_id %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientListCertificatesBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID, ReadHSMCustomCertificate)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing certificates list JSON response for site_id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesInvalidRequest(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientListCertificatesInvalidRequest")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"id-info":"13007","site_id":"1234"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID, ReadHSMCustomCertificate)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting custom certificates list for site_id %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listCertificatesResponse == nil {
		t.Errorf("Should have received a listCertificatesResponse instance")
	}
}

////////////////////////////////////////////////////////////////
// EditCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientEditCertificateBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientEditCertificateBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	certificate := "foo"
	privateKey := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, privateKey, passphrase, "RSA", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error editing custom certificate for site_id: %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

func TestClientEditCertificateBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientEditCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	privateKey := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, privateKey, passphrase, "RSA", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing edit custom certificarte JSON response for site_id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

func TestClientEditCertificateValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientEditCertificateValidRule")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	privateKey := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, privateKey, passphrase, "RSA", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if editCertificateResponse == nil {
		t.Errorf("Should not have received a nil editCertificateResponse instance")
	}

	if editCertificateResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteCertificateBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID, "RSA")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting custom certificate for site_id: %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteCertificateBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID, "RSA")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting custom certificate for site_id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteCertificateInvalidSiteID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateInvalidSiteID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"id-info":"13008","site_id":"1234"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID, "RSA")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting custom certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteCertificateValidSite(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG]Running test client_certificate_test.TestClientDeleteCertificateValidSite")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID, "RSA")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
