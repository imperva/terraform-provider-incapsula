package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClientCreateSiemConnectionBadConfig(t *testing.T) {
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

	siemConnection, _, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        RandomNumbersExcludingZeroString(10),
		ConnectionName: RandomLetterAndNumberString(20),
		StorageType:    RandomCapitalLetterString(10),
		ConnectionInfo: S3ConnectionInfo{
			AccessKey: RandomCapitalLetterAndNumberString(20),
			SecretKey: RandomLetterAndNumberString(40),
			Path:      RandomLowLetterString(20) + "/" + RandomLowLetterString(10),
		},
	}}})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error from Incapsula service when executing %s operation on SIEM connection:", CreateSiemConnection)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siemConnection != nil {
		t.Errorf("Should have received a nil SiemConnection instance")
	}
}

func TestClientCreateSiemConnectionBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	assetId := RandomNumbersExcludingZeroString(10)

	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemConnection, assetId)

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

	siemConnection, _, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        assetId,
		ConnectionName: RandomLetterAndNumberString(20),
		StorageType:    RandomCapitalLetterString(10),
		ConnectionInfo: S3ConnectionInfo{
			AccessKey: RandomCapitalLetterAndNumberString(20),
			SecretKey: RandomLetterAndNumberString(40),
			Path:      RandomLowLetterString(20) + "/" + RandomLowLetterString(10),
		},
	}}})

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("error obtained")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if siemConnection != nil {
		t.Errorf("Should have received a nil SiemConnection instance")
	}
}

func TestClientCreateSiemConnectionBadRequest(t *testing.T) {
	responseStatusCode := 400
	m := make(map[string]string)
	m["Already existing connection"] = fmt.Sprintf(`{
					"errors": [
						{
							"status": %d,
							"id": "%s",
							"source": {
								"pointer": "/%s"
							},
							"title": "Bad Request"
						}
					]
				}`, responseStatusCode, RandomLowLetterAndNumberString(24), endpointSiemConnection)
	m["Wrong connectionInfo credentials"] = fmt.Sprintf(`{
					"errors": [
						{
							"status": %d,
							"id": "%s",
							"code": "1000",
							"source": {
								"pointer": "/%s"
							},
							"title": "Bad Request"
						}
					]
				}`, responseStatusCode, RandomLowLetterAndNumberString(24), endpointSiemConnection)
	m["Wrong storageType"] = ""
	for k, v := range m {
		log.Printf("======================== BEGIN TEST ========================")
		log.Printf("[DEBUG] Executing failure scenario (%s) with response code: %d and message: %s", k, responseStatusCode, v)

		apiID := RandomCapitalLetterAndNumberString(20)
		apiKey := RandomLetterAndNumberString(40)
		assetId := RandomNumbersExcludingZeroString(10)
		endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemConnection, assetId)

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

		storageType := ""
		switch k {
		case "Already existing connection":
			storageType = "CUSTOMER_S3"
		case "Wrong connectionInfo credentials":
			storageType = "CUSTOMER_S3_ARN"
		case "Wrong storageType":
			storageType = "CUSTOMER_S3_" + RandomLetterAndNumberString(1)
		}

		siemConnection, _, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
			AssetID:        assetId,
			ConnectionName: RandomLetterAndNumberString(20),
			StorageType:    storageType,
			ConnectionInfo: S3ConnectionInfo{
				AccessKey: RandomCapitalLetterAndNumberString(20),
				SecretKey: RandomLetterAndNumberString(40),
				Path:      RandomLowLetterString(20) + "/" + RandomLowLetterString(10),
			},
		}}})

		if err == nil {
			t.Errorf("Should have received an error")
		}
		if !strings.HasPrefix(err.Error(), fmt.Sprintf("received failure response for operation: %s on SIEM connection\nstatus code: %d",
			CreateSiemConnection, responseStatusCode)) {
			t.Errorf("Should have received a response body for all responses different from %d, but received %s", 200, err)
		}
		if siemConnection != nil {
			t.Errorf("Should have received a nil SiemConnection instance")
		}
	}
}

func TestClientCreateValidSiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	assetId := RandomNumbersExcludingZeroString(10)
	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemConnection, assetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(201)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "CUSTOMER_S3",
														"connectionInfo": {
															"accessKey": "%s",
															"secretKey": "%s",
															"path": "%s/%s"
														}
													}
												]
											}`,
			RandomLowLetterAndNumberString(24),
			RandomLetterAndNumberString(20),
			RandomNumbersExcludingZeroString(10),
			RandomCapitalLetterAndNumberString(20),
			RandomLetterAndNumberString(40),
			RandomLowLetterString(20),
			RandomLowLetterString(10))))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemConnection, _, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        assetId,
		ConnectionName: RandomLetterAndNumberString(20),
		StorageType:    RandomCapitalLetterString(10),
		ConnectionInfo: S3ConnectionInfo{
			AccessKey: RandomCapitalLetterAndNumberString(20),
			SecretKey: RandomLetterAndNumberString(40),
			Path:      RandomLowLetterString(20) + "/" + RandomLowLetterString(10),
		},
	}}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if siemConnection == nil {
		t.Errorf("Should not have received a nil SiemConnection instance")
	}
}

func TestClientReadExistingSiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	ID := RandomLowLetterAndNumberString(25)

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemConnection, ID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "CUSTOMER_S3",
														"connectionInfo": {
															"accessKey": "%s",
															"secretKey": "%s",
															"path": "%s/%s"
														}
													}
												]
											}`,
			RandomLowLetterAndNumberString(24),
			RandomLetterAndNumberString(20),
			RandomNumbersExcludingZeroString(10),
			RandomCapitalLetterAndNumberString(20),
			RandomLetterAndNumberString(40),
			RandomLowLetterString(20),
			RandomLowLetterString(10))))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemConnection, _, err := client.ReadSiemConnection(ID, "")

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}
}

func TestClientReadNonExistingSiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	ID := RandomLowLetterAndNumberString(25)
	responseStatusCode := 404

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemConnection, ID)

	var responseBody = fmt.Sprintf(`{
											"errors": [
												{
													"status": %d,
													"id": "%s",
													"source": {
														"pointer": "/%s/%s"
													},
													"title": "Not Found"
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

	siemConnection, _, err := client.ReadSiemConnection(ID, "")

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("received failure response for operation: %s on SIEM connection\nstatus code: %d",
		ReadSiemConnection, responseStatusCode)) {
		t.Errorf("Should have received a response body for all responses different from %d, but received %s", 200, err)
	}
	if siemConnection != nil {
		t.Errorf("Should not have received a SiemConnection instance")
	}
}

func TestClientUpdateExistingS3SiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	responseID := RandomLowLetterAndNumberString(24)
	responseConnectionName := RandomLetterAndNumberString(20)
	responseStorageType := "CUSTOMER_S3"
	responseAssetId := RandomNumbersExcludingZeroString(10)
	responseAccessKey := RandomCapitalLetterAndNumberString(20)
	responseSecretKey := RandomLetterAndNumberString(40)
	responsePath1 := RandomLowLetterString(20)
	responsePath2 := RandomLowLetterString(10)

	endpoint := fmt.Sprintf("/%s/%s?caid=%s", endpointSiemConnection, responseID, responseAssetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",	
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "%s",
														"connectionInfo": {
															"accessKey": "%s",
															"secretKey": "%s",
															"path": "%s/%s"
														}
													}
												]
											}`,
			responseID,
			responseConnectionName,
			responseAssetId,
			responseStorageType,
			responseAccessKey,
			responseSecretKey,
			responsePath1,
			responsePath2)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	sent := SiemConnectionData{
		ID:             responseID,
		AssetID:        responseAssetId,
		ConnectionName: responseConnectionName,
		StorageType:    responseStorageType,
		ConnectionInfo: S3ConnectionInfo{
			AccessKey: RandomCapitalLetterAndNumberString(20),
			SecretKey: RandomLetterAndNumberString(40),
			Path:      responsePath1 + "/" + responsePath2,
		},
	}

	siemConnection, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{sent}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}

	var received = siemConnection.Data[0]
	if (received.ID != sent.ID) || (received.AssetID != sent.AssetID) || (received.StorageType != sent.StorageType) || (received.ConnectionName != sent.ConnectionName) || (received.ConnectionInfo.(S3ConnectionInfo).Path != sent.ConnectionInfo.(S3ConnectionInfo).Path) || (received.ConnectionInfo.(S3ConnectionInfo).AccessKey == sent.ConnectionInfo.(S3ConnectionInfo).AccessKey) || (received.ConnectionInfo.(S3ConnectionInfo).SecretKey == sent.ConnectionInfo.(S3ConnectionInfo).SecretKey) {
		t.Errorf("Returned data should be same as sent data with different accessKey and secretKey")
	}
}

