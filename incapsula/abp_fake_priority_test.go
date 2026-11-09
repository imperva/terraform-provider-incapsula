package incapsula

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Helpers shared by the two priority ordering resources. Both are a view on a
// server-side order rather than an object of their own, so most of what the
// tests assert is a list of ids: the one in state, and the one the service
// ended up with.

const (
	abpPriorityDomainsResource = "incapsula_abp_site_domain_priority.test"
	abpPrioritySitesResource   = "incapsula_abp_account_site_priority.test"
)

func abpPriorityQuoted(ids []string) string {
	quoted := make([]string, 0, len(ids))
	for _, id := range ids {
		quoted = append(quoted, strconv.Quote(id))
	}
	return strings.Join(quoted, ", ")
}

func abpPriorityDomainsConfig(siteId string, domainIds ...string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_site_domain_priority" "test" {
  site_id    = %q
  domain_ids = [%s]
}
`, siteId, abpPriorityQuoted(domainIds))
}

func abpPrioritySitesConfig(accountId string, siteIds ...string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_account_site_priority" "test" {
  account_id = %q
  site_ids   = [%s]
}
`, accountId, abpPriorityQuoted(siteIds))
}

// abpPriorityStateList reads a list attribute out of the flat state
// representation.
func abpPriorityStateList(attrs map[string]string, key string) []string {
	count, _ := strconv.Atoi(attrs[key+".#"])
	ids := make([]string, 0, count)
	for i := 0; i < count; i++ {
		ids = append(ids, attrs[fmt.Sprintf("%s.%d", key, i)])
	}
	return ids
}

func abpPriorityCheckIds(label string, got, want []string) error {
	if strings.Join(got, ",") != strings.Join(want, ",") {
		return fmt.Errorf("%s is %v, want %v", label, got, want)
	}
	return nil
}

// abpPriorityCheckAttrList asserts a list attribute of one resource in state.
func abpPriorityCheckAttrList(t *testing.T, s *terraform.State, address, key string, want []string) error {
	t.Helper()
	return abpPriorityCheckIds(key, abpPriorityStateList(abpAttrs(t, s, address), key), want)
}

// abpPrioritySeedSite seeds a Site with `count` Domains attached to it, as an
// earlier workspace or the UI would have left them.
func abpPrioritySeedSite(f *abpFake, name string, count int) (string, []string) {
	siteId := f.seedSite(name)
	domainIds := make([]string, 0, count)
	for i := 0; i < count; i++ {
		domainIds = append(domainIds, f.seedDomain(siteId, fmt.Sprintf("%s-%d.example.com", name, i), ""))
	}
	return siteId, domainIds
}

// abpPrioritySeedSiteIn seeds a Site belonging to a specific Account, which
// seedSite cannot express.
func abpPrioritySeedSiteIn(f *abpFake, accountId, name string) string {
	return f.seed(abpFakeSites, map[string]any{
		"account_id": accountId,
		"name":       name,
		"selectors":  []any{},
	})
}
