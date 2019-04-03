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
	dcID := 42
	addDataCenterResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center server for dcID %d", dcID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
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
	dcID := 42
	addDataCenterResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add data center server JSON response for dcID %d", dcID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
	}
}

func TestClientAddDataCenterServersInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":0,"res":1}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := 42
	addDataCenterResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center server for dcID %d", dcID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
	}
}

func TestClientAddDataCenterServersValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":123,"res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := 42
	addDataCenterResponse, err := client.AddDataCenterServers(dcID, "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDataCenterResponse == nil {
		t.Errorf("Should not have received a nil addDataCenterResponse instance")
	}
	if addDataCenterResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// ListDataCenterServers Tests
////////////////////////////////////////////////////////////////

func TestClientListDataCenterServersBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42
	listDataCentersResponse, err := client.ListDataCenterServers(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting data centers servers (site_id: %d)", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCenterServersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersList, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 42
	listDataCentersResponse, err := client.ListDataCenterServers(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing data centers servers list JSON response (site_id: %d)", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCenterServersInvalidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersList, req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 42
	listDataCentersResponse, err := client.ListDataCenterServers(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting data centers servers list (site_id: %d)", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCenterServersValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersList, req.URL.String())
		}
		// todo: what is response
		rw.Write([]byte(`{"foo":1527885500000, "bar":[], "res":0}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := 42
	listDataCentersResponse, err := client.ListDataCenterServers(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if listDataCentersResponse == nil {
		t.Errorf("Should not have received a nil listDataCentersResponse instance")
	}
	// todo: test response properties
	if listDataCentersResponse.Res != 0 {
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
		rw.Write([]byte(`{"rule_id":0,"res":1}`))
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
		t.Errorf("Should have received a bad site error, got: %s", err)
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
		rw.Write([]byte(`{"rule_id":123,"res":0}`))
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
	if editDataCenterResponse.Res != 0 {
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
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
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
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteDataCenterServersValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterServersDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterServersDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
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
