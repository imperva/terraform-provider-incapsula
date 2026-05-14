package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// //////////////////////////////////////////////////////////////
// Verify Tests
// //////////////////////////////////////////////////////////////
func TestClientVerifyBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	_, err := client.Verify()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error checking account") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientVerifyBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountVerify) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountVerify, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.Verify()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error parsing account JSON response") {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientVerifyInvalidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountVerify) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountVerify, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.Verify()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error from Incapsula service when checking account") {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
}

func TestClientVerifyValidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountVerify) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountVerify, req.URL.String())
		}
		rw.Write([]byte(`{"account_type":"Reseller Customer","account_id":52219722,"parent_id":51632845,"account_name":"test account","plan_id":"ent100","plan_name":"ENTERPRISE","res":0,"res_message":"OK","debug_info":{"id-info":"999999"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.Verify()
	if err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
}
