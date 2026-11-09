package incapsula

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These tests cover incapsula_abp_preflight: the snapshot of the Account
// configuration that incapsula_abp_publish activates. What matters most is what
// happens to a snapshot that can no longer be used — it was discarded, the
// Account was published from elsewhere, or the publish itself failed — because
// such a snapshot used to wedge every later apply instead of being replaced.

const abpPreflightResource = "incapsula_abp_preflight.current"

// abpPreflightPipelineConfig is the pipeline the sample configuration wires up.
var abpPreflightPipelineConfig = fmt.Sprintf(`
data "incapsula_abp_pending_changes" "current" {}

resource "incapsula_abp_preflight" "current" {
  account_id   = %q
  pending_hash = data.incapsula_abp_pending_changes.current.hash
}

resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
`, abpFakeAccountId)

// abpPreflightLiteralConfig is the same pipeline with the change-detection
// value written out, so a test can decide when a new snapshot is taken.
func abpPreflightLiteralConfig(accountId, pendingHash string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_preflight" "current" {
  account_id   = %q
  pending_hash = %q
}

resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
`, accountId, pendingHash)
}

// abpPreflightCheckPublished asserts that the snapshot in state is the one the
// service last published, and that the publish points at it.
func abpPreflightCheckPublished(t *testing.T, f *abpFake) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		preflightId := abpAttrs(t, s, abpPreflightResource)["id"]
		published := f.publishedIds()
		if len(published) == 0 {
			return fmt.Errorf("nothing was published, the snapshot in state is %s", preflightId)
		}
		if latest := published[len(published)-1]; latest != preflightId {
			return fmt.Errorf("the published snapshot is %s but state holds %s", latest, preflightId)
		}
		if publishId := abpAttrs(t, s, abpPublishResource)["preflight_id"]; publishId != preflightId {
			return fmt.Errorf("the publish points at %s but state holds %s", publishId, preflightId)
		}
		return nil
	}
}

// abpPreflightCapture remembers the snapshot id a step ended up with.
func abpPreflightCapture(t *testing.T, into *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*into = abpAttrs(t, s, abpPreflightResource)["id"]
		return nil
	}
}

// abpPreflightCheckReplaced fails unless the snapshot in state differs from the
// remembered one.
func abpPreflightCheckReplaced(t *testing.T, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if got := abpAttrs(t, s, abpPreflightResource)["id"]; got == *previous {
			return fmt.Errorf("the snapshot %s was kept", got)
		}
		return nil
	}
}

// abpPreflightCheckKept fails unless the snapshot in state is the remembered one.
func abpPreflightCheckKept(t *testing.T, previous *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if got := abpAttrs(t, s, abpPreflightResource)["id"]; got != *previous {
			return fmt.Errorf("the snapshot %s was replaced by %s", *previous, got)
		}
		return nil
	}
}

// --- create, re-apply, replacement ------------------------------------------

func TestAbpPreflight_CreateThenPlanNothing(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					resource.TestCheckResourceAttr(abpPreflightResource, "account_id", abpFakeAccountId),
					resource.TestCheckResourceAttrSet(abpPreflightResource, "created_at"),
					resource.TestCheckResourceAttr(abpPreflightResource, "warnings.#", "0"),
				),
			},
			{
				Config:   abpPreflightPipelineConfig,
				PlanOnly: true,
			},
		},
	})

	// An unchanged configuration takes one snapshot and publishes it once.
	if count := f.preflightCount(); count != 1 {
		t.Errorf("expected one snapshot to have been taken, got %d", count)
	}
	if published := f.publishedIds(); len(published) != 1 {
		t.Errorf("expected one publish, got %d", len(published))
	}
}

// pending_hash is the whole change-detection mechanism: a new value takes a new
// snapshot and publishes it, the same value does nothing.
func TestAbpPreflight_PendingHashDrivesReplacement(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightLiteralConfig(abpFakeAccountId, "hash-1"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					abpPreflightCapture(t, &first),
				),
			},
			{
				Config:   abpPreflightLiteralConfig(abpFakeAccountId, "hash-1"),
				PlanOnly: true,
			},
			{
				Config: abpPreflightLiteralConfig(abpFakeAccountId, "hash-2"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckReplaced(t, &first),
					abpPreflightCheckPublished(t, f),
					resource.TestCheckResourceAttr(abpPreflightResource, "pending_hash", "hash-2"),
					func(s *terraform.State) error {
						// Replacing the snapshot replaced the publish with it.
						if published := f.publishedIds(); len(published) != 2 {
							return fmt.Errorf("expected two publishes, got %d", len(published))
						}
						// The snapshot Terraform forgot is left on the server.
						if !f.exists(abpFakePreflights, first) {
							return fmt.Errorf("the previous snapshot %s was deleted", first)
						}
						return nil
					},
				),
			},
		},
	})
}

// A snapshot belongs to one Account, so pointing the resource at another one
// replaces it and leaves the old snapshot alone.
func TestAbpPreflight_AccountIdChangeReplaces(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightLiteralConfig(abpFakeAccountId, "hash-1"),
				Check:  abpPreflightCapture(t, &first),
			},
			{
				Config: abpPreflightLiteralConfig(abpFakeOtherAccountId, "hash-1"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckReplaced(t, &first),
					resource.TestCheckResourceAttr(abpPreflightResource, "account_id", abpFakeOtherAccountId),
					func(s *terraform.State) error {
						if !f.exists(abpFakePreflights, first) {
							return fmt.Errorf("the previous snapshot %s was deleted", first)
						}
						if deleted := f.deletedIds(); len(deleted) != 0 {
							return fmt.Errorf("the API was asked to delete %v", deleted)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestAbpPreflight_InvalidAccountIdIsRejected(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPreflightLiteralConfig("not-a-uuid", "hash-1"),
				ExpectError: regexp.MustCompile(`expected "account_id" to be a valid UUID`),
			},
		},
	})
}

// --- a snapshot that can no longer be used ----------------------------------

// A snapshot overtaken by a publish from elsewhere can never be published
// again, so the refresh has to drop it and the next apply has to take a new one
// — without -replace and without any manual step.
func TestAbpPreflight_UnpublishableSnapshotIsReplaced(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					abpPreflightCapture(t, &first),
				),
			},
			{
				PreConfig: f.publishElsewhere,
				Config:    abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckReplaced(t, &first),
					abpPreflightCheckPublished(t, f),
				),
			},
		},
	})
}

// A snapshot the service has discarded as it aged must not keep looking
// healthy: the refresh has to notice it is gone and plan a new one.
func TestAbpPreflight_DiscardedSnapshotIsReplaced(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightPipelineConfig,
				Check:  abpPreflightCapture(t, &first),
			},
			{
				PreConfig: func() { f.drop(abpFakePreflights, first) },
				Config:    abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckReplaced(t, &first),
					abpPreflightCheckPublished(t, f),
				),
			},
		},
	})
}

// The wedge the recovery was written for: an apply takes a snapshot, the
// publish fails and never reaches state, then someone else publishes. The next
// apply has to recover on its own.
func TestAbpPreflight_FailedPublishRecoversOnNextApply(t *testing.T) {
	f := newAbpFake(t)
	f.setRejectPublish(true)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPreflightPipelineConfig,
				ExpectError: regexp.MustCompile("publishing preflight"),
			},
			{
				// Another operator publishes while this configuration is
				// holding a snapshot it never managed to publish.
				PreConfig: f.publishElsewhere,
				Config:    abpPreflightPipelineConfig,
				Check:     abpPreflightCheckPublished(t, f),
			},
		},
	})
}

// The refresh adopts what the service says the snapshot belongs to, so a
// snapshot that turns out to be for another Account is replaced rather than
// silently kept.
func TestAbpPreflight_AccountDriftPlansReplacement(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightLiteralConfig(abpFakeAccountId, "hash-1"),
				Check:  abpPreflightCapture(t, &first),
			},
			{
				PreConfig: func() {
					f.mutate(t, abpFakePreflights, first, func(doc map[string]any) {
						doc["account_id"] = abpFakeOtherAccountId
					})
				},
				Config:             abpPreflightLiteralConfig(abpFakeAccountId, "hash-1"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: abpPreflightLiteralConfig(abpFakeAccountId, "hash-1"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckReplaced(t, &first),
					resource.TestCheckResourceAttr(abpPreflightResource, "account_id", abpFakeAccountId),
				),
			},
		},
	})
}

// --- a snapshot that is merely outdated -------------------------------------
// An outdated snapshot is still publishable, so it is kept on purpose:
// replacing it would publish changes Terraform does not manage, and on a busy
// Account it would never stop replacing. The reason is exposed instead.
func TestAbpPreflight_OutdatedSnapshotIsKept(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPreflightResource, "warnings.#", "0"),
					abpPreflightCapture(t, &first),
				),
			},
			{
				PreConfig: func() { f.setPreflightOutdated(true) },
				Config:    abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckKept(t, &first),
					resource.TestCheckResourceAttr(abpPreflightResource, "warnings.#", "1"),
					resource.TestCheckResourceAttr(abpPreflightResource, "warnings.0.type", "outdated"),
					resource.TestCheckResourceAttrSet(abpPreflightResource, "warnings.0.message"),
					func(s *terraform.State) error {
						if count := f.preflightCount(); count != 1 {
							return fmt.Errorf("expected one snapshot to have been taken, got %d", count)
						}
						return nil
					},
				),
			},
		},
	})
}

// The service reports no warnings at all when it cannot work them out, which
// must not be mistaken for an unusable snapshot.
func TestAbpPreflight_MissingWarningsAreNotDrift(t *testing.T) {
	f := newAbpFake(t)
	f.setPreflightNoWarnings(true)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					resource.TestCheckResourceAttr(abpPreflightResource, "warnings.#", "0"),
				),
			},
			{
				Config:   abpPreflightPipelineConfig,
				PlanOnly: true,
			},
		},
	})
}

// --- the cost of omitting pending_hash --------------------------------------

// Without pending_hash the snapshot is taken once and nothing in the
// configuration ever triggers a new one, so a changed ABP resource is never
// published. This is a documented foot-gun rather than an intention: the
// attribute is Optional, and leaving it out is what makes the pipeline silent.
func TestAbpPreflight_WithoutPendingHashNothingTriggersANewSnapshot(t *testing.T) {
	f := newAbpFake(t)
	var first string

	config := func(description string) string {
		return fmt.Sprintf(`
