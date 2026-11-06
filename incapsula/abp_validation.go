package incapsula

import (
	"fmt"
	"strings"
)

// Plan-time validation for ABP attributes whose value the API stores in a
// normalised form.
//
// The provider deliberately does not re-derive the API's normalisation (it
// would have to track every change to it, and any divergence brings the
// problem back). What it does instead is refuse the input shapes that could
// never converge: the API would store something other than what the
// configuration says, so every plan would report the same change and applying
// it would not help. Rejecting those while the configuration is planned keeps
// the provider out of the business of guessing the normalised form, and the
// message names the value to use instead.

// validateAbpName rejects a name the API would not store verbatim. Site,
// Policy and Proof Of Work Configuration names share one type on the API side,
// which trims the value, refuses it once empty and caps it at 100 characters.
// A padded name is the harmful case: it never matches the value read back, so
// the plan never converges.
func validateAbpName(i any, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return nil, []error{fmt.Errorf("%q must not be empty or consist only of whitespace", k)}
	}
	if trimmed != v {
		return nil, []error{fmt.Errorf("%q must not have leading or trailing whitespace, which the API trims away, leaving a permanent diff; use %q", k, trimmed)}
	}
	return nil, nil
}

// validateAbpDomainNamePart rejects a Domain criteria value the API would not
// store verbatim. The API trims domain name parts, so a padded value never
// matches the value read back; because criteria forces replacement, that diff
// is not a harmless permanent one — every apply would destroy the Domain, along
// with its encryption keys and its slot in the Site's Domain priority order,
// and create a new one. An empty value means "match on another criteria
// instead" and is left to the "exactly one of" check.
func validateAbpDomainNamePart(i any, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	trimmed := strings.TrimSpace(v)
	if trimmed == v {
		return nil, nil
	}
	if trimmed == "" {
		return nil, []error{fmt.Errorf("%q must not consist only of whitespace, which the API trims away before rejecting the empty value; omit it to match on another criteria instead", k)}
	}
	return nil, []error{fmt.Errorf("%q must not have leading or trailing whitespace, which the API trims away, so that every apply would replace the Domain; use %q", k, trimmed)}
}

// validateAbpCookiescope rejects a cookiescope the API would not store
// verbatim: it trims the value, so a padded one leaves a permanent diff. A
// whitespace-only value trims to the empty string, which the API accepts, so
// there the fix is to write an empty string.
func validateAbpCookiescope(i any, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	if trimmed := strings.TrimSpace(v); trimmed != v {
		return nil, []error{fmt.Errorf("%q must not have leading or trailing whitespace, which the API trims away, leaving a permanent diff; use %q", k, trimmed)}
	}
	return nil, nil
}
