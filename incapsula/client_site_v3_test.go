package incapsula

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddV3SiteWithNameAndType(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_v3_test.TestAddV3SiteWithName")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", endpointSiteV3+"?caid=123") {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpointSiteV3+"?caid=123", req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n  \"data\": [\n    {\n      \"id\": 462102065,\n      \"name\": \"de3affdrere.inddcapcwafteam.net\",\n      \"type\": \"CLOUD_WAF\",\n      \"accountId\": 51999737,\n      \"creationTime\": 1717588301055,\n      \"cname\": \"mhhp8q4.ng.impervadnsstage.net\"\n    }\n  ]\n}"))
	}))

	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	accountID := "123"
	siteV3Request := SiteV3Request{}
	siteV3Request.Name = "de3affdrere.inddcapcwafteam.net"

	siteV3Request.SiteType = "CLOUD_WAF"
	siteV3Response, diags := client.AddV3Site(&siteV3Request, accountID)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to add v3 site to Account ID: %s, %v\n", accountID, diags)
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to add v3 site to Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)
	}

	checkResponse(t, siteV3Response, siteV3Request, 51999737, 1717588301055, 462102065, "mhhp8q4.ng.impervadnsstage.net")
}

func checkResponse(t *testing.T, siteV3Response *SiteV3Response, siteV3Request SiteV3Request, AccountId int, CreationTime int64, Id int, Cname string) {
	if siteV3Response.Data[0].Name != siteV3Request.Name {
		t.Errorf("Should have  %s site name. Got: %s", siteV3Request.Name, siteV3Response.Data[0].Name)
	}

	if siteV3Response.Data[0].SiteType != siteV3Request.SiteType {
		t.Errorf("Should have  %s site type. Got: %s", siteV3Request.SiteType, siteV3Response.Data[0].SiteType)
	}

	if siteV3Response.Data[0].AccountId != AccountId {
		t.Errorf("Should have  %d site type. Got: %d", AccountId, siteV3Response.Data[0].AccountId)
	}

	if siteV3Response.Data[0].CreationTime != CreationTime {
		t.Errorf("Should have  %d site creation time. Got: %d", CreationTime, siteV3Response.Data[0].CreationTime)
	}

	if siteV3Response.Data[0].Id != Id {
		t.Errorf("Should have  %d site id time. Got: %d", Id, siteV3Response.Data[0].Id)
	}

	if siteV3Response.Data[0].Cname != Cname {
		t.Errorf("Should have  %s site cname. Got: %s", Cname, siteV3Response.Data[0].Cname)
	}
}

func TestSiteV3RequestMarshalsIsLoadBalancerSite(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_v3_test.TestSiteV3RequestMarshalsIsLoadBalancerSite")

	unset := SiteV3Request{Name: "example.com", SiteType: "PUBLIC_CLOUD"}
	unsetJSON, err := json.Marshal(unset)
	if err != nil {
		t.Fatalf("unexpected marshal error: %s", err)
	}
	if strings.Contains(string(unsetJSON), "isLoadBalancerSite") {
		t.Errorf("Unset is_load_balancer_site should be omitted. Got: %s", string(unsetJSON))
	}

	trueVal := true
	set := SiteV3Request{Name: "example.com", SiteType: "PUBLIC_CLOUD", CloudType: "GCP", IsLoadBalancerSite: &trueVal}
	setJSON, err := json.Marshal(set)
	if err != nil {
		t.Fatalf("unexpected marshal error: %s", err)
	}
	if !strings.Contains(string(setJSON), "\"isLoadBalancerSite\":true") {
		t.Errorf("Explicit true is_load_balancer_site should be sent as true. Got: %s", string(setJSON))
	}

	falseVal := false
	setFalse := SiteV3Request{Name: "example.com", IsLoadBalancerSite: &falseVal}
	setFalseJSON, err := json.Marshal(setFalse)
	if err != nil {
		t.Fatalf("unexpected marshal error: %s", err)
	}
	if !strings.Contains(string(setFalseJSON), "\"isLoadBalancerSite\":false") {
		t.Errorf("Explicit false is_load_balancer_site should be sent as false. Got: %s", string(setFalseJSON))
	}
}

