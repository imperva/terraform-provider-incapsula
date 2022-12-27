package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClientCreateSiemLogConfigurationBadConfig(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	badPrefix := RandomLowLetterString(10)
	config := &Config{
		APIID:       RandomCapitalLetterAndNumberString(20),
		APIKey:      RandomLetterAndNumberString(40),
		BaseURL:     badPrefix + ".incapsula.com",
		BaseURLRev2: badPrefix + ".incapsula.com",
		BaseURLAPI:  badPrefix + ".incapsula.com",
	}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}

	siemLogConfiguration, _, err := client.CreateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		AssetID:           RandomNumbersExcludingZeroString(10),
		ConfigurationName: RandomLetterAndNumberString(20),
		Provider:          RandomCapitalLetterString(6),
		Datasets:          append(make([]interface{}, 0), RandomCapitalLetterString(6), RandomCapitalLetterString(6)),
		Enabled:           true,
		ConnectionId:      RandomLowLetterAndNumberString(24),
	}}})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error from Incapsula service when executing %s operation on SIEM log configuration:", CreateSiemLogConfiguration)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siemLogConfiguration != nil {
		t.Errorf("Should have received a nil SiemLogConfiguration instance")
	}
}

func TestClientCreateSiemLogConfigurationBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	assetId := RandomNumbersExcludingZeroString(10)

	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemLogConfiguration, assetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(201)
		_, err := rw.Write([]byte(`{` + RandomLetterAndNumberString(20)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemLogConfiguration, _, err := client.CreateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		AssetID:           assetId,
		ConfigurationName: RandomLetterAndNumberString(20),
		Provider:          RandomCapitalLetterString(6),
		Datasets:          append(make([]interface{}, 0), RandomCapitalLetterString(6), RandomCapitalLetterString(6)),
		Enabled:           true,
		ConnectionId:      RandomLowLetterAndNumberString(24),
	}}})

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error obtained")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if siemLogConfiguration != nil {
		t.Errorf("Should have received a nil SiemLogConfiguration instance")
	}
}

func TestClientCreateSiemLogConfigurationBadRequest(t *testing.T) {
	responseStatusCode := 400
	m := make(map[string]string)
	m["Wrong connectionId"] = fmt.Sprintf(`{
					"errors": [
						{
							"status": %d,
							"id": "%s",
							"code": "1005",
							"source": {
								"pointer": "/%s"
							},
							"title": "Bad Request"
						}
					]
				}`, responseStatusCode, RandomLowLetterAndNumberString(24), endpointSiemLogConfiguration)
	m["Wrong configuration name"] = ""
	for k, v := range m {
		log.Printf("======================== BEGIN TEST ========================")
		log.Printf("[DEBUG] Executing failure scenario (%s) with response code: %d and message: %s", k, responseStatusCode, v)

		apiID := RandomCapitalLetterAndNumberString(20)
		apiKey := RandomLetterAndNumberString(40)
		assetId := RandomNumbersExcludingZeroString(10)
		endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemLogConfiguration, assetId)

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.String() != endpoint {
				t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
			}
			rw.WriteHeader(responseStatusCode)
			_, err := rw.Write([]byte(v))
			if err != nil {
				return
			}
		}))

		defer server.Close()

		config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
		client := &Client{config: config, httpClient: &http.Client{}}

		siemLogConfiguration, _, err := client.CreateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
			AssetID:           assetId,
			ConfigurationName: RandomLetterAndNumberString(20),
			Provider:          RandomCapitalLetterString(6),
			Datasets:          append(make([]interface{}, 0), RandomCapitalLetterString(6), RandomCapitalLetterString(6)),
			Enabled:           true,
			ConnectionId:      RandomLowLetterAndNumberString(24),
		}}})

		if err == nil {
			t.Errorf("Should have received an error")
		}
		if strings.Compare(err.Error(), fmt.Sprintf("received failure response for operation: %s on SIEM log configuration\nstatus code: %d\nbody: %s",
			CreateSiemLogConfiguration, responseStatusCode, v)) != 0 {
			t.Errorf("Should have received a response body for all responses different from %d, but received %s", 200, err)
		}
		if siemLogConfiguration != nil {
			t.Errorf("Should have received a nil SiemLogConfiguration instance")
		}
	}
}

