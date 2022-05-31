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

func TestClientAddCertificateHsmBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_hsm_test.TestClientAddCertificateHsmBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	hSMDataDTO := getFakeHsmDataDto()
	addCertificateResponse, err := client.AddHsmCertificate(siteID, "dfgdfg", &hSMDataDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error from Imperva service when adding HSM certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func getFakeHsmDataDto() HSMDataDTO {
	var hsmDetailList []HSMDetailsDTO
	hSMDetailsDTO := HSMDetailsDTO{
		KeyId:    "sdfdhrthrth",
		ApiKey:   "sdfsgdrgrg",
		HostName: "dfgrgergerg",
	}
	hsmDetailList = append(hsmDetailList, hSMDetailsDTO)
	hSMDataDTO := HSMDataDTO{
		Certificate:    "sdfsdfsdf",
		HsmDetailsList: hsmDetailList,
	}
	return hSMDataDTO
}

func TestClientAddCertificateHsmBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_certificate_hsm_test.TestClientAddCertificateHsmBadJSON")
	siteID := "1234"
	url := getHsmUrlForServerMock(siteID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != url {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", url, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	hSMDataDTO := getFakeHsmDataDto()
	addCertificateResponse, err := client.AddHsmCertificate(siteID, "bla", &hSMDataDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add HSM certificate JSON response for siteId %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func getHsmUrlForServerMock(siteID string) string {
	//url := fmt.Sprintf("/sites/%s/%s?input_hash=bla", siteID, endpointHsmCertificateAdd)
	url := fmt.Sprintf("/sites/%s/%s", siteID, endpointHsmCertificateAdd)
	return url
}

func TestClientAddCertificateHsmInvalidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_hsm_test.TestClientAddCertificateHsmInvalidRule")
	siteID := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != getHsmUrlForServerMock(siteID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", getHsmUrlForServerMock(siteID), req.URL.String())
		}
		rw.Write([]byte(`{"res":3015,"res_message":"Internal error","debug_info":{"id-info":"13008","Error":"Unexpected error occurred"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	hsmDataFakeDto := getFakeHsmDataDto()
	addCertificateResponse, err := client.AddHsmCertificate(siteID, "bla", &hsmDataFakeDto)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error adding HSM certificate- res not 0. siteId: %s", siteID)) {
		t.Errorf("Should have received a bad certificate error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateHsmValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_hsm_test.TestClientAddCertificateHsmValidRule")
	siteID := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != getHsmUrlForServerMock(siteID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", getHsmUrlForServerMock(siteID), req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	hsmDataFateDto := getFakeHsmDataDto()
	addCertificateResponse, err := client.AddHsmCertificate(siteID, "bla", &hsmDataFateDto)
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
// DeleteCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteCertificateHsmBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_hsm_test.TestClientDeleteCertificateHsmBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	err := client.DeleteHsmCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error deleting HSM certificate while sending request. siteId: %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteCertificateHsmBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_hsm_test.TestClientDeleteCertificateHsmBadJSON")
	siteID := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != getHsmUrlForServerMock(siteID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", getHsmUrlForServerMock(siteID), req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteHsmCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error deleting HSM certificate, json parse error. siteId: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteCertificateHsmInvalidSiteID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateInvalidSiteID")
	siteID := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != getHsmUrlForServerMock(siteID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", getHsmUrlForServerMock(siteID), req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"id-info":"13008","site_id":"1234"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteHsmCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error deleting HSM certificate- res not 0. siteId: %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteCertificateHsmValidSite(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG]Running test client_certificate_hsm_test.TestClientDeleteCertificateHsmValidSite")
	siteID := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != getHsmUrlForServerMock(siteID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", getHsmUrlForServerMock(siteID), req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev2: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteHsmCertificate(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
