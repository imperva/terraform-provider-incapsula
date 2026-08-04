package incapsula

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAbpPreflight() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPreflightCreate,
		ReadContext:   resourceAbpPreflightRead,
		DeleteContext: resourceAbpPreflightDelete,

		Description: `Creates a new ABP preflight. A preflight is a snapshot of the entire
account configuration and is required in order to publish that configuration to
the Analysis Host (see ` + "`incapsula_abp_publish`" + `).

A preflight may become invalid and be deleted by the server if the account
configuration or external dependent state changes, or if too much time has
passed since it was created. Publishing an older preflight can act as a
rollback as long as it is still consistent with any external state.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this preflight belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"pending_hash": {
				Description: "Any string to facilitate change detection, using the hash from `data.incapsula_abp_pending_changes` causes a replacement of this resource (and thereby a new preflight)\n",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAbpPreflightCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	preflight, diags := client.CreateAbpPreflight(accountId)
	if diags.HasError() {
		return diags
	}

	log.Printf("[INFO] Created ABP preflight %s for account %s", preflight.Id, accountId)
	data.SetId(preflight.Id)
	return diags
}

func resourceAbpPreflightRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceAbpPreflightDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	// Preflights don't need to be deleted, forgetting it is fine
	data.SetId("")
	return nil
}
