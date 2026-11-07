package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"sync"
)

// Preserving the fields this provider does not model
// --------------------------------------------------
//
// Every ABP update endpoint is a full replace: the API applies the request body
// as the new document and resets whatever the body does not mention. The client
// structs in this package, on the other hand, are only a snapshot of the wire
// model as it looked when they were written, so every
// field the backend adds later -- or that nobody got round to modelling -- is
// silently cleared by the next unrelated update. That is how `industry`,
// selector `purpose` and `path_prefix.case_insensitive` used to disappear from
// an ABP Site when all the user changed was its name.
//
// The update path is therefore a read-modify-write: the document the API
// currently holds is fetched, the payload this provider built is overlaid onto
// it, and the merged document is sent. Keys the provider does not model ride
// along untouched; keys it does model are dictated by the payload, including
// by their absence, so clearing a value still clears it.
//
// What makes that safe here is that the update endpoints ignore unknown keys,
// so the keys a read returns but an update model lacks (`id`, `created_at`, a
// Domain's `criteria`, ...) are ignored
// rather than rejected. Verified against the public API: PUT accepts a
// verbatim GET body on both /v1/site/{id} and /v1/domain/{id}.
//
// Three rules keep the merge from inventing configuration of its own:
//
//   - A key is "owned" when the payload type models it. Owned keys are taken
//     from the payload and from nowhere else, so an omitted optional still
//     means "clear this".
//   - A "oneOf" object (selector criteria, captcha settings, ...) that switched
//     variant has its leftovers dropped instead of merged. Merging them would
//     be worse than clearing: the API picks the variant by key presence in a
//     fixed order, so a body carrying both `path_prefix` and `path_regex`
//     silently keeps the old prefix.
//   - A list is merged element-wise only when its elements carry an `id`.
//     Position is not identity: policy directives have no id, and index-merging
//     them could move an inherited directive's `origin_policy_id` onto a
//     directive of the policy's own, which the API responds to by dropping it.

// abpWireOneOf describes the keys of a payload type that marshals as a JSON
// "oneOf" -- an object in which one key selects the shape of the rest. Such a
// type has no usable struct tags, so it declares its keys itself.
type abpWireOneOf struct {
	// Selectors are the keys that choose the variant. A change in which of
	// them is present means the value switched to another variant.
	Selectors []string
	// Extra are keys the type may emit alongside whichever variant is
	// selected, e.g. `case_insensitive` next to `path_prefix`.
	Extra []string
}

// abpWireVariant is implemented by the client's oneOf types.
type abpWireVariant interface {
	abpWireOneOf() abpWireOneOf
}

// abpWireServerOnlyKeys are keys the API returns but never accepts. They are
// dropped instead of being echoed back, to keep request bodies to what the
// endpoint actually reads. Keys that are merely read-only (`id`, `created_at`,
// `derived_id`, ...) need no special case: they are modelled, so the payload
// decides, and update payloads leave them empty.
var abpWireServerOnlyKeys = map[string]bool{
	"_links": true, // HAL navigation links
}

// abpWireModel records which JSON keys a payload type owns, so the merge can
// tell "the provider left this out" from "the provider knows nothing about
// this".
type abpWireModel struct {
	// fields maps a JSON key to the model of its value. A key present with a
	// nil model is owned but its contents are not described, so the payload
	// replaces them wholesale.
	fields map[string]*abpWireModel
	// elem models the elements of an array.
	elem *abpWireModel
	// oneOf is set for variant objects; Selectors identifies the variant.
	oneOf *abpWireOneOf
}

var abpWireModelCache sync.Map // reflect.Type -> *abpWireModel

// abpWireModelOf describes t, which must be the type of the value the client
// marshals into a request body.
func abpWireModelOf(t reflect.Type) *abpWireModel {
	if cached, ok := abpWireModelCache.Load(t); ok {
		return cached.(*abpWireModel)
	}
	m := buildAbpWireModel(t, map[reflect.Type]bool{})
	abpWireModelCache.Store(t, m)
	return m
}

var abpWireVariantType = reflect.TypeOf((*abpWireVariant)(nil)).Elem()

// abpWireVariantOf returns t's oneOf description, if t is a variant type.
func abpWireVariantOf(t reflect.Type) (abpWireVariant, bool) {
	zero := reflect.New(t)
	if t.Implements(abpWireVariantType) {
		return zero.Elem().Interface().(abpWireVariant), true
	}
	if reflect.PointerTo(t).Implements(abpWireVariantType) {
		return zero.Interface().(abpWireVariant), true
	}
	return nil, false
}

func buildAbpWireModel(t reflect.Type, seen map[reflect.Type]bool) *abpWireModel {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if variant, ok := abpWireVariantOf(t); ok {
		keys := variant.abpWireOneOf()
		fields := map[string]*abpWireModel{}
		for _, k := range slices.Concat(keys.Selectors, keys.Extra) {
			fields[k] = nil
		}
		return &abpWireModel{fields: fields, oneOf: &keys}
	}

	switch t.Kind() {
	case reflect.Struct:
		if seen[t] {
			return nil
		}
		seen[t] = true
		defer delete(seen, t)

		fields := map[string]*abpWireModel{}
		for i := range t.NumField() {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
			if name == "-" {
				continue
			}
			if name == "" {
				name = field.Name
			}
			fields[name] = buildAbpWireModel(field.Type, seen)
		}
		return &abpWireModel{fields: fields}
	case reflect.Slice, reflect.Array:
		if t.Elem().Kind() == reflect.Uint8 {
			return nil // []byte marshals as a string
		}
		return &abpWireModel{elem: buildAbpWireModel(t.Elem(), seen)}
	default:
		// Scalars, maps and anything else opaque: the payload replaces it.
		return nil
	}
}

