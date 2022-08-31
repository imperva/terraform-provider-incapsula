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
// AddAccount Tests
////////////////////////////////////////////////////////////////

func TestClientAddAccountBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	email := "example@example.com"
	addAccountResponse, err := client.AddAccount(email, "", "", "", "", "", 0, 0)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding account for email %s", email)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addAccountResponse != nil {
		t.Errorf("Should have received a nil addAccountResponse instance")
	}
}

func TestClientAddAccountBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	email := "example@example.com"
	addAccountResponse, err := client.AddAccount(email, "", "", "", "", "", 0, 0)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add account JSON response for email %s", email)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addAccountResponse != nil {
		t.Errorf("Should have received a nil addAccountResponse instance")
	}
}

func TestClientAddAccountInvalidParent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountAdd, req.URL.String())
		}
		rw.Write([]byte(`{"parent_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	email := "example@example.com"
	addAccountResponse, err := client.AddAccount(email, "", "", "", "", "", 0, 0)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding account for email %s", email)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if addAccountResponse != nil {
		t.Errorf("Should have received a nil addAccountResponse instance")
	}
}

func TestClientAddAccountValidParent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountAdd, req.URL.String())
		}
		rw.Write([]byte(`{"Account": {"parent_id":123},"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	email := "example@example.com"
	addAccountResponse, err := client.AddAccount(email, "", "", "", "", "", 0, 0)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addAccountResponse == nil {
		t.Errorf("Should not have received a nil addAccountResponse instance")
	}
	if addAccountResponse.Account.ParentID != 123 {
		t.Errorf("Parent ID doesn't match")
	}
	if addAccountResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// AccountStatus Tests
////////////////////////////////////////////////////////////////

func TestClientAccountStatusBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	accountStatusResponse, err := client.AccountStatus(accountID, ReadAccount)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting account status for account id %d", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if accountStatusResponse != nil {
		t.Errorf("Should have received a nil accountStatusResponse instance")
	}
}

func TestClientAccountStatusBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountStatus) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountStatus, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	accountStatusResponse, err := client.AccountStatus(accountID, ReadAccount)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing account status JSON response for account id %d", accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if accountStatusResponse != nil {
		t.Errorf("Should have received a nil accountStatusResponse instance")
	}
}

func TestClientAccountStatusInvalidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountStatus) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountStatus, req.URL.String())
		}
		rw.Write([]byte(`{"Account": {"trial_end_date":"0"}, "res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	accountStatusResponse, err := client.AccountStatus(accountID, ReadAccount)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting account status for account id %d", accountID)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if accountStatusResponse == nil {
		t.Errorf("Should have received a accountStatusResponse instance")
	}
}

func TestClientAccountStatusValidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountStatus) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountStatus, req.URL.String())
		}
		rw.Write([]byte(`{"Account": {"trial_end_date":"May 28, 2014"}, "res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	accountStatusResponse, err := client.AccountStatus(accountID, ReadAccount)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if accountStatusResponse == nil {
		t.Errorf("Should not have received a nil accountStatusResponse instance")
	}
	if accountStatusResponse.Account.TrialEndDate != "May 28, 2014" {
		t.Errorf("Account trial end date doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// UpdateAccount Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateAccountBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := "42"
	param := "error_page_template"
	value := "ABC123"
	updateAccountResponse, err := client.UpdateAccount(accountID, param, value)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating param (%s) with value (%s) on account_id: %s", param, value, accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updateAccountResponse != nil {
		t.Errorf("Should have received a nil updateAccountResponse instance")
	}
}

func TestClientUpdateAccountBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountUpdate, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "42"
	updateAccountResponse, err := client.UpdateAccount(accountID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update account JSON response for accountID %s", accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if updateAccountResponse != nil {
		t.Errorf("Should have received a nil updateAccountResponse instance")
	}
}

func TestClientUpdateAccountInvalidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"account_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "42"
	updateAccountResponse, err := client.UpdateAccount(accountID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when updating account for accountID %s", accountID)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if updateAccountResponse != nil {
		t.Errorf("Should have received a nil updateAccountResponse instance")
	}
}

func TestClientUpdateAccountValidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountUpdate) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountUpdate, req.URL.String())
		}
		rw.Write([]byte(`{"account_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := "42"
	updateAccountResponse, err := client.UpdateAccount(accountID, "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if updateAccountResponse == nil {
		t.Errorf("Should not have received a nil updateAccountResponse instance")
	}
	if updateAccountResponse.AccountID != 123 {
		t.Errorf("Account ID doesn't match")
	}
	if updateAccountResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteAccount Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteAccountBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	err := client.DeleteAccount(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting account id: %d", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteAccountBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	err := client.DeleteAccount(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete account JSON response for account id: %d", accountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteAccountInvalidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	err := client.DeleteAccount(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting account id: %d", accountID)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
}

func TestClientDeleteAccountValidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointAccountDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointAccountDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	err := client.DeleteAccount(accountID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
