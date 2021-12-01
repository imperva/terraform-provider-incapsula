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
	interface_var := make([]interface{},1)
	interface_var[0] = 0
	UserAddResponse, err := client.AddUser(0, email, interface_var, "", "")
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
server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s?api_id=%s&api_key=%s", endpointUserAdd, "foo" ,"bar" ) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointUserAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	email := "example@example.com"
	interface_var := make([]interface{},1)
	interface_var[0] = 10
	UserAddResponse, err := client.AddUser(accountID, email, interface_var, "f", "l")
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
	UserStatusResponse, err := client.UserStatus(accountID, email)
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
		if req.URL.String() != fmt.Sprintf("/%s?api_id=%s&api_key=%s&accountId=%d&userEmail=%s", endpointUserStatus, "foo" ,"bar",accountID, email) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointUserStatus, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	UserStatusResponse, err := client.UserStatus(accountID, email)
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
	err := client.DeleteUser(accountID, email)
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
		if req.URL.String() != fmt.Sprintf("/%s?api_id=%s&api_key=%s&accountId=%d&userEmail=%s", endpointUserDelete, "foo" ,"bar",accountID, email) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointUserDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteUser(accountID,email)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete user JSON response for user %s", email)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}
