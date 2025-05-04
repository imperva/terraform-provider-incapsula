package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	apiID    = "foo"
	apiKey   = "bar"
	siteId   = "123"
	endpoint = fmt.Sprintf(endpointShortRenewalCycleConfiguration, siteId)
)

func TestClient_GetShortRenewalCycleConfigurationConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.Method == http.MethodGet && req.URL.String() == endpoint {
			rw.Write([]byte(`{"data": [{"shortRenewalCycle": false}]}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	response, err := client.GetShortRenewalCycleConfiguration(siteId, "")
	if err != nil {
		t.Errorf("unexcpected error: %s", err)
	}

	var shortRenewalCycle = response.Data[0].ShortRenewalCycle
	if shortRenewalCycle {
		t.Errorf("expected to get shortRenewalCycle = false but got %t", shortRenewalCycle)
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		t.Errorf("expected no errors, got %s", response.Errors[0].Detail)
	}
}

func TestClient_EnableShortRenewalCycleConfigurationConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.Method == http.MethodPost && req.URL.String() == endpoint {
			rw.Write([]byte(`{"data": [{"shortRenewalCycle": true}]}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	response, err := client.EnableShortRenewalCycleConfiguration(siteId, "")
	if err != nil {
		t.Errorf("unexcpected error: %s", err)
	}

	var shortRenewalCycle = response.Data[0].ShortRenewalCycle
	if !shortRenewalCycle {
		t.Errorf("expected to get shortRenewalCycle = true but got %t", shortRenewalCycle)
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		t.Errorf("expected no errors, got %s", response.Errors[0].Detail)
	}
}

func TestDeleteShortRenewalCycleConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.Method == http.MethodGet && req.URL.String() == endpoint {
			rw.Write([]byte(`{"data": [{"shortRenewalCycle": false}]}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	response, err := client.DeleteShortRenewalCycleConfiguration(siteId, "")
	if err != nil {
		t.Errorf("unexcpected error: %s", err)
	}

	var shortRenewalCycle = response.Data[0].ShortRenewalCycle
	if shortRenewalCycle {
		t.Errorf("expected to get shortRenewalCycle = false but got %t", shortRenewalCycle)
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		t.Errorf("expected no errors, got %s", response.Errors[0].Detail)
	}
}

func TestDeleteShortRenewalCycleConfigurationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.Method == http.MethodDelete && req.URL.String() == endpoint {
			rw.Write([]byte(`{
    "errors": [
        {
            "status": 404,
            "id": "850c4b95fffac0e1",
            "source": {
                "pointer": "/v3/sites/1027036767/certificate/shortRenewalCycle"
            },
            "title": "Not Found",
            "detail": "account 1111 is not allowed to manage short renewal cycle configuration"
        }
    ]
}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.DeleteShortRenewalCycleConfiguration(siteId, "")
	if err == nil {
		t.Errorf("expected to get error")
	}
	var expectedError = "[ERROR] bad status code 404 from Incapsula service when deleting short renewal cycle configuration. site id: 123"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected to get error: %s", expectedError)
	}
}

func TestEnableShortRenewalCycleConfigurationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.Method == http.MethodPost && req.URL.String() == endpoint {
			rw.Write([]byte(`{
    "errors": [
        {
            "status": 404,
            "id": "850c4b95fffac0e1",
            "source": {
                "pointer": "/v3/sites/1027036767/certificate/shortRenewalCycle"
            },
            "title": "Not Found",
            "detail": "account 1111 is not allowed to manage short renewal cycle configuration"
        }
    ]
}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.EnableShortRenewalCycleConfiguration(siteId, "")
	if err == nil {
		t.Errorf("expected to get error")
	}
	var expectedError = "[ERROR] bad status code 404 from Incapsula service when enabling short renewal cycle configuration. site id: 123"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected to get error: %s", expectedError)
	}
}
