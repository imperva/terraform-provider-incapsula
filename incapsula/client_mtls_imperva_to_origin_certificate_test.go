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
// GetSiteMtlsCertificateAssociation Tests
////////////////////////////////////////////////////////////////
func TestGetMTLSCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	certificateID := "42"

	mTLSCertificateResponse, err := client.GetMTLSCertificate(certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when reading mTLS Imperva to Origin Certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSCertificateResponse instance")
	}

}

func TestGetMTLSCertificateBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	endpoint := fmt.Sprintf("%s/%s", endpointMTLSCertificate, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificateResponse, err := client.GetMTLSCertificate(certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error parsing mutual TLS Imperva to Origin Certificate JSON response for certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSCertificateResponse instance")
	}
}

func TestGetMTLSCertificateInvalidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	endpoint := fmt.Sprintf("%s/%s", endpointMTLSCertificate, certificateID)

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

	mTLSCertificateResponse, err := client.GetMTLSCertificate(certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 500 from Incapsula service on fetching mutual TLS Imperva to Origin certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSCertificateResponse instance")
	}
}

func TestGetMTLSCertificateValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	endpoint := fmt.Sprintf("%s/%s", endpointMTLSCertificate, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "data": [
        {
            "certificateId": 1,
            "issuedTo": "company",
            "issuedBy": "ca company",
            "validFrom": 1655011098999,
            "validUntil": 168456887638629,
            "chain": "ca cert",
            "name": "mtls-cert",
            "creationDate": 1685036052000,
            "lastUpdate": 1685036052000,
            "serialNumber": "4345667",
            "appliedSitesDetails": [
                {
                    "externalId": 122222,
                    "name": "site.incapsula.co",
                    "accountId": 2121221
                }
            ]
        }
    ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificateResponse, err := client.GetMTLSCertificate(certificateID)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), mTLSCertificateResponse)
	}
	if mTLSCertificateResponse == nil {
		t.Errorf("Should not have received a nil mTLSCertificateResponse instance")
	}
	if mTLSCertificateResponse.Id != 1 {
		t.Errorf("Certifcate ID doesn't match. Actual : %v", mTLSCertificateResponse.Id)
	}
}

////////////////////////////////////////////////////////////////
// AddMTLSCertificate, UpdateMTLSCertificate Tests
////////////////////////////////////////////////////////////////
func TestEditMTLSCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	certificate := []byte("cert")
	privateKey := []byte("pkey")
	passphrase := "passphrase"
	certificateName := "certificateName"
	inputHash := "inputHash"

	mTLSCertificateResponse, err := client.AddMTLSCertificate(certificate, privateKey, passphrase, certificateName, inputHash)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error while Create mTLS Imperva to Origin Certificate: ")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSCertificateResponse instance")
	}

}

func TestEditMTLSCertificateBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"

	certificate := []byte("cert")
	privateKey := []byte("pkey")
	passphrase := "passphrase"
	certificateName := "certificateName"
	inputHash := "inputHash"
	certificateID := "21"
	endpoint := fmt.Sprintf("%s/%s", endpointMTLSCertificate, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificateResponse, err := client.UpdateMTLSCertificate(certificateID, certificate, privateKey, passphrase, certificateName, inputHash)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprint("[ERROR] Error parsing mutual TLS Imperva to Origin Certificate JSON response:")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if mTLSCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSCertificateResponse instance")
	}
}

func TestEditMTLSCertificateApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	endpoint := fmt.Sprintf("%s", endpointMTLSCertificate)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "errors": [
        {
            "status": 400,
            "id": "16d37a3dfb2b3aff",
            "source": {
                "pointer": "/v3/mtls-origin/certificates"
            },
            "title": "Bad Request",
            "detail": "handleRequest - Got response headers:org.springframework.web.reactive.function.client.DefaultClientResponse$DefaultHeaders@20c80d50, status: 400 BAD_REQUEST, body: {\"errors\":[{\"status\":400,\"id\":\"de31602becdf6d4b\",\"source\":{\"pointer\":\"/mtls-origin/certificates\"},\"title\":\"Bad Request\",\"detail\":\"Certificate already exists\"}]}"
        }
    ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	certificate := []byte("cert")
	privateKey := []byte("pkey")
	passphrase := "passphrase"
	certificateName := "certificateName"
	inputHash := "inputHash"

	mTLSCertificateResponse, err := client.AddMTLSCertificate(certificate, privateKey, passphrase, certificateName, inputHash)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprint("[ERROR] Error status code 400 from Incapsula service on Create mutual TLS Imperva to Origin certificate")) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if mTLSCertificateResponse != nil {
		t.Errorf("Should have received a nil mTLSCertificateResponse instance")
	}
}

func TestEditMTLSCertificateValidApiConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"

	certificate := []byte("cert")
	privateKey := []byte("pkey")
	passphrase := "passphrase"
	certificateName := "certificateName"
	inputHash := "inputHash"

	endpoint := fmt.Sprintf("%s", endpointMTLSCertificate)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "data": [
        {
            "certificateId": 1,
            "issuedTo": "company",
            "issuedBy": "ca company",
            "validFrom": 1655011098999,
            "validUntil": 168456887638629,
            "chain": "ca cert",
            "name": "mtls-cert",
            "creationDate": 1685036052000,
            "lastUpdate": 1685036052000,
            "serialNumber": "4345667",
            "appliedSitesDetails": [
                {
                    "externalId": 122222,
                    "name": "site.incapsula.co",
                    "accountId": 2121221
                }
            ]
        }
    ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	mTLSCertificateResponse, err := client.AddMTLSCertificate(certificate, privateKey, passphrase, certificateName, inputHash)

	if err != nil {
		t.Errorf("Should not have received an error : %s\n, %v", err.Error(), mTLSCertificateResponse)
	}
	if mTLSCertificateResponse == nil {
		t.Errorf("Should not have received a nil mTLSCertificateResponse instance")
	}
	if mTLSCertificateResponse.Id != 1 {
		t.Errorf("Certifcate ID doesn't match. Actual : %v", mTLSCertificateResponse.Id)
	}
}

////////////////////////////////////////////////////////////////
// DeleteMTLSCertificate Tests
////////////////////////////////////////////////////////////////
func TestDeleteMTLSCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	certificateID := "42"

	err := client.DeleteMTLSCertificate(certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error from Incapsula service when deleting mTLS Imperva to Origin Certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestDeleteMTLSCertificateInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	endpoint := fmt.Sprintf("%s/%s", endpointMTLSCertificate, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
    "errors": [
        {
            "status": 400,
            "id": "bd2f35a9b684a7cf",
            "source": {
                "pointer": "/v3/mtls-origin/certificates/1"
            },
            "title": "Bad Request",
            "detail": "handleRequest - Got response headers:org.springframework.web.reactive.function.client.DefaultClientResponse$DefaultHeaders@7d57d69b, status: 400 BAD_REQUEST, body: {\"errors\":[{\"status\":400,\"id\":\"b352b66ace051df4\",\"source\":{\"pointer\":\"/mtls-origin/certificates/1\"},\"title\":\"Bad Request\",\"detail\":\"Certificate Id does not exist\"}]}"
        }
    ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteMTLSCertificate(certificateID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("[ERROR] Error status code 400 from Incapsula service on deleting mutual TLS Imperva to Origin certificate ID %s", certificateID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}

}

func TestDeleteMTLSCertificateValidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	certificateID := "42"
	endpoint := fmt.Sprintf("%s/%s", endpointMTLSCertificate, certificateID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)

		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	err := client.DeleteMTLSCertificate(certificateID)

	if err != nil {
		t.Errorf("Should not have received an error : %v\n", err.Error())
	}
}