func TestSiteV3ResponseUnmarshalsAuthorityHeader(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_v3_test.TestSiteV3ResponseUnmarshalsAuthorityHeader")

	body := "{\"data\":[{\"id\":123,\"name\":\"lb.example.com\",\"type\":\"PUBLIC_CLOUD\",\"cloud\":\"GCP\",\"isLoadBalancerSite\":true,\"authorityHeader\":\"abc-51999737.example.net\"}]}"
	var resp SiteV3Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("unexpected unmarshal error: %s", err)
	}
	if resp.Data[0].IsLoadBalancerSite == nil || !*resp.Data[0].IsLoadBalancerSite {
		t.Errorf("Expected isLoadBalancerSite=true. Got: %v", resp.Data[0].IsLoadBalancerSite)
	}
	if resp.Data[0].AuthorityHeader != "abc-51999737.example.net" {
		t.Errorf("Expected authorityHeader to be parsed. Got: %s", resp.Data[0].AuthorityHeader)
	}
}

func TestUpdateV3SiteWithName(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_v3_test.TestAddV3SiteWithName")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", endpointSiteV3+"/111?caid=123") {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpointSiteV3+"/111?caid=123", req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n  \"data\": [\n    {\n      \"id\": 462102065,\n      \"name\": \"de3affdrere.inddcapcwafteam.net\",\n      \"type\": \"CLOUD_WAF\",\n      \"accountId\": 51999737,\n      \"creationTime\": 1717588301055,\n      \"cname\": \"mhhp8q4.ng.impervadnsstage.net\"\n    }\n  ]\n}"))
	}))

	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	accountID := "123"
	siteV3Request := SiteV3Request{}
	siteV3Request.Name = "de3affdrere.inddcapcwafteam.net"
	siteV3Request.Id = 111
	siteV3Response, diags := client.UpdateV3Site(&siteV3Request, accountID)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to add v3 site to Account ID: %s, %v\n", accountID, diags)
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to add v3 site to Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)
	}
	siteV3Request.SiteType = "CLOUD_WAF"
	checkResponse(t, siteV3Response, siteV3Request, 51999737, 1717588301055, 462102065, "mhhp8q4.ng.impervadnsstage.net")
}
func TestDeleteV3Site(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_v3_test.TestAddV3SiteWithName")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s", endpointSiteV3+"/1234?caid=123") {
			t.Errorf("Should have hit %s endpoint. Got: %s", endpointSiteV3+"/1234?caid=123", req.URL.String())
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\n  \"data\": [\n    {\n      \"id\": 462102065,\n      \"name\": \"de3affdrere.inddcapcwafteam.net\",\n      \"type\": \"CLOUD_WAF\",\n      \"accountId\": 51999737,\n      \"creationTime\": 1717588301055,\n      \"cname\": \"mhhp8q4.ng.impervadnsstage.net\"\n    }\n  ]\n}"))
	}))

	defer server.Close()
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	accountID := "123"
	siteV3Request := SiteV3Request{}
	siteV3Request.Name = "de3affdrere.inddcapcwafteam.net"
	siteV3Request.SiteType = "CLOUD_WAF"
	siteV3Request.Id = 1234

	siteV3Response, diags := client.DeleteV3Site(&siteV3Request, accountID)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] failed to delete v3 site of Account ID: %s, %v\n", accountID, diags)
	} else if siteV3Response.Errors != nil {
		log.Printf("[ERROR] Failed to delete v3 site of Account ID: %s, %v\n", accountID, siteV3Response.Errors[0].Detail)

		checkResponse(t, siteV3Response, siteV3Request, 51999737, 1717588301055, 462102065, "mhhp8q4.ng.impervadnsstage.net")

	}
}
