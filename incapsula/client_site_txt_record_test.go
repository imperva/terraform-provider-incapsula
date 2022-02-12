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
// GetTXTRecords Tests
////////////////////////////////////////////////////////////////

func TestClientGetTXTRecordsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 123
	txtRecords, err := client.ReadTXTRecords(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading TXT record(s) for siteID: %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if txtRecords != nil {
		t.Errorf("Should have received a nil txtRecords instance")
	}
}

func TestClientGetTXTRecordsBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/general/additionalTxtRecords", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	txtSettings, err := client.ReadTXTRecords(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap TXT record(s) JSON response for siteID: %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if txtSettings != nil {
		t.Errorf("Should have received a nil txtSettings instance")
	}
}

func TestClientGetTXTRecordsInvalidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/general/additionalTxtRecords", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	txtSettings, err := client.ReadTXTRecords(siteID)
	log.Print(err)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code %d from Incapsula service when reading TXT record(s) for siteID: %d", 404, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if txtSettings != nil {
		t.Errorf("Should have received a nil txtSettings instance")
	}
}

func TestClientGetTXTRecordsValidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/general/additionalTxtRecords", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123456,"txt_record_value_one":"test1","txt_record_value_two":"test2","txt_record_value_three":"test3","txt_record_value_four":"test4","txt_record_value_five":"test5"}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	txtSettings, err := client.ReadTXTRecords(siteID)
	log.Printf("Test Results: SiteID=%d, "+
		"txt_record_value_one=%s, "+
		"txt_record_value_two=%s, "+
		"txt_record_value_three=%s, "+
		"txt_record_value_four=%s, "+
		"txt_record_value_five=%s", txtSettings.SiteID, txtSettings.TxtRecordValueOne,
		txtSettings.TxtRecordValueTwo, txtSettings.TxtRecordValueThree, txtSettings.TxtRecordValueFour,
		txtSettings.TxtRecordValueFive)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if txtSettings == nil {
		t.Errorf("Should not have received a nil txtRecords instance")
	}
	if txtSettings.SiteID != 123456 {
		t.Errorf("Site ID doesn't match")
	}
	if txtSettings.TxtRecordValueOne != "test1" {
		t.Errorf("Text value one doesn't match")
	}
	if txtSettings.TxtRecordValueTwo != "test2" {
		t.Errorf("Text value two doesn't match")
	}
	if txtSettings.TxtRecordValueThree != "test3" {
		t.Errorf("Text value three doesn't match")
	}
	if txtSettings.TxtRecordValueFour != "test4" {
		t.Errorf("Text value four doesn't match")
	}
	if txtSettings.TxtRecordValueFive != "test5" {
		t.Errorf("Text value five doesn't match")
	}

}

////////////////////////////////////////////////////////////////
// UpdateTXTRecords Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateTXTRecordsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 123
	_, err := client.UpdateTXTRecord(siteID, "test1", "test2", "test3", "test4", "test5")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating TXT record(s) for siteID: %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdateTXTRecordsInvalidSite(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/general/additionalTxtRecords", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.UpdateTXTRecord(siteID, "", "test2", "test3", "test4", "test5")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when updating TXT record(s) for siteID: %d", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientUpdateTXTRecordsValidSite(t *testing.T) {
	log.Print("Start TestClientUpdateTXTRecordsValidSite")
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/general/additionalTxtRecords", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	txtResponse, err := client.UpdateTXTRecord(siteID, "test1", "test2", "test3", "test4", "test5")
	log.Print(err)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if txtResponse.Res != 0 {
		t.Errorf("Should not have received a nil txtRecords instance")
	}
	//if txtSettings.SiteID != 123456 {
	//	t.Errorf("Site ID doesn't match")
	//}
	//if txtSettings.TxtRecordValueOne != "test1" {
	//	t.Errorf("Text value one doesn't match")
	//}
	//if txtSettings.TxtRecordValueTwo != "test2" {
	//	t.Errorf("Text value two doesn't match")
	//}
	//if txtSettings.TxtRecordValueThree != "test3" {
	//	t.Errorf("Text value three doesn't match")
	//}
	//if txtSettings.TxtRecordValueFour != "test4" {
	//	t.Errorf("Text value four doesn't match")
	//}
	//if txtSettings.TxtRecordValueFive != "test5" {
	//	t.Errorf("Text value five doesn't match")
	//}
}
