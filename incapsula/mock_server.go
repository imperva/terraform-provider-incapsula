// Mock Imperva API Server
//
// This file implements a mock server for the Imperva API, enabling tests to run
// without requiring real API credentials.
//
// API Documentation:
//   - Cloud v1 API: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/cloud-v1-api-definition.htm
//   - CSP API: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
//
// Implemented endpoints are documented in README.md under "Mock Server for Testing".

package incapsula

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// MockImpervaServer provides a mock implementation of the Imperva API for testing
type MockImpervaServer struct {
	Server *httptest.Server

	// In-memory storage
	mu       sync.RWMutex
	accounts map[int]*MockAccount
	sites    map[int]*MockSite

	// CSP domain storage: map[siteID]map[domain]*MockCSPDomain
	cspDomains map[int]map[string]*MockCSPDomain

	// ID generators
	nextAccountID int
	nextSiteID    int
}

// MockAccount represents an account in the mock server
type MockAccount struct {
	AccountID    int         `json:"account_id"`
	Email        string      `json:"email"`
	ParentID     int         `json:"parent_id"`
	AccountName  string      `json:"account_name"`
	PlanID       string      `json:"plan_id"`
	RefID        string      `json:"ref_id"`
	UserName     string      `json:"user_name"`
	TrialEndDate string      `json:"trial_end_date"`
	Logins       []MockLogin `json:"logins"`
}

// MockLogin represents a login entry for an account
type MockLogin struct {
	LoginID       float64 `json:"login_id"`
	Email         string  `json:"email"`
	EmailVerified bool    `json:"email_verified"`
}

// MockSite represents a site in the mock server
type MockSite struct {
	SiteID     int    `json:"site_id"`
	AccountID  int    `json:"account_id"`
	Domain     string `json:"domain"`
	Status     string `json:"status"`
	SiteType   string `json:"site_type"`
	RefID      string `json:"ref_id"`
	DnsARecord string `json:"dns_a_record"`
	DnsCname   string `json:"dns_cname_record"`
}

// MockCSPDomain represents a CSP pre-approved domain in the mock server
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
type MockCSPDomain struct {
	Domain                   string         `json:"domain"`
	Subdomains               bool           `json:"subdomains"`
	ReferenceID              string         `json:"referenceId"`
	ApplyToAllOnboardedPaths bool           `json:"applyToAllOnboardedPaths"`
	Notes                    []MockCSPNote  `json:"notes,omitempty"`
	Status                   *MockCSPStatus `json:"status,omitempty"`
}

// MockCSPNote represents a note on a CSP domain (FullNote in API docs)
type MockCSPNote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
	Date   int64  `json:"date"`
}

// MockCSPStatus represents the authorization status of a CSP domain (AuthorizationStatus in API docs)
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
type MockCSPStatus struct {
	Blocked     *bool  `json:"blocked,omitempty"`
	Reviewed    *bool  `json:"reviewed,omitempty"`
	Note        string `json:"note,omitempty"`
	Author      string `json:"author,omitempty"`
	ReviewedAt  int64  `json:"reviewedAt,omitempty"`
	LastNoteAt  int64  `json:"lastNoteAt,omitempty"`
	ForceChange bool   `json:"forceChange,omitempty"`
}

// NewMockImpervaServer creates a new mock server instance
func NewMockImpervaServer() *MockImpervaServer {
	mock := &MockImpervaServer{
		accounts:      make(map[int]*MockAccount),
		sites:         make(map[int]*MockSite),
		cspDomains:    make(map[int]map[string]*MockCSPDomain),
		nextAccountID: 1000,
		nextSiteID:    10000,
	}

	// Create the HTTP server with the router
	mock.Server = httptest.NewServer(http.HandlerFunc(mock.router))

	return mock
}

// Close shuts down the mock server
func (m *MockImpervaServer) Close() {
	m.Server.Close()
}

// URL returns the mock server's URL
func (m *MockImpervaServer) URL() string {
	return m.Server.URL
}

// ServeHTTP implements http.Handler interface for standalone server usage
func (m *MockImpervaServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router(w, r)
}

