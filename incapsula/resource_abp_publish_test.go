package incapsula

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// These tests cover incapsula_abp_publish on its own: what it does with the
// snapshot it is given, what it reports when that snapshot can no longer be
// used, and what it does not do — it neither tracks the outcome of the publish
// nor undoes it.

const abpPublishResource = "incapsula_abp_publish.publish"

func abpPublishConfig(preflightId string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_publish" "publish" {
  preflight_id = %q
}
`, preflightId)
}

// --- create, re-apply, replacement ------------------------------------------

// The resource's id is the id of the published snapshot, not of the publish
// itself: nothing in state points at the publish record.
func TestAbpPublish_PublishesTheGivenSnapshot(t *testing.T) {
	f := newAbpFake(t)
	preflightId := f.abpPreflightSeedSnapshot()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPublishConfig(preflightId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPublishResource, "id", preflightId),
					resource.TestCheckResourceAttr(abpPublishResource, "preflight_id", preflightId),
					func(s *terraform.State) error {
						if published := f.publishedIds(); len(published) != 1 || published[0] != preflightId {
							return fmt.Errorf("expected %s to have been published, got %v", preflightId, published)
						}
						return nil
					},
				),
			},
			{
				Config:   abpPublishConfig(preflightId),
				PlanOnly: true,
			},
		},
	})

	// A re-apply must not publish again.
	if count := f.requestCount(http.MethodPost, "/publish"); count != 1 {
		t.Errorf("expected one publish request, got %d", count)
	}
}

// Pointing the resource at another snapshot is how a further publish is
// triggered: the previous publish is forgotten rather than undone.
func TestAbpPublish_SnapshotChangeRepublishes(t *testing.T) {
	f := newAbpFake(t)
	first := f.abpPreflightSeedSnapshot()
	second := f.abpPreflightSeedSnapshot()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPublishConfig(first)},
			{
				Config: abpPublishConfig(second),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPublishResource, "id", second),
					func(s *terraform.State) error {
						if published := f.publishedIds(); len(published) != 2 ||
							published[0] != first || published[1] != second {
							return fmt.Errorf("expected %s then %s to have been published, got %v",
								first, second, published)
						}
						if deleted := f.deletedIds(); len(deleted) != 0 {
							return fmt.Errorf("the replacement asked the API to delete %v", deleted)
						}
						return nil
					},
				),
			},
		},
	})
}

// The publish keeps no state of its own and never re-reads anything, so a
// snapshot that disappears afterwards, or an Account published by somebody
// else, is invisible here. That is a limitation rather than an intention: the
// publish is asynchronous and its outcome is not tracked or retried.
func TestAbpPublish_DoesNotTrackTheOutcome(t *testing.T) {
	f := newAbpFake(t)
	preflightId := f.abpPreflightSeedSnapshot()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPublishConfig(preflightId)},
			{
				PreConfig: func() {
					f.publishElsewhere()
					f.drop(abpFakePreflights, preflightId)
				},
				Config:   abpPublishConfig(preflightId),
				PlanOnly: true,
				Check:    resource.TestCheckResourceAttr(abpPublishResource, "id", preflightId),
			},
		},
	})
}

// --- errors -----------------------------------------------------------------

func TestAbpPublish_InvalidSnapshotIdIsRejected(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPublishConfig("not-a-uuid"),
				ExpectError: regexp.MustCompile(`expected "preflight_id" to be a valid UUID`),
			},
		},
	})
}

// A snapshot that is gone and one that was overtaken fail differently, and both
// have to say that running again recovers.
func TestAbpPublish_MissingSnapshotIsReported(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPublishConfig(abpFakeAbsentId),
				ExpectError: regexp.MustCompile("no longer exists"),
			},
		},
	})

	if published := f.publishedIds(); len(published) != 0 {
		t.Errorf("something was published anyway: %v", published)
	}
}

func TestAbpPublish_OvertakenSnapshotIsReported(t *testing.T) {
	f := newAbpFake(t)
	preflightId := f.abpPreflightSeedSnapshot()
	f.publishElsewhere()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPublishConfig(preflightId),
				ExpectError: regexp.MustCompile("is not publishable"),
			},
		},
	})

	// The pre-check keeps the request the API would have rejected from being
	// sent at all.
	if count := f.requestCount(http.MethodPost, "/publish"); count != 0 {
		t.Errorf("expected no publish request, got %d", count)
	}
}

// A publish the service refuses for any other reason is reported with the
// action that failed.
func TestAbpPublish_ServiceFailureIsReported(t *testing.T) {
	f := newAbpFake(t)
	preflightId := f.abpPreflightSeedSnapshot()
	f.setRejectPublish(true)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpPublishConfig(preflightId),
				ExpectError: regexp.MustCompile("publishing preflight"),
			},
			// The failure was transient, so the same configuration succeeds.
			{
				Config: abpPublishConfig(preflightId),
				Check:  resource.TestCheckResourceAttr(abpPublishResource, "id", preflightId),
			},
		},
	})
}

// --- destroy and import -----------------------------------------------------

// Destroying the publish drops it from state and nothing else: the Account
// keeps running what was published.
func TestAbpPublish_DestroyUndoesNothing(t *testing.T) {
	f := newAbpFake(t)
	preflightId := f.abpPreflightSeedSnapshot()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps:             []resource.TestStep{{Config: abpPublishConfig(preflightId)}},
		CheckDestroy: func(s *terraform.State) error {
			if published := f.publishedIds(); len(published) != 1 || published[0] != preflightId {
				return fmt.Errorf("the destroy changed what is published: %v", published)
			}
			if deleted := f.deletedIds(); len(deleted) != 0 {
				return fmt.Errorf("the destroy asked the API to delete %v", deleted)
			}
			return nil
		},
	})
}

func TestAbpPublish_ImportIsNotSupported(t *testing.T) {
	f := newAbpFake(t)
	preflightId := f.abpPreflightSeedSnapshot()

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: abpPublishConfig(preflightId)},
			{
				Config:            abpPublishConfig(preflightId),
				ResourceName:      abpPublishResource,
				ImportState:       true,
				ImportStateVerify: true,
				ExpectError:       regexp.MustCompile("doesn't support import"),
			},
		},
	})
}

// --- client level -----------------------------------------------------------

// A snapshot can still become unusable between the refresh and the publish, so
// the two failure modes have to be told apart in the message the user sees.
func TestAbpPublish_UnusableSnapshotErrorsAreSpecific(t *testing.T) {
	f := newAbpFake(t)
	client := f.client()

	publish := func(preflightId string) string {
		data := schema.TestResourceDataRaw(t, resourceAbpPublish().Schema, map[string]any{
			"preflight_id": preflightId,
		})
		diags := resourceAbpPublishCreate(context.Background(), data, client)
		if !diags.HasError() {
			t.Fatalf("publishing %s unexpectedly succeeded", preflightId)
		}
		return diags[0].Summary + ": " + diags[0].Detail
	}

	discarded, diags := client.CreateAbpPreflight(abpFakeAccountId)
	if diags.HasError() {
		t.Fatalf("could not take a snapshot: %v", diags)
	}
	f.drop(abpFakePreflights, discarded.Id)
	if message := publish(discarded.Id); !regexp.MustCompile("no longer exists.*apply` again").MatchString(message) {
		t.Errorf("publishing a discarded snapshot reported %q", message)
	}

	overtaken, diags := client.CreateAbpPreflight(abpFakeAccountId)
	if diags.HasError() {
		t.Fatalf("could not take a snapshot: %v", diags)
	}
	f.publishElsewhere()
	if message := publish(overtaken.Id); !regexp.MustCompile("not publishable.*apply` again").MatchString(message) {
		t.Errorf("publishing an overtaken snapshot reported %q", message)
	}
}
