package incapsula

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientPatchAPIClientBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	_, err := client.PatchAPIClient(context.Background(), "test-client", &APIClientUpdateRequest{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "PATCH request failed") {
		t.Errorf("Should have received a PATCH request error, got: %s", err)
	}
}

func TestClientPatchAPIClientBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.PatchAPIClient(context.Background(), "test-client", &APIClientUpdateRequest{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "failed to unmarshal PATCH response") {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientPatchAPIClientInvalidStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.PatchAPIClient(context.Background(), "test-client", &APIClientUpdateRequest{})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "PATCH /v3/api-client/test-client failed") {
		t.Errorf("Should have received a bad status error, got: %s", err)
	}
}

func TestClientPatchAPIClientValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"api_client_id":"1234","api_key":"key-abc","enabled":true,"expiration_date":"2026-01-01T00:00:00Z","last_used_at":"2025-09-08T10:15:00Z","grace_period":2000}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	resp, err := client.PatchAPIClient(context.Background(), "test-client", &APIClientUpdateRequest{})
	if err != nil {
		t.Errorf("Should not have received an error: %s", err)
	}
	if resp.APIClientID != "1234" || resp.APIKey != "key-abc" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestClientGetAPIClientBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	_, err := client.GetAPIClient(context.Background(), "test-client")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "GET request failed") {
		t.Errorf("Should have received a GET request error, got: %s", err)
	}
}

func TestClientGetAPIClientBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.GetAPIClient(context.Background(), "test-client")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "failed to unmarshal GET response") {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientGetAPIClientInvalidStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte(`{"error":"not found"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.GetAPIClient(context.Background(), "test-client")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "GET /v3/api-client/test-client failed") {
		t.Errorf("Should have received a bad status error, got: %s", err)
	}
}

func TestClientGetAPIClientValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"api_client_id":"1234","api_key":"key-abc","enabled":true,"expiration_date":"2026-01-01T00:00:00Z","last_used_at":"2025-09-08T10:15:00Z","grace_period":2000}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	resp, err := client.GetAPIClient(context.Background(), "test-client")
	if err != nil {
		t.Errorf("Should not have received an error: %s", err)
	}
	if resp.APIClientID != "1234" || resp.APIKey != "key-abc" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestClientDeleteAPIClientBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	err := client.DeleteAPIClient(context.Background(), "test-client")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "DELETE request failed") {
		t.Errorf("Should have received a DELETE request error, got: %s", err)
	}
}

func TestClientDeleteAPIClientBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteAPIClient(context.Background(), "test-client")
	if err != nil {
		t.Errorf("Should not have received an error for DELETE (body ignored): %s", err)
	}
}

func TestClientDeleteAPIClientInvalidStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteAPIClient(context.Background(), "test-client")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.Contains(err.Error(), "DELETE /v3/api-client/test-client failed") {
		t.Errorf("Should have received a bad status error, got: %s", err)
	}
}

func TestClientDeleteAPIClientValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"result":"ok"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteAPIClient(context.Background(), "test-client")
	if err != nil {
		t.Errorf("Should not have received an error: %s", err)
	}
}
