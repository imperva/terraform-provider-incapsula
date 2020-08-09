package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientConfigureAdvanceCacheBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	configureCacheHeadersResponse, err := client.ConfigureAdvanceCaching(siteID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding  Advance WAF Caching for site id %d", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if configureCacheHeadersResponse != nil {
		t.Errorf("Should have received a nil configureCacheHeadersResponse instance")
	}
}

func TestClientConfigureAdvanceCacheBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", advancedCacheEndpoint) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", advancedCacheEndpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 42
	configureCacheHeadersResponse, err := client.ConfigureAdvanceCaching(siteID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add Advances Cache JSON response for site id %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if configureCacheHeadersResponse != nil {
		t.Errorf("Should have received a nil configureCacheHeadersResponse instance")
	}
}

func TestClientConfigureAdvanceCacheInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", advancedCacheEndpoint) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", advancedCacheEndpoint, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":"0","res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 42
	configureCacheHeadersResponse, err := client.ConfigureAdvanceCaching(siteID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding Advance Caching for site id %d: %s", siteID, string(`{"site_id":"0","res":1}`))) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if configureCacheHeadersResponse != nil {
		t.Errorf("Should have received a nil configureCacheHeadersResponse instance")
	}
}

func TestClientConfigureAdvanceCacheValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", advancedCacheEndpoint) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", advancedCacheEndpoint, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":"123","res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 42
	configureCacheHeadersResponse, err := client.ConfigureAdvanceCaching(siteID, "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}

	if configureCacheHeadersResponse == nil {
		t.Errorf("Should not have received a nil configureCacheHeadersResponse instance")
	}
	if configureCacheHeadersResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}
