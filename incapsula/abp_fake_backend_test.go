package incapsula

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// abpFake is an in-memory stand-in for the ABP API, good enough to drive every
// incapsula_abp_* resource and data source end to end without credentials and
// without a server to talk to. The whole ABP test suite runs against it, so
// resources that reference each other (a Policy Directive's Condition List, a
// Domain's Site, a Site's priority order) can be exercised together.
//
// Objects are kept as the JSON documents they arrived as, so fields the
// provider does not model ride along untouched, the way the real service
// behaves. On top of that the fake reproduces the semantics that shape the
// Terraform model:
//
//   - a literal Condition's `code` is reformatted before it is stored, so a
//     read rarely returns the value that was sent, and reformatting is
//     idempotent
//   - a Policy owns its Directives' Condition Lists: they are allocated server
//     side, a Directive submitted without one gets a new empty list, and a list
//     the update no longer references is deleted together with its entries
//   - `proof_of_work_configuration_id` and `skip_condition_id` are rejected on
//     any action other than `proof_of_work`
//   - a Site always ends with a catch-all `path_prefix = "/"` Selector, and a
//     Selector's criteria is a oneOf resolved by key presence
//   - a Domain's `encryption_key_id` is consumed on create and never returned
//   - a Credential's secret is server generated and cannot be changed
//   - priority endpoints only accept a list covering exactly the objects that
//     are attached, and objects never prioritised sort last
//   - publishing a preflight makes every earlier preflight unpublishable
//
// Policy, Site and Proof of Work Configuration names and descriptions, and
// Domain name parts, are trimmed before they are stored, which is what makes a
// padded value impossible to converge on. Condition names are stored verbatim.
type abpFake struct {
	mu     sync.Mutex
	server *httptest.Server

	seq   int
	clock int

	// store holds every object, by collection and id, and order remembers
	// the sequence they were created in so lists are stable.
	store map[string]map[string]map[string]any
	order map[string][]string

	requests []abpFakeRequest
	deleted  []string

	// encryptionKeyOf records the `encryption_key_id` each Domain was created
	// with, which is all the fake does with a value the API consumes once.
	encryptionKeyOf map[string]string

	// priority of a Site within its Account and of a Domain within its Site.
	// Zero means "never prioritised", which the API sorts last.
	sitePriority   map[string]int
	domainPriority map[string]int

	// globalPolicyId is the Policy the Account's global_policy endpoint
	// returns, when one has been seeded.
	globalPolicyId string

	// preflight publishing
	preflightSeq        map[string]int
	publishedSeq        int
	published           []string
	rejectPublish       bool
	preflightOutdated   bool
	preflightNoWarnings bool
}

// abpFakeRequest is one request the fake served, for tests that need to see
// what the provider actually sent.
type abpFakeRequest struct {
	Method string
	Path   string
	Body   string
}

// Ids the tests use. They are valid UUIDs because the provider validates
// references as such.
const (
	abpFakeAccountId      = "11111111-1111-1111-1111-111111111111"
	abpFakeOtherAccountId = "99999999-9999-9999-9999-999999999999"
	abpFakeAbsentId       = "00000000-0000-4000-8000-000000000000"
)

// The collections the fake stores, with the id prefix each one hands out.
var abpFakeIdPrefix = map[string]string{
	abpFakeSites:       "51e00000",
	abpFakeConditions:  "c0d00000",
	abpFakePolicies:    "b0110000",
	abpFakePowConfigs:  "b0b00000",
	abpFakeDomains:     "d0000000",
	abpFakeCredentials: "cbed0000",
	abpFakeKeys:        "ecc00000",
	abpFakePreflights:  "efef0000",
}

const (
	abpFakeSites       = "site"
	abpFakeConditions  = "condition"
	abpFakePolicies    = "policy"
	abpFakePowConfigs  = "proof_of_work_configuration"
	abpFakeDomains     = "domain"
	abpFakeCredentials = "credential"
	abpFakeKeys        = "encryptionkey"
	abpFakePreflights  = "preflight"
)

// abpFakeCriteriaKeys is the declaration order a Selector criteria oneOf is
// resolved in.
var abpFakeCriteriaKeys = []string{"postback", "path_prefix", "path_regex", "api_security_group"}

func newAbpFake(t *testing.T) *abpFake {
	f := &abpFake{
		store:           map[string]map[string]map[string]any{},
		order:           map[string][]string{},
		encryptionKeyOf: map[string]string{},
		sitePriority:    map[string]int{},
		domainPriority:  map[string]int{},
		preflightSeq:    map[string]int{},
	}
	for collection := range abpFakeIdPrefix {
		f.store[collection] = map[string]map[string]any{}
	}
	f.server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)
	return f
}

// abpRunTerraformTests opts into the ABP tests that drive a real terraform
// binary. They take minutes rather than seconds, so a plain `go test
// ./incapsula` skips them and only `-abp` (or ABP_TF_TESTS=1, for CI and for
// shells that keep it exported) runs them.
var abpRunTerraformTests = flag.Bool("abp", false,
	"run the ABP tests that drive the terraform CLI (same as ABP_TF_TESTS=1)")

// requireTerraform skips the calling test unless the ABP terraform tests were
// asked for.
func requireTerraform(t *testing.T) {
	t.Helper()
	if !*abpRunTerraformTests && os.Getenv("ABP_TF_TESTS") == "" {
		t.Skip("skipping ABP terraform test: pass -abp or set ABP_TF_TESTS=1 to run it")
	}
}

// providers points the provider at the fake. The credentials are never
// checked, but the provider insists on having some.
//
// Every test that drives terraform goes through here, which is what makes this
// the one place the -abp gate has to be applied.
func (f *abpFake) providers(t *testing.T) map[string]func() (*schema.Provider, error) {
	requireTerraform(t)
	t.Setenv("TF_ACC", "1")
	t.Setenv("INCAPSULA_API_ID", "fake-id")
	t.Setenv("INCAPSULA_API_KEY", "fake-key")
	t.Setenv("INCAPSULA_BASE_URL", f.server.URL)
	t.Setenv("INCAPSULA_BASE_URL_REV_2", f.server.URL)
	t.Setenv("INCAPSULA_BASE_URL_REV_3", f.server.URL)
	t.Setenv("INCAPSULA_BASE_URL_API", f.server.URL)
	return map[string]func() (*schema.Provider, error){
		"incapsula": func() (*schema.Provider, error) { return Provider(), nil },
	}
}

// client returns a client talking to the fake, for tests that call the API
// layer directly.
func (f *abpFake) client() *Client {
	return &Client{
		config:     &Config{APIID: "fake-id", APIKey: "fake-key", BaseURLAPI: f.server.URL},
		httpClient: &http.Client{},
	}
}

// --- routing ----------------------------------------------------------------

