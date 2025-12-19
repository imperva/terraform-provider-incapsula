package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestMockServerAccountLifecycle(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	// Test account creation
	resp, err := http.PostForm(mock.URL()+"/accounts/add", url.Values{
		"email":        {"test@example.com"},
		"account_name": {"Test Account"},
		"parent_id":    {"0"},
	})
	if err != nil {
		t.Fatalf("Failed to create account: %v", err)
	}
	defer resp.Body.Close()

	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	if createResp["res"].(float64) != 0 {
		t.Errorf("Expected res=0, got %v", createResp["res"])
	}

	account := createResp["Account"].(map[string]interface{})
	accountID := int(account["account_id"].(float64))

	if accountID < 1000 {
		t.Errorf("Expected account_id >= 1000, got %d", accountID)
	}
	if account["email"].(string) != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got %s", account["email"])
	}

	// Test account status
	resp, err = http.PostForm(mock.URL()+"/account", url.Values{
		"account_id": {fmt.Sprintf("%d", accountID)},
	})
	if err != nil {
		t.Fatalf("Failed to get account status: %v", err)
	}
	defer resp.Body.Close()

	// Test account delete
	resp, err = http.PostForm(mock.URL()+"/accounts/delete", url.Values{
		"account_id": {fmt.Sprintf("%d", accountID)},
	})
	if err != nil {
		t.Fatalf("Failed to delete account: %v", err)
	}
	defer resp.Body.Close()

	var deleteResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		t.Fatalf("Failed to decode delete response: %v", err)
	}

	if deleteResp["res"].(float64) != 0 {
		t.Errorf("Expected res=0, got %v", deleteResp["res"])
	}

	// Verify account is gone
	if mock.GetAccount(accountID) != nil {
		t.Errorf("Account should have been deleted")
	}
}

func TestMockServerSiteLifecycle(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	// Test site creation
	resp, err := http.PostForm(mock.URL()+"/sites/add", url.Values{
		"domain":     {"example.com"},
		"account_id": {"1000"},
		"site_type":  {"api"},
	})
	if err != nil {
		t.Fatalf("Failed to create site: %v", err)
	}
	defer resp.Body.Close()

	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	if createResp["res"].(float64) != 0 {
		t.Errorf("Expected res=0, got %v", createResp["res"])
	}

	siteID := int(createResp["site_id"].(float64))
	if siteID < 10000 {
		t.Errorf("Expected site_id >= 10000, got %d", siteID)
	}

	if createResp["domain"].(string) != "example.com" {
		t.Errorf("Expected domain='example.com', got %s", createResp["domain"])
	}

	// Verify DNS records in response
	dns := createResp["dns"].([]interface{})
	if len(dns) != 2 {
		t.Errorf("Expected 2 DNS records, got %d", len(dns))
	}

	// Verify site was stored
	site := mock.GetSite(siteID)
	if site == nil {
		t.Fatalf("Site should exist in mock server")
	}
	if site.Domain != "example.com" {
		t.Errorf("Expected domain='example.com', got %s", site.Domain)
	}

	// Test site status
	resp, err = http.PostForm(mock.URL()+"/sites/status", url.Values{
		"site_id": {fmt.Sprintf("%d", siteID)},
	})
	if err != nil {
		t.Fatalf("Failed to get site status: %v", err)
	}
	defer resp.Body.Close()

	// Test site delete
	resp, err = http.PostForm(mock.URL()+"/sites/delete", url.Values{
		"site_id": {fmt.Sprintf("%d", siteID)},
	})
	if err != nil {
		t.Fatalf("Failed to delete site: %v", err)
	}
	defer resp.Body.Close()

	var deleteResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&deleteResp); err != nil {
		t.Fatalf("Failed to decode delete response: %v", err)
	}

	if deleteResp["res"].(float64) != 0 {
		t.Errorf("Expected res=0, got %v", deleteResp["res"])
	}

	// Verify site is gone
	if mock.GetSite(siteID) != nil {
		t.Errorf("Site should have been deleted")
	}
}

func TestMockServerUnknownEndpoint(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	resp, err := http.Post(mock.URL()+"/unknown/endpoint", "application/x-www-form-urlencoded", strings.NewReader(""))
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	defer resp.Body.Close()

	var errorResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResp["res"].(float64) != 9999 {
		t.Errorf("Expected res=9999 for unknown endpoint, got %v", errorResp["res"])
	}
}

func TestMockServerReset(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	// Add an account
	mock.AddAccount(&MockAccount{Email: "test@example.com"})

	// Verify account exists
	if len(mock.accounts) != 1 {
		t.Errorf("Expected 1 account, got %d", len(mock.accounts))
	}

	// Reset
	mock.Reset()

	// Verify account is gone
	if len(mock.accounts) != 0 {
		t.Errorf("Expected 0 accounts after reset, got %d", len(mock.accounts))
	}
}

