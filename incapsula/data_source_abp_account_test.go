package incapsula

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// incapsula_abp_account takes no arguments: it resolves the Account UUID that
// belongs to the configured API credentials and exposes it as both `id` and
// `account_id`, which is how every other ABP resource gets its `account_id`.

const abpAccountDataSource = "data.incapsula_abp_account.current"

// abpAccountConfig is the lookup the sample configuration starts with, feeding
// the resolved UUID into the preflight.
const abpAccountConfig = `
data "incapsula_abp_account" "current" {}

resource "incapsula_abp_preflight" "current" {
  account_id   = data.incapsula_abp_account.current.id
  pending_hash = "hash-1"
}
`

func TestAbpAccount_ResolvesTheAccountUuid(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpAccountConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(abpAccountDataSource, "id", abpFakeAccountId),
					resource.TestCheckResourceAttr(abpAccountDataSource, "account_id", abpFakeAccountId),
					// The resolved UUID is usable wherever an account_id is
					// required.
					resource.TestCheckResourceAttr(abpPreflightResource, "account_id", abpFakeAccountId),
				),
			},
			{
				Config:   abpAccountConfig,
				PlanOnly: true,
			},
		},
	})
}

// The lookup is a read of the API root, repeated on every refresh, and it never
// writes anything.
func TestAbpAccount_IsReadOnlyAndRefreshed(t *testing.T) {
	f := newAbpFake(t)
	config := `data "incapsula_abp_account" "current" {}`

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if _, ok := abpAttrs(t, s, abpAccountDataSource)["account_id"]; !ok {
						return fmt.Errorf("account_id was not set")
					}
					return nil
				},
			},
			{Config: config},
		},
	})

	if reads := f.requestCount(http.MethodGet, "/botmanagement/v1/"); reads < 2 {
		t.Errorf("expected the API root to be read on every refresh, got %d reads", reads)
	}
	for _, write := range f.writes() {
		// Only the credential check the provider does when it builds its
		// client is expected.
		if !strings.HasSuffix(write.Path, "/account/verify") {
			t.Errorf("the lookup sent %s %s", write.Method, write.Path)
		}
	}
}

// A service that cannot answer the API root has to be reported with what was
// being read, rather than as an empty Account UUID that later resources would
// reject as an invalid reference.
func TestAbpAccount_ServiceErrorIsReported(t *testing.T) {
	_, err := abpAccountBrokenClient(t, http.StatusBadRequest, `{"message":"down"}`).ReadAbpAccountId()
	if err == nil {
		t.Fatal("a failing API root was reported as success")
	}
	for _, want := range []string{"400", "ABP Account ID", "down"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("the error %q does not mention %q", err.Error(), want)
		}
	}

	_, err = abpAccountBrokenClient(t, http.StatusOK, `not json`).ReadAbpAccountId()
	if err == nil {
		t.Fatal("an unparsable API root was reported as success")
	}
	if !strings.Contains(err.Error(), "error parsing ABP Account ID") {
		t.Errorf("the error %q does not say what could not be parsed", err.Error())
	}
}
