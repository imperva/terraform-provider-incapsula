package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Compression struct {
	FileCompression  bool   `json:"file_compression"`
	CompressionType  string `json:"compression_type"`
	MinifyJs         bool   `json:"minify_js"`
	MinifyCss        bool   `json:"minify_css"`
	MinifyStaticHtml bool   `json:"minify_static_html"`
}

type ImageCompression struct {
	CompressJpeg              bool `json:"compress_jpeg"`
	ProgressiveImageRendering bool `json:"progressive_image_rendering"`
	AggressiveCompression     bool `json:"aggressive_compression"`
	CompressPng               bool `json:"compress_png"`
}

type Network struct {
	TcpPrePooling         bool    `json:"tcp_pre_pooling"`
	OriginConnectionReuse bool    `json:"origin_connection_reuse"`
	SupportNonSniClients  bool    `json:"support_non_sni_clients"`
	EnableHttp2           *bool   `json:"enable_http2"`
	Http2ToOrigin         *bool   `json:"http2_to_origin"`
	Port                  Port    `json:"port"`
	SslPort               SslPort `json:"ssl_port"`
}

type Port struct {
	To string `json:"to"`
}

type SslPort struct {
	To string `json:"to"`
}

type Redirection struct {
	RedirectNakedToFull bool `json:"redirect_naked_to_full"`
	RedirectHttpToHttps bool `json:"redirect_http_to_https"`
}

type CustomErrorPageTemplates struct {
	ErrorConnectionTimeout       string `json:"error.type.connection_timeout,omitempty"`
	ErrorAccessDenied            string `json:"error.type.access_denied,omitempty"`
	ErrorParseReqError           string `json:"error.type.parse_req_error,omitempty"`
	ErrorParseRespError          string `json:"error.type.parse_resp_error,omitempty"`
	ErrorConnectionFailed        string `json:"error.type.connection_failed,omitempty"`
	ErrorSslFailed               string `json:"error.type.ssl_failed,omitempty"`
	ErrorDenyAndCaptcha          string `json:"error.type.deny_and_captcha,omitempty"`
	ErrorTypeNoSslConfig         string `json:"error.type.no_ssl_config,omitempty"`
	ErrorAbpIdentificationFailed string `json:"error.type.abp_identification_failed,omitempty"`
}

type CustomErrorPage struct {
	DefaultErrorPage         string                   `json:"error_page_template,omitempty"`
	CustomErrorPageTemplates CustomErrorPageTemplates `json:"custom_error_page_templates"`
}

type ApplicationDelivery struct {
	Compression      Compression      `json:"compression"`
	ImageCompression ImageCompression `json:"image_compression"`
	Network          Network          `json:"network"`
	Redirection      Redirection      `json:"redirection"`
}

func (c *Client) GetApplicationDelivery(siteID int) (*ApplicationDelivery, diag.Diagnostics) {
	log.Printf("[INFO] Getting Incapsula Application Delivery for Site ID %d", siteID)
	return CrudApplicationDelivery("Read", siteID, http.MethodGet, nil, c)
}

func (c *Client) UpdateApplicationDelivery(siteID int, applicationDelivery *ApplicationDelivery) (*ApplicationDelivery, diag.Diagnostics) {
	log.Printf("[INFO] Updating Incapsula Application Delivery for Site ID %d", siteID)
	var diags diag.Diagnostics
	applicationDeliveryJSON, err := json.Marshal(applicationDelivery)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed JSON marshalling Application Delivery resource",
			Detail:   fmt.Sprintf("Failed to JSON marshal Application Delivery for SiteID %d: %s", siteID, err),
		})
		return nil, diags
	}
	return CrudApplicationDelivery("Update", siteID, http.MethodPut, applicationDeliveryJSON, c)
}

func (c *Client) DeleteApplicationDelivery(siteID int) (*ApplicationDelivery, diag.Diagnostics) {
	log.Printf("[INFO] Deleting Incapsula Application Delivery for Site ID %d", siteID)
	return CrudApplicationDelivery("Delete", siteID, http.MethodDelete, nil, c)
}

