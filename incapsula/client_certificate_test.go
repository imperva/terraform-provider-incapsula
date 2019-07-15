package incapsula

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

////////////////////////////////////////////////////////////////
// AddCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientAddCertificateBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAddCertificateBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "abc", "def", "efg")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding custom certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test Running test client_certificate_test.TestClientAddCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateAdd, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing add custom certificate JSON response for site_id %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateInvalidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAddCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateAdd, req.URL.String())
		}
		rw.Write([]byte(`{"res":3015,"res_message":"Internal error","debug_info":{"id-info":"13008","Error":"Unexpected error occurred"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when adding custom certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a bad certificate error, got: %s", err)
	}
	if addCertificateResponse != nil {
		t.Errorf("Should have received a nil addCertificateResponse instance")
	}
}

func TestClientAddCertificateValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAddCertificateValidRule")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateAdd) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateAdd, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	//certificate := "-----BEGIN CERTIFICATE-----\nMIIDgjCCAmoCCQCk3MsAS5x+UjANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlTYW4gRGllZ28xCzAJBgNVBAoMAlNF\nMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFzaC5iZWVyLmNlbnRlcjEdMBsGCSqG\nSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wHhcNMTkwNzA4MTU0MjQ0WhcNMjAwNzA3\nMTU0MjQ0WjCBgjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlT\nYW4gRGllZ28xCzAJBgNVBAoMAlNFMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFz\naC5iZWVyLmNlbnRlcjEdMBsGCSqGSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCj0rYKUhVtNKQ/oKZdCfxvLKhQ\nLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1DhPRSDGgONdHe\n2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkPQ8dsPpE945lW\nsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607DqdZUS/O4XSx+d0\n5kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1wqRdwGJsUdq6\nkC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2uWCBSY/OLAgMB\nAAEwDQYJKoZIhvcNAQELBQADggEBABfNZcItHdsSpfp8h+1EP5BnRuoKj+l42EI5\nE9dVlqdOZ25+V5Ee899sn2Nj8h+/zVU3+IDO2abUPrDd2xZHaHdf0p69htSwFTHs\nEwUdPUUsKRSys7fVP1clHcKWswTcoWIzQiPZsDMoOQw/pzN05cXSzdo8wSWuEeBK\ncqRNd5BKPeeXbFa4i5TFzT/+pl8V075k16tzHSbT7QDk5fuZWYv/2jImw/lgS/nx\nDWtlprrgG6AX1FzovDs/NnNq/e7vZtn8sdOoO2pCSVymNvctNLV2tFcS8sPQDl5M\nIpnZa3kktAegjsCln1JvD0AFigXrF8wjK+FKGI8SPJfbTQ149+A=\n-----END CERTIFICATE-----"
	//private_key := "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCj0rYKUhVtNKQ/\noKZdCfxvLKhQLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1D\nhPRSDGgONdHe2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkP\nQ8dsPpE945lWsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607Dqd\nZUS/O4XSx+d05kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1\nwqRdwGJsUdq6kC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2u\nWCBSY/OLAgMBAAECggEAfDPprkNzWTta95594vYKO+OYrEknnRUwRV0LF/ziae2P\nLR1EX0AeXKpIXwwwDzpXITJS7mu1c8uZwTIQ5g/f6D4nopULVYwlJZhbjXd49hpx\nhmGfk8227te5BqnVS3IPvRx5vjz+r8obYFZb4JZDGa/v9okAlI04FS0hR/Bl4ckD\naIsztf4R+AO2dP6BxYZGIwcq3jkbf0BdyQpkw4Ds7pdKbSa+PsobseyI2NqR2ryX\n4HH4b89HZj8lfiniIN3tPV6uIvpPS6jJklLKy6zdkIFOng/OGwxXomGkrk9ZjBHm\nJx5yA5YfwPidyt80wO9/26wClXYidfKQC8mDN21owQKBgQDPQbNr/sGiI2QzTOpb\nYTx0FWzWMnn9N2XiQm5rcr9kM5WsXh+anlqP54MeXDGZ2f6L8+aGrghZ/78WbG9J\nDbtEc7qTSRw5LFRglqn32a3ppHToEzOVxsA3g/OBJT5lJJwGMTdeKEXtLMmkm/sz\n1ClFnYJ1I8rNcueI9936odDWKwKBgQDKWgGwWTbqVa3wVIOFvluxolQzo6TEBFbf\nQTJo7byO2iRZvhrZUUk8539Uz2px0Ilzxx61CszhNWDVNwgqsN7FtuzXuCwz9GzU\nyBWkzPKGzvK12aFMYoj/cPbcRfMpYWNoK/YfEKfTRkJJfrJSbWP2XlyEr69te8s7\nB/zxOtUIIQKBgEjoJcOhtF/i70aUkgRfKjLzrnuS+hK3QCHdmJY3oVgQRWCDI77y\nYY0ptZgielhStRZqT/eklM+EBaZPsr4SFIQ56bISD9mU3IG1vkivzFvaPD2/M3BG\noCtnQWt2vII75J7RBVcb9609ChnbvPw4b+RLSi8GzjqDZytpdi7KaXpNAoGAS2Ym\nYvObRs4ONhMHvvojaJk4DtXXO0Lyq9W7VuXe8MvP57CyiG+FfrAz/gIbg7VUwlNb\n2dHgbbpaDpim7mFhYQK8VdVGg0V8l/zGM9Y6OIk8Xw5sz+2XZrdNBN77sFudkt9u\nojyujEcNxBz1jUk9iju29aoREBakr6ZWVfy6DIECgYEAtXxrOsbMsbHhVGqgeGXy\nhLXIltR+7NIUaxpLHhYCMzK9SbyZvx/Hd6m34oTw9ws+tHFpeCyiVU+wQgmx0ARD\ncDLKOPIHTGYhq/H8Oc6/Dzfxs1L/hH34mw5u7hVtAaA+q8iaRGVZ797dTVSxw4U0\nRm+BCDRhDcvaG7qpvFj8T6k=\n-----END PRIVATE KEY-----"
	// passphrase := "webco123"
	addCertificateResponse, err := client.AddCertificate(siteID, "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if addCertificateResponse == nil {
		t.Errorf("Should not have received a nil addCertificateResponse instance")
	}
	if addCertificateResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// ListCertificates Tests
////////////////////////////////////////////////////////////////

func TestClientListCertificatesBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientListCertificatesBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting custom certificates for site_id %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientListCertificatesBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing certificates list JSON response for site_id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

func TestClientListCertificatesInvalidRequest(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientListCertificatesInvalidRequest")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"id-info":"13007","site_id":"1234"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	listCertificatesResponse, err := client.ListCertificates(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when getting custom certificates list for site_id %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
	if listCertificatesResponse != nil {
		t.Errorf("Should have received a nil listCertificatesResponse instance")
	}
}

//func TestClientListCertificatesValidRequest(t *testing.T) {
//	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateList) {
//			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateList, req.URL.String())
//		}
//		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13008"},"site_id":1234,"active":"active"}`))
//	}))
//	defer server.Close()
//
//	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
//	client := &Client{config: config, httpClient: &http.Client{}}
//	siteID := "1234"
//	listCertificatesResponse, err := client.ListCertificates(siteID)
//
//	if err != nil {
//		t.Errorf("Should not have received an error")
//	}
//	if listCertificatesResponse == nil {
//		t.Errorf("Should not have received a nil listCertificatesResponse instance")
//	}
//
//	if listCertificatesResponse.Res != 0 {
//		t.Errorf("Response code doesn't match")
//	}
//}