// router handles incoming requests and routes them to the appropriate handler
func (m *MockImpervaServer) router(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Remove leading slash for matching
	path = strings.TrimPrefix(path, "/")

	// Account endpoints
	switch {
	case path == "accounts/add" && r.Method == http.MethodPost:
		m.handleAccountAdd(w, r)
	case path == "account" && r.Method == http.MethodPost:
		m.handleAccountStatus(w, r)
	case path == "accounts/configure" && r.Method == http.MethodPost:
		m.handleAccountUpdate(w, r)
	case path == "accounts/delete" && r.Method == http.MethodPost:
		m.handleAccountDelete(w, r)
	case path == "accounts/data-privacy/show" && r.Method == http.MethodPost:
		m.handleDataPrivacyShow(w, r)
	case path == "accounts/data-privacy/set-region-default" && r.Method == http.MethodPost:
		m.handleDataPrivacySetRegionDefault(w, r)

	// Site endpoints
	case path == "sites/add" && r.Method == http.MethodPost:
		m.handleSiteAdd(w, r)
	case path == "sites/status" && r.Method == http.MethodPost:
		m.handleSiteStatus(w, r)
	case path == "sites/configure" && r.Method == http.MethodPost:
		m.handleSiteUpdate(w, r)
	case path == "sites/delete" && r.Method == http.MethodPost:
		m.handleSiteDelete(w, r)

	// CSP API endpoints
	case strings.HasPrefix(path, "csp-api/v1/sites/"):
		m.handleCSPAPI(w, r, path)

	default:
		// Return 404 for unimplemented endpoints
		m.writeErrorResponse(w, 9999, fmt.Sprintf("Endpoint not implemented: %s %s", r.Method, path))
	}
}

// writeJSONResponse writes a JSON response with the given data
func (m *MockImpervaServer) writeJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// writeErrorResponse writes an error response in Imperva's format
func (m *MockImpervaServer) writeErrorResponse(w http.ResponseWriter, resCode int, message string) {
	response := map[string]interface{}{
		"res":         resCode,
		"res_message": message,
	}
	m.writeJSONResponse(w, response)
}

// writeSuccessResponse writes a success response in Imperva's format
func (m *MockImpervaServer) writeSuccessResponse(w http.ResponseWriter, data interface{}) {
	// If data is already a map, add res: 0 to it
	if mapData, ok := data.(map[string]interface{}); ok {
		mapData["res"] = 0
		mapData["res_message"] = "OK"
		m.writeJSONResponse(w, mapData)
	} else {
		m.writeJSONResponse(w, data)
	}
}

// parseFormValue extracts a form value from the request
func (m *MockImpervaServer) parseFormValue(r *http.Request, key string) string {
	if r.Form == nil {
		r.ParseForm()
	}
	return r.FormValue(key)
}

// parseFormInt extracts an integer form value from the request
func (m *MockImpervaServer) parseFormInt(r *http.Request, key string) int {
	val := m.parseFormValue(r, key)
	if val == "" {
		return 0
	}
	intVal, _ := strconv.Atoi(val)
	return intVal
}

// Account Handlers

// handleAccountAdd handles POST /accounts/add
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/cloud-v1-api-definition.htm
func (m *MockImpervaServer) handleAccountAdd(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r.ParseForm()
	email := m.parseFormValue(r, "email")
	userName := m.parseFormValue(r, "user_name")

	accountID := m.nextAccountID
	m.nextAccountID++

	logins := []MockLogin{
		{
			LoginID:       float64(accountID),
			Email:         email,
			EmailVerified: true,
		},
	}

	account := &MockAccount{
		AccountID:    accountID,
		Email:        email,
		ParentID:     m.parseFormInt(r, "parent_id"),
		AccountName:  m.parseFormValue(r, "account_name"),
		PlanID:       m.parseFormValue(r, "plan_id"),
		RefID:        m.parseFormValue(r, "ref_id"),
		UserName:     userName,
		TrialEndDate: "",
		Logins:       logins,
	}

	m.accounts[accountID] = account

	response := map[string]interface{}{
		"res": 0,
		"account": map[string]interface{}{
			"account_id":   accountID,
			"email":        email,
			"parent_id":    account.ParentID,
			"account_name": account.AccountName,
			"plan_id":      account.PlanID,
			"user_name":    account.UserName,
			"logins":       logins,
		},
	}
	m.writeJSONResponse(w, response)
}

