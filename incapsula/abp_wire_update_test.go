package incapsula

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeAbpFullReplaceEndpoint serves one object: a GET returns `current`, and a
// PUT records the body it was given. That is all the update path needs in
// order to be observable, and the body is where the regression lives -- a
// field missing from it is a field the next full replace clears.
type fakeAbpFullReplaceEndpoint struct {
	current    string
	putBody    map[string]any
	rawPutBody string
	gets       int
}

func (f *fakeAbpFullReplaceEndpoint) client(t *testing.T) *Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			f.gets++
			w.Write([]byte(f.current))
		case http.MethodPut:
			if err := json.NewDecoder(r.Body).Decode(&f.putBody); err != nil {
				t.Errorf("the update request body does not parse: %v", err)
			}
			raw, _ := json.Marshal(f.putBody)
			f.rawPutBody = string(raw)
			w.Write([]byte(f.current))
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	t.Cleanup(server.Close)

	return &Client{
		config:     &Config{APIID: "id", APIKey: "key", BaseURLAPI: server.URL},
		httpClient: &http.Client{},
	}
}

// TestAbpUpdatePreservesUnmodelledFieldsEndToEnd runs every full-replace
// update endpoint against an object that carries a field the provider does not
// model, and checks that the field is still in the request body. Each case
// fails on the pre-merge code, where the body was built from the provider's
// struct alone.
func TestAbpUpdatePreservesUnmodelledFieldsEndToEnd(t *testing.T) {
	cases := []struct {
		name string
		// current is what a read of the object returns.
		current string
		// update performs the provider's update of that object.
		update func(c *Client) error
		// unmodelled is the path that must survive, as a chain of object keys
		// (an element of a list is addressed by its position).
		unmodelled []string
		want       any
	}{
		{
			name: "Site keeps a field the Site model lacks",
			current: `{
				"id": "11111111-1111-1111-1111-111111111111",
				"name": "site",
				"mx_hostname_id": "mx-7",
				"selectors": [{"id": "sel-1", "criteria": {"path_prefix": "/"}, "analysis_settings": {"rate_limiting": "none"}}]
			}`,
			update: func(c *Client) error {
				_, err := c.UpdateAbpSite("11111111-1111-1111-1111-111111111111", AbpSite{Name: "renamed"})
				return err
			},
			unmodelled: []string{"mx_hostname_id"},
			want:       "mx-7",
		},
		{
			name: "Site keeps a field a Selector model lacks",
			current: `{
				"id": "11111111-1111-1111-1111-111111111111",
				"name": "site",
				"selectors": [{
					"id": "sel-1",
					"criteria": {"path_prefix": "/login"},
					"analysis_settings": {"rate_limiting": "none"},
					"future_selector_field": "keep me"
				}]
			}`,
			update: func(c *Client) error {
				prefix := "/login"
				_, err := c.UpdateAbpSite("11111111-1111-1111-1111-111111111111", AbpSite{
					Name: "site",
					Selectors: []AbpSelector{{
						Id:               "sel-1",
						Criteria:         AbpSelectorCriteria{PathPrefix: &prefix},
						AnalysisSettings: AbpAnalysisSettings{RateLimiting: AbpRateLimiting{Mode: AbpRateLimitingModeNone}},
					}},
				})
				return err
			},
			unmodelled: []string{"selectors", "0", "future_selector_field"},
			want:       "keep me",
		},
		{
			name: "Domain keeps a field the Domain update model lacks",
			current: `{
				"id": "22222222-2222-2222-2222-222222222222",
				"site_id": "33333333-3333-3333-3333-333333333333",
				"criteria": {"exact": "example.com"},
				"cookiescope": "domain",
				"future_domain_field": "keep me"
			}`,
			update: func(c *Client) error {
				_, err := c.UpdateAbpDomain("22222222-2222-2222-2222-222222222222", AbpDomain{
					SiteId:      "33333333-3333-3333-3333-333333333333",
					Cookiescope: "domain",
				})
				return err
			},
			unmodelled: []string{"future_domain_field"},
			want:       "keep me",
		},
		{
			name: "Policy keeps a field the Policy model lacks",
			current: `{
				"id": "44444444-4444-4444-4444-444444444444",
				"name": "policy",
				"directives": [{"action": "block", "condition_id": "55555555-5555-5555-5555-555555555555"}],
				"future_policy_field": "keep me"
			}`,
			update: func(c *Client) error {
				conditionId := "55555555-5555-5555-5555-555555555555"
				_, err := c.UpdateAbpPolicy("44444444-4444-4444-4444-444444444444", AbpPolicy{
					Name:       "policy",
					Directives: []AbpDirective{{Action: "block", ConditionId: &conditionId}},
				})
				return err
			},
			unmodelled: []string{"future_policy_field"},
			want:       "keep me",
		},
		{
			name: "Condition keeps a field the Condition model lacks",
			current: `{
				"id": "66666666-6666-6666-6666-666666666666",
				"name": "condition",
				"description": "",
				"code": "true",
				"owner": "account",
				"future_condition_field": "keep me"
			}`,
			update: func(c *Client) error {
				_, err := c.UpdateAbpCondition("66666666-6666-6666-6666-666666666666", AbpCondition{
					Kind: AbpConditionKindLiteral,
					Name: "condition",
					Code: "true",
				})
				return err
			},
			unmodelled: []string{"future_condition_field"},
			want:       "keep me",
		},
		{
			name: "Proof of Work Configuration keeps a field its model lacks",
			current: `{
				"id": "77777777-7777-7777-7777-777777777777",
				"name": "pow",
				"difficulty": 5,
				"future_pow_field": "keep me"
			}`,
			update: func(c *Client) error {
				_, err := c.UpdateAbpProofOfWorkConfiguration("77777777-7777-7777-7777-777777777777", AbpProofOfWorkConfiguration{
					Name:       "pow",
					Difficulty: 7,
				})
				return err
			},
			unmodelled: []string{"future_pow_field"},
			want:       "keep me",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			endpoint := &fakeAbpFullReplaceEndpoint{current: test.current}
			if err := test.update(endpoint.client(t)); err != nil {
				t.Fatalf("updating: %v", err)
			}
			if endpoint.gets == 0 {
				t.Error("the update did not read the current object first")
			}
			if got := lookupJsonPath(t, endpoint.putBody, test.unmodelled); got != test.want {
				t.Errorf("%v in the update request body = %#v, want %#v\nbody: %s",
					test.unmodelled, got, test.want, endpoint.rawPutBody)
			}
		})
	}
}

