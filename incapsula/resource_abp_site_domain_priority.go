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

		Description: `Manages the priority order of Domains within an ABP Site. The first ID in
` + "`domain_ids`" + ` has the highest priority. The list must include every Domain
attached to the Site exactly once.

NOTE: this resource wraps a singleton property of the parent Site. Destroying
it removes the priority from Terraform state but leaves the server-side order
unchanged.`,

		Schema: map[string]*schema.Schema{
			"site_id": {
				Description:  "Site UUID whose Domain priority order is being managed.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"domain_ids": {
				Description: "Domain UUIDs in priority order (highest first). Must include every Domain attached to the Site exactly once.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsUUID,
				},
			},
		},
	}
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

	updated, err := client.UpdateAbpSiteDomainPriority(siteId, dp)
	if err != nil {
		return diag.FromErr(err)
	}
	if updated == nil {
		return diag.Errorf("ABP Site %s not found when setting domain_priority", siteId)
	}

	data.SetId(siteId)
	if err := data.Set("domain_ids", updated.DomainIds); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Set ABP Site %s domain_priority to %d entries", siteId, len(updated.DomainIds))
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
	if err := data.Set("domain_ids", dp.DomainIds); err != nil {
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
	if err := data.Set("domain_ids", dp.DomainIds); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
