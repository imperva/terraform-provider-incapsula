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
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", IncapRuleBaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	addIncapRuleResponse, err := client.AddIncapRule(siteID, &rule)
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

func TestClientAddIncapRuleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/rules?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	addIncapRuleResponse, err := client.AddIncapRule(siteID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap Rule JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestClientAddIncapRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/rules?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	rule := IncapRule{
		Name: "some_name",
	}

	addIncapRuleResponse, err := client.AddIncapRule(siteID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 406 from Incapsula service when adding Incap Rule for Site ID %s", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if addIncapRuleResponse != nil {
		t.Errorf("Should have received a nil addIncapRuleResponse instance")
	}
}

func TestClientAddIncapRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/rules?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"filter":"Full-URL == \"/someurl\"","rule_id":290109,"name":"myfirstcoolrule","action":"RULE_ACTION_ALERT","enabled":true}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	addIncapRuleResponse, err := client.AddIncapRule(siteID, &rule)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil addIncapRuleResponse instance")
	}
	if addIncapRuleResponse.RuleID == 0 {
		t.Errorf("Should not have received an empty rule ID")
	}
}

////////////////////////////////////////////////////////////////
// ReadIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientReadIncapRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", IncapRuleBaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 62

	readIncapRuleResponse, _, err := client.ReadIncapRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if readIncapRuleResponse != nil {
		t.Errorf("Should have received a nil readIncapRuleResponse instance")
	}
}

func TestClientReadIncapRuleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 62

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, _, err := client.ReadIncapRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap Rule %d JSON response for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if readIncapRuleResponse != nil {
		t.Errorf("Should have received a nil readIncapRuleResponse instance")
	}
}

func TestClientReadIncapRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 29010333

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"ruleId":"invalid rule : 29010333","id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, statusCode, err := client.ReadIncapRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when reading Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if statusCode != 404 {
		t.Errorf("Should have received a 404 status code")
	}
	if readIncapRuleResponse != nil {
		t.Errorf("Should have received a nil readIncapRuleResponse instance")
	}
}

func TestClientReadIncapRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 290109

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"filter":"Full-URL == \"/someurl\"","rule_id":290109,"name":"myfirstcoolrule","action":"RULE_ACTION_ALERT","enabled":true}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readIncapRuleResponse, statusCode, err := client.ReadIncapRule(siteID, ruleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if readIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil readIncapRuleResponse instance")
	}
	if statusCode != 200 {
		t.Errorf("Should not have received a 200 status code")
	}
	if readIncapRuleResponse.RuleID == 0 {
		t.Errorf("Should not have received an empty rule ID")
	}
}

////////////////////////////////////////////////////////////////
// UpdateIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateIncapRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", IncapRuleBaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 62

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	updateIncapRuleResponse, err := client.UpdateIncapRule(siteID, ruleID, &rule)
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

func TestClientUpdateIncapRuleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 62

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	updateIncapRuleResponse, err := client.UpdateIncapRule(siteID, ruleID, &rule)
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

func TestClientUpdateIncapRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 11111

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"ruleId":"invalid rule : 11111","id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	updateIncapRuleResponse, err := client.UpdateIncapRule(siteID, ruleID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when updating Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if updateIncapRuleResponse != nil {
		t.Errorf("Should have received a nil updateIncapRuleResponse instance")
	}
}

func TestClientUpdateIncapRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 290109

	rule := IncapRule{
		Name:   "myfirstcoolrule",
		Action: "RULE_ACTION_ALERT",
		Filter: "Full-URL == \"/someurl\"",
	}

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"filter":"Full-URL == \"/someurl\"","rule_id":290109,"name":"myfirstcoolrule","action":"RULE_ACTION_ALERT","enabled":true}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	updateIncapRuleResponse, err := client.UpdateIncapRule(siteID, ruleID, &rule)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if updateIncapRuleResponse == nil {
		t.Errorf("Should not have received a nil updateIncapRuleResponse instance")
	}
	if updateIncapRuleResponse.RuleID == 0 {
		t.Errorf("Should not have received an empty rule ID")
	}
}

////////////////////////////////////////////////////////////////
// DeleteIncapRule Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteIncapRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", IncapRuleBaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 62

	err := client.DeleteIncapRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteIncapRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 11111

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteIncapRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when deleting Incap Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
}

func TestClientDeleteIncapRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 290109

	endpoint := fmt.Sprintf("/sites/%s/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, IncapRuleBaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteIncapRule(siteID, ruleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
