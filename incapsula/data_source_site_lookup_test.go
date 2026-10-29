package incapsula

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func readSiteLookup(t *testing.T, body string) *schema.ResourceData {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte(body))
	}))
	defer server.Close()
	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	d := schema.TestResourceDataRaw(t, dataSourceSite().Schema, map[string]interface{}{"name": "lb.example.com"})
	if diags := dataSourceSiteRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected diags: %v", diags)
	}
	return d
}

func TestSiteLookupReturnsLoadBalancerSiteAttributes(t *testing.T) {
	d := readSiteLookup(t, `{"data":[{"id":123,"name":"lb.example.com","type":"PUBLIC_CLOUD","cloud":"GCP","cname":"abc.impervadns.net","accountId":1,"isLoadBalancerSite":true,"authorityHeader":"abc-1.example.net"}]}`)
	if d.Get("cloud_type") != "GCP" {
		t.Errorf("Expected cloud_type GCP. Got: %v", d.Get("cloud_type"))
	}
	if d.Get("is_load_balancer_site") != true {
		t.Errorf("Expected is_load_balancer_site true. Got: %v", d.Get("is_load_balancer_site"))
	}
	if d.Get("authority_header") != "abc-1.example.net" {
		t.Errorf("Expected authority_header abc-1.example.net. Got: %v", d.Get("authority_header"))
	}
}

func TestSiteLookupReturnsEmptyGoogleCloudAttributesForAwsSite(t *testing.T) {
	d := readSiteLookup(t, `{"data":[{"id":124,"name":"lb.example.com","type":"PUBLIC_CLOUD","cloud":"AWS","cname":"abc.impervadns.net","accountId":1}]}`)
	if d.Get("cloud_type") != "AWS" {
		t.Errorf("Expected cloud_type AWS. Got: %v", d.Get("cloud_type"))
	}
	if d.Get("is_load_balancer_site") != false {
		t.Errorf("Expected is_load_balancer_site false. Got: %v", d.Get("is_load_balancer_site"))
	}
	if d.Get("authority_header") != "" {
		t.Errorf("Expected empty authority_header. Got: %v", d.Get("authority_header"))
	}
}
