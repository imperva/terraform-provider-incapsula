package incapsula

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPscTargetDecodesDataField(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_gcp_psc_target_test.TestGetPscTargetDecodesDataField")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/anywhere-provisioner/v1/igc/psc-target" {
			t.Errorf("Should have hit the GCP psc-target endpoint. Got path: %s", req.URL.Path)
		}
		if req.URL.Query().Get("region") != "us-east1" {
			t.Errorf("Should pass region=us-east1. Got: %s", req.URL.Query().Get("region"))
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":\"projects/p/regions/us-east1/serviceAttachments/sa\"}"))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetPscTarget("us-east1", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if resp.Address != "projects/p/regions/us-east1/serviceAttachments/sa" {
		t.Errorf("Expected address decoded from data field. Got: %s", resp.Address)
	}
}

func TestGetPscTargetUrlEscapesQueryValues(t *testing.T) {
	url := getPscTargetUrl("https://api.example.com", "us east1&x=1", "12 3")
	expected := "https://api.example.com/anywhere-provisioner/v1/igc/psc-target?region=us+east1%26x%3D1&caid=12+3"
	if url != expected {
		t.Errorf("Expected %s. Got: %s", expected, url)
	}
}

func TestGetPscTargetReturnsApiErrorDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(422)
		rw.Write([]byte(`{"errors":[{"status":422,"title":"Site has no parameters","detail":"this site has no parameters attached"}]}`))
	}))
	defer server.Close()
	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	_, err := client.GetPscTarget("us-east1", "")
	if err == nil {
		t.Fatalf("Expected an error")
	}
	if !strings.Contains(err.Error(), "this site has no parameters attached") || strings.Contains(err.Error(), "{") {
		t.Errorf("Expected the API detail without the raw body. Got: %s", err)
	}
}

func TestGetPscTargetFallsBackToRawBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(502)
		rw.Write([]byte("bad gateway"))
	}))
	defer server.Close()
	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	_, err := client.GetPscTarget("us-east1", "")
	if err == nil || !strings.Contains(err.Error(), "bad gateway") {
		t.Errorf("Expected the raw body in the error. Got: %v", err)
	}
}
