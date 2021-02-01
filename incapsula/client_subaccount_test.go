package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//////////////////////////////////////////////////////////////
/// 	AddSubAccount Tests
//////////////////////////////////////////////////////////////

func TestClientAddSubAccountBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	subAccountAddResponse, err := client.AddSubAccount("", "", "", 0, 0)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error adding subaccount")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if subAccountAddResponse != nil {
		t.Errorf("Should have received a nil addSubAccountResponse instance")
	}
}

func TestClientAddSubAccountBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSubAccountAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSubAccountAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	subaccountName := "testsubaccount"
	subAccountAddResponse, err := client.AddSubAccount(subaccountName, "", "", 0, 0)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add subaccount JSON response")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if subAccountAddResponse != nil {
		t.Errorf("Should have received a nil addSubAccountResponse instance")
	}
}

func TestClientAddSubAccountInvalidParent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSubAccountAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSubAccountAdd, req.URL.String())
		}
		rw.Write([]byte(`{"parent_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	subAccountAddResponse, err := client.AddSubAccount("", "", "", 0, 0)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding subaccount")) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
	if subAccountAddResponse != nil {
		t.Errorf("Should have received a nil addSubAccountResponse instance")
	}
}

////////////////////////////////////////////////////////////////
/// 	DeleteSubAccount Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteSubAccountBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	subAccountID := 123
	err := client.DeleteSubAccount(subAccountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting subaccount id: %d", subAccountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteSubAccountBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSubAccountDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSubAccountDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	subAccountID := 123
	err := client.DeleteSubAccount(subAccountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete account JSON response for subaccount id: %d", subAccountID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteSubAccountInvalidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSubAccountDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSubAccountDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	subAccountID := 123
	err := client.DeleteSubAccount(subAccountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting subaccount id: %d", subAccountID)) {
		t.Errorf("Should have received a bad account error, got: %s", err)
	}
}

func TestClientDeleteSubAccountValidAccount(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointSubAccountDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointSubAccountDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	subAccountID := 123
	err := client.DeleteSubAccount(subAccountID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}