// handleAccountStatus handles POST /account
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/cloud-v1-api-definition.htm
func (m *MockImpervaServer) handleAccountStatus(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r.ParseForm()
	accountID := m.parseFormInt(r, "account_id")

	// If no account_id specified, return default account status (for credential verification)
	if accountID == 0 {
		defaultLogins := []map[string]interface{}{
			{"login_id": 1000.0, "email": "test@example.com", "email_verified": true},
		}
		response := map[string]interface{}{
			"res":         0,
			"res_message": "OK",
			"account_id":  1000,
			"email":       "test@example.com",
			"plan_name":   "Enterprise",
			"plan_id":     "entTrial",
			"user_name":   "test",
			"logins":      defaultLogins,
			"account": map[string]interface{}{
				"account_id":                           1000,
				"email":                                "test@example.com",
				"plan_name":                            "Enterprise",
				"plan_id":                              "entTrial",
				"user_name":                            "test",
				"trial_end_date":                       "",
				"logins":                               defaultLogins,
				"support_level":                        "Standard",
				"supprt_all_tls_versions":              false,
				"wildcard_san_for_new_sites":           "Default",
				"naked_domain_san_for_new_www_sites":   true,
				"inactivity_timeout":                   15,
				"enable_http2_for_new_sites":           true,
				"enable_http2_to_origin_for_new_sites": false,
			},
		}
		m.writeJSONResponse(w, response)
		return
	}

	account, exists := m.accounts[accountID]
	if !exists {
		m.writeErrorResponse(w, 9413, "Unknown/unauthorized account_id")
		return
	}

	loginsInterface := make([]map[string]interface{}, len(account.Logins))
	for i, login := range account.Logins {
		loginsInterface[i] = map[string]interface{}{
			"login_id":       login.LoginID,
			"email":          login.Email,
			"email_verified": login.EmailVerified,
		}
	}

	response := map[string]interface{}{
		"res":          0,
		"res_message":  "OK",
		"account_id":   account.AccountID,
		"email":        account.Email,
		"parent_id":    account.ParentID,
		"account_name": account.AccountName,
		"plan_name":    "Enterprise",
		"plan_id":      account.PlanID,
		"user_name":    account.UserName,
		"ref_id":       account.RefID,
		"account_type": "SubAccount",
		"logins":       loginsInterface,
		"account": map[string]interface{}{
			"account_id":                           account.AccountID,
			"email":                                account.Email,
			"parent_id":                            account.ParentID,
			"account_name":                         account.AccountName,
			"plan_name":                            "Enterprise",
			"plan_id":                              account.PlanID,
			"ref_id":                               account.RefID,
			"user_name":                            account.UserName,
			"trial_end_date":                       account.TrialEndDate,
			"logins":                               loginsInterface,
			"support_level":                        "Standard",
			"supprt_all_tls_versions":              false,
			"wildcard_san_for_new_sites":           "Default",
			"naked_domain_san_for_new_www_sites":   true,
			"inactivity_timeout":                   15,
			"enable_http2_for_new_sites":           true,
			"enable_http2_to_origin_for_new_sites": false,
		},
	}
	m.writeJSONResponse(w, response)
}

// handleAccountUpdate handles POST /accounts/configure
// Uses param/value pattern as per API documentation
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/cloud-v1-api-definition.htm
func (m *MockImpervaServer) handleAccountUpdate(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r.ParseForm()
	accountID := m.parseFormInt(r, "account_id")

	account, exists := m.accounts[accountID]
	if !exists {
		m.writeErrorResponse(w, 9413, "Unknown/unauthorized account_id")
		return
	}

	// API uses param/value pattern for updates
	param := m.parseFormValue(r, "param")
	value := m.parseFormValue(r, "value")

	if param != "" {
		switch param {
		case "account_name", "name":
			account.AccountName = value
		case "email":
			account.Email = value
		case "ref_id":
			account.RefID = value
		case "plan_id":
			account.PlanID = value
		case "user_name":
			account.UserName = value
		case "support_all_tls_versions",
			"naked_domain_san_for_new_www_sites",
			"wildcard_san_for_new_sites",
			"enable_http2_for_new_sites",
			"enable_http2_to_origin_for_new_sites",
			"error_page_template",
			"consent_required",
			"inactivity_timeout":
			// Accept these parameters silently - mock server just acknowledges them
		default:
			// Unknown param - return error code 6001 as per docs
			m.writeErrorResponse(w, 6001, fmt.Sprintf("Invalid parameter: %s", param))
			return
		}
	}

	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
		"account_id":  account.AccountID,
	}
	m.writeJSONResponse(w, response)
}

