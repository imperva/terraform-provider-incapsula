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

func TestClientAbpWebsitesReadBadConnection(t *testing.T) {
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

func TestClientAbpWebsitesReadBadJSON(t *testing.T) {
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

func TestClientAbpWebsitesReadBadRequest(t *testing.T) {
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

func TestClientAbpWebsitesReadOk(t *testing.T) {
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

func TestClientAbpWebsitesCreateBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountId := 1234
	abpWebsitesResponse, diags := client.CreateAbpWebsites(accountId, AbpTerraformAccount{})
	if len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when creating ABP Websites for Account ID %d", accountId)) {
		t.Errorf("Should have received a client error, got: %+v", diags)
	}
	if abpWebsitesResponse != nil {
		t.Errorf("Should have received a nil abpWebsitesResponse instance")
	}
}

func TestClientAbpWebsitesCreateBadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountId := 1234
	abpWebsitesResponse, diags := client.CreateAbpWebsites(accountId, AbpTerraformAccount{})
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

func TestClientAbpWebsitesCreateBadRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`some error`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountId := 1234
	abpWebsitesResponse, diags := client.CreateAbpWebsites(accountId, AbpTerraformAccount{})
	if len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code 400 from Incapsula service when creating ABP Websites for Account ID %d: some error", accountId)) {
		t.Errorf("Should have received a client error, got: %+v", diags)
	}
	if abpWebsitesResponse != nil {
		t.Errorf("Should have received a nil abpWebsitesResponse instance")
	}
}

func TestClientAbpWebsitesCreateOk(t *testing.T) {
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
		rw.WriteHeader(http.StatusCreated)
		rw.Write(b)
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	accountId := 1234
	abpWebsitesResponse, diags := client.CreateAbpWebsites(accountId, abpWebsites)
	if len(diags) != 0 {
		t.Errorf("Should not have received an error %+v", diags)
		return
	}

	if !reflect.DeepEqual(*abpWebsitesResponse, abpWebsites) {
		t.Errorf("Unexpected abpWebsitesResponse: %+v", abpWebsitesResponse)
	}
}
