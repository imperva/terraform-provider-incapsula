package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
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

// //////////////////////////////////////////////////////////////
// executeRequest retry Tests (UM-13362)
// //////////////////////////////////////////////////////////////

// withShortRetries temporarily shrinks the retry window so retry-exhaustion
// tests run fast, restoring the original value when the test finishes.
func withShortRetries(t *testing.T) {
	t.Helper()
	original := durationOfRetriesInSeconds
	durationOfRetriesInSeconds = 1
	t.Cleanup(func() { durationOfRetriesInSeconds = original })
}

// TestExecuteRequestRetriesExhaustedReturnsError is the regression test for the
// nil pointer crash in GetPerformanceSettings: when a "read" request keeps
// getting 502 until retries are exhausted, executeRequest must surface a
// non-nil error rather than swallowing the resource.Retry failure. The caller's
// `if err != nil` guard then prevents the nil/usable response from reaching
// `defer resp.Body.Close()`.
func TestExecuteRequestRetriesExhaustedReturnsError(t *testing.T) {
	withShortRetries(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&calls, 1)
		rw.WriteHeader(http.StatusBadGateway) // 502 on every attempt
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	req, err := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error preparing request: %s", err)
	}
	SetHeaders(client, req, contentTypeApplicationJson, ReadSitePerformance, nil)

	resp, err := client.executeRequest(req)
	if err == nil {
		t.Errorf("Should have received an error after retries were exhausted on 502")
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Errorf("Expected the request to be retried at least twice, got %d call(s)", atomic.LoadInt32(&calls))
	}
	if resp != nil {
		defer resp.Body.Close()
	}
}

// TestExecuteRequestRetriesThenSucceeds verifies the retry path still recovers:
// a 502 followed by a 200 should return the successful response and no error.
func TestExecuteRequestRetriesThenSucceeds(t *testing.T) {
	withShortRetries(t)

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.WriteHeader(http.StatusBadGateway) // first attempt fails
			return
		}
		rw.WriteHeader(http.StatusOK) // subsequent attempt succeeds
		rw.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	req, err := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error preparing request: %s", err)
	}
	SetHeaders(client, req, contentTypeApplicationJson, ReadSitePerformance, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Errorf("Should not have received an error once a retry succeeds, got: %s", err)
	}
	if resp == nil {
		t.Fatalf("Should have received a non-nil response on success")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 after successful retry, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Errorf("Expected at least one retry before success, got %d call(s)", atomic.LoadInt32(&calls))
	}
}

// TestGetPerformanceSettingsNoPanicOnPersistent502 reproduces the original
// crash scenario end-to-end: persistent 502s during a terraform import refresh
// must produce a clean error from GetPerformanceSettings, not a panic.
func TestGetPerformanceSettingsNoPanicOnPersistent502(t *testing.T) {
	withShortRetries(t)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	performanceSettings, err := client.GetPerformanceSettings("76050311")
	if err == nil {
		t.Errorf("Should have received an error on persistent 502 instead of a nil response")
	}
	if performanceSettings != nil {
		t.Errorf("Should have received a nil performanceSettings instance on failure")
	}
}
