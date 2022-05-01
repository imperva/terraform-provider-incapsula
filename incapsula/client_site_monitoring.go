package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MonitoringParameters struct {
	FailedRequestsPercentage    int    `json:"failedRequestsPercentage"`
	FailedRequestsMinNumber     int    `json:"failedRequestsMinNumber"`
	FailedRequestsDuration      int    `json:"failedRequestsDuration"`
	FailedRequestsDurationUnits string `json:"failedRequestsDurationUnits"`
}

type FailedRequestCriteria struct {
	HttpRequestTimeout      int    `json:"httpRequestTimeout"`
	HttpRequestTimeoutUnits string `json:"httpRequestTimeoutUnits"`
	HttpResponseError       string `json:"httpResponseError"`
}

type UpDownVerification struct {
	UseVerificationForDown bool   `json:"useVerificationForDown"`
	MonitoringUrl          string `json:"monitoringUrl"`
	ExpectedReceivedString string `json:"expectedReceivedString"`
	UpChecksInterval       int    `json:"upChecksInterval"`
	UpChecksIntervalUnits  string `json:"upChecksIntervalUnits"`
	UpCheckRetries         int    `json:"upCheckRetries"`
}

type Notifications struct {
	AlarmOnStandsByFailover bool   `json:"alarmOnStandsByFailover"`
	AlarmOnDcFailover       bool   `json:"alarmOnDcFailover"`
	AlarmOnServerFailover   bool   `json:"alarmOnServerFailover"`
	RequiredMonitors        string `json:"requiredMonitors"`
}

type SiteMonitoring struct {
	MonitoringParameters  MonitoringParameters  `json:"monitoringParameters"`
	FailedRequestCriteria FailedRequestCriteria `json:"failedRequestCriteria"`
	UpDownVerification    UpDownVerification    `json:"upDownVerification"`
	Notifications         Notifications         `json:"notifications"`
}

type SiteMonitoringResponse struct {
	Data []SiteMonitoring `json:"data"`
}

func (c *Client) GetSiteMonitoring(siteID int) (*SiteMonitoringResponse, error) {
	log.Printf("[INFO] Getting Incapsula Site Monitoring for Site ID %d", siteID)
	return Crud("Read", siteID, http.MethodGet, nil, c)
}

func (c *Client) UpdateSiteMonitoring(siteID int, siteMonitoring *SiteMonitoring) (*SiteMonitoringResponse, error) {
	log.Printf("[INFO] Updating Incapsula Site Monitoring for Site ID %d", siteID)
	siteMopnitoringJSON, err := json.Marshal(siteMonitoring)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Site Monitoring for SiteID: %s", err)
	}
	return Crud("Update", siteID, http.MethodPost, siteMopnitoringJSON, c)
}

func Crud(action string, siteID int, hhtpMethod string, data []byte, c *Client) (*SiteMonitoringResponse, error) {
	url := fmt.Sprintf("%s/appdlv-site-settings/v2/site/%d/monitoring", c.config.BaseURLAPI, siteID)

	resp, err := c.DoJsonRequestWithHeaders(hhtpMethod, url, data, "")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when %s Site Monitoring for Site ID %d: %s", strings.ToLower(action)+"ing", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula %s Site Monitoring JSON response: %s\n", action, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when %s Site Monitoring for Site ID %d: %s", resp.StatusCode, strings.ToLower(action)+"ing", siteID, string(responseBody))
	}

	// Dump JSON
	var siteMonitoringResponse SiteMonitoringResponse
	err = json.Unmarshal([]byte(responseBody), &siteMonitoringResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Site Monitoring Response JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}
	return &siteMonitoringResponse, nil
}