func (m *MockImpervaServer) handleAccountDelete(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r.ParseForm()
	accountID := m.parseFormInt(r, "account_id")

	if _, exists := m.accounts[accountID]; !exists {
		m.writeErrorResponse(w, 9413, "Unknown/unauthorized account_id")
		return
	}

	delete(m.accounts, accountID)

	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
	}
	m.writeJSONResponse(w, response)
}

func (m *MockImpervaServer) handleDataPrivacyShow(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	response := map[string]interface{}{
		"res":                    0,
		"res_message":            "OK",
		"region":                 "US",
		"waf_log_setup_link":     "",
		"regions_on_setup_alarm": false,
	}
	m.writeJSONResponse(w, response)
}

func (m *MockImpervaServer) handleDataPrivacySetRegionDefault(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
	}
	m.writeJSONResponse(w, response)
}

// Site Handlers

func (m *MockImpervaServer) handleSiteAdd(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r.ParseForm()
	domain := m.parseFormValue(r, "domain")

	// Generate new site ID
	siteID := m.nextSiteID
	m.nextSiteID++

	// Create site
	site := &MockSite{
		SiteID:     siteID,
		AccountID:  m.parseFormInt(r, "account_id"),
		Domain:     domain,
		Status:     "pending",
		SiteType:   m.parseFormValue(r, "site_type"),
		RefID:      m.parseFormValue(r, "ref_id"),
		DnsARecord: "1.2.3.4",
		DnsCname:   fmt.Sprintf("%d.incapdns.net", siteID),
	}

	m.sites[siteID] = site

	// Return response
	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
		"site_id":     siteID,
		"status":      site.Status,
		"domain":      domain,
		"dns": []map[string]interface{}{
			{
				"dns_record_name": domain,
				"set_type_to":     "A",
				"set_data_to":     []string{site.DnsARecord},
			},
			{
				"dns_record_name": domain,
				"set_type_to":     "CNAME",
				"set_data_to":     []string{site.DnsCname},
			},
		},
	}
	m.writeJSONResponse(w, response)
}

func (m *MockImpervaServer) handleSiteStatus(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r.ParseForm()
	siteID := m.parseFormInt(r, "site_id")

	site, exists := m.sites[siteID]
	if !exists {
		m.writeErrorResponse(w, 9413, "Unknown/unauthorized site_id")
		return
	}

	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
		"site_id":     site.SiteID,
		"status":      site.Status,
		"domain":      site.Domain,
		"account_id":  site.AccountID,
		"dns": []map[string]interface{}{
			{
				"dns_record_name": site.Domain,
				"set_type_to":     "A",
				"set_data_to":     []string{site.DnsARecord},
			},
		},
	}
	m.writeJSONResponse(w, response)
}

// handleSiteUpdate handles POST /sites/configure
// Uses param/value pattern as per API documentation
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/cloud-v1-api-definition.htm
func (m *MockImpervaServer) handleSiteUpdate(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r.ParseForm()
	siteID := m.parseFormInt(r, "site_id")

	site, exists := m.sites[siteID]
	if !exists {
		m.writeErrorResponse(w, 9413, "Unknown/unauthorized site_id")
		return
	}

	// API uses param/value pattern for updates
	param := m.parseFormValue(r, "param")
	value := m.parseFormValue(r, "value")

	if param != "" {
		switch param {
		case "domain":
			site.Domain = value
		case "ref_id":
			site.RefID = value
		case "site_type":
			site.SiteType = value
		default:
			// Unknown param - return error
			m.writeErrorResponse(w, 6001, fmt.Sprintf("Invalid parameter: %s", param))
			return
		}
	}

	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
		"site_id":     site.SiteID,
	}
	m.writeJSONResponse(w, response)
}

