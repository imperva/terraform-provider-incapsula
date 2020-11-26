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
// GetAccountDataStorageRegion Tests
////////////////////////////////////////////////////////////////

func TestClientGetAccountDataStorageRegionBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := "123"
	dataStorageRegionResponse, err := client.GetAccountDataStorageRegion(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting default data storage region for account id: %s", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientGetAccountDataStorageRegionBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDataStorageRegionGet) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDataStorageRegionGet, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "123"
	dataStorageRegionResponse, err := client.GetAccountDataStorageRegion(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing default data storage region JSON response for account id: %s", accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientGetAccountDataStorageRegionInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDataStorageRegionGet) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDataStorageRegionGet, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized account_id","debug_info":{"account_id":"7289383","id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "7289383"
	dataStorageRegionResponse, err := client.GetAccountDataStorageRegion(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting default data storage region for account id: %s", accountID)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if dataStorageRegionResponse == nil {
		t.Errorf("Should have received a dataStorageRegionResponse instance")
	}
}

func TestClientGetAccountDataStorageRegionValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDataStorageRegionGet) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDataStorageRegionGet, req.URL.String())
		}
		rw.Write([]byte(`{"region":"US","res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "123"
	dataStorageRegionResponse, err := client.GetAccountDataStorageRegion(accountID)
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

func TestClientUpdateAccountDataStorageRegionBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := "42"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateAccountDataStorageRegion(accountID, region)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating data storage region with value (%s) on account_id: %s", region, accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientUpdateAccountDataStorageRegionBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDataStorageRegionUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDataStorageRegionUpdate, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "42"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateAccountDataStorageRegion(accountID, region)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update default data storage region JSON response for accountID %s", accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientUpdateAccountDataStorageRegionInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDataStorageRegionUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDataStorageRegionUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized account_id","debug_info":{"account_id":"7293873","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "7293873"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateAccountDataStorageRegion(accountID, region)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating default data storage region for accountID %s", accountID)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if dataStorageRegionResponse != nil {
		t.Errorf("Should have received a nil dataStorageRegionResponse instance")
	}
}

func TestClientUpdateAccountDataStorageRegionValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDataStorageRegionUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDataStorageRegionUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"region":"US","res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "7293873"
	region := "US"
	dataStorageRegionResponse, err := client.UpdateAccountDataStorageRegion(accountID, region)
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
