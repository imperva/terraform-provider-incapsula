package incapsula

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetGcpRegions(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_gcp_regions_test.TestGetGcpRegions")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/anywhere-provisioner/v1/igc/regions" {
			t.Errorf("Should have hit the GCP regions endpoint. Got path: %s", req.URL.Path)
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[\"us-east1\",\"europe-west1\"]}"))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetGcpRegions("")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(resp.Data) != 2 || resp.Data[0] != "us-east1" || resp.Data[1] != "europe-west1" {
		t.Errorf("Unexpected regions parsed: %+v", resp.Data)
	}
}

func TestGetGcpRegionsWithAccountId(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_gcp_regions_test.TestGetGcpRegionsWithAccountId")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Query().Get("caid") != "456" {
			t.Errorf("Should pass caid=456 when account id set. Got: %s", req.URL.Query().Get("caid"))
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[\"us-east1\"]}"))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	if _, err := client.GetGcpRegions("456"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestGcpRegionsDataSourceSortsRegions(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_gcp_regions_test.TestGcpRegionsDataSourceSortsRegions")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[\"us-east1\",\"asia-south1\",\"europe-west1\"]}"))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	d := schema.TestResourceDataRaw(t, dataSourceIncapsulaGcpRegions().Schema, map[string]interface{}{})
	if diags := dataSourceIncapsulaGcpRegionsRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	regions := d.Get("regions").([]interface{})
	expected := []string{"asia-south1", "europe-west1", "us-east1"}
	if len(regions) != len(expected) {
		t.Fatalf("Expected %d regions. Got: %v", len(expected), regions)
	}
	for i, r := range expected {
		if regions[i] != r {
			t.Errorf("Expected sorted regions %v. Got: %v", expected, regions)
			break
		}
	}
}
