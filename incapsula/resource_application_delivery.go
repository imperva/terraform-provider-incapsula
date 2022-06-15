package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

const defaultPortTo = 80
const defaultSslPortTo = 443

func resourceApplicationDelivery() *schema.Resource {
	return &schema.Resource{
		Create: resourceApplicationDeliveryFlatUpdate,
		Read:   resourceApplicationDeliveryFlatRead,
		Update: resourceApplicationDeliveryFlatUpdate,
		Delete: resourceApplicationDeliveryFlatDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID, err := strconv.Atoi(d.Id())
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command for Application Delivery, actual value: %s, expected numeric id", d.Id())
				}

				d.Set("site_id", siteID)
				log.Printf("[DEBUG] Import Application Delivery JSON for Site ID %d", siteID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on. ",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},

			"file_compression": {
				Type:        schema.TypeBool,
				Description: "When this option is enabled, any textual resource, such as Javascript, CSS and HTML, is compressed using Gzip as it is being transferred, and then unzipped within the browser. All modern browsers support this feature.",
				Optional:    true,
				Default:     true,
			},
			"minify_js": {
				Type:        schema.TypeBool,
				Description: "Minify JavaScript. Minification removes characters that are not necessary for rendering the page, such as whitespace and comments. This makes the files smaller and therefore reduces their access time. Minification has no impact on the functionality of the Javascript, CSS, and HTML files.",
				Optional:    true,
				Default:     true,
			},
			"minify_css": {
				Type:        schema.TypeBool,
				Description: "Content minification can applied only to cached Javascript, CSS and HTML content.",
				Optional:    true,
				Default:     true,
			},
			"minify_static_html": {
				Type:        schema.TypeBool,
				Description: "Minify static HTML",
				Optional:    true,
				Default:     true,
			},

			"compress_jpeg": {
				Type:        schema.TypeBool,
				Description: "Compress JPEG images. Compression reduces download time by reducing the file size.",
				Optional:    true,
				Default:     true,
			},
			"progressive_image_rendering": {
				Type:        schema.TypeBool,
				Description: "The image is rendered with progressively finer resolution, potentially causing a pixelated effect until the final image is rendered with no loss of quality. This option reduces page load times and allows images to gradually load after the page is rendered.",
				Optional:    true,
				Default:     false,
			},
			"aggressive_compression": {
				Type:        schema.TypeBool,
				Description: "A more aggressive method of compression is applied with the goal of minimizing the image file size, possibly impacting the final quality of the image displayed. Applies to JPEG compression only.",
				Optional:    true,
				Default:     false,
				//DefaultFunc: func() (interface{}, error) {
				//	//return true
				//},
			},
			"compress_png": {
				Type:        schema.TypeBool,
				Description: "Compress PNG images. Compression reduces download time by reducing the file size. PNG compression removes only image meta-data with no impact on quality.",
				Optional:    true,
				Default:     true,
			},

			"tcp_pre_pooling": {
				Type:        schema.TypeBool,
				Description: "Maintain a set of idle TCP connections to the origin server to eliminate the latency associated with opening new connections or new requests (TCP handshake).",
				Optional:    true,
				Default:     true,
			},
			"origin_connection_reuse": {
				Type:        schema.TypeBool,
				Description: "TCP connections that are opened for a client request remain open for a short time to handle additional requests that may arrive.",
				Optional:    true,
				Default:     true,
			},
			"support_non_sni_clients": {
				Type:        schema.TypeBool,
				Description: "By default, non-SNI clients are supported. Disable this option to block non-SNI clients.",
				Optional:    true,
				Default:     true,
			},
			"enable_http2": {
				Type:        schema.TypeBool,
				Description: "Allows supporting browsers to take advantage of the performance enhancements provided by HTTP/2 for your website. Non-supporting browsers can connect via HTTP/1.0 or HTTP/1.1.",
				Optional:    true,
				Default:     false,
			},
			"http2_to_origin": {
				Type:        schema.TypeBool,
				Description: "Enables HTTP/2 for the connection between Imperva and your origin server. (HTTP/2 must also be supported by the origin server.)",
				Optional:    true,
				Default:     false,
			},
			"port_to": {
				Type:        schema.TypeInt,
				Description: "The port number. If field is set to 80 (the default value), rewrite port will be removed.",
				Optional:    true,
				Default:     defaultPortTo,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" && old == strconv.Itoa(defaultPortTo) {
						return true
					}
					return false
				},
			},
			"ssl_port_to": {
				Type:        schema.TypeInt,
				Description: "The port number to rewrite default SSL port to. if field is set to 443 (the default value), rewrite SSL port will be removed.",
				Optional:    true,
				Default:     defaultSslPortTo,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == "" && old == strconv.Itoa(defaultSslPortTo) {
						return true
					}
					return false
				},
			},
			"redirect_naked_to_full": {
				Type:        schema.TypeBool,
				Description: "Redirect all visitors to your site’s full domain (which includes www). This option is displayed only for a naked domain.",
				Optional:    true,
				Default:     false,
			},
			"redirect_http_to_https": {
				Type:        schema.TypeBool,
				Description: "Sites that require an HTTPS connection force all HTTP requests to be redirected to HTTPS. This option is displayed only for an SSL site.",
				Optional:    true,
				Default:     false,
			},

			"default_error_page_template": {
				Type:        schema.TypeString,
				Description: "The default error page HTML template. $TITLE$ and $BODY$ placeholders are required.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_connection_timeout": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Connection Timeout' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_access_denied": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Access Denied' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_parse_req_error": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Unable to parse request' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_parse_resp_error": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Unable to parse response' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_connection_failed": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Unable to connect to origin server' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_ssl_failed": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Unable to establish SSL connection' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_deny_and_captcha": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Initial connection denied - CAPTCHA required' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
			"error_no_ssl_config": {
				Type:        schema.TypeString,
				Description: "The HTML template for 'Site not configured for SSL' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.",
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						log.Printf("will supress error_ssl_failed¬")
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceApplicationDeliveryFlatRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	siteIdStr := strconv.Itoa(siteID)

	applicationDelivery, err := client.GetApplicationDelivery(siteID)
	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula Application Delivery for Site Id: %d - %s\n", siteID, err)
		return err
	}

	d.SetId(siteIdStr)

	d.Set("file_compression", applicationDelivery.Compression.FileCompression)
	d.Set("minify_js", applicationDelivery.Compression.MinifyJs)
	d.Set("minify_css", applicationDelivery.Compression.MinifyCss)
	d.Set("minify_static_html", applicationDelivery.Compression.MinifyStaticHtml)

	d.Set("compress_jpeg", applicationDelivery.ImageCompression.CompressJpeg)
	d.Set("progressive_image_rendering", applicationDelivery.ImageCompression.ProgressiveImageRendering)
	d.Set("aggressive_compression", applicationDelivery.ImageCompression.AggressiveCompression)
	d.Set("compress_png", applicationDelivery.ImageCompression.CompressPng)

	d.Set("tcp_pre_pooling", applicationDelivery.Network.TcpPrePooling)
	d.Set("origin_connection_reuse", applicationDelivery.Network.OriginConnectionReuse)
	d.Set("support_non_sni_clients", applicationDelivery.Network.SupportNonSniClients)
	d.Set("enable_http2", applicationDelivery.Network.EnableHttp2)
	d.Set("http2_to_origin", applicationDelivery.Network.Http2ToOrigin)
	d.Set("port_to", getPortTo(applicationDelivery.Network.Port.To, defaultPortTo))
	d.Set("ssl_port_to", getPortTo(applicationDelivery.Network.SslPort.To, defaultSslPortTo))

	d.Set("redirect_naked_to_full", applicationDelivery.Redirection.RedirectNakedToFull)
	d.Set("redirect_http_to_https", applicationDelivery.Redirection.RedirectHttpToHttps)

	d.Set("default_error_page_template", strings.ReplaceAll(applicationDelivery.CustomErrorPage.DefaultErrorPage, "'", "\""))
	d.Set("error_connection_timeout", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorConnectionTimeout, "'", "\""))
	d.Set("error_access_denied", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorAccessDenied, "'", "\""))
	d.Set("error_parse_req_error", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorParseReqError, "'", "\""))
	d.Set("error_parse_resp_error", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorParseRespError, "'", "\""))
	d.Set("error_connection_failed", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorConnectionFailed, "'", "\""))
	d.Set("error_ssl_failed", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorSslFailed, "'", "\""))
	d.Set("error_deny_and_captcha", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorDenyAndCaptcha, "'", "\""))
	d.Set("error_no_ssl_config", strings.ReplaceAll(applicationDelivery.CustomErrorPage.CustomErrorPageTemplates.ErrorTypeNoSslConfig, "'", "\""))

	return nil
}

func resourceApplicationDeliveryFlatUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)

	compression := Compression{
		FileCompression:  d.Get("file_compression").(bool),
		MinifyJs:         d.Get("minify_js").(bool),
		MinifyCss:        d.Get("minify_css").(bool),
		MinifyStaticHtml: d.Get("minify_static_html").(bool),
	}

	imageCompression := ImageCompression{
		CompressJpeg:              d.Get("compress_jpeg").(bool),
		ProgressiveImageRendering: d.Get("progressive_image_rendering").(bool),
		AggressiveCompression:     d.Get("aggressive_compression").(bool),
		CompressPng:               d.Get("compress_png").(bool),
	}

	networkd := Network{
		TcpPrePooling:         d.Get("tcp_pre_pooling").(bool),
		OriginConnectionReuse: d.Get("origin_connection_reuse").(bool),
		SupportNonSniClients:  d.Get("support_non_sni_clients").(bool),
		EnableHttp2:           d.Get("enable_http2").(bool),
		Http2ToOrigin:         d.Get("http2_to_origin").(bool),
		Port:                  Port{To: strconv.Itoa(d.Get("port_to").(int))},
		SslPort:               SslPort{To: strconv.Itoa(d.Get("ssl_port_to").(int))},
	}

	redirection := Redirection{
		RedirectNakedToFull: d.Get("redirect_naked_to_full").(bool),
		RedirectHttpToHttps: d.Get("redirect_http_to_https").(bool),
	}
	customErrorPage := CustomErrorPage{
		DefaultErrorPage: d.Get("default_error_page_template").(string),
		CustomErrorPageTemplates: CustomErrorPageTemplates{
			ErrorConnectionTimeout: d.Get("error_connection_timeout").(string),
			ErrorAccessDenied:      d.Get("error_access_denied").(string),
			ErrorParseReqError:     d.Get("error_parse_req_error").(string),
			ErrorParseRespError:    d.Get("error_parse_resp_error").(string),
			ErrorConnectionFailed:  d.Get("error_connection_failed").(string),
			ErrorSslFailed:         d.Get("error_ssl_failed").(string),
			ErrorDenyAndCaptcha:    d.Get("error_deny_and_captcha").(string),
			ErrorTypeNoSslConfig:   d.Get("error_no_ssl_config").(string),
		},
	}

	payload := ApplicationDelivery{
		Compression:      compression,
		ImageCompression: imageCompression,
		Network:          networkd,
		Redirection:      redirection,
		CustomErrorPage:  customErrorPage,
	}

	_, err := client.UpdateApplicationDelivery(siteID, &payload)

	if err != nil {
		log.Printf("[ERROR] Could not get Incapsula Application Delivery for Site Id: %d - %s\n", siteID, err)
		return err
	}
	return resourceApplicationDeliveryFlatRead(d, m)
}

func resourceApplicationDeliveryFlatDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)

	_, err := client.DeleteApplicationDelivery(siteID)

	if err != nil {
		log.Printf("[ERROR] Could delete Incapsula Application Delivery for Site Id: %d - %s\n", siteID, err)
		return err
	}

	d.SetId("")
	return nil
}

func getPortTo(port string, defaultValue int) int {
	portTo := defaultValue
	if port != "-" {
		portToInt, err := strconv.Atoi(port)
		if err == nil {
			portTo = portToInt
		}
	}
	return portTo
}
