package incapsula

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/mail"
	"strconv"
	"strings"
)

func resourceApiClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApiClientCreate,
		ReadContext:   resourceApiClientRead,
		UpdateContext: resourceApiClientUpdate,
		DeleteContext: resourceApiClientDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceApiClientImport,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account ID",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"user_email": {
				Description: "Email address of the user that the api client belongs to",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					email := val.(string)
					if _, err := mail.ParseAddress(email); err != nil {
						errs = append(errs, fmt.Errorf("%q is invalid, got: %s", key, email))
					}
					return
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the API client.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Description of the API client.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the API client is enabled.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated API key for client authentication.",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Expiration date of the API key (YYYY-MM-DD format only). Must be a future date. Changing this value will cause regeneration of the key.",
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last used timestamp.",
			},
		},
	}
}

func resourceApiClientCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	accountID := d.Get("account_id").(int)

	apiClient := APIClientUpdateRequest{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		ExpirationDate: d.Get("expiration_date").(string),
		Enabled:        Bool(d.Get("enabled").(bool)),
	}

	resp, err := client.CreateAPIClient(
		accountID,
		d.Get("user_email").(string),
		&apiClient,
	)

	if err != nil {
		return diag.Errorf("[ERROR] Could not create API client, %s\n", err)
	}

	// Set the User ID
	d.SetId(strconv.Itoa(resp.APIClientID))
	log.Printf("[INFO] Created Incapsula API client with ID: %d ", resp.APIClientID)

	// Set the account_id in state
	if accountID != 0 {
		d.Set("account_id", accountID)
	}

	if resp.APIKey != "" {
		d.Set("api_key", resp.APIKey)
	}

	return resourceApiClientRead(ctx, d, m)
}

func resourceApiClientUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	id := d.Id()
	accountID := d.Get("account_id").(int)

	request := &APIClientUpdateRequest{}
	changed := false

	if d.HasChange("expiration_date") {
		request.ExpirationDate = d.Get("expiration_date").(string)
		request.Regenerate = true
		changed = true
	}
	if d.HasChange("enabled") {
		request.Enabled = Bool(d.Get("enabled").(bool))
		changed = true
	}

	if d.HasChange("name") {
		request.Name = d.Get("name").(string)
		changed = true
	}

	if d.HasChange("description") {
		request.Description = d.Get("description").(string)
		changed = true
	}

	if !changed {
		// Set the account_id in state even if no changes
		if accountID != 0 {
			d.Set("account_id", accountID)
		}
		return resourceApiClientRead(ctx, d, meta)
	}

	resp, err := client.PatchAPIClient(accountID, id, request)
	if err != nil {
		return diag.Errorf("error updating API client: %v", err)
	}

	// Set the account_id in state
	if accountID != 0 {
		d.Set("account_id", accountID)
	}

	if resp.APIKey != "" {
		d.Set("api_key", resp.APIKey)
	}

	return resourceApiClientRead(ctx, d, meta)
}

// After: Implemented resourceApiClientRead function
func resourceApiClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	id := d.Id()
	if id == "" {
		return nil
	}
	accountID := d.Get("account_id").(int)
	resp, err := client.GetAPIClient(accountID, id)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Set the account_id in state
	if accountID != 0 {
		d.Set("account_id", accountID)
	}

	d.Set("name", resp.Name)
	d.Set("description", resp.Description)
	d.Set("enabled", resp.Enabled)
	d.Set("expiration_date", resp.ExpirationDate)
	d.Set("last_used_at", resp.LastUsedAt)
	d.Set("user_email", resp.UserEmail)
	return nil
}

// DELETE /v3/api-client/{client_id}
func resourceApiClientDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)
	id := d.Id()
	if err := client.DeleteAPIClient(d.Get("account_id").(int), id); err != nil {
		return diag.Errorf("error deleting API client: %v", err)
	}
	d.SetId("")
	return nil
}

// Supports "<resource_id>" OR "<account_id>/<resource_id>".
// If account_id is omitted, we do NOT set the "account_id" attribute.
func resourceApiClientImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	raw := strings.TrimSpace(d.Id())
	if raw == "" {
		return nil, fmt.Errorf("expected import ID to be '<api_client_id>' or '<account_id>/<api_client_id>'")
	}

	var accountID string
	resourceID := raw

	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid import ID %q: want '<api_client_id>' or '<account_id>/<api_client_id>'", raw)
		}
		accountID = strings.TrimSpace(parts[0])
		resourceID = strings.TrimSpace(parts[1])
	}

	// Set the canonical Terraform ID to just the resource ID.
	d.SetId(resourceID)

	// Only set account_id if it was provided.
	if accountID != "" {
		caid, _ := strconv.Atoi(accountID)
		if err := d.Set("account_id", caid); err != nil {
			return nil, fmt.Errorf("setting account_id: %w", err)
		}
	}

	return []*schema.ResourceData{d}, nil
}
