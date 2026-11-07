package incapsula

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Extra helpers for the preflight → publish pipeline and the Account data
// source, on top of what the shared fake offers.

// abpPreflightSeedSnapshot stores a snapshot that was taken outside Terraform,
// so a publish can be exercised on its own. It takes part in publish ordering
// like any other snapshot.
func (f *abpFake) abpPreflightSeedSnapshot() string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return abpFakeString(f.newPreflight(abpFakeAccountId)["id"])
}

// abpAccountBrokenClient returns a client talking to a service that answers
// every request with the given status and body, for the failure paths the
// shared fake never takes.
func abpAccountBrokenClient(t *testing.T, status int, body string) *Client {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return &Client{
		config:     &Config{APIID: "fake-id", APIKey: "fake-key", BaseURLAPI: server.URL},
		httpClient: &http.Client{},
	}
}
