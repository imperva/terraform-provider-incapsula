package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// //////////////////////////////////////////////////////////////
// GetApiSecurityEndpointConfig Tests
// //////////////////////////////////////////////////////////////
func TestGetApiSecurityEndpointConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	//siteID := 42
	apiID := 100
	endpointId := "92"
	//
	apiSecurityEndpointConfigGetResponse, err := client.GetApiSecurityEndpointConfig(apiID, endpointId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service while reading Api-Security Endpoint Config for API ID %d and Endpoint ID %s:", apiID, endpointId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecurityEndpointConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}

}

func TestGetApiSecurityEndpointConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	apiConfigID := 100
	endpointId := "92"
	endpoint := fmt.Sprintf("%s%d/%s", endpointConfigUrl, apiConfigID, endpointId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiSecurityEndpointConfigGetResponse, err := client.GetApiSecurityEndpointConfig(apiConfigID, endpointId)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing GET Api-Security Endpoint Config JSON response for API ID %d and endpoint ID %s", apiConfigID, endpointId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecurityEndpointConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestGetApiSecurityEndpointConfigInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	apiConfigID := 100
	endpointId := "92"
	endpoint := fmt.Sprintf("%s%d/%s", endpointConfigUrl, apiConfigID, endpointId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
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

	apiSecurityEndpointConfigGetResponse, err := client.GetApiSecurityEndpointConfig(apiConfigID, endpointId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service when reading Api-Security Endpoint Config for API ID %d and Endpoint ID %s", apiConfigID, endpointId)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiSecurityEndpointConfigGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestGetApiSecurityEndpointConfigValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	apiConfigID := 100
	endpointId := "92"
	endpoint := fmt.Sprintf("%s%d/%s", endpointConfigUrl, apiConfigID, endpointId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`	{
    "value": {
        "id": 2,
        "path": "/api/test",
        "method": "GET",
        "violationActions": {
            "missingParamViolationAction": "DEFAULT",
            "invalidParamValueViolationAction": "DEFAULT",
            "invalidParamNameViolationAction": "DEFAULT"
        },
        "specificationViolationAction": "DEFAULT"
    },
    "isError": false
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	apiSecurityEndpointConfigGetResponse, err := client.GetApiSecurityEndpointConfig(apiConfigID, endpointId)
	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if apiSecurityEndpointConfigGetResponse == nil {
		t.Errorf("Should not have received a nil apiConfigGetResponse instance")
	}
	if apiSecurityEndpointConfigGetResponse.Value.Id != 2 {
		t.Errorf("Should not have received an empty site config ID")
	}
}

// //////////////////////////////////////////////////////////////
// PostApiSecurityEndpointConfig Tests
// //////////////////////////////////////////////////////////////
func TestPostApiSecurityEndpointConfigBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	apiConfigID := 100
	endpointId := 92
	payload := ApiSecurityEndpointConfigPostPayload{
		ViolationActions: UserViolationActions{
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "BLOCK_USER",
			InvalidParamValueViolationAction: "BLOCK_IP",
		},
		SpecificationViolationAction: "BLOCK_REQUEST",
	}

	apiSecurityEndpointConfigPostResponse, err := client.PostApiSecurityEndpointConfig(apiConfigID, endpointId, &payload)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service while updating Api Security Endpoint Configuration for API Config Id %d, API Config Id %d", apiConfigID, endpointId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecurityEndpointConfigPostResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}

}

func TestPostApiSecurityEndpointConfigBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	apiConfigID := 100
	endpointId := 92
	endpoint := fmt.Sprintf("%s%d"+"/"+"%d", endpointConfigUrl, apiConfigID, endpointId)

	payload := ApiSecurityEndpointConfigPostPayload{
		ViolationActions: UserViolationActions{
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "BLOCK_USER",
			InvalidParamValueViolationAction: "BLOCK_IP",
		},
		SpecificationViolationAction: "BLOCK_REQUEST",
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()
	//
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	//
	apiSecurityEndpointConfigPostResponse, err := client.PostApiSecurityEndpointConfig(apiConfigID, endpointId, &payload)
	//
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing api-security JSON response for create/update Api Security Endpoint Configuration for API Config Id %d, Endpoint Config Id %d", apiConfigID, endpointId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiSecurityEndpointConfigPostResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestPostApiSecurityEndpointConfigInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	apiConfigID := 100
	endpointId := 92
	endpoint := fmt.Sprintf("%s%d"+"/"+"%d", endpointConfigUrl, apiConfigID, endpointId)

	payload := ApiSecurityEndpointConfigPostPayload{
		ViolationActions: UserViolationActions{
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "BLOCK_USER",
			InvalidParamValueViolationAction: "BLOCK_IP",
		},
		SpecificationViolationAction: "BLOCK_REQUEST",
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(403)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`<html>

<head>
	<title>Error</title>
</head>

<body>Bad Request</body>

</html>`))
	}))
	defer server.Close()
	//
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	//
	apiSecurityEndpointConfigPostResponse, err := client.PostApiSecurityEndpointConfig(apiConfigID, endpointId, &payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 403 from Incapsula service while updating Api Security Endpoint configuration for API Config Id %d, Endpoint Config Id: %d. Error:", apiConfigID, endpointId)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiSecurityEndpointConfigPostResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestPostApiSecurityEndpointConfigValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	apiConfigID := 100
	endpointId := 92
	endpoint := fmt.Sprintf("%s%d"+"/"+"%d", endpointConfigUrl, apiConfigID, endpointId)

	payload := ApiSecurityEndpointConfigPostPayload{
		ViolationActions: UserViolationActions{
			MissingParamViolationAction:      "IGNORE",
			InvalidParamNameViolationAction:  "BLOCK_USER",
			InvalidParamValueViolationAction: "BLOCK_IP",
		},
		SpecificationViolationAction: "BLOCK_REQUEST",
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "value": {
        "endpointId": 92
    },
    "isError": false
}`))
	}))
	defer server.Close()
	//
	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	//
	apiSecurityEndpointConfigPostResponse, err := client.PostApiSecurityEndpointConfig(apiConfigID, endpointId, &payload)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if apiSecurityEndpointConfigPostResponse == nil {
		t.Errorf("Should not have received a nil apiSecurityEndpointConfigPostResponse instance")
	}
	if apiSecurityEndpointConfigPostResponse.Value.EndpointId == 92 {
		t.Errorf("Should not have received an empty site config ID")
	}
}
