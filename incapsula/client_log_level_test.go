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
// UpdateLogLevel Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateLogLevelBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	logLevel := "full"
	logsAccountId := "123"
	err := client.UpdateLogLevel(siteID, logLevel, logsAccountId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating log level (%s) on site_id: %s", logLevel, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdateLogLevelBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteLogLevel) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteLogLevel, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	logLevel := "full"
	logsAccountId := "123"
	err := client.UpdateLogLevel(siteID, logLevel, logsAccountId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update log level JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientUpdateLogLevelInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteLogLevel) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteLogLevel, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"1","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	logLevel := "full"
	logsAccountId := "123"
	err := client.UpdateLogLevel(siteID, logLevel, logsAccountId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating log level for siteID %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientUpdateLogLevelValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteLogLevel) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteLogLevel, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"log_level":"full","id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	logLevel := "full"
	logsAccountId := "123"
	err := client.UpdateLogLevel(siteID, logLevel, logsAccountId)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
