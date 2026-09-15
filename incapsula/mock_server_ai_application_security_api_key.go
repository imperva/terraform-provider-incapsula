// Mock Imperva API Server — AI Application Security API key endpoints
//
// Implements the stateful behaviour of the AI Application Security management service's
// api-key endpoints needed by the incapsula_ai_application_security_api_key acceptance tests:
//   - POST   /v3/applications/{applicationId}/api-keys        (create, app-scoped)
//   - DELETE /v3/applications/{applicationId}/api-keys/{id}   (delete, app-scoped)
//   - GET    /v3/api-keys                                     (account-level list, used by Read/import)
//
// Requests and responses use the backend ImpervaApiBody<T> envelope ({"data": …}).
// The backend caps an account at 5 API keys; the mock mirrors that with a 400.

package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const mockAiApplicationSecurityApiKeyMaxPerAccount = 5

// MockAiApplicationSecurityApiKey represents an AI Application Security API key in the mock server.
// Field shapes match ApiKeyDto.
type MockAiApplicationSecurityApiKey struct {
	Id            int64  `json:"id"`
	MaskedApiKey  string `json:"maskedApiKey"`
	AccountId     int64  `json:"accountId"`
	Name          string `json:"name"`
	CreatedAt     int64  `json:"createdAt"`
	Active        bool   `json:"active"`
	LastUsedAt    int64  `json:"lastUsedAt,omitempty"`
	ApplicationId string `json:"applicationId"`
}

// handleAiApplicationSecurityApiKeys routes AI Application Security api-key requests. It distinguishes the
// account-level list (GET /v3/api-keys) from the application-scoped create/delete paths
// (/v3/applications/{applicationId}/api-keys[/{apiKeyId}]).
func (m *MockImpervaServer) handleAiApplicationSecurityApiKeys(w http.ResponseWriter, r *http.Request, path string) {
	// Account-level list: /v3/api-keys (no /applications/ segment).
	if !strings.Contains(path, "v3/applications/") {
		if r.Method == http.MethodGet {
			m.handleAiApplicationSecurityApiKeyList(w, r)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Method not allowed on api-keys list: %s", r.Method))
		return
	}

	// Application-scoped: rest = {applicationId}/api-keys[/{apiKeyId}]
	marker := "v3/applications/"
	idx := strings.Index(path, marker)
	rest := path[idx+len(marker):]
	segments := strings.Split(strings.Trim(rest, "/"), "/")
	if len(segments) < 2 || segments[1] != "api-keys" {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid api-key path: %s", path))
		return
	}
	applicationID := segments[0]
	apiKeyID := ""
	if len(segments) >= 3 {
		apiKeyID = segments[2]
	}

	switch r.Method {
	case http.MethodPost:
		m.handleAiApplicationSecurityApiKeyCreate(w, r, applicationID)
	case http.MethodDelete:
		m.handleAiApplicationSecurityApiKeyDelete(w, apiKeyID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// handleAiApplicationSecurityApiKeyCreate handles POST /v3/applications/{applicationId}/api-keys.
func (m *MockImpervaServer) handleAiApplicationSecurityApiKeyCreate(w http.ResponseWriter, r *http.Request, applicationID string) {
	req, err := m.decodeAiApplicationSecurityApiKeyRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid request body: %s", err))
		return
	}

	accountID := mockCaid(r)

	m.mu.Lock()
	// Enforce the max-5-keys-per-account limit, matching the backend.
	count := 0
	for _, k := range m.aiApplicationSecurityApiKeys {
		if k.AccountId == accountID {
			count++
		}
	}
	if count >= mockAiApplicationSecurityApiKeyMaxPerAccount {
		m.mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("An account may hold at most %d API keys", mockAiApplicationSecurityApiKeyMaxPerAccount))
		return
	}

	id := m.nextAiApplicationSecurityApiKeyID
	m.nextAiApplicationSecurityApiKeyID++

	key := &MockAiApplicationSecurityApiKey{
		Id:           id,
		MaskedApiKey: fmt.Sprintf("****%04d", id),
		AccountId:    accountID,
		Name:         req.Name,
		CreatedAt:    1700000000000,
		// active is always false on the real backend: AIFW-1131 dropped the entity's active
		// column but left the field on the shared ApiKeyDto, so MapStruct silently zero-values
		// it on every response. Mirroring false here keeps the mock honest about that.
		Active:        false,
		ApplicationId: applicationID,
	}
	m.aiApplicationSecurityApiKeys[id] = key
	m.mu.Unlock()

	fullKey := fmt.Sprintf("aifw-%012d-plaintext-secret", id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	m.encodeAiApplicationSecurityData(w, map[string]interface{}{
		"apiKey":  key,
		"fullKey": fullKey,
	})
}

// handleAiApplicationSecurityApiKeyList handles GET /v3/api-keys?caid=, returning the account's keys.
func (m *MockImpervaServer) handleAiApplicationSecurityApiKeyList(w http.ResponseWriter, r *http.Request) {
	accountID := mockCaid(r)

	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*MockAiApplicationSecurityApiKey, 0)
	for _, k := range m.aiApplicationSecurityApiKeys {
		if accountID == 0 || k.AccountId == accountID {
			list = append(list, k)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	m.encodeAiApplicationSecurityData(w, map[string]interface{}{
		"apiKeys":    list,
		"totalCount": len(list),
	})
}

// handleAiApplicationSecurityApiKeyDelete handles DELETE /v3/applications/{applicationId}/api-keys/{apiKeyId}.
func (m *MockImpervaServer) handleAiApplicationSecurityApiKeyDelete(w http.ResponseWriter, apiKeyID string) {
	id, err := strconv.ParseInt(apiKeyID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid api-key id: %s", apiKeyID))
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.aiApplicationSecurityApiKeys[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("API key not found: %d", id))
		return
	}

	delete(m.aiApplicationSecurityApiKeys, id)
	w.WriteHeader(http.StatusNoContent)
}

// decodeAiApplicationSecurityApiKeyRequest unwraps the {"data": {…}} envelope into an api-key request.
func (m *MockImpervaServer) decodeAiApplicationSecurityApiKeyRequest(r *http.Request) (*AiApplicationSecurityApiKeyRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var envelope AiApplicationSecurityImpervaApiBody
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	var req AiApplicationSecurityApiKeyRequest
	if err := json.Unmarshal(envelope.Data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// encodeAiApplicationSecurityData writes payload wrapped in the {"data": …} envelope.
func (m *MockImpervaServer) encodeAiApplicationSecurityData(w http.ResponseWriter, payload interface{}) {
	json.NewEncoder(w).Encode(map[string]interface{}{"data": payload})
}
