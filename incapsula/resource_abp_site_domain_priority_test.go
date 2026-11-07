package incapsula

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// incapsula_abp_site_domain_priority ranks the Domains of one Site. The API
// only accepts an order covering exactly the Domains that are attached, so the
// resource takes a subset — a prefix of the order — and the provider pads it
// with the Domains the configuration does not name. These tests are mostly
// about what happens to the Domains it does not rank.

// --- unit tests for the order arithmetic ------------------------------------

// The padding puts the ranked Domains first and keeps the relative order of the
// rest, which is what makes a Domain created elsewhere the lowest priority one.
func TestAbpPriorityCompleteOrder(t *testing.T) {
	for _, test := range []struct {
		name     string
		ranked   []string
		attached []string
		want     []string
	}{
		{"a subset is padded", []string{"c", "a"}, []string{"a", "b", "c"}, []string{"c", "a", "b"}},
		{"a full list is kept", []string{"b", "a"}, []string{"a", "b"}, []string{"b", "a"}},
		{"nothing ranked keeps the server order", nil, []string{"a", "b"}, []string{"a", "b"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if err := abpPriorityCheckIds("the completed order", abpPriorityCompleteOrder(test.ranked, test.attached), test.want); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// Reading projects the server's order onto the ranked Domains only: pulling an
// unranked one into `domain_ids` would propose a removal the API cannot do.
func TestAbpPriorityProjectAndPrefix(t *testing.T) {
	order := []string{"c", "a", "b"}
	if err := abpPriorityCheckIds("the projection", abpPriorityProject(order, []string{"a", "c"}), []string{"c", "a"}); err != nil {
		t.Fatal(err)
	}
	if err := abpPriorityCheckIds("the unattached ids", abpPriorityUnattachedIds([]string{"a", "z"}, order), []string{"z"}); err != nil {
		t.Fatal(err)
	}
	if !abpPriorityHasPrefix(order, []string{"c", "a"}) {
		t.Fatal("the ranked Domains are the highest priority ones, but that was not detected")
	}
	if abpPriorityHasPrefix(order, []string{"a", "b"}) {
		t.Fatal("an unranked Domain in front was accepted as a prefix")
	}
}

// --- create -----------------------------------------------------------------

// A list that ranks only some of the Site's Domains applies, puts those Domains
// on top in the order given, and leaves the rest behind them.
func TestAbpPriorityDomains_SubsetIsPadded(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "subset", 3)
	config := abpPriorityDomainsConfig(siteId, domains[2], domains[0])
	effective := []string{domains[2], domains[0], domains[1]}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "domain_ids", []string{domains[2], domains[0]}); err != nil {
						return err
					}
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "effective_domain_ids", effective); err != nil {
						return err
					}
					if got := abpAttrs(t, s, abpPriorityDomainsResource)["id"]; got != siteId {
						return fmt.Errorf("the resource ID is %s, want the Site ID %s", got, siteId)
					}
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(siteId), effective)
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// The full-list form still works, and then `effective_domain_ids` is just
// `domain_ids`.
func TestAbpPriorityDomains_FullListIsItsOwnEffectiveOrder(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "full", 3)
	want := []string{domains[1], domains[2], domains[0]}
	config := abpPriorityDomainsConfig(siteId, want...)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "effective_domain_ids", want); err != nil {
						return err
					}
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(siteId), want)
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Domains created in the same apply can be ranked straight away: their IDs are
// only known once they exist, so the order is written after them.
func TestAbpPriorityDomains_RanksDomainsCreatedInTheSameApply(t *testing.T) {
	f := newAbpFake(t)
	config := fmt.Sprintf(`
resource "incapsula_abp_site" "site" {
  account_id = %[1]q
  name       = "ranked-in-one-apply"
}

resource "incapsula_abp_domain" "first" {
  account_id  = %[1]q
  site_id     = incapsula_abp_site.site.id
  cookiescope = "first.example.com"

  criteria {
    exact = "first.example.com"
  }
}

resource "incapsula_abp_domain" "second" {
  account_id  = %[1]q
  site_id     = incapsula_abp_site.site.id
  cookiescope = "second.example.com"

  criteria {
    exact = "second.example.com"
  }
}

resource "incapsula_abp_site_domain_priority" "test" {
  site_id    = incapsula_abp_site.site.id
  domain_ids = [incapsula_abp_domain.second.id]
}
`, abpFakeAccountId)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					attrs := abpAttrs(t, s, abpPriorityDomainsResource)
					first := abpAttrs(t, s, "incapsula_abp_domain.first")["id"]
					second := abpAttrs(t, s, "incapsula_abp_domain.second")["id"]
					want := []string{second, first}
					if err := abpPriorityCheckIds("effective_domain_ids", abpPriorityStateList(attrs, "effective_domain_ids"), want); err != nil {
						return err
					}
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(attrs["site_id"]), want)
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// --- update and replacement -------------------------------------------------

