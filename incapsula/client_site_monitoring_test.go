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
// UpdateSiteMonitoring Tests
////////////////////////////////////////////////////////////////
func TestUpdateSiteMonitoringBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	payload := SiteMonitoring{
		MonitoringParameters:  MonitoringParameters{},
		FailedRequestCriteria: FailedRequestCriteria{},
		UpDownVerification:    UpDownVerification{},
		Notifications:         Notifications{},
	}

	siteMonitoringResponse, err := client.UpdateSiteMonitoring(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error from Incapsula service when updateing Site Monitoring for Site ID") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteMonitoringResponse != nil {
		t.Errorf("Should have received a nil siteMonitoringResponse instance")
	}
}

func TestUpdateSiteMonitoringBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/appdlv-site-settings/v2/site/%d/monitoring", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	payload := SiteMonitoring{
		MonitoringParameters:  MonitoringParameters{},
		FailedRequestCriteria: FailedRequestCriteria{},
		UpDownVerification:    UpDownVerification{},
		Notifications:         Notifications{},
	}

	siteMonitoringResponse, err := client.UpdateSiteMonitoring(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error parsing Site Monitoring Response JSON response for Site ID") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteMonitoringResponse != nil {
		t.Errorf("Should have received a nil siteMonitoringResponse instance")
	}
}

func TestUpdateSiteMonitoringInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/appdlv-site-settings/v2/site/%d/monitoring", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "errors": [
    {
      "id": null,
      "status": 404,
      "source": {
        "pointer": "/appdlv-site-settings/v2/site/42/monitoring"
      },
      "title": "Not Found"
    }
  ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	//invalid payload
	payload := SiteMonitoring{
		//MonitoringParameters: MonitoringParameters{},
		//FailedRequestCriteria: FailedRequestCriteria{},
		//UpDownVerification: UpDownVerification{},
		//Notifications: Notifications{},
	}

	siteMonitoringResponse, err := client.UpdateSiteMonitoring(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when updateing Site Monitoring for Site ID")) {
		t.Errorf("Should have received a bad Site Monitoring error, got: %s", err)
	}
	if siteMonitoringResponse != nil {
		t.Errorf("Should have received a nil siteMonitoringResponse instance")
	}
}

func TestUpdateSiteMonitoringConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/appdlv-site-settings/v2/site/%d/monitoring", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "data": [
    {
      "monitoringParameters": {
        "failedRequestsPercentage": 30,
        "failedRequestsMinNumber": 3,
        "failedRequestsDuration": 30,
        "failedRequestsDurationUnits": "SECONDS"
      },
      "failedRequestCriteria": {
        "httpRequestTimeout": 35,
        "httpRequestTimeoutUnits": "SECONDS",
        "httpResponseError": "501-510,530"
      },
      "upDownVerification": {
        "useVerificationForDown": false,
        "monitoringUrl": "/health",
        "expectedReceivedString": "Am Alive",
        "upChecksInterval": 20,
        "upChecksIntervalUnits": "SECONDS",
        "upCheckRetries": 30
      },
      "notifications": {
        "alarmOnStandsByFailover": false,
        "alarmOnDcFailover": false,
        "alarmOnServerFailover": true,
        "requiredMonitors": "MORE"
      }
    }
  ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	payload := SiteMonitoring{
		MonitoringParameters:  MonitoringParameters{FailedRequestsMinNumber: 3},
		FailedRequestCriteria: FailedRequestCriteria{HttpRequestTimeout: 20},
		UpDownVerification:    UpDownVerification{},
		Notifications:         Notifications{},
	}

	siteMonitoringResponse, err := client.UpdateSiteMonitoring(
		siteID,
		&payload)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if siteMonitoringResponse == nil {
		t.Errorf("Should not have received a nil apiSecuritySiteConfigPostResponse instance")
	}
	if siteMonitoringResponse.Data[0].MonitoringParameters.FailedRequestsMinNumber != 3 {
		t.Errorf("Should have received a FailedRequestsMinNumber equal 3")
	}
}

////////////////////////////////////////////////////////////////
// ReadSiteMonitoring Tests
////////////////////////////////////////////////////////////////
func TestReadSiteMonitoringBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	siteMonitoringResponse, err := client.GetSiteMonitoring(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error from Incapsula service when reading Site Monitoring for Site ID") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteMonitoringResponse != nil {
		t.Errorf("Should have received a nil siteMonitoringResponse instance")
	}
}

func TestUpdatReadSiteMonitoringBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/appdlv-site-settings/v2/site/%d/monitoring", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteMonitoringResponse, err := client.GetSiteMonitoring(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error parsing Site Monitoring Response JSON response for Site ID") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if siteMonitoringResponse != nil {
		t.Errorf("Should have received a nil siteMonitoringResponse instance")
	}

}

func TestReadSiteMonitoringInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/appdlv-site-settings/v2/site/%d/monitoring", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(404)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "errors": [
    {
      "id": null,
      "status": 404,
      "source": {
        "pointer": "/appdlv-site-settings/v2/site/42/monitoring"
      },
      "title": "Not Found"
    }
  ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteMonitoringResponse, err := client.GetSiteMonitoring(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 404 from Incapsula service when reading Site Monitoring for Site ID")) {
		t.Errorf("Should have received a bad Site Monitoring error, got: %s", err)
	}
	if siteMonitoringResponse != nil {
		t.Errorf("Should have received a nil siteMonitoringResponse instance")
	}
}

func TestReadSiteMonitoringConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("/appdlv-site-settings/v2/site/%d/monitoring", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
  "data": [
    {
      "monitoringParameters": {
        "failedRequestsPercentage": 30,
        "failedRequestsMinNumber": 3,
        "failedRequestsDuration": 30,
        "failedRequestsDurationUnits": "SECONDS"
      },
      "failedRequestCriteria": {
        "httpRequestTimeout": 35,
        "httpRequestTimeoutUnits": "SECONDS",
        "httpResponseError": "501-510,530"
      },
      "upDownVerification": {
        "useVerificationForDown": false,
        "monitoringUrl": "/health",
        "expectedReceivedString": "Am Alive",
        "upChecksInterval": 20,
        "upChecksIntervalUnits": "SECONDS",
        "upCheckRetries": 30
      },
      "notifications": {
        "alarmOnStandsByFailover": false,
        "alarmOnDcFailover": false,
        "alarmOnServerFailover": true,
        "requiredMonitors": "MORE"
      }
    }
  ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	siteMonitoringResponse, err := client.GetSiteMonitoring(siteID)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if siteMonitoringResponse == nil {
		t.Errorf("Should not have received a nil apiSecuritySiteConfigPostResponse instance")
	}
	if siteMonitoringResponse.Data[0].MonitoringParameters.FailedRequestsMinNumber != 3 {
		t.Errorf("Should have received a FailedRequestsMinNumber equal 3")
	}
}