func (f *abpFake) serve(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)

	f.mu.Lock()
	defer f.mu.Unlock()

	f.requests = append(f.requests, abpFakeRequest{Method: r.Method, Path: r.URL.Path, Body: string(raw)})

	// The provider verifies its credentials while constructing its client.
	if strings.HasSuffix(r.URL.Path, "/account/verify") {
		f.write(w, http.StatusOK, map[string]any{"res": 0, "res_message": "OK", "account_id": 1234})
		return
	}

	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/botmanagement/v1"), "/")
	var segments []string
	if path != "" {
		segments = strings.Split(path, "/")
	}

	switch {
	case len(segments) == 0:
		// The API root carries the Account the credentials belong to.
		f.write(w, http.StatusOK, map[string]any{"account_id": abpFakeAccountId})
	case segments[0] == "account" && len(segments) == 3:
		f.serveAccount(w, r, segments[1], segments[2], raw)
	case len(segments) == 2:
		f.serveObject(w, r, segments[0], segments[1], raw)
	case len(segments) == 3:
		f.serveSubObject(w, r, segments[0], segments[1], segments[2], raw)
	default:
		f.fail(w, http.StatusNotFound, "no route for "+r.URL.Path)
	}
}

func (f *abpFake) serveAccount(w http.ResponseWriter, r *http.Request, accountId, kind string, raw []byte) {
	switch kind {
	case "site_priority":
		switch r.Method {
		case http.MethodGet:
			f.write(w, http.StatusOK, map[string]any{"site_ids": f.siteOrder(accountId)})
		case http.MethodPut:
			f.updatePriority(w, raw, "site_ids", f.siteOrder(accountId), f.sitePriority, func() []string { return f.siteOrder(accountId) })
		default:
			f.fail(w, http.StatusMethodNotAllowed, "unsupported")
		}
	case "global_policy":
		policy := f.store[abpFakePolicies][f.globalPolicyId]
		if r.Method != http.MethodGet || policy == nil {
			f.fail(w, http.StatusNotFound, "the account has no global policy")
			return
		}
		f.write(w, http.StatusOK, policy)
	case abpFakePreflights:
		if r.Method != http.MethodPost {
			f.fail(w, http.StatusMethodNotAllowed, "unsupported")
			return
		}
		f.write(w, http.StatusCreated, f.newPreflight(accountId))
	case abpFakeSites, abpFakeConditions, abpFakePolicies, abpFakePowConfigs, abpFakeDomains, abpFakeCredentials:
		switch r.Method {
		case http.MethodGet:
			f.write(w, http.StatusOK, map[string]any{"items": f.listOf(kind, accountId), "_links": map[string]any{}})
		case http.MethodPost:
			f.create(w, kind, accountId, raw)
		default:
			f.fail(w, http.StatusMethodNotAllowed, "unsupported")
		}
	default:
		f.fail(w, http.StatusNotFound, "unknown account collection "+kind)
	}
}

func (f *abpFake) serveObject(w http.ResponseWriter, r *http.Request, collection, id string, raw []byte) {
	if _, ok := abpFakeIdPrefix[collection]; !ok {
		f.fail(w, http.StatusNotFound, "unknown collection "+collection)
		return
	}
	stored := f.store[collection][id]
	if stored == nil {
		f.fail(w, http.StatusNotFound, fmt.Sprintf("%s %s not found", collection, id))
		return
	}

	switch r.Method {
	case http.MethodGet:
		if collection == abpFakePreflights {
			f.write(w, http.StatusOK, f.preflightEntity(id))
			return
		}
		f.write(w, http.StatusOK, stored)
	case http.MethodPut:
		f.update(w, collection, stored, raw)
	case http.MethodDelete:
		f.remove(collection, id)
		f.write(w, http.StatusOK, map[string]any{})
	default:
		f.fail(w, http.StatusMethodNotAllowed, "unsupported")
	}
}

func (f *abpFake) serveSubObject(w http.ResponseWriter, r *http.Request, collection, id, sub string, raw []byte) {
	switch {
	case collection == abpFakeSites && sub == "domain_priority":
		if f.store[abpFakeSites][id] == nil {
			f.fail(w, http.StatusNotFound, "site "+id+" not found")
			return
		}
		switch r.Method {
		case http.MethodGet:
			f.write(w, http.StatusOK, map[string]any{"domain_ids": f.domainOrder(id)})
		case http.MethodPut:
			f.updatePriority(w, raw, "domain_ids", f.domainOrder(id), f.domainPriority, func() []string { return f.domainOrder(id) })
		default:
			f.fail(w, http.StatusMethodNotAllowed, "unsupported")
		}
	case collection == abpFakeDomains && sub == abpFakeKeys:
		if f.store[abpFakeDomains][id] == nil {
			f.fail(w, http.StatusNotFound, "domain "+id+" not found")
			return
		}
		switch r.Method {
		case http.MethodGet:
			f.write(w, http.StatusOK, map[string]any{"items": f.keysOfDomain(id)})
		case http.MethodPost:
			f.createEncryptionKey(w, id, raw)
		default:
			f.fail(w, http.StatusMethodNotAllowed, "unsupported")
		}
	case collection == abpFakePreflights && sub == "publish":
		if r.Method != http.MethodPost {
			f.fail(w, http.StatusMethodNotAllowed, "unsupported")
			return
		}
		f.publish(w, id)
	default:
		f.fail(w, http.StatusNotFound, fmt.Sprintf("no route for %s/%s/%s", collection, id, sub))
	}
}

// --- create and update ------------------------------------------------------

func (f *abpFake) create(w http.ResponseWriter, collection, accountId string, raw []byte) {
	body, ok := f.decode(w, raw)
	if !ok {
		return
	}
	switch collection {
	case abpFakeConditions:
		f.createCondition(w, accountId, body)
	case abpFakePolicies:
		f.createPolicy(w, accountId, body)
	case abpFakeSites:
		f.createSite(w, accountId, body)
	case abpFakePowConfigs:
		f.createPowConfig(w, accountId, body)
	case abpFakeDomains:
		f.createDomain(w, accountId, body)
	case abpFakeCredentials:
		f.createCredential(w, accountId)
	default:
		f.fail(w, http.StatusNotFound, "cannot create "+collection)
	}
}

func (f *abpFake) update(w http.ResponseWriter, collection string, stored map[string]any, raw []byte) {
	body, ok := f.decode(w, raw)
	if !ok {
		return
	}
	switch collection {
	case abpFakeConditions:
		f.updateCondition(w, stored, body)
	case abpFakePolicies:
		f.updatePolicy(w, stored, body)
	case abpFakeSites:
		f.updateSite(w, stored, body)
	case abpFakePowConfigs:
		f.updatePowConfig(w, stored, body)
	case abpFakeDomains:
		f.updateDomain(w, stored, body)
	default:
		f.fail(w, http.StatusMethodNotAllowed, collection+" has no update endpoint")
	}
}

// --- Conditions, Condition Lists and Condition List Entries -----------------

// abpFakeConditionKind reports which of the three ConditionV1 variants a
// document is, the way the API's oneOf does: by which discriminator key is
// present, not by its value.
func abpFakeConditionKind(doc map[string]any) AbpConditionKind {
	if _, ok := doc["reference"]; ok {
		return AbpConditionKindReference
	}
	if _, ok := doc["code"]; ok {
		return AbpConditionKindLiteral
	}
	return AbpConditionKindList
}

