package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientUpdateAccountSSlSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	dto := AccountSSLSettingsDTO{}
	updateAccountSSLSettingsResponse, diag := client.UpdateAccountSSLSettings(&dto, "")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
	if updateAccountSSLSettingsResponse != nil {
		t.Errorf("Should have received a nil addAccountResponse instance")
	}
}

func TestClientUpdateAccountSSlSettingsInternalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	val := false
	imp := ImpervaCertificate{
		AddNakedDomainSanForWWWSites: &val,
	}
	dto := AccountSSLSettingsDTO{
		ImpervaCertificate: &imp,
	}
	_, diag := client.UpdateAccountSSLSettings(&dto, "")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500, error") {
		t.Errorf("Should have received an error")
	}
}

func TestClientUpdateAccountSSlSettingsErrorsInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n    \"errors\": [\n        {\n            \"status\": 400,\n            \"id\": \"ad7e0ec7b8b45cc7\",\n            \"source\": {\n                \"pointer\": \"/v3/account/ssl-settings\"\n            },\n            \"title\": \"Bad Request\",\n            \"detail\": \"JSON parse error: Unexpected end-of-input: expected close marker for Object \"\n        }\n    ]\n}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	val := false
	imp := ImpervaCertificate{
		AddNakedDomainSanForWWWSites: &val,
	}
	dto := AccountSSLSettingsDTO{
		ImpervaCertificate: &imp,
	}
	accountSSLSettingsResponse, diag := client.UpdateAccountSSLSettings(&dto, "")
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if accountSSLSettingsResponse.Errors == nil || accountSSLSettingsResponse.Errors[0].Status != 400 || accountSSLSettingsResponse.Data != nil {
		t.Errorf("Got unexpected response")
	}
}

func TestClientUpdateAccountSSlSettingsDataInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n    \"data\": [\n        {\n            \"impervaCertificate\": {\n                \"delegation\": {\n                    \"valueForCNAMEValidation\": \"_f508e2d5de256715aee9ca0a68bccb93.rzhoumyrbmsgczpavmvbgubdqumsvqjq.validation.incaptest.co\",\n                    \"allowedDomainsForCNAMEValidation\": [],\n                    \"allowCNAMEValidation\": false\n                },\n                \"useWildCardSanInsteadOfFQDN\": true,\n                \"addNakedDomainSanForWWWSites\": true\n            }\n        }\n    ]\n}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	val := false
	imp := ImpervaCertificate{
		AddNakedDomainSanForWWWSites: &val,
	}
	dto := AccountSSLSettingsDTO{
		ImpervaCertificate: &imp,
	}
	accountSSLSettingsResponse, diag := client.UpdateAccountSSLSettings(&dto, "")
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if accountSSLSettingsResponse.Errors != nil || accountSSLSettingsResponse.Data == nil || !*accountSSLSettingsResponse.Data[0].ImpervaCertificate.UseWildCardSanInsteadOfFQDN {
		t.Errorf("Got unexpected response")
	}
}

func TestClientGetAccountSSlSettingsDataInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n    \"data\": [\n        {\n            \"impervaCertificate\": {\n                \"delegation\": {\n                    \"valueForCNAMEValidation\": \"_f508e2d5de256715aee9ca0a68bccb93.rzhoumyrbmsgczpavmvbgubdqumsvqjq.validation.incaptest.co\",\n                    \"allowedDomainsForCNAMEValidation\": [],\n                    \"allowCNAMEValidation\": false\n                },\n                \"useWildCardSanInsteadOfFQDN\": true,\n                \"addNakedDomainSanForWWWSites\": true\n            }\n        }\n    ]\n}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountSSLSettingsResponse, diag := client.GetAccountSSLSettings("")
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if accountSSLSettingsResponse.Errors != nil || accountSSLSettingsResponse.Data == nil || !*accountSSLSettingsResponse.Data[0].ImpervaCertificate.UseWildCardSanInsteadOfFQDN {
		t.Errorf("Got unexpected response")
	}
}

func TestClientGetAccountSSlSettingsErrorsInBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n    \"errors\": [\n        {\n            \"status\": 400,\n            \"id\": \"ad7e0ec7b8b45cc7\",\n            \"source\": {\n                \"pointer\": \"/v3/account/ssl-settings\"\n            },\n            \"title\": \"Bad Request\",\n            \"detail\": \"JSON parse error: Unexpected end-of-input: expected close marker for Object \"\n        }\n    ]\n}"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	accountSSLSettingsResponse, diag := client.GetAccountSSLSettings("")
	if diag != nil {
		t.Errorf("Should not received an error")
	}
	if accountSSLSettingsResponse.Errors == nil || accountSSLSettingsResponse.Errors[0].Status != 400 || accountSSLSettingsResponse.Data != nil {
		t.Errorf("Got unexpected response")
	}
}

func TestClientGetAccountSSlSettingsErrorFromMY(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	_, diag := client.GetAccountSSLSettings("")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500, error") {
		t.Errorf("Should have received an error")
	}
}

func TestClientGetAccountSSlSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	updateAccountSSLSettingsResponse, diag := client.GetAccountSSLSettings("")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
	if updateAccountSSLSettingsResponse != nil {
		t.Errorf("Should have received a nil addAccountResponse instance")
	}
}

func TestClientDeleteAccountSSlSettingsBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "http://badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	diag := client.DeleteAccountSSLSettings("")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "Timeout exceeded while awaiting") {
		t.Errorf("Should have received an time out error")
	}
}

func TestClientDeleteAccountSSlSettingsErrorFromMY(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(500)
		rw.Write([]byte("error"))
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	diag := client.DeleteAccountSSLSettings("")
	if diag == nil || !diag.HasError() || !strings.Contains(diag[0].Detail, "got response status 500") {
		t.Errorf("Should have received an error")
	}
}

func TestClientDeleteAccountSSlSettings200FromMY(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", accountSSLSettingsUrl) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", accountSSLSettingsUrl, req.URL.String())
		}
		rw.WriteHeader(200)
	}))
	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	diag := client.DeleteAccountSSLSettings("")
	if diag != nil || diag.HasError() {
		t.Errorf("Should not received an error")
	}
}
