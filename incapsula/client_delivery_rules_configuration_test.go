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
// ReadIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientAReadIncaRulePriorityBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"

	addIncapRuleResponse, _, diags := client.ReadIncapRulePriorities(siteID, "Test")
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when adding Incap Rule for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", diags[0].Detail)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}
func TestClientAReadIncaRulePriorityBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	endpoint := fmt.Sprintf("/sites/%s/delivery-rules-configuration?category=%s", siteID, "REDIRECT")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev3: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, _, diags := client.ReadIncapRulePriorities(siteID, "REDIRECT")
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if readIncapRuleResponse != nil {
		t.Errorf("Should have received a nil readIncapRuleResponse instance")
	}
}

func TestClient_ReadIncapRulePrioritiesValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/delivery-rules-configuration?category=%s", siteID, "REDIRECT")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`"{data": [{"filter":"Full-URL == \"/someurl\"","id":290109,"rule_name":"rule","action":"RULE_ACTION_REDIRECT","enabled":true }]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, statusCode, err := client.ReadIncapRulePriorities(siteID, "REDIRECT")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if readIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil readIncapRuleResponse instance")
	}
	if statusCode != 200 {
		t.Errorf("Should not have received a 200 status code")
	}

	if readIncapRuleResponse.RuleDetails[0].Enabled != true {
		t.Errorf("Should not have received disabled rule")
	}
}

////////////////////////////////////////////////////////////////
// UpdateIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateIncapRulePriorityBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev3: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 62

	isEnable := false
	rule := []DeliveryRuleDto{
		{
			From:    "http://www.a.com",
			To:      "http://www.b.com",
			Filter:  "",
			Enabled: isEnable,
		},
	}

	updateIncapRuleResponse, err := client.UpdateIncapRulePriorities(siteID, "REWRITE", rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err[0].Detail, fmt.Sprintf("Error from Incapsula service when updating Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err[0].Detail)
	}
	if updateIncapRuleResponse != nil {
		t.Errorf("Should have received a nil updateIncapRuleResponse instance")
	}
}

func TestClientUpdateIncapRulePriorityBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 62

	isEnable := false
	rule := []DeliveryRuleDto{
		{
			From:    "http://www.a.com",
			To:      "http://www.b.com",
			Filter:  "",
			Enabled: isEnable,
		},
	}

	endpoint := fmt.Sprintf("/sites/%s/delivery-rules-configuration?category=%s", siteID, "REWRITE")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	updateIncapRuleResponse, diags := client.UpdateIncapRulePriorities(siteID, "REWRITE", rule)
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Incap Rule %d JSON response for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if updateIncapRuleResponse != nil {
		t.Errorf("Should have received a nil updateIncapRuleResponse instance")
	}
}