func TestClientUpdateExistingS3ArnSiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	responseID := RandomLowLetterAndNumberString(24)
	responseConnectionName := RandomLetterAndNumberString(20)
	responseStorageType := "CUSTOMER_S3_ARN"
	responseAssetId := RandomNumbersExcludingZeroString(10)
	responsePath1 := RandomLowLetterString(20)
	responsePath2 := RandomLowLetterString(10)

	endpoint := fmt.Sprintf("/%s/%s?caid=%s", endpointSiemConnection, responseID, responseAssetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "%s",
														"connectionInfo": {
															"path": "%s/%s"
														}
													}
												]
											}`,
			responseID,
			responseConnectionName,
			responseAssetId,
			responseStorageType,
			responsePath1,
			responsePath2)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	sent := SiemConnectionData{
		ID:             responseID,
		AssetID:        responseAssetId,
		ConnectionName: responseConnectionName,
		StorageType:    responseStorageType,
		ConnectionInfo: S3ConnectionInfo{
			Path: responsePath1 + "/" + responsePath2,
		},
	}

	siemConnection, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{sent}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}

	var received = siemConnection.Data[0]
	if !reflect.DeepEqual(sent, received) {
		t.Errorf("Returned data should be same as sent data")
	}
}

func TestClientDeleteSiemConnectionSuccess(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	responseStatusCode := 200
	log.Printf("[DEBUG] Executing success scenario with response code: %d and no message", responseStatusCode)
	err := ClientDeleteSiemConnectionBase(t, responseStatusCode)
	if err != nil {
		t.Errorf("Error is not expected for status %d", responseStatusCode)
	}
}

func TestClientDeleteSiemConnectionFailure(t *testing.T) {
	m := make(map[int]string)
	m[400] = "Bad Request"
	m[401] = "Invalid Bearer token"
	m[404] = "Not Found"

	for k, v := range m {
		log.Printf("======================== BEGIN TEST ========================")
		log.Printf("[DEBUG] Executing failure scenario with response code: %d and message: %s", k, v)
		err := ClientDeleteSiemConnectionBase(t, k)
		if err == nil {
			t.Errorf("Error is expected for status %d", k)
		}

		if !strings.Contains(err.Error(), v) {
			t.Errorf("Should have received a %s response, got: %s", v, err)
		}
	}
}

func ClientDeleteSiemConnectionBase(t *testing.T, responseStatusCode int) error {
	log.Printf("[DEBUG] Running test client_siem_connection.TestDeleteSiemConnection with response status code: %d", responseStatusCode)
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	responseID := RandomLowLetterAndNumberString(24)
	ID := RandomLowLetterAndNumberString(25)

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemConnection, ID)
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
				}`, responseStatusCode, responseID, endpointSiemConnection, ID)))
			if err != nil {
				return
			}
		case 401:
			_, err := rw.Write([]byte(`{
				"errMsg": "Invalid Bearer token",
				"errCode": 10001
			}`))
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
			}`, responseStatusCode, responseID, endpointSiemConnection, ID)))
			if err != nil {
				return
			}

		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, err := client.DeleteSiemConnection(ID, "")
	return err
}

func TestClientCreateValidSiemSplunkConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	host := "test.splunk.com"
	port := 8080
	token, _ := uuid.GenerateUUID()

	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	assetId := RandomNumbersExcludingZeroString(10)
	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemConnection, assetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(201)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "CUSTOMER_SPLUNK",
														"connectionInfo": {
															"host": "%s",
															"port": %d,
															"token": "%s",
															"disableCertVerification": %t
														}
													}
												]
											}`,
			RandomLowLetterAndNumberString(24),
			RandomLetterAndNumberString(20),
			RandomNumbersExcludingZeroString(10),
			host,
			port,
			token,
			true)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemConnection, _, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        assetId,
		ConnectionName: RandomLetterAndNumberString(20),
		StorageType:    RandomCapitalLetterString(10),
		ConnectionInfo: SplunkConnectionInfo{
			Host:                    host,
			Port:                    port,
			Token:                   token,
			DisableCertVerification: true,
		},
	}}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if siemConnection == nil {
		t.Errorf("Should not have received a nil SiemConnection instance")
	}
}

func TestClientReadExistingSiemSplunkConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	host := "test.splunk.com"
	port := 8080
	token, _ := uuid.GenerateUUID()
	disableCertVerification := true

	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	ID := RandomLowLetterAndNumberString(25)

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemConnection, ID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "CUSTOMER_SPLUNK",
														"connectionInfo": {
															"host": "%s",
															"port": %d,
															"token": "%s",
															"disableCertVerification": %t
														}
													}
												]
											}`,
			RandomLowLetterAndNumberString(24),
			RandomLetterAndNumberString(20),
			RandomNumbersExcludingZeroString(10),
			host,
			port,
			token,
			disableCertVerification)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemConnection, _, err := client.ReadSiemConnection(ID, "")

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}
}

func TestClientUpdateExistingSplunkSiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	responseID := RandomLowLetterAndNumberString(24)
	responseConnectionName := RandomLetterAndNumberString(20)
	responseStorageType := "CUSTOMER_SPLUNK"
	responseAssetId := RandomNumbersExcludingZeroString(10)
	responseHost := "splunk.test.com"
	responsePort := 8080
	responseToken, _ := uuid.GenerateUUID()
	responseDisableCert := true

	endpoint := fmt.Sprintf("/%s/%s?caid=%s", endpointSiemConnection, responseID, responseAssetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
												"data": [
													{
														"id": "%s",	
														"connectionName": "%s",
														"assetId": "%s",
														"storageType": "%s",
														"connectionInfo": {
															"host": "%s",
															"port": %d,
															"token": "%s",
															"disableCertVerification": %t
														}
													}
												]
											}`,
			responseID,
			responseConnectionName,
			responseAssetId,
			responseStorageType,
			responseHost,
			responsePort,
			responseToken,
			responseDisableCert)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	sent := SiemConnectionData{
		ID:             responseID,
		AssetID:        responseAssetId,
		ConnectionName: responseConnectionName,
		StorageType:    responseStorageType,
		ConnectionInfo: SplunkConnectionInfo{
			Host:                    responseHost,
			Port:                    responsePort,
			Token:                   responseToken,
			DisableCertVerification: true,
		},
	}

	siemConnection, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{sent}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}

	var received = siemConnection.Data[0]
	if (received.ID != sent.ID) ||
		(received.AssetID != sent.AssetID) ||
		(received.StorageType != sent.StorageType) ||
		(received.ConnectionName != sent.ConnectionName) ||
		(received.ConnectionInfo.(SplunkConnectionInfo).Host != sent.ConnectionInfo.(SplunkConnectionInfo).Host) ||
		(received.ConnectionInfo.(SplunkConnectionInfo).Port != sent.ConnectionInfo.(SplunkConnectionInfo).Port) ||
		(received.ConnectionInfo.(SplunkConnectionInfo).Token != sent.ConnectionInfo.(SplunkConnectionInfo).Token) ||
		(received.ConnectionInfo.(SplunkConnectionInfo).DisableCertVerification != sent.ConnectionInfo.(SplunkConnectionInfo).DisableCertVerification) {
		t.Errorf("Returned data should be same as sent data")
	}
}

func TestClientCreateValidSiemSftpConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	host := "test.sftp.com"
	path := "/data/testing/accounts/53548213"
	username := "testuser"
	password := "testpassword"

	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	assetId := RandomNumbersExcludingZeroString(10)
	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemConnection, assetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(201)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
            "data": [
                {
                    "id": "%s",
                    "connectionName": "%s",
                    "assetId": "%s",
                    "storageType": "CUSTOMER_SFTP",
                    "connectionInfo": {
                        "host": "%s",
                        "path": "%s",
                        "username": "%s",
                        "password": "%s"
                    }
                }
            ]
        }`,
			RandomLowLetterAndNumberString(24),
			RandomLetterAndNumberString(20),
			RandomNumbersExcludingZeroString(10),
			host,
			path,
			username,
			password)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemConnection, _, err := client.CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        assetId,
		ConnectionName: RandomLetterAndNumberString(20),
		StorageType:    "CUSTOMER_SFTP",
		ConnectionInfo: SftpConnectionInfo{
			Host:     host,
			Path:     path,
			Username: username,
			Password: password,
		},
	}}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if siemConnection == nil {
		t.Errorf("Should not have received a nil SiemConnection instance")
	}
}

func TestClientReadExistingSiemSftpConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	host := "test.sftp.com"
	path := "/data/testing/accounts/53548213"
	username := "testuser"
	password := "testpassword"

	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)
	ID := RandomLowLetterAndNumberString(25)

	endpoint := fmt.Sprintf("/%s/%s", endpointSiemConnection, ID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
            "data": [
                {
                    "id": "%s",
                    "connectionName": "%s",
                    "assetId": "%s",
                    "storageType": "CUSTOMER_SFTP",
                    "connectionInfo": {
                        "host": "%s",
                        "path": "%s",
                        "username": "%s",
                        "password": "%s"
                    }
                }
            ]
        }`,
			RandomLowLetterAndNumberString(24),
			RandomLetterAndNumberString(20),
			RandomNumbersExcludingZeroString(10),
			host,
			path,
			username,
			password)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siemConnection, _, err := client.ReadSiemConnection(ID, "")

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}
}

func TestClientUpdateExistingSftpSiemConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	apiID := RandomCapitalLetterAndNumberString(20)
	apiKey := RandomLetterAndNumberString(40)

	responseID := RandomLowLetterAndNumberString(24)
	responseConnectionName := RandomLetterAndNumberString(20)
	responseStorageType := "CUSTOMER_SFTP"
	responseAssetId := RandomNumbersExcludingZeroString(10)
	responseHost := "sftp.test.com"
	responsePath := "/data/testing/accounts/53548213"
	responseUsername := "updateduser"
	responsePassword := "updatedpassword"

	endpoint := fmt.Sprintf("/%s/%s?caid=%s", endpointSiemConnection, responseID, responseAssetId)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.WriteHeader(200)
		_, err := rw.Write([]byte(fmt.Sprintf(`{
            "data": [
                {
                    "id": "%s",
                    "connectionName": "%s",
                    "assetId": "%s",
                    "storageType": "%s",
                    "connectionInfo": {
                        "host": "%s",
                        "path": "%s",
                        "username": "%s",
                        "password": "%s"
                    }
                }
            ]
        }`,
			responseID,
			responseConnectionName,
			responseAssetId,
			responseStorageType,
			responseHost,
			responsePath,
			responseUsername,
			responsePassword)))
		if err != nil {
			return
		}
	}))

	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	sent := SiemConnectionData{
		ID:             responseID,
		AssetID:        responseAssetId,
		ConnectionName: responseConnectionName,
		StorageType:    responseStorageType,
		ConnectionInfo: SftpConnectionInfo{
			Host:     responseHost,
			Path:     responsePath,
			Username: responseUsername,
			Password: responsePassword,
		},
	}

	siemConnection, _, err := client.UpdateSiemConnection(&SiemConnection{Data: []SiemConnectionData{sent}})

	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if (siemConnection == nil) || (siemConnection != nil && len(siemConnection.Data) != 1) {
		t.Errorf("Should have received only one SiemConnection")
	}

	var received = siemConnection.Data[0]
	if !reflect.DeepEqual(sent, received) {
		t.Errorf("Returned data should be same as sent data")
	}
}