func TestClientCreateValidSiemLogConfiguration(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	assetId := RandomNumbersExcludingZeroString(10)
	configurationName := RandomLetterAndNumberString(20)
	provider := RandomCapitalLetterString(6)
	dataset := RandomCapitalLetterString(6)
	connectionId := RandomLowLetterAndNumberString(24)

	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemLogConfiguration, assetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(201)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"assetId": "%s",
														"configurationName": "%s",
														"provider": "%s",
														"datasets": [
															"%s"
														],
														"enabled": true,
														"connectionId": "%s"
													}
												]
											}`,
			RandomLowLetterAndNumberString(24),
			assetId,
			configurationName,
			provider,
			dataset,
			connectionId,
		)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemLogConfiguration, _, err := client.CreateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{{
		AssetID:           assetId,
		ConfigurationName: configurationName,
		Provider:          provider,
		Datasets:          append(make([]interface{}, 0), dataset),
		Enabled:           true,
		ConnectionId:      RandomLowLetterAndNumberString(24),
	}}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if siemLogConfiguration == nil {
		t.Errorf("Should not have received a nil SiemLogConfiguration instance")
	}
}

func TestClientReadExistingSiemLogConfiguration(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	ID := RandomLowLetterAndNumberString(25)
	assetId := RandomNumbersExcludingZeroString(10)
	configurationName := RandomLetterAndNumberString(20)
	provider := RandomCapitalLetterString(6)
	dataset := RandomCapitalLetterString(6)
	connectionId := RandomLowLetterAndNumberString(24)

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemLogConfiguration, ID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"assetId": "%s",
														"configurationName": "%s",
														"provider": "%s",
														"datasets": [
															"%s"
														],
														"enabled": true,
														"connectionId": "%s"
													}
												]
											}`,
			RandomLowLetterAndNumberString(24),
			assetId,
			configurationName,
			provider,
			dataset,
			connectionId,
		)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemLogConfiguration, _, err := client.ReadSiemLogConfiguration(ID, "")

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemLogConfiguration == nil) || (siemLogConfiguration != nil && len(siemLogConfiguration.Data) != 1) {
		t.Errorf("Should have received only one SiemLogConfiguration")
	}
}

func TestClientReadNonExistingSiemLogConfiguration(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	ID := RandomLowLetterAndNumberString(25)
	responseStatusCode := 400

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemLogConfiguration, ID)

	var responseBody = fmt.Sprintf(`{
											"errors": [
												{
													"status": %d,
													"id": "%s",
													"source": {
														"pointer": "/%s/%s"
													},
													"title": "Bad Request"
												}
											]
										}`, responseStatusCode, RandomLowLetterAndNumberString(24), endpointSiemConnection, ID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(responseStatusCode)
		_, err := rw.Write([]byte(responseBody))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemLogConfiguration, _, err := client.ReadSiemLogConfiguration(ID, "")

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if strings.Compare(err.Error(), fmt.Sprintf("received failure response for operation: %s on SIEM log configuration\nstatus code: %d\nbody: %s",
		ReadSiemLogConfiguration, responseStatusCode, responseBody)) != 0 {
		t.Errorf("Should have received a response body for all responses different from %d, but received %s", 200, err)
	}
	if siemLogConfiguration != nil {
		t.Errorf("Should not have received a SiemLogConfiguration instance")
	}
}

