---
layout: "incapsula"
page_title: "Incapsula: application_delivery"
sidebar_current: "docs-incapsula-resource-application_delivery"
description: |-
  Provides a Incapsula Application Delivery resource.
---

# incapsula_application_delivery

Configure delivery options to help you optimize your content delivery and improve performance by providing faster loading of your web pages.
Note that the destroy command returns the configuration to the default values.

## Example Usage

### Basic Usage - Application Delivery

```hcl
resource "incapsula_application_delivery" "example_application_delivery" {
	site_id = incapsula_site.testacc-terraform-site.id
	file_compression            = true
	compression_type            = "GZIP"
	minify_css                  = true
	minify_js                   = true
	minify_static_html          = false
	default_error_page_template = "<html><body><h1>$TITLE$</h1><p>$BODY$</p><div>1</div></body></html>"
	error_access_denied         = "<html><body><h1>$TITLE$</h1><p>$BODY$</p><div>1</div></body></html>"
	error_connection_failed     = "${file("error_page_example.txt")}"
	aggressive_compression      = true
	compress_jpeg               = false
	compress_png                = true
	progressive_image_rendering = true
	enable_http2                = false
	http2_to_origin             = false
	origin_connection_reuse     = false
	port_to                     = 225
	ssl_port_to                 = 555
	support_non_sni_clients     = false
	tcp_pre_pooling             = false
	redirect_http_to_https      = false
	redirect_naked_to_full      = false			
}
```

## Argument Reference

The following arguments are supported:

* `site_id` - (Required) Numeric identifier of the site to operate on.
* `file_compression` - (Optional) When this option is enabled, files such as JavaScript, CSS and HTML are dynamically compressed using the selected format as they are transferred. They are automatically unzipped within the browser. If Brotli is not supported by the browser, files are automatically sent in Gzip. Default: true
* `compression_type` - (Optional) BROTLI (recommended for more efficient compression). Default: GZIP
* `minify_js` - (Optional) Minify JavaScript. Minification removes characters that are not necessary for rendering the page, such as whitespace and comments. This makes the files smaller and therefore reduces their access time. Minification has no impact on the functionality of the Javascript, CSS, and HTML files. Default: true
* `minify_css` - (Optional) Content minification can applied only to cached Javascript, CSS and HTML content. Default: true.
* `minify_static_html` - (Optional) Minify static HTML. Default: true.
* `compress_jpeg` - (Optional) Compress JPEG images. Compression reduces download time by reducing the file size. Default: true
* `progressive_image_rendering` - (Optional) The image is rendered with progressively finer resolution, potentially causing a pixelated effect until the final image is rendered with no loss of quality. This option reduces page load times and allows images to gradually load after the page is rendered. Default: false.
* `aggressive_compression` - (Optional) A more aggressive method of compression is applied with the goal of minimizing the image file size, possibly impacting the final quality of the image displayed. Applies to JPEG compression only. Default: false.
* `compress_png` - (Optional) Compress PNG images. Compression reduces download time by reducing the file size. PNG compression removes only image meta-data with no impact on quality. Default: true.
* `tcp_pre_pooling` - (Optional) Maintain a set of idle TCP connections to the origin server to eliminate the latency associated with opening new connections or new requests (TCP handshake). Default: true
* `origin_connection_reuse` - (Optional) TCP connections that are opened for a client request remain open for a short time to handle additional requests that may arrive. Default: true
* `support_non_sni_clients` - (Optional) By default, non-SNI clients are supported. Disable this option to block non-SNI clients. Default: true
* `enable_http2` - (Optional) Allows supporting browsers to take advantage of the performance enhancements provided by HTTP/2 for your website. Non-supporting browsers can connect via HTTP/1.0 or HTTP/1.1.
* `http2_to_origin` - (Optional) Enables HTTP/2 for the connection between Imperva and your origin server. (HTTP/2 must also be supported by the origin server.)
* `port_to` - (Optional) The port number.
* `ssl_port_to` - (Optional) The port number to rewrite default SSL port to.
* `redirect_naked_to_full` - (Optional) Redirect all visitors to your site’s full domain (which includes www). This option is displayed only for a naked domain. Default: false
* `redirect_http_to_https`- (Optional) Sites that require an HTTPS connection force all HTTP requests to be redirected to HTTPS. This option is displayed only for an SSL site. Default: false
* `default_error_page_template` - (Optional) The default error page HTML template. $TITLE$ and $BODY$ placeholders are required.
* `error_connection_timeout`- (Optional) The HTML template for 'Connection Timeout' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_access_denied`- (Optional) The HTML template for 'Access Denied' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_parse_req_error`- (Optional) The HTML template for 'Unable to parse request' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_parse_resp_error`- (Optional) The HTML template for 'Unable to parse response' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_connection_failed`- (Optional) The HTML template for 'Unable to connect to origin server' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_ssl_failed`- (Optional) The HTML template for 'Unable to establish SSL connection' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_deny_and_captcha`- (Optional) The HTML template for 'Initial connection denied - CAPTCHA required' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_no_ssl_config`- (Optional) The HTML template for 'Site not configured for SSL' error. $TITLE$ and $BODY$ placeholders are required. Set empty value to return to default.
* `error_abp_identification_failed`- (Optional) The HTML template for 'ABP identification failed' error. Only HTML elements located inside the body tag are supported. Set empty value to return to default.


## Attributes Reference

The following attributes are exported:

* `id` - Unique identifier in the API for the application delivery configuration. The id is identical to Site id.

## Import

Application Delivery configuration can be imported using the `id`, e.g.:

```
$ terraform import incapsula_application_delivery.example_application_delivery 1234
```
