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
// ConfigureACLSecurityRule Tests
////////////////////////////////////////////////////////////////

func TestClientConfigureACLSecurityRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID, ruleID := 42, "42"
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding ACL for rule id %s and site id %d", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if configureACLSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureACLSecurityRuleResponse instance")
	}
}

func TestClientConfigureACLSecurityRuleBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointACLRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointACLRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID, ruleID := 42, "42"
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add ACL rule JSON response for rule id %s and site id %d", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if configureACLSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureACLSecurityRuleResponse instance")
	}
}

func TestClientConfigureACLSecurityRuleInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointACLRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointACLRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID, ruleID := 42, "42"
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when configuring ACL rule for rule id %s and site id %d: %s", ruleID, siteID, string(`{"site_id":0,"res":1}`))) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if configureACLSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureACLSecurityRuleResponse instance")
	}
}

func TestClientConfigureACLSecurityRuleContinentsValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointACLRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointACLRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID, ruleID := 42, blacklistedCountries
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "Africa", "Australia", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if configureACLSecurityRuleResponse == nil {
		t.Errorf("Should not have received a nil configureACLSecurityRuleResponse instance")
	}
	if configureACLSecurityRuleResponse.SiteID != 123 {
		t.Errorf("Site ID doesn't match")
	}
}

func TestClientConfigureACLSecurityRuleIPsValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointACLRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointACLRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID, ruleID := 42, blacklistedIPs
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "", "", "44.55.66.77", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if configureACLSecurityRuleResponse == nil {
		t.Errorf("Should not have received a nil configureACLSecurityRuleResponse instance")
	}
	if configureACLSecurityRuleResponse.SiteID != 123 {
		t.Errorf("Site ID doesn't match")
	}
}

func TestClientConfigureACLSecurityRuleURLsValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointACLRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointACLRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID, ruleID := 42, blacklistedURLs
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "", "", "", "/alpha,/bravo", "CONTAINS,EQUALS")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if configureACLSecurityRuleResponse == nil {
		t.Errorf("Should not have received a nil configureACLSecurityRuleResponse instance")
	}
	if configureACLSecurityRuleResponse.SiteID != 123 {
		t.Errorf("Site ID doesn't match")
	}
}

func TestClientConfigureACLSecurityRuleResultCodeStringValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointACLRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointACLRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"site_id":123,"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID, ruleID := 42, blacklistedCountries
	configureACLSecurityRuleResponse, err := client.ConfigureACLSecurityRule(siteID, ruleID, "Africa", "Australia", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if configureACLSecurityRuleResponse == nil {
		t.Errorf("Should not have received a nil configureACLSecurityRuleResponse instance")
	}
	if configureACLSecurityRuleResponse.SiteID != 123 {
		t.Errorf("Site ID doesn't match")
	}
}