// owns reports whether the payload type models key.
func (m *abpWireModel) owns(key string) bool {
	if m == nil {
		return false
	}
	_, ok := m.fields[key]
	return ok
}

func (m *abpWireModel) child(key string) *abpWireModel {
	if m == nil {
		return nil
	}
	return m.fields[key]
}

// variantSwitched reports whether base and overlay select different variants
// of a oneOf object.
func (m *abpWireModel) variantSwitched(base, overlay map[string]any) bool {
	if m == nil || m.oneOf == nil {
		return false
	}
	present := func(obj map[string]any) []string {
		var keys []string
		for _, k := range m.oneOf.Selectors {
			if v, ok := obj[k]; ok && v != nil {
				keys = append(keys, k)
			}
		}
		slices.Sort(keys)
		return keys
	}
	return !slices.Equal(present(base), present(overlay))
}

// abpMergeWireDocuments overlays overlay onto base. model describes the type
// overlay was marshalled from.
func abpMergeWireDocuments(base, overlay []byte, model *abpWireModel) ([]byte, error) {
	baseValue, err := decodeAbpWireDocument(base)
	if err != nil {
		return nil, fmt.Errorf("decoding the current document: %w", err)
	}
	overlayValue, err := decodeAbpWireDocument(overlay)
	if err != nil {
		return nil, fmt.Errorf("decoding the request body: %w", err)
	}

	merged, err := json.Marshal(mergeAbpWireValue(baseValue, overlayValue, model))
	if err != nil {
		return nil, fmt.Errorf("encoding the merged request body: %w", err)
	}
	return merged, nil
}

// decodeAbpWireDocument decodes into any, keeping numbers verbatim so that
// round-tripping a large integer cannot turn it into 1.234e+18.
func decodeAbpWireDocument(data []byte) (any, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	return value, nil
}

func mergeAbpWireValue(base, overlay any, model *abpWireModel) any {
	switch overlayValue := overlay.(type) {
	case map[string]any:
		baseValue, ok := base.(map[string]any)
		if !ok {
			return overlayValue // the shape changed; there is nothing to keep
		}
		return mergeAbpWireObject(baseValue, overlayValue, model)
	case []any:
		baseValue, ok := base.([]any)
		if !ok {
			return overlayValue
		}
		return mergeAbpWireList(baseValue, overlayValue, model.elemModel())
	default:
		return overlay
	}
}

func (m *abpWireModel) elemModel() *abpWireModel {
	if m == nil {
		return nil
	}
	return m.elem
}

func mergeAbpWireObject(base, overlay map[string]any, model *abpWireModel) map[string]any {
	merged := make(map[string]any, len(base)+len(overlay))

	if !model.variantSwitched(base, overlay) {
		for key, value := range base {
			if model.owns(key) || abpWireServerOnlyKeys[key] {
				continue
			}
			merged[key] = value
		}
	}

	for key, value := range overlay {
		if baseValue, ok := base[key]; ok {
			merged[key] = mergeAbpWireValue(baseValue, value, model.child(key))
			continue
		}
		merged[key] = value
	}
	return merged
}

// mergeAbpWireList merges elements that both sides identify with the same
// `id`, keeping the overlay's order and dropping elements the overlay left
// out. Elements without an `id` are taken from the overlay as they are.
func mergeAbpWireList(base, overlay []any, elem *abpWireModel) []any {
	byId := make(map[string]map[string]any, len(base))
	for _, item := range base {
		if object, ok := item.(map[string]any); ok {
			if id, ok := object["id"].(string); ok && id != "" {
				byId[id] = object
			}
		}
	}

	merged := make([]any, 0, len(overlay))
	for _, item := range overlay {
		object, ok := item.(map[string]any)
		if !ok {
			merged = append(merged, item)
			continue
		}
		id, ok := object["id"].(string)
		if !ok || byId[id] == nil {
			merged = append(merged, item)
			continue
		}
		merged = append(merged, mergeAbpWireObject(byId[id], object, elem))
	}
	return merged
}

// abpPreserveUnmodelledFields fetches the document at url and overlays payload
// onto it, so that the full-replace update about to be sent does not clear the
// fields this provider does not model. model is the type payload was
// marshalled from.
//
// A payload that cannot be reconciled with what the API holds is an error:
// sending it anyway is what used to destroy configuration. A 404 is not, since
// the update itself reports a vanished object better than a read can.
func (c *Client) abpPreserveUnmodelledFields(url string, payload []byte, model reflect.Type, resourceName, operation string) ([]byte, error) {
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, url, nil, operation)
	if err != nil {
		return nil, fmt.Errorf("error reading current %s before updating it: %w", resourceName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return payload, nil
	}

	current, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body of the current %s: %w", resourceName, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading current %s before updating it: %s", resp.StatusCode, resourceName, string(current))
	}

	merged, err := abpMergeWireDocuments(current, payload, abpWireModelOf(model))
	if err != nil {
		return nil, fmt.Errorf("error preserving unmodelled fields of %s: %w", resourceName, err)
	}
	return merged, nil
}
