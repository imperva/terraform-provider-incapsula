package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Extra seeders and checks the Policy and Directive tests share.

// abpPolicySeedGlobalPolicy seeds the Account's global Policy and keeps it out
// of the Account's Policy list, which is what the API does: the global Policy
// is only reachable through the Account's own global_policy endpoint, so
// `incapsula_abp_policies` never reports it.
func (f *abpFake) abpPolicySeedGlobalPolicy(name string, actions ...string) string {
	id := f.seedGlobalPolicy(name, actions...)

	f.mu.Lock()
	defer f.mu.Unlock()
	kept := make([]string, 0, len(f.order[abpFakePolicies]))
	for _, existing := range f.order[abpFakePolicies] {
		if existing != id {
			kept = append(kept, existing)
		}
	}
	f.order[abpFakePolicies] = kept
	return id
}

// abpPolicySeedWithDirectives seeds a Policy whose Directives are given
// verbatim, for the attributes seedPolicy does not cover.
func (f *abpFake) abpPolicySeedWithDirectives(name, description string, directives ...map[string]any) string {
	raw := make([]any, 0, len(directives))
	for _, directive := range directives {
		raw = append(raw, directive)
	}
	return f.seed(abpFakePolicies, map[string]any{
		"name":        name,
		"description": description,
		"directives":  raw,
	})
}

// abpPolicyStoredDirectives returns the Directives the fake holds for a Policy.
func abpPolicyStoredDirectives(t *testing.T, f *abpFake, policyId string) []map[string]any {
	t.Helper()
	doc := f.doc(abpFakePolicies, policyId)
	if doc == nil {
		t.Fatalf("policy %s is not stored", policyId)
	}
	return abpFakeList(doc["directives"])
}

// abpPolicyNotReported fails unless each attribute is absent or empty, which is
// how the data sources report a value they deliberately do not fill in.
func abpPolicyNotReported(address string, keys ...string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		found := s.RootModule().Resources[address]
		if found == nil {
			return fmt.Errorf("%s is not in state", address)
		}
		for _, key := range keys {
			if value := found.Primary.Attributes[key]; value != "" {
				return fmt.Errorf("%s.%s is %q, expected it not to be reported", address, key, value)
			}
		}
		return nil
	}
}
