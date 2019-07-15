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
// AddDataCenter Tests
////////////////////////////////////////////////////////////////

func TestClientAddDataCenterBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	addDataCenterResponse, err := client.AddDataCenter(siteID, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
	}
}

func TestClientAddDataCenterBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	addDataCenterResponse, err := client.AddDataCenter(siteID, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add data center JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
	}
}

func TestClientAddDataCenterInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"0","res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	addDataCenterResponse, err := client.AddDataCenter(siteID, "", "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding data center for siteID %s", siteID)) {
		t.Errorf("Should have received a bad incap rule error, got: %s", err)
	}
	if addDataCenterResponse != nil {
		t.Errorf("Should have received a nil addDataCenterResponse instance")
	}
}

func TestClientAddDataCenterValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterAdd, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":"123","res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	addDataCenterResponse, err := client.AddDataCenter(siteID, "", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addDataCenterResponse == nil {
		t.Errorf("Should not have received a nil addDataCenterResponse instance")
	}
	if addDataCenterResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// ListDataCenters Tests
////////////////////////////////////////////////////////////////

func TestClientListDataCentersBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting data centers for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCentersBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterList, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing data centers list JSON response for siteID: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCentersInvalidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting data centers list (site_id: %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listDataCentersResponse != nil {
		t.Errorf("Should have received a nil listDataCentersResponse instance")
	}
}

func TestClientListDataCentersValidRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterList, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "42"
	listDataCentersResponse, err := client.ListDataCenters(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if listDataCentersResponse == nil {
		t.Errorf("Should not have received a nil listDataCentersResponse instance")
	}

	if listDataCentersResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// EditDataCenter Tests
////////////////////////////////////////////////////////////////

func TestClientEditDataCenterBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	dcID := "411"
	name := "foo"
	isStandBy := "yes"
	isContent := "yes"
	isActive := "yes"
	editDataCenterResponse, err := client.EditDataCenter(dcID, name, isStandBy, isContent, isActive)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error editing data center  for dcID: %s", dcID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if editDataCenterResponse != nil {
		t.Errorf("Should have received a nil editDataCenterResponse instance")
	}
}

func TestClientEditDataCenterBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "411"
	name := "foo"
	isStandBy := "yes"
	isContent := "yes"
	isActive := "yes"
	editDataCenterResponse, err := client.EditDataCenter(dcID, name, isStandBy, isContent, isActive)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing edit dta center JSON response for dcID %s", dcID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editDataCenterResponse != nil {
		t.Errorf("Should have received a nil editDataCenterResponse instance")
	}
}

func TestClientEditDataCenterInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":0,"res":"1"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "411"
	name := "foo"
	isStandBy := "yes"
	isContent := "yes"
	isActive := "yes"
	editDataCenterResponse, err := client.EditDataCenter(dcID, name, isStandBy, isContent, isActive)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when editing data center for dcID %s", dcID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if editDataCenterResponse != nil {
		t.Errorf("Should have received a nil editDataCenterResponse instance")
	}
}

func TestClientEditDataCenterValidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterEdit, req.URL.String())
		}
		rw.Write([]byte(`{"rule_id":123,"res":"0"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "411"
	name := "foo"
	isStandBy := "yes"
	isContent := "yes"
	isActive := "yes"
	editDataCenterResponse, err := client.EditDataCenter(dcID, name, isStandBy, isContent, isActive)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if editDataCenterResponse == nil {
		t.Errorf("Should not have received a nil editDataCenterResponse instance")
	}

	if editDataCenterResponse.Res != "0" {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteDataCenter Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteDataCenterBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	dcID := "42"
	err := client.DeleteDataCenter(dcID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting data center (dc_id: %s)", dcID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteDataCenterBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	err := client.DeleteDataCenter(dcID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing delete data center JSON response (dc_id: %s)", dcID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteDataCenterInvalidRule(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"1","res_message":"fail"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	err := client.DeleteDataCenter(dcID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting data center (dc_id: %s)", dcID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteDataCenterValidSite(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointDataCenterDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointDataCenterDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	dcID := "42"
	err := client.DeleteDataCenter(dcID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