// Reordering, extending and shrinking the ranked list are in-place updates: the
// resource ID is the Site's, so it never changes.
func TestAbpPriorityDomains_UpdatedInPlace(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "update", 3)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPriorityDomainsConfig(siteId, domains[0], domains[1])},
			{
				// Reordered and extended.
				Config: abpPriorityDomainsConfig(siteId, domains[2], domains[1], domains[0]),
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPriorityDomainsResource)["id"]; got != siteId {
						return fmt.Errorf("the resource ID became %s, want the Site ID %s", got, siteId)
					}
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(siteId), []string{domains[2], domains[1], domains[0]})
				},
			},
			{
				// Shrunk back to one ranked Domain: the two it no longer ranks
				// keep the order they had.
				Config: abpPriorityDomainsConfig(siteId, domains[1]),
				Check: func(s *terraform.State) error {
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(siteId), []string{domains[1], domains[2], domains[0]})
				},
			},
		},
	})
}

// `site_id` is ForceNew: the replacement forgets the old Site's order instead
// of resetting it.
func TestAbpPriorityDomains_SiteIdForcesReplacement(t *testing.T) {
	f := newAbpFake(t)
	firstSite, firstDomains := abpPrioritySeedSite(f, "first-site", 2)
	secondSite, secondDomains := abpPrioritySeedSite(f, "second-site", 2)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPriorityDomainsConfig(firstSite, firstDomains[1], firstDomains[0])},
			{
				Config: abpPriorityDomainsConfig(secondSite, secondDomains[1]),
				Check: func(s *terraform.State) error {
					if got := abpAttrs(t, s, abpPriorityDomainsResource)["id"]; got != secondSite {
						return fmt.Errorf("the resource ID is %s, want the new Site ID %s", got, secondSite)
					}
					if err := abpPriorityCheckIds("the old Site's order", f.domainOrderOf(firstSite), []string{firstDomains[1], firstDomains[0]}); err != nil {
						return err
					}
					return abpPriorityCheckIds("the new Site's order", f.domainOrderOf(secondSite), []string{secondDomains[1], secondDomains[0]})
				},
			},
		},
	})
}

// --- drift ------------------------------------------------------------------

// The finding this resource was wedged by: a Domain attached to the Site but
// absent from `domain_ids`. It must not produce a diff, let alone an apply that
// fails with a count mismatch on every run.
func TestAbpPriorityDomains_DomainAttachedOutOfBandDoesNotWedge(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "out-of-band", 2)
	config := abpPriorityDomainsConfig(siteId, domains[0], domains[1])
	var appeared string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				// A Domain created in the ABP UI, by another workspace or by an
				// incapsula_abp_domain whose ID was never added to the list.
				PreConfig: func() { appeared = f.seedDomain(siteId, "appeared.example.com", "") },
				Config:    config,
				PlanOnly:  true,
			},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "domain_ids", domains); err != nil {
						return err
					}
					want := append(append([]string{}, domains...), appeared)
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "effective_domain_ids", want); err != nil {
						return err
					}
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(siteId), want)
				},
			},
		},
	})
}

