package incapsula

import (
	"encoding/json"
	"reflect"
	"testing"
)

// mergeForTest merges overlay onto base the way the update path does, and
// returns the result as a generic value so tests can compare it against a
// literal document.
func mergeForTest(t *testing.T, model any, base, overlay string) map[string]any {
	t.Helper()

	merged, err := abpMergeWireDocuments([]byte(base), []byte(overlay), abpWireModelOf(reflect.TypeOf(model)))
	if err != nil {
		t.Fatalf("merging: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(merged, &got); err != nil {
		t.Fatalf("the merged document does not parse: %v; body: %s", err, merged)
	}
	return got
}

func wantDocument(t *testing.T, got map[string]any, want string) {
	t.Helper()

	var expected map[string]any
	if err := json.Unmarshal([]byte(want), &expected); err != nil {
		t.Fatalf("the expected document does not parse: %v", err)
	}
	if !reflect.DeepEqual(got, expected) {
		gotJson, _ := json.Marshal(got)
		t.Errorf("merged document mismatch\n got: %s\nwant: %s", gotJson, want)
	}
}

// TestAbpWireMergePreservesUnmodelledFields is the heart of it: a key the
// payload type knows nothing about survives a full-replace update, while a key
// it does model is dictated by the payload even when the payload omits it.
func TestAbpWireMergePreservesUnmodelledFields(t *testing.T) {
	got := mergeForTest(t, AbpSite{},
		`{
			"id": "site-1",
			"name": "old",
			"industry": "travel",
			"mx_hostname_id": "mx-7",
			"selectors": [],
			"default_max_requests_per_minute": 12,
			"created_at": "2026-01-01T00:00:00Z",
			"_links": {"site": {"href": "/v1/site/site-1"}}
		}`,
		`{"name": "new", "selectors": []}`,
	)

	wantDocument(t, got, `{
		"name": "new",
		"selectors": [],
		"mx_hostname_id": "mx-7"
	}`)
}

func TestAbpWireMergeModelledFieldsFollowThePayload(t *testing.T) {
	t.Run("an omitted optional still clears", func(t *testing.T) {
		got := mergeForTest(t, AbpSite{},
			`{"name": "site", "industry": "travel", "selectors": []}`,
			`{"name": "site", "selectors": []}`,
		)
		wantDocument(t, got, `{"name": "site", "selectors": []}`)
	})

	t.Run("a new value wins", func(t *testing.T) {
		got := mergeForTest(t, AbpSite{},
			`{"name": "site", "industry": "travel", "selectors": []}`,
			`{"name": "site", "industry": "retail", "selectors": []}`,
		)
		wantDocument(t, got, `{"name": "site", "industry": "retail", "selectors": []}`)
	})
}

func TestAbpWireMergeKeepsUnmodelledSelectorFields(t *testing.T) {
	got := mergeForTest(t, AbpSite{},
		`{
			"name": "site",
			"selectors": [
				{
					"id": "sel-1",
					"policy_id": null,
					"criteria": {"path_prefix": "/login", "case_insensitive": true},
					"purpose": "login",
					"future_field": "keep me",
					"derived_id": "derived-1",
					"created_at": "2026-01-01T00:00:00Z"
				}
			]
		}`,
		`{
			"name": "site",
			"selectors": [
				{
					"id": "sel-1",
					"policy_id": null,
					"criteria": {"path_prefix": "/login"},
					"analysis_settings": {"rate_limiting": "none"}
				}
			]
		}`,
	)

	wantDocument(t, got, `{
		"name": "site",
		"selectors": [
			{
				"id": "sel-1",
				"policy_id": null,
				"criteria": {"path_prefix": "/login"},
				"analysis_settings": {"rate_limiting": "none"},
				"future_field": "keep me"
			}
		]
	}`)
}

// TestAbpWireMergeDropsTheLeftoversOfASwitchedVariant guards the one case
// where merging would be worse than clearing: the API selects a criteria
// variant by key presence in a fixed order, so a body carrying both
// `path_prefix` and `path_regex` silently keeps the prefix.
func TestAbpWireMergeDropsTheLeftoversOfASwitchedVariant(t *testing.T) {
	got := mergeForTest(t, AbpSite{},
		`{
			"name": "site",
			"selectors": [
				{
					"id": "sel-1",
					"criteria": {"path_prefix": "/login", "case_insensitive": true, "future_prefix_option": 1}
				}
			]
		}`,
		`{
			"name": "site",
			"selectors": [
				{"id": "sel-1", "criteria": {"path_regex": "^/login"}}
			]
		}`,
	)

	wantDocument(t, got, `{
		"name": "site",
		"selectors": [
			{"id": "sel-1", "criteria": {"path_regex": "^/login"}}
		]
	}`)
}

func TestAbpWireMergeKeepsVariantExtrasWhenTheVariantIsUnchanged(t *testing.T) {
	got := mergeForTest(t, AbpSite{},
		`{
			"name": "site",
			"selectors": [
				{"id": "sel-1", "criteria": {"path_prefix": "/a", "future_prefix_option": 1}}
			]
		}`,
		`{
			"name": "site",
			"selectors": [
				{"id": "sel-1", "criteria": {"path_prefix": "/b", "case_insensitive": true}}
			]
		}`,
	)

	wantDocument(t, got, `{
		"name": "site",
		"selectors": [
			{
				"id": "sel-1",
				"criteria": {"path_prefix": "/b", "case_insensitive": true, "future_prefix_option": 1}
			}
		]
	}`)
}

func TestAbpWireMergeIdentifiesListElementsById(t *testing.T) {
	t.Run("order comes from the payload and removed elements go", func(t *testing.T) {
		got := mergeForTest(t, AbpSite{},
			`{
				"name": "site",
				"selectors": [
					{"id": "a", "purpose": "login", "future": "a"},
					{"id": "b", "future": "b"},
					{"id": "c", "future": "c"}
				]
			}`,
			`{
				"name": "site",
				"selectors": [
					{"id": "c", "criteria": {"path_prefix": "/c"}},
					{"id": "a", "criteria": {"path_prefix": "/a"}}
				]
			}`,
		)

		wantDocument(t, got, `{
			"name": "site",
			"selectors": [
				{"id": "c", "criteria": {"path_prefix": "/c"}, "future": "c"},
				{"id": "a", "criteria": {"path_prefix": "/a"}, "future": "a"}
			]
		}`)
	})

	t.Run("an element without an id is taken as it is", func(t *testing.T) {
		got := mergeForTest(t, AbpSite{},
			`{"name": "site", "selectors": [{"id": "a", "future": "a"}]}`,
			`{"name": "site", "selectors": [{"criteria": {"path_prefix": "/new"}}, {"id": "a"}]}`,
		)

		wantDocument(t, got, `{
			"name": "site",
			"selectors": [
				{"criteria": {"path_prefix": "/new"}},
				{"id": "a", "future": "a"}
			]
		}`)
	})
}

// TestAbpWireMergeReplacesKeylessLists pins down why position is not identity.
// Policy directives have no id, and the API drops a directive whose
// `origin_policy_id` names another policy, so index-merging an inherited
// directive's marker onto a directive of this policy's own would delete it.
func TestAbpWireMergeReplacesKeylessLists(t *testing.T) {
	got := mergeForTest(t, AbpPolicy{},
		`{
			"name": "policy",
			"directives": [
				{"action": "allow", "condition_id": "c1", "origin_policy_id": "global-policy"},
				{"action": "block", "condition_id": "c2", "origin_policy_id": "this-policy"}
			]
		}`,
		`{
			"name": "policy",
			"directives": [
				{"action": "block", "condition_id": "c2"}
			]
		}`,
	)

	wantDocument(t, got, `{
		"name": "policy",
		"directives": [{"action": "block", "condition_id": "c2"}]
	}`)
}

// TestAbpWireMergeKeepsNumbersVerbatim doubles as the Domain case: `criteria`
// is not part of UpdateDomainV1, so it is carried over untouched -- and a
// CloudWAF id must not become 9.007199254740992e+15 on the way.
func TestAbpWireMergeKeepsNumbersVerbatim(t *testing.T) {
	got := mergeForTest(t, abpDomainUpdatePayload{},
		`{"site_id": "s", "criteria": {"cloudwaf_id": 9007199254740993}}`,
		`{"site_id": "s"}`,
	)

	wantDocument(t, got, `{"site_id": "s", "criteria": {"cloudwaf_id": 9007199254740993}}`)
}

func TestAbpWireMergeTakesThePayloadWhenTheShapeChanges(t *testing.T) {
	got := mergeForTest(t, abpDomainUpdatePayload{},
		`{"site_id": "s", "captcha_settings": {"geetest": {"geetest_captcha_id": "id", "geetest_private_key": "key"}}}`,
		`{"site_id": "s", "captcha_settings": "none"}`,
	)

	wantDocument(t, got, `{"site_id": "s", "captcha_settings": "none"}`)
}

func TestAbpWireModelDescribesOnlyModelledKeys(t *testing.T) {
	site := abpWireModelOf(reflect.TypeOf(AbpSite{}))

	for _, key := range []string{"id", "account_id", "name", "industry", "selectors", "created_at", "modified_at"} {
		if !site.owns(key) {
			t.Errorf("AbpSite should model %q", key)
		}
	}
	for _, key := range []string{"mx_hostname_id", "future_field"} {
		if site.owns(key) {
			t.Errorf("AbpSite should not model %q", key)
		}
	}

	selector := site.child("selectors").elemModel()
	for _, key := range []string{"id", "policy_id", "criteria", "purpose", "analysis_settings", "derived_id"} {
		if !selector.owns(key) {
			t.Errorf("AbpSelector should model %q", key)
		}
	}

	criteria := selector.child("criteria")
	if criteria.oneOf == nil {
		t.Fatal("selector criteria should be described as a oneOf")
	}
	for _, key := range []string{"path_prefix", "path_regex", "postback", "api_security_group", "case_insensitive"} {
		if !criteria.owns(key) {
			t.Errorf("selector criteria should model %q", key)
		}
	}
}
