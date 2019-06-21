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
// AddDataCenterServers Tests
////////////////////////////////////////////////////////////////

func TestClientAddDataCenterServersBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	dcID := "42"
	addDataCenterServersResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center server for dcID %s", dcID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addDataCenterServersResponse != nil {
		t.Errorf("Should have received a nil addDataCenterServersResponse instance")
	}
}

func TestClientAddDataCenterServersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	addDataCenterServersResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add data center server JSON response for dcID %s", dcID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addDataCenterServersResponse != nil {
		t.Errorf("Should have received a nil addDataCenterServersResponse instance")
	}
}

func TestClientAddDataCenterServersInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"0","res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	addDataCenterServersResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center server for dcID %s", dcID)) {
		t.Errorf("Should have received a data center servers error, got: %s", err)
	}
	if addDataCenterServersResponse != nil {
		t.Errorf("Should have received a nil addDataCenterServersResponse instance")
	}
}

func TestClientAddDataCenterServersValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"123","res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	addDataCenterServersResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDataCenterServersResponse == nil {
		t.Errorf("Should not have received a nil addDataCenterServersResponse instance")
	}
	if addDataCenterServersResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// EditDataCenterServers Tests
////////////////////////////////////////////////////////////////

func TestClientEditDataCenterServersBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServers(serverID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error editing data center server for serverID: %s: ", serverID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if editDataCenterResponse != nil {
		t.Errorf("Should have received a nil editDataCenterResponse instance")
	}
}

func TestClientEditDataCenterServersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServers(serverID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing edit data center server JSON response for serverID %s", serverID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editDataCenterResponse != nil {
		t.Errorf("Should have received a nil editDataCenterResponse instance")
	}
}

func TestClientEditDataCenterServersInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":0,"res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServers(serverID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when editing data center server for serverID %s", serverID)) {
		t.Errorf("Should have received a bad data center error, got: %s", err)
	}
	if editDataCenterResponse != nil {
		t.Errorf("Should have received a nil editDataCenterResponse instance")
	}
}

func TestClientEditDataCenterServersValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":123,"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServers(serverID, "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if editDataCenterResponse == nil {
		t.Errorf("Should not have received a nil editDataCenterResponse instance")
	}
	// todo: test response properties
	if editDataCenterResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteDataCenterServers Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteDataCenterServersBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	serverID := 42
	err := client.DeleteDataCenterServers(serverID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting data center server (server_id: %d)", serverID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := 42
	err := client.DeleteDataCenterServers(serverID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete data center server JSON response (server_id: %d)", serverID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServersInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := 42
	err := client.DeleteDataCenterServers(serverID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting data center server (server_id: %d)", serverID)) {
		t.Errorf("Should have received a bad data center error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServersValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := 42
	err := client.DeleteDataCenterServers(serverID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
