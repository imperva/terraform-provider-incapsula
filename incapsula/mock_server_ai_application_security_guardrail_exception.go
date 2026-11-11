// Mock Imperva API Server — AI Application Security guardrail exception endpoints
//
// POST/GET/PUT/DELETE /v3/applications/{applicationId}/policies/{policyId}/guardian-exceptions[/{id}], mirroring the
// management service's GuardianExceptionController: guardianId is required and must be a guardrail of the policy
// (404 otherwise), exceptionPattern and note are trimmed when saved, the reason must be one of the allowed codes,
// and a guardrail exception of another policy is not found.

package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// MockAiApplicationSecurityGuardrailException is a guardrail exception in the mock server (GuardianExceptionResponseDto).
type MockAiApplicationSecurityGuardrailException struct {
	Id               string  `json:"id"`
	GuardrailId      int64   `json:"guardianId"`
	PolicyId         string  `json:"policyId"`
	ExceptionPattern string  `json:"exceptionPattern"`
	ExceptionReason  string  `json:"exceptionReason"`
	Note             *string `json:"note"`
	Active           bool    `json:"active"`
	Author           string  `json:"author"`
	UpdatedBy        *string `json:"updatedBy"`
}

// mockGuardrailExceptionBody is the create/update body; pointers tell a missing field from a zero value.
type mockGuardrailExceptionBody struct {
	GuardrailId      *int64  `json:"guardianId"`
	ExceptionPattern *string `json:"exceptionPattern"`
	ExceptionReason  *string `json:"exceptionReason"`
	Note             *string `json:"note"`
	Active           *bool   `json:"active"`
}

const mockGuardrailExceptionUser = "mock-user@example.com"

func (m *MockImpervaServer) handleAiApplicationSecurityGuardrailExceptions(w http.ResponseWriter, r *http.Request, path string) {
	// .../v3/applications/{applicationId}/policies/{policyId}/guardian-exceptions[/{id}]
	segments := strings.Split(strings.Trim(path[strings.Index(path, "v3/applications/")+len("v3/applications/"):], "/"), "/")
	if len(segments) < 4 || segments[1] != "policies" || segments[3] != "guardian-exceptions" {
		m.writeMockGuardrailExceptionError(w, http.StatusNotFound, fmt.Sprintf("Invalid guardrail exception path: %s", path))
		return
	}
	applicationID, policyID, id := segments[0], segments[2], ""
	if len(segments) >= 5 {
		id = segments[4]
	}

	var body mockGuardrailExceptionBody
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		raw, _ := ioutil.ReadAll(r.Body)
		var envelope AiApplicationSecurityImpervaApiBody
		if json.Unmarshal(raw, &envelope) != nil || json.Unmarshal(envelope.Data, &body) != nil {
			m.writeMockGuardrailExceptionError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	policy, ok := m.aiApplicationSecurityPolicies[policyID]
	if !ok || policy.ApplicationId != applicationID {
		m.writeMockGuardrailExceptionError(w, http.StatusNotFound, fmt.Sprintf("Policy with ID %s not found for application %s", policyID, applicationID))
		return
	}

	if r.Method == http.MethodPost {
		if status, problem := checkMockGuardrailExceptionBody(policy, &body, false); status != 0 {
			m.writeMockGuardrailExceptionError(w, status, problem)
			return
		}
		gex := &MockAiApplicationSecurityGuardrailException{
			Id:       fmt.Sprintf("bbbbbbbb-0000-0000-0000-%012d", m.nextAiApplicationSecurityGuardrailExceptionID),
			PolicyId: policyID,
			Author:   mockGuardrailExceptionUser,
			Active:   body.Active == nil || *body.Active, // optional on create, defaults to enabled
		}
		m.nextAiApplicationSecurityGuardrailExceptionID++
		applyMockGuardrailExceptionBody(gex, &body)
		m.aiApplicationSecurityGuardrailExceptions[gex.Id] = gex
		m.writeMockGuardrailException(w, http.StatusCreated, gex)
		return
	}

	gex, ok := m.aiApplicationSecurityGuardrailExceptions[id]
	if !ok || gex.PolicyId != policyID {
		m.writeMockGuardrailExceptionError(w, http.StatusNotFound, "Guardian exception not found")
		return
	}

	switch r.Method {
	case http.MethodGet:
		m.writeMockGuardrailException(w, http.StatusOK, gex)
	case http.MethodPut:
		if status, problem := checkMockGuardrailExceptionBody(policy, &body, true); status != 0 {
			m.writeMockGuardrailExceptionError(w, status, problem)
			return
		}
		applyMockGuardrailExceptionBody(gex, &body)
		gex.Active = *body.Active
		updatedBy := mockGuardrailExceptionUser
		gex.UpdatedBy = &updatedBy
		m.writeMockGuardrailException(w, http.StatusOK, gex)
	case http.MethodDelete:
		delete(m.aiApplicationSecurityGuardrailExceptions, id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"data": map[string]string{"message": "Guardian exception deleted successfully"}})
	default:
		m.writeMockGuardrailExceptionError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// checkMockGuardrailExceptionBody applies the backend's checks: guardianId required (400) and in the policy (404),
// active required on update, 3-200 words in the trimmed pattern, an allowed reason (400). Returns 0 when valid.
func checkMockGuardrailExceptionBody(policy *MockAiApplicationSecurityPolicy, body *mockGuardrailExceptionBody, isUpdate bool) (int, string) {
	if body.GuardrailId == nil {
		return http.StatusBadRequest, "Guardian ID cannot be null"
	}
	inPolicy := false
	for _, g := range append(append([]AiApplicationSecurityGuardrail{}, policy.Request...), policy.Response...) {
		inPolicy = inPolicy || g.Id == *body.GuardrailId
	}
	if !inPolicy {
		return http.StatusNotFound, fmt.Sprintf("Guardian (id: %d) not found", *body.GuardrailId)
	}
	if isUpdate && body.Active == nil {
		return http.StatusBadRequest, "Active state must be boolean"
	}
	if body.ExceptionPattern == nil {
		return http.StatusBadRequest, "Exception pattern cannot be null"
	}
	if words := len(strings.Fields(*body.ExceptionPattern)); words < 3 || words > 200 {
		return http.StatusBadRequest, fmt.Sprintf("Exception pattern must contain between 3 and 200 words, got %d.", words)
	}
	if body.ExceptionReason == nil || !stringInSlice(aiApplicationSecurityGuardrailExceptionReasons, strings.TrimSpace(*body.ExceptionReason)) {
		return http.StatusBadRequest, fmt.Sprintf("Invalid exception reason. Allowed values: %s.", strings.Join(aiApplicationSecurityGuardrailExceptionReasons, ", "))
	}
	return 0, ""
}

// applyMockGuardrailExceptionBody stores the body's fields, trimmed like the backend.
func applyMockGuardrailExceptionBody(gex *MockAiApplicationSecurityGuardrailException, body *mockGuardrailExceptionBody) {
	gex.GuardrailId = *body.GuardrailId
	gex.ExceptionPattern = strings.TrimSpace(*body.ExceptionPattern)
	gex.ExceptionReason = strings.TrimSpace(*body.ExceptionReason)
	gex.Note = nil
	if body.Note != nil {
		note := strings.TrimSpace(*body.Note)
		gex.Note = &note
	}
}

func (m *MockImpervaServer) writeMockGuardrailException(w http.ResponseWriter, status int, gex *MockAiApplicationSecurityGuardrailException) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": gex})
}

func (m *MockImpervaServer) writeMockGuardrailExceptionError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	m.writeAiApplicationSecurityError(w, message)
}
