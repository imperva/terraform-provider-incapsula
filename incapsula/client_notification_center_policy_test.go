package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var notificationPolicyFullDto = NotificationPolicyFullDto{
	AccountId:   1234,
	PolicyName:  "The best policy in the world",
	SubCategory: "ACCOUNT_NOTIFICATIONS",
}

func TestClientCreateNotificationCenterPolicyBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	notificationCenterPolicyAddResponse, err := client.AddNotificationCenterPolicy(&notificationPolicyFullDto)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from NotificationCenter service when adding policy")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if notificationCenterPolicyAddResponse != nil {
		t.Errorf("Should have received a nil addSubAccountResponse instance")
	}
}

func TestClientNotificationCenterPolicyBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.String(), fmt.Sprintf("/%s", endPointNotificationCenterPolicy)) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endPointNotificationCenterPolicy, req.URL.String())
		}
		rw.Write([]byte(`{ bad json`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	notificationCenterPolicyAddResponse, err := client.AddNotificationCenterPolicy(&notificationPolicyFullDto)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing NotificationCenterPolicy JSON ")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if notificationCenterPolicyAddResponse != nil {
		t.Errorf("Should have received a nil instance")
	}
}

func TestClientAddNotificationCenterPolicyInvalidResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if !strings.HasPrefix(req.URL.String(), fmt.Sprintf("/%s", endPointNotificationCenterPolicy)) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endPointNotificationCenterPolicy, req.URL.String())
		}
		rw.Write([]byte(`{"error" : "cant add new policy""}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	notificationPolicyResponse, err := client.AddNotificationCenterPolicy(&notificationPolicyFullDto)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from NotificationCenter service when adding policy")) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if notificationPolicyResponse != nil {
		t.Errorf("Should have received a nil addNotificationCenterPolicyResponse instance")
	}
}

func TestAddNotificationCenterPolicyValidPolicy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.String(), fmt.Sprintf("/%s", endPointNotificationCenterPolicy)) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endPointNotificationCenterPolicy, req.URL.String())
		}
		rw.Write([]byte(`
			{
			    "data":
			    {
			        "accountId": 777,
					"policyId": 888,
			        "policyName": "Terraform acceptance test- policy account without assets",
			        "status": "ENABLE",
			        "subCategory": "ACCOUNT_NOTIFICATIONS",
			        "notificationChannelList":
			        [
			            {
			                "recipientToList":
			                [
			                    {
			                        "recipientType": "External",
			                        "displayName": "john.mcclane@externalemail.com"
			                    },
			                    {
			                        "recipientType": "External",
			                        "displayName": "another.exernal.email@gmail.com"
			                    }
			                ],
			                "channelType": "email"
			            }
			        ],
			        "assetList": [],
			        "applyToNewAssets": "FALSE",
			        "policyType": "ACCOUNT",
			        "subAccountPolicyInfo":
			        {
			            "applyToNewSubAccounts": "FALSE",
			            "subAccountList": []
			        }
			    }
			}
			`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	addNotificationPolicyResponse, err := client.AddNotificationCenterPolicy(&notificationPolicyFullDto)
	if err != nil {
		t.Errorf("Should not have received an error, the error: %s", err)
	}
	if addNotificationPolicyResponse == nil {
		t.Errorf("Should not have received a nil notificationPolicy instance")
	}
	if addNotificationPolicyResponse.Data.PolicyId != 888 {
		t.Errorf("Should not have received an empty policy Id")
	}
}