func (m *MockImpervaServer) handleSiteDelete(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	r.ParseForm()
	siteID := m.parseFormInt(r, "site_id")

	if _, exists := m.sites[siteID]; !exists {
		m.writeErrorResponse(w, 9413, "Unknown/unauthorized site_id")
		return
	}

	delete(m.sites, siteID)

	response := map[string]interface{}{
		"res":         0,
		"res_message": "OK",
	}
	m.writeJSONResponse(w, response)
}

// CSP API Handlers

// handleCSPAPI routes CSP API requests to appropriate handlers
func (m *MockImpervaServer) handleCSPAPI(w http.ResponseWriter, r *http.Request, path string) {
	// Pattern: csp-api/v1/sites/{siteId}/preapprovedlist[/{domainRef}]
	// Pattern: csp-api/v1/sites/{siteId}/domains/{domainRef}/status
	// Pattern: csp-api/v1/sites/{siteId}/domains/{domainRef}/notes

	preapprovedListPattern := regexp.MustCompile(`^csp-api/v1/sites/(\d+)/preapprovedlist(?:/(.*))?$`)
	domainStatusPattern := regexp.MustCompile(`^csp-api/v1/sites/(\d+)/domains/([^/]+)/status$`)
	domainNotesPattern := regexp.MustCompile(`^csp-api/v1/sites/(\d+)/domains/([^/]+)/notes$`)

	if matches := preapprovedListPattern.FindStringSubmatch(path); matches != nil {
		siteID, _ := strconv.Atoi(matches[1])
		domainRef := matches[2]
		m.handleCSPPreapprovedList(w, r, siteID, domainRef)
		return
	}

	if matches := domainStatusPattern.FindStringSubmatch(path); matches != nil {
		siteID, _ := strconv.Atoi(matches[1])
		domainRef := matches[2]
		m.handleCSPDomainStatus(w, r, siteID, domainRef)
		return
	}

	if matches := domainNotesPattern.FindStringSubmatch(path); matches != nil {
		siteID, _ := strconv.Atoi(matches[1])
		domainRef := matches[2]
		m.handleCSPDomainNotes(w, r, siteID, domainRef)
		return
	}

	m.writeErrorResponse(w, 9999, fmt.Sprintf("CSP endpoint not implemented: %s %s", r.Method, path))
}

// handleCSPPreapprovedList handles GET/POST/DELETE for pre-approved domains
func (m *MockImpervaServer) handleCSPPreapprovedList(w http.ResponseWriter, r *http.Request, siteID int, domainRef string) {
	switch r.Method {
	case http.MethodGet:
		if domainRef == "" {
			m.handleCSPListAllDomains(w, siteID)
		} else {
			m.handleCSPGetDomain(w, siteID, domainRef)
		}
	case http.MethodPost:
		m.handleCSPAddDomain(w, r, siteID)
	case http.MethodDelete:
		m.handleCSPDeleteDomain(w, siteID, domainRef)
	default:
		m.writeErrorResponse(w, 9999, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// handleCSPListAllDomains returns all pre-approved domains for a site
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
func (m *MockImpervaServer) handleCSPListAllDomains(w http.ResponseWriter, siteID int) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	domains := make([]map[string]interface{}, 0)
	if siteDomains, exists := m.cspDomains[siteID]; exists {
		for _, domain := range siteDomains {
			domains = append(domains, map[string]interface{}{
				"domain":                   domain.Domain,
				"subdomains":               domain.Subdomains,
				"referenceId":              domain.ReferenceID,
				"applyToAllOnboardedPaths": domain.ApplyToAllOnboardedPaths,
			})
		}
	}

	m.writeJSONResponse(w, domains)
}

// handleCSPGetDomain returns a specific pre-approved domain
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
func (m *MockImpervaServer) handleCSPGetDomain(w http.ResponseWriter, siteID int, domainRef string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	domain := m.decodeDomainRef(domainRef)
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Invalid domain reference")
		return
	}

	siteDomains, exists := m.cspDomains[siteID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		m.writeErrorResponse(w, 404, "Site not found")
		return
	}

	cspDomain, exists := siteDomains[domain]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		m.writeErrorResponse(w, 404, "Domain not found")
		return
	}

	response := map[string]interface{}{
		"domain":                   cspDomain.Domain,
		"subdomains":               cspDomain.Subdomains,
		"referenceId":              cspDomain.ReferenceID,
		"applyToAllOnboardedPaths": cspDomain.ApplyToAllOnboardedPaths,
	}
	m.writeJSONResponse(w, response)
}

