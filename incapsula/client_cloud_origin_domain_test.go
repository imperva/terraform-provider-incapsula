package incapsula

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientCloudOriginDomainCreate(t *testing.T) {
	tests := map[string]struct {
		statusCode      int
		responseBody    string
		expectedErr     bool
		expectedErrMsg  string
	}{
		"BadConnection": {
			statusCode:     0,
			responseBody:   "",
			expectedErr:    true,
			expectedErrMsg: "Error from Incapsula service",
		},
		"BadJSON": {
			statusCode:     201,
			responseBody:   "invalid json",
			expectedErr:    true,
			expectedErrMsg: "Error parsing cloud origin domain JSON response",
		},
		"InvalidResponse": {
			statusCode:     400,
			responseBody:   `{"isError": true, "message": "Bad request"}`,
			expectedErr:    true,
			expectedErrMsg: "Error status code 400",
		},
		"ValidResponse": {
			statusCode: 201,
			responseBody: `{
				"value": {
					"originId": 12345,
					"domain": "api.example.com",
					"region": "us-east-1",
					"port": 443,
					"impervaOriginDomain": "api-example-com.c123.imperva.com",
					"status": "PENDING",
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-01T00:00:00Z"
				},
				"isError": false
			}`,
			expectedErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Skip test setup if it's a BadConnection case (won't use server)
			var baseURL string
			if test.statusCode != 0 {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(test.statusCode)
					w.Write([]byte(test.responseBody))
				}))
				defer server.Close()
				baseURL = server.URL
			} else {
				// Use invalid URL for connection error
				baseURL = "http://invalid.test:99999"
			}

			client := &Client{
				config: &Config{
					BaseURLRev3: baseURL,
					APIID:       "test_id",
					APIKey:      "test_key",
				},
				httpClient: &http.Client{},
			}

			response, err := client.CreateCloudOriginDomain(1, 0, "api.example.com", "us-east-1", 443)

			if test.expectedErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !test.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if test.expectedErr && err != nil && test.expectedErrMsg != "" {
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Errorf("Expected error containing %q, got %q", test.expectedErrMsg, err.Error())
				}
			}
			if !test.expectedErr && response == nil {
				t.Errorf("Expected response, got nil")
			}
			if !test.expectedErr && response != nil {
				if response.Value.OriginID != 12345 {
					t.Errorf("Expected OriginID 12345, got %d", response.Value.OriginID)
				}
				if response.Value.Domain != "api.example.com" {
					t.Errorf("Expected Domain api.example.com, got %s", response.Value.Domain)
				}
			}
		})
	}
}

func TestClientCloudOriginDomainGet(t *testing.T) {
	tests := map[string]struct {
		statusCode      int
		responseBody    string
		expectedErr     bool
		expectedErrMsg  string
	}{
		"BadConnection": {
			statusCode:     0,
			responseBody:   "",
			expectedErr:    true,
			expectedErrMsg: "Error from Incapsula service",
		},
		"BadJSON": {
			statusCode:     200,
			responseBody:   "invalid json",
			expectedErr:    true,
			expectedErrMsg: "Error parsing cloud origin domain JSON response",
		},
		"InvalidResponse": {
			statusCode:     404,
			responseBody:   `{"isError": true, "message": "Not found"}`,
			expectedErr:    true,
			expectedErrMsg: "Error status code 404",
		},
		"ValidResponse": {
			statusCode: 200,
			responseBody: `{
				"value": {
					"originId": 12345,
					"domain": "api.example.com",
					"region": "us-east-1",
					"port": 443,
					"impervaOriginDomain": "api-example-com.c123.imperva.com",
					"status": "ACTIVE",
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-02T00:00:00Z"
				},
				"isError": false
			}`,
			expectedErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var baseURL string
			if test.statusCode != 0 {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(test.statusCode)
					w.Write([]byte(test.responseBody))
				}))
				defer server.Close()
				baseURL = server.URL
			} else {
				baseURL = "http://invalid.test:99999"
			}

			client := &Client{
				config: &Config{
					BaseURLRev3: baseURL,
					APIID:       "test_id",
					APIKey:      "test_key",
				},
				httpClient: &http.Client{},
			}

			response, err := client.GetCloudOriginDomain(1, 0, 12345)

			if test.expectedErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !test.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if test.expectedErr && err != nil && test.expectedErrMsg != "" {
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Errorf("Expected error containing %q, got %q", test.expectedErrMsg, err.Error())
				}
			}
			if !test.expectedErr && response == nil {
				t.Errorf("Expected response, got nil")
			}
			if !test.expectedErr && response != nil {
				if response.Value.Status != "ACTIVE" {
					t.Errorf("Expected Status ACTIVE, got %s", response.Value.Status)
				}
			}
		})
	}
}

