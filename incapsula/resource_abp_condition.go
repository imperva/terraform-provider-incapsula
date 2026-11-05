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

func resourceAbpCondition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpConditionCreate,
		ReadContext:   resourceAbpConditionRead,
		UpdateContext: resourceAbpConditionUpdate,
		DeleteContext: resourceAbpConditionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpConditionImport,
		},

		Description: `Provides an ABP Condition resource. A Condition contains MOI expression
evaluated against incoming requests. Conditions are referenced from Policies
(directly or via Condition Lists) to drive Directive actions.

NOTE: the API stores and returns a formatted and optimized form of the Condition
code, so the provider keeps the ` + "`code`" + ` you wrote and exposes the
server's form separately as ` + "`code_normalized`" + `. This has two
consequences. A code change made outside Terraform (for example in the UI) is
**not** reported by ` + "`terraform plan`" + `. And the ` + "`code`" + ` from
your configuration is sent on every update, so any later change to this
resource — even to ` + "`name`" + ` alone — silently replaces the out-of-band
code with the configured one.

It is therefore recommended to copy ` + "`code_normalized`" + ` into
` + "`code`" + ` after running ` + "`terraform import`" + `, and to make all code
changes through Terraform.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "ABP account UUID this Condition belongs to. Changing it replaces " +
					"the Condition — Terraform deletes the existing one and creates a new one — " +
					"because the API cannot move a Condition to another account.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Human-readable name of the Condition. 1..100 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description: "Description of the Condition. Required; pass an empty string for no description.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"code": {
				Description: "MOI expression evaluated against the request. The server stores a " +
					"normalized form internally; see `code_normalized` for the server's view. " +
					"Out-of-band edits to this field (via the UI or API) are not detected as drift, " +
					"and the next update of this resource — whatever it changes — overwrites them " +
					"with the value configured here.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"code_normalized": {
				Description: "Server-side normalized/optimized form of `code`. Useful for " +
					"diagnostics and for seeding `.tf` after `terraform import`.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"template": {
				Description: "Editor template the server uses to validate `code`. Open enumeration; " +
					"new values may be added, so the provider does not restrict the set.\n\n" +
					"Currently available values:\n" +
					"  - `custom` (default)\n" +
					"  - `ai_tool`\n" +
					"  - `ip_set`\n" +
					"  - `header`\n" +
					"  - `tag`\n" +
					"  - `flag`\n" +
					"  - `challenge`\n" +
					"  - `captcha_cleared`\n" +
					"  - `proof_of_work_succeeded`\n" +
					"  - `rate_limiting`\n" +
					"  - `compound_rate_limiting`\n" +
					"  - `identify_eventually`\n\n" +
					"Omitting `template` is the same as setting `custom`.\n\n" +
					"**WARNING:** for any value other than `custom`, the server requires `code` to " +
					"match that template's expected structure. Mismatched `code` makes Condition " +
					"create/update fail.\n\n" +
					"**WARNING:** the template is fixed when the Condition is created. The API " +
					"ignores a new value on update, so changing it appears to succeed while nothing " +
					"changes and the next plan shows the same change again. To change the template, " +
					"recreate the Condition, for example with " +
					"`terraform apply -replace=incapsula_abp_condition.<name>`.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"last_change_by": {
				Description: "Identifier of the user who last changed this Condition.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Condition was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Condition was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func extractAbpCondition(data *schema.ResourceData) AbpCondition {
	return AbpCondition{
		Kind:        AbpConditionKindLiteral,
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
		Code:        data.Get("code").(string),
		Template:    data.Get("template").(string),
	}
}

func serializeAbpCondition(data *schema.ResourceData, condition *AbpCondition) error {
	if condition.Kind != AbpConditionKindLiteral {
		return fmt.Errorf("ABP Condition %s is not a literal variant (it is a %s)", condition.Id, condition.Kind)
	}
	if err := data.Set("name", condition.Name); err != nil {
		return err
	}

	if err := data.Set("description", condition.Description); err != nil {
		return err
	}

	// Intentionally do not overwrite "code": the server stores a normalized form
	// that would otherwise show as perpetual drift. Expose it via "code_normalized".
	if err := data.Set("code_normalized", condition.Code); err != nil {
		return err
	}

	// The server defaults an unset template to "custom" and echoes it back; the
	// field is Computed so this resolves the empty-vs-"custom" diff for the user.
	if err := data.Set("template", condition.Template); err != nil {
		return err
	}

	if condition.AccountId == "" {
		return fmt.Errorf("Managed conditions are not supported: account_id of condition %s is empty", condition.Id)
	}

	if err := data.Set("account_id", condition.AccountId); err != nil {
		return err
	}

	if err := data.Set("last_change_by", condition.LastChangeBy); err != nil {
		return err
	}

	if err := data.Set("created_at", condition.CreatedAt); err != nil {
		return err
	}

	if err := data.Set("modified_at", condition.ModifiedAt); err != nil {
		return err
	}

	return nil
}

func resourceAbpConditionCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	created, err := client.CreateAbpCondition(accountId, extractAbpCondition(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Condition create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpCondition(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Condition %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpConditionRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	condition, err := client.ReadAbpCondition(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if condition == nil {
		log.Printf("[INFO] ABP Condition %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpCondition(data, condition); err != nil {
		return diag.FromErr(err)
	}

	// On import there is no prior state for "code"; seed it from the server's
	// normalized value so the user has something to copy into their .tf.
	if data.Get("code").(string) == "" {
		if err := data.Set("code", condition.Code); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceAbpConditionUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	updated, err := client.UpdateAbpCondition(id, extractAbpCondition(data))
	if err != nil {
		return diag.FromErr(err)
	}

	// Condition not found
	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpCondition(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpConditionDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpCondition(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpConditionImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<condition_id>'")
	}

	client := m.(*Client)
	condition, err := client.ReadAbpCondition(id)
	if err != nil {
		return nil, err
	}
	if condition == nil {
		return nil, fmt.Errorf("ABP Condition %s not found", id)
	}
	if condition.AccountId == "" {
		return nil, fmt.Errorf("ABP Condition %s is a managed condition and cannot be imported; only account-owned conditions are supported", id)
	}

	data.SetId(id)
	if err := data.Set("account_id", condition.AccountId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}
