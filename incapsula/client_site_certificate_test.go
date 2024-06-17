package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientRequestSiteCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	requestSiteCertificateResponse, diag := client.RequestSiteCertificate(123, "DNS")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
	if requestSiteCertificateResponse != nil {
		t.Errorf("Should have received a nil SiteCertificateDTO instance")
	}
}

func TestClientRequestSiteCertificateInternalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, diag := client.RequestSiteCertificate(123, "DNS")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500, error") {
		t.Errorf("Should have received an error")
	}
}

func TestClientRequestSiteCertificateErrorsInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n    \"errors\": [\n        {\n            \"status\": 400,\n            \"id\": \"ad7e0ec7b8b45cc7\",\n            \"source\": {\n                \"pointer\": \"/certificates-ui/v3/sites//certificates/managed\"\n            },\n            \"title\": \"Bad Request\",\n            \"detail\": \"JSON parse error: Unexpected end-of-input: expected close marker for Object \"\n        }\n    ]\n}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountSSLSettingsResponse, diag := client.RequestSiteCertificate(123, "DNS")
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if accountSSLSettingsResponse.Errors == nil || accountSSLSettingsResponse.Errors[0].Status != 400 || accountSSLSettingsResponse.Data != nil {
		t.Errorf("Got unexpected response")
	}
}

func TestClientRequestSiteCertificateDataInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[{\"siteId\":1088225110,\"defaultValidationMethod\":\"CNAME\",\"certificateDetails\":[{\"id\":3476,\"name\":\"ATLAS_845-1717566889024\",\"status\":\"IN_PROCESS\",\"expirationDate\":1749113691000,\"inRenewal\":false,\"sans\":[{\"sanId\":34,\"sanValue\":\"domain-8226.com\",\"validationMethod\":\"CNAME\",\"expirationDate\":null,\"status\":\"PENDING_USER_ACTION\",\"statusDate\":1717577707000,\"numSitesCovered\":1,\"verificationCode\":\"globalsign-domain-verification=A477A8393D17A55ECB2964402\",\"cnameValidationValue\":\"dh8pzxg.devImpervaDns.com\",\"autoValidation\":false,\"approverFqdn\":\"domain-8226.com\",\"validationEmail\":null,\"domainIds\":[7]}],\"level\":\"SITE\"}]}]}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	requestSiteCertificateResponse, diag := client.RequestSiteCertificate(123, "DNS")
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if requestSiteCertificateResponse.Errors != nil || requestSiteCertificateResponse.Data == nil || requestSiteCertificateResponse.Data[0].DefaultValidationMethod != "CNAME" {
		t.Errorf("Got unexpected response")
	}
}
func TestClientGetSiteCertificateRequestStatusDataInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[{\"siteId\":1088225110,\"defaultValidationMethod\":\"CNAME\",\"certificateDetails\":[{\"id\":3476,\"name\":\"ATLAS_845-1717566889024\",\"status\":\"IN_PROCESS\",\"expirationDate\":1749113691000,\"inRenewal\":false,\"sans\":[{\"sanId\":34,\"sanValue\":\"domain-8226.com\",\"validationMethod\":\"CNAME\",\"expirationDate\":null,\"status\":\"PENDING_USER_ACTION\",\"statusDate\":1717577707000,\"numSitesCovered\":1,\"verificationCode\":\"globalsign-domain-verification=A477A8393D17A55ECB2964402\",\"cnameValidationValue\":\"dh8pzxg.devImpervaDns.com\",\"autoValidation\":false,\"approverFqdn\":\"domain-8226.com\",\"validationEmail\":null,\"domainIds\":[7]}],\"level\":\"SITE\"}]}]}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteCertificateResponse, diag := client.GetSiteCertificateRequestStatus(123)
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if siteCertificateResponse.Errors != nil || siteCertificateResponse.Data == nil || siteCertificateResponse.Data[0].DefaultValidationMethod != "CNAME" {
		t.Errorf("Got unexpected response")
	}
}

func TestClientGetSiteCertificateRequestStatusErrorsInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n    \"errors\": [\n        {\n            \"status\": 400,\n            \"id\": \"ad7e0ec7b8b45cc7\",\n            \"source\": {\n                \"pointer\": \"/v3/account/ssl-settings\"\n            },\n            \"title\": \"Bad Request\",\n            \"detail\": \"JSON parse error: Unexpected end-of-input: expected close marker for Object \"\n        }\n    ]\n}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteCertificateResponse, diag := client.GetSiteCertificateRequestStatus(123)
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if siteCertificateResponse.Errors == nil || siteCertificateResponse.Errors[0].Status != 400 || siteCertificateResponse.Data != nil {
		t.Errorf("Got unexpected response")
	}
}

func TestClientGetSiteCertificateRequestStatusError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, diag := client.GetSiteCertificateRequestStatus(123)
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500, error") {
		t.Errorf("Should have received an error")
	}
}

func TestClientGetSiteCertificateRequestStatusBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	updateAccountSSLSettingsResponse, diag := client.GetSiteCertificateRequestStatus(123)
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
	if updateAccountSSLSettingsResponse != nil {
		t.Errorf("Should have received a nil addAccountResponse instance")
	}
}

func TestClientDeleteRequestSiteCertificateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	_, diag := client.DeleteRequestSiteCertificate(123)
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
}

func TestClientDeleteRequestSiteCertificateError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, diag := client.DeleteRequestSiteCertificate(123)
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500") {
		t.Errorf("Should have received an error")
	}
}

func TestClientValidateDomainBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	diag := client.ValidateDomains(123, []int{222})
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
}
func TestClientValidateDomainError500(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSSLValidationSuffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSSLValidationSuffix, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	diag := client.ValidateDomains(123, []int{222})
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500") {
		t.Errorf("Should have received an error")
	}
}

func TestClientValidateDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s%d%s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix) {
			t.Errorf("Should have hit /%s%d%s endpoint. Got: %s", endpointSiteCertV3BasePath, 123, endpointSiteCertV3Suffix, req.URL.String())
		}
		rw.Write([]byte("{\"data\":[{\"siteId\":1088225110,\"defaultValidationMethod\":\"CNAME\"}]}"))
		rw.WriteHeader(200)
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	deleteSiteCertificateResponse, diag := client.DeleteRequestSiteCertificate(123)
	if diag != nil || diag.HasError() {
		t.Errorf("Should not received an error")
	}
	if deleteSiteCertificateResponse.Errors != nil || deleteSiteCertificateResponse.Data == nil || deleteSiteCertificateResponse.Data[0].DefaultValidationMethod != "CNAME" {
		t.Errorf("Got unexpected response")
	}
}
