package incapsula

import (
	"fmt"
	"io"
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
// executeRequest retry Tests
// //////////////////////////////////////////////////////////////

// withShortRetries overrides package-level retry vars for fast tests and returns a restore function.
func withShortRetries() func() {
	origMax := maxRetries
	origMin := retryWaitMinSeconds
	origMaxWait := retryWaitMaxSeconds
	maxRetries = 1
	retryWaitMinSeconds = 0
	retryWaitMaxSeconds = 0
	return func() {
		maxRetries = origMax
		retryWaitMinSeconds = origMin
		retryWaitMaxSeconds = origMaxWait
	}
}

// newTestClient creates a Client with minimal retry settings for fast tests.
func newTestClient(serverURL string) *Client {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: serverURL}
	return &Client{
		config:     config,
		httpClient: &http.Client{},
	}
}

// TestExecuteRequestRetriesExhaustedReturnsResponse verifies that when a "read"
// request keeps getting 502 until retries are exhausted, executeRequest returns
// the last response (resp, nil) so callers can handle the status code themselves.
func TestExecuteRequestRetriesExhaustedReturnsResponse(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&calls, 1)
		rw.WriteHeader(http.StatusBadGateway) // 502 on every attempt
	}))
	defer server.Close()

	client := newTestClient(server.URL)

	req, err := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("unexpected error preparing request: %s", err)
	}
	SetHeaders(client, req, contentTypeApplicationJson, ReadSitePerformance, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Errorf("Should not have received an error (response is returned to caller), got: %s", err)
	}
	if resp == nil {
		t.Fatal("Expected non-nil response after retries exhausted")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected status 502 in last response, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 { // 1 initial + 1 retry
		t.Errorf("Expected 2 total calls (1 + 1 retry), got %d", atomic.LoadInt32(&calls))
	}
}

// TestExecuteRequestRetriesThenSucceeds verifies the retry path still recovers:
// a 502 followed by a 200 should return the successful response and no error.
func TestExecuteRequestRetriesThenSucceeds(t *testing.T) {
	restore := withShortRetries()
	defer restore()

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

	client := newTestClient(server.URL)

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
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 total calls (1 fail + 1 success), got %d", atomic.LoadInt32(&calls))
	}
}

// TestGetPerformanceSettingsNoPanicOnPersistent502 reproduces the original
// crash scenario end-to-end: persistent 502s during a terraform import refresh
// must produce a clean error from GetPerformanceSettings, not a panic.
func TestGetPerformanceSettingsNoPanicOnPersistent502(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := NewClient(config)
	client.httpClient = &http.Client{}

	performanceSettings, err := client.GetPerformanceSettings("76050311")
	if err == nil {
		t.Errorf("Should have received an error on persistent 502 instead of a nil response")
	}
	if performanceSettings != nil {
		t.Errorf("Should have received a nil performanceSettings instance on failure")
	}
}

// //////////////////////////////////////////////////////////////
// Retry on transient failures (HTML, 503, 504, 429)
// //////////////////////////////////////////////////////////////

func TestRetryOn503ReadOperation(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls (1 fail + 1 success), got %d", atomic.LoadInt32(&calls))
	}
}

func TestRetryOn504ReadOperation(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.WriteHeader(http.StatusGatewayTimeout)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls, got %d", atomic.LoadInt32(&calls))
	}
}

func TestRetryOn429(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.WriteHeader(http.StatusTooManyRequests)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodPost, server.URL, []byte(`{"data":"test"}`))
	SetHeaders(client, req, contentTypeApplicationJson, CreateSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls (429 retried even on write), got %d", atomic.LoadInt32(&calls))
	}
}

func TestRetryOnHTMLResponse(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.Header().Set("Content-Type", "text/html")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`<html><body>The operation you requested could not be completed by Management Console</body></html>`))
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls (HTML retried), got %d", atomic.LoadInt32(&calls))
	}
}

