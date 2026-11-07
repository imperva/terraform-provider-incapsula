package incapsula

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// incapsula_abp_account_site_priority ranks the Sites of one Account. Unlike
// its Domain sibling it takes no subset: the API requires the whole Account's
// Site list and the provider passes the configured list straight through, so
// one Site it does not name blocks every apply. These tests pin both halves —
// the list that covers everything, and what happens when it stops doing so.

// --- create and update ------------------------------------------------------

func TestAbpPrioritySites_CoveringListApplies(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	config := abpPrioritySitesConfig(abpFakeAccountId, second, first)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPrioritySitesResource)
					if attrs["id"] != abpFakeAccountId {
						return fmt.Errorf("the resource ID is %s, want the Account ID %s", attrs["id"], abpFakeAccountId)
					}
					if err := abpPriorityCheckIds("site_ids", abpPriorityStateList(attrs, "site_ids"), []string{second, first}); err != nil {
						return err
					}
					return abpPriorityCheckIds("the Account's order", f.siteOrderOf(abpFakeAccountId), []string{second, first})
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Reordering the list is an in-place update; the resource ID is the Account's,
// so it never changes.
func TestAbpPrioritySites_UpdatedInPlace(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	third := f.seedSite("site-three")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPrioritySitesConfig(abpFakeAccountId, first, second, third)},
			{
				Config: abpPrioritySitesConfig(abpFakeAccountId, third, first, second),
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPrioritySitesResource)["id"]; got != abpFakeAccountId {
						return fmt.Errorf("the resource ID became %s, want the Account ID %s", got, abpFakeAccountId)
					}
					return abpPriorityCheckIds("the Account's order", f.siteOrderOf(abpFakeAccountId), []string{third, first, second})
				},
			},
		},
	})
}

// `account_id` is ForceNew: the replacement forgets the old Account's order
// instead of resetting it.
func TestAbpPrioritySites_AccountIdForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	elsewhere := abpPrioritySeedSiteIn(f, abpFakeOtherAccountId, "other-account-site")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPrioritySitesConfig(abpFakeAccountId, second, first)},
			{
				Config: abpPrioritySitesConfig(abpFakeOtherAccountId, elsewhere),
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPrioritySitesResource)["id"]; got != abpFakeOtherAccountId {
						return fmt.Errorf("the resource ID is %s, want the new Account ID %s", got, abpFakeOtherAccountId)
					}
					return abpPriorityCheckIds("the old Account's order", f.siteOrderOf(abpFakeAccountId), []string{second, first})
				},
			},
		},
	})
}

// --- the list has to cover the Account --------------------------------------

// A Site the list does not name fails the apply with the API's own message,
// which says neither how many Sites were expected nor which ones are missing.
// That opaqueness is a limitation, not an intention: the provider forwards the
// list without checking it against the Account.
func TestAbpPrioritySites_MustCoverEverySite(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	f.seedSite("site-two")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPrioritySitesConfig(abpFakeAccountId, first),
				ExpectError: regexp.MustCompile("(?s)Wrong number of"),
			},
		},
	})
}

// A Site created by anyone else is enough to make a literal list fail on every
// apply, which is why the documentation tells the user not to write one.
func TestAbpPrioritySites_SiteCreatedOutOfBandWedgesALiteralList(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	config := abpPrioritySitesConfig(abpFakeAccountId, second, first)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig:   func() { f.seedSite("created-in-the-ui") },
				Config:      config,
				ExpectError: regexp.MustCompile("(?s)Wrong number of"),
			},
		},
	})
}

// The form the sample configuration prescribes: the Sites whose order matters
// first, then every other Site of the Account appended from the
// incapsula_abp_sites data source. A Site created out of band is picked up by
// the next plan instead of blocking it.
func TestAbpPrioritySites_ConcatWithTheDataSourceSurvivesANewSite(t *testing.T) {
	f := newAbpFake(t)
	existing := f.seedSite("pre-existing")
	config := fmt.Sprintf(`
data "incapsula_abp_sites" "all" {
  account_id = %[1]q
}

resource "incapsula_abp_site" "owned" {
  account_id = %[1]q
  name       = "owned-by-terraform"
}

locals {
  managed_sites = [incapsula_abp_site.owned.id]
}

resource "incapsula_abp_account_site_priority" "test" {
  account_id = %[1]q
  site_ids = concat(
    local.managed_sites,
    [for site in data.incapsula_abp_sites.all.sites : site.id if !contains(local.managed_sites, site.id)]
  )
}
`, abpFakeAccountId)

	var appeared string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					owned := abpAttrs(t, s, "incapsula_abp_site.owned")["id"]
					return abpPriorityCheckIds("the Account's order", f.siteOrderOf(abpFakeAccountId), []string{owned, existing})
				},
			},
			{
				PreConfig: func() { appeared = f.seedSite("created-in-the-ui") },
				Config:    config,
				Check: func(s *terraform.State) error {
					owned := abpAttrs(t, s, "incapsula_abp_site.owned")["id"]
					want := []string{owned, existing, appeared}
					if err := abpPriorityCheckAttrList(t, s, abpPrioritySitesResource, "site_ids", want); err != nil {
						return err
					}
					return abpPriorityCheckIds("the Account's order", f.siteOrderOf(abpFakeAccountId), want)
				},
			},
		},
	})
}

