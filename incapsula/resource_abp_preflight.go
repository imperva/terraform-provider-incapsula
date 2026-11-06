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
rollback as long as it is still consistent with any external state.

NOTE: ` + "`pending_hash`" + ` is what makes the whole pipeline work, even
though it is optional. It should always be set to
` + "`data.incapsula_abp_pending_changes.<name>.hash`" + `: that both triggers a
new preflight when the configuration changes and orders the preflight after the
resources it is supposed to snapshot. Omitting it means the preflight is created
once, in parallel with those resources, and is never refreshed again. Only the
resources reachable through that data source's ` + "`depends_on`" + ` take part
in this, so every ABP resource must be listed there (most simply: keep them all
in one module and depend on that module).

NOTE: this resource does not re-read the preflight from the server, so a
preflight that the server has invalidated or deleted still looks healthy in
state and ` + "`terraform plan`" + ` reports no changes. Warnings the server
attaches to a preflight — including the one saying that changes made after the
snapshot will not be published — are not surfaced either. If publishing starts
failing with "ABP preflight is not publishable", force a new one with
` + "`terraform apply -replace=incapsula_abp_preflight.<name>`" + `. Importing an
existing preflight is not supported, and ` + "`terraform destroy`" + ` only drops
it from state: nothing is deleted or published, so an account destroyed through
Terraform is left with the deleted objects still live until something publishes
again.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "ABP account UUID this preflight belongs to. Changing it replaces " +
					"the resource: a preflight belongs to one account, so Terraform forgets this " +
					"one (it stays on the server) and creates a preflight for the new account.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"pending_hash": {
				Description: "Change-detection value; set it to " +
					"`data.incapsula_abp_pending_changes.<name>.hash`. Any change to it replaces " +
					"the resource, which is how a new snapshot of the configuration is taken. " +
					"Strongly recommended: see the notes above for what omitting it costs.",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
