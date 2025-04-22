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
	endpoint = fmt.Sprintf(endpointFastRenewalConfiguration, siteId)
)

func TestGetFastRenewalConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.Method == http.MethodGet && req.URL.String() == endpoint {
			rw.Write([]byte(`{"data": [{"fastRenewal": false}]}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	response, err := client.GetFastRenewalConfiguration(siteId, "")
	if err != nil {
		t.Errorf("unexcpected error: %s", err)
	}

	var fastRenewal = response.Data[0].FastRenewal
	if fastRenewal {
		t.Errorf("expected to get fastRnewal = false but got %t", fastRenewal)
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		t.Errorf("expected no errors, got %s", response.Errors[0].Detail)
	}
}

func TestEnableFastRenewalConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.Method == http.MethodPost && req.URL.String() == endpoint {
			rw.Write([]byte(`{"data": [{"fastRenewal": true}]}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	response, err := client.EnableFastRenewalConfiguration(siteId, "")
	if err != nil {
		t.Errorf("unexcpected error: %s", err)
	}

	var fastRenewal = response.Data[0].FastRenewal
	if !fastRenewal {
		t.Errorf("expected to get fastRnewal = true but got %t", fastRenewal)
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		t.Errorf("expected no errors, got %s", response.Errors[0].Detail)
	}
}

func TestDeleteFastRenewalConfiguration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.Method == http.MethodGet && req.URL.String() == endpoint {
			rw.Write([]byte(`{"data": [{"fastRenewal": false}]}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	response, err := client.DeleteFastRenewalConfiguration(siteId, "")
	if err != nil {
		t.Errorf("unexcpected error: %s", err)
	}

	var fastRenewal = response.Data[0].FastRenewal
	if fastRenewal {
		t.Errorf("expected to get fastRnewal = false but got %t", fastRenewal)
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		t.Errorf("expected no errors, got %s", response.Errors[0].Detail)
	}
}

func TestDeleteFastRenewalConfigurationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.Method == http.MethodDelete && req.URL.String() == endpoint {
			rw.Write([]byte(`{
    "errors": [
        {
            "status": 404,
            "id": "850c4b95fffac0e1",
            "source": {
                "pointer": "/v3/sites/1027036767/certificate/fastRenewal"
            },
            "title": "Not Found",
            "detail": "account 1111 is not allowed to manage fast renewal configuration"
        }
    ]
}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.DeleteFastRenewalConfiguration(siteId, "")
	if err == nil {
		t.Errorf("expected to get error")
	}
	var expectedError = "[ERROR] bad status code 404 from Incapsula service when deleting fast renewal configuration. site ID: 123"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected to get error: %s", expectedError)
	}
}

func TestEnableFastRenewalConfigurationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.Method == http.MethodPost && req.URL.String() == endpoint {
			rw.Write([]byte(`{
    "errors": [
        {
            "status": 404,
            "id": "850c4b95fffac0e1",
            "source": {
                "pointer": "/v3/sites/1027036767/certificate/fastRenewal"
            },
            "title": "Not Found",
            "detail": "account 1111 is not allowed to manage fast renewal configuration"
        }
    ]
}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, err := client.EnableFastRenewalConfiguration(siteId, "")
	if err == nil {
		t.Errorf("expected to get error")
	}
	var expectedError = "[ERROR] bad status code 404 from Incapsula service when enabling fast renewal configuration. site ID: 123"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected to get error: %s", expectedError)
	}
}
