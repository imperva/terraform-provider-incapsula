package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////
// GetApiSecurityApiConfig Tests
////////////////////////////////////////////////////////////////
func TestGetApiSecurityApiConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	apiID := 100

	apiConfigGetResponse, err := client.GetApiSecurityApiConfig(siteID, apiID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when reading Api-Security Api Config for Api ID %d", apiID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}

}

func TestGetApiSecurityApiConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := 100
	endpoint := fmt.Sprintf("%s%d/%d", apiConfigUrl, siteID, apiConfigID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.GetApiSecurityApiConfig(siteID, apiConfigID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing GET Api-Security Api Config JSON response for API ID %d", apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestGetApiSecurityApiConfigInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := 100
	endpoint := fmt.Sprintf("%s%d/%d", apiConfigUrl, siteID, apiConfigID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "value": "An internal error occurred. Contact support specifying your account ID and site ID.",
    "isError": true
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.GetApiSecurityApiConfig(siteID, apiConfigID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing GET Api-Security Api Config JSON response for API ID %d", apiConfigID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}

}

func TestGetApiSecurityApiConfigValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := 100
	endpoint := fmt.Sprintf("%s%d/%d", apiConfigUrl, siteID, apiConfigID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
   "value": {
       "id": 123,
       "siteId": 222,
       "siteName": "www.abc.incaptest.co",
       "hostName": "api.imperva.com",
       "basePath": "/api-security",
       "description": "third api-security collection",
       "lastModified": 1630392889000,
       "creationTime": 1630392888000,
       "apiSource": "USER",
       "violationActions": {
           "invalidUrlViolationAction": "IGNORE",
           "invalidMethodViolationAction": "IGNORE",
           "missingParamViolationAction": "IGNORE",
           "invalidParamValueViolationAction": "IGNORE",
           "invalidParamNameViolationAction": "IGNORE"
       },
       "specificationViolationAction": "IGNORE"
   },
   "isError": false
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.GetApiSecurityApiConfig(siteID, apiConfigID)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), apiConfigGetResponse)
	}
	if apiConfigGetResponse == nil {
		t.Errorf("Should not have received a nil apiSecuritySiteConfigGetResponse instance")
	}
	if apiConfigGetResponse.Value.SiteId != 222 {
		t.Errorf("Site ID doesn't match. Actual : %v", apiConfigGetResponse.Value.Id)
	}
}

////////////////////////////////////////////////////////////////
// CreateApiSecurityApiConfig Tests
////////////////////////////////////////////////////////////////
func TestCreateApiSecurityApiConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}
	apiConfigGetResponse, err := client.CreateApiSecurityApiConfig(siteID, &payload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error adding API Security API Config for site %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestCreateApiSecurityApiConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s%d", apiConfigUrl, siteID)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.CreateApiSecurityApiConfig(siteID, &payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing add API Security API Config JSON response for site id %d", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestCreateApiSecurityApiConfigInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	//apiConfigID := 100
	endpoint := fmt.Sprintf("%s%d", apiConfigUrl, siteID)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "value": "apiSpecification field is missing",
    "isError": true
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.CreateApiSecurityApiConfig(siteID, &payload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 400 from Incapsula service while creating API Security API Config for Site ID %d", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestCreateApiSecurityApiConfigValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("%s%d", apiConfigUrl, siteID)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`	{
		"value": {
		"apiId": 8016,
			"resultMessage": "API 8016 was added successfully.",
			"duplicateEndpointsList": []
	},
	"isError": false
	}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.CreateApiSecurityApiConfig(siteID, &payload)
	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if apiConfigGetResponse == nil {
		t.Errorf("Should not have received a nil apiConfigGetResponse instance")
	}
	if apiConfigGetResponse.Value.ApiId == 0 {
		t.Errorf("Should not have received an empty site config ID")
	}
}

////////////////////////////////////////////////////////////////
// UpdateApiSecurityApiConfig Tests
////////////////////////////////////////////////////////////////

func TestUpdateApiSecurityApiConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	apiConfigID := "100"

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}
	apiConfigGetResponse, err := client.UpdateApiSecurityApiConfig(siteID, apiConfigID, &payload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error updating API Security API Config for site id %d, API id %s", siteID, apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestUpdateApiSecurityApiConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := "100"
	endpoint := fmt.Sprintf("%s%d/%s", apiConfigUrl, siteID, apiConfigID)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.UpdateApiSecurityApiConfig(siteID, apiConfigID, &payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing update API Security API Config JSON response for Site ID %d, API id %s", siteID, apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestUpdateApiSecurityApiConfigInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := "100"
	endpoint := fmt.Sprintf("%s%d/%s", apiConfigUrl, siteID, apiConfigID)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "value": "An API does not exist for account id: 50759581, site id: 77002718, api id: 7971. Please use the appropriate command to add it",
    "isError": true
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.UpdateApiSecurityApiConfig(siteID, apiConfigID, &payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service while updating API Security API for siteId %d, API id %s", siteID, apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestUpdateApiSecurityApiConfigValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := "100"
	endpoint := fmt.Sprintf("%s%d/%s", apiConfigUrl, siteID, apiConfigID)

	payload := ApiSecurityApiConfigPostPayload{
		ValidateHost:     true,
		Description:      "test api config",
		ApiSpecification: "/Users/katrin.polit/api_security_files/api-security-swagger.yaml",
		BasePath:         "/api-security",
		ViolationActions: ViolationActions{
			InvalidUrlViolationAction:        "IGNORE",
			InvalidMethodViolationAction:     "IGNORE",
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "IGNORE",
			InvalidParamValueViolationAction: "IGNORE",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "value": {
        "apiId": 100,
        "resultMessage": "API 100 was updated successfully.",
        "duplicateEndpointsList": []
    },
    "isError": false
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiConfigGetResponse, err := client.UpdateApiSecurityApiConfig(siteID, apiConfigID, &payload)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if apiConfigGetResponse == nil {
		t.Errorf("Should not have received a nil apiConfigGetResponse instance")
	}
	if apiConfigGetResponse.Value.ApiId == 0 {
		t.Errorf("Should not have received an empty site config ID")
	}
}

////////////////////////////////////////////////////////////////
// DeleteApiSecurityApiConfig Tests
////////////////////////////////////////////////////////////////

func TestDeleteApiSecurityApiConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	apiConfigID := "100"

	err := client.DeleteApiSecurityApiConfig(siteID, apiConfigID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when deleting API Secirity API Config with Site ID %d, API ID %s", siteID, apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestDeleteApiSecurityApiConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := "100"
	endpoint := fmt.Sprintf("%s%d/%s", apiConfigUrl, siteID, apiConfigID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteApiSecurityApiConfig(siteID, apiConfigID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing delete API Secirity API Config JSON response for Site ID %d, API Config ID %s", siteID, apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestDeleteApiSecurityApiConfigInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := "100"
	endpoint := fmt.Sprintf("%s%d/%s", apiConfigUrl, siteID, apiConfigID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "value": "API id: 0 (account id: 0, site id: 0) was not found. Delete operation was aborted.",
    "isError": true
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteApiSecurityApiConfig(siteID, apiConfigID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code %d from Incapsula service when deleting API Security API Config for Site ID %d, API Config ID %s", 400, siteID, apiConfigID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestDeleteApiSecurityApiConfigValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	apiConfigID := "100"
	endpoint := fmt.Sprintf("%s%d/%s", apiConfigUrl, siteID, apiConfigID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
		"value": "Deletion successful",
		"isError": false
	}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteApiSecurityApiConfig(siteID, apiConfigID)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
