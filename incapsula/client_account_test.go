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

////////////////////////////////////////////////////////////////
// ListSubAccounts Tests
////////////////////////////////////////////////////////////////

func TestClientListSubAccountsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accounts, err := client.ListSubAccounts()
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "Error listing sub accounts") {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if accounts != nil {
		t.Errorf("Should have received nil accounts")
	}
}

func TestClientListSubAccountsBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}
	_, err := client.ListSubAccounts()
	if err == nil || !strings.HasPrefix(err.Error(), "Error parsing sub accounts list JSON response") {
		t.Errorf("Should have received a JSON parse error, got: %v", err)
	}
}

func TestClientListSubAccountsErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"message":"not eligible"}`))
	}))
	defer server.Close()

	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}
	_, err := client.ListSubAccounts()
	if err == nil || !strings.HasPrefix(err.Error(), "Error from Incapsula service when listing sub accounts") {
		t.Errorf("Should have received a service error, got: %v", err)
	}
}

func TestClientListSubAccountsPaging(t *testing.T) {
	total := subAccountsListPageSize + 5
	var pagesRequested []string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/"+endpointSubAccountsList {
			t.Errorf("Should have hit /%s endpoint. Got: %s", endpointSubAccountsList, req.URL.Path)
		}
		if req.Method != http.MethodGet {
			t.Errorf("Should have used GET. Got: %s", req.Method)
		}
		if op := req.Header.Get("x-tf-operation"); !strings.HasPrefix(op, "read") {
			t.Errorf("Operation should start with read. Got: %s", op)
		}
		pageNum := req.URL.Query().Get("page_num")
		pagesRequested = append(pagesRequested, pageNum)
		count := subAccountsListPageSize
		first := 1
		if pageNum == "1" {
			count = 5
			first = subAccountsListPageSize + 1
		}
		var items []string
		for i := 0; i < count; i++ {
			items = append(items, fmt.Sprintf(`{"sub_account_id":%d,"sub_account_name":"acc%d","parent_id":7}`, first+i, first+i))
		}
		rw.Write([]byte(fmt.Sprintf(`{"data":[%s],"pagination":{"page_num":%s,"page_size":%d,"total_count":%d}}`,
			strings.Join(items, ","), pageNum, subAccountsListPageSize, total)))
	}))
	defer server.Close()

	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}
	accounts, err := client.ListSubAccounts()
	if err != nil {
		t.Fatalf("Should not have received an error: %v", err)
	}
	if len(accounts) != total {
		t.Errorf("Should have received %d accounts, got %d", total, len(accounts))
	}
	if strings.Join(pagesRequested, ",") != "0,1" {
		t.Errorf("Should have requested pages 0,1. Got: %v", pagesRequested)
	}
	if accounts[0].SubAccountID != 1 || accounts[0].ParentID != 7 || accounts[0].SubAccountName != "acc1" {
		t.Errorf("Unexpected first account: %+v", accounts[0])
	}
}