// handleCSPAddDomain adds a new pre-approved domain
// Request body: ShallowPreApprovedDomain, Response: PreApprovedDomain
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
func (m *MockImpervaServer) handleCSPAddDomain(w http.ResponseWriter, r *http.Request, siteID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Failed to read request body")
		return
	}

	var domainReq struct {
		Domain                   string `json:"domain"`
		Subdomains               bool   `json:"subdomains"`
		ApplyToAllOnboardedPaths bool   `json:"applyToAllOnboardedPaths"`
	}
	if err := json.Unmarshal(body, &domainReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Invalid JSON")
		return
	}

	if m.cspDomains[siteID] == nil {
		m.cspDomains[siteID] = make(map[string]*MockCSPDomain)
	}

	refID := base64.RawURLEncoding.EncodeToString([]byte(domainReq.Domain))
	cspDomain := &MockCSPDomain{
		Domain:                   domainReq.Domain,
		Subdomains:               domainReq.Subdomains,
		ReferenceID:              refID,
		ApplyToAllOnboardedPaths: domainReq.ApplyToAllOnboardedPaths,
		Notes:                    []MockCSPNote{},
		Status:                   &MockCSPStatus{},
	}
	m.cspDomains[siteID][domainReq.Domain] = cspDomain

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"domain":                   cspDomain.Domain,
		"subdomains":               cspDomain.Subdomains,
		"referenceId":              cspDomain.ReferenceID,
		"applyToAllOnboardedPaths": cspDomain.ApplyToAllOnboardedPaths,
	}
	m.writeJSONResponse(w, response)
}

// handleCSPDeleteDomain deletes a pre-approved domain
func (m *MockImpervaServer) handleCSPDeleteDomain(w http.ResponseWriter, siteID int, domainRef string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	domain := m.decodeDomainRef(domainRef)
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Invalid domain reference")
		return
	}

	siteDomains, exists := m.cspDomains[siteID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		m.writeErrorResponse(w, 404, "Site not found")
		return
	}

	if _, exists := siteDomains[domain]; !exists {
		w.WriteHeader(http.StatusNotFound)
		m.writeErrorResponse(w, 404, "Domain not found")
		return
	}

	delete(siteDomains, domain)
	w.WriteHeader(http.StatusNoContent)
}

