package incapsula

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Tests for incapsula_abp_domain_encryption_key. The key material is immutable,
// so every change replaces the resource, and the order in which the replacement
// happens matters: the Account configuration has to be published before the
// analysis host accepts a new key, so the old key must outlive the new one's
// creation.

const abpDomainKeyTestAddress = "incapsula_abp_domain_encryption_key.test"

// abpDomainKeyConfig renders a Domain and a key on it, with the key's body
// given verbatim.
func abpDomainKeyConfig(siteId, fqdn, body string) string {
	return abpDomainMinimalConfig(siteId, fqdn) + fmt.Sprintf(`
resource "incapsula_abp_domain_encryption_key" "test" {
  domain_id = incapsula_abp_domain.test.id
%s
}
`, body)
}

// abpDomainKeyOnConfig renders a key on a Domain that already exists.
func abpDomainKeyOnConfig(domainId, key string) string {
	return fmt.Sprintf(`
resource "incapsula_abp_domain_encryption_key" "test" {
  domain_id = %q
  key       = %q
}
`, domainId, key)
}

// --- create -----------------------------------------------------------------

// A key is created with the material the configuration gives, and is not
// published yet.
func TestAbpDomainEncryptionKey_Create(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-create")
	config := abpDomainKeyConfig(siteId, "key-create.example.com", fmt.Sprintf(`  key = %q`, abpDomainTestKeyMaterial))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainKeyTestAddress, "key", abpDomainTestKeyMaterial),
					resource.TestCheckResourceAttrSet(abpDomainKeyTestAddress, "created_at"),
					// Empty until the Account configuration is published.
					resource.TestCheckResourceAttr(abpDomainKeyTestAddress, "first_published_at", ""),
					resource.TestCheckResourceAttrPair(abpDomainKeyTestAddress, "domain_id", abpDomainTestAddress, "id"),
					func(s *terraform.State) error {
						stored := f.only(t, abpFakeKeys)
						if got := abpFakeString(stored["key"]); got != abpDomainTestKeyMaterial {
							return fmt.Errorf("stored key material is %q", got)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// `first_published_at` is filled in by the service once the configuration has
// been published, and the provider reports it without planning a change.
func TestAbpDomainEncryptionKey_FirstPublishedAtIsAdoptedWithoutADiff(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-published")
	config := abpDomainKeyConfig(siteId, "key-published.example.com", fmt.Sprintf(`  key = %q`, abpDomainTestKeyMaterial))

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					f.mutate(t, abpFakeKeys, abpAttrs(t, s, abpDomainKeyTestAddress)["id"], func(doc map[string]any) {
						doc["first_published_at"] = "2026-02-03T04:05:06Z"
					})
					return nil
				},
			},
			{
				Config: config,
				Check:  resource.TestCheckResourceAttr(abpDomainKeyTestAddress, "first_published_at", "2026-02-03T04:05:06Z"),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// --- replacement ------------------------------------------------------------

// The key material cannot be modified, so changing it replaces the resource.
// By default that deletes the old key before creating the new one, which leaves
// the Domain without a usable key in between.
func TestAbpDomainEncryptionKey_KeyChangeReplacesDestroyFirst(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-rotate")

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpDomainKeyConfig(siteId, "key-rotate.example.com", fmt.Sprintf(`  key = %q`, abpDomainTestKeyMaterial)),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpDomainKeyConfig(siteId, "key-rotate.example.com", fmt.Sprintf(`  key = %q`, abpDomainTestOtherKeyMaterial)),
				Check: func(s *terraform.State) error {
					id := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					if id == created {
						return fmt.Errorf("changing the key material did not replace key %s", created)
					}
					if f.exists(abpFakeKeys, created) {
						return fmt.Errorf("the replaced key %s is still stored", created)
					}
					if deleted := f.deletedIds(); !slices.Equal(deleted, []string{created}) {
						return fmt.Errorf("the replacement deleted %v, want just %s", deleted, created)
					}
					want := []string{http.MethodPost, http.MethodDelete, http.MethodPost}
					if got := abpDomainWritesTo(f, "encryptionkey"); !slices.Equal(got, want) {
						return fmt.Errorf("the key endpoints saw %v, want %v", got, want)
					}
					return nil
				},
			},
		},
	})
}

