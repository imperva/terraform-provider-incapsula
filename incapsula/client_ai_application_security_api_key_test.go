package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAiApplicationSecurityApiKeyBadConnection(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: 1}}

	if _, err := client.CreateAiApplicationSecurityApiKey(55, "app-id", &AiApplicationSecurityApiKeyRequest{Name: "k"}); err == nil {
		t.Errorf("Should have received an error")
	}
	if _, err := client.GetAiApplicationSecurityApiKey(55, 1); err == nil {
		t.Errorf("Should have received an error")
	}
	if err := client.DeleteAiApplicationSecurityApiKey(55, "app-id", 1); err == nil {
		t.Errorf("Should have received an error")
	}
}

func TestAiApplicationSecurityApiKeyCreate(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	expectedEndpoint := fmt.Sprintf("%s?caid=%d", fmt.Sprintf(aiApplicationSecurityApiKeyAppEndpoint, applicationID), accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != expectedEndpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", expectedEndpoint, req.URL.String())
		}
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST, got %s", req.Method)
		}

		var envelope AiApplicationSecurityImpervaApiBody
		if err := json.NewDecoder(req.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode request body envelope: %s", err)
		}
		// The request body must be exactly {"data":{"name":"..."}}.
		if strings.TrimSpace(string(envelope.Data)) != `{"name":"my-key"}` {
			t.Errorf("Unexpected request data payload: %s", string(envelope.Data))
		}

		rw.WriteHeader(201)
		rw.Write([]byte(`{"data":{"apiKey":{"id":42,"name":"my-key","maskedApiKey":"****0042","accountId":55,"active":true,"createdAt":1700000000000,"applicationId":"11111111-1111-1111-1111-111111111111"},"fullKey":"aifw-000000000042-plaintext-secret"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.CreateAiApplicationSecurityApiKey(accountID, applicationID, &AiApplicationSecurityApiKeyRequest{Name: "my-key"})
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp.FullKey != "aifw-000000000042-plaintext-secret" {
		t.Errorf("Unexpected fullKey: %s", resp.FullKey)
	}
	if resp.ApiKey.Id != 42 {
		t.Errorf("Unexpected nested apiKey.id: %d", resp.ApiKey.Id)
	}
}

func TestAiApplicationSecurityApiKeyCreateMaxLimitError(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`{"errors":[{"detail":"maximum number of API keys reached"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.CreateAiApplicationSecurityApiKey(55, "app-id", &AiApplicationSecurityApiKeyRequest{Name: "k"})
	if err == nil {
		t.Fatalf("Should have received an error for a 400 response")
	}
	// The error must mention the 5-key limit and preserve the backend message.
	if !strings.Contains(err.Error(), "5 API keys") {
		t.Errorf("Error should mention the 5-key limit, got: %s", err)
	}
	if !strings.Contains(err.Error(), "maximum number of API keys reached") {
		t.Errorf("Error should preserve the backend message, got: %s", err)
	}
}

func TestAiApplicationSecurityApiKeyRead(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	expectedEndpoint := fmt.Sprintf("%s?caid=%d", aiApplicationSecurityApiKeyAccountEndpoint, accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != expectedEndpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", expectedEndpoint, req.URL.String())
		}
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", req.Method)
		}
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":{"apiKeys":[{"id":7,"name":"other","applicationId":"aaaa"},{"id":42,"name":"my-key","maskedApiKey":"****0042","accountId":55,"active":true,"createdAt":1700000000000,"lastUsedAt":1700000009999,"applicationId":"11111111-1111-1111-1111-111111111111"}],"totalCount":2}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityApiKey(accountID, 42)
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp == nil {
		t.Fatalf("Expected to find the API key with id 42")
	}
	if resp.Name != "my-key" || resp.ApplicationId != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("Unexpected api key resolved: %+v", resp)
	}
	if resp.MaskedApiKey != "****0042" || !resp.Active || resp.LastUsedAt != 1700000009999 {
		t.Errorf("Unexpected metadata: %+v", resp)
	}
}

func TestAiApplicationSecurityApiKeyReadNotFound(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":{"apiKeys":[{"id":7,"name":"other"}],"totalCount":1}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityApiKey(55, 42)
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp != nil {
		t.Errorf("Expected nil when id not present in list, got: %+v", resp)
	}
}

func TestAiApplicationSecurityApiKeyReadHttp404(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte(`{"errors":[{"status":"404"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityApiKey(55, 42)
	if err != nil {
		t.Fatalf("Should not have received an error for 404, got: %s", err)
	}
	if resp != nil {
		t.Errorf("Expected nil response for 404, got: %+v", resp)
	}
}

func TestAiApplicationSecurityApiKeyDelete(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	apiKeyID := int64(42)
	expectedEndpoint := fmt.Sprintf("%s/%d?caid=%d", fmt.Sprintf(aiApplicationSecurityApiKeyAppEndpoint, applicationID), apiKeyID, accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != expectedEndpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", expectedEndpoint, req.URL.String())
		}
		if req.Method != http.MethodDelete {
			t.Errorf("Expected DELETE, got %s", req.Method)
		}
		rw.WriteHeader(204)
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	if err := client.DeleteAiApplicationSecurityApiKey(accountID, applicationID, apiKeyID); err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
}

func TestAiApplicationSecurityApiKeyErrorResponse(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		rw.Write([]byte(`Server error`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	if _, err := client.CreateAiApplicationSecurityApiKey(55, "app-id", &AiApplicationSecurityApiKeyRequest{Name: "k"}); err == nil {
		t.Errorf("Should have received an error")
	}
	if _, err := client.GetAiApplicationSecurityApiKey(55, 1); err == nil {
		t.Errorf("Should have received an error")
	}
	if err := client.DeleteAiApplicationSecurityApiKey(55, "app-id", 1); err == nil {
		t.Errorf("Should have received an error")
	}
}
