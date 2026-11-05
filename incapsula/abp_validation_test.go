package incapsula

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// attributeValidator resolves a dotted attribute path to the ValidateFunc the
// provider's schema actually uses, so these tests fail if an attribute is
// validated by a different function or by none at all.
func attributeValidator(t *testing.T, resource, path string) schema.SchemaValidateFunc {
	t.Helper()

	r, ok := Provider().ResourcesMap[resource]
	if !ok {
		t.Fatalf("resource %s is not registered", resource)
	}

	current := r.Schema
	keys := strings.Split(path, ".")
	for _, key := range keys[:len(keys)-1] {
		nested, ok := current[key]
		if !ok {
			t.Fatalf("%s has no attribute %s", resource, key)
		}
		elem, ok := nested.Elem.(*schema.Resource)
		if !ok {
			t.Fatalf("%s.%s is not a block", resource, key)
		}
		current = elem.Schema
	}

	attribute, ok := current[keys[len(keys)-1]]
	if !ok {
		t.Fatalf("%s has no attribute %s", resource, path)
	}
	if attribute.ValidateFunc == nil {
		t.Fatalf("%s.%s has no ValidateFunc, so a value the API would not store verbatim is accepted", resource, path)
	}
	return attribute.ValidateFunc
}

// The API normalises these attributes, so a value it would not store verbatim
// leaves a plan that applying cannot resolve — and for the Domain's criteria,
// which forces replacement, one that destroys and recreates the Domain on every
// apply. Each one must be rejected while the configuration is planned.
func TestAbpNormalisedAttributesRejectUnstorableValues(t *testing.T) {
	cases := []struct {
		resource  string
		attribute string
		value     any
		want      string
	}{
		{"incapsula_abp_policy", "name", "  policy  ", "leading or trailing whitespace"},
		{"incapsula_abp_policy", "name", "   ", "must not be empty or consist only of whitespace"},
		{"incapsula_abp_site", "name", "  site  ", "leading or trailing whitespace"},
		{"incapsula_abp_site", "name", "   ", "must not be empty or consist only of whitespace"},
		{"incapsula_abp_site", "name", strings.Repeat("s", 101), "expected length of name to be in the range (1 - 100)"},
		{"incapsula_abp_proof_of_work_configuration", "name", "  pow  ", "leading or trailing whitespace"},
		{"incapsula_abp_proof_of_work_configuration", "name", "   ", "must not be empty or consist only of whitespace"},
		{"incapsula_abp_proof_of_work_configuration", "name", strings.Repeat("p", 101), "expected length of name to be in the range (1 - 100)"},
		{"incapsula_abp_domain", "cookiescope", "  example.com  ", "leading or trailing whitespace"},
		{"incapsula_abp_domain", "criteria.exact", "  example.com  ", "every apply would replace the Domain"},
		{"incapsula_abp_domain", "criteria.prefix", "example.com ", "every apply would replace the Domain"},
		{"incapsula_abp_domain", "criteria.suffix", "\texample.com", "every apply would replace the Domain"},
		{"incapsula_abp_domain", "criteria.exact", "   ", "must not consist only of whitespace"},
	}

	for _, test := range cases {
		t.Run(test.resource+"."+test.attribute+"="+test.value.(string), func(t *testing.T) {
			name := test.attribute[strings.LastIndex(test.attribute, ".")+1:]
			_, errs := attributeValidator(t, test.resource, test.attribute)(test.value, name)
			if len(errs) == 0 {
				t.Fatalf("%q was accepted, want it rejected with %q", test.value, test.want)
			}
			if !strings.Contains(errs[0].Error(), test.want) {
				t.Errorf("rejected with %q, want it to mention %q", errs[0], test.want)
			}
		})
	}
}

// The same attributes must accept the values they are meant to carry. An empty
// Domain criteria means "match on another criteria instead", and an empty
// cookiescope is a value the API stores as such, so neither is this
// validation's business.
func TestAbpNormalisedAttributesAcceptStorableValues(t *testing.T) {
	cases := []struct {
		resource  string
		attribute string
		value     any
	}{
		{"incapsula_abp_policy", "name", "policy"},
		{"incapsula_abp_site", "name", "site"},
		{"incapsula_abp_site", "name", strings.Repeat("s", 100)},
		{"incapsula_abp_proof_of_work_configuration", "name", "pow"},
		{"incapsula_abp_domain", "cookiescope", "example.com"},
		{"incapsula_abp_domain", "cookiescope", ""},
		{"incapsula_abp_domain", "criteria.exact", "example.com"},
		{"incapsula_abp_domain", "criteria.exact", ""},
		{"incapsula_abp_domain", "criteria.prefix", ""},
		{"incapsula_abp_domain", "criteria.suffix", ""},
	}

	for _, test := range cases {
		t.Run(test.resource+"."+test.attribute+"="+test.value.(string), func(t *testing.T) {
			name := test.attribute[strings.LastIndex(test.attribute, ".")+1:]
			_, errs := attributeValidator(t, test.resource, test.attribute)(test.value, name)
			if len(errs) != 0 {
				t.Fatalf("%q was rejected with %v", test.value, errs)
			}
		})
	}
}
