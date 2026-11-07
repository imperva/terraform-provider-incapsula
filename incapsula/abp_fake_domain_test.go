package incapsula

import (
	"encoding/base64"
	"testing"
)

// Fixtures and fake-backend helpers shared by the Domain, Domain encryption key
// and Credential tests.

const (
	// abpDomainTestKeyMaterial is base64 of 48 bytes, inside the 24..64 the API
	// accepts.
	abpDomainTestKeyMaterial      = "U2VjcmV0IGtleSB1c2luZyBzdGF0ZS1vZi10aGUtYXJ0IGJhc2U2NCBlbmNyeXB0aW9u"
	abpDomainTestOtherKeyMaterial = "QW5vdGhlciBzZWNyZXQga2V5LCBqdXN0IGFzIGxvbmcgYXMgdGhlIG90aGVyIG9uZS4="

	// Ids of encryption keys that live on some other Domain, for
	// `encryption_key_id`, which the API only ever consumes.
	abpDomainTestKeyId      = "ecc00000-0000-4000-8000-000000009001"
	abpDomainTestOtherKeyId = "ecc00000-0000-4000-8000-000000009002"
)

// dropDomainWithKeys deletes a Domain out of band the way the API does: the
// Domain's encryption keys are cascade-deleted with it, which is what makes a
// Domain replacement destructive.
func (f *abpFake) dropDomainWithKeys(domainId string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, id := range append([]string{}, f.order[abpFakeKeys]...) {
		if key := f.store[abpFakeKeys][id]; key != nil && abpFakeString(key["domain_id"]) == domainId {
			f.remove(abpFakeKeys, id)
		}
	}
	f.remove(abpFakeDomains, domainId)
}

// abpDomainTestBase64 renders n bytes as base64, for the key length checks.
func abpDomainTestBase64(t *testing.T, n int) string {
	t.Helper()
	raw := make([]byte, n)
	for i := range raw {
		raw[i] = byte('a' + i%26)
	}
	return base64.StdEncoding.EncodeToString(raw)
}
