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
func TestUpdateApplicationDeliveryBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	payload := ApplicationDelivery{
		Compression:      Compression{},
		ImageCompression: ImageCompression{},
		Network:          Network{},
		Redirection:      Redirection{},
		CustomErrorPage:  CustomErrorPage{},
	}

	applicationDeliveryResponse, err := client.UpdateApplicationDelivery(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error from Incapsula service when trying to update Application Delivery for Site ID 42") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}
}

func TestUpdateApplicationDeliveryBadJSON(t *testing.T) {
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

	payload := ApplicationDelivery{
		Compression:      Compression{},
		ImageCompression: ImageCompression{},
		Network:          Network{},
		Redirection:      Redirection{},
		CustomErrorPage:  CustomErrorPage{},
	}

	applicationDeliveryResponse, err := client.UpdateApplicationDelivery(
		siteID,
		&payload)

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

func TestUpdateApplicationDeliveryInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
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

	applicationDeliveryResponse, err := client.UpdateApplicationDelivery(
		siteID,
		&payload)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from Incapsula service when Updating Application Delivery for Site ID")) {
		t.Errorf("Should have received a bad Application Delivery error, got: %s", err)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}
}

func TestUpdateApplicationDeliveryConfig(t *testing.T) {
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
	payload := ApplicationDelivery{
		Compression:      Compression{},
		ImageCompression: ImageCompression{},
		Network:          Network{},
		Redirection:      Redirection{},
		CustomErrorPage:  CustomErrorPage{},
	}

	applicationDeliveryResponse, err := client.UpdateApplicationDelivery(
		siteID,
		&payload)

	if err != nil {
		t.Errorf("Should not have received an error : %s", err.Error())
	}
	if applicationDeliveryResponse.Network.SupportNonSniClients != true {
		t.Errorf("Should have received a SupportNonSniClients equal true\n%v", applicationDeliveryResponse)
	}
}

////////////////////////////////////////////////////////////////
// ReadSiteMonitoring Tests
////////////////////////////////////////////////////////////////
func TestReadApplicationDeliveryBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := 42

	applicationDeliveryResponse, err := client.GetApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), "[ERROR] Error from Incapsula service when trying to read Application Delivery for Site ID 42") {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
	}
}

func TestReadApplicationDeliveryBadJSON(t *testing.T) {
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

	applicationDeliveryResponse, err := client.GetApplicationDelivery(siteID)

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

func TestReadApplicationDeliveryInvalidConfig(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := 42

	endpoint := fmt.Sprintf("/sites/%d/settings/delivery", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(500)
		if req.URL.String() != endpoint {
			t.Errorf("Should have have hit %s endpoint. Got: %s", endpoint, req.URL.String())
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

	applicationDeliveryResponse, err := client.GetApplicationDelivery(siteID)

	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error status code 500 from Incapsula service when Reading Application Delivery for Site ID")) {
		t.Errorf("Should have received an internal server error for Application Delivery, got: %s", err)
	}

	if applicationDeliveryResponse != nil {
		t.Errorf("Should have received a nil applicationDeliveryResponse instance")
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

////////////////////////////////////////////////////////////////
// DeleteSiteMonitoring Tests
////////////////////////////////////////////////////////////////
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

	_, err := client.DeleteApplicationDelivery(
		siteID)

	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