// Reordering the ranked Domains out of band is drift the configuration can
// express, so it is reported through `domain_ids` and restored.
func TestAbpPriorityDomains_ReorderedRankedDomainsAreRestored(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "reordered", 3)
	config := abpPriorityDomainsConfig(siteId, domains[0], domains[1], domains[2])

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() { f.reorderDomains(domains[2], domains[1], domains[0]) },
				Config:    config,
				Check: func(s *terraform.State) error {
					return abpPriorityCheckIds("the restored order", f.domainOrderOf(siteId), domains)
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A Domain the configuration does not rank, moved ahead of the ranked ones out
// of band, is invisible in `domain_ids` by construction. It still has to be
// planned and corrected, otherwise the drift would be adopted silently.
func TestAbpPriorityDomains_UnrankedDomainAheadIsPushedBack(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "pushed-back", 3)
	config := abpPriorityDomainsConfig(siteId, domains[0], domains[1])

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() { f.reorderDomains(domains[2], domains[0], domains[1]) },
				Config:    config,
				Check: func(s *terraform.State) error {
					return abpPriorityCheckIds("the restored order", f.domainOrderOf(siteId), domains)
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A ranked Domain deleted out of band leaves a configuration that cannot be
// applied as written — in a real configuration the ID comes from an
// incapsula_abp_domain, which is recreated first. The error names the Domain
// rather than reporting a count mismatch.
func TestAbpPriorityDomains_DeletedRankedDomainIsNamed(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "deleted-domain", 3)
	config := abpPriorityDomainsConfig(siteId, domains[0], domains[1])

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig:   func() { f.drop(abpFakeDomains, domains[1]) },
				Config:      config,
				ExpectError: regexp.MustCompile("(?s)not attached to ABP Site"),
			},
		},
	})
}

// The Site itself disappearing drops the order from state, and the recreate
// that follows cannot write an order for a Site that is gone. A real
// configuration recreates the Site first.
func TestAbpPriorityDomains_DeletedSiteIsReported(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "deleted-site", 2)
	config := abpPriorityDomainsConfig(siteId, domains[0])

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig:          func() { f.drop(abpFakeSites, siteId) },
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      config,
				ExpectError: regexp.MustCompile("(?s)not found when setting domain_priority"),
			},
		},
	})
}

// --- validation and plan-time errors ----------------------------------------

// Duplicates are rejected by the provider: the API accepts them and silently
// keeps the last position of each repeated ID.
func TestAbpPriorityDomains_DuplicateIdsAreRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "duplicates", 2)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPriorityDomainsConfig(siteId, domains[0], domains[0]),
				ExpectError: regexp.MustCompile("(?s)is a duplicate of an earlier entry"),
			},
		},
	})

	if count := f.requestCount(http.MethodPut, "domain_priority"); count != 0 {
		t.Fatalf("the rejected configuration still wrote %d order(s)", count)
	}
}

// An ID that names a Domain of another Site is rejected before anything is
// written, by name, instead of letting the API answer with an ID-less count
// mismatch.
func TestAbpPriorityDomains_DomainOfAnotherSiteIsRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "here", 2)
	_, elsewhere := abpPrioritySeedSite(f, "elsewhere", 1)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPriorityDomainsConfig(siteId, domains[0], elsewhere[0]),
				ExpectError: regexp.MustCompile("(?s)not attached to ABP Site"),
			},
		},
	})

	if count := f.requestCount(http.MethodPut, "domain_priority"); count != 0 {
		t.Fatalf("the rejected configuration still wrote %d order(s)", count)
	}
}

// An ID that names nothing at all is the same error, so a typo does not reach
// the API either.
func TestAbpPriorityDomains_UnknownIdIsRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "unknown", 2)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPriorityDomainsConfig(siteId, domains[0], abpFakeAbsentId),
				ExpectError: regexp.MustCompile("(?s)not attached to ABP Site"),
			},
		},
	})
}

