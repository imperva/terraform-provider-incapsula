package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	testGexAppID    = "11111111-1111-1111-1111-111111111111"
	testGexPolicyID = "22222222-2222-2222-2222-222222222222"
	testGexID       = "33333333-3333-3333-3333-333333333333"
)

func testGexClient(serverURL string) *Client {
	return &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: serverURL}, httpClient: &http.Client{}}
}

// Create posts the guardrail id and every field to the policy-scoped path.
func TestAiApplicationSecurityGuardrailExceptionCreate(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	wantURL := fmt.Sprintf(aiApplicationSecurityGuardrailExceptionEndpoint, testGexAppID, testGexPolicyID) + "?caid=55"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost || req.URL.String() != wantURL {
			t.Errorf("Expected POST %s, got %s %s", wantURL, req.Method, req.URL.String())
		}
		var envelope AiApplicationSecurityImpervaApiBody
		json.NewDecoder(req.Body).Decode(&envelope)
		want := `{"guardianId":7,"exceptionPattern":"some allowed text\n","exceptionReason":"OTHER","note":"","active":true}`
		if strings.TrimSpace(string(envelope.Data)) != want {
			t.Errorf("Unexpected body:\n got  %s\n want %s", envelope.Data, want)
		}
		rw.WriteHeader(201)
		rw.Write([]byte(`{"data":{"id":"` + testGexID + `","guardianId":7,"exceptionPattern":"some allowed text","exceptionReason":"OTHER","note":null,"active":true}}`))
	}))
	defer server.Close()

	created, err := testGexClient(server.URL).CreateAiApplicationSecurityGuardrailException(55, testGexAppID, testGexPolicyID,
		&AiApplicationSecurityGuardrailExceptionRequest{GuardrailId: 7, ExceptionPattern: "some allowed text\n", ExceptionReason: "OTHER", Active: true})
	if err != nil || created.Id != testGexID || created.GuardrailId != 7 || created.Note != nil {
		t.Errorf("Unexpected create result: %+v, %v", created, err)
	}
}

// The service's error text reaches the user.
func TestAiApplicationSecurityGuardrailExceptionCreateError(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		rw.Write([]byte(`{"errors":[{"status":404,"title":"Guardian (id: 99) not found"}]}`))
	}))
	defer server.Close()

	_, err := testGexClient(server.URL).CreateAiApplicationSecurityGuardrailException(55, testGexAppID, testGexPolicyID,
		&AiApplicationSecurityGuardrailExceptionRequest{GuardrailId: 99})
	if err == nil || !strings.Contains(err.Error(), "Guardian (id: 99) not found") {
		t.Errorf("Expected the service's message, got: %v", err)
	}
}

// Read uses the policy-scoped item path; 404 reads as nil so Terraform drops the resource from state.
func TestAiApplicationSecurityGuardrailExceptionRead(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	wantURL := fmt.Sprintf(aiApplicationSecurityGuardrailExceptionEndpoint, testGexAppID, testGexPolicyID) + "/" + testGexID + "?caid=55"
	found := true
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet || req.URL.String() != wantURL {
			t.Errorf("Expected GET %s, got %s %s", wantURL, req.Method, req.URL.String())
		}
		if !found {
			rw.WriteHeader(404)
			return
		}
		rw.Write([]byte(`{"data":{"id":"` + testGexID + `","guardianId":8,"exceptionPattern":"some allowed text","exceptionReason":"OTHER","active":false}}`))
	}))
	defer server.Close()

	read, err := testGexClient(server.URL).GetAiApplicationSecurityGuardrailException(55, testGexAppID, testGexPolicyID, testGexID)
	if err != nil || read == nil || read.GuardrailId != 8 || read.Active {
		t.Errorf("Unexpected read result: %+v, %v", read, err)
	}

	found = false
	read, err = testGexClient(server.URL).GetAiApplicationSecurityGuardrailException(55, testGexAppID, testGexPolicyID, testGexID)
	if err != nil || read != nil {
		t.Errorf("Expected (nil, nil) on 404, got %+v, %v", read, err)
	}
}

// Update is a PUT on the item path; delete treats 404 as already gone.
func TestAiApplicationSecurityGuardrailExceptionUpdateAndDelete(t *testing.T) {
	restore := withShortRetries()
	defer restore()

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPut:
			rw.Write([]byte(`{"data":{"id":"` + testGexID + `","guardianId":8}}`))
		case http.MethodDelete:
			rw.WriteHeader(404)
		default:
			t.Errorf("Unexpected %s", req.Method)
		}
	}))
	defer server.Close()

	client := testGexClient(server.URL)
	updated, err := client.UpdateAiApplicationSecurityGuardrailException(55, testGexAppID, testGexPolicyID, testGexID, &AiApplicationSecurityGuardrailExceptionRequest{GuardrailId: 8})
	if err != nil || updated.GuardrailId != 8 {
		t.Errorf("Unexpected update result: %+v, %v", updated, err)
	}
	if err := client.DeleteAiApplicationSecurityGuardrailException(55, testGexAppID, testGexPolicyID, testGexID); err != nil {
		t.Errorf("Expected success on 404, got: %s", err)
	}
}