func TestMockServerHelperMethods(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	// Test AddAccount with auto-generated ID
	account := &MockAccount{
		Email:       "test@example.com",
		AccountName: "Test Account",
	}
	mock.AddAccount(account)

	if account.AccountID < 1000 {
		t.Errorf("Expected auto-generated account ID >= 1000, got %d", account.AccountID)
	}

	// Test GetAccount
	retrieved := mock.GetAccount(account.AccountID)
	if retrieved == nil {
		t.Fatalf("Expected to retrieve account")
	}
	if retrieved.Email != "test@example.com" {
		t.Errorf("Expected email='test@example.com', got %s", retrieved.Email)
	}

	// Test AddSite with auto-generated ID
	site := &MockSite{
		Domain:    "example.com",
		AccountID: account.AccountID,
	}
	mock.AddSite(site)

	if site.SiteID < 10000 {
		t.Errorf("Expected auto-generated site ID >= 10000, got %d", site.SiteID)
	}

	// Test GetSite
	retrievedSite := mock.GetSite(site.SiteID)
	if retrievedSite == nil {
		t.Fatalf("Expected to retrieve site")
	}
	if retrievedSite.Domain != "example.com" {
		t.Errorf("Expected domain='example.com', got %s", retrievedSite.Domain)
	}
}

func TestMockServerCSPDomainLifecycle(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	siteID := 12345
	domain := "example.com"
	domainRef := "ZXhhbXBsZS5jb20" // base64 URL-safe encoding of "example.com"

	// Test add domain
	domainJSON := `{"domain":"example.com","subdomains":true}`
	resp, err := http.Post(
		fmt.Sprintf("%s/csp-api/v1/sites/%d/preapprovedlist", mock.URL(), siteID),
		"application/json",
		strings.NewReader(domainJSON),
	)
	if err != nil {
		t.Fatalf("Failed to add domain: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	var createResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&createResp); err != nil {
		t.Fatalf("Failed to decode create response: %v", err)
	}

	if createResp["domain"].(string) != domain {
		t.Errorf("Expected domain='example.com', got %s", createResp["domain"])
	}
	if !createResp["subdomains"].(bool) {
		t.Errorf("Expected subdomains=true")
	}

	// Test get domain
	resp, err = http.Get(fmt.Sprintf("%s/csp-api/v1/sites/%d/preapprovedlist/%s", mock.URL(), siteID, domainRef))
	if err != nil {
		t.Fatalf("Failed to get domain: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test list all domains
	resp, err = http.Get(fmt.Sprintf("%s/csp-api/v1/sites/%d/preapprovedlist", mock.URL(), siteID))
	if err != nil {
		t.Fatalf("Failed to list domains: %v", err)
	}
	defer resp.Body.Close()

	var listResp []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		t.Fatalf("Failed to decode list response: %v", err)
	}

	if len(listResp) != 1 {
		t.Errorf("Expected 1 domain, got %d", len(listResp))
	}

	// Test delete domain
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/csp-api/v1/sites/%d/preapprovedlist/%s", mock.URL(), siteID, domainRef), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete domain: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}

	// Verify domain is deleted
	if mock.GetCSPDomain(siteID, domain) != nil {
		t.Errorf("Domain should have been deleted")
	}
}

func TestMockServerCSPDomainStatus(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	siteID := 12345
	domain := "example.com"
	domainRef := "ZXhhbXBsZS5jb20"

	// First add the domain
	mock.AddCSPDomain(siteID, &MockCSPDomain{
		Domain:     domain,
		Subdomains: true,
	})

	// Test update status
	statusJSON := `{"blocked":true,"reviewed":true}`
	req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/csp-api/v1/sites/%d/domains/%s/status", mock.URL(), siteID, domainRef), strings.NewReader(statusJSON))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to update status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test get status
	resp, err = http.Get(fmt.Sprintf("%s/csp-api/v1/sites/%d/domains/%s/status", mock.URL(), siteID, domainRef))
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	defer resp.Body.Close()

	var statusResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		t.Fatalf("Failed to decode status response: %v", err)
	}

	if !statusResp["blocked"].(bool) {
		t.Errorf("Expected blocked=true")
	}
	if !statusResp["reviewed"].(bool) {
		t.Errorf("Expected reviewed=true")
	}
}

func TestMockServerCSPDomainNotes(t *testing.T) {
	mock := NewMockImpervaServer()
	defer mock.Close()

	siteID := 12345
	domain := "example.com"
	domainRef := "ZXhhbXBsZS5jb20"

	// First add the domain
	mock.AddCSPDomain(siteID, &MockCSPDomain{
		Domain:     domain,
		Subdomains: true,
	})

	// Test add note
	noteText := "This is a test note"
	resp, err := http.Post(
		fmt.Sprintf("%s/csp-api/v1/sites/%d/domains/%s/notes", mock.URL(), siteID, domainRef),
		"text/plain",
		strings.NewReader(noteText),
	)
	if err != nil {
		t.Fatalf("Failed to add note: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.StatusCode)
	}

	// Test get notes
	resp, err = http.Get(fmt.Sprintf("%s/csp-api/v1/sites/%d/domains/%s/notes", mock.URL(), siteID, domainRef))
	if err != nil {
		t.Fatalf("Failed to get notes: %v", err)
	}
	defer resp.Body.Close()

	var notesResp []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&notesResp); err != nil {
		t.Fatalf("Failed to decode notes response: %v", err)
	}

	if len(notesResp) != 1 {
		t.Errorf("Expected 1 note, got %d", len(notesResp))
	}
	if notesResp[0]["text"].(string) != noteText {
		t.Errorf("Expected note text='%s', got %s", noteText, notesResp[0]["text"])
	}

	// Test delete notes
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/csp-api/v1/sites/%d/domains/%s/notes", mock.URL(), siteID, domainRef), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete notes: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", resp.StatusCode)
	}

	// Verify notes are deleted
	cspDomain := mock.GetCSPDomain(siteID, domain)
	if len(cspDomain.Notes) != 0 {
		t.Errorf("Expected 0 notes, got %d", len(cspDomain.Notes))
	}
}
