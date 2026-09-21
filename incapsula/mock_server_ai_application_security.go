// Mock Imperva API Server — AI Application Security application endpoints
//
// Implements the stateful CRUD behaviour of the AI Application Security management service
// (/v3/applications) needed by the incapsula_ai_application_security_application
// acceptance tests. Requests and responses are wrapped in the backend
// ImpervaApiBody<T> envelope ({"data": …}), matching ApplicationController.

package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// MockAiApplicationSecurityApplication represents an AI Application Security application in the mock server.
// Field shapes match ApplicationDetailsDto / ApplicationRequestDto.
type MockAiApplicationSecurityApplication struct {
	ApplicationId     string                                  `json:"applicationId"`
	Name              string                                  `json:"name"`
	AccountId         int64                                   `json:"accountId"`
	Region            string                                  `json:"region"`
	Status            string                                  `json:"status"`
	StatusDescription string                                  `json:"statusDescription"`
	ApplicationType   string                                  `json:"applicationType"`
	Configuration     *AiApplicationSecurityApplicationConfig `json:"configuration,omitempty"`
}

// handleAiApplicationSecurityApplications routes AI Application Security application requests by method
// and whether an {applicationId} path segment is present.
func (m *MockImpervaServer) handleAiApplicationSecurityApplications(w http.ResponseWriter, r *http.Request, path string) {
	applicationID := strings.TrimPrefix(path, "v3/applications")
	applicationID = strings.TrimPrefix(applicationID, "/")

	switch r.Method {
	case http.MethodPost:
		m.handleAiApplicationSecurityApplicationCreate(w, r)
	case http.MethodGet:
		m.handleAiApplicationSecurityApplicationRead(w, r)
	case http.MethodPatch:
		m.handleAiApplicationSecurityApplicationUpdate(w, r, applicationID)
	case http.MethodDelete:
		m.handleAiApplicationSecurityApplicationDelete(w, applicationID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// decodeAiApplicationSecurityRequest unwraps the {"data": {…}} envelope into an application request.
func (m *MockImpervaServer) decodeAiApplicationSecurityRequest(r *http.Request) (*AiApplicationSecurityApplicationRequest, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var envelope AiApplicationSecurityImpervaApiBody
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	var req AiApplicationSecurityApplicationRequest
	if err := json.Unmarshal(envelope.Data, &req); err != nil {
		return nil, err
	}
	return &req, nil
}

// handleAiApplicationSecurityApplicationCreate handles POST /v3/applications.
func (m *MockImpervaServer) handleAiApplicationSecurityApplicationCreate(w http.ResponseWriter, r *http.Request) {
	req, err := m.decodeAiApplicationSecurityRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid request body: %s", err))
		return
	}

	// Capture the account the application belongs to from the caid query param
	// so Read (used on import) can echo it back via ApplicationDetailsDto.accountId.
	var accountID int64
	if caid := r.URL.Query().Get("caid"); caid != "" {
		if parsed, err := strconv.ParseInt(caid, 10, 64); err == nil {
			accountID = parsed
		}
	}

	m.mu.Lock()
	appID := fmt.Sprintf("00000000-0000-0000-0000-%012d", m.nextAiApplicationSecurityAppID)
	m.nextAiApplicationSecurityAppID++

	region := req.Region
	if region == "" {
		// Backend resolves an unset region from account-management, falling back to US.
		region = "US"
	}

	app := &MockAiApplicationSecurityApplication{
		ApplicationId:     appID,
		Name:              req.Name,
		AccountId:         accountID,
		Region:            region,
		Status:            "CONFIGURED",
		StatusDescription: "Application configured",
		ApplicationType:   req.ApplicationType,
		Configuration:     req.Configuration,
	}
	m.aiApplicationSecurityApplications[appID] = app
	m.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	m.writeAiApplicationSecurityObject(w, app)
}

// handleAiApplicationSecurityApplicationRead handles GET /v3/applications?applicationId=&caid=,
// returning ImpervaApiBody<List<ApplicationDetailsDto>>.
func (m *MockImpervaServer) handleAiApplicationSecurityApplicationRead(w http.ResponseWriter, r *http.Request) {
	applicationID := r.URL.Query().Get("applicationId")

	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]*MockAiApplicationSecurityApplication, 0)
	if applicationID != "" {
		if app, ok := m.aiApplicationSecurityApplications[applicationID]; ok {
			list = append(list, app)
		}
	} else {
		for _, app := range m.aiApplicationSecurityApplications {
			list = append(list, app)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": list})
}

// handleAiApplicationSecurityApplicationUpdate handles PATCH /v3/applications/{applicationId}.
func (m *MockImpervaServer) handleAiApplicationSecurityApplicationUpdate(w http.ResponseWriter, r *http.Request, applicationID string) {
	req, err := m.decodeAiApplicationSecurityRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Invalid request body: %s", err))
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	app, ok := m.aiApplicationSecurityApplications[applicationID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Application not found: %s", applicationID))
		return
	}

	// PATCH is partial: apply only fields present in the request.
	if req.Name != "" {
		app.Name = req.Name
	}
	if req.Region != "" {
		app.Region = req.Region
	}
	if req.Configuration != nil {
		app.Configuration = req.Configuration
	}

	w.WriteHeader(http.StatusOK)
	m.writeAiApplicationSecurityObject(w, app)
}

// handleAiApplicationSecurityApplicationDelete handles DELETE /v3/applications/{applicationId}.
func (m *MockImpervaServer) handleAiApplicationSecurityApplicationDelete(w http.ResponseWriter, applicationID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.aiApplicationSecurityApplications[applicationID]; !ok {
		w.WriteHeader(http.StatusNotFound)
		m.writeAiApplicationSecurityError(w, fmt.Sprintf("Application not found: %s", applicationID))
		return
	}

	delete(m.aiApplicationSecurityApplications, applicationID)
	w.WriteHeader(http.StatusNoContent)
}

// writeAiApplicationSecurityObject writes a single application wrapped in the {"data": {…}} envelope.
func (m *MockImpervaServer) writeAiApplicationSecurityObject(w http.ResponseWriter, app *MockAiApplicationSecurityApplication) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"data": app})
}

// writeAiApplicationSecurityError writes an ImpervaApiBody-style error body.
func (m *MockImpervaServer) writeAiApplicationSecurityError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"errors": []map[string]interface{}{
			{"detail": message},
		},
	})
}
