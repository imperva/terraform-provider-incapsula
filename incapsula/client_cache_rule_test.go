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
// AddCacheRule Tests
////////////////////////////////////////////////////////////////

func TestClientAddCacheRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	addCacheRuleResponse, err := client.AddCacheRule(siteID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding Cache Rule for Site ID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addCacheRuleResponse != nil {
		t.Errorf("Should have received a nil addCacheRuleResponse instance")
	}
}

func TestClientAddCacheRuleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	addCacheRuleResponse, err := client.AddCacheRule(siteID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Cache Rule JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addCacheRuleResponse != nil {
		t.Errorf("Should have received a nil addCacheRuleResponse instance")
	}
}

func TestClientAddCacheRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"action":"Unknown action","id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	rule := CacheRule{
		Name: "myfirstcoolrule",
	}

	addCacheRuleResponse, err := client.AddCacheRule(siteID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 406 from Incapsula service when adding Cache Rule for Site ID %s", siteID)) {
		t.Errorf("Should have received a bad cache rule error, got: %s", err)
	}
	if addCacheRuleResponse != nil {
		t.Errorf("Should have received a nil addCacheRuleResponse instance")
	}
}

func TestClientAddCacheRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules?api_id=%s&api_key=%s", siteID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":66770,"res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	addCacheRuleResponse, err := client.AddCacheRule(siteID, &rule)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addCacheRuleResponse == nil {
		t.Errorf("Should not have received a nil addCacheRuleResponse instance")
	}
	if addCacheRuleResponse.RuleID == 0 {
		t.Errorf("Should not have received an empty rule ID")
	}
}

////////////////////////////////////////////////////////////////
// ReadCacheRule Tests
////////////////////////////////////////////////////////////////

func TestClientReadCacheRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 62

	readCacheRuleResponse, _, err := client.ReadCacheRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when reading Cache Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if readCacheRuleResponse != nil {
		t.Errorf("Should have received a nil readCacheRuleResponse instance")
	}
}

func TestClientReadCacheRuleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 62

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readCacheRuleResponse, _, err := client.ReadCacheRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Cache Rule %d JSON response for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if readCacheRuleResponse != nil {
		t.Errorf("Should have received a nil readCacheRuleResponse instance")
	}
}

func TestClientReadCacheRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 29010333

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"ruleId":"invalid rule : 29010333","id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readCacheRuleResponse, statusCode, err := client.ReadCacheRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when reading Cache Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if statusCode != 404 {
		t.Errorf("Should have received a 404 status code")
	}
	if readCacheRuleResponse != nil {
		t.Errorf("Should have received a nil readCacheRuleResponse instance")
	}
}

func TestClientReadCacheRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 66772

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":66772,"action":"HTTP_CACHE_MAKE_STATIC","enabled":true,"filter":"isMobile == Yes","name":"test","ttl":300}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readCacheRuleResponse, statusCode, err := client.ReadCacheRule(siteID, ruleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if readCacheRuleResponse == nil {
		t.Errorf("Should not have received a nil readCacheRuleResponse instance")
	}
	if statusCode != 200 {
		t.Errorf("Should not have received a 200 status code")
	}
	if readCacheRuleResponse.RuleID == 0 {
		t.Errorf("Should not have received an empty rule ID")
	}
}

////////////////////////////////////////////////////////////////
// UpdateCacheRule Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateCacheRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 66772

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	err := client.UpdateCacheRule(siteID, ruleID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating Cache Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdateCacheRuleBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 66772

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.UpdateCacheRule(siteID, ruleID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Cache Rule %d JSON response for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientUpdateCacheRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 11111

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2002,"res_message":"Object is not found","debug_info":{"rule_id":"rule with id 11111 not found","id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.UpdateCacheRule(siteID, ruleID, &rule)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when updating Cache Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
}

func TestClientUpdateCacheRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 66772

	rule := CacheRule{
		Name:    "myfirstcoolrule",
		Action:  "HTTP_CACHE_MAKE_STATIC",
		Filter:  "Full-URL == \"/someurl\"",
		Enabled: true,
	}

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.UpdateCacheRule(siteID, ruleID, &rule)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}

////////////////////////////////////////////////////////////////
// DeleteCacheRule Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteCacheRuleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	ruleID := 62

	err := client.DeleteCacheRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting Cache Rule %d for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteCacheRuleInvalidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 11111

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2002,"res_message":"Object is not found","debug_info":{"rule_id":"rule with id 11111 not found","id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteCacheRule(siteID, ruleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting Cache Rule %d JSON response for Site ID %s", ruleID, siteID)) {
		t.Errorf("Should have received a bad cache rule error, got: %s", err)
	}
}

func TestClientDeleteCacheRuleValidRule(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	ruleID := 290109

	endpoint := fmt.Sprintf("/sites/%s/settings/cache/rules/%d?api_id=%s&api_key=%s", siteID, ruleID, apiID, apiKey)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13017"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteCacheRule(siteID, ruleID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
