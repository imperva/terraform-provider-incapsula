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
// AddSecurityRuleException Tests
////////////////////////////////////////////////////////////////

func TestClientAddSecurityRuleExceptionBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientAddSecurityRuleExceptionBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	addSecurityRuleExceptionResponse, err := client.AddSecurityRuleException(siteID, ruleID, "", "", "", "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error configuring security rule exception rule_id (%s) for site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if addSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

func TestClientAddSecurityRuleExceptionBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_security_rule_exception.TestClientAddSecurityRuleExceptionBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	addSecurityRuleExceptionResponse, err := client.AddSecurityRuleException(siteID, ruleID, "", "", "", "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing SecurityRuleExceptionCreateResponse JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientAddSecurityRuleExceptionInvalidRuleID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientAddSecurityRuleExceptionInvalidRuleID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"id-info":"13007","Unknown rule_id. The value should be selected out of the ACLs or security rules ids":"bad_rule_id"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "bad_rule_id"
	addSecurityRuleExceptionResponse, err := client.AddSecurityRuleException(siteID, ruleID, "", "", "", "AN,AS", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error - invalid security rule exception rule_id (%s)", ruleID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if addSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

func TestClientAddSecurityRuleExceptionInvalidParam(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientAddSecurityRuleExceptionInvalidRuleParam")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"Unexpected error","debug_info":{"id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	badIps := "1234"
	addSecurityRuleExceptionResponse, err := client.AddSecurityRuleException(siteID, ruleID, "", "", "", "", badIps, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing SecurityRuleExceptionCreateResponse JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if addSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

////////////////////////////////////////////////////////////////
// EditSecurityRuleException Tests
////////////////////////////////////////////////////////////////

func TestClientEditSecurityRuleExceptionBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientEditSecurityRuleExceptionBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	editSecurityRuleExceptionResponse, err := client.EditSecurityRuleException(siteID, ruleID, "", "", "", "", "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error configuring security rule exception rule_id (%s) for site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if editSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

func TestClientEditecurityRuleExceptionBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_security_rule_exception.TestClientEditecurityRuleExceptionBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	editSecurityRuleExceptionResponse, err := client.EditSecurityRuleException(siteID, ruleID, "", "", "", "", "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing configure security rule exception JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil configureWAFSecurityRuleResponse instance")
	}
}

func TestClientEditSecurityRuleExceptionInvalidRuleID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientEditSecurityRuleExceptionInvalidRuleID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"id-info":"13007","Unknown rule_id. The value should be selected out of the ACLs or security rules ids":"bad_rule_id"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "bad_rule_id"
	editSecurityRuleExceptionResponse, err := client.EditSecurityRuleException(siteID, ruleID, "", "", "", "", "", "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error - invalid security rule exception rule_id (%s)", ruleID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if editSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

func TestClientEditSecurityRuleExceptionInvalidWhitelistID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientEditSecurityRuleExceptionInvalidParam")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"id-info":"13007","Whitelist id should be a long number":"19280691s"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	badIps := "1.2.3.4,1.2.4"
	badWhitelistID := "1234"
	editSecurityRuleExceptionResponse, err := client.EditSecurityRuleException(siteID, ruleID, "", "", "", "", badIps, "", "", "", "", badWhitelistID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding security rule exception for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if editSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

func TestClientEditSecurityRuleExceptionInvalidParam(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientEditSecurityRuleExceptionInvalidParam")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"Unexpected error","debug_info":{"id-info":"13007"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	badIps := "1234"
	editSecurityRuleExceptionResponse, err := client.EditSecurityRuleException(siteID, ruleID, "", "", "", "", badIps, "", "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding security rule exception for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
	if editSecurityRuleExceptionResponse != nil {
		t.Errorf("Should have received a nil addSecurityRuleExceptionResponse instance")
	}
}

////////////////////////////////////////////////////////////////
// DeleteSecurityRuleException Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteSecurityRuleExceptionBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientDeleteSecurityRuleExceptionBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	whitelistID := "12345"
	err := client.DeleteSecurityRuleException(siteID, ruleID, whitelistID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting security rule exception whitelist_id (%s) for rule_id (%s) for site_id (%d)", whitelistID, ruleID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
}

func TestClientDeleteSecurityRuleExceptionBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_security_rule_exception.TestClientDeleteSecurityRuleExceptionBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	whitelistID := "12345"
	err := client.DeleteSecurityRuleException(siteID, ruleID, whitelistID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete security rule exception JSON response for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteSecurityRuleExceptionInvalidRuleID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientDeleteSecurityRuleExceptionInvalidRuleID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"id-info":"13008","Unknown rule_id. The value should be selected out of the ACLs or security rules ids":"api.threats.sql_injection2"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "bad_rule_id"
	whitelistID := "12345"
	err := client.DeleteSecurityRuleException(siteID, ruleID, whitelistID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting security rule exception for rule_id (%s)", ruleID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
}

func TestClientDeleteSecurityRuleExceptionInvalidWhitelistID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientDeleteSecurityRuleExceptionInvalidWhitelistID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"Whitelist id does not exist, if you're trying to add a new one leave this field empty":"192806912","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	badWhitelistID := "abc"
	err := client.DeleteSecurityRuleException(siteID, ruleID, badWhitelistID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting security rule exception for rule_id (%s) and site_id (%d)", ruleID, siteID)) {
		t.Errorf("Should have received a bad WAF security error, got: %s", err)
	}
}

func TestClientDeleteSecurityRuleExceptionValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_security_rule_exception.TestClientDeleteSecurityRuleExceptionValidRule")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointExceptionConfigure) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointExceptionConfigure, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 1234
	ruleID := "api.threats.backdoor"
	badWhitelistID := "abc"
	err := client.DeleteSecurityRuleException(siteID, ruleID, badWhitelistID)
	if err != nil {
		t.Errorf("Should have received an error")
	}
}
