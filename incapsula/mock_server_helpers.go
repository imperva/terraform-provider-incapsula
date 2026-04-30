package incapsula

import (
	"os"
	"testing"
)

// MockEnvironment holds the original environment variables for restoration
type MockEnvironment struct {
	originalAPIID       string
	originalAPIKey      string
	originalBaseURL     string
	originalBaseURLRev2 string
	originalBaseURLRev3 string
	originalBaseURLAPI  string
	wasAPIIDSet         bool
	wasAPIKeySet        bool
	wasBaseURLSet       bool
	wasBaseURLRev2Set   bool
	wasBaseURLRev3Set   bool
	wasBaseURLAPISet    bool
}

// SetupMockEnvironment configures environment variables to point to the mock server
// Returns a MockEnvironment that can be used to restore the original environment
func SetupMockEnvironment(mockURL string) *MockEnvironment {
	env := &MockEnvironment{}

	// Save original values
	env.originalAPIID, env.wasAPIIDSet = os.LookupEnv("INCAPSULA_API_ID")
	env.originalAPIKey, env.wasAPIKeySet = os.LookupEnv("INCAPSULA_API_KEY")
	env.originalBaseURL, env.wasBaseURLSet = os.LookupEnv("INCAPSULA_BASE_URL")
	env.originalBaseURLRev2, env.wasBaseURLRev2Set = os.LookupEnv("INCAPSULA_BASE_URL_REV_2")
	env.originalBaseURLRev3, env.wasBaseURLRev3Set = os.LookupEnv("INCAPSULA_BASE_URL_REV_3")
	env.originalBaseURLAPI, env.wasBaseURLAPISet = os.LookupEnv("INCAPSULA_BASE_URL_API")

	// Set mock values
	os.Setenv("INCAPSULA_API_ID", "mock-api-id")
	os.Setenv("INCAPSULA_API_KEY", "mock-api-key")
	os.Setenv("INCAPSULA_BASE_URL", mockURL)
	os.Setenv("INCAPSULA_BASE_URL_REV_2", mockURL)
	os.Setenv("INCAPSULA_BASE_URL_REV_3", mockURL)
	os.Setenv("INCAPSULA_BASE_URL_API", mockURL)

	return env
}

// Restore restores the original environment variables
func (e *MockEnvironment) Restore() {
	if e.wasAPIIDSet {
		os.Setenv("INCAPSULA_API_ID", e.originalAPIID)
	} else {
		os.Unsetenv("INCAPSULA_API_ID")
	}

	if e.wasAPIKeySet {
		os.Setenv("INCAPSULA_API_KEY", e.originalAPIKey)
	} else {
		os.Unsetenv("INCAPSULA_API_KEY")
	}

	if e.wasBaseURLSet {
		os.Setenv("INCAPSULA_BASE_URL", e.originalBaseURL)
	} else {
		os.Unsetenv("INCAPSULA_BASE_URL")
	}

	if e.wasBaseURLRev2Set {
		os.Setenv("INCAPSULA_BASE_URL_REV_2", e.originalBaseURLRev2)
	} else {
		os.Unsetenv("INCAPSULA_BASE_URL_REV_2")
	}

	if e.wasBaseURLRev3Set {
		os.Setenv("INCAPSULA_BASE_URL_REV_3", e.originalBaseURLRev3)
	} else {
		os.Unsetenv("INCAPSULA_BASE_URL_REV_3")
	}

	if e.wasBaseURLAPISet {
		os.Setenv("INCAPSULA_BASE_URL_API", e.originalBaseURLAPI)
	} else {
		os.Unsetenv("INCAPSULA_BASE_URL_API")
	}
}

// MockTestContext provides a complete test context with mock server and environment
type MockTestContext struct {
	Server *MockImpervaServer
	Env    *MockEnvironment
}

// WithMockServer creates a mock server, sets up the environment, and executes the test function.
// The mock server and environment are automatically cleaned up after the test.
func WithMockServer(t *testing.T, testFunc func(ctx *MockTestContext)) {
	server := NewMockImpervaServer()
	defer server.Close()

	env := SetupMockEnvironment(server.URL())
	defer env.Restore()

	ctx := &MockTestContext{
		Server: server,
		Env:    env,
	}

	testFunc(ctx)
}

// CreateTestAccount creates a test account with default values and returns it
func (ctx *MockTestContext) CreateTestAccount() *MockAccount {
	account := &MockAccount{
		Email:       "test@example.com",
		AccountName: "Test Account",
		ParentID:    0,
		PlanID:      "enterprise",
	}
	ctx.Server.AddAccount(account)
	return account
}

// CreateTestSite creates a test site with default values and returns it
func (ctx *MockTestContext) CreateTestSite(accountID int) *MockSite {
	site := &MockSite{
		AccountID:  accountID,
		Domain:     "test.example.com",
		Status:     "active",
		SiteType:   "api",
		DnsARecord: "1.2.3.4",
		DnsCname:   "test.incapdns.net",
	}
	ctx.Server.AddSite(site)
	return site
}

// CreateTestCSPDomain creates a test CSP domain with default values and returns it
func (ctx *MockTestContext) CreateTestCSPDomain(siteID int, domain string) *MockCSPDomain {
	cspDomain := &MockCSPDomain{
		Domain:     domain,
		Subdomains: true,
		Notes:      []MockCSPNote{},
		Status:     &MockCSPStatus{},
	}
	ctx.Server.AddCSPDomain(siteID, cspDomain)
	return cspDomain
}

// ShouldUseMockServer returns true if tests should use the mock server
// This allows tests to be run against either real API or mock server
func ShouldUseMockServer() bool {
	_, hasAPIID := os.LookupEnv("INCAPSULA_API_ID")
	useMock := os.Getenv("USE_MOCK_SERVER")

	// Use mock if explicitly requested or if no API credentials are set
	return useMock == "true" || useMock == "1" || !hasAPIID
}

// SkipIfNoMockAndNoCredentials skips the test if neither mock server nor real credentials are available
func SkipIfNoMockAndNoCredentials(t *testing.T) {
	if !ShouldUseMockServer() {
		_, hasAPIID := os.LookupEnv("INCAPSULA_API_ID")
		_, hasAPIKey := os.LookupEnv("INCAPSULA_API_KEY")
		if !hasAPIID || !hasAPIKey {
			t.Skip("Skipping test: neither mock server enabled nor real API credentials available")
		}
	}
}
