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
// Verify Tests
////////////////////////////////////////////////////////////////
func TestClientVerifyBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	err := client.Verify()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error checking account") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientVerifyBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccount) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccount, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.Verify()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error parsing account JSON response") {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientVerifyInvalidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccount) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccount, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.Verify()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from Incapsula service when checking account") {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
}

func TestClientVerifyValidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccount) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccount, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.Verify()
	if err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
}

////////////////////////////////////////////////////////////////
// AddSite Tests
////////////////////////////////////////////////////////////////

func TestClientAddSiteBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	domain := "foo.com"
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "", "", "")
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
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "", "", "")
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
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "", "", "")
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
	addSiteResponse, err := client.AddSite(domain, "", "", "", "", "", "", "")
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
