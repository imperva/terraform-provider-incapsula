package incapsula

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////
// AddDataCenter Tests
////////////////////////////////////////////////////////////////

func TestClientAddDataCentersConfigurationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	requestDTO := DataCentersConfigurationDTO{}
	responseDTO, err := client.PutDataCentersConfiguration(siteID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error executing update Data Centers configuration request for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

/*
func TestClientAddDataCentersConfigurationBadJSON(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/v3/sites/%s/data-centers-configurations", siteID) {
			t.Errorf("Should have have hit /v3/sites/%s/data-centers-configurations endpoint. Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := DataCentersConfigurationDTO{}
	responseDTO, err := client.PutDataCentersConfiguration(siteID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update Data Centers configuration JSON request for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientAddDataCenterInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"0","res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	addDataCenterResponse, err := client.AddDataCenter(siteID, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center for siteID %s", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
	}
}

func TestClientAddDataCenterValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"123","res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	addDataCenterResponse, err := client.AddDataCenter(siteID, "", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDataCenterResponse == nil {
		t.Errorf("Should not have received a nil addDataCenterResponse instance")
	}
	if addDataCenterResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// ListDataCenters Tests
////////////////////////////////////////////////////////////////

func TestClientListDataCentersBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting data centers for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCentersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterList, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing data centers list JSON response for siteID: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCentersInvalidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting data centers list (site_id: %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listDataCentersResponse == nil {
		t.Errorf("Should have received a listDataCentersResponse instance")
	}
}

func TestClientListDataCentersValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if listDataCentersResponse == nil {
		t.Errorf("Should not have received a nil listDataCentersResponse instance")
	}

	if listDataCentersResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
} */
