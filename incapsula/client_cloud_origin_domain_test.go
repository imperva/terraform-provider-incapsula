package incapsula

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var validCreateResponse = `{
	"data": [{
		"id": 12345,
		"siteId": 1,
		"originDomain": "api.example.com",
		"region": "us-east-1",
		"impervaOriginDomain": "api-example-com.c123.imperva.com",
		"originConfig": {"port": 443},
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-01T00:00:00Z"
	}]
}`

var validGetResponse = `{
	"data": [{
		"id": 12345,
		"siteId": 1,
		"originDomain": "api.example.com",
		"region": "us-east-1",
		"impervaOriginDomain": "api-example-com.c123.imperva.com",
		"originConfig": {"port": 443},
		"createdAt": "2024-01-01T00:00:00Z",
		"updatedAt": "2024-01-02T00:00:00Z"
	}]
}`

var errorResponse = `{"errors": [{"status": 400, "detail": "Bad request"}]}`

func TestClientCloudOriginDomainCreate(t *testing.T) {
	tests := map[string]struct {
		statusCode     int
		responseBody   string
		expectedErr    bool
		expectedErrMsg string
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
			responseBody:   errorResponse,
			expectedErr:    true,
			expectedErrMsg: "Error status code 400",
		},
		"ValidResponse": {
			statusCode:   201,
			responseBody: validCreateResponse,
			expectedErr:  false,
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

			response, err := client.CreateCloudOriginDomain(1, "", "api.example.com", "us-east-1", 443, "TLS_1_2")

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
				if len(response.Data) == 0 {
					t.Fatal("Expected data in response, got empty")
				}
				if response.Data[0].ID != 12345 {
					t.Errorf("Expected ID 12345, got %d", response.Data[0].ID)
				}
				if response.Data[0].OriginDomain != "api.example.com" {
					t.Errorf("Expected OriginDomain api.example.com, got %s", response.Data[0].OriginDomain)
				}
			}
		})
	}
}

func TestClientCloudOriginDomainGet(t *testing.T) {
	tests := map[string]struct {
		statusCode     int
		responseBody   string
		expectedErr    bool
		expectedErrMsg string
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
			responseBody:   `{"errors": [{"status": 404, "detail": "Not found"}]}`,
			expectedErr:    true,
			expectedErrMsg: "Error status code 404",
		},
		"ValidResponse": {
			statusCode:   200,
			responseBody: validGetResponse,
			expectedErr:  false,
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

			response, err := client.GetCloudOriginDomain(1, 12345, "")

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
				if len(response.Data) == 0 {
					t.Fatal("Expected data in response, got empty")
				}
				if response.Data[0].Region != "us-east-1" {
					t.Errorf("Expected Region us-east-1, got %s", response.Data[0].Region)
				}
			}
		})
	}
}

func TestClientCloudOriginDomainDelete(t *testing.T) {
	tests := map[string]struct {
		statusCode     int
		responseBody   string
		expectedErr    bool
		expectedErrMsg string
	}{
		"BadConnection": {
			statusCode:     0,
			responseBody:   "",
			expectedErr:    true,
			expectedErrMsg: "Error from Incapsula service",
		},
		"InvalidResponse": {
			statusCode:     404,
			responseBody:   `{"errors": [{"status": 404, "detail": "Not found"}]}`,
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

			err := client.DeleteCloudOriginDomain(1, 12345, "")

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
	createReq := CloudOriginDomainCreateRequest{
		OriginDomain: "api.example.com",
		Region:       "us-east-1",
		DomainConfig: &CloudOriginDomainConfig{Port: 443},
	}

	data, err := json.Marshal(createReq)
	if err != nil {
		t.Fatalf("Failed to marshal create request: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"originDomain"`) {
		t.Errorf("Expected originDomain field, got %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"domainConfig"`) {
		t.Errorf("Expected domainConfig field, got %s", jsonStr)
	}

	var unmarshaled CloudOriginDomainCreateRequest
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal create request: %v", err)
	}

	if unmarshaled.OriginDomain != "api.example.com" {
		t.Errorf("Expected domain api.example.com, got %s", unmarshaled.OriginDomain)
	}
	if unmarshaled.Region != "us-east-1" {
		t.Errorf("Expected region us-east-1, got %s", unmarshaled.Region)
	}
	if unmarshaled.DomainConfig.Port != 443 {
		t.Errorf("Expected port 443, got %d", unmarshaled.DomainConfig.Port)
	}
}