// `create_before_destroy` is what the rotation procedure needs: the new key is
// created first and the previous one stays alive until the publish completes.
func TestAbpDomainEncryptionKey_CreateBeforeDestroyRotation(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-cbd")
	config := func(key string) string {
		return abpDomainKeyConfig(siteId, "key-cbd.example.com", fmt.Sprintf(`  key = %q

  lifecycle {
    create_before_destroy = true
  }`, key))
	}

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config(abpDomainTestKeyMaterial),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
			{
				Config: config(abpDomainTestOtherKeyMaterial),
				Check: func(s *terraform.State) error {
					if id := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]; id == created {
						return fmt.Errorf("changing the key material did not replace key %s", created)
					}
					want := []string{http.MethodPost, http.MethodPost, http.MethodDelete}
					if got := abpDomainWritesTo(f, "encryptionkey"); !slices.Equal(got, want) {
						return fmt.Errorf("the key endpoints saw %v, want the new key created before the old one is deleted (%v)", got, want)
					}
					return nil
				},
			},
		},
	})
}

// A key belongs to one Domain, so pointing it at another one replaces it.
func TestAbpDomainEncryptionKey_DomainIdChangeReplaces(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-move")
	first := f.seedDomain(siteId, "first.example.com", "")
	second := f.seedDomain(siteId, "second.example.com", "")

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: abpDomainKeyOnConfig(first, abpDomainTestKeyMaterial),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
			{
				Config: abpDomainKeyOnConfig(second, abpDomainTestKeyMaterial),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainKeyTestAddress, "domain_id", second),
					func(s *terraform.State) error {
						if id := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]; id == created {
							return fmt.Errorf("changing domain_id did not replace key %s", created)
						}
						if deleted := f.deletedIds(); !slices.Equal(deleted, []string{created}) {
							return fmt.Errorf("the replacement deleted %v, want just %s", deleted, created)
						}
						return nil
					},
				),
			},
		},
	})
}

// --- drift and delete -------------------------------------------------------

// A key deleted outside Terraform is recreated.
func TestAbpDomainEncryptionKey_DeletedOutOfBandIsRecreated(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-vanished")
	domainId := f.seedDomain(siteId, "key-vanished.example.com", "")
	config := abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial)

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
			{
				PreConfig: func() { f.drop(abpFakeKeys, created) },
				Config:    config,
				Check: func(s *terraform.State) error {
					id := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					if id == created {
						return fmt.Errorf("the vanished key %s was not recreated", created)
					}
					if !f.exists(abpFakeKeys, id) {
						return fmt.Errorf("key %s is not stored", id)
					}
					return nil
				},
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// A second key on the same Domain is not mistaken for the managed one: the read
// has to pick the key out of the Domain's list by id.
func TestAbpDomainEncryptionKey_OtherKeyOnTheDomainIsIgnored(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-siblings")
	domainId := f.seedDomain(siteId, "key-siblings.example.com", "")
	sibling := f.seedEncryptionKey(domainId, abpDomainTestOtherKeyMaterial)
	config := abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(abpDomainKeyTestAddress, "key", abpDomainTestKeyMaterial),
					func(s *terraform.State) error {
						if id := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]; id == sibling {
							return fmt.Errorf("the resource adopted the key %s that it does not manage", sibling)
						}
						return nil
					},
				),
			},
			{Config: config, PlanOnly: true},
		},
	})
}

// Destroying the key deletes it through the key endpoint and leaves the Domain
// alone.
func TestAbpDomainEncryptionKey_Destroy(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-destroy")
	domainId := f.seedDomain(siteId, "key-destroy.example.com", "")

	var created string
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		CheckDestroy: func(s *terraform.State) error {
			if f.exists(abpFakeKeys, created) {
				return fmt.Errorf("key %s survived the destroy", created)
			}
			if !f.exists(abpFakeDomains, domainId) {
				return fmt.Errorf("destroying the key deleted Domain %s", domainId)
			}
			if deleted := f.deletedIds(); !slices.Equal(deleted, []string{created}) {
				return fmt.Errorf("the destroy deleted %v, want just the key %s", deleted, created)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial),
				Check: func(s *terraform.State) error {
					created = abpAttrs(t, s, abpDomainKeyTestAddress)["id"]
					return nil
				},
			},
		},
	})
}

// --- import -----------------------------------------------------------------

// The read endpoint needs the parent Domain, so the import id carries both.
func TestAbpDomainEncryptionKey_Import(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-import")
	domainId := f.seedDomain(siteId, "key-import.example.com", "")
	config := abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{Config: config},
			{
				Config:            config,
				ResourceName:      abpDomainKeyTestAddress,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					return domainId + "/" + abpAttrs(t, s, abpDomainKeyTestAddress)["id"], nil
				},
			},
		},
	})
}

