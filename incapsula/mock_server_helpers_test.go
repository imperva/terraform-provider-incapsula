package incapsula

import (
	"os"
	"testing"
)

func TestSetupMockEnvironment(t *testing.T) {
	// Save current state
	originalAPIID, wasSet := os.LookupEnv("INCAPSULA_API_ID")

	mock := NewMockImpervaServer()
	defer mock.Close()

	// Setup mock environment
	env := SetupMockEnvironment(mock.URL())

	// Verify environment variables are set
	if os.Getenv("INCAPSULA_API_ID") != "mock-api-id" {
		t.Errorf("Expected INCAPSULA_API_ID='mock-api-id', got %s", os.Getenv("INCAPSULA_API_ID"))
	}
	if os.Getenv("INCAPSULA_API_KEY") != "mock-api-key" {
		t.Errorf("Expected INCAPSULA_API_KEY='mock-api-key', got %s", os.Getenv("INCAPSULA_API_KEY"))
	}
	if os.Getenv("INCAPSULA_BASE_URL") != mock.URL() {
		t.Errorf("Expected INCAPSULA_BASE_URL='%s', got %s", mock.URL(), os.Getenv("INCAPSULA_BASE_URL"))
	}
	if os.Getenv("INCAPSULA_BASE_URL_API") != mock.URL() {
		t.Errorf("Expected INCAPSULA_BASE_URL_API='%s', got %s", mock.URL(), os.Getenv("INCAPSULA_BASE_URL_API"))
	}

	// Restore environment
	env.Restore()

	// Verify original values are restored
	currentAPIID, isSet := os.LookupEnv("INCAPSULA_API_ID")
	if wasSet != isSet {
		t.Errorf("INCAPSULA_API_ID set state changed: was %v, now %v", wasSet, isSet)
	}
	if wasSet && currentAPIID != originalAPIID {
		t.Errorf("INCAPSULA_API_ID value changed: was %s, now %s", originalAPIID, currentAPIID)
	}
}

func TestWithMockServer(t *testing.T) {
	testExecuted := false

	WithMockServer(t, func(ctx *MockTestContext) {
		testExecuted = true

		// Verify server is running
		if ctx.Server == nil {
			t.Errorf("Expected Server to be set")
		}
		if ctx.Env == nil {
			t.Errorf("Expected Env to be set")
		}

		// Verify environment is set up
		if os.Getenv("INCAPSULA_API_ID") != "mock-api-id" {
			t.Errorf("Expected environment to be set up")
		}
	})

	if !testExecuted {
		t.Errorf("Test function was not executed")
	}
}

func TestMockTestContextHelpers(t *testing.T) {
	WithMockServer(t, func(ctx *MockTestContext) {
		// Test CreateTestAccount
		account := ctx.CreateTestAccount()
		if account == nil {
			t.Fatalf("Expected account to be created")
		}
		if account.AccountID < 1000 {
			t.Errorf("Expected account ID >= 1000, got %d", account.AccountID)
		}

		// Verify account is in server
		retrieved := ctx.Server.GetAccount(account.AccountID)
		if retrieved == nil {
			t.Errorf("Expected account to be stored in server")
		}

		// Test CreateTestSite
		site := ctx.CreateTestSite(account.AccountID)
		if site == nil {
			t.Fatalf("Expected site to be created")
		}
		if site.SiteID < 10000 {
			t.Errorf("Expected site ID >= 10000, got %d", site.SiteID)
		}
		if site.AccountID != account.AccountID {
			t.Errorf("Expected site account ID=%d, got %d", account.AccountID, site.AccountID)
		}

		// Verify site is in server
		retrievedSite := ctx.Server.GetSite(site.SiteID)
		if retrievedSite == nil {
			t.Errorf("Expected site to be stored in server")
		}

		// Test CreateTestCSPDomain
		domain := ctx.CreateTestCSPDomain(site.SiteID, "test-csp.example.com")
		if domain == nil {
			t.Fatalf("Expected CSP domain to be created")
		}
		if domain.Domain != "test-csp.example.com" {
			t.Errorf("Expected domain='test-csp.example.com', got %s", domain.Domain)
		}

		// Verify CSP domain is in server
		retrievedCSP := ctx.Server.GetCSPDomain(site.SiteID, "test-csp.example.com")
		if retrievedCSP == nil {
			t.Errorf("Expected CSP domain to be stored in server")
		}
	})
}

func TestShouldUseMockServer(t *testing.T) {
	// Save original values
	origAPIID, wasAPIIDSet := os.LookupEnv("INCAPSULA_API_ID")
	origUseMock, wasUseMockSet := os.LookupEnv("USE_MOCK_SERVER")
	defer func() {
		if wasAPIIDSet {
			os.Setenv("INCAPSULA_API_ID", origAPIID)
		} else {
			os.Unsetenv("INCAPSULA_API_ID")
		}
		if wasUseMockSet {
			os.Setenv("USE_MOCK_SERVER", origUseMock)
		} else {
			os.Unsetenv("USE_MOCK_SERVER")
		}
	}()

	// Test: No API ID, no USE_MOCK_SERVER -> should use mock
	os.Unsetenv("INCAPSULA_API_ID")
	os.Unsetenv("USE_MOCK_SERVER")
	if !ShouldUseMockServer() {
		t.Errorf("Expected ShouldUseMockServer()=true when no API ID set")
	}

	// Test: API ID set, no USE_MOCK_SERVER -> should use real API
	os.Setenv("INCAPSULA_API_ID", "real-api-id")
	os.Unsetenv("USE_MOCK_SERVER")
	if ShouldUseMockServer() {
		t.Errorf("Expected ShouldUseMockServer()=false when API ID is set")
	}

	// Test: API ID set, USE_MOCK_SERVER=true -> should use mock
	os.Setenv("INCAPSULA_API_ID", "real-api-id")
	os.Setenv("USE_MOCK_SERVER", "true")
	if !ShouldUseMockServer() {
		t.Errorf("Expected ShouldUseMockServer()=true when USE_MOCK_SERVER=true")
	}

	// Test: API ID set, USE_MOCK_SERVER=1 -> should use mock
	os.Setenv("USE_MOCK_SERVER", "1")
	if !ShouldUseMockServer() {
		t.Errorf("Expected ShouldUseMockServer()=true when USE_MOCK_SERVER=1")
	}
}
