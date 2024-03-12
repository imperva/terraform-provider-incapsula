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

func TestClientADReadIncaRulePriorityBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev3: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	category := "Test"

	addIncapRuleResponse, diags := client.ReadDeliveryRuleConfiguration(siteID, category)
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when reading Delivery Rules of category %s for Site ID %s:", category, siteID)) {
		t.Errorf("Should have received a client error, got: %s", diags[0].Detail)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestClientADReadIncaRulePriorityBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	category := "REDIRECT"
	endpoint := fmt.Sprintf("/sites/%s/delivery-rules-configuration?category=%s", siteID, category)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev3: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, diags := client.ReadDeliveryRuleConfiguration(siteID, category)
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Delivery Rules JSON response of categorie %s for Site ID %s", category, siteID)) {
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
	category := "REDIRECT"

	endpoint := fmt.Sprintf("/sites/%s/delivery-rules-configuration?category=%s", siteID, category)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"data": [{"filter":"Full-URL == \"/someurl\"","rule_name":"rule","action":"RULE_ACTION_REDIRECT","enabled":true }]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLRev3: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, err := client.ReadDeliveryRuleConfiguration(siteID, "REDIRECT")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if readIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil readIncapRuleResponse instance")
	}
	if readIncapRuleResponse != nil && readIncapRuleResponse.Errors != nil {
		t.Errorf("Should not be error response")
	}

	if readIncapRuleResponse != nil && readIncapRuleResponse.RulesList[0].Enabled != true {
		t.Errorf("Should not have received disabled rule")
	}
}

////////////////////////////////////////////////////////////////
// UpdateIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateIncapRulePriorityBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev3: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	category := "REDIRECT"

	isEnable := false
	rule := []DeliveryRuleDto{
		{
			From:    "http://www.a.com",
			To:      "http://www.b.com",
			Filter:  "",
			Enabled: isEnable,
		},
	}
	rulesList := DeliveryRulesListDTO{
		RulesList: rule,
	}

	updateIncapRuleResponse, diags := client.UpdateDeliveryRuleConfiguration(siteID, category, &rulesList)
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when updating delivery rules category %s for Site ID %s", category, siteID)) {
		t.Errorf("Should have received an client error, got: %s", diags[0].Detail)
	}
	if updateIncapRuleResponse != nil {
		t.Errorf("Should have received a nil updateIncapRuleResponse instance")
	}
}

func TestClientUpdateIncapRulePriorityBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	category := "REDIRECT"

	rule := []DeliveryRuleDto{
		{
			From:    "http://www.a.com",
			To:      "http://www.b.com",
			Filter:  "",
			Enabled: false,
		},
	}
	rulesList := DeliveryRulesListDTO{
		RulesList: rule,
	}

	endpoint := fmt.Sprintf("/sites/%s/delivery-rules-configuration?category=%s", siteID, category)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURLRev3: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	updateIncapRuleResponse, diags := client.UpdateDeliveryRuleConfiguration(siteID, category, &rulesList)
	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Delivery Rules JSON response of categorie %s for Site ID %s", category, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if updateIncapRuleResponse != nil {
		t.Errorf("Should have received a nil updateIncapRuleResponse instance")
	}
}