// handleCSPDomainStatus handles GET/PUT for domain status
func (m *MockImpervaServer) handleCSPDomainStatus(w http.ResponseWriter, r *http.Request, siteID int, domainRef string) {
	domain := m.decodeDomainRef(domainRef)
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Invalid domain reference")
		return
	}

	switch r.Method {
	case http.MethodGet:
		m.handleCSPGetDomainStatus(w, siteID, domain)
	case http.MethodPut:
		m.handleCSPUpdateDomainStatus(w, r, siteID, domain)
	default:
		m.writeErrorResponse(w, 9999, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// handleCSPGetDomainStatus returns the authorization status of a domain
// Returns AuthorizationStatus object as per API documentation
// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
func (m *MockImpervaServer) handleCSPGetDomainStatus(w http.ResponseWriter, siteID int, domain string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	siteDomains, exists := m.cspDomains[siteID]
	if !exists {
		m.writeJSONResponse(w, map[string]interface{}{})
		return
	}

	cspDomain, exists := siteDomains[domain]
	if !exists {
		m.writeJSONResponse(w, map[string]interface{}{})
		return
	}

	if cspDomain.Status == nil {
		m.writeJSONResponse(w, map[string]interface{}{})
		return
	}

	response := map[string]interface{}{}
	if cspDomain.Status.Blocked != nil {
		response["blocked"] = *cspDomain.Status.Blocked
	}
	if cspDomain.Status.Reviewed != nil {
		response["reviewed"] = *cspDomain.Status.Reviewed
	}
	if cspDomain.Status.Note != "" {
		response["note"] = cspDomain.Status.Note
	}
	if cspDomain.Status.Author != "" {
		response["author"] = cspDomain.Status.Author
	}
	if cspDomain.Status.ReviewedAt != 0 {
		response["reviewedAt"] = cspDomain.Status.ReviewedAt
	}
	if cspDomain.Status.LastNoteAt != 0 {
		response["lastNoteAt"] = cspDomain.Status.LastNoteAt
	}
	response["forceChange"] = cspDomain.Status.ForceChange

	m.writeJSONResponse(w, response)
}

// handleCSPUpdateDomainStatus updates the status of a domain
func (m *MockImpervaServer) handleCSPUpdateDomainStatus(w http.ResponseWriter, r *http.Request, siteID int, domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Failed to read request body")
		return
	}

	var statusReq MockCSPStatus
	if err := json.Unmarshal(body, &statusReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Invalid JSON")
		return
	}

	if m.cspDomains[siteID] == nil {
		m.cspDomains[siteID] = make(map[string]*MockCSPDomain)
	}

	cspDomain, exists := m.cspDomains[siteID][domain]
	if !exists {
		refID := base64.RawURLEncoding.EncodeToString([]byte(domain))
		cspDomain = &MockCSPDomain{
			Domain:      domain,
			ReferenceID: refID,
			Notes:       []MockCSPNote{},
			Status:      &MockCSPStatus{},
		}
		m.cspDomains[siteID][domain] = cspDomain
	}

	if cspDomain.Status == nil {
		cspDomain.Status = &MockCSPStatus{}
	}

	// Update all AuthorizationStatus fields from request
	// See: https://docs-cybersec-be.thalesgroup.com/api/bundle/api-docs/page/csp-api-definition.htm
	if statusReq.Blocked != nil {
		cspDomain.Status.Blocked = statusReq.Blocked
	}
	if statusReq.Reviewed != nil {
		cspDomain.Status.Reviewed = statusReq.Reviewed
	}
	if statusReq.Note != "" {
		cspDomain.Status.Note = statusReq.Note
	}
	if statusReq.Author != "" {
		cspDomain.Status.Author = statusReq.Author
	}
	if statusReq.ReviewedAt != 0 {
		cspDomain.Status.ReviewedAt = statusReq.ReviewedAt
	}
	if statusReq.LastNoteAt != 0 {
		cspDomain.Status.LastNoteAt = statusReq.LastNoteAt
	}
	cspDomain.Status.ForceChange = statusReq.ForceChange

	// Build response with all AuthorizationStatus fields
	response := map[string]interface{}{}
	if cspDomain.Status.Blocked != nil {
		response["blocked"] = *cspDomain.Status.Blocked
	}
	if cspDomain.Status.Reviewed != nil {
		response["reviewed"] = *cspDomain.Status.Reviewed
	}
	if cspDomain.Status.Note != "" {
		response["note"] = cspDomain.Status.Note
	}
	if cspDomain.Status.Author != "" {
		response["author"] = cspDomain.Status.Author
	}
	if cspDomain.Status.ReviewedAt != 0 {
		response["reviewedAt"] = cspDomain.Status.ReviewedAt
	}
	if cspDomain.Status.LastNoteAt != 0 {
		response["lastNoteAt"] = cspDomain.Status.LastNoteAt
	}
	response["forceChange"] = cspDomain.Status.ForceChange
	m.writeJSONResponse(w, response)
}

// handleCSPDomainNotes handles GET/POST/DELETE for domain notes
func (m *MockImpervaServer) handleCSPDomainNotes(w http.ResponseWriter, r *http.Request, siteID int, domainRef string) {
	domain := m.decodeDomainRef(domainRef)
	if domain == "" {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Invalid domain reference")
		return
	}

	switch r.Method {
	case http.MethodGet:
		m.handleCSPGetDomainNotes(w, siteID, domain)
	case http.MethodPost:
		m.handleCSPAddDomainNote(w, r, siteID, domain)
	case http.MethodDelete:
		m.handleCSPDeleteDomainNotes(w, siteID, domain)
	default:
		m.writeErrorResponse(w, 9999, fmt.Sprintf("Method not allowed: %s", r.Method))
	}
}

// handleCSPGetDomainNotes returns all notes for a domain
func (m *MockImpervaServer) handleCSPGetDomainNotes(w http.ResponseWriter, siteID int, domain string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	siteDomains, exists := m.cspDomains[siteID]
	if !exists {
		m.writeJSONResponse(w, []MockCSPNote{})
		return
	}

	cspDomain, exists := siteDomains[domain]
	if !exists {
		m.writeJSONResponse(w, []MockCSPNote{})
		return
	}

	m.writeJSONResponse(w, cspDomain.Notes)
}

// handleCSPAddDomainNote adds a note to a domain
func (m *MockImpervaServer) handleCSPAddDomainNote(w http.ResponseWriter, r *http.Request, siteID int, domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		m.writeErrorResponse(w, 400, "Failed to read request body")
		return
	}

	noteText := string(body)

	if m.cspDomains[siteID] == nil {
		m.cspDomains[siteID] = make(map[string]*MockCSPDomain)
	}

	cspDomain, exists := m.cspDomains[siteID][domain]
	if !exists {
		refID := base64.RawURLEncoding.EncodeToString([]byte(domain))
		cspDomain = &MockCSPDomain{
			Domain:      domain,
			ReferenceID: refID,
			Notes:       []MockCSPNote{},
			Status:      &MockCSPStatus{},
		}
		m.cspDomains[siteID][domain] = cspDomain
	}

	note := MockCSPNote{
		Text:   noteText,
		Author: "mock-user@example.com",
		Date:   time.Now().Unix(),
	}
	cspDomain.Notes = append(cspDomain.Notes, note)

	w.WriteHeader(http.StatusCreated)
	m.writeJSONResponse(w, cspDomain.Notes)
}

