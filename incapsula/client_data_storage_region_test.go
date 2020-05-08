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
// GetDataStorageRegion Tests
////////////////////////////////////////////////////////////////

func TestClientGetDataStorageRegionBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "123"
	dataStorageRegionResponse, err := client.GetDataStorageRegion(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting data storage region for site id: %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientGetDataStorageRegionBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataStorageRegionGet) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataStorageRegionGet, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "123"
	dataStorageRegionResponse, err := client.GetDataStorageRegion(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing site data storage region JSON response for site id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientGetDataStorageRegionInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataStorageRegionGet) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataStorageRegionGet, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"7289383","id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "7289383"
	dataStorageRegionResponse, err := client.GetDataStorageRegion(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting site data storage region for site id: %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if dataStorageRegionResponse == nil {
		t.Errorf("Should have received a dataStorageRegionResponse instance")
	}
}

func TestClientGetDataStorageRegionValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataStorageRegionGet) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataStorageRegionGet, req.URL.String())
		}
		rw.Write([]byte(`{"region":"US","res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "123"
	dataStorageRegionResponse, err := client.GetDataStorageRegion(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if dataStorageRegionResponse == nil {
		t.Errorf("Should not have received a nil dataStorageRegionResponse instance")
	}
	if dataStorageRegionResponse.Region != "US" {
		t.Errorf("Data storage region doesn't match")
	}
	if dataStorageRegionResponse.Res != 0 {
		t.Errorf("Data storage result code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// UpdateSite Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateDataStorageRegionBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateDataStorageRegion(siteID, region)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating data storage region with value (%s) on site_id: %s", region, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientUpdateDataStorageRegionBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataStorageRegionUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataStorageRegionUpdate, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateDataStorageRegion(siteID, region)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update site data storage region JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientUpdateDataStorageRegionInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataStorageRegionUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataStorageRegionUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"7293873","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "7293873"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateDataStorageRegion(siteID, region)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating site data storage region for siteID %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientUpdateDataStorageRegionValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataStorageRegionUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataStorageRegionUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"region":"US","res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "7293873"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateDataStorageRegion(siteID, region)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if dataStorageRegionResponse == nil {
		t.Errorf("Should not have received a nil dataStorageRegionResponse instance")
	}
	if dataStorageRegionResponse.Region != "US" {
		t.Errorf("Region doesn't match")
	}
	if dataStorageRegionResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}
