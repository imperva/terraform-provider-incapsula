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
// AddRole Tests
////////////////////////////////////////////////////////////////

func TestClientAddRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	requestDTO := RoleDetailsCreateDTO{}
	RoleAddResponse, err := client.AddAccountRole(requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding account role")) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if RoleAddResponse != nil {
		t.Errorf("Should have received a nil RoleAddResponse instance")
	}
}

func TestClientAddRoleBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointRoleAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointRoleAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := RoleDetailsCreateDTO{}
	RoleAddResponse, err := client.AddAccountRole(requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add account role JSON response")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if RoleAddResponse != nil {
		t.Errorf("Should have received a nil RoleAddResponse instance")
	}
}

// ////////////////////////////////////////////////////////////////
// // GetRole Tests
// ////////////////////////////////////////////////////////////////

func TestClientGetRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	roleID := 123
	RoleStatusResponse, err := client.GetAccountRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error executing get Account Role")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if RoleStatusResponse != nil {
		t.Errorf("Should have received a nil RoleStatusResponse instance")
	}
}

func TestClientGetRoleBadJSON(t *testing.T) {
	roleID := 123
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s/%d", endpointRoleGet, roleID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointRoleGet, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	GetRoleResponse, err := client.GetAccountRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Account Role JSON response")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if GetRoleResponse != nil {
		t.Errorf("Should have received a nil GetRoleResponse instance")
	}
}

////////////////////////////////////////////////////////////////
// DeleteRole Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	roleID := 123
	err := client.DeleteAccountRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error executing delete Account Role")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteRoleBadJSON(t *testing.T) {
	roleID := 123
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s/%d", endpointRoleDelete, roleID) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointRoleDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	err := client.DeleteAccountRole(roleID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Account Role JSON response")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

////////////////////////////////////////////////////////////////
// UpdateUser Tests
////////////////////////////////////////////////////////////////

func TestClientUpdateRoleBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	roleID := 123
	requestDTO := RoleDetailsBasicDTO{}
	updateRoleResponse, err := client.UpdateAccountRole(roleID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating account role with Id %d", roleID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if updateRoleResponse != nil {
		t.Errorf("Should have received a nil updateRoleResponse instance")
	}
}

func TestClientUpdateRoleBadJSON(t *testing.T) {
	roleID := 123
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s/%d", endpointRoleUpdate, roleID) {
			t.Errorf("Should have have hit /%s/%d endpoint. Got: %s", endpointAccountUpdate, roleID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := RoleDetailsBasicDTO{}
	updateRoleResponse, err := client.UpdateAccountRole(roleID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update account role JSON response")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if updateRoleResponse != nil {
		t.Errorf("Should have received a nil updateRoleResponse instance")
	}
}