func CrudApplicationDelivery(action string, siteID int, httpMethod string, applicationDeliveyData []byte, c *Client) (*ApplicationDelivery, diag.Diagnostics) {
	var diags diag.Diagnostics

	if applicationDeliveyData != nil {
		log.Printf("[DEBUG] Incapsula %s Application Delivery JSON request: %s\n", action, string(applicationDeliveyData))
	}

	applicationDeliveryUrl := fmt.Sprintf("%s/sites/%d/settings/delivery", c.config.BaseURLRev2, siteID)

	resp, err := c.DoJsonRequestWithHeaders(httpMethod, applicationDeliveryUrl, applicationDeliveyData, strings.ToLower(action)+"_application_delivery")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error from Incapsula service when trying to %s Application Delivery", strings.ToLower(action)),
			Detail:   fmt.Sprintf("Error from Incapsula service when trying to %s Application Delivery for Site ID %d: %s", strings.ToLower(action), siteID, err),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula %s Application Delivery JSON response: (%d) %s\n", action, resp.StatusCode, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error status code %d from Incapsula service when %s Application Delivery", resp.StatusCode, strings.TrimSuffix(action, "e")+"ing"),
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when %s Application Delivery for Site ID %d: %s", resp.StatusCode, strings.TrimSuffix(action, "e")+"ing", siteID, string(responseBody)),
		})
		return nil, diags
	}

	// Dump JSON
	var applicationDelivery ApplicationDelivery
	err = json.Unmarshal([]byte(responseBody), &applicationDelivery)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error parsing Application Delivery Response JSON response"),
			Detail:   fmt.Sprintf("Error parsing Application Delivery Response JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody)),
		})
		return nil, diags
	}

	return &applicationDelivery, nil
}

func (c *Client) GetErrorPages(siteID int) (*CustomErrorPage, diag.Diagnostics) {
	log.Printf("[INFO] Getting Incapsula Error Pages for Site ID %d", siteID)
	return CrudErrorPages("Read", siteID, http.MethodGet, nil, c)
}

func (c *Client) UpdateErrorPages(siteID int, errorPages *CustomErrorPage) (*CustomErrorPage, diag.Diagnostics) {
	log.Printf("[INFO] Updating Incapsula Application Delivery for Site ID %d", siteID)
	var diags diag.Diagnostics
	errorPagesJSON, err := json.Marshal(errorPages)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed JSON marshalling Error Pages resource",
			Detail:   fmt.Sprintf("Failed to JSON marshal Error Pages for SiteID %d: %s", siteID, err),
		})
		return nil, diags
	}
	return CrudErrorPages("Update", siteID, http.MethodPut, errorPagesJSON, c)
}

func (c *Client) DeleteErrorPages(siteID int) (*CustomErrorPage, diag.Diagnostics) {
	log.Printf("[INFO] Deleting Incapsula Application Delivery for Site ID %d", siteID)
	customErrorPage := CustomErrorPage{}
	errorPagesJSON, _ := json.Marshal(customErrorPage)
	return CrudErrorPages("Delete", siteID, http.MethodPut, errorPagesJSON, c)
}

func CrudErrorPages(action string, siteID int, httpMethod string, errorPagesData []byte, c *Client) (*CustomErrorPage, diag.Diagnostics) {
	var diags diag.Diagnostics

	if errorPagesData != nil {
		log.Printf("[DEBUG] Incapsula %s Error Pages JSON request: %s\n", action, string(errorPagesData))
	}

	errorPagesUrl := fmt.Sprintf("%s/sites/%d/settings/delivery/error-pages", c.config.BaseURLRev2, siteID)

	resp, err := c.DoJsonRequestWithHeaders(httpMethod, errorPagesUrl, errorPagesData, strings.ToLower(action)+"_error_pages")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error from Incapsula service when trying to %s Error Pages", strings.ToLower(action)),
			Detail:   fmt.Sprintf("Error from Incapsula service when trying to %s Error Pages for Site ID %d: %s", strings.ToLower(action), siteID, err),
		})
		return nil, diags
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, _ := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula %s Error pages JSON response: (%d) %s\n", action, resp.StatusCode, string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error status code %d from Incapsula service when %s Error Pages", resp.StatusCode, strings.TrimSuffix(action, "e")+"ing"),
			Detail:   fmt.Sprintf("Error status code %d from Incapsula service when %s Error Pages for Site ID %d: %s", resp.StatusCode, strings.TrimSuffix(action, "e")+"ing", siteID, string(responseBody)),
		})
		return nil, diags
	}

	// Dump JSON
	var customErrorPage CustomErrorPage
	err = json.Unmarshal([]byte(responseBody), &customErrorPage)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Error parsing Error Pages Response JSON response for Site ID %d", siteID),
			Detail:   fmt.Sprintf("Error parsing Error Pages Response JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody)),
		})
		return nil, diags
	}

	return &customErrorPage, nil
}