func TestClientUpdateExistingSiemLogConfiguration(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	ID := RandomLowLetterAndNumberString(24)
	assetId := RandomNumbersExcludingZeroString(10)
	configurationName := RandomLetterAndNumberString(20)
	provider := RandomCapitalLetterString(6)
	dataset := RandomCapitalLetterString(6)
	connectionId := RandomLowLetterAndNumberString(24)

	endpoint := fmt.Sprintf("/%s/%s?caid=%s", endpointSiemLogConfiguration, ID, assetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"assetId": "%s",
														"configurationName": "%s",
														"provider": "%s",
														"datasets": [
															"%s"
														],
														"enabled": true,
														"connectionId": "%s"
													}
												]
											}`,
			ID,
			assetId,
			configurationName,
			provider,
			dataset,
			connectionId,
		)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	sent := SiemLogConfigurationData{
		ID:                ID,
		AssetID:           assetId,
		ConfigurationName: configurationName,
		Provider:          provider,
		Datasets:          append(make([]interface{}, 0), dataset),
		Enabled:           true,
		ConnectionId:      connectionId,
	}

	siemLogConfiguration, _, err := client.UpdateSiemLogConfiguration(&SiemLogConfiguration{Data: []SiemLogConfigurationData{sent}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemLogConfiguration == nil) || (siemLogConfiguration != nil && len(siemLogConfiguration.Data) != 1) {
		t.Errorf("Should have received only one SiemLogConfiguration")
	}

	var received = siemLogConfiguration.Data[0]
	if !reflect.DeepEqual(sent, received) {
		t.Errorf("Returned data should be same as sent data")
	}
}

func TestClientDeleteSiemLogConfigurationSuccess(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	responseStatusCode := 200
	log.Printf("[DEBUG] Executing success scenario with response code: %d and no message", responseStatusCode)
	err := ClientDeleteSiemLogConfigurationBase(t, RandomLowLetterAndNumberString(25), responseStatusCode)
	if err != nil {
		t.Errorf("Error is not expected for status %d", responseStatusCode)
	}
}

func TestClientDeleteSiemLogConfigurationFailure(t *testing.T) {
	m := make(map[int]string)
	m[400] = "Bad Request"
	m[404] = "Not Found"

	for k, v := range m {
		log.Printf("======================== BEGIN TEST ========================")
		log.Printf("[DEBUG] Executing failure scenario with response code: %d and message: %s", k, v)

		ID := ""
		if k == 400 {
			ID = RandomLowLetterAndNumberString(26)
		} else {
			ID = RandomLowLetterAndNumberString(25)
		}

		err := ClientDeleteSiemLogConfigurationBase(t, ID, k)
		if err == nil {
			t.Errorf("Error is expected for status %d", k)
		}

		if !strings.Contains(err.Error(), fmt.Sprintf(v)) {
			t.Errorf("Should have received a %s response, got: %s", v, err)
		}
	}
}

func ClientDeleteSiemLogConfigurationBase(t *testing.T, ID string, responseStatusCode int) error {
	log.Printf("[DEBUG] Running test client_siem_connection.TestDeleteSiemLogConfiguration with response status code: %d", responseStatusCode)
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	responseID := RandomLowLetterAndNumberString(24)

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemLogConfiguration, ID)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}

		rw.WriteHeader(responseStatusCode)

		switch responseStatusCode {
		case 400:
			_, err := rw.Write([]byte(fmt.Sprintf(`{
				"errors": [
					{
						"status": %d,
						"id": "%s",
						"source": {
							"pointer": "/%s/id/%s"
						},
						"title": "Bad Request"
				}
				]
			}`, responseStatusCode, responseID, endpointSiemLogConfiguration, ID)))
			if err != nil {
				return
			}
		case 404:
			_, err := rw.Write([]byte(fmt.Sprintf(`{
				"errors": [
					{
						"status": %d,
						"id": "%s",
						"source": {
							"pointer": "/%s/id/%s"
						},
						"title": "Not Found"
				}
				]
			}`, responseStatusCode, responseID, endpointSiemLogConfiguration, ID)))
			if err != nil {
				return
			}
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.DeleteSiemLogConfiguration(ID, "")
	return err
}