// handleCSPDeleteDomainNotes deletes all notes from a domain
func (m *MockImpervaServer) handleCSPDeleteDomainNotes(w http.ResponseWriter, siteID int, domain string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	siteDomains, exists := m.cspDomains[siteID]
	if !exists {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cspDomain, exists := siteDomains[domain]
	if !exists {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cspDomain.Notes = []MockCSPNote{}
	w.WriteHeader(http.StatusNoContent)
}

// decodeDomainRef decodes a base64 URL-encoded domain reference
func (m *MockImpervaServer) decodeDomainRef(domainRef string) string {
	decoded, err := base64.RawURLEncoding.DecodeString(domainRef)
	if err != nil {
		return ""
	}
	return string(decoded)
}

// Helper methods for tests

// GetAccount returns an account by ID (for test assertions)
func (m *MockImpervaServer) GetAccount(accountID int) *MockAccount {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.accounts[accountID]
}

// GetSite returns a site by ID (for test assertions)
func (m *MockImpervaServer) GetSite(siteID int) *MockSite {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sites[siteID]
}

// AddAccount adds an account directly (for test setup)
func (m *MockImpervaServer) AddAccount(account *MockAccount) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if account.AccountID == 0 {
		account.AccountID = m.nextAccountID
		m.nextAccountID++
	}
	m.accounts[account.AccountID] = account
}

// AddSite adds a site directly (for test setup)
func (m *MockImpervaServer) AddSite(site *MockSite) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if site.SiteID == 0 {
		site.SiteID = m.nextSiteID
		m.nextSiteID++
	}
	m.sites[site.SiteID] = site
}

// Reset clears all data from the mock server
func (m *MockImpervaServer) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.accounts = make(map[int]*MockAccount)
	m.sites = make(map[int]*MockSite)
	m.cspDomains = make(map[int]map[string]*MockCSPDomain)
	m.nextAccountID = 1000
	m.nextSiteID = 10000
}

// GetCSPDomain returns a CSP domain by site ID and domain name (for test assertions)
func (m *MockImpervaServer) GetCSPDomain(siteID int, domain string) *MockCSPDomain {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if siteDomains, exists := m.cspDomains[siteID]; exists {
		return siteDomains[domain]
	}
	return nil
}

// AddCSPDomain adds a CSP domain directly (for test setup)
func (m *MockImpervaServer) AddCSPDomain(siteID int, domain *MockCSPDomain) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cspDomains[siteID] == nil {
		m.cspDomains[siteID] = make(map[string]*MockCSPDomain)
	}
	if domain.ReferenceID == "" {
		domain.ReferenceID = base64.RawURLEncoding.EncodeToString([]byte(domain.Domain))
	}
	m.cspDomains[siteID][domain.Domain] = domain
}

// Suppress unused import warning
var _ = ioutil.ReadAll