// An empty list is rejected when the configuration is validated. A Site with no
// Domains therefore cannot be expressed, which is a limitation rather than an
// intention: the order of nothing is well defined.
func TestAbpPriorityDomains_EmptyListIsRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId, _ := abpPrioritySeedSite(f, "empty", 1)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPriorityDomainsConfig(siteId),
				ExpectError: regexp.MustCompile("(?s)domain_ids requires 1 item minimum"),
			},
		},
	})
}

// Both IDs are validated as UUIDs, so a mistyped reference fails before any
// request is made.
func TestAbpPriorityDomains_NonUuidIdsAreRejected(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "non-uuid", 1)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPriorityDomainsConfig("not-a-uuid", domains[0]),
				ExpectError: regexp.MustCompile(`(?s)expected "site_id" to be a valid UUID`),
			},
			{
				Config:      abpPriorityDomainsConfig(siteId, "not-a-uuid"),
				ExpectError: regexp.MustCompile("(?s)to be a valid UUID"),
			},
		},
	})
}

// --- import and destroy -----------------------------------------------------

// Importing adopts the Site's whole order, so the imported resource re-plans
// empty against a configuration that ranks every Domain.
func TestAbpPriorityDomains_ImportAdoptsTheFullOrder(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "import", 3)
	want := []string{domains[1], domains[2], domains[0]}
	config := abpPriorityDomainsConfig(siteId, want...)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpPriorityDomainsResource,
				ImportState:       true,
				ImportStateId:     siteId,
				ImportStateVerify: true,
			},
		},
	})
}

// Importing a Site whose order was never managed adopts it as it stands, and
// dropping a Domain from the list afterwards leaves it as the lowest priority
// one rather than failing the apply.
func TestAbpPriorityDomains_ImportThenRankASubset(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "import-subset", 3)
	f.reorderDomains(domains[1], domains[2], domains[0])
	imported := []string{domains[1], domains[2], domains[0]}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             abpPriorityDomainsConfig(siteId, imported...),
				ResourceName:       abpPriorityDomainsResource,
				ImportState:        true,
				ImportStateId:      siteId,
				ImportStatePersist: true,
				Check: func(s *terraform.State) error {
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "domain_ids", imported); err != nil {
						return err
					}
					return abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "effective_domain_ids", imported)
				},
			},
			{Config: abpPriorityDomainsConfig(siteId, imported...), PlanOnly: true},
			{
				Config: abpPriorityDomainsConfig(siteId, domains[0]),
				Check: func(s *terraform.State) error {
					want := []string{domains[0], domains[1], domains[2]}
					if err := abpPriorityCheckAttrList(t, s, abpPriorityDomainsResource, "effective_domain_ids", want); err != nil {
						return err
					}
					return abpPriorityCheckIds("the Site's order", f.domainOrderOf(siteId), want)
				},
			},
		},
	})
}

// Importing a Site that does not exist fails rather than producing an empty
// order.
func TestAbpPriorityDomains_ImportOfAbsentSiteFails(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "import-absent", 1)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        abpPriorityDomainsConfig(siteId, domains[0]),
				ResourceName:  abpPriorityDomainsResource,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId,
				ExpectError:   regexp.MustCompile("(?s)not found"),
			},
		},
	})
}

// The order is a property of the Site, not an object of its own: destroying the
// resource only drops it from state, and the Site keeps the order it was given.
func TestAbpPriorityDomains_DestroyLeavesTheOrderInPlace(t *testing.T) {
	f := newAbpFake(t)
	siteId, domains := abpPrioritySeedSite(f, "destroy", 3)
	want := []string{domains[2], domains[0], domains[1]}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if count := f.requestCount(http.MethodDelete, "domain_priority"); count != 0 {
				return fmt.Errorf("destroying the order sent %d DELETE(s)", count)
			}
			return abpPriorityCheckIds("the Site's order after destroy", f.domainOrderOf(siteId), want)
		},
		Steps: []resource.TestStep{
			{Config: abpPriorityDomainsConfig(siteId, domains[2], domains[0])},
		},
	})

	if f.count(abpFakeDomains) != 3 {
		t.Fatalf("destroying the order deleted Domains: %d left of 3", f.count(abpFakeDomains))
	}
}