func newSiemHttpConnectionTestClient(serverURL string) *Client {
	config := &Config{APIID: RandomCapitalLetterAndNumberString(20), APIKey: RandomLetterAndNumberString(40), BaseURL: serverURL, BaseURLRev2: serverURL, BaseURLAPI: serverURL}
	return &Client{config: config, httpClient: &http.Client{}}
}

func siemHttpConnectionResponse(connectionInfo string) string {
	return fmt.Sprintf(`{"data":[{"id":"%s","connectionName":"%s","assetId":"%s","storageType":"CUSTOMER_HTTP","connectionInfo":%s}]}`,
		RandomLowLetterAndNumberString(24), RandomLetterAndNumberString(20), RandomNumbersExcludingZeroString(10), connectionInfo)
}

func captureSiemHttpConnectionCreateRequest(t *testing.T, connectionInfo HttpConnectionInfo) map[string]interface{} {
	assetId := RandomNumbersExcludingZeroString(10)
	endpoint := fmt.Sprintf("/%s/?caid=%s", endpointSiemConnection, assetId)
	var requestBody map[string]interface{}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
			t.Errorf("Request body should be valid JSON: %s", err)
		}
		rw.WriteHeader(201)
		rw.Write([]byte(siemHttpConnectionResponse(`{"url":"https://example.com/ingest","authentication":{"type":"BASIC_AUTH","username":"admin","password":"cipher"},"disableCertVerification":false}`)))
	}))
	defer server.Close()

	_, _, err := newSiemHttpConnectionTestClient(server.URL).CreateSiemConnection(&SiemConnection{Data: []SiemConnectionData{{
		AssetID:        assetId,
		ConnectionName: RandomLetterAndNumberString(20),
		StorageType:    StorageTypeCustomerHttp,
		ConnectionInfo: connectionInfo,
	}}})
	if err != nil {
		t.Fatalf("Should not have received an error, got: %s", err)
	}
	return requestBody
}

func requestConnectionInfo(requestBody map[string]interface{}) map[string]interface{} {
	return requestBody["data"].([]interface{})[0].(map[string]interface{})["connectionInfo"].(map[string]interface{})
}

func TestClientCreateSiemHttpConnectionWithHeadersAuthSendsOrderedAuthHeaders(t *testing.T) {
	requestBody := captureSiemHttpConnectionCreateRequest(t, HttpConnectionInfo{
		Url: "https://example.com/ingest",
		Authentication: HttpAuthenticationInfo{Type: HttpAuthTypeHeaders, AuthHeaders: []AuthHeaderInfo{
			{HeaderName: "X-goog-api-key", SecretFormat: "{secret}", Secret: "first"},
			{HeaderName: "X-Webhook-Access-Key", SecretFormat: "Bearer {secret}", Secret: "second"},
		}},
		Headers: map[string]string{"Content-Type": "application/json"},
	})

	if storageType := requestBody["data"].([]interface{})[0].(map[string]interface{})["storageType"]; storageType != StorageTypeCustomerHttp {
		t.Errorf("Should have sent storageType %s, got: %v", StorageTypeCustomerHttp, storageType)
	}
	authentication := requestConnectionInfo(requestBody)["authentication"].(map[string]interface{})
	authHeaders := authentication["authHeaders"].([]interface{})
	if len(authHeaders) != 2 || authHeaders[0].(map[string]interface{})["headerName"] != "X-goog-api-key" || authHeaders[1].(map[string]interface{})["headerName"] != "X-Webhook-Access-Key" {
		t.Errorf("Should have sent both auth headers in declaration order, got: %v", authHeaders)
	}
	if _, hasUsername := authentication["username"]; hasUsername {
		t.Errorf("Should not have sent username for HEADERS auth")
	}
	if _, hasPassword := authentication["password"]; hasPassword {
		t.Errorf("Should not have sent password for HEADERS auth")
	}
}

