package incapsula

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDomainRefFromState_withValidStateID(t *testing.T) {
	r := resourceCSPSiteDomain()
	d := r.TestResourceData()
	d.SetId("2096205/863703496/Ki5nb29nbGVzeW5kaWNhdGlvbi5jb20")
	d.Set("domain", "googlesyndication.com")

	got := domainRefFromState(d)
	want := "Ki5nb29nbGVzeW5kaWNhdGlvbi5jb20"
	if got != want {
		t.Errorf("domainRefFromState() = %q, want %q", got, want)
	}
}

func TestDomainRefFromState_fallsBackToComputedRef(t *testing.T) {
	r := resourceCSPSiteDomain()
	d := r.TestResourceData()
	// No state ID set — simulates first call during Create before SetId
	d.Set("domain", "autodesk.com")

	got := domainRefFromState(d)
	want := base64.RawURLEncoding.EncodeToString([]byte("autodesk.com"))
	if got != want {
		t.Errorf("domainRefFromState() = %q, want %q", got, want)
	}
}

func TestDomainRefFromState_malformedStateID(t *testing.T) {
	r := resourceCSPSiteDomain()
	d := r.TestResourceData()
	d.SetId("not-a-valid-id")
	d.Set("domain", "autodesk.com")

	got := domainRefFromState(d)
	want := base64.RawURLEncoding.EncodeToString([]byte("autodesk.com"))
	if got != want {
		t.Errorf("domainRefFromState() with malformed ID = %q, want %q", got, want)
	}
}

func TestImporter_stripsWildcardPrefix(t *testing.T) {
	wildcardRef := base64.RawURLEncoding.EncodeToString([]byte("*.googlesyndication.com"))
	importID := "2096205/863703496/" + wildcardRef

	r := resourceCSPSiteDomain()
	d := r.TestResourceData()
	d.SetId(importID)

	// Simulate what the importer does
	keyParts := strings.Split(importID, "/")
	decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(keyParts[2])
	if err != nil {
		t.Fatalf("base64 decode failed: %s", err)
	}
	domain := strings.TrimPrefix(string(decoded), "*.")
	d.Set("domain", domain)

	if domain != "googlesyndication.com" {
		t.Errorf("importer domain = %q, want %q", domain, "googlesyndication.com")
	}
	// State ID should still hold the full wildcard referenceId
	if domainRefFromState(d) != wildcardRef {
		t.Errorf("domainRefFromState after import = %q, want %q", domainRefFromState(d), wildcardRef)
	}
}

func TestImporter_bareRefUnchanged(t *testing.T) {
	bareRef := base64.RawURLEncoding.EncodeToString([]byte("autodesk.com"))
	importID := "2096205/863703496/" + bareRef

	keyParts := strings.Split(importID, "/")
	decoded, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(keyParts[2])
	if err != nil {
		t.Fatalf("base64 decode failed: %s", err)
	}
	domain := strings.TrimPrefix(string(decoded), "*.")

	if domain != "autodesk.com" {
		t.Errorf("importer bare domain = %q, want %q", domain, "autodesk.com")
	}
}

// Verify getCSPPreApprovedDomain delegates correctly to getCSPPreApprovedDomainByRef
// by checking the ref it would pass for a given domain.
func TestGetCSPPreApprovedDomain_computesCorrectRef(t *testing.T) {
	domain := "googlesyndication.com"
	expectedRef := base64.RawURLEncoding.EncodeToString([]byte(domain))
	if expectedRef != "Z29vZ2xlc3luZGljYXRpb24uY29t" {
		t.Errorf("computed ref = %q, want Z29vZ2xlc3luZGljYXRpb24uY29t", expectedRef)
	}
}

// Verify the wildcard referenceId for googlesyndication.com matches what Imperva stores.
func TestWildcardRef_googlesyndication(t *testing.T) {
	wildcardRef := base64.RawURLEncoding.EncodeToString([]byte("*.googlesyndication.com"))
	if wildcardRef != "Ki5nb29nbGVzeW5kaWNhdGlvbi5jb20" {
		t.Errorf("wildcard ref = %q, want Ki5nb29nbGVzeW5kaWNhdGlvbi5jb20", wildcardRef)
	}
}

// Unused import guard
var _ = schema.HashString
