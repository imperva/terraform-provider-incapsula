// Mock Imperva API Server — AI Application Security policy endpoints
//
// Implements the stateful CRUD behaviour of the AI Application Security management service
// (/v3/applications/{applicationId}/policies) needed by the
// incapsula_ai_application_security_policy acceptance tests. Requests and responses are wrapped
// in the backend ImpervaApiBody<T> envelope ({"data": …}), matching PolicyController.
//
// The backend enforces a single policy per application; the mock mirrors that.

package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// mockCaid extracts the caid query parameter as an int64 (0 if absent/invalid), so the
// mock can echo the owning account back on reads (used on import).
func mockCaid(r *http.Request) int64 {
	if caid := r.URL.Query().Get("caid"); caid != "" {
		if parsed, err := strconv.ParseInt(caid, 10, 64); err == nil {
			return parsed
		}
	}
	return 0
}

// MockAiApplicationSecurityPolicy represents an AI Application Security policy in the mock server.
// Field shapes match PolicyDto (guardrails split into request/response phase buckets).
type MockAiApplicationSecurityPolicy struct {
	Id            string                           `json:"id"`
	AccountId     int64                            `json:"accountId"`
	ApplicationId string                           `json:"applicationId"`
	Name          string                           `json:"name"`
	Description   string                           `json:"description"`
	Active        bool                             `json:"active"`
	Request       []AiApplicationSecurityGuardrail `json:"request"`
	Response      []AiApplicationSecurityGuardrail `json:"response"`
}

