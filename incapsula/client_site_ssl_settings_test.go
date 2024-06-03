package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUpdateSiteSSLSettingsHandleBadConnection(t *testing.T) {
	// arrange
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev3: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	sslSettingsDTO := getUpdateSiteSSLSettingsDTO()

	// act
	var res, err = client.UpdateSiteSSLSettings(123, 1234, sslSettingsDTO)

	// assert
	if err == nil {
		t.Errorf("Should have received error from Incapsula service when updating Site SSL settings: %s", err)
	}

	if res != nil {
		t.Errorf("Should have received error updating site SSL settings and not a valid response: %s", err)
	}
}

func TestUpdateSiteSSLSettingsHandleResponseCodeNotSuccess(t *testing.T) {
	// arrange
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 1234

	endpoint := fmt.Sprintf("/sites-mgmt/v3/sites/%d/settings/TLSConfiguration", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		rw.Write([]byte(`{"res":406,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := getClientTestConfig(apiID, apiKey, server)

	client := &Client{config: config, httpClient: &http.Client{}}
	var dto = getUpdateSiteSSLSettingsDTO()

	// act
	_, err := client.UpdateSiteSSLSettings(siteID, accountID, dto)

	// assert
	if err == nil {
		t.Errorf("Should have received an error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error status code 406 from Incapsula service when updating Site SSL settings")) {
		t.Errorf("Should have received an update error, got: %s", err)
	}
}

func TestUpdateSiteSSLSettingsHandleInvalidResponseBody(t *testing.T) {
	// arrange
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 1234

	endpoint := fmt.Sprintf("/sites-mgmt/v3/sites/%d/settings/TLSConfiguration", siteID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		var invalidResponse = getValidJSONResponse() + "}"

		rw.Write([]byte(invalidResponse))
	}))
	defer server.Close()

	config := getClientTestConfig(apiID, apiKey, server)
	client := &Client{config: config, httpClient: &http.Client{}}
	var dto = getUpdateSiteSSLSettingsDTO()

	// act
	_, err := client.UpdateSiteSSLSettings(siteID, accountID, dto)

	// assert
	if err == nil {
		t.Errorf("Should have received an error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing Incap Site settings JSON response for Site ID %d", siteID)) {
		t.Errorf("Should have received a bad json response %s", err)
	}
}

func TestUpdateSiteSSLSettingsSuccess(t *testing.T) {
	// arrange
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 1234

	validResponse := getValidJSONResponse()

	endpoint := fmt.Sprintf("/sites-mgmt/v3/sites/%d/settings/TLSConfiguration", siteID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		rw.Write([]byte(validResponse))
	}))
	defer server.Close()

	config := getClientTestConfig(apiID, apiKey, server)
	client := &Client{config: config, httpClient: &http.Client{}}
	var dto = getUpdateSiteSSLSettingsDTO()

	// act
	_, err := client.UpdateSiteSSLSettings(siteID, accountID, dto)

	// assert
	if err != nil {
		t.Errorf("Should not have received an Error %s", err)
	}
}

func TestReadSiteSSLSettingsHandleRequestError(t *testing.T) {
	// arrange
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLRev3: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}

	// act
	var res, statusCode, err = client.ReadSiteSSLSettings(123, 1234)

	// assert
	if err == nil {
		t.Errorf("Should have received error from Incapsula service when reading Site SSL settings: %s", err)
	}

	if res != nil {
		t.Errorf("Should have recieved an error with no response")
	}

	if statusCode != 0 {
		t.Errorf("Should have received status code 0")
	}
}

func TestReadSiteSSLSettingsHandleResponseCodeNotSuccess(t *testing.T) {
	// arrange
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 1234

	endpoint := fmt.Sprintf("/sites-mgmt/v3/sites/%d/settings/TLSConfiguration", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"site_id":"42","id-info":"13008"}}`))
	}))
	defer server.Close()

	config := getClientTestConfig(apiID, apiKey, server)

	client := &Client{config: config, httpClient: &http.Client{}}

	// act
	_, statusCode, err := client.ReadSiteSSLSettings(siteID, accountID)

	// assert
	if err == nil {
		t.Errorf("Should have received an error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error status code %d from Incapsula service when reading SSL settings", statusCode)) {
		t.Errorf("Should have received a reading error, got: %s", err)
	}
}

