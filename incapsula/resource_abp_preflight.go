package incapsula

import (
	"context"
	"fmt"
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

NOTE: a snapshot that can no longer be used is replaced automatically. On every
refresh the snapshot in state is re-read, and if it has been deleted by the
server, or can no longer be published because the account has been published
from elsewhere in the meantime, Terraform plans a new one — and a new publish
with it. So a publish that failed, or that someone else overtook, is recovered
by the next ` + "`terraform apply`" + ` instead of failing forever.

A snapshot that is merely *outdated* — the configuration changed after it was
taken — is deliberately kept, because changes Terraform does not manage do not
trigger a publish. Such a snapshot is reported as a warning when it is created,
and the reason is also available in ` + "`warnings`" + `.

` + "`terraform destroy`" + ` only drops the snapshot from state: nothing is
deleted or published, so an account destroyed through Terraform is left with the
deleted objects still live until something publishes again. Importing an
existing preflight is not supported.`,

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
			"created_at": {
				Description: "When the snapshot was taken.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"warnings": {
				Description: "Reasons why publishing this snapshot may not have the expected " +
					"effect, as reported when it was taken. Each entry has a `type` and a " +
					"human-readable `message`; an empty list means there is nothing to report.",
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Description: "Short identifier of the warning, for example `outdated`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"message": {
							Description: "Explanation of the warning.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func flattenAbpPreflightWarnings(warnings []AbpPreflightWarning) []any {
	flattened := make([]any, 0, len(warnings))
	for _, warning := range warnings {
		flattened = append(flattened, map[string]any{
			"type":    warning.Type,
			"message": warning.Message,
		})
	}
	return flattened
}

func setAbpPreflightAttributes(data *schema.ResourceData, preflight *AbpPreflight) diag.Diagnostics {
	var diags diag.Diagnostics
	for key, value := range map[string]any{
		"account_id": preflight.AccountId,
		"created_at": preflight.CreatedAt,
		"warnings":   flattenAbpPreflightWarnings(preflight.Warnings),
	} {
		if err := data.Set(key, value); err != nil {
			diags = append(diags, diag.FromErr(err)...)
		}
	}
	return diags
}

// abpPreflightWarningDiagnostics surfaces what the service says about a
// snapshot it has just taken. These are warnings rather than errors: the
// snapshot is publishable, it just may not do everything the user expects.
func abpPreflightWarningDiagnostics(preflight *AbpPreflight) diag.Diagnostics {
	var diags diag.Diagnostics
	for _, warning := range preflight.Warnings {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  fmt.Sprintf("ABP preflight %s was created with a warning: %s", preflight.Id, warning.Type),
			Detail:   warning.Message,
		})
	}
	return diags
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
	diags = append(diags, setAbpPreflightAttributes(data, preflight)...)
	return append(diags, abpPreflightWarningDiagnostics(preflight)...)
}

func resourceAbpPreflightRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)

	preflight, diags := client.GetAbpPreflight(data.Id())
	if diags.HasError() {
		return diags
	}

	// A preflight that is gone, or that can no longer be published, is of no
	// use to anyone: forget it so that a new one is planned instead of every
	// later publish failing on this one.
	if preflight == nil {
		log.Printf("[INFO] ABP preflight %s no longer exists, a new one will be created", data.Id())
		data.SetId("")
		return diags
	}

	if !preflight.CanPublish() {
		log.Printf("[INFO] ABP preflight %s can no longer be published, a new one will be created", data.Id())
		data.SetId("")
		return diags
	}

	return append(diags, setAbpPreflightAttributes(data, preflight)...)
}

func resourceAbpPreflightDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	// Preflights don't need to be deleted, forgetting it is fine
	data.SetId("")
	return nil
}