// abpFakeFormatCode stands in for the API's expression formatter: it wraps the
// expression in parentheses and terminates it with a newline, and leaves an
// already formatted expression alone.
func abpFakeFormatCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return code
	}
	if !strings.HasPrefix(code, "(") || !strings.HasSuffix(code, ")") {
		code = "(" + code + ")"
	}
	return code + "\n"
}

func (f *abpFake) createCondition(w http.ResponseWriter, accountId string, body map[string]any) {
	doc := abpFakeClone(body)

	switch abpFakeConditionKind(body) {
	case AbpConditionKindReference:
		reference := abpFakeString(body["reference"])
		parent := abpFakeString(body["parent"])
		if f.store[abpFakeConditions][reference] == nil {
			f.fail(w, http.StatusBadRequest, "Unknown condition "+reference)
			return
		}
		if f.store[abpFakeConditions][parent] == nil {
			f.fail(w, http.StatusBadRequest, "Unknown parent condition list "+parent)
			return
		}
		if doc["tags"] == nil {
			doc["tags"] = []any{}
		}
		if abpFakeString(doc["state"]) == "" {
			doc["state"] = abpConditionStateInactive
		}
	case AbpConditionKindLiteral:
		if !f.checkName(w, doc, false) {
			return
		}
		doc["code"] = abpFakeFormatCode(abpFakeString(body["code"]))
		if abpFakeString(doc["template"]) == "" {
			doc["template"] = "custom"
		}
	case AbpConditionKindList:
		if !f.checkName(w, doc, false) {
			return
		}
	}

	doc["id"] = f.nextId(abpFakeConditions)
	doc["account_id"] = accountId
	doc["last_change_by"] = "fake"
	doc["created_at"] = f.now()
	doc["modified_at"] = doc["created_at"]

	f.insert(abpFakeConditions, doc)
	f.write(w, http.StatusCreated, doc)
}

// updateCondition replaces the stored document with the body, as a
// full-replace endpoint does, keeping the fields an update cannot change.
func (f *abpFake) updateCondition(w http.ResponseWriter, stored, body map[string]any) {
	if abpFakeConditionKind(body) != abpFakeConditionKind(stored) {
		f.fail(w, http.StatusBadRequest, fmt.Sprintf("cannot turn a %s condition into a %s one",
			abpFakeConditionKind(stored), abpFakeConditionKind(body)))
		return
	}

	updated := abpFakeClone(body)
	switch abpFakeConditionKind(body) {
	case AbpConditionKindReference:
		reference := abpFakeString(body["reference"])
		if f.store[abpFakeConditions][reference] == nil {
			f.fail(w, http.StatusBadRequest, "Unknown condition "+reference)
			return
		}
		if updated["tags"] == nil {
			updated["tags"] = []any{}
		}
	case AbpConditionKindLiteral:
		if !f.checkName(w, updated, false) {
			return
		}
		updated["code"] = abpFakeFormatCode(abpFakeString(body["code"]))
		// The template is fixed when the Condition is created.
		updated["template"] = stored["template"]
	case AbpConditionKindList:
		if !f.checkName(w, updated, false) {
			return
		}
	}

	updated["id"] = stored["id"]
	updated["account_id"] = stored["account_id"]
	updated["created_at"] = stored["created_at"]
	updated["last_change_by"] = "fake"
	updated["modified_at"] = f.now()

	f.replace(abpFakeConditions, updated)
	f.write(w, http.StatusOK, updated)
}

// --- Policies ---------------------------------------------------------------

// abpFakeValidateDirective mirrors the API's directive validation.
func abpFakeValidateDirective(directive map[string]any) error {
	action := strings.TrimSpace(abpFakeString(directive["action"]))
	if action == "" {
		return fmt.Errorf("Action must not be empty")
	}
	if len(action) > 100 || !abpDirectiveActionRegexp.MatchString(action) {
		return fmt.Errorf("Action must be snake_case")
	}
	pow := abpFakeString(directive["proof_of_work_configuration_id"]) != ""
	skip := abpFakeString(directive["skip_condition_id"]) != ""
	if (pow || skip) && action != "proof_of_work" {
		return fmt.Errorf("proof_of_work_configuration_id is only allowed for proof_of_work action")
	}
	return nil
}

