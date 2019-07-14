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
// ConfigureWAFSecurityRule Tests
////////////////////////////////////////////////////////////////

func TestClientConfigureWAFSecurityRuleBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	security_rule_action := "badRuleAction"
	// siteID, ruleID, security_rule_action, activation_mode, ddos_traffic_threshold, block_bad_bots, challenge_suspected_bots
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, security_rule_action, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error configuring WAF security rule rule_id (%s) for site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	security_rule_action := "badRuleAction"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, security_rule_action, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing configure WAF rule JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleInvalidRuleID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleInvalidRuleID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":13003,"res_message":"api.response.message.13003","debug_info":{"id-info":"13007","Unknown security rule id":"bad_rule_id"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "bad_rule_id"
	security_rule_action := "bad_rule_action"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, security_rule_action, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error - invalid WAF security rule rule_id (%s)", ruleID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleInvalidRuleAction(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleInvalidRuleAction")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":13004,"res_message":"api.response.message.13004","debug_info":{"id-info":"13007","Unknown security rule action":"bad_rule_action"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	security_rule_action := "bad_rule_action"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, security_rule_action, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleInvalidRule_activationMode(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleInvalidRuleAction")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"Unexpected error","debug_info":{"id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	activation_mode := "api.threats.ddos.activation_mode.on"
	ddos_traffic_threshold := "123"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, "", activation_mode, ddos_traffic_threshold, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleInvalidRule_ddosThreshold(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleInvalidRuleAction")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"The provided DDoS Traffic Threshold is not allowed":"123","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	activation_mode := "api.threats.ddos.activation_mode.on"
	ddos_traffic_threshold := "123"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, "", activation_mode, ddos_traffic_threshold, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleInvalidRule_blockBadBots(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleInvalidRuleAction")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":6001,"res_message":"Invalid configuration parameter name","debug_info":{"id-info":"13007","Param value is invalid":"123"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	challenge_suspected_bots := "true"
	block_bad_bots := "123"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, "", "", "", block_bad_bots, challenge_suspected_bots)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientConfigureWAFSecurityRuleInvalidRule_challengeSuspectedBots(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleInvalidRuleAction")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":6001,"res_message":"Invalid configuration parameter name","debug_info":{"id-info":"13007","Param value is invalid":"123"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	challenge_suspected_bots := "123"
	block_bad_bots := "true"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, "", "", "", block_bad_bots, challenge_suspected_bots)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding WAF rule for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if configureWAFSecurityRuleResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

////////////////////////////////////////////////////////////////
// ConfigureWAFSecurityRule Tests for each ruleID
////////////////////////////////////////////////////////////////

func TestClientConfigureWAFSecurityRuleValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_waf_security_rule.TestClientConfigureWAFSecurityRuleValidRule")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointWAFRuleConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointWAFRuleConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	security_rule_action := "api.threats.action.quarantine_url"
	configureWAFSecurityRuleResponse, err := client.ConfigureWAFSecurityRule(siteID, ruleID, security_rule_action, "", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if configureWAFSecurityRuleResponse == nil {
		t.Errorf("Should not have received a nil configureWAFSecurityRuleResponse instance")
	}
	if configureWAFSecurityRuleResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}