// handleAiApplicationSecurityPolicies routes AI Application Security policy requests. It parses the
// {applicationId} and optional {policyId} path segments robustly, tolerating the
// "/ai-application-security" prefix the client prepends to the mock base URL.
func (m *MockImpervaServer) handleAiApplicationSecurityPolicies(w http.ResponseWriter, r *http.Request, path string) {
	marker := "v3/applications/"
	idx := strings.Index(path, marker)
	if idx < 0 {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid policy path: %s", path))
		return
	}
	// rest = {applicationId}/policies[/{policyId}]
	rest := path[idx+len(marker):]
	segments := strings.Split(strings.Trim(rest, "/"), "/")
	if len(segments) < 2 || segments[1] != "policies" {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid policy path: %s", path))
		return
	}
	applicationID := segments[0]
	policyID := ""
	if len(segments) >= 3 {
		policyID = segments[2]
	}

	switch r.Method {
	case http.MethodPost:
		m.handleAiApplicationSecurityPolicyCreate(w, r, applicationID)
	case http.MethodGet:
		m.handleAiApplicationSecurityPolicyRead(w, policyID)
	case http.MethodPatch:
		m.handleAiApplicationSecurityPolicyUpdate(w, r, policyID)
	case http.MethodDelete:
		m.handleAiApplicationSecurityPolicyDelete(w, policyID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// decodeAiApplicationSecurityPolicyRequest unwraps the {"data": {…}} envelope into a policy request.
func (m *MockImpervaServer) decodeAiApplicationSecurityPolicyRequest(r *http.Request) (*AiApplicationSecurityPolicyRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var envelope AiApplicationSecurityImpervaApiBody
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	var req AiApplicationSecurityPolicyRequest
	if err := json.Unmarshal(envelope.Data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// handleAiApplicationSecurityPolicyCreate handles POST /v3/applications/{applicationId}/policies.
func (m *MockImpervaServer) handleAiApplicationSecurityPolicyCreate(w http.ResponseWriter, r *http.Request, applicationID string) {
	req, err := m.decodeAiApplicationSecurityPolicyRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid request body: %s", err))
		return
	}

	accountID := mockCaid(r)

	m.mu.Lock()
	// Enforce a single policy per application, matching the backend.
	for _, p := range m.aiApplicationSecurityPolicies {
		if p.ApplicationId == applicationID {
			m.mu.Unlock()
			w.WriteHeader(http.StatusBadRequest)
			m.writeAiApplicationSecurityError(w, fmt.Sprintf("A policy already exists for application %s", applicationID))
			return
		}
	}

	policyID := fmt.Sprintf("aaaaaaaa-0000-0000-0000-%012d", m.nextAiApplicationSecurityPolicyID)
	m.nextAiApplicationSecurityPolicyID++

	policy := &MockAiApplicationSecurityPolicy{
		Id:            policyID,
		AccountId:     accountID,
		ApplicationId: applicationID,
		Name:          req.Name,
		Description:   req.Description,
		Active:        req.Active,
		Request:       normalizeGuardrails(req.Request),
		Response:      normalizeGuardrails(req.Response),
	}
	m.aiApplicationSecurityPolicies[policyID] = policy
	m.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	m.writeAiApplicationSecurityPolicyObject(w, policy)
}

// handleAiApplicationSecurityPolicyRead handles GET /v3/applications/{applicationId}/policies/{policyId}.
// The backend looks the policy up by id regardless of the path's application segment.
func (m *MockImpervaServer) handleAiApplicationSecurityPolicyRead(w http.ResponseWriter, policyID string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	policy, ok := m.aiApplicationSecurityPolicies[policyID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Policy not found: %s", policyID))
		return
	}

	w.WriteHeader(http.StatusOK)
	m.writeAiApplicationSecurityPolicyObject(w, policy)
}

// handleAiApplicationSecurityPolicyUpdate handles PATCH /v3/applications/{applicationId}/policies/{policyId}.
// The policy's scalar fields are replaced, but the guardrails are *mutated in place*, mirroring
// PolicyService.mutateExistingGuardians: see mutateAiApplicationSecurityGuardrails.
func (m *MockImpervaServer) handleAiApplicationSecurityPolicyUpdate(w http.ResponseWriter, r *http.Request, policyID string) {
	req, err := m.decodeAiApplicationSecurityPolicyRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid request body: %s", err))
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	policy, ok := m.aiApplicationSecurityPolicies[policyID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Policy not found: %s", policyID))
		return
	}

	policy.Name = req.Name
	policy.Description = req.Description
	policy.Active = req.Active
	policy.Request = mutateAiApplicationSecurityGuardrails(policy.Request, req.Request)
	policy.Response = mutateAiApplicationSecurityGuardrails(policy.Response, req.Response)

	w.WriteHeader(http.StatusOK)
	m.writeAiApplicationSecurityPolicyObject(w, policy)
}

// handleAiApplicationSecurityPolicyDelete handles DELETE /v3/applications/{applicationId}/policies/{policyId}.
func (m *MockImpervaServer) handleAiApplicationSecurityPolicyDelete(w http.ResponseWriter, policyID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.aiApplicationSecurityPolicies[policyID]; !ok {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Policy not found: %s", policyID))
		return
	}

	delete(m.aiApplicationSecurityPolicies, policyID)
	w.WriteHeader(http.StatusNoContent)
}

// writeAiApplicationSecurityPolicyObject writes a single policy wrapped in the {"data": {…}} envelope.
func (m *MockImpervaServer) writeAiApplicationSecurityPolicyObject(w http.ResponseWriter, policy *MockAiApplicationSecurityPolicy) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": policy})
}

// normalizeGuardrails returns a non-nil slice so the {"data": …} envelope always carries
// request/response arrays (never JSON null), matching the backend's serialization.
func normalizeGuardrails(guardrails []AiApplicationSecurityGuardrail) []AiApplicationSecurityGuardrail {
	if guardrails == nil {
		return []AiApplicationSecurityGuardrail{}
	}
	return guardrails
}

// mutateAiApplicationSecurityGuardrails applies a PATCH body's guardrails to the ones a policy
// already holds in that phase, mirroring PolicyService.mutateExistingGuardians. PATCH is
// update-only:
//
//   - guardrails are matched by type within the phase (the path already fixes the phase);
//   - an incoming guardrail whose type has no existing counterpart is silently ignored — PATCH
//     never inserts;
//   - an existing guardrail whose type is absent from the body is left untouched — PATCH never
//     deletes;
//   - if the body repeats a type, the first occurrence wins (Collectors.toMap(…, (a, b) -> a)).
//
// Only mode, active and config are mutable. This is deliberately *not* a full replace: modelling
// it as one would let acceptance tests pass against the mock while failing against the real
// service, which is exactly the divergence these tests exist to catch.
func mutateAiApplicationSecurityGuardrails(existing, incoming []AiApplicationSecurityGuardrail) []AiApplicationSecurityGuardrail {
	byType := make(map[string]AiApplicationSecurityGuardrail, len(incoming))
	for _, g := range incoming {
		if _, seen := byType[g.GuardrailType]; !seen {
			byType[g.GuardrailType] = g
		}
	}

	mutated := normalizeGuardrails(existing)
	for i := range mutated {
		update, ok := byType[mutated[i].GuardrailType]
		if !ok {
			continue
		}
		mutated[i].GuardrailMode = update.GuardrailMode
		mutated[i].Active = update.Active
		mutated[i].Config = update.Config
	}
	return mutated
}
