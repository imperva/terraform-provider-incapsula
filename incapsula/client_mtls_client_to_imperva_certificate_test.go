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
// GetClientCaCertificate Tests
////////////////////////////////////////////////////////////////
func TestGetClientCaCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	certificateID := "42"
	accountID := "100"

	mTLSClientToImpervaCertificateResponse, _, err := client.GetClientCaCertificate(accountID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when reading mTLS Client CA to Imperva Certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSClientToImpervaCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSClientToImpervaCertificateResponse instance")
	}

}

func TestGetClientCaCertificateBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	accountID := "100"
	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates/%s", accountID, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSClientToImpervaCertificateResponse, _, err := client.GetClientCaCertificate(accountID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing mutual GET TLS Client To Imperva Certificate for Account ID %s", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSClientToImpervaCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSClientToImpervaCertificateResponse instance")
	}
}

func TestGetClientCaCertificateInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	accountID := "100"
	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates/%s", accountID, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"errors": [
	{
		"status": 500,
		"id": "cca667c1371c31ff",
		"source": {
		"pointer": "/v3/mtls-origin/certificates/111"
	},
		"title": "Internal Server Error",
		"detail": "Internal Server Error"
	}
]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSClientToImpervaCertificateResponse, _, err := client.GetClientCaCertificate(accountID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on fetching TLS Client to Imperva certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSClientToImpervaCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSClientToImpervaCertificateResponse instance")
	}
}

func TestGetClientCaCertificateValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	accountID := "100"
	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates/%s", accountID, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "id": 1,
  "name": "string",
  "serialNumber": "string",
  "issuer": "string",
  "creationDate": "string"
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSClientToImpervaCertificateResponse, _, err := client.GetClientCaCertificate(accountID, certificateID)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), mTLSClientToImpervaCertificateResponse)
	}
	if mTLSClientToImpervaCertificateResponse == nil {
		t.Errorf("Should not have received a nil mTLSCertificateResponse instance")
	}
	if mTLSClientToImpervaCertificateResponse.Id != 1 {
		t.Errorf("Certifcate ID doesn't match. Actual : %v", mTLSClientToImpervaCertificateResponse.Id)
	}
}

////////////////////////////////////////////////////////////////
// AddClientCaCertificate Tests
////////////////////////////////////////////////////////////////
func TestAddClientCaCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	certificateName := "test certificate"
	accountID := "100"
	certificate := []byte{}

	mTLSClientToImpervaCertificateResponse, err := client.AddClientCaCertificate(certificate, accountID, certificateName)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula while creating mutual TLS Client To Imperva Certificate for Account ID %s", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSClientToImpervaCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSClientToImpervaCertificateResponse instance")
	}

}

func TestAddClientCaCertificateBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateName := "test certificate"
	accountID := "100"
	certificate := []byte{}

	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates", accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSClientToImpervaCertificateResponse, err := client.AddClientCaCertificate(certificate, accountID, certificateName)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing ADD mutual TLS Client To Imperva Certificate for Account ID %s", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSClientToImpervaCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSClientToImpervaCertificateResponse instance")
	}
}

func TestAddClientCaCertificateInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateName := "test certificate"
	accountID := "100"
	certificate := []byte{}

	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates", accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
			"errors": [
	{
		"status": 500,
		"id": "cca667c1371c31ff",
		"source": {
		"pointer": "/v3/mtls-origin/certificates/111"
	},
		"title": "Internal Server Error",
		"detail": "Internal Server Error"
	}
]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSClientToImpervaCertificateResponse, err := client.AddClientCaCertificate(certificate, accountID, certificateName)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on create mutual TLS Client To Imperva certificate for account ID %s", accountID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSClientToImpervaCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSClientToImpervaCertificateResponse instance")
	}
}

func TestAddClientCaCertificateValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateName := "test certificate"
	accountID := "100"
	certificate := []byte{}

	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates", accountID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`[{
  "id": 1,
  "name": "string",
  "serialNumber": "string",
  "issuer": "string",
  "creationDate": "string"
}]`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSClientToImpervaCertificateResponse, err := client.AddClientCaCertificate(certificate, accountID, certificateName)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), mTLSClientToImpervaCertificateResponse)
	}
	if mTLSClientToImpervaCertificateResponse == nil {
		t.Errorf("Should not have received a nil mTLSClientToImpervaCertificateResponse instance")
	}
	if mTLSClientToImpervaCertificateResponse.Id != 1 {
		t.Errorf("Certifcate ID doesn't match. Actual : %v", mTLSClientToImpervaCertificateResponse.Id)
	}
}

////////////////////////////////////////////////////////////////
// DeleteClientCaCertificate Tests
////////////////////////////////////////////////////////////////
func TestDeleteClientCaCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	certificateID := "42"
	accountID := "100"

	err := client.DeleteClientCaCertificate(accountID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when deletingmutual TLS Client To Imperva Certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestDeleteClientCaCertificateInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	accountID := "100"
	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates/%s", accountID, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(406)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "id": "aWNB2TU91Og",
    "code": 406,
    "message": "Client CA certificate assigned to sites [65887086]. Can't delete assigned certificate"
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteClientCaCertificate(accountID, certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 406 from Incapsula service on deleting mutual TLS Client To Imperva Certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

}

func TestDeleteClientCaCertificateValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	accountID := "100"

	endpoint := fmt.Sprintf("/certificate-manager/v2/accounts/%s/client-certificates/%s", accountID, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteClientCaCertificate(accountID, certificateID)

	if err != nil {
		t.Errorf("Should not have received an error : %v\n", err.Error())
	}
}
