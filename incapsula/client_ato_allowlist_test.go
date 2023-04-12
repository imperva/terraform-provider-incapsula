package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

const (
	ATOSitePath      = "/v2/sites"
	ATOAllowlistPath = "/allowlist"
)

const atoSiteAllowlistResourceType = "incapsula_ato_site_allowlist"
const atoSiteAllowlistConfigName = "testacc-terraform-ato-site-allowlist"
const atoSiteAllowlistConfigResource = atoSiteAllowlistResourceType + "." + atoSiteAllowlistConfigName

func TestATOSiteAllowlistConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteId := 42
	accountId := 55

	ret, err := client.GetAtoSiteAllowlist(accountId, siteId)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[Error] Error executing get ATO allowlist request") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteAllowlist(&AtoAllowlistDTO{})

	if err == nil {
		t.Errorf("Should have received an error")
	}

	// Site ID is not present and we should produce this error
	if !strings.Contains(err.Error(), "site_id is not specified in updating ATO allowlist") {
		t.Errorf("Should have received an client error, got: %s", err)
	}

	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestATOSiteAllowlistConfigErrorResponse(t *testing.T) {
	apiId := "foo"
	apiKey := "bar"
	accountId := 55
	siteId := 42
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d", ATOSitePath, siteId, ATOAllowlistPath, accountId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`Server error`))
	}))

	defer server.Close()

	config := &Config{APIID: apiId, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	ret, err := client.GetAtoSiteAllowlist(accountId, siteId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("Error parsing ATO allowlist response for site with ID: %d", siteId)) {
		t.Errorf("Should have received 'Error parsing ATO allowlist response for site with ID: %d', got: %s", siteId, err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteAllowlist(&AtoAllowlistDTO{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "site_id is not specified in updating ATO allowlist") {
		t.Errorf("Should have received 'site_id is not specified in updating ATO allowlist', got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestATOSiteAllowlistConfigInvalidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteId := 42
	accountId := 55
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d", ATOSitePath, siteId, ATOAllowlistPath, accountId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`site-config=wrong-value`))
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	ret, err := client.GetAtoSiteAllowlist(accountId, siteId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("Error parsing ATO allowlist response for site with ID: %d", siteId)) {
		t.Errorf("Should have received 'Error parsing ATO allowlist response for site with ID: %d', got: %s", siteId, err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}

	err = client.UpdateATOSiteAllowlist(&AtoAllowlistDTO{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "site_id is not specified in updating ATO allowlist") {
		t.Errorf("Should have received 'site_id is not specified in updating ATO allowlist', got: %s", err)
	}
	if ret != nil {
		t.Errorf("Should have received a nil response")
	}
}

func TestATOSiteAllowlistConfigResponse(t *testing.T) {
	apiId := "foo"
	apiKey := "bar"
	accountId := 55
	siteId := 42
	endpoint := fmt.Sprintf("%s/%d%s?caid=%d", ATOSitePath, siteId, ATOAllowlistPath, accountId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`[
    {
        "Ip": "192.10.20.0",
        "Mask": "24",
        "Desc": "Test IP 1",
        "Updated": 1632530998076
    }, {
        "Ip": "192.10.20.1",
        "Mask": "8",
        "Desc": "Test IP 2",
        "Updated": 1632530998077
    }
]`))
	}))

	defer server.Close()

	// Initialize config and client
	config := &Config{APIID: apiId, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	// Fetch the allowlist for the site
	response, err := client.GetAtoSiteAllowlist(accountId, siteId)

	// Check for no value edge cases
	if err != nil {
		t.Errorf("Should have not received an error")
	}
	if response == nil {
		t.Errorf("Should have received a response")
	}

	if response.allowlist == nil {
		t.Errorf("Allowlist should not be nil")
	}

	// Verify that there are 2 items in the allowlist
	if len(response.allowlist) != 2 {
		t.Errorf("Size of Allowlist should be 2, received : %d", len(response.allowlist))
	}

	// Use the first item for testing values
	allowlistItem := response.allowlist[0]

	if allowlistItem.Ip != "192.10.20.0" {
		t.Errorf("Expected allowlist IP : 192.10.20.0, received : %s", allowlistItem.Ip)
	}

	if allowlistItem.Mask != "24" {
		t.Errorf("Expected allowlist Mask to be 24, received : %s", allowlistItem.Mask)
	}

	if allowlistItem.Desc != "Test IP 1" {
		t.Errorf("Expected allowlist description to be 'Test IP 1', received : %s", allowlistItem.Desc)
	}

	if allowlistItem.Updated != 1632530998076 {
		t.Errorf("Expected allowlist Updated at time to be 1632530998076, received : %d", allowlistItem.Updated)
	}

	// Verify that both the allowlist items are not the same
	if allowlistItem.Ip == response.allowlist[1].Ip {
		t.Errorf("Allowlist IPs are not expected to be identical with a value of %s", allowlistItem.Ip)
	}
	if allowlistItem.Mask == response.allowlist[1].Mask {
		t.Errorf("Allowlist Mask are not expected to be identical with a value of %s", allowlistItem.Mask)
	}
	if allowlistItem.Desc == response.allowlist[1].Desc {
		t.Errorf("Allowlist descriptions are not expected to be identical with a value of %s", allowlistItem.Ip)
	}
	if allowlistItem.Ip == response.allowlist[1].Ip {
		t.Errorf("Allowlist IPs are not expected to be identical with a value of %s", allowlistItem.Ip)
	}

}
