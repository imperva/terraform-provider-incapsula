package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func aiApplicationSecuritySampleGuardrails() ([]AiApplicationSecurityGuardrail, []AiApplicationSecurityGuardrail) {
	request := []AiApplicationSecurityGuardrail{
		{GuardrailType: "PROMPT_INJECTION", GuardrailMode: "BLOCK", GuardrailPhase: "PROMPT", Active: true, Config: json.RawMessage(`{"type":"PROMPT_INJECTION"}`)},
	}
	response := []AiApplicationSecurityGuardrail{
		{GuardrailType: "MODERATION", GuardrailMode: "ALERT", GuardrailPhase: "RESPONSE", Active: true, Config: json.RawMessage(`{"type":"MODERATION"}`)},
	}
	return request, response
}

func TestAiApplicationSecurityPolicyBadConnection(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: 1}}

	req, resp := aiApplicationSecuritySampleGuardrails()
	body := &AiApplicationSecurityPolicyRequest{Name: "p", Active: true, Request: req, Response: resp}

	if _, err := client.CreateAiApplicationSecurityPolicy(55, "app-id", body); err == nil {
		t.Errorf("Should have received an error")
	}
	if _, err := client.GetAiApplicationSecurityPolicy(55, "app-id", "policy-id"); err == nil {
		t.Errorf("Should have received an error")
	}
	if _, err := client.UpdateAiApplicationSecurityPolicy(55, "app-id", "policy-id", body); err == nil {
		t.Errorf("Should have received an error")
	}
	if err := client.DeleteAiApplicationSecurityPolicy(55, "app-id", "policy-id"); err == nil {
		t.Errorf("Should have received an error")
	}
}

