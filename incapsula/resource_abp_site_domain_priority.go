package incapsula

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAbpSiteDomainPriority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpSiteDomainPriorityUpsert,
		ReadContext:   resourceAbpSiteDomainPriorityRead,
		UpdateContext: resourceAbpSiteDomainPriorityUpsert,
		DeleteContext: resourceAbpSiteDomainPriorityDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpSiteDomainPriorityImport,
		},

		CustomizeDiff: resourceAbpSiteDomainPriorityDiff,

		Description: `Manages the priority order of Domains within an ABP Site. The first ID in
` + "`domain_ids`" + ` has the highest priority.

` + "`domain_ids`" + ` does not have to cover every Domain of the Site. The
Domains it lists become the highest priority ones, in the order given, and every
other Domain attached to the Site keeps its relative order behind them. A Domain
created elsewhere — in the ABP UI, by another workspace or by another module —
therefore becomes the lowest priority Domain of the Site and does not have to be
added here. ` + "`effective_domain_ids`" + ` reports the resulting full order.

NOTE: this resource wraps a singleton property of the parent Site. Destroying
it removes the priority from Terraform state but leaves the server-side order
unchanged.`,

		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Site UUID whose Domain priority order is being managed. Changing it " +
					"replaces the resource: Terraform forgets the order of the old Site (leaving " +
					"it unchanged on the server) and applies this order to the new Site.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"domain_ids": {
				Description: "Domain UUIDs in priority order (highest first). Every ID must belong " +
					"to a Domain attached to the Site, but the list does not have to cover all of " +
					"them: the Domains it leaves out keep their relative order behind the ones it " +
					"lists.",
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
			"effective_domain_ids": {
				Description: "Every Domain attached to the Site in the priority order the Site now " +
					"has: the Domains listed in `domain_ids` first, then the ones this resource " +
					"does not rank.",
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

// abpPriorityIdList reads a list-of-strings attribute value, skipping empty
// entries.
func abpPriorityIdList(value any) []string {
	raw, _ := value.([]any)
	ids := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, _ := v.(string); s != "" {
			ids = append(ids, s)
		}
	}
	return ids
}

// abpPriorityCompleteOrder puts the ranked ids first, in the order given, and
// appends every other attached id behind them, keeping their current relative
// order. The result covers `attached` exactly, which is what the API requires
// of a priority list. The ranked ids must all be attached.
func abpPriorityCompleteOrder(ranked, attached []string) []string {
	isRanked := make(map[string]struct{}, len(ranked))
	for _, id := range ranked {
		isRanked[id] = struct{}{}
	}

	order := make([]string, 0, len(attached))
	order = append(order, ranked...)
	for _, id := range attached {
		if _, ok := isRanked[id]; !ok {
			order = append(order, id)
		}
	}
	return order
}

// abpPriorityUnattachedIds returns the ranked ids that are not among the
// attached ones, in the order they were ranked.
func abpPriorityUnattachedIds(ranked, attached []string) []string {
	isAttached := make(map[string]struct{}, len(attached))
	for _, id := range attached {
		isAttached[id] = struct{}{}
	}

	var unattached []string
	for _, id := range ranked {
		if _, ok := isAttached[id]; !ok {
			unattached = append(unattached, id)
		}
	}
	return unattached
}

// abpPriorityProject returns the ids of `keep` in the order they appear in
// `order`, dropping everything else. It is how the server's full order is read
// back into an attribute that only ranks part of it.
func abpPriorityProject(order, keep []string) []string {
	isKept := make(map[string]struct{}, len(keep))
	for _, id := range keep {
		isKept[id] = struct{}{}
	}

	projected := make([]string, 0, len(keep))
	for _, id := range order {
		if _, ok := isKept[id]; ok {
			projected = append(projected, id)
		}
	}
	return projected
}

// abpPriorityHasPrefix reports whether order starts with prefix.
func abpPriorityHasPrefix(order, prefix []string) bool {
	if len(order) < len(prefix) {
		return false
	}
	for i, id := range prefix {
		if order[i] != id {
			return false
		}
	}
	return true
}

// resourceAbpSiteDomainPriorityDiff plans an update when the ranked Domains are
// no longer the highest priority ones. A reordering among them is already
// reported through `domain_ids`, but an unranked Domain that was moved ahead of
// them only shows in `effective_domain_ids`, and without this nothing would ask
// Terraform to put the ranked Domains back on top.
func resourceAbpSiteDomainPriorityDiff(ctx context.Context, diff *schema.ResourceDiff, m any) error {
	if diff.Id() == "" {
		return nil
	}

	ranked := abpPriorityIdList(diff.Get("domain_ids"))
	effective := abpPriorityIdList(diff.Get("effective_domain_ids"))
	if len(ranked) == 0 || abpPriorityHasPrefix(effective, ranked) {
		return nil
	}
	return diff.SetNewComputed("effective_domain_ids")
}

