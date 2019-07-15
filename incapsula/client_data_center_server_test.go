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
// AddDataCenterServer Tests
////////////////////////////////////////////////////////////////

func TestClientAddDataCenterServerBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	dcID := "42"
	addDataCenterServerResponse, err := client.AddDataCenterServer(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center server for dcID %s", dcID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addDataCenterServerResponse != nil {
		t.Errorf("Should have received a nil addDataCenterServerResponse instance")
	}
}

func TestClientAddDataCenterServerBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	addDataCenterServerResponse, err := client.AddDataCenterServer(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add data center server JSON response for dcID %s", dcID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addDataCenterServerResponse != nil {
		t.Errorf("Should have received a nil addDataCenterServerResponse instance")
	}
}

func TestClientAddDataCenterServerInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"0","res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	addDataCenterServerResponse, err := client.AddDataCenterServer(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center server for dcID %s", dcID)) {
		t.Errorf("Should have received a data center server error, got: %s", err)
	}
	if addDataCenterServerResponse != nil {
		t.Errorf("Should have received a nil addDataCenterServerResponse instance")
	}
}

func TestClientAddDataCenterServerValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"123","res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	addDataCenterServerResponse, err := client.AddDataCenterServer(dcID, "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDataCenterServerResponse == nil {
		t.Errorf("Should not have received a nil addDataCenterServerResponse instance")
	}
	if addDataCenterServerResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// EditDataCenterServer Tests
////////////////////////////////////////////////////////////////

func TestClientEditDataCenterServerBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServer(serverID, "", "", "")
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

func TestClientEditDataCenterServerBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServer(serverID, "", "", "")
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

func TestClientEditDataCenterServerInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":0,"res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServer(serverID, "", "", "")
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

func TestClientEditDataCenterServerValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":123,"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "411"
	editDataCenterResponse, err := client.EditDataCenterServer(serverID, "", "", "")
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
// DeleteDataCenterServer Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteDataCenterServerBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	serverID := "42"
	err := client.DeleteDataCenterServer(serverID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting data center server (server_id: %s)", serverID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServerBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "42"
	err := client.DeleteDataCenterServer(serverID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete data center server JSON response (server_id: %s)", serverID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServerInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "42"
	err := client.DeleteDataCenterServer(serverID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting data center server (server_id: %s)", serverID)) {
		t.Errorf("Should have received a bad data center error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServerValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServerDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServerDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	serverID := "42"
	err := client.DeleteDataCenterServer(serverID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
