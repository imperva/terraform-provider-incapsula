package incapsula

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSiteDomains(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_test.TestGetSiteDomains")

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if !strings.HasPrefix(req.URL.Path, endpointSiteDomainsV3+"999/domains") {
			t.Errorf("Should have hit the v3 site domains endpoint. Got path: %s", req.URL.Path)
		}
		if req.URL.Query().Get("excludeAutoDiscovered") != "false" {
			t.Errorf("Should pass excludeAutoDiscovered=false. Got: %s", req.URL.Query().Get("excludeAutoDiscovered"))
		}
		if req.URL.Query().Get("caid") != "123" {
			t.Errorf("Should pass caid=123. Got: %s", req.URL.Query().Get("caid"))
		}
		if req.URL.Query().Get("pageSize") != "100" || req.URL.Query().Get("pageNumber") != "0" {
			t.Errorf("Should request pageSize=100&pageNumber=0. Got: %s", req.URL.RawQuery)
		}
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[{\"id\":1,\"siteId\":999,\"domain\":\"a.example.com\",\"status\":\"PROTECTED\",\"cname\":\"abc.impervadns.net\",\"creationDate\":1727654400000,\"aRecords\":[\"1.2.3.4\"]},{\"id\":2,\"siteId\":999,\"domain\":\"b.example.com\",\"cname\":\"abc.impervadns.net\",\"creationDate\":1727654400001}],\"meta\":{\"page\":0,\"size\":100,\"totalElements\":2,\"totalPages\":1}}"))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetSiteDomains("999", "123")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("Expected 2 domains. Got: %d", len(resp.Data))
	}
	first := resp.Data[0]
	if first.Domain != "a.example.com" || first.Status != "PROTECTED" || first.Cname != "abc.impervadns.net" || first.CreationDate != 1727654400000 || len(first.ARecords) != 1 {
		t.Errorf("Unexpected first domain: %+v", first)
	}
	if resp.Data[1].Status != "" {
		t.Errorf("Absent status should decode as empty. Got: %+v", resp.Data[1])
	}
}

func TestGetSiteDomainsFollowsPagination(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_test.TestGetSiteDomainsFollowsPagination")

	var requestedPages []string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("pageNumber")
		requestedPages = append(requestedPages, page)
		rw.WriteHeader(200)
		switch page {
		case "0":
			rw.Write([]byte("{\"data\":[{\"id\":1,\"siteId\":999,\"domain\":\"a.example.com\"}],\"meta\":{\"page\":0,\"size\":1,\"totalElements\":2,\"totalPages\":2}}"))
		case "1":
			rw.Write([]byte("{\"data\":[{\"id\":2,\"siteId\":999,\"domain\":\"b.example.com\"}],\"meta\":{\"page\":1,\"size\":1,\"totalElements\":2,\"totalPages\":2}}"))
		default:
			t.Errorf("Unexpected page requested: %s", page)
			rw.Write([]byte("{\"data\":[],\"meta\":{\"page\":2,\"size\":1,\"totalElements\":2,\"totalPages\":2}}"))
		}
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetSiteDomains("999", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(requestedPages) != 2 {
		t.Errorf("Expected 2 page requests. Got: %v", requestedPages)
	}
	if len(resp.Data) != 2 || resp.Data[0].Domain != "a.example.com" || resp.Data[1].Domain != "b.example.com" {
		t.Errorf("Expected domains from both pages. Got: %+v", resp.Data)
	}
}

func TestGetSiteDomainsStopsOnEmptyPage(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_site_domains_test.TestGetSiteDomainsStopsOnEmptyPage")

	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		calls++
		rw.WriteHeader(200)
		rw.Write([]byte("{\"data\":[],\"meta\":{\"page\":0,\"size\":100,\"totalElements\":0,\"totalPages\":5}}"))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	resp, err := client.GetSiteDomains("999", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if calls != 1 || len(resp.Data) != 0 {
		t.Errorf("Expected a single request and no domains. Got calls=%d, data=%+v", calls, resp.Data)
	}
}
