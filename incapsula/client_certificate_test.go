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
// AddCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientAddCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteId := "1234"
	addCertificateResponse, err := client.AddCertificate(siteId, "abc", "def", "efg")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding custom certificate for siteID %s", siteId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateBadJSON(t *testing.T) {
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
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add custom certificate JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateInvalidRule(t *testing.T) {
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
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding custom certificate for siteID %s", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateValidRule(t *testing.T) {
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
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "")
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
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting custom certificates for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesBadJSON(t *testing.T) {
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
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing cutsom certificates list JSON response for siteID: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesInvalidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting custom certificate list (site_id: %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if listCertificatesResponse == nil {
		t.Errorf("Should not have received a nil listCertificatesResponse instance")
	}

	if listCertificatesResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// EditCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientEditCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	certificate := "foo"
	private_key := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, private_key, passphrase)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error editing custom certificate for siteID: %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

func TestClientEditCertificateBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	private_key := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, private_key, passphrase)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing edit custom certificate JSON response for siteID: %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

func TestClientEditCertificateInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":0,"res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	private_key := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, private_key, passphrase)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when editing custom certificate for siteID %d", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

func TestClientEditCertificateValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":123,"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	private_key := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, private_key, passphrase)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if editCertificateResponse == nil {
		t.Errorf("Should not have received a nil editCertificateResponse instance")
	}

	if editCertificateResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting custom certificate (site_id: %s)", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteCertificateBadJSON(t *testing.T) {
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
	err := client.DeleteCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete custom certificate JSON response (site_id: %s)", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteCertificateInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting custom certificate (site_id: %s)", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteCertificateValidSite(t *testing.T) {
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
	err := client.DeleteCertificate(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
