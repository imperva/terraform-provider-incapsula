package incapsula

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMissingCredentials(t *testing.T) {
	config := Config{}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingAPIIDMessage {
		t.Errorf("Should have received missing API ID message, got: %s", err)
	}
}

func TestMissingAPIID(t *testing.T) {
	config := Config{APIID: "", APIKey: "foo"}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingAPIIDMessage {
		t.Errorf("Should have received missing API ID message, got: %s", err)
	}
}

func TestMissingAPIKey(t *testing.T) {
	config := Config{APIID: "foo", APIKey: ""}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingAPIKeyMessage {
		t.Errorf("Should have received missing API key message, got: %s", err)
	}
}

func TestMissingBaseURL(t *testing.T) {
	config := Config{APIID: "foo", APIKey: "bar", BaseURL: ""}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingBaseURLMessage {
		t.Errorf("Should have received missing base URL message, got: %s", err)
	}
}

func TestMissingBaseURLRev2(t *testing.T) {
	config := Config{APIID: "foo", APIKey: "bar", BaseURL: "foobar.com", BaseURLRev2: "", BaseURLRev3: "foobar.com"}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingBaseURLRev2Message {
		t.Errorf("Should have received missing Base URL Revision 3 message, got: %s", err)
	}
}

func TestMissingBaseURLRev3(t *testing.T) {
	config := Config{APIID: "foo", APIKey: "bar", BaseURL: "foobar.com", BaseURLRev2: "foobar.com", BaseURLRev3: ""}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingBaseURLRev3Message {
		t.Errorf("Should have received missing Base URL Revision 2 message, got: %s", err)
	}
}

func TestMissingBaseURLAPI(t *testing.T) {
	config := Config{APIID: "foo", APIKey: "bar", BaseURL: "foobar.com", BaseURLRev2: "foobar.com", BaseURLRev3: "foobar.com", BaseURLAPI: ""}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if err.Error() != missingBaseURLAPIMessage {
		t.Errorf("Should have received missing Base URL API message, got: %s", err)
	}
}

func TestInvalidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/account" {
			t.Errorf("Should have have hit /account endpoint. Got: %s", req.URL.String())
		}
		rw.Write([]byte(`{"res":1,"res_message":"fail"}`))
	}))
	defer server.Close()

	config := Config{APIID: "bad", APIKey: "bad", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLRev3: server.URL, BaseURLAPI: server.URL, CustomTestDomain: ".example.com"}
	client, err := config.Client()
	if err == nil {
		t.Errorf("Should have received an error, got a client: %q", client)
	}
	if !strings.HasPrefix(err.Error(), "Error from Incapsula service when checking account") {
		t.Errorf("Should have received Incapsula service error, got: %s", err)
	}
}

func TestValidCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() != "/account" {
			t.Errorf("Should have have hit /account endpoint. Got: %s", req.URL.String())
		}
		rw.Write([]byte(`{"res":0,"res_message":"OK"}`))
	}))
	defer server.Close()

	config := Config{APIID: "good", APIKey: "good", BaseURL: server.URL, BaseURLRev2: server.URL, BaseURLRev3: server.URL, BaseURLAPI: server.URL, CustomTestDomain: ".example.com"}
	client, err := config.Client()
	if err != nil {
		t.Errorf("Should not have received an error, got: %s", err)
	}
	if client == nil {
		t.Error("Client should not be nil")
	}
}
