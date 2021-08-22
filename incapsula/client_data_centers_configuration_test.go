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

func TestClientPutDataCentersConfigurationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	requestDTO := DataCentersConfigurationDTO{}
	responseDTO, err := client.PutDataCentersConfiguration(siteID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error executing update Data Centers configuration request for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientPutDataCentersConfigurationBadJSON(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/data-centers-configuration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/data-centers-configurations endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := DataCentersConfigurationDTO{}
	responseDTO, err := client.PutDataCentersConfiguration(siteID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update Data Centers configuration JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientPutDataCenterInvalidDcConfiguration(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/data-centers-configuration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/data-centers-configurations endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"errors":[{"status": "406"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := DataCentersConfigurationDTO{}
	responseDTO, err := client.PutDataCentersConfiguration(siteID, requestDTO)
	if err != nil {
		t.Errorf("Should not receive an error. Got: %s", err.Error())
	}
	if responseDTO == nil || responseDTO.Errors == nil || len(responseDTO.Errors) < 1 {
		t.Errorf("Should have received a response DTO instance with at least one error item")
		return
	}
	if responseDTO.Errors[0].Status != "406" {
		t.Errorf("Should have received a bad DC configuration error, got: %s", err)
	}
}

func TestClientPutDataCenterValidDcConfiguration(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/data-centers-configuration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/data-centers-configurations endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"data":[{"dataCenterMode":"SINGLE_DC","dataCenters":[{"name":"New DC","servers":[{"address":"1.2.3.4"}]}]}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := DataCentersConfigurationDTO{}
	responseDTO, err := client.PutDataCentersConfiguration(siteID, requestDTO)
	if err != nil {
		t.Errorf("Should not have received an error. Got: %s", err.Error())
	}
	if responseDTO == nil {
		t.Errorf("Should not have received a nil response DTO instance")
		return
	}
	if responseDTO.Data == nil || len(responseDTO.Data) < 1 || responseDTO.Data[0].DataCenters == nil ||
		len(responseDTO.Data[0].DataCenters) < 1 || responseDTO.Data[0].DataCenters[0].OriginServers == nil ||
		len(responseDTO.Data[0].DataCenters[0].OriginServers) < 1 {
		t.Errorf("Response must contain one Data Center, which contains one Origin Server. Items: %d", len(responseDTO.Data))
		t.Errorf("Response must contain one Data Center, which contains one Origin Server. Data Centers: %d",
			len(responseDTO.Data[0].DataCenters))
		t.Errorf("Response must contain one Data Center, which contains one Origin Server. Origin Servers: %d",
			len(responseDTO.Data[0].DataCenters[0].OriginServers))
	}
}

////////////////////////////////////////////////////////////////
// ListDataCenters Tests
////////////////////////////////////////////////////////////////

func TestClientGetDataCentersConfigurationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	responseDTO, err := client.GetDataCentersConfiguration(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf(
		"Error executing get Data Centers configuration request for siteID %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientGetDataCentersConfigurationBadJSON(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/data-centers-configuration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/data-centers-configurations endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	responseDTO, err := client.GetDataCentersConfiguration(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing data centers list JSON response for siteID: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientGetDataCentersConfigurationInvalidRequest(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/data-centers-configuration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/data-centers-configurations endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"errors":[{"status": "404"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	responseDTO, err := client.GetDataCentersConfiguration(siteID)
	if err != nil {
		t.Errorf("Should not receive an error. Got: %s", err.Error())
	}
	if responseDTO == nil || responseDTO.Errors == nil || len(responseDTO.Errors) < 1 {
		t.Errorf("Should have received a response DTO instance with at least one error item")
		return
	}
	if responseDTO.Errors[0].Status != "404" {
		t.Errorf("Should have received a bad DC configuration error, got: %s", err)
	}
}

func TestClientGetDataCentersConfigurationValidRequest(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/data-centers-configuration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/data-centers-configurations endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"data":[{"dataCenterMode":"SINGLE_DC","dataCenters":[{"name":"New DC","servers":[{"address":"1.2.3.4"}]}]}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	responseDTO, err := client.GetDataCentersConfiguration(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if responseDTO == nil {
		t.Errorf("Should not have received a nil responseDTO instance")
	}

	if responseDTO.Data == nil || len(responseDTO.Data) < 1 || responseDTO.Data[0].DataCenters == nil ||
		len(responseDTO.Data[0].DataCenters) < 1 || responseDTO.Data[0].DataCenters[0].OriginServers == nil ||
		len(responseDTO.Data[0].DataCenters[0].OriginServers) < 1 {
		t.Errorf("Response must contain one Data Center, which contains one Origin Server. Items: %d", len(responseDTO.Data))
		t.Errorf("Response must contain one Data Center, which contains one Origin Server. Data Centers: %d",
			len(responseDTO.Data[0].DataCenters))
		t.Errorf("Response must contain one Data Center, which contains one Origin Server. Origin Servers: %d",
			len(responseDTO.Data[0].DataCenters[0].OriginServers))
	}
}