////////////////////////////////////////////////////////////////
// EditCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientEditCertificateBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientEditCertificateBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	certificate := "foo"
	privateKey := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, privateKey, passphrase)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error editing custom certificate for site_id: %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

func TestClientEditCertificateBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientEditCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	privateKey := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, privateKey, passphrase)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing edit custom certificarte JSON response for site_id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if editCertificateResponse != nil {
		t.Errorf("Should have received a nil editCertificateResponse instance")
	}
}

//func TestClientEditCertificateInvalidRule(t *testing.T) {
//	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
//		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
//			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
//		}
//		rw.Write([]byte(`{"rule_id":0,"res":"1"}`))
//	}))
//	defer server.Close()
//
//	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
//	client := &Client{config: config, httpClient: &http.Client{}}
//	siteID := "1234"
//	certificate := "foo"
//	private_key := "bar"
//	passphrase := "loremipsum"
//	editCertificateResponse, err := client.EditCertificate(siteID, certificate, private_key, passphrase)
//	if err == nil {
//		t.Errorf("Should have received an error")
//	}
//	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when editing custom certificate for siteID %s", siteID)) {
//		t.Errorf("Should have received a bad site error, got: %s", err)
//	}
//	if editCertificateResponse != nil {
//		t.Errorf("Should have received a nil editCertificateResponse instance")
//	}
//}

func TestClientEditCertificateValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientEditCertificateValidRule")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateEdit) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateEdit, req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK","debug_info":{"id-info":"13008"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	certificate := "foo"
	privateKey := "bar"
	passphrase := "loremipsum"
	editCertificateResponse, err := client.EditCertificate(siteID, certificate, privateKey, passphrase)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if editCertificateResponse == nil {
		t.Errorf("Should not have received a nil editCertificateResponse instance")
	}

	if editCertificateResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}

////////////////////////////////////////////////////////////////
// DeleteCertificate Tests
////////////////////////////////////////////////////////////////

func TestClientDeleteCertificateBadConnection(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateBadConnection")
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting custom certificate for site_id: %s", siteID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientDeleteCertificateBadJSON(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateBadJSON")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error deleting custom certificate for site_id: %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
}

func TestClientDeleteCertificateInvalidSiteID(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientDeleteCertificateInvalidSiteID")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":9413,"res_message":"Unknown/unauthorized site_id","debug_info":{"id-info":"13008","site_id":"1234"}}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting custom certificate for site_id %s", siteID)) {
		t.Errorf("Should have received a bad site error, got: %s", err)
	}
}

func TestClientDeleteCertificateValidSite(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG]Running test client_certificate_test.TestClientDeleteCertificateValidSite")
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("/%s", endpointCertificateDelete) {
			t.Errorf("Should have have hit /%s endpoint. Got: %s", endpointCertificateDelete, req.URL.String())
		}
		rw.Write([]byte(`{"res":"0","res_message":"OK"}`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	siteID := "1234"
	err := client.DeleteCertificate(siteID)
	if err != nil {
		t.Errorf("Should not have received an error")
	}
}
