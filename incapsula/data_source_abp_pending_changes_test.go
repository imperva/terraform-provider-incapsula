package incapsula

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// incapsula_abp_pending_changes is a marker, not a lookup: it takes no
// arguments, makes no API call and its `hash` is a constant. It only becomes a
// trigger because Terraform postpones reading a data source whose `depends_on`
// has a planned change, which leaves `pending_hash` unknown and so replaces the
// preflight that consumes it. These tests pin down both halves: the constant,
// and the triggering it enables.

const (
	abpPendingChangesDataSource = "data.incapsula_abp_pending_changes.current"
	abpPendingChangesHash       = "incapsula_abp_pending_changes"
)

// abpPendingChangesPipelineConfig is the wiring of the sample configuration:
// the tracked resources, the data source that depends on them, the preflight
// and the publish. The untracked Condition stands for a resource left outside
// the `depends_on` graph.
func abpPendingChangesPipelineConfig(tracked, untracked string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_condition" "tracked" {
  account_id  = %[1]q
  name        = "tracked"
  description = %[2]q
  code        = "(all flags.no_token)"
}

resource "incapsula_abp_condition" "untracked" {
  account_id  = %[1]q
  name        = "untracked"
  description = %[3]q
  code        = "(all flags.no_token)"
}

data "incapsula_abp_pending_changes" "current" {
  depends_on = [incapsula_abp_condition.tracked]
}

resource "incapsula_abp_preflight" "current" {
  account_id   = %[1]q
  pending_hash = data.incapsula_abp_pending_changes.current.hash
}

resource "incapsula_abp_publish" "publish" {
  preflight_id = incapsula_abp_preflight.current.id
}
`, abpFakeAccountId, tracked, untracked)
}

// The marker carries no information about the Account: every instance has the
// same value, and reading it talks to nothing.
func TestAbpPendingChanges_HashIsConstantAndReadsNothing(t *testing.T) {
	f := newAbpFake(t)
	config := `
data "incapsula_abp_pending_changes" "current" {}

data "incapsula_abp_pending_changes" "second" {}
`

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpPendingChangesDataSource, "hash", abpPendingChangesHash),
					resource.TestCheckResourceAttr(abpPendingChangesDataSource, "id", abpPendingChangesHash),
					resource.TestCheckResourceAttrPair(
						abpPendingChangesDataSource, "hash",
						"data.incapsula_abp_pending_changes.second", "hash"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})

	if reads := f.requestCount(http.MethodGet, "botmanagement"); reads != 0 {
		t.Errorf("the data source made %d ABP API requests", reads)
	}
}

// The pipeline end to end: a change to a tracked ABP resource takes a new
// snapshot and publishes it, and a re-apply that changes nothing does neither.
func TestAbpPendingChanges_TrackedChangeTriggersPreflightAndPublish(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPendingChangesPipelineConfig("before", "before"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					resource.TestCheckResourceAttr(abpPendingChangesDataSource, "hash", abpPendingChangesHash),
					abpPreflightCapture(t, &first),
				),
			},
			{
				Config:   abpPendingChangesPipelineConfig("before", "before"),
				PlanOnly: true,
			},
			{
				Config: abpPendingChangesPipelineConfig("after", "before"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckReplaced(t, &first),
					abpPreflightCheckPublished(t, f),
					// The marker itself never changes: the trigger is the
					// postponed read, not a new value.
					resource.TestCheckResourceAttr(abpPendingChangesDataSource, "hash", abpPendingChangesHash),
					func(s *terraform.State) error {
						if count := f.preflightCount(); count != 2 {
							return fmt.Errorf("expected two snapshots to have been taken, got %d", count)
						}
						if published := f.publishedIds(); len(published) != 2 {
							return fmt.Errorf("expected two publishes, got %v", published)
						}
						return nil
					},
				),
			},
		},
	})
}

// A resource outside the `depends_on` graph does not trigger a publish of its
// own: the change is applied and then stays pending until something else
// publishes the Account. This is the documented consequence of the design, and
// the reason every ABP resource has to be reachable through `depends_on`.
func TestAbpPendingChanges_UntrackedChangeTriggersNothing(t *testing.T) {
	f := newAbpFake(t)
	var first string

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpPendingChangesPipelineConfig("before", "before"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckPublished(t, f),
					abpPreflightCapture(t, &first),
				),
			},
			{
				Config: abpPendingChangesPipelineConfig("before", "after"),
				Check: resource.ComposeTestCheckFunc(
					abpPreflightCheckKept(t, &first),
					func(s *terraform.State) error {
						// The Condition was updated...
						for _, condition := range f.all(abpFakeConditions) {
							if abpFakeString(condition["name"]) == "untracked" &&
								abpFakeString(condition["description"]) != "after" {
								return fmt.Errorf("the untracked Condition was not updated: %v", condition)
							}
						}
						// ... but nothing was snapshotted or published for it.
						if count := f.preflightCount(); count != 1 {
							return fmt.Errorf("expected one snapshot to have been taken, got %d", count)
						}
						if published := f.publishedIds(); len(published) != 1 {
							return fmt.Errorf("expected one publish, got %v", published)
						}
						return nil
					},
				),
			},
		},
	})
}
