package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// //////////////////////////////////////////////////////////////
// UpdateSiteMonitoring Tests
// //////////////////////////////////////////////////////////////
func TestUpdateApplicationDeliveryBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	applicationDeliveryPayload := ApplicationDelivery{
		Compression:      Compression{},
		ImageCompression: ImageCompression{},
		Network:          Network{},
		Redirection:      Redirection{},
	}

	applicationDeliveryResponse, err := client.UpdateApplicationDelivery(
		siteID,
		&applicationDeliveryPayload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

	customErrorPagesPayload := CustomErrorPage{}

	errorPages, err := client.UpdateErrorPages(
		siteID,
		&customErrorPagesPayload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if errorPages != nil {
		t.Errorf("Should have received a nil errorPages instance")
	}
}

func TestUpdateApplicationDeliveryBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	applicationDeliveryEndpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)
	errorPagesEndpoint := fmt.Sprintf("/sites/%d/settings/delivery/error-pages", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != applicationDeliveryEndpoint && req.URL.String() != errorPagesEndpoint {
			t.Errorf("Should have have hit delivery settings endpoint. Got: %s", req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	applicationDeliveryPayload := ApplicationDelivery{
		Compression:      Compression{},
		ImageCompression: ImageCompression{},
		Network:          Network{},
		Redirection:      Redirection{},
	}

	applicationDeliveryResponse, err := client.UpdateApplicationDelivery(
		siteID,
		&applicationDeliveryPayload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

	customErrorPagesPayload := CustomErrorPage{}

	errorPages, err := client.UpdateErrorPages(
		siteID,
		&customErrorPagesPayload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if errorPages != nil {
		t.Errorf("Should have received a nil errorPages instance")
	}
}

func TestUpdateApplicationDeliveryInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	applicationDeliveryEndpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)
	errorPagesEndpoint := fmt.Sprintf("/sites/%d/settings/delivery/error-pages", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != applicationDeliveryEndpoint && req.URL.String() != errorPagesEndpoint {
			t.Errorf("Should have have hit delivery settings endpoint. Got: %s", req.URL.String())
		}
		rw.Write([]byte(`{
  "errors": [
    {
      "id": null,
      "status": 500,
      "source": {
        "pointer": "/sites/32/settings/delivery"
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
	payload := ApplicationDelivery{}

	applicationDeliveryResponse, diags := client.UpdateApplicationDelivery(
		siteID,
		&payload)

	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when %s Application Delivery for Site ID %d", 500, siteID)) {
		t.Errorf("Should have received a bad Application Delivery error, got: %s", diags[0].Detail)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

	customErrorPagesPayload := CustomErrorPage{}

	errorPages, diags := client.UpdateErrorPages(
		siteID,
		&customErrorPagesPayload)

	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when %s Error Pages for Site ID %d", 500, siteID)) {
		t.Errorf("Should have received a bad Application Delivery error, got: %s", diags[0].Detail)
	}
	if errorPages != nil {
		t.Errorf("Should have received a nil errorPages instance")
	}
}

func TestUpdateApplicationDeliveryConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	applicationDeliveryEndpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)
	errorPagesEndpoint := fmt.Sprintf("/sites/%d/settings/delivery/error-pages", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != applicationDeliveryEndpoint && req.URL.String() != errorPagesEndpoint {
			t.Errorf("Should have have hit delivery settings endpoint. Got: %s", req.URL.String())
		}
		if req.URL.String() == applicationDeliveryEndpoint {
			rw.Write([]byte(`
				{
				"compression": {
					"file_compression": true,
					"compression_type": "GZIP",
					"minify_js": true,
					"minify_css": false,
					"minify_static_html": true
				},
				"image_compression": {
					"compress_jpeg": true,
					"progressive_image_rendering": true,
					"aggressive_compression": false,
					"compress_png": true
				},
				"network": {
					"tcp_pre_pooling": true,
					"origin_connection_reuse": false,
					"support_non_sni_clients": true,
					"port": {
					"to": "8080"
					},
					"ssl_port": {
					"to": "9001"
					}
				},
				"redirection": {
					"redirect_naked_to_full": false,
					"redirect_http_to_https": true
				}
				}
			`))
		} else {
			rw.Write([]byte(`
				{
					"custom_error_page": {
					"error_page_template": "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>",
					"custom_error_page_templates": {
						"error.type.connection_timeout": "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>",
						"error.type.access_denied": "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>"
					}
					}
				}
			`))
		}
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	payload := ApplicationDelivery{
		Compression:      Compression{},
		ImageCompression: ImageCompression{},
		Network:          Network{},
		Redirection:      Redirection{},
	}

	customErrorPagesPayload := CustomErrorPage{}

	applicationDeliveryResponse, diags := client.UpdateApplicationDelivery(
		siteID,
		&payload)

	if diags != nil {
		t.Errorf("Should not have received an error : %s", diags[0].Detail)
	}
	if applicationDeliveryResponse.Network.SupportNonSniClients != true {
		t.Errorf("Should have received a SupportNonSniClients equal true\n%v", applicationDeliveryResponse)
	}

	errorPages, diags := client.UpdateErrorPages(
		siteID,
		&customErrorPagesPayload)
	if diags != nil {
		t.Errorf("Should not have received an error : %s", diags[0].Detail)
	}
	if errorPages.DefaultErrorPage == "" || errorPages.CustomErrorPageTemplates.ErrorConnectionTimeout == "" || errorPages.CustomErrorPageTemplates.ErrorConnectionFailed != "" || errorPages.DefaultErrorPage != "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>" {
		t.Errorf("unexpected error page response: %v", errorPages)
	}
}

// //////////////////////////////////////////////////////////////
// ReadSiteMonitoring Tests
// //////////////////////////////////////////////////////////////
func TestReadApplicationDeliveryBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	applicationDeliveryResponse, err := client.GetApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

	errorPages, err := client.GetErrorPages(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if errorPages != nil {
		t.Errorf("Should have received a nil errorPages instance")
	}
}

func TestReadApplicationDeliveryBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	applicationDeliveryEndpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)
	errorPagesEndpoint := fmt.Sprintf("/sites/%d/settings/delivery/error-pages", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != applicationDeliveryEndpoint && req.URL.String() != errorPagesEndpoint {
			t.Errorf("Should have have hit delivery settings endpoint. Got: %s", req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	applicationDeliveryResponse, err := client.GetApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

	errorPages, err := client.GetErrorPages(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if errorPages != nil {
		t.Errorf("Should have received a nil errorPages instance")
	}

}

func TestReadApplicationDeliveryInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	applicationDeliveryEndpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)
	errorPagesEndpoint := fmt.Sprintf("/sites/%d/settings/delivery/error-pages", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != applicationDeliveryEndpoint && req.URL.String() != errorPagesEndpoint {
			t.Errorf("Should have have hit delivery settings endpoint. Got: %s", req.URL.String())
		}
		rw.Write([]byte(`{
  "errors": [
    {
      "id": null,
      "status": 500,
      "source": {
        "pointer": "/sites/%d/settings/delivery"
      },
      "title": "Not Found"
    }
  ]
}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	applicationDeliveryResponse, diags := client.GetApplicationDelivery(siteID)

	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when %s Application Delivery for Site ID %d", 500, siteID)) {
		t.Errorf("Should have received a bad Application Delivery error, got: %s", diags[0].Detail)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

	errorPages, diags := client.GetErrorPages(siteID)

	if diags == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when %s Error Pages for Site ID %d", 500, siteID)) {
		t.Errorf("Should have received a bad Application Delivery error, got: %s", diags[0].Detail)
	}
	if errorPages != nil {
		t.Errorf("Should have received a nil errorPages instance")
	}
}

func TestReadApplicationDeliveryConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42
	endpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`
{
  "compression": {
    "file_compression": true,
    "compression_type": "GZIP",
    "minify_js": true,
    "minify_css": false,
    "minify_static_html": true
  },
  "image_compression": {
    "compress_jpeg": true,
    "progressive_image_rendering": true,
    "aggressive_compression": false,
    "compress_png": true
  },
  "network": {
    "tcp_pre_pooling": true,
    "origin_connection_reuse": false,
    "support_non_sni_clients": true,
    "port": {
      "to": "8080"
    },
    "ssl_port": {
      "to": "9001"
    }
  },
  "redirection": {
    "redirect_naked_to_full": false,
    "redirect_http_to_https": true
  },
  "custom_error_page": {
    "error_page_template": "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>",
    "custom_error_page_templates": {
      "error.type.connection_timeout": "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>",
      "error.type.access_denied": "<html><body><h1>$TITLE$</h1><div>$BODY$</div></body></html>"
    }
  }
}
`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	applicationDeliveryResponse, err := client.GetApplicationDelivery(siteID)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if applicationDeliveryResponse == nil {
		t.Errorf("Should not have received a nil applicationDeliveryResponse instance")
	}
	if applicationDeliveryResponse.Network.SupportNonSniClients != true {
		t.Errorf("Should have received a SupportNonSniClients equal true\n%v", applicationDeliveryResponse)
	}
}

// //////////////////////////////////////////////////////////////
// DeleteSiteMonitoring Tests
// //////////////////////////////////////////////////////////////
func TestDeleteApplicationDeliveryBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	applicationDeliveryResponse, err := client.DeleteApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error from Incapsula service when trying to delete Application Delivery for Site ID 42:") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}
}

func TestDeleteApplicationDeliveryBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	applicationDeliveryResponse, err := client.DeleteApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error parsing Application Delivery Response JSON response for Site ID") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}

}

func TestDeleteApplicationDeliveryInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{"res":2,"res_message":"Invalid input","debug_info":{"ssl_port":["ssl is not supported for your site"],"id-info":"999999"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	applicationDeliveryResponse, err := client.DeleteApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from Incapsula service when Deleting Application Delivery for Site ID")) {
		t.Errorf("Should have received an internal server error for Application Delivery, got: %s", err)
	}

	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}
}

func TestDeleteApplicationDeliveryConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
		}
		rw.Write([]byte(`{
		"value": "Deletion successful",
		"isError": false
	}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	_, _, err := client.DeleteApplicationDelivery(
		siteID)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