func TestClientCreateSiemHttpConnectionWithBasicAuthOmitsAuthHeaders(t *testing.T) {
	requestBody := captureSiemHttpConnectionCreateRequest(t, HttpConnectionInfo{
		Url:            "https://example.com/ingest",
		Authentication: HttpAuthenticationInfo{Type: HttpAuthTypeBasicAuth, Username: "admin", Password: "pass"},
	})

	authentication := requestConnectionInfo(requestBody)["authentication"].(map[string]interface{})
	if authentication["username"] != "admin" || authentication["password"] != "pass" {
		t.Errorf("Should have sent username and password, got: %v", authentication)
	}
	if _, hasAuthHeaders := authentication["authHeaders"]; hasAuthHeaders {
		t.Errorf("Should not have sent authHeaders for BASIC_AUTH")
	}
}

func TestClientCreateSiemHttpConnectionWithEmptyHeadersOmitsHeaders(t *testing.T) {
	requestBody := captureSiemHttpConnectionCreateRequest(t, HttpConnectionInfo{
		Url:            "https://example.com/ingest",
		Authentication: HttpAuthenticationInfo{Type: HttpAuthTypeBasicAuth, Username: "admin", Password: "pass"},
		Headers:        map[string]string{},
	})

	if _, hasHeaders := requestConnectionInfo(requestBody)["headers"]; hasHeaders {
		t.Errorf("Should not have sent headers when the map is empty")
	}
}

func TestSiemConnectionUnmarshalCustomerHttp(t *testing.T) {
	testCases := []struct {
		name           string
		connectionInfo string
		expected       HttpConnectionInfo
	}{
		{
			name:           "headers auth without headers map and with ciphertext secret",
			connectionInfo: `{"url":"https://example.com/ingest","authentication":{"type":"HEADERS","authHeaders":[{"headerName":"X-Key","secretFormat":"{secret}","secret":"ZW5jcnlwdGVk=="}]},"disableCertVerification":true}`,
			expected: HttpConnectionInfo{
				Url:                     "https://example.com/ingest",
				Authentication:          HttpAuthenticationInfo{Type: HttpAuthTypeHeaders, AuthHeaders: []AuthHeaderInfo{{HeaderName: "X-Key", SecretFormat: "{secret}", Secret: "ZW5jcnlwdGVk=="}}},
				DisableCertVerification: true,
			},
		},
		{
			name:           "basic auth with headers map",
			connectionInfo: `{"url":"https://example.com/ingest","authentication":{"type":"BASIC_AUTH","username":"admin","password":"cipher"},"headers":{"Content-Type":"application/json"},"disableCertVerification":false}`,
			expected: HttpConnectionInfo{
				Url:            "https://example.com/ingest",
				Authentication: HttpAuthenticationInfo{Type: HttpAuthTypeBasicAuth, Username: "admin", Password: "cipher"},
				Headers:        map[string]string{"Content-Type": "application/json"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			var connection SiemConnection
			if err := json.Unmarshal([]byte(siemHttpConnectionResponse(testCase.connectionInfo)), &connection); err != nil {
				t.Fatalf("Should not have received an error, got: %s", err)
			}
			if !reflect.DeepEqual(connection.Data[0].ConnectionInfo, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, connection.Data[0].ConnectionInfo)
			}
		})
	}
}

func TestRedactSiemConnectionSecretsMasksSecretAndPassword(t *testing.T) {
	redacted := redactSiemConnectionSecrets([]byte(`{"data":[{"connectionInfo":{"url":"https://example.com","authentication":{"authHeaders":[{"headerName":"X-Key","secret":"top-secret"}],"password":"hunter2"}}}]}`))

	if strings.Contains(redacted, "top-secret") || strings.Contains(redacted, "hunter2") {
		t.Errorf("Should have masked secret values, got: %s", redacted)
	}
	if !strings.Contains(redacted, "https://example.com") {
		t.Errorf("Should have kept non-secret values, got: %s", redacted)
	}
}
