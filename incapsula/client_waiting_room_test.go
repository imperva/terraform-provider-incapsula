package incapsula

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// Create Waiting Room tests

func TestClientCreateWaitingRoomBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when creating Waiting Room for Site ID %s", siteID)) {
		t.Errorf("Should have received a client error, got: %s", diags[0].Detail)
	}
	if createWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil createWaitingRoomResponse instance")
	}
}

func TestClientCreateWaitingRoomBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "POST" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "POST", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if createWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil createWaitingRoomResponse instance")
	}
}

func TestClientCreateWaitingRoomBadStatusCodeWithEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "POST" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "POST", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when creating Waiting Room for Site ID %s", 401, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if createWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil createWaitingRoomResponse instance")
	}
}

func TestClientCreateWaitingRoomBadStatusCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "POST" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "POST", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"errors": [{"status":404,"message":"not found"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when creating Waiting Room for Site ID %s", 404, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if createWaitingRoomResponse == nil || createWaitingRoomResponse.Errors[0].Status != 404 {
		t.Errorf("Should have received an error DTO")
	}
}

func TestClientCreateWaitingRoomEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "POST" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "POST", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if createWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil createWaitingRoomResponse instance")
	}
}

func TestClientCreateWaitingRoomInvalidJSONValue(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "POST" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "POST", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(`{"data": [{"name":"waiting room 1","description":"","enabled":true,"htmlTemplateBase64":"","filter":"","botsActionInQueuingMode":"","queueInactivityTimeout":5,"isEntranceRateEnabled":true,"entranceRateThreshold":0,"isConcurrentSessionsEnabled":5,"concurrentSessionsThreshold":0,"inactivityTimeout":0}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room JSON response for Site ID %s", siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if createWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil createWaitingRoomResponse instance")
	}
}

func TestClientCreateWaitingRoomValidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms", siteID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "POST" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "POST", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(`{"data": [{"id": 1, "accountId":1234, "name":"waiting room 1","description":"","enabled":true,"htmlTemplateBase64":"","filter":"","botsActionInQueuingMode":"WAIT_IN_LINE","queueInactivityTimeout":5,"isEntranceRateEnabled":true,"entranceRateThreshold":500,"isConcurrentSessionsEnabled":false,"concurrentSessionsThreshold":0,"inactivityTimeout":5}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
		EntranceRateThreshold:  500,
	}

	createWaitingRoomResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags != nil {
		t.Errorf("Should not have received an error")
	}
	if len(createWaitingRoomResponse.Data) != 1 {
		t.Errorf("Waiting Rooms list size doesn't match")
	}
	if createWaitingRoomResponse.Data[0].Id != 1 {
		t.Errorf("Waiting Room ID doesn't match")
	}
	if createWaitingRoomResponse.Data[0].AccountId != 1234 {
		t.Errorf("Account ID doesn't match")
	}
	if !createWaitingRoomResponse.Data[0].Enabled {
		t.Errorf("Waiting Room should be enabled")
	}
	if createWaitingRoomResponse.Data[0].Name != "waiting room 1" {
		t.Errorf("Waiting Room name doesn't match")
	}
	if !createWaitingRoomResponse.Data[0].EntranceRateEnabled || createWaitingRoomResponse.Data[0].ConcurrentSessionsEnabled || createWaitingRoomResponse.Data[0].EntranceRateThreshold != 500 || createWaitingRoomResponse.Data[0].ConcurrentSessionsThreshold != 0 || createWaitingRoomResponse.Data[0].InactivityTimeout != 5 {
		t.Errorf("Thresholds don't match")
	}
}

// Read Waiting Room tests

func TestClientReadWaitingRoomBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	waitingRoomID := int64(1)

	readWaitingRoomResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when reading Waiting Room %d for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", diags[0].Detail)
	}
	if readWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil readWaitingRoomResponse instance")
	}
}

func TestClientReadWaitingRoomBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "GET" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "GET", endpoint, req.Method, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readWaitingRoomResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if readWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil readWaitingRoomResponse instance")
	}
}

func TestClientReadWaitingRoomBadStatusCodeWithEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "GET" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "GET", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readWaitingRoomResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when reading Waiting Room %d for Site ID %s", 404, waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if readWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil readWaitingRoomResponse instance")
	}
}

func TestClientReadWaitingRoomBadStatusCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "GET" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "GET", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"errors": [{"status":404,"message":"not found"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readWaitingRoomResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when reading Waiting Room %d for Site ID %s", 404, waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if readWaitingRoomResponse == nil || readWaitingRoomResponse.Errors[0].Status != 404 {
		t.Errorf("Should have received an error DTO")
	}
}

func TestClientReadWaitingRoomEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "GET" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "GET", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readWaitingRoomResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if readWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil readWaitingRoomResponse instance")
	}
}

func TestClientReadWaitingRoomValidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "GET" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "GET", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"data": [{"id": 1, "accountId":1234, "name":"waiting room 1","description":"","enabled":true,"htmlTemplateBase64":"","filter":"","botsActionInQueuingMode":"WAIT_IN_LINE","queueInactivityTimeout":5,"isEntranceRateEnabled":true,"entranceRateThreshold":500,"concurrentSessionsThreshold":0,"inactivityTimeout":5}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	readWaitingRoomResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if diags != nil {
		t.Errorf("Should not have received an error")
	}
	if len(readWaitingRoomResponse.Data) != 1 {
		t.Errorf("Waiting Rooms list size doesn't match")
	}
	if readWaitingRoomResponse.Data[0].Id != 1 {
		t.Errorf("Waiting Room ID doesn't match")
	}
	if readWaitingRoomResponse.Data[0].AccountId != 1234 {
		t.Errorf("Account ID doesn't match")
	}
	if !readWaitingRoomResponse.Data[0].Enabled {
		t.Errorf("Waiting Room should be enabled")
	}
	if readWaitingRoomResponse.Data[0].Name != "waiting room 1" {
		t.Errorf("Waiting Room name doesn't match")
	}
	if !readWaitingRoomResponse.Data[0].EntranceRateEnabled || readWaitingRoomResponse.Data[0].ConcurrentSessionsEnabled || readWaitingRoomResponse.Data[0].EntranceRateThreshold != 500 || readWaitingRoomResponse.Data[0].ConcurrentSessionsThreshold != 0 || readWaitingRoomResponse.Data[0].InactivityTimeout != 5 {
		t.Errorf("Thresholds don't match")
	}
}

// Update Waiting Room tests

func TestClientUpdateWaitingRoomBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	waitingRoomID := int64(1)

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	updateWaitingRoomResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when updating Waiting Room %d for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", diags[0].Detail)
	}
	if updateWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil updateWaitingRoomResponse instance")
	}
}

func TestClientUpdateWaitingRoomBadJSON(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "PUT" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "PUT", endpoint, req.Method, req.URL.String())
		}
		rw.Write([]byte(`{`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	updateWaitingRoomResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if updateWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil updateWaitingRoomResponse instance")
	}
}

func TestClientUpdateWaitingRoomBadStatusCodeWithEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "PUT" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "PUT", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	updateWaitingRoomResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when updating Waiting Room %d for Site ID %s", 404, waitingRoomID, siteID)) {
		t.Errorf("Should have received a Status Code error, got: %s", diags[0].Detail)
	}
	if updateWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil updateWaitingRoomResponse instance")
	}
}

func TestClientUpdateWaitingRoomBadStatusCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "PUT" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "PUT", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"errors": [{"status":404,"message":"not found"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	updateWaitingRoomResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when updating Waiting Room %d for Site ID %s", 404, waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if updateWaitingRoomResponse == nil || updateWaitingRoomResponse.Errors[0].Status != 404 {
		t.Errorf("Should have received an error DTO")
	}
}

func TestClientUpdateWaitingRoomEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "PUT" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "PUT", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
	}

	updateWaitingRoomResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error parsing Waiting Room %d JSON response for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a JSON parse error, got: %s", diags[0].Detail)
	}
	if updateWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil updateWaitingRoomResponse instance")
	}
}

