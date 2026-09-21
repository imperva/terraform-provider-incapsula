package incapsula

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAiApplicationSecurityApiKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAiApplicationSecurityApiKeyCreate,
		ReadContext:   resourceAiApplicationSecurityApiKeyRead,
		DeleteContext: resourceAiApplicationSecurityApiKeyDelete,
		// No UpdateContext: only name is settable and the backend has no update endpoint,
		// so every input is ForceNew.
		Importer: &schema.ResourceImporter{
			StateContext: resourceAiApplicationSecurityApiKeyImport,
		},

		Description: "Manages an AI Application Security API key for an application. The plaintext key is " +
			"returned only once, on create, in the sensitive api_key attribute and cannot be " +
			"recovered afterward (it is empty after import). All inputs are immutable.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "The Imperva account ID that owns the application. Defaults to the account of the API credentials when omitted.",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			"application_id": {
				Type:         schema.TypeString,
				Description:  "The UUID of the AI Application Security application this API key belongs to.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The API key name (1-50 chars, alphanumeric, underscore or hyphen).",
				Required:    true,
				ForceNew:    true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 50),
					validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_-]+$`), "must be alphanumeric, underscore or hyphen"),
				),
			},
			"api_key": {
				Type:        schema.TypeString,
				Description: "The plaintext API key. Returned only once on creation and unrecoverable afterward; empty after import.",
				Computed:    true,
				Sensitive:   true,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the API key is active.",
				Computed:    true,
			},
		},
	}
}

func resourceAiApplicationSecurityApiKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Get("application_id").(string)

	req := &AiApplicationSecurityApiKeyRequest{
		Name: d.Get("name").(string),
	}

	created, err := client.CreateAiApplicationSecurityApiKey(accountID, applicationID, req)
	if err != nil {
		return diag.Errorf("Failed to create AI Application Security API key for application %s: %s", applicationID, err)
	}

	// The plaintext key is available only in this response; persist it now — Read never touches it.
	d.Set("api_key", created.FullKey)
	d.SetId(strconv.FormatInt(created.ApiKey.Id, 10))
	log.Printf("[INFO] Created AI Application Security API key with ID: %d", created.ApiKey.Id)

	return resourceAiApplicationSecurityApiKeyRead(ctx, d, m)
}

func resourceAiApplicationSecurityApiKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)

	apiKeyID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Invalid AI Application Security API key ID %q: %s", d.Id(), err)
	}

	apiKey, err := client.GetAiApplicationSecurityApiKey(accountID, apiKeyID)
	if err != nil {
		return diag.Errorf("Failed to read AI Application Security API key %d: %s", apiKeyID, err)
	}
	if apiKey == nil {
		log.Printf("[INFO] AI Application Security API key %d not found, removing from state", apiKeyID)
		d.SetId("")
		return nil
	}

	// Note: api_key (plaintext) is intentionally NOT set here — it is unavailable on read
	// and must retain its create-time value (empty after import).
	if apiKey.AccountId != 0 {
		d.Set("account_id", int(apiKey.AccountId))
	}
	// application_id is Required + ForceNew. The account-level list DTO tags applicationId
	// omitempty, so guard against an empty value overwriting the configured UUID — otherwise a
	// missing field would trigger a destroy/recreate whose Delete targets an empty path.
	if apiKey.ApplicationId != "" {
		d.Set("application_id", apiKey.ApplicationId)
	}
	d.Set("name", apiKey.Name)
	d.Set("active", apiKey.Active)

	return nil
}

func resourceAiApplicationSecurityApiKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Get("application_id").(string)

	apiKeyID, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.Errorf("Invalid AI Application Security API key ID %q: %s", d.Id(), err)
	}

	if err := client.DeleteAiApplicationSecurityApiKey(accountID, applicationID, apiKeyID); err != nil {
		return diag.Errorf("Failed to delete AI Application Security API key %d: %s", apiKeyID, err)
	}

	d.SetId("")
	return nil
}

func resourceAiApplicationSecurityApiKeyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	raw := strings.TrimSpace(d.Id())
	if raw == "" {
		return nil, fmt.Errorf("expected import ID to be '<api_key_id>' or '<account_id>/<api_key_id>'")
	}

	var accountID string
	resourceID := raw

	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid import ID %q: want '<api_key_id>' or '<account_id>/<api_key_id>'", raw)
		}
		accountID = strings.TrimSpace(parts[0])
		resourceID = strings.TrimSpace(parts[1])
	}

	// Set the canonical Terraform ID to just the resource ID.
	d.SetId(resourceID)

	// Only set account_id if it was provided.
	if accountID != "" {
		caid, err := strconv.Atoi(accountID)
		if err != nil {
			return nil, fmt.Errorf("invalid account_id %q: must be numeric", accountID)
		}
		if err := d.Set("account_id", caid); err != nil {
			return nil, fmt.Errorf("setting account_id: %w", err)
		}
	}

	return []*schema.ResourceData{d}, nil
}
