package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Compression struct {
	FileCompression  bool `json:"file_compression"`
	MinifyJs         bool `json:"minify_js"`
	MinifyCss        bool `json:"minify_css"`
	MinifyStaticHtml bool `json:"minify_static_html"`
}
type CompressionStr struct {
	FileCompression  string `json:"file_compression"`
	MinifyJs         string `json:"minify_js"`
	MinifyCss        string `json:"minify_css"`
	MinifyStaticHtml string `json:"minify_static_html"`
}

type ImageCompression struct {
	CompressJpeg              bool `json:"compress_jpeg"`
	ProgressiveImageRendering bool `json:"progressive_image_rendering"`
	AggressiveCompression     bool `json:"aggressive_compression"`
	CompressPng               bool `json:"compress_png"`
}

type ImageCompressionStr struct {
	CompressJpeg              string
	ProgressiveImageRendering string
	AggressiveCompression     string
	CompressPng               string
}

type Network struct {
	TcpPrePooling         bool    `json:"tcp_pre_pooling"`
	OriginConnectionReuse bool    `json:"origin_connection_reuse"`
	SupportNonSniClients  bool    `json:"support_non_sni_clients"`
	EnableHttp2           bool    `json:"enable_http2"`
	Http2ToOrigin         bool    `json:"http2_to_origin"`
	Port                  Port    `json:"port"`
	SslPort               SslPort `json:"ssl_port"`
}

type NetworkStr struct {
	TcpPrePooling         string  `json:"tcp_pre_pooling"`
	OriginConnectionReuse string  `json:"origin_connection_reuse"`
	SupportNonSniClients  string  `json:"support_non_sni_clients"`
	EnableHttp2           string  `json:"enable_http2"`
	Http2ToOrigin         string  `json:"http2_to_origin"`
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

type RedirectionStr struct {
	RedirectNakedToFullStr string `json:"redirect_naked_to_full"`
	RedirectHttpToHttpsStr string `json:"redirect_http_to_https"`
}

type CustomErrorPageTemplates struct {
	ErrorConnectionTimeout string `json:"error.type.connection_timeout"`
	ErrorAccessDenied      string `json:"error.type.access_denied"`
	ErrorParseReqError     string `json:"error.type.parse_req_error"`
	ErrorParseRespError    string `json:"error.type.parse_resp_error"`
	ErrorConnectionFailed  string `json:"error.type.connection_failed"`
	ErrorSslFailed         string `json:"error.type.ssl_failed"`
	ErrorDenyAndCaptcha    string `json:"error.type.deny_and_captcha"`
	ErrorTypeNoSslConfig   string `json:"error.type.no_ssl_config"`
}

type CustomErrorPage struct {
	DefaultErrorPage         string                   `json:"error_page_template"`
	CustomErrorPageTemplates CustomErrorPageTemplates `json:"custom_error_page_templates"`
}

type ApplicationDelivery struct {
	Compression      Compression      `json:"compression"`
	ImageCompression ImageCompression `json:"image_compression"`
	Network          Network          `json:"network"`
	Redirection      Redirection      `json:"redirection"`
	CustomErrorPage  CustomErrorPage  `json:"custom_error_page"`
}

func (c *Client) GetApplicationDelivery(siteID int) (*ApplicationDelivery, error) {
	log.Printf("[INFO] Getting Incapsula Application Delivery for Site ID %d", siteID)
	return CrudApplicationDelivery("Read", siteID, http.MethodGet, nil, c)
}

func (c *Client) UpdateApplicationDelivery(siteID int, applicationDelivery *ApplicationDelivery) (*ApplicationDelivery, error) {
	log.Printf("[INFO] Updating Incapsula Application Delivery for Site ID %d", siteID)
	applicationDeliveryJSON, err := json.Marshal(applicationDelivery)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Application Delivery for SiteID: %s", err)
	}
	return CrudApplicationDelivery("Update", siteID, http.MethodPut, applicationDeliveryJSON, c)
}

func (c *Client) DeleteApplicationDelivery(siteID int) (*ApplicationDelivery, error) {
	log.Printf("[INFO] Deleting Incapsula Application Delivery for Site ID %d", siteID)
	return CrudApplicationDelivery("Delete", siteID, http.MethodDelete, nil, c)
}

func CrudApplicationDelivery(action string, siteID int, hhtpMethod string, data []byte, c *Client) (*ApplicationDelivery, error) {
	url := fmt.Sprintf("%s/sites/%d/settings/delivery", c.config.BaseURLRev2, siteID)
	//todo tolowerCase for operation name
	resp, err := c.DoJsonRequestWithHeaders(hhtpMethod, url, data, strings.ToLower(action)+"_application_delivery")
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when %s Application Delivery for Site ID %d: %s", strings.ToLower(action)+"ing", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula %s Application Delivery JSON response: %s\n", action, string(responseBody))

	////todo ????????? leave it? what subscription do we need?
	//if resp.StatusCode == 404 {
	//	return nil, fmt.Errorf("Missing Load Balancing subscription for Site ID %d: %s", siteID, string(responseBody))
	//}

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when %s Application Delivery for Site ID %d: %s", resp.StatusCode, strings.TrimSuffix(action, "e")+"ing", siteID, string(responseBody))
	}

	// Dump JSON
	var applicationDelivery ApplicationDelivery
	err = json.Unmarshal([]byte(responseBody), &applicationDelivery)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing Application Delivery Response JSON response for Site ID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}
	//todo check if data.length is >0
	return &applicationDelivery, nil
}
