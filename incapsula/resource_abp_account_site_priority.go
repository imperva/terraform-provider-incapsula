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

func resourceAbpAccountSitePriority() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpAccountSitePriorityUpsert,
		ReadContext:   resourceAbpAccountSitePriorityRead,
		UpdateContext: resourceAbpAccountSitePriorityUpsert,
		DeleteContext: resourceAbpAccountSitePriorityDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpAccountSitePriorityImport,
		},

		Description: `Manages the priority order of Sites within an ABP Account. The first ID in
` + "`site_ids`" + ` has the highest priority. The list must include every Site
attached to the Account exactly once.

NOTE: this resource wraps a singleton property of the parent Account.
Destroying it removes the priority from Terraform state but leaves the
server-side order unchanged.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "Account UUID whose Site priority order is being managed.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"site_ids": {
				Description: "Site UUIDs in priority order (highest first). Must include every Site attached to the Account exactly once.",
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

func extractAbpAccountSitePriority(data *schema.ResourceData) (AbpAccountSitePriority, error) {
	raw := data.Get("site_ids").([]any)
	ids := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for i, v := range raw {
		s, _ := v.(string)
		if s == "" {
			return AbpAccountSitePriority{}, fmt.Errorf("site_ids[%d] is empty", i)
		}
		if _, dup := seen[s]; dup {
			return AbpAccountSitePriority{}, fmt.Errorf("site_ids[%d] is a duplicate of an earlier entry: %s", i, s)
		}
		seen[s] = struct{}{}
		ids = append(ids, s)
	}
	return AbpAccountSitePriority{SiteIds: ids}, nil
}

func resourceAbpAccountSitePriorityUpsert(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	sp, err := extractAbpAccountSitePriority(data)
	if err != nil {
		return diag.FromErr(err)
	}

	updated, err := client.UpdateAbpAccountSitePriority(accountId, sp)
	if err != nil {
		return diag.FromErr(err)
	}
	if updated == nil {
		return diag.Errorf("ABP Account %s not found when setting site_priority", accountId)
	}

	data.SetId(accountId)
	if err := data.Set("site_ids", updated.SiteIds); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Set ABP Account %s site_priority to %d entries", accountId, len(updated.SiteIds))
	return nil
}

func resourceAbpAccountSitePriorityRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Id()

	sp, err := client.ReadAbpAccountSitePriority(accountId)
	if err != nil {
		return diag.FromErr(err)
	}
	if sp == nil {
		log.Printf("[INFO] ABP Account %s not found, removing site_priority from state", accountId)
		data.SetId("")
		return nil
	}

	if err := data.Set("account_id", accountId); err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("site_ids", sp.SiteIds); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// resourceAbpAccountSitePriorityDelete leaves the server-side priority order
// untouched. There is no DELETE endpoint for site_priority — it is a
// singleton property of the Account — so unmanaging it just drops state.
func resourceAbpAccountSitePriorityDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	data.SetId("")
	return nil
}

func resourceAbpAccountSitePriorityImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	accountId := strings.TrimSpace(data.Id())
	if accountId == "" {
		return nil, fmt.Errorf("expected import ID to be '<account_id>'")
	}

	client := m.(*Client)
	sp, err := client.ReadAbpAccountSitePriority(accountId)
	if err != nil {
		return nil, err
	}
	if sp == nil {
		return nil, fmt.Errorf("ABP Account %s not found", accountId)
	}

	data.SetId(accountId)
	if err := data.Set("account_id", accountId); err != nil {
		return nil, err
	}
	if err := data.Set("site_ids", sp.SiteIds); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
