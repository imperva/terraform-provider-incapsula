package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAiApplicationSecurityApplicationBadConnection(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: 1}}

	_, err := client.CreateAiApplicationSecurityApplication(55, &AiApplicationSecurityApplicationRequest{Name: "app"})
	if err == nil {
		t.Errorf("Should have received an error")
	}

	_, err = client.GetAiApplicationSecurityApplication(55, "app-id")
	if err == nil {
		t.Errorf("Should have received an error")
	}

	_, err = client.UpdateAiApplicationSecurityApplication(55, "app-id", &AiApplicationSecurityApplicationRequest{})
	if err == nil {
		t.Errorf("Should have received an error")
	}

	err = client.DeleteAiApplicationSecurityApplication(55, "app-id")
	if err == nil {
		t.Errorf("Should have received an error")
	}
}

func TestAiApplicationSecurityApplicationCreate(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	expectedEndpoint := fmt.Sprintf("%s?caid=%d", aiApplicationSecurityApplicationEndpoint, accountID)

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
		var payload AiApplicationSecurityApplicationRequest
		if err := json.Unmarshal(envelope.Data, &payload); err != nil {
			t.Fatalf("Failed to decode envelope data: %s", err)
		}
		if payload.Name != "my-app" {
			t.Errorf("Expected name my-app, got %s", payload.Name)
		}

		rw.WriteHeader(201)
		rw.Write([]byte(`{"data":{"applicationId":"11111111-1111-1111-1111-111111111111","name":"my-app","applicationType":"API","region":"US","status":"CONFIGURED"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.CreateAiApplicationSecurityApplication(accountID, &AiApplicationSecurityApplicationRequest{Name: "my-app", ApplicationType: "API"})
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp.ApplicationId != "11111111-1111-1111-1111-111111111111" {
		t.Errorf("Unexpected application id: %s", resp.ApplicationId)
	}
}

func TestAiApplicationSecurityApplicationRead(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	expectedEndpoint := fmt.Sprintf("%s?applicationId=%s&caid=%d", aiApplicationSecurityApplicationEndpoint, applicationID, accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != expectedEndpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", expectedEndpoint, req.URL.String())
		}
		if req.Method != http.MethodGet {
			t.Errorf("Expected GET, got %s", req.Method)
		}
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":[{"applicationId":"11111111-1111-1111-1111-111111111111","name":"my-app","applicationType":"API","region":"US","status":"CONFIGURED"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityApplication(accountID, applicationID)
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp == nil || resp.Name != "my-app" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestAiApplicationSecurityApplicationReadEmptyList(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":[]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityApplication(55, "does-not-exist")
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp != nil {
		t.Errorf("Expected nil response for empty list, got: %+v", resp)
	}
}

func TestAiApplicationSecurityApplicationReadNotFound(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte(`{"errors":[{"status":"404"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetAiApplicationSecurityApplication(55, "does-not-exist")
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp != nil {
		t.Errorf("Expected nil response for 404, got: %+v", resp)
	}
}

func TestAiApplicationSecurityApplicationUpdate(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	expectedEndpoint := fmt.Sprintf("%s/%s?caid=%d", aiApplicationSecurityApplicationEndpoint, applicationID, accountID)

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
		rw.Write([]byte(`{"data":{"applicationId":"11111111-1111-1111-1111-111111111111","name":"renamed-app","applicationType":"API","region":"EU","status":"CONFIGURED"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.UpdateAiApplicationSecurityApplication(accountID, applicationID, &AiApplicationSecurityApplicationRequest{Name: "renamed-app", Region: "EU"})
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	if resp.Name != "renamed-app" || resp.Region != "EU" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestAiApplicationSecurityApplicationDelete(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	accountID := 55
	applicationID := "11111111-1111-1111-1111-111111111111"
	expectedEndpoint := fmt.Sprintf("%s/%s?caid=%d", aiApplicationSecurityApplicationEndpoint, applicationID, accountID)

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

	err := client.DeleteAiApplicationSecurityApplication(accountID, applicationID)
	if err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
}

func TestAiApplicationSecurityApplicationErrorResponse(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		rw.Write([]byte(`Server error`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.CreateAiApplicationSecurityApplication(55, &AiApplicationSecurityApplicationRequest{Name: "app"})
	if err == nil {
		t.Errorf("Should have received an error")
	}

	_, err = client.GetAiApplicationSecurityApplication(55, "app-id")
	if err == nil {
		t.Errorf("Should have received an error")
	}

	_, err = client.UpdateAiApplicationSecurityApplication(55, "app-id", &AiApplicationSecurityApplicationRequest{})
	if err == nil {
		t.Errorf("Should have received an error")
	}

	err = client.DeleteAiApplicationSecurityApplication(55, "app-id")
	if err == nil {
		t.Errorf("Should have received an error")
	}
}
