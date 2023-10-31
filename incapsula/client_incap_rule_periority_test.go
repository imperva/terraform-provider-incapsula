package incapsula

import (
	"fmt"
	"github.com/terraform-providers/terraform-provider-incapsula/utils"
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

	addIncapRuleResponse, _, err := client.ReadIncapRulePriorities(siteID, 8)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding Incap Rule for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
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

	readIncapRuleResponse, _, err := client.ReadIncapRulePriorities(siteID, 60)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing JSON response for Site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if readIncapRuleResponse != nil {
		t.Errorf("Should have received a nil readIncapRuleResponse instance")
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
	rule := []utils.RuleDetails{
		{
			From:    "http://www.a.com",
			To:      "http://www.b.com",
			Filter:  "",
			Enabled: &isEnable,
		},
	}

	updateIncapRuleResponse, err := client.UpdateIncapRulePriorities(siteID, utils.REWRITE, rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
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
	rule := []utils.RuleDetails{
		{From: "http://www.a.com",
			To:      "http://www.b.com",
			Filter:  "",
			Enabled: &isEnable,
		},
	}

	endpoint := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", siteID, utils.REWRITE.String())

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	updateIncapRuleResponse, err := client.UpdateIncapRulePriorities(siteID, utils.REWRITE, rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap Rule %d JSON response for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if updateIncapRuleResponse != nil {
		t.Errorf("Should have received a nil updateIncapRuleResponse instance")
	}
}

// todo
//fix json body
//add more test thats check the son body.