// TestAbpUpdateSendsModelledFieldsFromThePayload is the other half: the merge
// must not resurrect a value the provider deliberately left out, or clearing
// an optional would stop working.
func TestAbpUpdateSendsModelledFieldsFromThePayload(t *testing.T) {
	endpoint := &fakeAbpFullReplaceEndpoint{current: `{
		"id": "11111111-1111-1111-1111-111111111111",
		"name": "site",
		"industry": "travel",
		"selectors": [{"id": "sel-1", "criteria": {"path_prefix": "/login", "case_insensitive": true}, "purpose": "login"}]
	}`}

	prefix := "/login"
	if _, err := endpoint.client(t).UpdateAbpSite("11111111-1111-1111-1111-111111111111", AbpSite{
		Name: "site",
		Selectors: []AbpSelector{{
			Id:               "sel-1",
			Criteria:         AbpSelectorCriteria{PathPrefix: &prefix},
			AnalysisSettings: AbpAnalysisSettings{RateLimiting: AbpRateLimiting{Mode: AbpRateLimitingModeNone}},
		}},
	}); err != nil {
		t.Fatalf("updating: %v", err)
	}

	if _, ok := endpoint.putBody["industry"]; ok {
		t.Errorf("industry should have been cleared, body: %s", endpoint.rawPutBody)
	}
	selector := endpoint.putBody["selectors"].([]any)[0].(map[string]any)
	if _, ok := selector["purpose"]; ok {
		t.Errorf("purpose should have been cleared, body: %s", endpoint.rawPutBody)
	}
	criteria := selector["criteria"].(map[string]any)
	if _, ok := criteria["case_insensitive"]; ok {
		t.Errorf("case_insensitive should have been cleared, body: %s", endpoint.rawPutBody)
	}
}

// TestAbpUpdateFailsWhenTheCurrentObjectCannotBeRead keeps the fallback
// honest: a full replace built from an unknown starting point is how
// configuration used to get destroyed, so it is not attempted.
func TestAbpUpdateFailsWhenTheCurrentObjectCannotBeRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message":"Database Error"}`))
			return
		}
		t.Errorf("the update should not have been sent")
	}))
	defer server.Close()

	client := &Client{
		config:     &Config{APIID: "id", APIKey: "key", BaseURLAPI: server.URL},
		httpClient: &http.Client{},
	}

	_, err := client.UpdateAbpSite("11111111-1111-1111-1111-111111111111", AbpSite{Name: "site"})
	if err == nil {
		t.Fatal("expected an error")
	}
	if got := err.Error(); !strings.Contains(got, "Database Error") || !strings.Contains(got, "ABP Site") {
		t.Errorf("unhelpful error: %v", err)
	}
}

func lookupJsonPath(t *testing.T, document map[string]any, path []string) any {
	t.Helper()

	var value any = document
	for _, key := range path {
		switch container := value.(type) {
		case map[string]any:
			value = container[key]
		case []any:
			index := int(key[0] - '0')
			if index >= len(container) {
				return nil
			}
			value = container[index]
		default:
			return nil
		}
	}
	return value
}
