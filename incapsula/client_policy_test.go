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
// GetAllPoliciesForAccount Tests
////////////////////////////////////////////////////////////////
func TestGetAllPoliciesForAccountBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountId := "92"

	getAllPolicieGetResponse, err := client.GetAllPoliciesForAccount(accountId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when reading All Policies for Account ID %s", accountId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if getAllPolicieGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestGetAllPoliciesForAccountBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	accountId := "92"

	endpoint := fmt.Sprintf("/policies/v2/policies?extended=true")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	getAllPolicieGetResponse, err := client.GetAllPoliciesForAccount(accountId)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing All Policies JSON response for Account ID %s", accountId)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if getAllPolicieGetResponse != nil {
		t.Errorf("Should have received a nil apiSecuritySiteConfigGetResponse instance")
	}
}

func TestGetAllPoliciesForAccountInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	accountId := "92"

	endpoint := fmt.Sprintf("/policies/v2/policies?extended=true")

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

	getAllPolicieGetResponse, err := client.GetAllPoliciesForAccount(accountId)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service when reading All Policies for Account ID %s", accountId)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if getAllPolicieGetResponse != nil {
		t.Errorf("Should have received a nil apiConfigGetResponse instance")
	}
}

func TestGetAllPoliciesForAccountValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	accountId := "92"

	endpoint := fmt.Sprintf("/policies/v2/policies?extended=true")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`	{
    "value": [
       {
            "defaultPolicyConfig": [],
            "policySettings": [
                {
                    "id": 7890,
                    "policyId": 7854390,
                    "settingsAction": "ALLOW",
                    "policySettingType": "IP",
                    "data": {
                        "ips": [
                            "1.2.3.4"
                        ]
                    },
                    "policyDataExceptions": []
                }
            ],
            "id": 7890,
            "name": "WHITELIST test for subaccount",
            "description": "WHITELIST test for subaccount",
            "enabled": true,
            "accountId": 1234,
            "policyType": "WHITELIST",
            "lastModifiedBy": 987439875
        }
    ],
    "isError": false
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	getAllPolicieGetResponse, err := client.GetAllPoliciesForAccount(accountId)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if getAllPolicieGetResponse == nil {
		t.Errorf("Should not have received a nil getAllPolicieGetResponse instance")
	}
	if (*getAllPolicieGetResponse)[0].AccountID != 1234 {
		t.Errorf("Should not have received an empty site config ID")
	}
}
