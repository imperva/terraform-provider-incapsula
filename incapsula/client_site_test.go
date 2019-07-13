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
// AddSite Tests
////////////////////////////////////////////////////////////////

func TestClientAddSiteBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	domain := "foo.com"
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding site for domain %s", domain)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addSiteResponse != nil {
		t.Errorf("Should have received a nil addSiteResponse instance")
	}
}

func TestClientAddSiteBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add site JSON response for domain %s", domain)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addSiteResponse != nil {
		t.Errorf("Should have received a nil addSiteResponse instance")
	}
}

func TestClientAddSiteInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteAdd, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding site for domain %s", domain)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if addSiteResponse != nil {
		t.Errorf("Should have received a nil addSiteResponse instance")
	}
}

func TestClientAddSiteValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteAdd, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addSiteResponse == nil {
		t.Errorf("Should not have received a nil addSiteResponse instance")
	}
	if addSiteResponse.SiteID != 123 {
		t.Errorf("Site ID doesn't match")
	}
	if addSiteResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// SiteStatus Tests
////////////////////////////////////////////////////////////////

func TestClientSiteStatusBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	domain := "foo.com"
	siteID := 123
	siteStatusResponse, err := client.SiteStatus(domain, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting site status for domain %s (site id: %d)", domain, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteStatusResponse != nil {
		t.Errorf("Should have received a nil siteStatusResponse instance")
	}
}

func TestClientSiteStatusBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteStatus) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteStatus, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	siteID := 123
	siteStatusResponse, err := client.SiteStatus(domain, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing site status JSON response for domain %s (site id: %d)", domain, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if siteStatusResponse != nil {
		t.Errorf("Should have received a nil siteStatusResponse instance")
	}
}

func TestClientSiteStatusInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteStatus) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteStatus, req.URL.String())
		}
		rw.Write([]byte(`{"site_creation_date":0, "dns":[], "res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	siteID := 123
	siteStatusResponse, err := client.SiteStatus(domain, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting site status for domain %s (site id: %d)", domain, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if siteStatusResponse != nil {
		t.Errorf("Should have received a nil siteStatusResponse instance")
	}
}

func TestClientSiteStatusValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteStatus) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteStatus, req.URL.String())
		}
		rw.Write([]byte(`{"site_creation_date":1527885500000, "dns":[], "res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	siteID := 123
	siteStatusResponse, err := client.SiteStatus(domain, siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if siteStatusResponse == nil {
		t.Errorf("Should not have received a nil siteStatusResponse instance")
	}
	if siteStatusResponse.SiteCreationDate != 1527885500000 {
		t.Errorf("Site creation date doesn't match")
	}
	if len(siteStatusResponse.DNS) != 0 {
		t.Errorf("DNS records are not empty")
	}
	if siteStatusResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// UpdateSite Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateSiteBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	param := "active"
	value := "bypass"
	updateSiteResponse, err := client.UpdateSite(siteID, param, value)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating param (%s) with value (%s) on site_id: %s", param, value, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updateSiteResponse != nil {
		t.Errorf("Should have received a nil updateSiteResponse instance")
	}
}

func TestClientUpdateSiteBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteUpdate, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	updateSiteResponse, err := client.UpdateSite(siteID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update site JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if updateSiteResponse != nil {
		t.Errorf("Should have received a nil updateSiteResponse instance")
	}
}

func TestClientUpdateSiteInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	updateSiteResponse, err := client.UpdateSite(siteID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating site for siteID %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if updateSiteResponse != nil {
		t.Errorf("Should have received a nil updateSiteResponse instance")
	}
}

func TestClientUpdateSiteValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	addSiteResponse, err := client.UpdateSite(siteID, "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addSiteResponse == nil {
		t.Errorf("Should not have received a nil addSiteResponse instance")
	}
	if addSiteResponse.SiteID != 123 {
		t.Errorf("Site ID doesn't match")
	}
	if addSiteResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteSite Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteSiteBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	domain := "foo.com"
	siteID := 123
	err := client.DeleteSite(domain, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting site for domain %s (site id: %d)", domain, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteSiteBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	siteID := 123
	err := client.DeleteSite(domain, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete site JSON response for domain %s (site id: %d)", domain, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteSiteInvalidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	siteID := 123
	err := client.DeleteSite(domain, siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting site for domain %s (site id: %d)", domain, siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteSiteValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSiteDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSiteDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	domain := "foo.com"
	siteID := 123
	err := client.DeleteSite(domain, siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