func (f *abpFake) createPolicy(w http.ResponseWriter, accountId string, body map[string]any) {
	if !f.trimName(w, body) {
		return
	}
	directives := abpFakeList(body["directives"])
	if len(directives) == 0 {
		f.fail(w, http.StatusBadRequest, "Policy has no directives")
		return
	}
	for _, directive := range directives {
		if err := abpFakeValidateDirective(directive); err != nil {
			f.fail(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	doc := abpFakeClone(body)
	doc["id"] = f.nextId(abpFakePolicies)
	doc["account_id"] = accountId
	doc["created_at"] = f.now()
	doc["modified_at"] = doc["created_at"]
	doc["_links"] = map[string]any{}

	// CreateDirectiveV1 has no condition_id at all: the Condition Lists are
	// always allocated server side.
	out := make([]any, 0, len(directives))
	for _, directive := range directives {
		action := strings.TrimSpace(abpFakeString(directive["action"]))
		fresh := map[string]any{"action": action, "condition_id": f.newConditionList(accountId, action)}
		if pow := abpFakeString(directive["proof_of_work_configuration_id"]); pow != "" {
			fresh["proof_of_work_configuration_id"] = pow
		}
		if action == "proof_of_work" {
			fresh["skip_condition_id"] = f.newConditionList(accountId, action+"_skip")
		}
		out = append(out, fresh)
	}
	doc["directives"] = out

	f.insert(abpFakePolicies, doc)
	f.write(w, http.StatusCreated, doc)
}

func (f *abpFake) updatePolicy(w http.ResponseWriter, stored, body map[string]any) {
	if !f.trimName(w, body) {
		return
	}
	directives := abpFakeList(body["directives"])
	if len(directives) == 0 {
		f.fail(w, http.StatusBadRequest, "Policy directives must not be empty")
		return
	}
	for _, directive := range directives {
		if err := abpFakeValidateDirective(directive); err != nil {
			f.fail(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Submitted Condition Lists must already belong to this Policy, and must
	// not repeat.
	submitted := map[string]bool{}
	for _, directive := range directives {
		for _, key := range []string{"condition_id", "skip_condition_id"} {
			id := abpFakeString(directive[key])
			if id == "" {
				continue
			}
			if submitted[id] {
				f.fail(w, http.StatusBadRequest, "Policy has duplicate directives")
				return
			}
			submitted[id] = true
		}
	}
	owned := abpFakePolicyConditionLists(stored)
	for id := range submitted {
		if !owned[id] {
			f.fail(w, http.StatusBadRequest, "Provided directives are not owned by this policy")
			return
		}
	}

	accountId := abpFakeString(stored["account_id"])
	updated := abpFakeClone(body)
	updated["id"] = stored["id"]
	updated["account_id"] = stored["account_id"]
	updated["created_at"] = stored["created_at"]
	updated["modified_at"] = f.now()
	updated["_links"] = map[string]any{}

	out := make([]any, 0, len(directives))
	for _, directive := range directives {
		action := strings.TrimSpace(abpFakeString(directive["action"]))
		fresh := map[string]any{"action": action}
		if id := abpFakeString(directive["condition_id"]); id != "" {
			fresh["condition_id"] = id
		} else {
			fresh["condition_id"] = f.newConditionList(accountId, action)
		}
		if pow := abpFakeString(directive["proof_of_work_configuration_id"]); pow != "" {
			fresh["proof_of_work_configuration_id"] = pow
		}
		switch {
		case abpFakeString(directive["skip_condition_id"]) != "":
			fresh["skip_condition_id"] = abpFakeString(directive["skip_condition_id"])
		case action == "proof_of_work":
			fresh["skip_condition_id"] = f.newConditionList(accountId, action+"_skip")
		}
		out = append(out, fresh)
	}
	updated["directives"] = out

	// Condition Lists the update no longer references are deleted, together
	// with every entry inside them.
	for id := range owned {
		if !submitted[id] {
			f.removeConditionList(id)
		}
	}

	f.replace(abpFakePolicies, updated)
	f.write(w, http.StatusOK, updated)
}

// abpFakePolicyConditionLists returns the ids of every Condition List the
// Policy owns.
func abpFakePolicyConditionLists(policy map[string]any) map[string]bool {
	out := map[string]bool{}
	for _, directive := range abpFakeList(policy["directives"]) {
		for _, key := range []string{"condition_id", "skip_condition_id"} {
			if id := abpFakeString(directive[key]); id != "" {
				out[id] = true
			}
		}
	}
	return out
}

// newConditionList allocates a Condition List for a Directive. The real API
// stores it as an ordinary Condition of the list variant, which is what makes
// it usable as an entry's parent.
func (f *abpFake) newConditionList(accountId, name string) string {
	id := f.nextId(abpFakeConditions)
	f.insert(abpFakeConditions, map[string]any{
		"id":          id,
		"account_id":  accountId,
		"name":        name,
		"description": "",
		"created_at":  f.now(),
		"modified_at": f.now(),
	})
	return id
}

// removeConditionList deletes a Condition List and every entry attached to it,
// as the API does when the Policy no longer references it.
func (f *abpFake) removeConditionList(id string) {
	f.purge(abpFakeConditions, id)
}

// --- Sites ------------------------------------------------------------------

// abpFakeResolveCriteria keeps the first criteria key present in declaration
// order, as the real oneOf resolution does.
func abpFakeResolveCriteria(selector map[string]any) map[string]any {
	criteria, _ := selector["criteria"].(map[string]any)
	resolved := map[string]any{}
	if criteria == nil {
		return resolved
	}
	for _, key := range abpFakeCriteriaKeys {
		if value, ok := criteria[key]; ok {
			resolved[key] = value
			// case_insensitive only applies to path_prefix, and the API omits
			// it when false.
			if key == "path_prefix" {
				if insensitive, ok := criteria["case_insensitive"].(bool); ok && insensitive {
					resolved["case_insensitive"] = true
				}
			}
			break
		}
	}
	return resolved
}

func (f *abpFake) normalizeSelectors(site map[string]any) {
	for _, selector := range abpFakeList(site["selectors"]) {
		if abpFakeString(selector["id"]) == "" {
			selector["id"] = f.nextId("selector")
			selector["created_at"] = f.now()
		}
		criteria := abpFakeResolveCriteria(selector)
		selector["criteria"] = criteria
		selector["modified_at"] = f.now()
		encoded, _ := json.Marshal(criteria)
		selector["derived_id"] = fmt.Sprintf("%x", sha1.Sum(encoded))
		if _, ok := selector["policy_id"]; !ok {
			selector["policy_id"] = nil
		}
	}
}

// addDefaultSelector appends the catch-all Selector the API adds to every Site
// it creates. It is appended unconditionally, as the real one is: a Site whose
// Selectors already end with a `/` prefix gets a second one anyway. The
// catch-all points at the Policy the API creates along with the Site.
//
// Only a create appends it. An update stores the Selectors it is given, so a
// Site can be left without a catch-all, which is how the real API behaves.
func (f *abpFake) addDefaultSelector(site map[string]any) {
	raw, _ := site["selectors"].([]any)
	site["selectors"] = append(raw, map[string]any{
		"criteria":          map[string]any{"path_prefix": "/"},
		"analysis_settings": map[string]any{"rate_limiting": "per_site"},
		"policy_id":         f.newDefaultPolicy(abpFakeString(site["account_id"]), abpFakeString(site["name"])),
	})
	f.normalizeSelectors(site)
}

// newDefaultPolicy creates the `"<site> Default"` Policy the backend creates
// for every Site, with Directives in an order the caller does not control.
func (f *abpFake) newDefaultPolicy(accountId, siteName string) string {
	actions := []string{"block", "allow", "captcha"}
	directives := []any{}
	for _, action := range actions {
		directives = append(directives, map[string]any{
			"action":       action,
			"condition_id": f.newConditionList(accountId, action),
		})
	}
	id := f.nextId(abpFakePolicies)
	f.insert(abpFakePolicies, map[string]any{
		"id":          id,
		"account_id":  accountId,
		"name":        siteName + " Default",
		"description": "",
		"directives":  directives,
		"created_at":  f.now(),
		"modified_at": f.now(),
		"_links":      map[string]any{},
	})
	return id
}

func (f *abpFake) createSite(w http.ResponseWriter, accountId string, body map[string]any) {
	if !f.trimName(w, body) {
		return
	}
	doc := abpFakeClone(body)
	doc["id"] = f.nextId(abpFakeSites)
	doc["account_id"] = accountId
	doc["created_at"] = f.now()
	doc["modified_at"] = doc["created_at"]
	f.normalizeSelectors(doc)
	f.addDefaultSelector(doc)

	f.insert(abpFakeSites, doc)
	f.write(w, http.StatusCreated, doc)
}

func (f *abpFake) updateSite(w http.ResponseWriter, stored, body map[string]any) {
	if !f.trimName(w, body) {
		return
	}
	updated := abpFakeClone(body)
	updated["id"] = stored["id"]
	updated["account_id"] = stored["account_id"]
	updated["created_at"] = stored["created_at"]
	updated["modified_at"] = f.now()
	f.normalizeSelectors(updated)

	f.replace(abpFakeSites, updated)
	f.write(w, http.StatusOK, updated)
}

// --- Proof of Work Configurations -------------------------------------------

func (f *abpFake) createPowConfig(w http.ResponseWriter, accountId string, body map[string]any) {
	if !f.trimName(w, body) {
		return
	}
	doc := abpFakeClone(body)
	// The API has no default algorithm: an omitted one is simply not stored,
	// and reads back as absent.
	if abpFakeString(doc["algorithm"]) == "" {
		delete(doc, "algorithm")
	}
	doc["id"] = f.nextId(abpFakePowConfigs)
	doc["account_id"] = accountId
	doc["created_at"] = f.now()
	doc["modified_at"] = doc["created_at"]

	f.insert(abpFakePowConfigs, doc)
	f.write(w, http.StatusCreated, doc)
}

func (f *abpFake) updatePowConfig(w http.ResponseWriter, stored, body map[string]any) {
	if !f.trimName(w, body) {
		return
	}
	updated := abpFakeClone(body)
	if abpFakeString(updated["algorithm"]) == "" {
		delete(updated, "algorithm")
	}
	updated["id"] = stored["id"]
	updated["account_id"] = stored["account_id"]
	updated["created_at"] = stored["created_at"]
	updated["modified_at"] = f.now()

	f.replace(abpFakePowConfigs, updated)
	f.write(w, http.StatusOK, updated)
}

// --- Domains ----------------------------------------------------------------

// abpFakeTrimDomain mirrors the API, which trims the cookiescope and the
// criteria domain name parts before storing them and rejects a part that is
// empty once trimmed.
func abpFakeTrimDomain(domain map[string]any) error {
	if _, ok := domain["cookiescope"]; ok {
		domain["cookiescope"] = strings.TrimSpace(abpFakeString(domain["cookiescope"]))
	}
	criteria, _ := domain["criteria"].(map[string]any)
	for _, key := range []string{"exact", "prefix", "suffix"} {
		value, ok := criteria[key]
		if !ok {
			continue
		}
		trimmed := strings.TrimSpace(abpFakeString(value))
		if trimmed == "" {
			return fmt.Errorf("Domain name part must not be empty")
		}
		criteria[key] = trimmed
	}
	return nil
}

func (f *abpFake) createDomain(w http.ResponseWriter, accountId string, body map[string]any) {
	doc := abpFakeClone(body)
	if err := abpFakeTrimDomain(doc); err != nil {
		f.fail(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(abpFakeMap(doc["criteria"])) == 0 {
		f.fail(w, http.StatusBadRequest, "Domain criteria must select exactly one variant")
		return
	}
	siteId := abpFakeString(doc["site_id"])
	if f.store[abpFakeSites][siteId] == nil {
		f.fail(w, http.StatusBadRequest, "Unknown site "+siteId)
		return
	}

	doc["id"] = f.nextId(abpFakeDomains)
	doc["account_id"] = accountId
	doc["created_at"] = f.now()
	doc["modified_at"] = doc["created_at"]
	// encryption_key_id is consumed on create and is never part of the Domain.
	if key := abpFakeString(doc["encryption_key_id"]); key != "" {
		f.encryptionKeyOf[abpFakeString(doc["id"])] = key
	}
	delete(doc, "encryption_key_id")

	f.insert(abpFakeDomains, doc)
	f.write(w, http.StatusCreated, doc)
}

func (f *abpFake) updateDomain(w http.ResponseWriter, stored, body map[string]any) {
	if _, ok := body["encryption_key_id"]; ok {
		f.fail(w, http.StatusBadRequest, "encryption_key_id is not part of UpdateDomainV1")
		return
	}
	updated := abpFakeClone(body)
	if err := abpFakeTrimDomain(updated); err != nil {
		f.fail(w, http.StatusBadRequest, err.Error())
		return
	}
	if siteId := abpFakeString(updated["site_id"]); f.store[abpFakeSites][siteId] == nil {
		f.fail(w, http.StatusBadRequest, "Unknown site "+siteId)
		return
	}
	// criteria is not settable by an update, so it keeps its stored value.
	updated["criteria"] = stored["criteria"]
	updated["id"] = stored["id"]
	updated["account_id"] = stored["account_id"]
	updated["created_at"] = stored["created_at"]
	updated["modified_at"] = f.now()

	f.replace(abpFakeDomains, updated)
	f.write(w, http.StatusOK, updated)
}

// --- Credentials ------------------------------------------------------------

func (f *abpFake) createCredential(w http.ResponseWriter, accountId string) {
	f.write(w, http.StatusCreated, f.newCredential(accountId))
}

func (f *abpFake) newCredential(accountId string) map[string]any {
	id := f.nextId(abpFakeCredentials)
	doc := map[string]any{
		"id":         id,
		"account_id": accountId,
		// The secret is server generated and never changes: the only way to
		// get a new one is to create a new Credential.
		"secret":     "secret-" + id,
		"created_at": f.now(),
	}
	f.insert(abpFakeCredentials, doc)
	return doc
}

// --- Domain encryption keys -------------------------------------------------

func (f *abpFake) createEncryptionKey(w http.ResponseWriter, domainId string, raw []byte) {
	body, ok := f.decode(w, raw)
	if !ok {
		return
	}
	material := abpFakeString(body["key"])
	decoded, err := base64.StdEncoding.DecodeString(material)
	if err != nil {
		f.fail(w, http.StatusBadRequest, "key must be valid base64")
		return
	}
	if len(decoded) < 24 || len(decoded) > 64 {
		f.fail(w, http.StatusBadRequest, "key must decode to between 24 and 64 bytes")
		return
	}

	doc := map[string]any{
		"id":                 f.nextId(abpFakeKeys),
		"domain_id":          domainId,
		"key":                material,
		"created_at":         f.now(),
		"first_published_at": nil,
	}
	f.insert(abpFakeKeys, doc)
	f.write(w, http.StatusCreated, doc)
}

func (f *abpFake) keysOfDomain(domainId string) []map[string]any {
	items := []map[string]any{}
	for _, id := range f.order[abpFakeKeys] {
		key := f.store[abpFakeKeys][id]
		if key != nil && abpFakeString(key["domain_id"]) == domainId {
			items = append(items, key)
		}
	}
	return items
}

// --- Priority ---------------------------------------------------------------

// abpFakeOrder returns the attached ids the way a priority read does:
// prioritised first by priority, then the ones never prioritised in
// attachment order.
func abpFakeOrder(attached []string, priority map[string]int) []string {
	out := append([]string{}, attached...)
	sort.SliceStable(out, func(i, j int) bool {
		left, right := priority[out[i]], priority[out[j]]
		if (left == 0) != (right == 0) {
			return right == 0
		}
		return left < right
	})
	return out
}

func (f *abpFake) siteOrder(accountId string) []string {
	var attached []string
	for _, id := range f.order[abpFakeSites] {
		if site := f.store[abpFakeSites][id]; site != nil && abpFakeString(site["account_id"]) == accountId {
			attached = append(attached, id)
		}
	}
	return abpFakeOrder(attached, f.sitePriority)
}

func (f *abpFake) domainOrder(siteId string) []string {
	var attached []string
	for _, id := range f.order[abpFakeDomains] {
		if domain := f.store[abpFakeDomains][id]; domain != nil && abpFakeString(domain["site_id"]) == siteId {
			attached = append(attached, id)
		}
	}
	return abpFakeOrder(attached, f.domainPriority)
}

// updatePriority applies a priority write. The API insists on a list that
// covers exactly the attached objects: a list of the wrong length or one
// naming something that is not attached is rejected. Duplicates are not
// rejected, which is why the provider checks for them itself.
func (f *abpFake) updatePriority(w http.ResponseWriter, raw []byte, field string, attached []string, priority map[string]int, reread func() []string) {
	body, ok := f.decode(w, raw)
	if !ok {
		return
	}
	var ids []string
	for _, item := range abpFakeStrings(body[field]) {
		ids = append(ids, item)
	}

	if len(ids) != len(attached) {
		f.fail(w, http.StatusBadRequest, fmt.Sprintf("Wrong number of %s supplied", field))
		return
	}
	isAttached := map[string]bool{}
	for _, id := range attached {
		isAttached[id] = true
	}
	for _, id := range ids {
		if !isAttached[id] {
			f.fail(w, http.StatusBadRequest, "Unknown ID "+id)
			return
		}
	}
	for i, id := range ids {
		priority[id] = i + 1
	}
	f.write(w, http.StatusOK, map[string]any{field: reread()})
}

// --- Preflight and publish --------------------------------------------------

func (f *abpFake) newPreflight(accountId string) map[string]any {
	id := f.nextId(abpFakePreflights)
	f.insert(abpFakePreflights, map[string]any{"id": id, "account_id": accountId})
	f.preflightSeq[id] = len(f.preflightSeq) + 1
	return f.preflightEntity(id)
}

// preflightEntity renders a preflight the way the API does: the
// `create_publish` link is offered only while the preflight can still be
// published, and a preflight older than the last published one says so.
func (f *abpFake) preflightEntity(id string) map[string]any {
	canPublish := f.preflightSeq[id] >= f.publishedSeq

	var warnings any
	if !f.preflightNoWarnings {
		list := []any{}
		switch {
		case !canPublish:
			list = append(list, map[string]any{
				"type":    "older_than_most_recent_published",
				"message": "This preflight is older than the most recently published preflight. You must create a new preflight to publish.",
			})
		case f.preflightOutdated:
			list = append(list, map[string]any{
				"type":    "outdated",
				"message": "There have been changes since this preflight was created. Those changes will not be published.",
			})
		}
		warnings = list
	}

	links := map[string]any{}
	if canPublish {
		links["create_publish"] = map[string]any{"href": fmt.Sprintf("/v1/preflight/%s/publish", id)}
	}

	return map[string]any{
		"id":          id,
		"account_id":  abpFakeString(f.store[abpFakePreflights][id]["account_id"]),
		"created_at":  "2026-01-01T00:00:00Z",
		"modified_at": "2026-01-01T00:00:00Z",
		"warnings":    warnings,
		"_links":      links,
	}
}

func (f *abpFake) publish(w http.ResponseWriter, id string) {
	if f.store[abpFakePreflights][id] == nil {
		f.fail(w, http.StatusNotFound, "id "+id+" not found")
		return
	}
	if f.rejectPublish {
		f.rejectPublish = false
		f.fail(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	if f.preflightSeq[id] < f.publishedSeq {
		f.fail(w, http.StatusBadRequest, "Preflight is older than the latest published preflight")
		return
	}

	f.publishedSeq = f.preflightSeq[id]
	f.published = append(f.published, id)
	f.write(w, http.StatusCreated, map[string]any{"status": "in_progress"})
}

// --- store plumbing ---------------------------------------------------------

func (f *abpFake) nextId(collection string) string {
	f.seq++
	prefix, ok := abpFakeIdPrefix[collection]
	if !ok {
		// Selectors and other nested objects are not collections of their own.
		prefix = "5e1ec700"
	}
	return fmt.Sprintf("%s-0000-4000-8000-%012d", prefix, f.seq)
}

func (f *abpFake) now() string {
	f.clock++
	return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).
		Add(time.Duration(f.clock) * time.Second).Format(time.RFC3339)
}

func (f *abpFake) insert(collection string, doc map[string]any) {
	id := abpFakeString(doc["id"])
	f.store[collection][id] = doc
	f.order[collection] = append(f.order[collection], id)
}

func (f *abpFake) replace(collection string, doc map[string]any) {
	f.store[collection][abpFakeString(doc["id"])] = doc
}

func (f *abpFake) remove(collection, id string) {
	if f.store[collection][id] == nil {
		return
	}
	f.deleted = append(f.deleted, id)
	f.purge(collection, id)
}

// purge deletes an object and everything the service deletes along with it,
// without recording it as something the API was asked to delete.
func (f *abpFake) purge(collection, id string) {
	doc := f.store[collection][id]
	if doc == nil {
		return
	}
	delete(f.store[collection], id)
	kept := make([]string, 0, len(f.order[collection]))
	for _, existing := range f.order[collection] {
		if existing != id {
			kept = append(kept, existing)
		}
	}
	f.order[collection] = kept

	switch {
	case collection == abpFakePolicies:
		// Deleting a Policy deletes its Directives' Condition Lists.
		for list := range abpFakePolicyConditionLists(doc) {
			f.purge(abpFakeConditions, list)
		}
	case collection == abpFakeDomains:
		// Deleting a Domain deletes its encryption keys.
		for _, key := range f.keysOfDomain(id) {
			f.purge(abpFakeKeys, abpFakeString(key["id"]))
		}
	case collection == abpFakeConditions && abpFakeConditionKind(doc) == AbpConditionKindList:
		// Deleting a Condition List deletes every entry inside it. The
		// referenced Conditions are owned separately and are left alone.
		for _, entryId := range append([]string{}, f.order[abpFakeConditions]...) {
			entry := f.store[abpFakeConditions][entryId]
			if entry != nil && abpFakeString(entry["parent"]) == id {
				f.purge(abpFakeConditions, entryId)
			}
		}
	}
}

func (f *abpFake) listOf(collection, accountId string) []map[string]any {
	ids := f.order[collection]
	// Sites are listed most to least significant priority, not in creation
	// order.
	if collection == abpFakeSites {
		ids = f.siteOrder(accountId)
	}

	items := []map[string]any{}
	for _, id := range ids {
		doc := f.store[collection][id]
		if doc == nil {
			continue
		}
		// The Account's global Policy is only reachable through the Account's
		// own global_policy endpoint, never through the Policy list.
		if collection == abpFakePolicies && id == f.globalPolicyId {
			continue
		}
		owner := abpFakeString(doc["account_id"])
		// Managed objects carry no account and are visible from every Account.
		if owner == accountId || owner == "" {
			items = append(items, doc)
		}
	}
	return items
}

// checkName rejects a name that is empty once trimmed, and trims `name` and
// `description` when the endpoint does. Condition names are stored verbatim,
// unlike a Policy's, a Site's or a Proof of Work Configuration's.
func (f *abpFake) checkName(w http.ResponseWriter, doc map[string]any, trim bool) bool {
	if _, ok := doc["name"]; ok {
		name := strings.TrimSpace(abpFakeString(doc["name"]))
		if name == "" {
			f.fail(w, http.StatusBadRequest, "Name must not be empty")
			return false
		}
		if trim {
			doc["name"] = name
		}
	}
	if _, ok := doc["description"]; ok && trim {
		doc["description"] = strings.TrimSpace(abpFakeString(doc["description"]))
	}
	return true
}

// trimName is checkName for the endpoints that trim.
func (f *abpFake) trimName(w http.ResponseWriter, doc map[string]any) bool {
	return f.checkName(w, doc, true)
}

func (f *abpFake) decode(w http.ResponseWriter, raw []byte) (map[string]any, bool) {
	body := map[string]any{}
	if len(raw) == 0 {
		return body, true
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		f.fail(w, http.StatusBadRequest, err.Error())
		return nil, false
	}
	return body, true
}

func (f *abpFake) write(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func (f *abpFake) fail(w http.ResponseWriter, code int, message string) {
	f.write(w, code, map[string]any{"message": message})
}

// --- value helpers ----------------------------------------------------------

func abpFakeString(value any) string {
	s, _ := value.(string)
	return s
}

func abpFakeMap(value any) map[string]any {
	m, _ := value.(map[string]any)
	return m
}

// abpFakeList returns the elements of a JSON array of objects.
func abpFakeList(value any) []map[string]any {
	raw, _ := value.([]any)
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func abpFakeStrings(value any) []string {
	raw, _ := value.([]any)
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		out = append(out, abpFakeString(item))
	}
	return out
}

// abpFakeClone deep-copies a document so that what a test reads can never be
// changed under it by a later request.
func abpFakeClone(doc map[string]any) map[string]any {
	encoded, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}
	var out map[string]any
	if err := json.Unmarshal(encoded, &out); err != nil {
		panic(err)
	}
	return out
}

// --- helpers the whole ABP test suite shares --------------------------------

// abpAttrs returns the state attributes of one resource or data source.
func abpAttrs(t *testing.T, s *terraform.State, address string) map[string]string {
	t.Helper()
	found := s.RootModule().Resources[address]
	if found == nil {
		t.Fatalf("%s is not in state", address)
	}
	return found.Primary.Attributes
}

// abpTestRsaKeys returns two reusable RSA key pairs, so a test can check what
// the provider encrypted and what changing the key did.
var abpTestRsaKeys = sync.OnceValue(func() [2]*rsa.PrivateKey {
	var keys [2]*rsa.PrivateKey
	for i := range keys {
		key, err := rsa.GenerateKey(cryptorand.Reader, 3072)
		if err != nil {
			panic(err)
		}
		keys[i] = key
	}
	return keys
})

// abpTestRsaPublicPem returns one of the test keys as a PEM-encoded SPKI
// public key, the form the provider accepts.
func abpTestRsaPublicPem(t *testing.T, index int) string {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(&abpTestRsaKeys()[index].PublicKey)
	if err != nil {
		t.Fatalf("marshalling test public key: %s", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}

// abpDecryptSecret reverses what the provider stored in `encrypted_secret`.
func abpDecryptSecret(t *testing.T, index int, encrypted string) string {
	t.Helper()
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		t.Fatalf("decoding encrypted_secret: %s", err)
	}
	plain, err := rsa.DecryptOAEP(sha256.New(), cryptorand.Reader, abpTestRsaKeys()[index], raw, []byte(abpCredentialEncryptionLabel))
	if err != nil {
		t.Fatalf("decrypting encrypted_secret: %s", err)
	}
	return string(plain)
}

// --- what tests use to seed and inspect the fake ----------------------------

// seed stores an object that was created outside Terraform, the way the UI or
// an earlier workspace would have. The returned id is filled into the document
// unless it already carries one.
func (f *abpFake) seed(collection string, doc map[string]any) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.seedLocked(collection, doc)
}

func (f *abpFake) seedLocked(collection string, doc map[string]any) string {
	stored := abpFakeClone(doc)
	if abpFakeString(stored["id"]) == "" {
		stored["id"] = f.nextId(collection)
	}
	if _, ok := stored["account_id"]; !ok && collection != abpFakeKeys {
		stored["account_id"] = abpFakeAccountId
	}
	if _, ok := stored["created_at"]; !ok {
		stored["created_at"] = f.now()
	}
	if _, ok := stored["modified_at"]; !ok && collection != abpFakeKeys && collection != abpFakeCredentials {
		stored["modified_at"] = stored["created_at"]
	}
	f.insert(collection, stored)
	return abpFakeString(stored["id"])
}

// seedSite stores a Site with the Selectors given as resolved criteria maps,
// plus the catch-all the backend maintains.
func (f *abpFake) seedSite(name string, criteria ...map[string]any) string {
	f.mu.Lock()
	defer f.mu.Unlock()

	selectors := []any{}
	for _, c := range criteria {
		selectors = append(selectors, map[string]any{
			"criteria":          c,
			"analysis_settings": map[string]any{"rate_limiting": "per_site"},
		})
	}
	site := map[string]any{"name": name, "selectors": selectors}
	id := f.seedLocked(abpFakeSites, site)
	stored := f.store[abpFakeSites][id]
	f.normalizeSelectors(stored)
	f.addDefaultSelector(stored)
	return id
}

// seedCondition stores a literal Condition. An empty accountId makes it a
// managed (Imperva-owned) Condition.
func (f *abpFake) seedCondition(accountId, name, description, code string) string {
	return f.seed(abpFakeConditions, map[string]any{
		"account_id":     accountId,
		"name":           name,
		"description":    description,
		"code":           abpFakeFormatCode(code),
		"template":       "custom",
		"last_change_by": "fake",
	})
}

// seedConditionList stores a Condition List.
func (f *abpFake) seedConditionList(accountId, name, description string) string {
	return f.seed(abpFakeConditions, map[string]any{
		"account_id":  accountId,
		"name":        name,
		"description": description,
	})
}

// seedEntry attaches a Condition to a Condition List, the way
// incapsula_abp_condition_list_entry does.
func (f *abpFake) seedEntry(parent, reference, state string, tags ...string) string {
	list := make([]any, 0, len(tags))
	for _, tag := range tags {
		list = append(list, tag)
	}
	return f.seed(abpFakeConditions, map[string]any{
		"reference": reference,
		"parent":    parent,
		"tags":      list,
		"state":     state,
	})
}

// seedPolicy stores a Policy with one Directive per action, each with its own
// Condition List.
func (f *abpFake) seedPolicy(name, description string, actions ...string) string {
	f.mu.Lock()
	defer f.mu.Unlock()

	directives := []any{}
	for _, action := range actions {
		directive := map[string]any{"action": action, "condition_id": f.newConditionList(abpFakeAccountId, action)}
		if action == "proof_of_work" {
			directive["skip_condition_id"] = f.newConditionList(abpFakeAccountId, action+"_skip")
		}
		directives = append(directives, directive)
	}
	return f.seedLocked(abpFakePolicies, map[string]any{
		"name":        name,
		"description": description,
		"directives":  directives,
	})
}

// seedGlobalPolicy makes a Policy the one the Account's global policy endpoint
// returns. It is reachable by id and through that endpoint, but never through
// the Account's Policy list.
func (f *abpFake) seedGlobalPolicy(name string, actions ...string) string {
	id := f.seedPolicy(name, "the account's global policy", actions...)
	f.mu.Lock()
	defer f.mu.Unlock()
	f.globalPolicyId = id
	return id
}

// seedDomain attaches a Domain to a Site, cloning the given encryption key
// exactly as the backend would. An empty encryptionKeyId means none.
func (f *abpFake) seedDomain(siteId, fqdn, encryptionKeyId string) string {
	id := f.seed(abpFakeDomains, map[string]any{
		"site_id":                  siteId,
		"criteria":                 map[string]any{"exact": fqdn},
		"cookiescope":              fqdn,
		"challenge_ip_lookup_mode": "none",
		"analysis_ip_lookup_mode":  "none",
		"captcha_settings":         "none",
		"log_region":               "us",
		"cookie_mode":              "on",
		"no_js_injection_paths":    []any{},
		"unmasked_headers":         []any{},
		"proxy_flags":              []any{},
		"enable_mitigation":        true,
	})
	if encryptionKeyId != "" {
		f.mu.Lock()
		f.encryptionKeyOf[id] = encryptionKeyId
		f.mu.Unlock()
	}
	return id
}

// seedCredential stores a Credential with a server-generated secret.
func (f *abpFake) seedCredential() map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	return abpFakeClone(f.newCredential(abpFakeAccountId))
}

// seedEncryptionKey attaches key material to a Domain.
func (f *abpFake) seedEncryptionKey(domainId, key string) string {
	return f.seed(abpFakeKeys, map[string]any{
		"domain_id":          domainId,
		"key":                key,
		"first_published_at": nil,
	})
}

// seedPowConfig stores a Proof of Work Configuration.
func (f *abpFake) seedPowConfig(name string, difficulty int) string {
	return f.seed(abpFakePowConfigs, map[string]any{
		"name":       name,
		"difficulty": difficulty,
		"algorithm":  "bbs",
	})
}

// doc returns a copy of one stored object, or nil when it does not exist.
func (f *abpFake) doc(collection, id string) map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	stored := f.store[collection][id]
	if stored == nil {
		return nil
	}
	return abpFakeClone(stored)
}

// exists reports whether an object is still stored.
func (f *abpFake) exists(collection, id string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.store[collection][id] != nil
}

// all returns copies of every object in a collection, in creation order.
func (f *abpFake) all(collection string) []map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := []map[string]any{}
	for _, id := range f.order[collection] {
		if doc := f.store[collection][id]; doc != nil {
			out = append(out, abpFakeClone(doc))
		}
	}
	return out
}

// only returns the single object in a collection, failing the test otherwise.
func (f *abpFake) only(t *testing.T, collection string) map[string]any {
	t.Helper()
	docs := f.all(collection)
	if len(docs) != 1 {
		t.Fatalf("expected exactly one %s, got %d", collection, len(docs))
	}
	return docs[0]
}

// count returns how many objects a collection holds.
func (f *abpFake) count(collection string) int {
	return len(f.all(collection))
}

// mutate changes a stored object out of band, standing in for the UI or
// another workspace editing it.
func (f *abpFake) mutate(t *testing.T, collection, id string, change func(doc map[string]any)) {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	stored := f.store[collection][id]
	if stored == nil {
		t.Fatalf("%s %s is not stored", collection, id)
	}
	change(stored)
}

// drop deletes an object out of band.
func (f *abpFake) drop(collection, id string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.remove(collection, id)
}

// deletedIds returns the ids the API was asked to delete, in order.
func (f *abpFake) deletedIds() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.deleted...)
}

// writes returns every request that changed something, in order.
func (f *abpFake) writes() []abpFakeRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := []abpFakeRequest{}
	for _, request := range f.requests {
		if request.Method != http.MethodGet {
			out = append(out, request)
		}
	}
	return out
}

// lastWrite returns the body of the most recent request that changed
// something, parsed, for diagnosing what the provider actually sent.
func (f *abpFake) lastWrite(t *testing.T) map[string]any {
	t.Helper()
	writes := f.writes()
	for i := len(writes) - 1; i >= 0; i-- {
		if writes[i].Body == "" {
			continue
		}
		var body map[string]any
		if err := json.Unmarshal([]byte(writes[i].Body), &body); err != nil {
			t.Fatalf("recorded request is not JSON: %s", err)
		}
		return body
	}
	t.Fatal("no request with a body was recorded")
	return nil
}

// requestCount returns how many requests of a method had a path containing the
// given fragment.
func (f *abpFake) requestCount(method, contains string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	count := 0
	for _, request := range f.requests {
		if request.Method == method && strings.Contains(request.Path, contains) {
			count++
		}
	}
	return count
}

// encryptionKeyIdOf returns the `encryption_key_id` a Domain was created with.
func (f *abpFake) encryptionKeyIdOf(domainId string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.encryptionKeyOf[domainId]
}

// siteOrderOf returns an Account's Site priority order.
func (f *abpFake) siteOrderOf(accountId string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.siteOrder(accountId)
}

// domainOrderOf returns a Site's Domain priority order.
func (f *abpFake) domainOrderOf(siteId string) []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.domainOrder(siteId)
}

// reorderSites changes an Account's Site priority out of band. The list must
// cover the Account's Sites, as the API requires.
func (f *abpFake) reorderSites(ids ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, id := range ids {
		f.sitePriority[id] = i + 1
	}
}

// reorderDomains changes a Site's Domain priority out of band.
func (f *abpFake) reorderDomains(ids ...string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i, id := range ids {
		f.domainPriority[id] = i + 1
	}
}

// publishElsewhere takes a snapshot and publishes it without Terraform,
// standing in for another operator or pipeline publishing the Account. Every
// preflight taken before it becomes unpublishable.
func (f *abpFake) publishElsewhere() {
	f.mu.Lock()
	defer f.mu.Unlock()
	entity := f.newPreflight(abpFakeAccountId)
	id := abpFakeString(entity["id"])
	f.publishedSeq = f.preflightSeq[id]
	f.published = append(f.published, id)
}

// setRejectPublish makes the next publish attempt fail, standing in for a
// transient failure or an interrupted apply.
func (f *abpFake) setRejectPublish(reject bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.rejectPublish = reject
}

// setPreflightOutdated adds the `outdated` warning to every preflight,
// standing in for an Account changed after the snapshot was taken.
func (f *abpFake) setPreflightOutdated(outdated bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.preflightOutdated = outdated
}

// setPreflightNoWarnings makes the service report `warnings: null`.
func (f *abpFake) setPreflightNoWarnings(none bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.preflightNoWarnings = none
}

// publishedIds returns the preflight ids that were published, in order.
func (f *abpFake) publishedIds() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]string{}, f.published...)
}

// preflightCount returns how many preflights were taken.
func (f *abpFake) preflightCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.preflightSeq)
}