func TestAiApplicationSecurityPolicyCreate(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	expectedEndpoint := fmt.Sprintf("%s?caid=%d", fmt.Sprintf(aiApplicationSecurityPolicyEndpoint, applicationID), accountID)

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
		var payload AiApplicationSecurityPolicyRequest
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			t.Fatalf("Failed to decode envelope data: %s", err)
		}
		if payload.Name != "my-policy" {
			t.Errorf("Expected name my-policy, got %s", payload.Name)
		}
		if len(payload.Request) != 1 || payload.Request[0].GuardrailType != "PROMPT_INJECTION" {
			t.Errorf("Unexpected request guardrails: %+v", payload.Request)
		}
		if len(payload.Response) != 1 || payload.Response[0].GuardrailType != "MODERATION" {
			t.Errorf("Unexpected response guardrails: %+v", payload.Response)
		}

		rw.WriteHeader(201)
		rw.Write([]byte(`{"data":{"id":"aaaaaaaa-0000-0000-0000-000000000001","accountId":55,"applicationId":"11111111-1111-1111-1111-111111111111","name":"my-policy","active":true,"request":[{"guardrailType":"PROMPT_INJECTION","guardrailMode":"BLOCK","guardrailPhase":"PROMPT","active":true,"config":{"type":"PROMPT_INJECTION"}}],"response":[{"guardrailType":"MODERATION","guardrailMode":"ALERT","guardrailPhase":"RESPONSE","active":true,"config":{"type":"MODERATION"}}]}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	reqG, respG := aiApplicationSecuritySampleGuardrails()
	resp, err := client.CreateAiApplicationSecurityPolicy(accountID, applicationID, &AiApplicationSecurityPolicyRequest{Name: "my-policy", Active: true, Request: reqG, Response: respG})
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp.Id != "aaaaaaaa-0000-0000-0000-000000000001" {
		t.Errorf("Unexpected policy id: %s", resp.Id)
	}
	if len(resp.Request) != 1 || len(resp.Response) != 1 {
		t.Errorf("Unexpected guardrails in response: %+v", resp)
	}
}

func TestAiApplicationSecurityPolicyCreateError(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`{"errors":[{"detail":"A policy already exists for this application"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	reqG, respG := aiApplicationSecuritySampleGuardrails()
	_, err := client.CreateAiApplicationSecurityPolicy(55, "app-id", &AiApplicationSecurityPolicyRequest{Name: "p", Active: true, Request: reqG, Response: respG})
	if err == nil {
		t.Fatalf("Should have received an error for a 400 response")
	}
}

func TestAiApplicationSecurityPolicyRead(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	policyID := "aaaaaaaa-0000-0000-0000-000000000001"
	expectedEndpoint := fmt.Sprintf("%s/%s?caid=%d", fmt.Sprintf(aiApplicationSecurityPolicyEndpoint, applicationID), policyID, accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != expectedEndpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", expectedEndpoint, req.URL.String())
		}
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", req.Method)
		}
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":{"id":"aaaaaaaa-0000-0000-0000-000000000001","accountId":55,"applicationId":"11111111-1111-1111-1111-111111111111","name":"my-policy","active":true,"request":[{"guardrailType":"PROMPT_INJECTION","guardrailMode":"BLOCK","guardrailPhase":"PROMPT","active":true,"config":{"type":"PROMPT_INJECTION"}}],"response":[]}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityPolicy(accountID, applicationID, policyID)
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp == nil || resp.Name != "my-policy" {
		t.Errorf("Unexpected response: %+v", resp)
	}
	if resp.ApplicationId != applicationID {
		t.Errorf("Expected applicationId %s, got %s", applicationID, resp.ApplicationId)
	}
}

func TestAiApplicationSecurityPolicyReadNotFound(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte(`{"errors":[{"status":"404"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityPolicy(55, "app-id", "does-not-exist")
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp != nil {
		t.Errorf("Expected nil response for 404, got: %+v", resp)
	}
}

func TestAiApplicationSecurityPolicyUpdate(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	policyID := "aaaaaaaa-0000-0000-0000-000000000001"
	expectedEndpoint := fmt.Sprintf("%s/%s?caid=%d", fmt.Sprintf(aiApplicationSecurityPolicyEndpoint, applicationID), policyID, accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != expectedEndpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", expectedEndpoint, req.URL.String())
		}
		if req.Method != http.MethodPatch {
			t.Errorf("Expected PATCH, got %s", req.Method)
		}

		var envelope AiApplicationSecurityImpervaApiBody
		if err := json.NewDecoder(req.Body).Decode(&envelope); err != nil {
			t.Fatalf("Failed to decode request body envelope: %s", err)
		}

		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":{"id":"aaaaaaaa-0000-0000-0000-000000000001","accountId":55,"applicationId":"11111111-1111-1111-1111-111111111111","name":"renamed-policy","active":false,"request":[],"response":[]}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	reqG, respG := aiApplicationSecuritySampleGuardrails()
	resp, err := client.UpdateAiApplicationSecurityPolicy(accountID, applicationID, policyID, &AiApplicationSecurityPolicyRequest{Name: "renamed-policy", Active: false, Request: reqG, Response: respG})
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp.Name != "renamed-policy" || resp.Active {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestAiApplicationSecurityPolicyDelete(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	policyID := "aaaaaaaa-0000-0000-0000-000000000001"
	expectedEndpoint := fmt.Sprintf("%s/%s?caid=%d", fmt.Sprintf(aiApplicationSecurityPolicyEndpoint, applicationID), policyID, accountID)

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

	if err := client.DeleteAiApplicationSecurityPolicy(accountID, applicationID, policyID); err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
}

func TestAiApplicationSecurityPolicyErrorResponse(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		rw.Write([]byte(`Server error`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	reqG, respG := aiApplicationSecuritySampleGuardrails()
	body := &AiApplicationSecurityPolicyRequest{Name: "p", Active: true, Request: reqG, Response: respG}

	if _, err := client.CreateAiApplicationSecurityPolicy(55, "app-id", body); err == nil {
		t.Errorf("Should have received an error")
	}
	if _, err := client.GetAiApplicationSecurityPolicy(55, "app-id", "policy-id"); err == nil {
		t.Errorf("Should have received an error")
	}
	if _, err := client.UpdateAiApplicationSecurityPolicy(55, "app-id", "policy-id", body); err == nil {
		t.Errorf("Should have received an error")
	}
	if err := client.DeleteAiApplicationSecurityPolicy(55, "app-id", "policy-id"); err == nil {
		t.Errorf("Should have received an error")
	}
}