resource "incapsula_abp_condition" "tracked" {
  account_id  = %[1]q
  name        = "tracked"
  description = %[2]q
  code        = "(all flags.no_token)"
}

data "incapsula_abp_pending_changes" "current" {
  depends_on = [incapsula_abp_condition.tracked]
}

resource "incapsula_abp_preflight" "current" {
  account_id = %[1]q
}

resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
`, abpFakeAccountId, description)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config("before"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					resource.TestCheckNoResourceAttr(abpPreflightResource, "pending_hash"),
					abpPreflightCapture(t, &first),
				),
			},
			{
				Config: config("after"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckKept(t, &first),
					func(s *terraform.State) error {
						if published := f.publishedIds(); len(published) != 1 {
							return fmt.Errorf("expected the single publish of %s, got %v", first, published)
						}
						return nil
					},
				),
			},
		},
	})
}

// --- destroy and import -----------------------------------------------------

// Destroying the pipeline only drops the snapshot from state: nothing is
// deleted and nothing is published, so the Account keeps running the
// configuration that was published last.
func TestAbpPreflight_DestroyDropsStateOnly(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPreflightPipelineConfig,
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					abpPreflightCapture(t, &first),
				),
			},
		},
		CheckDestroy: func(s *terraform.State) error {
			if !f.exists(abpFakePreflights, first) {
				return fmt.Errorf("snapshot %s was deleted by the destroy", first)
			}
			if deleted := f.deletedIds(); len(deleted) != 0 {
				return fmt.Errorf("the destroy asked the API to delete %v", deleted)
			}
			if published := f.publishedIds(); len(published) != 1 {
				return fmt.Errorf("the destroy published something: %v", published)
			}
			return nil
		},
	})
}

// Importing a snapshot is not supported: it is a point-in-time object, and an
// imported one carries no publish.
func TestAbpPreflight_ImportIsNotSupported(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPreflightPipelineConfig},
			{
				Config:            abpPreflightPipelineConfig,
				ResourceName:      abpPreflightResource,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("doesn't support import"),
			},
		},
	})
}

// --- client level -----------------------------------------------------------

// A read that fails for any other reason must not be mistaken for a discarded
// snapshot: only a 404 means the snapshot is gone, everything else has to be
// reported, or the resource would quietly drop a snapshot that still exists.
func TestAbpPreflightClient_ServiceErrorIsReported(t *testing.T) {
	client := abpAccountBrokenClient(t, http.StatusBadRequest, `{"message":"broken"}`)

	preflight, diags := client.GetAbpPreflight(abpFakeAbsentId)
	if !diags.HasError() {
		t.Errorf("a failing read was reported as the snapshot being gone: %+v", preflight)
	}
	if _, diags := client.CreateAbpPreflight(abpFakeAccountId); !diags.HasError() {
		t.Error("a failing create was reported as success")
	}
	if diags := client.PublishAbpPreflight(abpFakeAbsentId); !diags.HasError() {
		t.Error("a failing publish was reported as success")
	}
}

// GetAbpPreflight has to report a snapshot that is gone as an ordinary outcome
// rather than an error, because that is how the resource decides to take a new
// one, and it has to derive publishability from the publish transition the API
// offers.
func TestAbpPreflightClient_ReportsAbsenceAndPublishability(t *testing.T) {
	f := newAbpFake(t)
	client := f.client()

	absent, diags := client.GetAbpPreflight(abpFakeAbsentId)
	if diags.HasError() {
		t.Fatalf("reading an absent snapshot errored: %v", diags)
	}
	if absent != nil {
		t.Errorf("reading an absent snapshot returned %+v", absent)
	}

	created, diags := client.CreateAbpPreflight(abpFakeAccountId)
	if diags.HasError() {
		t.Fatalf("could not take a snapshot: %v", diags)
	}
	if created.AccountId != abpFakeAccountId || created.CreatedAt == "" {
		t.Errorf("a fresh snapshot came back as %+v", created)
	}
	if !created.CanPublish() {
		t.Error("a fresh snapshot is not publishable")
	}
	if len(created.Warnings) != 0 {
		t.Errorf("a fresh snapshot carries warnings %+v", created.Warnings)
	}

	if diags := client.PublishAbpPreflight(created.Id); diags.HasError() {
		t.Fatalf("publishing errored: %v", diags)
	}

	// A snapshot taken before the last publish stays unpublishable forever, and
	// says why.
	f.publishElsewhere()
	overtaken, diags := client.GetAbpPreflight(created.Id)
	if diags.HasError() {
		t.Fatalf("re-reading the snapshot errored: %v", diags)
	}
	if overtaken.CanPublish() {
		t.Error("an overtaken snapshot still claims to be publishable")
	}
	if len(overtaken.Warnings) != 1 || overtaken.Warnings[0].Type != "older_than_most_recent_published" {
		t.Errorf("an overtaken snapshot reports warnings %+v", overtaken.Warnings)
	}
}
