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

func TestClientUpdateBotAccessControlConfigurationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	requestDTO := BotsConfigurationDTO{}
	responseDTO, err := client.UpdateBotAccessControlConfiguration(siteID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error executing update Bot Access Control configuration request for siteID %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientUpdateBotAccessControlConfigurationBadJSON(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/settings/botConfiguration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/settings/botConfiguration endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := BotsConfigurationDTO{}
	responseDTO, err := client.UpdateBotAccessControlConfiguration(siteID, requestDTO)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Update Bot Access Control configuration JSON response for siteID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientUpdateBotAccessControlInvalidBotConfiguration(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/settings/botConfiguration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/settings/botConfiguration endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"errors":[{"status": "406"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := BotsConfigurationDTO{}
	responseDTO, err := client.UpdateBotAccessControlConfiguration(siteID, requestDTO)
	if err != nil {
		t.Errorf("Should not receive an error. Got: %s", err.Error())
	}
	if responseDTO == nil || responseDTO.Errors == nil || len(responseDTO.Errors) < 1 {
		t.Errorf("Should have received a response DTO instance with at least one error item")
		return
	}
	if responseDTO.Errors[0].Status != "406" {
		t.Errorf("Should have received a bad BOT configuration error, got: %s", err)
	}
}

func TestClientUpdateBotAccessControlValidBotConfiguration(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/settings/botConfiguration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/settings/botConfiguration endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"data":[{"canceledGoodBots":[{"displayName":"Googlebot (Site Helper)"}], "badBots":[{"displayName":"Googlebot (Site Helper)"}]}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestDTO := BotsConfigurationDTO{}
	responseDTO, err := client.UpdateBotAccessControlConfiguration(siteID, requestDTO)
	if err != nil {
		t.Errorf("Should not have received an error. Got: %s", err.Error())
	}
	if responseDTO == nil {
		t.Errorf("Should not have received a nil response DTO instance")
		return
	}
	if responseDTO.Data == nil || len(responseDTO.Data) < 1 ||
		responseDTO.Data[0].BadBots == nil || len(responseDTO.Data[0].BadBots) < 1 ||
		responseDTO.Data[0].CanceledGoodBots == nil || len(responseDTO.Data[0].CanceledGoodBots) < 1 {
		t.Errorf("Response must contain one BadBots and one CanceledGoodBots. Items: %d", len(responseDTO.Data))
		t.Errorf("Response must contain one BadBots and one CanceledGoodBots. BadBots: %d", len(responseDTO.Data[0].BadBots))
		t.Errorf("Response must contain one BadBots and one CanceledGoodBots. CanceledGoodBots: %d", len(responseDTO.Data[0].CanceledGoodBots))
	}
}

////////////////////////////////////////////////////////////////
// ListDataCenters Tests
////////////////////////////////////////////////////////////////

func TestClientGetBotAccessControlConfigurationBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	responseDTO, err := client.GetBotAccessControlConfiguration(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf(
		"Error executing get Bot Access Control configuration request for siteID %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientGetBotAccessControlConfigurationBadJSON(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/settings/botConfiguration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/settings/botConfiguration endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	responseDTO, err := client.GetBotAccessControlConfiguration(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Bot Access Control list JSON response for siteID: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if responseDTO != nil {
		t.Errorf("Should have received a nil responseDTO instance")
	}
}

func TestClientGetBotAccessControlConfigurationInvalidRequest(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/settings/botConfiguration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/settings/botConfiguration endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"errors":[{"status": "404"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	responseDTO, err := client.GetBotAccessControlConfiguration(siteID)
	if err != nil {
		t.Errorf("Should not receive an error. Got: %s", err.Error())
	}
	if responseDTO == nil || responseDTO.Errors == nil || len(responseDTO.Errors) < 1 {
		t.Errorf("Should have received a response DTO instance with at least one error item")
		return
	}
	if responseDTO.Errors[0].Status != "404" {
		t.Errorf("Should have received a bad BOT configuration error, got: %s", err)
	}
}

func TestClientGetBotAccessControlConfigurationValidRequest(t *testing.T) {
	siteID := "42"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/api/prov/v3/sites/%s/settings/botConfiguration", siteID) {
			t.Errorf("Should have have hit /api/prov/v3/sites/%s/settings/botConfiguration endpoint. "+
				"Got: %s", siteID, req.URL.String())
		}
		rw.Write([]byte(`{"data":[{"canceledGoodBots":[{"displayName":"Googlebot (Site Helper)"}], "badBots":[{"displayName":"Googlebot (Site Helper)"}]}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL + "/api/prov/v1"}
	client := &Client{config: config, httpClient: &http.Client{}}
	responseDTO, err := client.GetBotAccessControlConfiguration(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if responseDTO == nil {
		t.Errorf("Should not have received a nil responseDTO instance")
	}
	if responseDTO.Data == nil || len(responseDTO.Data) < 1 ||
		responseDTO.Data[0].BadBots == nil || len(responseDTO.Data[0].BadBots) < 1 ||
		responseDTO.Data[0].CanceledGoodBots == nil || len(responseDTO.Data[0].CanceledGoodBots) < 1 {
		t.Errorf("Response must contain one BadBots and one CanceledGoodBots. Items: %d", len(responseDTO.Data))
		t.Errorf("Response must contain one BadBots and one CanceledGoodBots. BadBots: %d", len(responseDTO.Data[0].BadBots))
		t.Errorf("Response must contain one BadBots and one CanceledGoodBots. CanceledGoodBots: %d", len(responseDTO.Data[0].CanceledGoodBots))
	}
}
