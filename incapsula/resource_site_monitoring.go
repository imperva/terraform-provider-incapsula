package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strconv"
	"strings"
)

const defaultFailedRequestsMinNumber = 3
const defaultFailedRequestsPercentage = 40
const defaultfailedRequestDuration = 40
const defaultFailedRequestsDurationUnits = "SECONDS"

const defaultHttpRequestTimeout = 35
const defaultHttpRequestTimeoutUnits = "SECONDS"
const defaultHttpResponseError = "501-599"

const defaultUseVerificationForDown = true
const defaultMonitoringUrl = "/"
const defaultUpChecksInterval = 20
const defaultUpChecksIntervalUnits = "SECONDS"
const defaultUpCheckRetries = 3

const defaultAlarmOnStandsByFailover = true
const defaultAlarmOnDcFailover = true
const defaultAlarmOnServerFailover = false
const defaultRequiredMonitors = "MOST"

func resourceSiteMonitoring() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteMonitoringUpdate,
		Read:   resourceSiteMonitoringRead,
		Update: resourceSiteMonitoringUpdate,
		Delete: resourceSiteMonitoringDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID, err := strconv.Atoi(d.Id())
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", d.Id())
				}

				d.Set("site_id", siteID)
				log.Printf("[DEBUG] Import Site Config JSON for Site ID %d", siteID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"failed_requests_percentage": {
				Type:         schema.TypeInt,
				Description:  "The percentage of failed requests to the origin server",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 100),
				Default:      defaultFailedRequestsPercentage,
			},
			"failed_requests_min_number": {
				Type:         schema.TypeInt,
				Description:  "The minimum number of of failed requests to be considered as failure",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 500),
				Default:      defaultFailedRequestsMinNumber,
			},
			"failed_requests_duration": {
				Type:         schema.TypeInt,
				Description:  "The minimum duration of failures above the threshold to consider server as down. 20-180 SECONDS or 1-2 MINUTES. Default: 40",
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      defaultfailedRequestDuration,
			},
			"failed_requests_duration_units": {
				Type:         schema.TypeString,
				Description:  "Time unit. Possible values: SECONDS, MINUTES. Default: SECONDS",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SECONDS", "MINUTES"}, false),
				Default:      defaultFailedRequestsDurationUnits,
			},
			"http_request_timeout": {
				Type:         schema.TypeInt,
				Description:  "The maximum time to wait for an HTTP response. 1-200 SECONDS or 1-2 MINUTES",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 200),
				Default:      defaultHttpRequestTimeout,
			},
			"http_request_timeout_units": {
				Type:         schema.TypeString,
				Description:  "Time unit",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SECONDS", "MINUTES"}, false),
				Default:      defaultHttpRequestTimeoutUnits,
			},
			"http_response_error": {
				Type:        schema.TypeString,
				Description: "The HTTP response error codes or patterns that will be counted as request failures",
				Optional:    true,
				Default:     defaultHttpResponseError,
			},
			"use_verification_for_down": {
				Type:        schema.TypeBool,
				Description: "If Imperva determines that an origin server is down according to failed request criteria, it will initiate another request to verify that the origin server is down", //todo ??????
				Optional:    true,
				Default:     defaultUseVerificationForDown,
			},
			"monitoring_url": {
				Type:        schema.TypeString,
				Description: "The URL to use for monitoring your website.",
				Optional:    true,
				Default:     defaultMonitoringUrl,
			},
			"expected_received_string": {
				Type:        schema.TypeString,
				Description: "The expected string. If left empty, any response, except for the codes defined in the HTTP response error codes to be treated as Down parameter, will be considered successful. If the value is non-empty, then the defined value must appear within the response string for the response to be considered successful.",
				Optional:    true,
			},
			"up_checks_interval": {
				Type:         schema.TypeInt,
				Description:  "After an origin server was identified as down, Imperva will periodically test it to see whether it has recovered, according to the frequency defined in this parameter. 10-120 SECONDS or 1-2 MINUTES",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 200),
				Default:      defaultUpChecksInterval,
			},
			"up_checks_interval_units": {
				Type:         schema.TypeString,
				Description:  "Time unit. Default: SECONDS",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"SECONDS", "MINUTES"}, false),
				Default:      defaultUpChecksIntervalUnits,
			},
			"up_check_retries": {
				Type:         schema.TypeInt,
				Description:  "Every time an origin server is tested to see whether it’s back up, the test will be retried this number of times.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 50),
				Default:      defaultUpCheckRetries,
			},
			"required_monitors": {
				Type:         schema.TypeString,
				Description:  "Monitors required to report server / data center as down",
				Optional:     true,
				Default:      defaultRequiredMonitors,
				ValidateFunc: validation.StringInSlice([]string{"ONE", "MANY", "MOST", "ALL"}, false),
			},
			"alarm_on_stands_by_failover": {
				Type:        schema.TypeBool,
				Description: "Indicates whether or not an email will be sent upon failover to a standby data center",
				Optional:    true,
				Default:     defaultAlarmOnStandsByFailover,
			},
			"alarm_on_server_failover": {
				Type:        schema.TypeBool,
				Description: "Indicates whether or not an email will be sent upon data center failover",
				Optional:    true,
				Default:     defaultAlarmOnServerFailover,
			},
			"alarm_on_dc_failover": {
				Type:        schema.TypeBool,
				Description: "Indicates whether or not an email will be sent upon server failover",
				Optional:    true,
				Default:     defaultAlarmOnDcFailover,
			},
			//todo add description  MOST - More than 50%, ask backend for descr

		},
	}
}

func resourceSiteMonitoringUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	siteIDStr := strconv.Itoa(siteID)

	siteMonitoring := SiteMonitoring{
		MonitoringParameters: MonitoringParameters{
			FailedRequestsPercentage:    d.Get("failed_requests_percentage").(int),
			FailedRequestsMinNumber:     d.Get("failed_requests_min_number").(int),
			FailedRequestsDuration:      d.Get("failed_requests_duration").(int),
			FailedRequestsDurationUnits: d.Get("failed_requests_duration_units").(string),
		},
		FailedRequestCriteria: FailedRequestCriteria{
			HttpRequestTimeout:      d.Get("http_request_timeout").(int),
			HttpRequestTimeoutUnits: d.Get("http_request_timeout_units").(string),
			HttpResponseError:       d.Get("http_response_error").(string),
		},
		UpDownVerification: UpDownVerification{
			UseVerificationForDown: d.Get("use_verification_for_down").(bool),
			MonitoringUrl:          d.Get("monitoring_url").(string),
			ExpectedReceivedString: d.Get("expected_received_string").(string),
			UpChecksInterval:       d.Get("up_checks_interval").(int),
			UpChecksIntervalUnits:  d.Get("up_checks_interval_units").(string),
			UpCheckRetries:         d.Get("up_check_retries").(int),
		},
		Notifications: Notifications{
			AlarmOnStandsByFailover: d.Get("alarm_on_stands_by_failover").(bool),
			AlarmOnDcFailover:       d.Get("alarm_on_dc_failover").(bool),
			AlarmOnServerFailover:   d.Get("alarm_on_server_failover").(bool),
			RequiredMonitors:        d.Get("required_monitors").(string),
		},
	}

	siteMonitoringResponse, err := client.UpdateSiteMonitoring(siteID, &siteMonitoring)
	if strings.Contains(fmt.Sprint(err), "Missing Load Balancing subscription for Site ID") {
		log.Printf("[ERROR] Could not get Incapsula Site Monitoring for Site Id: %d - %s\n. Missing Load Balancing subscription for Site ID. The resource will be removed.", siteID, err)
		d.SetId("")
		return err
	}

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula Site Monitoring for Site Id: %d - %s\n", siteID, err)
		return err
	}

	siteMonitoringResult := siteMonitoringResponse.Data[0]

	d.SetId(siteIDStr)
	d.Set("failed_requests_percentage", siteMonitoringResult.MonitoringParameters.FailedRequestsPercentage)
	d.Set("failed_requests_min_number", siteMonitoringResult.MonitoringParameters.FailedRequestsMinNumber)
	d.Set("failed_requests_duration", siteMonitoringResult.MonitoringParameters.FailedRequestsDuration)
	d.Set("failed_requests_duration_units", siteMonitoringResult.MonitoringParameters.FailedRequestsDurationUnits)

	d.Set("http_request_timeout", siteMonitoringResult.FailedRequestCriteria.HttpRequestTimeout)
	d.Set("http_request_timeout_units", siteMonitoringResult.FailedRequestCriteria.HttpRequestTimeoutUnits)
	d.Set("http_response_error", siteMonitoringResult.FailedRequestCriteria.HttpResponseError)

	d.Set("use_verification_for_down", siteMonitoringResult.UpDownVerification.UseVerificationForDown)
	d.Set("monitoring_url", siteMonitoringResult.UpDownVerification.MonitoringUrl)
	d.Set("expected_received_string", siteMonitoringResult.UpDownVerification.ExpectedReceivedString)
	d.Set("up_checks_interval", siteMonitoringResult.UpDownVerification.UpChecksInterval)
	d.Set("up_checks_interval_units", siteMonitoringResult.UpDownVerification.UpChecksIntervalUnits)
	d.Set("up_check_retries", siteMonitoringResult.UpDownVerification.UpCheckRetries)

	d.Set("alarm_on_stands_by_failover", siteMonitoringResult.Notifications.AlarmOnStandsByFailover)
	d.Set("alarm_on_server_failover", siteMonitoringResult.Notifications.AlarmOnServerFailover)
	d.Set("alarm_on_dc_failover", siteMonitoringResult.Notifications.AlarmOnDcFailover)
	d.Set("required_monitors", siteMonitoringResult.Notifications.RequiredMonitors)

	return nil
}

func resourceSiteMonitoringRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	siteIdStr := strconv.Itoa(siteID)

	siteMonitoringResponse, err := client.GetSiteMonitoring(siteID)
	if strings.Contains(fmt.Sprint(err), "Missing Load Balancing subscription for Site ID") {
		log.Printf("[ERROR] Could not get Incapsula Site Monitoring for Site Id: %d - %s\n. Missing Load Balancing subscription for Site ID. The resource will be removed.", siteID, err)
		d.SetId("")
		return err
	}
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula Site Monitoring for Site Id: %d - %s\n", siteID, err)
		return err
	}

	d.SetId(siteIdStr)

	siteMonitoring := siteMonitoringResponse.Data[0]

	d.Set("failed_requests_percentage", siteMonitoring.MonitoringParameters.FailedRequestsPercentage)
	d.Set("failed_requests_min_number", siteMonitoring.MonitoringParameters.FailedRequestsMinNumber)
	d.Set("failed_requests_duration", siteMonitoring.MonitoringParameters.FailedRequestsDuration)
	d.Set("failed_requests_duration_units", siteMonitoring.MonitoringParameters.FailedRequestsDurationUnits)

	d.Set("http_request_timeout", siteMonitoring.FailedRequestCriteria.HttpRequestTimeout)
	d.Set("http_request_timeout_units", siteMonitoring.FailedRequestCriteria.HttpRequestTimeoutUnits)
	d.Set("http_response_error", siteMonitoring.FailedRequestCriteria.HttpResponseError)

	d.Set("use_verification_for_down", siteMonitoring.UpDownVerification.UseVerificationForDown)
	d.Set("monitoring_url", siteMonitoring.UpDownVerification.MonitoringUrl)
	d.Set("expected_received_string", siteMonitoring.UpDownVerification.ExpectedReceivedString)
	d.Set("up_checks_interval", siteMonitoring.UpDownVerification.UpChecksInterval)
	d.Set("up_checks_interval_units", siteMonitoring.UpDownVerification.UpChecksIntervalUnits)
	d.Set("up_check_retries", siteMonitoring.UpDownVerification.UpCheckRetries)

	d.Set("alarm_on_stands_by_failover", siteMonitoring.Notifications.AlarmOnStandsByFailover)
	d.Set("alarm_on_server_failover", siteMonitoring.Notifications.AlarmOnServerFailover)
	d.Set("alarm_on_dc_failover", siteMonitoring.Notifications.AlarmOnDcFailover)
	d.Set("required_monitors", siteMonitoring.Notifications.RequiredMonitors)

	return nil
}

func resourceSiteMonitoringDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")
	return nil
}