// A key created outside Terraform can be adopted, and the first apply after the
// import plans nothing.
func TestAbpDomainEncryptionKey_ImportThenApplyKeepsKey(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-import-apply")
	domainId := f.seedDomain(siteId, "key-import-apply.example.com", "")
	seeded := f.seedEncryptionKey(domainId, abpDomainTestKeyMaterial)
	config := abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:             config,
				ResourceName:       abpDomainKeyTestAddress,
				ImportState:        true,
				ImportStateId:      domainId + "/" + seeded,
				ImportStatePersist: true,
			},
			{Config: config, PlanOnly: true},
			{
				Config: config,
				Check: func(s *terraform.State) error {
					if id := abpAttrs(t, s, abpDomainKeyTestAddress)["id"]; id != seeded {
						return fmt.Errorf("the imported key %s was replaced by %s", seeded, id)
					}
					if deleted := f.deletedIds(); len(deleted) != 0 {
						return fmt.Errorf("the first apply after import deleted %v", deleted)
					}
					return nil
				},
			},
		},
	})
}

func TestAbpDomainEncryptionKey_ImportErrors(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-import-errors")
	domainId := f.seedDomain(siteId, "key-import-errors.example.com", "")
	seeded := f.seedEncryptionKey(domainId, abpDomainTestKeyMaterial)
	config := abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial)

	malformed := regexp.MustCompile(`expected import ID to be '<domain_id>/<key_id>'`)
	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:        config,
				ResourceName:  abpDomainKeyTestAddress,
				ImportState:   true,
				ImportStateId: seeded,
				ExpectError:   malformed,
			},
			{
				Config:        config,
				ResourceName:  abpDomainKeyTestAddress,
				ImportState:   true,
				ImportStateId: domainId + "/",
				ExpectError:   malformed,
			},
			{
				Config:        config,
				ResourceName:  abpDomainKeyTestAddress,
				ImportState:   true,
				ImportStateId: "/" + seeded,
				ExpectError:   malformed,
			},
			{
				Config:        config,
				ResourceName:  abpDomainKeyTestAddress,
				ImportState:   true,
				ImportStateId: domainId + "/" + abpFakeAbsentId,
				ExpectError:   regexp.MustCompile(`ABP Domain encryption key ` + abpFakeAbsentId + ` on domain ` + domainId + ` not found`),
			},
			{
				Config:        config,
				ResourceName:  abpDomainKeyTestAddress,
				ImportState:   true,
				ImportStateId: abpFakeAbsentId + "/" + seeded,
				ExpectError:   regexp.MustCompile(`ABP Domain encryption key ` + seeded + ` on domain ` + abpFakeAbsentId + ` not found`),
			},
		},
	})
}

// --- validation -------------------------------------------------------------

// The provider checks the length of the base64 text, which the API's 24..64
// decoded bytes translate to.
func TestAbpDomainEncryptionKey_KeyLengthIsRejectedAtPlanTime(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-length")
	domainId := f.seedDomain(siteId, "key-length.example.com", "")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpDomainKeyOnConfig(domainId, abpDomainTestBase64(t, 16)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected length of key to be in the range \(32 - 88\)`),
			},
			{
				Config:      abpDomainKeyOnConfig(domainId, abpDomainTestBase64(t, 96)),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected length of key to be in the range \(32 - 88\)`),
			},
			{Config: abpDomainKeyOnConfig(domainId, abpDomainTestKeyMaterial)},
		},
	})
}

// Known limitation: `key` is length-validated but not base64-validated, so a
// value of the right length that is not base64 is only rejected by the API,
// during apply. The same goes for a value whose length passes but whose decoded
// size does not.
func TestAbpDomainEncryptionKey_InvalidKeyIsRejectedAtApplyTime(t *testing.T) {
	f := newAbpFake(t)
	siteId := f.seedSite("key-base64")
	domainId := f.seedDomain(siteId, "key-base64.example.com", "")

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpDomainKeyOnConfig(domainId, "not base64 at all, but long enough to pass"),
				ExpectError: regexp.MustCompile(`key must be valid base64`),
			},
			{
				// 88 base64 characters decode to 66 bytes, two more than the
				// API accepts.
				Config:      abpDomainKeyOnConfig(domainId, abpDomainTestBase64(t, 66)),
				ExpectError: regexp.MustCompile(`key must decode to between 24 and 64 bytes`),
			},
		},
	})

	if f.count(abpFakeKeys) != 0 {
		t.Fatalf("a rejected key configuration created %d keys", f.count(abpFakeKeys))
	}
}

func TestAbpDomainEncryptionKey_UnknownDomainIsRejected(t *testing.T) {
	f := newAbpFake(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: f.providers(t),
		Steps: []resource.TestStep{
			{
				Config:      abpDomainKeyOnConfig("not-a-uuid", abpDomainTestKeyMaterial),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected "domain_id" to be a valid UUID`),
			},
			{
				Config:      abpDomainKeyOnConfig(abpFakeAbsentId, abpDomainTestKeyMaterial),
				ExpectError: regexp.MustCompile(`domain ` + abpFakeAbsentId + ` not found`),
			},
		},
	})
}
