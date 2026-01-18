package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientAddApiClientBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}

	email := "example@example.com"

	request := &APIClientUpdateRequest{}
	request.Name = "test"

	apiClientResponse, err := client.CreateAPIClient(1234, email, request)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error creating api_client")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiClientResponse != nil {
		t.Errorf("Should have received a nil apiClientResponse instance")
	}
}

func TestClientAddApiClientBadJSON(t *testing.T) {
	accountID := 123
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	email := "example@example.com"
	request := &APIClientUpdateRequest{}
	request.Name = "test"

	apiClientResponse, err := client.CreateAPIClient(accountID, email, request)
	if err == nil {
		t.Errorf("Should have received an error")
	}

	expectedErrorFragment := "Error status code 200 from Incapsula service"

	if !strings.Contains(err.Error(), expectedErrorFragment) {
		t.Errorf("Should have received response parse error, got: %s", err)
	}

	if apiClientResponse != nil {
		t.Errorf("Should have received a nil apiClientResponse instance")
	}
}

func TestClientApiClientStatusBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	apiClientResponse, err := client.GetAPIClient(accountID, "1234")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error getting api_client with id")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiClientResponse != nil {
		t.Errorf("Should have received a nil apiClientResponse instance")
	}
}

func TestClientApiClientStatusBadJSON(t *testing.T) {
	accountID := 123
	clientID := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s?caid=%d&id=%s", endpointAPIClient, accountID, clientID) {
			t.Errorf("Should have have hit %s?caid=%d&id=%s endpoint. Got: %s", endpointAPIClient, accountID, clientID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	apiClientResponse, err := client.GetAPIClient(accountID, "1234")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing api_client status JSON response for api_client id %s", clientID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiClientResponse != nil {
		t.Errorf("Should have received a nil apiClientResponse instance")
	}
}

func TestClientDeleteApiClientBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	err := client.DeleteAPIClient(accountID, "1234")
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error from Incapsula service when deleting api-client: %s", "1234")) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
}

func TestClientUpdateApiClientBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	accountID := 123
	clientID := "1234"
	name := "testest"
	request := &APIClientUpdateRequest{}
	request.Name = name
	apiClientResponse, err := client.PatchAPIClient(accountID, clientID, request)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error updating api_client with Id %s", clientID)) {
		t.Errorf("Should have received an client error, got: %s", err)
	}
	if apiClientResponse != nil {
		t.Errorf("Should have received a nil apiClientResponse instance")
	}
}

func TestClientUpdateApiClientBadJSON(t *testing.T) {
	accountID := 123
	clientID := "1234"

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != fmt.Sprintf("%s/%s?caid=%d", endpointAPIClient, clientID, accountID) {
			t.Errorf("Should have have hit %s/%s?caid=%d endpoint. Got: %s", endpointAPIClient, clientID, accountID, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	request := &APIClientUpdateRequest{}
	request.Name = "test"
	config := &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}
	apiClientResponse, err := client.PatchAPIClient(accountID, clientID, request)
	if err == nil {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(err.Error(), fmt.Sprintf("Error parsing update api_client JSON response for id %s", clientID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", err)
	}
	if apiClientResponse != nil {
		t.Errorf("Should have received a nil apiClientResponse instance")
	}
}