func TestReadSiteSSLSettingsHandleInvalidResponseBody(t *testing.T) {
	// arrange
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 1234

	endpoint := fmt.Sprintf("/sites-mgmt/v3/sites/%d/settings/TLSConfiguration", siteID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		var invalidResponse = getValidJSONResponse() + "}"

		rw.Write([]byte(invalidResponse))
	}))
	defer server.Close()

	config := getClientTestConfig(apiID, apiKey, server)
	client := &Client{config: config, httpClient: &http.Client{}}

	// act
	_, _, err := client.ReadSiteSSLSettings(siteID, accountID)

	// assert
	if err == nil {
		t.Errorf("Should have received an error")
	}

	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error parsing Site SSL settings JSON response for Site ID %d", siteID)) {
		t.Errorf("Should have received a bad json response %s", err)
	}
}

func TestReadSiteSSLSettingsSuccess(t *testing.T) {
	// arrange
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	accountID := 1234

	var validResponse = getValidJSONResponse()

	endpoint := fmt.Sprintf("/sites-mgmt/v3/sites/%d/settings/TLSConfiguration", siteID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		rw.Write([]byte(validResponse))
	}))
	defer server.Close()

	config := getClientTestConfig(apiID, apiKey, server)
	client := &Client{config: config, httpClient: &http.Client{}}

	// act
	_, statusCode, err := client.ReadSiteSSLSettings(siteID, accountID)

	// assert
	if err != nil {
		t.Errorf("Should not received an error")
	}

	if statusCode != 200 {
		t.Errorf("Status code should be 200")
	}
}

func getUpdateSiteSSLSettingsDTO() SSLSettingsResponse {
	var sslSettingsDTO = SSLSettingsResponse{
		Data: []SSLSettingsDTO{
			{
				HstsConfiguration: &HSTSConfiguration{
					PreLoaded:          true,
					MaxAge:             1237,
					SubDomainsIncluded: true,
					IsEnabled:          true,
				},
				InboundTLSSettingsConfiguration: &InboundTLSSettingsConfiguration{
					ConfigurationProfile: "CUSTOM",
					TLSConfigurations: []TLSConfiguration{
						{
							TLSVersion: "TLS 1.1",
							CiphersSupport: []string{
								"TLS_AES_128_GCM_SHA256",
								"TLS_AES_128_GCM_SHA256",
							},
						},
						{
							TLSVersion: "TLS 1.2",
							CiphersSupport: []string{
								"TLS_AES_128_GCM_SHA256",
								"TLS_AES_128_GCM_SHA256",
							},
						},
					},
				},
				// add more setting types here
			},
		},
	}

	return sslSettingsDTO
}

func getClientTestConfig(apiID string, apiKey string, server *httptest.Server) *Config {
	config := &Config{
		APIID:       apiID,
		APIKey:      apiKey,
		BaseURL:     server.URL,
		BaseURLRev2: server.URL,
		BaseURLRev3: server.URL,
		BaseURLAPI:  server.URL,
	}
	return config
}

func getValidJSONResponse() string {
	var validResponse = `{
			"data":[
				{
					"hstsConfiguration":{
						"isEnabled":true,
						"maxAge":31536000,
						"subDomainsIncluded":true,
						"preLoaded":false
					},
                    "inboundTlsSettings": {
                        "configurationProfile": "CUSTOM",
                        "tlsConfiguration": [
                            {
                                "tlsVersion": "TLS 1.1",
                                "ciphersSupport": [
                                    "TLS_AES_128_GCM_SHA256",
                                    "TLS_AES_128_GCM_SHA256"
                                ]
                            }
                        ]
                    }
				}
			]
		}`
	return validResponse
}