func TestClientUpdateWaitingRoomValidResponse(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "PUT" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "PUT", endpoint, req.Method, req.URL.String())
		}
		rw.Write([]byte(`{"data": [{"id": 1, "accountId":1234, "name":"waiting room 1","description":"","enabled":true,"htmlTemplateBase64":"","filter":"","botsActionInQueuingMode":"WAIT_IN_LINE","queueInactivityTimeout":5,"isEntranceRateEnabled":true,"entranceRateThreshold":500,"isConcurrentSessionsEnabled":false,"concurrentSessionsThreshold":0,"inactivityTimeout":5}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	waitingRoom := WaitingRoomDTO{
		Name:                   "waiting room 1",
		Enabled:                true,
		QueueInactivityTimeout: 5,
		EntranceRateEnabled:    true,
		EntranceRateThreshold:  500,
	}

	updateWaitingRoomResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if diags != nil {
		t.Errorf("Should not have received an error")
	}
	if len(updateWaitingRoomResponse.Data) != 1 {
		t.Errorf("Waiting Rooms list size doesn't match")
	}
	if updateWaitingRoomResponse.Data[0].Id != 1 {
		t.Errorf("Waiting Room ID doesn't match")
	}
	if updateWaitingRoomResponse.Data[0].AccountId != 1234 {
		t.Errorf("Account ID doesn't match")
	}
	if !updateWaitingRoomResponse.Data[0].Enabled {
		t.Errorf("Waiting Room should be enabled")
	}
	if updateWaitingRoomResponse.Data[0].Name != "waiting room 1" {
		t.Errorf("Waiting Room name doesn't match")
	}
	if !updateWaitingRoomResponse.Data[0].EntranceRateEnabled || updateWaitingRoomResponse.Data[0].ConcurrentSessionsEnabled || updateWaitingRoomResponse.Data[0].EntranceRateThreshold != 500 || updateWaitingRoomResponse.Data[0].ConcurrentSessionsThreshold != 0 || updateWaitingRoomResponse.Data[0].InactivityTimeout != 5 {
		t.Errorf("Thresholds don't match")
	}
}

// Delete Waiting Room tests

func TestClientDeleteWaitingRoomBadConnection(t *testing.T) {
	config := &Config{APIID: "foo", APIKey: "bar", BaseURL: "badness.incapsula.com", BaseURLRev2: "badness.incapsula.com", BaseURLAPI: "badness.incapsula.com"}
	client := &Client{config: config, httpClient: &http.Client{Timeout: time.Millisecond * 1}}
	siteID := "42"
	waitingRoomID := int64(1)

	deleteWaitingRoomResponse, diags := client.DeleteWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received a error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error from Incapsula service when deleting Waiting Room %d for Site ID %s", waitingRoomID, siteID)) {
		t.Errorf("Should have received a client error, got: %s", diags[0].Detail)
	}
	if deleteWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil deleteWaitingRoomResponse instance")
	}
}

func TestClientDeleteWaitingRoomBadStatusCodeWithEmptyBody(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "DELETE" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "DELETE", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	deleteWaitingRoomResponse, diags := client.DeleteWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when deleting Waiting Room %d for Site ID %s", 404, waitingRoomID, siteID)) {
		t.Errorf("Should have received a Status Code error, got: %s", diags[0].Detail)
	}
	if deleteWaitingRoomResponse != nil {
		t.Errorf("Should have received a nil deleteWaitingRoomResponse instance")
	}
}

func TestClientDeleteWaitingRoomBadStatusCode(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "DELETE" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "DELETE", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(`{"errors": [{"status":404,"message":"not found"}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	deleteWaitingRoomResponse, diags := client.DeleteWaitingRoom(siteID, waitingRoomID)
	if diags == nil || len(diags) == 0 {
		t.Errorf("Should have received an error")
	}
	if !strings.HasPrefix(diags[0].Detail, fmt.Sprintf("Error status code %d from Incapsula service when deleting Waiting Room %d for Site ID %s", 404, waitingRoomID, siteID)) {
		t.Errorf("Should have received a Status Code error, got: %s", diags[0].Detail)
	}
	if deleteWaitingRoomResponse == nil || deleteWaitingRoomResponse.Errors[0].Status != 404 {
		t.Errorf("Should have received an error DTO")
	}
}

func TestClientDeleteWaitingRoomValidWaitingRoom(t *testing.T) {
	apiID := "foo"
	apiKey := "bar"
	siteID := "42"
	waitingRoomID := int64(1)

	endpoint := fmt.Sprintf("/waiting-room-settings/v3/sites/%s/waiting-rooms/%d", siteID, waitingRoomID)

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != endpoint || req.Method != "DELETE" {
			t.Errorf("Should have have hit %s %s endpoint. Got: %s %s", "DELETE", endpoint, req.Method, req.URL.String())
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{"data": [{"id": 1, "accountId":1234, "name":"waiting room 1","description":"","enabled":true,"htmlTemplateBase64":"","filter":"","botsActionInQueuingMode":"WAIT_IN_LINE","queueInactivityTimeout":5,"isEntranceRateEnabled":true,"entranceRateThreshold":500,"concurrentSessionsThreshold":0,"inactivityTimeout":5}]}`))
	}))
	defer server.Close()

	config := &Config{APIID: apiID, APIKey: apiKey, BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLAPI: server.URL}
	client := &Client{config: config, httpClient: &http.Client{}}

	deleteWaitingRoomResponse, diags := client.DeleteWaitingRoom(siteID, waitingRoomID)
	if diags != nil {
		t.Errorf("Should not have received an error")
	}
	if deleteWaitingRoomResponse == nil || len(deleteWaitingRoomResponse.Data) == 0 {
		t.Errorf("Should have recived a response")
	}
}