func TestClientCloudOriginDomainUpdate(t *testing.T) {
	tests := map[string]struct {
		statusCode      int
		responseBody    string
		expectedErr     bool
		expectedErrMsg  string
	}{
		"BadConnection": {
			statusCode:     0,
			responseBody:   "",
			expectedErr:    true,
			expectedErrMsg: "Error from Incapsula service",
		},
		"BadJSON": {
			statusCode:     200,
			responseBody:   "invalid json",
			expectedErr:    true,
			expectedErrMsg: "Error parsing cloud origin domain JSON response",
		},
		"InvalidResponse": {
			statusCode:     400,
			responseBody:   `{"isError": true, "message": "Bad request"}`,
			expectedErr:    true,
			expectedErrMsg: "Error status code 400",
		},
		"ValidResponse": {
			statusCode: 200,
			responseBody: `{
				"value": {
					"originId": 12345,
					"domain": "api.example.com",
					"region": "eu-west-1",
					"port": 8443,
					"impervaOriginDomain": "api-example-com.c123.imperva.com",
					"status": "ACTIVE",
					"createdAt": "2024-01-01T00:00:00Z",
					"updatedAt": "2024-01-03T00:00:00Z"
				},
				"isError": false
			}`,
			expectedErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var baseURL string
			if test.statusCode != 0 {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodPut {
						t.Errorf("Expected PUT method, got %s", r.Method)
					}
					w.WriteHeader(test.statusCode)
					w.Write([]byte(test.responseBody))
				}))
				defer server.Close()
				baseURL = server.URL
			} else {
				baseURL = "http://invalid.test:99999"
			}

			client := &Client{
				config: &Config{
					BaseURLRev3: baseURL,
					APIID:       "test_id",
					APIKey:      "test_key",
				},
				httpClient: &http.Client{},
			}

			response, err := client.UpdateCloudOriginDomain(1, 0, 12345, "eu-west-1", 8443)

			if test.expectedErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !test.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if test.expectedErr && err != nil && test.expectedErrMsg != "" {
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Errorf("Expected error containing %q, got %q", test.expectedErrMsg, err.Error())
				}
			}
			if !test.expectedErr && response == nil {
				t.Errorf("Expected response, got nil")
			}
			if !test.expectedErr && response != nil {
				if response.Value.Region != "eu-west-1" {
					t.Errorf("Expected Region eu-west-1, got %s", response.Value.Region)
				}
				if response.Value.Port != 8443 {
					t.Errorf("Expected Port 8443, got %d", response.Value.Port)
				}
			}
		})
	}
}

func TestClientCloudOriginDomainDelete(t *testing.T) {
	tests := map[string]struct {
		statusCode      int
		responseBody    string
		expectedErr     bool
		expectedErrMsg  string
	}{
		"BadConnection": {
			statusCode:     0,
			responseBody:   "",
			expectedErr:    true,
			expectedErrMsg: "Error from Incapsula service",
		},
		"InvalidResponse": {
			statusCode:     404,
			responseBody:   `{"isError": true, "message": "Not found"}`,
			expectedErr:    true,
			expectedErrMsg: "Error status code 404",
		},
		"ValidResponse": {
			statusCode:   204,
			responseBody: "",
			expectedErr:  false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var baseURL string
			if test.statusCode != 0 {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodDelete {
						t.Errorf("Expected DELETE method, got %s", r.Method)
					}
					w.WriteHeader(test.statusCode)
					w.Write([]byte(test.responseBody))
				}))
				defer server.Close()
				baseURL = server.URL
			} else {
				baseURL = "http://invalid.test:99999"
			}

			client := &Client{
				config: &Config{
					BaseURLRev3: baseURL,
					APIID:       "test_id",
					APIKey:      "test_key",
				},
				httpClient: &http.Client{},
			}

			err := client.DeleteCloudOriginDomain(1, 0, 12345)

			if test.expectedErr && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !test.expectedErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if test.expectedErr && err != nil && test.expectedErrMsg != "" {
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Errorf("Expected error containing %q, got %q", test.expectedErrMsg, err.Error())
				}
			}
		})
	}
}

func TestCloudOriginDomainJSONMarshaling(t *testing.T) {
	// Test that request structs marshal correctly
	createReq := CloudOriginDomainCreateRequest{
		Domain: "api.example.com",
		Region: "us-east-1",
		Port:   443,
	}

	data, err := json.Marshal(createReq)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}

	var unmarshaled CloudOriginDomainCreateRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal create request: %v", err)
	}

	if unmarshaled.Domain != "api.example.com" {
		t.Errorf("Expected domain api.example.com, got %s", unmarshaled.Domain)
	}
	if unmarshaled.Region != "us-east-1" {
		t.Errorf("Expected region us-east-1, got %s", unmarshaled.Region)
	}
	if unmarshaled.Port != 443 {
		t.Errorf("Expected port 443, got %d", unmarshaled.Port)
	}
}
