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
// AddIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientAddIncapRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	dcID := "43"
	addIncapRuleResponse, err := client.AddIncapRule("", "", "", "", siteID, "", "", dcID, "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding incap rule for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestClientAddIncapRuleBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	dcID := "43"
	addIncapRuleResponse, err := client.AddIncapRule("", "", "", "", siteID, "", "", dcID, "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add incap rule JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestClientAddIncapRuleInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"0","res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	dcID := "43"
	addIncapRuleResponse, err := client.AddIncapRule("", "", "", "", siteID, "", "", dcID, "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding incap rule for siteID %s: %s", siteID, `{"rule_id":"0","res":"1"}`)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestClientAddIncapRuleValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"123","res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	dcID := "43"
	addIncapRuleResponse, err := client.AddIncapRule("", "", "", "", siteID, "", "", dcID, "", "", "", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil addIncapRuleResponse instance")
	}
	if addIncapRuleResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// ListIncapRules Tests
////////////////////////////////////////////////////////////////

func TestClientListIncapRulesBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	includeAdRules := "true"
	includeIncapRules := "true"
	listIncapRulesResponse, err := client.ListIncapRules(siteID, includeAdRules, includeIncapRules)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting incap rules (include_ad_rules: %s, include_incap_rules: %s)", includeAdRules, includeIncapRules)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if listIncapRulesResponse != nil {
		t.Errorf("Should have received a nil listIncapRulesResponse instance")
	}
}

func TestClientListIncapRulesBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleList, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	includeAdRules := "true"
	includeIncapRules := "true"
	listIncapRulesResponse, err := client.ListIncapRules(siteID, includeAdRules, includeIncapRules)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing incap rule list JSON response (include_ad_rules: %s, include_incap_rules: %s)", includeAdRules, includeIncapRules)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listIncapRulesResponse != nil {
		t.Errorf("Should have received a nil listIncapRulesResponse instance")
	}
}

func TestClientListIncapRulesInvalidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	includeAdRules := "true"
	includeIncapRules := "true"
	listIncapRulesResponse, err := client.ListIncapRules(siteID, includeAdRules, includeIncapRules)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting incap rule list (include_ad_rules: %s, include_incap_rules: %s)", includeAdRules, includeIncapRules)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listIncapRulesResponse != nil {
		t.Errorf("Should have received a nil listIncapRulesResponse instance")
	}
}

func TestClientListIncapRulesValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	includeAdRules := "true"
	includeIncapRules := "true"
	listIncapRulesResponse, err := client.ListIncapRules(siteID, includeAdRules, includeIncapRules)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if listIncapRulesResponse == nil {
		t.Errorf("Should not have received a nil listIncapRulesResponse instance")
	}

	if listIncapRulesResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// EditIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientEditIncapRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	name := "foo"
	siteID, ruleID := 42, 42
	editIncapRuleResponse, err := client.EditIncapRule(siteID, "", "", name, "", "", ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error editing incap rule name: %s for siteID: %d: ", name, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if editIncapRuleResponse != nil {
		t.Errorf("Should have received a nil editIncapRuleResponse instance")
	}
}

func TestClientEditIncapRuleBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	name := "foo"
	siteID, ruleID := 42, 42
	editIncapRuleResponse, err := client.EditIncapRule(siteID, "", "", name, "", "", ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing edit incap rule JSON response for siteID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editIncapRuleResponse != nil {
		t.Errorf("Should have received a nil editIncapRuleResponse instance")
	}
}

func TestClientEditIncapRuleInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":0,"res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	name := "foo"
	siteID, ruleID := 42, 42
	editIncapRuleResponse, err := client.EditIncapRule(siteID, "", "", name, "", "", ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when editing incap rule for siteID %d, ruleID: %d", siteID, ruleID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if editIncapRuleResponse != nil {
		t.Errorf("Should have received a nil editIncapRuleResponse instance")
	}
}

func TestClientEditIncapRuleValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":123,"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	name := "foo"
	siteID, ruleID := 42, 42
	editIncapRuleResponse, err := client.EditIncapRule(siteID, "", "", name, "", "", ruleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if editIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil editIncapRuleResponse instance")
	}

	if editIncapRuleResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteIncapRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	ruleID := "42"
	accountID := "43"
	err := client.DeleteIncapRule(ruleID, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting incap rule (rule_id: %s)", ruleID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteIncapRuleBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	ruleID := "42"
	accountID := "43"
	err := client.DeleteIncapRule(ruleID, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete incap rule JSON response (rule_id: %s)", ruleID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteIncapRuleInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	ruleID := "42"
	accountID := "43"
	err := client.DeleteIncapRule(ruleID, accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting incap rule (rule_id: %s)", ruleID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteIncapRuleValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointIncapRuleDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointIncapRuleDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	ruleID := "42"
	accountID := "43"
	err := client.DeleteIncapRule(ruleID, accountID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
