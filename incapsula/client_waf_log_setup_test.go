package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

//////////////////////////////////////////////////////////////
/// 	AddWAFLogSetup Tests
//////////////////////////////////////////////////////////////

func TestClientWAFLogS3SetupBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	wafLogSetupResponse, err := client.AddWAFLogSetupS3(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error activating WAF")) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientWAFLogSFTPSetupBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	wafLogSetupResponse, err := client.AddWAFLogSetupSFTP(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error activating WAF")) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientCreateWAFLogSetupS3BadJSONActivate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupS3(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing activate WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientCreateWAFLogSetupS3BadJSONEnable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else {
			rw.Write([]byte(`{`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupS3(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing change WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientCreateWAFLogSetupS3BadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupS3(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing create S3 WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientCreateWAFLogSetupSFTPBadJSONActivate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupSFTP(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing activate WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientCreateWAFLogSetupSFTPBadJSONEnable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else {
			rw.Write([]byte(`{`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupSFTP(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing change WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientCreateWAFLogSetupSFTPBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupSFTP(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing create SFTP WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientAddWAFLogSetupS3Invalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupS3(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientAddWAFLogSetupS3Valid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupS3(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if wafLogSetupResponse == nil {
		t.Errorf("Should have received a wafLogSetupResponse instance")
	}
}

func TestClientAddWAFLogSetupSFTPInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupSFTP(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientAddWAFLogSetupSFTPValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	wafLogSetupResponse, err := client.AddWAFLogSetupSFTP(&WAFLogSetupPayload{0, true, "", "", "", "", "", "", ""})
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if wafLogSetupResponse == nil {
		t.Errorf("Should have received a wafLogSetupResponse instance")
	}
}

////////////////////////////////////////////////////////////////
/// 	DeleteSubAccount Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteWAFLogSetupBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123

	wafLogSetupResponse, err := client.RestoreWAFLogSetupDefault(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error restoring WAF Log Setup")) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientDeleteWAFLogSetupBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	wafLogSetupResponse, err := client.RestoreWAFLogSetupDefault(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing create default WAF")) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientDeleteWAFLogSetupInvalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	wafLogSetupResponse, err := client.RestoreWAFLogSetupDefault(accountID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if wafLogSetupResponse != nil {
		t.Errorf("Should have received a nil wafLogSetupResponse instance")
	}
}

func TestClientDeleteWAFLogSetupValid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsActivate) {
			rw.Write([]byte(`{"res":0,"res_message":"OK", "logs_collector_config_id": 123}`))
		} else if req.URL.String() == fmt.Sprintf("/%s", endpointWAFLogsChangeStatus) {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		} else {
			rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountID := 123
	wafLogSetupResponse, err := client.RestoreWAFLogSetupDefault(accountID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if wafLogSetupResponse == nil {
		t.Errorf("Should have received a wafLogSetupResponse instance")
	}
}