func TestRetryWriteOnHTMLGatewayError(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.Header().Set("Content-Type", "text/html")
			rw.WriteHeader(http.StatusBadGateway)
			rw.Write([]byte(`<html><body>502 Bad Gateway</body></html>`))
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"site_id":123,"res":0}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodPost, server.URL, []byte(`{"domain":"example.com"}`))
	SetHeaders(client, req, contentTypeApplicationJson, CreateSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 after retry, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls (write retried on HTML gateway error), got %d", atomic.LoadInt32(&calls))
	}
}

func TestNoRetryWriteOn502WithJSON(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&calls, 1)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(`{"error":"backend processing failed"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodPost, server.URL, []byte(`{"domain":"example.com"}`))
	SetHeaders(client, req, contentTypeApplicationJson, CreateSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("Expected 502 returned without retry, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("Expected exactly 1 call (no retry for write with JSON 502), got %d", atomic.LoadInt32(&calls))
	}
}

func TestNoRetryOn400(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&calls, 1)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"res":1,"res_message":"Invalid parameter"}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 without retry, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 1 {
		t.Errorf("Expected exactly 1 call (no retry on 400), got %d", atomic.LoadInt32(&calls))
	}
}

func TestRetryWithRequestBodyIntact(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	var lastBody string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		body, _ := io.ReadAll(req.Body)
		lastBody = string(body)
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.WriteHeader(http.StatusTooManyRequests)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	payload := []byte(`{"domain":"example.com","account_id":12345}`)
	req, _ := PrepareJsonRequest(http.MethodPost, server.URL, payload)
	SetHeaders(client, req, contentTypeApplicationJson, CreateSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	if lastBody != string(payload) {
		t.Errorf("Request body was not preserved on retry.\nExpected: %s\nGot: %s", string(payload), lastBody)
	}
}

func TestRetryOnTransientNetworkError(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			conn, _, _ := rw.(http.Hijacker).Hijack()
			conn.Close()
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Expected retry to recover from network error, got: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 after retry, got %d", resp.StatusCode)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls (1 network error + 1 success), got %d", atomic.LoadInt32(&calls))
	}
}

func TestRetryOnNetworkErrorExhaustedReturnsError(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "http://192.0.2.1:1", BaseURLRev2: "http://192.0.2.1:1", BaseURLAPI: "http://192.0.2.1:1"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 50}}

	req, _ := PrepareJsonRequest(http.MethodGet, "http://192.0.2.1:1/test", nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	_, err := client.executeRequest(req)
	if err == nil {
		t.Fatal("Expected error after all retries exhausted on network failure")
	}
}

func TestVerifyRetriesOnTransientHTMLResponse(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			rw.Header().Set("Content-Type", "text/html")
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`<html><body>Management Console Error</body></html>`))
			return
		}
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"account_type":"Customer","account_id":123,"parent_id":1,"account_name":"test","plan_id":"ent100","plan_name":"ENTERPRISE","res":0,"res_message":"OK","debug_info":{"id-info":"999"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := NewClient(config)

	result, err := client.Verify()
	if err != nil {
		t.Fatalf("Verify should have succeeded after retry, got: %s", err)
	}
	if result == nil {
		t.Fatal("Expected non-nil account status response")
	}
	if result.AccountID != 123 {
		t.Errorf("Expected account_id 123, got %d", result.AccountID)
	}
	if atomic.LoadInt32(&calls) != 2 {
		t.Errorf("Expected 2 calls (1 HTML + 1 success), got %d", atomic.LoadInt32(&calls))
	}
}

func TestNormalJSONResponseBodyNotCorrupted(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	expectedBody := `{"site_id":123,"domain":"example.com","res":0,"res_message":"OK"}`
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(expectedBody))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	req, _ := PrepareJsonRequest(http.MethodGet, server.URL, nil)
	SetHeaders(client, req, contentTypeApplicationJson, ReadSite, nil)

	resp, err := client.executeRequest(req)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if string(body) != expectedBody {
		t.Errorf("Response body was corrupted.\nExpected: %s\nGot: %s", expectedBody, string(body))
	}
}
