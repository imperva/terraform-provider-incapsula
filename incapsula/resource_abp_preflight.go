package incapsula

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAbpPreflight() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPreflightCreate,
		ReadContext:   resourceAbpPreflightRead,
		DeleteContext: resourceAbpPreflightDelete,

		Description: "Creates a new ABP preflight\n",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account this preflight belongs to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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

func resourceAbpPreflightCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func resourceAbpPreflightRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func resourceAbpPreflightDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Preflights don't need to be deleted, forgetting it is fine
	data.SetId("")
	return nil
}
