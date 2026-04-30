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
// AddUser Tests
////////////////////////////////////////////////////////////////

func TestClientAddUserBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	email := "example@example.com"

	roleIds := make([]interface{}, 1)
	roleIds[0] = 0
	approvedIps := make([]interface{}, 0)
	UserAddResponse, err := client.AddAccountUser(0, email, "", "", roleIds, approvedIps)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding user email %s", email)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if UserAddResponse != nil {
		t.Errorf("Should have received a nil UserAddResponse instance")
	}
}

func TestClientAddUserBadJSON(t *testing.T) {
	accountID := 123
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s?caid=%d", endpointUserOperationNew, accountID) {
			t.Errorf("Should have have hit /%s?caid=%d endpoint. Got: %s", endpointUserOperationNew, accountId, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	email := "example@example.com"
	roleIds := make([]interface{}, 1)
	roleIds[0] = 10
	approvedIps := make([]interface{}, 0)
	UserAddResponse, err := client.AddAccountUser(accountID, email, "f", "l", roleIds, approvedIps)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add user JSON response for email %s", email)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if UserAddResponse != nil {
		t.Errorf("Should have received a nil UserAddResponse instance")
	}
}

// ////////////////////////////////////////////////////////////////
// // UserStatus Tests
// ////////////////////////////////////////////////////////////////

func TestClientUserStatusBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	email := "example@example.com"
	UserStatusResponse, err := client.GetAccountUser(accountID, email)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting user %s", email)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if UserStatusResponse != nil {
		t.Errorf("Should have received a nil UserStatusResponse instance")
	}
}

func TestClientUserStatusBadJSON(t *testing.T) {
	accountID := 123
	email := "example@example.com"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s/%s?caid=%d", endpointUserOperationNew, email, accountID) {
			t.Errorf("Should have have hit /%s/%s?caid=%d endpoint. Got: %s", endpointUserOperationNew, email, accountID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	UserStatusResponse, err := client.GetAccountUser(accountID, email)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing user status JSON response for user id %s", email)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if UserStatusResponse != nil {
		t.Errorf("Should have received a nil userStatusResponse instance")
	}
}

////////////////////////////////////////////////////////////////
// DeleteUser Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteUserBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	email := "example@example.com"
	err := client.DeleteAccountUser(accountID, email)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting USER: %s", email)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteUserBadJSON(t *testing.T) {
	accountID := 123
	email := "example@example.com"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s/%s?caid=%d", endpointUserOperationNew, email, accountID) {
			t.Errorf("Should have have hit /%s/%s?caid=%d endpoint. Got: %s", endpointUserOperationNew, email, accountID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteAccountUser(accountID, email)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete user JSON response for user %s", email)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

////////////////////////////////////////////////////////////////
// UpdateUser Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateUserBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	email := "example@example.com"
	roleIds := make([]interface{}, 1)
	roleIds[0] = 10
	approvedIps := make([]interface{}, 0)
	updateUserResponse, err := client.UpdateAccountUser(accountID, email, roleIds, approvedIps)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating user email %s", email)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updateUserResponse != nil {
		t.Errorf("Should have received a nil updateAccountResponse instance")
	}
}

func TestClientUpdateUserBadJSON(t *testing.T) {
	accountID := 123
	email := "example@example.com"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s/%s?caid=%d", endpointUserOperationNew, email, accountID) {
			t.Errorf("Should have have hit /%s/%s?caid=%d endpoint. Got: %s", endpointAccountUpdate, email, accountID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	roleIds := make([]interface{}, 1)
	roleIds[0] = 10
	approvedIps := make([]interface{}, 0)
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	updateUserResponse, err := client.UpdateAccountUser(accountID, email, roleIds, approvedIps)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update user JSON response for email %s", email)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if updateUserResponse != nil {
		t.Errorf("Should have received a nil updateUserResponse instance")
	}
}