func extractAbpSiteDomainPriority(data *schema.ResourceData) (AbpSiteDomainPriority, error) {
	raw := data.Get("domain_ids").([]any)
	ids := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for i, v := range raw {
		s, _ := v.(string)
		if s == "" {
			return AbpSiteDomainPriority{}, fmt.Errorf("domain_ids[%d] is empty", i)
		}
		if _, dup := seen[s]; dup {
			return AbpSiteDomainPriority{}, fmt.Errorf("domain_ids[%d] is a duplicate of an earlier entry: %s", i, s)
		}
		seen[s] = struct{}{}
		ids = append(ids, s)
	}
	return AbpSiteDomainPriority{DomainIds: ids}, nil
}

func resourceAbpSiteDomainPriorityUpsert(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	siteId := data.Get("site_id").(string)

	dp, err := extractAbpSiteDomainPriority(data)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API requires the submitted order to cover exactly the Domains attached
	// to the Site, so the ranked ones are completed with the rest from a fresh
	// read. A Domain created or deleted in between makes the write fail, which a
	// second read and write resolve, so the write gets one more attempt.
	var updated *AbpSiteDomainPriority
	for attempt := 1; attempt <= 2; attempt++ {
		attached, err := client.ReadAbpSiteDomainPriority(siteId)
		if err != nil {
			return diag.FromErr(err)
		}
		if attached == nil {
			return diag.Errorf("ABP Site %s not found when setting domain_priority", siteId)
		}

		if unattached := abpPriorityUnattachedIds(dp.DomainIds, attached.DomainIds); len(unattached) > 0 {
			return diag.Errorf("domain_ids lists %d Domain(s) that are not attached to ABP Site %s: %s. "+
				"A Domain can only be prioritised in the Site it belongs to, and this Site has %d Domain(s)",
				len(unattached), siteId, strings.Join(unattached, ", "), len(attached.DomainIds))
		}

		order := AbpSiteDomainPriority{DomainIds: abpPriorityCompleteOrder(dp.DomainIds, attached.DomainIds)}
		updated, err = client.UpdateAbpSiteDomainPriority(siteId, order)
		if err == nil {
			break
		}
		if attempt == 2 {
			return diag.FromErr(err)
		}
		log.Printf("[INFO] Setting ABP Site %s domain_priority failed (%s); its Domains may have changed, re-reading them and trying once more", siteId, err)
	}
	if updated == nil {
		return diag.Errorf("ABP Site %s not found when setting domain_priority", siteId)
	}

	data.SetId(siteId)
	if err := data.Set("domain_ids", dp.DomainIds); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("effective_domain_ids", updated.DomainIds); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Set ABP Site %s domain_priority: %d ranked Domain(s) of %d attached", siteId, len(dp.DomainIds), len(updated.DomainIds))
	return nil
}

func resourceAbpSiteDomainPriorityRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	siteId := data.Id()

	dp, err := client.ReadAbpSiteDomainPriority(siteId)
	if err != nil {
		return diag.FromErr(err)
	}
	if dp == nil {
		log.Printf("[INFO] ABP Site %s not found, removing domain_priority from state", siteId)
		data.SetId("")
		return nil
	}

	if err := data.Set("site_id", siteId); err != nil {
		return diag.FromErr(err)
	}
	// Only the ranked Domains belong in `domain_ids`. Projecting the server's
	// order onto them reports a reordering among them as drift, without pulling
	// a Domain the configuration does not rank into the configuration's list —
	// which would propose a removal the API cannot perform.
	ranked := abpPriorityIdList(data.Get("domain_ids"))
	if err := data.Set("domain_ids", abpPriorityProject(dp.DomainIds, ranked)); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("effective_domain_ids", dp.DomainIds); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// resourceAbpSiteDomainPriorityDelete leaves the server-side priority order
// untouched. There is no DELETE endpoint for domain_priority — it is a
// singleton property of the Site — so unmanaging it just drops state.
func resourceAbpSiteDomainPriorityDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	data.SetId("")
	return nil
}

func resourceAbpSiteDomainPriorityImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	siteId := strings.TrimSpace(data.Id())
	if siteId == "" {
		return nil, fmt.Errorf("expected import ID to be '<site_id>'")
	}

	client := m.(*Client)
	dp, err := client.ReadAbpSiteDomainPriority(siteId)
	if err != nil {
		return nil, err
	}
	if dp == nil {
		return nil, fmt.Errorf("ABP Site %s not found", siteId)
	}

	data.SetId(siteId)
	if err := data.Set("site_id", siteId); err != nil {
		return nil, err
	}
	// Importing adopts the Site's whole order, so every Domain starts out ranked.
	// Dropping one from `domain_ids` afterwards leaves it as the lowest priority
	// Domain rather than failing the apply.
	if err := data.Set("domain_ids", dp.DomainIds); err != nil {
		return nil, err
	}
	if err := data.Set("effective_domain_ids", dp.DomainIds); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
