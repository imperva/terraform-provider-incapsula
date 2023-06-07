package incapsula

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestClientAbpWebsitesBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountId := 1234
	abpWebsitesResponse, diags := client.ReadAbpWebsites(accountId)
	if len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when reading ABP Websites for Account ID %d", accountId)) {
		t.Errorf("Should have received a client error, got: %+v", diags)
	}
	if abpWebsitesResponse != nil {
		t.Errorf("Should have received a nil abpWebsitesResponse instance")
	}
}

func TestClientAbpWebsitesBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountId := 1234
	abpWebsitesResponse, diags := client.ReadAbpWebsites(accountId)
	if len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing ABP Websites JSON response for Account ID %d", accountId)) {
		t.Errorf("Should have received a client error, got: %+v", diags)
	}
	if abpWebsitesResponse != nil {
		t.Errorf("Should have received a nil abpWebsitesResponse instance")
	}
}

func TestClientAbpWebsitesBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`some error`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountId := 1234
	abpWebsitesResponse, diags := client.ReadAbpWebsites(accountId)
	if len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code 400 from Incapsula service when reading ABP Websites for Account ID %d: some error", accountId)) {
		t.Errorf("Should have received a client error, got: %+v", diags)
	}
	if abpWebsitesResponse != nil {
		t.Errorf("Should have received a nil abpWebsitesResponse instance")
	}
}

func TestClientAbpWebsitesOk(t *testing.T) {
	id1 := "id1"
	id2 := "id2"
	id3 := "id3"
	nameid1 := "nameid1"
	abpWebsites := AbpTerraformAccount{
		AutoPublish: true,
		WebsiteGroups: []AbpTerraformWebsiteGroup{
			{
				Id:     &id1,
				NameId: &nameid1,
				Name:   "name1",
				Websites: []AbpTerraformWebsite{{
					Id:               &id2,
					WebsiteId:        1,
					EnableMitigation: true,
				},
					{
						Id:               &id3,
						WebsiteId:        4,
						EnableMitigation: false,
					},
				},
			},
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		b, _ := json.Marshal(abpWebsites)
		rw.Write(b)
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountId := 1234
	abpWebsitesResponse, diags := client.ReadAbpWebsites(accountId)
	if len(diags) != 0 {
		t.Errorf("Should not have received an error %+v", diags)
	}

	if !reflect.DeepEqual(*abpWebsitesResponse, abpWebsites) {
		t.Errorf("Unexpected abpWebsitesResponse: %+v", abpWebsitesResponse)
	}
}

/*


func TestClientAbpWebsitesValidRule(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test client_certificate_test.TestClientAbpWebsitesValidRule")
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
	abpWebsitesResponse, err := client.AbpWebsites(siteID, "", "", "", "")
	if err != nil {
		t.Errorf("Should not have received an error")
	}
	if abpWebsitesResponse == nil {
		t.Errorf("Should not have received a nil abpWebsitesResponse instance")
	}
	if abpWebsitesResponse.Res != 0 {
		t.Errorf("Response code doesn't match")
	}
}
*/