// A Site deleted out of band wedges the list in the same way until it is
// dropped from it — the sample configuration's note that a destroyed Site
// takes one more apply to converge.
func TestAbpPrioritySites_DeletedSiteHasToLeaveTheList(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	third := f.seedSite("site-three")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPrioritySitesConfig(abpFakeAccountId, first, second, third)},
			{
				PreConfig:   func() { f.drop(abpFakeSites, second) },
				Config:      abpPrioritySitesConfig(abpFakeAccountId, first, second, third),
				ExpectError: regexp.MustCompile("(?s)Wrong number of"),
			},
			{
				Config: abpPrioritySitesConfig(abpFakeAccountId, third, first),
				Check: func(s *terraform.State) error {
					return abpPriorityCheckIds("the Account's order", f.siteOrderOf(abpFakeAccountId), []string{third, first})
				},
			},
		},
	})
}

// --- drift ------------------------------------------------------------------

// Reordering the Sites out of band is drift the configuration can express, so
// it is reported and restored.
func TestAbpPrioritySites_ReorderedSitesAreRestored(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	third := f.seedSite("site-three")
	config := abpPrioritySitesConfig(abpFakeAccountId, first, second, third)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() { f.reorderSites(third, second, first) },
				Config:    config,
				Check: func(s *terraform.State) error {
					return abpPriorityCheckIds("the restored order", f.siteOrderOf(abpFakeAccountId), []string{first, second, third})
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// --- validation and plan-time errors ----------------------------------------

// Duplicates are rejected by the provider: the API does not reject them, so
// without the guard a repeated ID would be written as an order.
func TestAbpPrioritySites_DuplicateIdsAreRejected(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	f.seedSite("site-two")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPrioritySitesConfig(abpFakeAccountId, first, first),
				ExpectError: regexp.MustCompile("(?s)is a duplicate of an earlier entry"),
			},
		},
	})

	if count := f.requestCount(http.MethodPut, "site_priority"); count != 0 {
		t.Fatalf("the rejected configuration still wrote %d order(s)", count)
	}
}

// An ID naming a Site of another Account is only caught by the API, and only
// because the provider has no Account-level check of its own: it reaches the
// service and comes back as an unknown ID.
func TestAbpPrioritySites_SiteOfAnotherAccountIsRejected(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	f.seedSite("site-two")
	elsewhere := abpPrioritySeedSiteIn(f, abpFakeOtherAccountId, "other-account-site")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPrioritySitesConfig(abpFakeAccountId, first, elsewhere),
				ExpectError: regexp.MustCompile("(?s)Unknown ID"),
			},
		},
	})
}

// An empty list is rejected when the configuration is validated, so an Account
// with no Sites cannot be expressed. A limitation rather than an intention.
func TestAbpPrioritySites_EmptyListIsRejected(t *testing.T) {
	f := newAbpFake(t)
	f.seedSite("site-one")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPrioritySitesConfig(abpFakeAccountId),
				ExpectError: regexp.MustCompile("(?s)site_ids requires 1 item minimum"),
			},
		},
	})
}

func TestAbpPrioritySites_NonUuidIdsAreRejected(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPrioritySitesConfig("not-a-uuid", first),
				ExpectError: regexp.MustCompile(`(?s)expected "account_id" to be a valid UUID`),
			},
			{
				Config:      abpPrioritySitesConfig(abpFakeAccountId, "not-a-uuid"),
				ExpectError: regexp.MustCompile("(?s)to be a valid UUID"),
			},
		},
	})
}

// --- import and destroy -----------------------------------------------------

// Importing adopts the Account's whole order, which is also the only list the
// resource accepts, so the imported resource re-plans empty.
func TestAbpPrioritySites_ImportAdoptsTheWholeOrder(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")
	third := f.seedSite("site-three")
	f.reorderSites(second, third, first)
	config := abpPrioritySitesConfig(abpFakeAccountId, second, third, first)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpPrioritySitesResource,
				ImportState:        true,
				ImportStateId:      abpFakeAccountId,
				ImportStatePersist: true,
				Check: func(s *terraform.State) error {
					return abpPriorityCheckAttrList(t, s, abpPrioritySitesResource, "site_ids", []string{second, third, first})
				},
			},
			{Config: config, PlanOnly: true},
			{Config: config, ResourceName: abpPrioritySitesResource, ImportState: true, ImportStateId: abpFakeAccountId, ImportStateVerify: true},
		},
	})
}

// An Account that does not exist is indistinguishable from one without Sites,
// because the order is read as a property of the Account rather than as an
// object: the import succeeds and yields an empty list, which no apply can then
// use (`site_ids` needs one entry). A limitation rather than an intention.
func TestAbpPrioritySites_ImportOfAnUnknownAccountYieldsAnEmptyList(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        abpPrioritySitesConfig(abpFakeAccountId, first),
				ResourceName:  abpPrioritySitesResource,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPrioritySitesResource)
					if attrs["id"] != abpFakeAbsentId {
						return fmt.Errorf("the imported ID is %s, want %s", attrs["id"], abpFakeAbsentId)
					}
					return abpPriorityCheckIds("the imported site_ids", abpPriorityStateList(attrs, "site_ids"), nil)
				},
			},
		},
	})
}

// The order is a property of the Account, not an object of its own: destroying
// the resource only drops it from state, and the Account keeps the order it was
// given.
func TestAbpPrioritySites_DestroyLeavesTheOrderInPlace(t *testing.T) {
	f := newAbpFake(t)
	first := f.seedSite("site-one")
	second := f.seedSite("site-two")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if count := f.requestCount(http.MethodDelete, "site_priority"); count != 0 {
				return fmt.Errorf("destroying the order sent %d DELETE(s)", count)
			}
			return abpPriorityCheckIds("the Account's order after destroy", f.siteOrderOf(abpFakeAccountId), []string{second, first})
		},
		Steps: []resource.TestStep{
			{Config: abpPrioritySitesConfig(abpFakeAccountId, second, first)},
		},
	})

	if f.count(abpFakeSites) != 2 {
		t.Fatalf("destroying the order deleted Sites: %d left of 2", f.count(abpFakeSites))
	}
}
