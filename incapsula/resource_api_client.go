package incapsula

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApiClient() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiClientCreate,
		Read:   resourceApiClientRead,
		Update: resourceApiClientUpdate,
		Delete: resourceApiClientDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the API client.",
			},
			"expiration_period": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration period for the API key (RFC3339 or duration, e.g. 30d).",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the API client is enabled.",
			},
			"grace_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Grace period in seconds.",
			},
			"regenerate_version": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Increment to trigger API key regeneration.",
			},
			"api_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Generated API key for client authentication.",
			},
			"api_client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the API client.",
			},
			"expiration_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Expiration date of the API key.",
			},
			"last_used_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last used timestamp.",
			},
		},
	}
}

// PATCH /v3/api-client/{client_id} for create, update, enable/disable, regenerate
func resourceApiClientCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	ctx := context.Background()

	name := d.Get("name").(string)
	expirationPeriod, expirationPeriodSet := d.GetOk("expiration_period")
	enabled, enabledSet := d.GetOk("enabled")
	gracePeriod, gracePeriodSet := d.GetOk("grace_period")
	regenerate := true // Always regenerate on create

	request := &APIClientUpdateRequest{
		Regenerate: &regenerate,
	}
	if expirationPeriodSet {
		str := expirationPeriod.(string)
		request.ExpirationPeriod = &str
	}
	if enabledSet {
		b := enabled.(bool)
		request.Enabled = &b
	}
	if gracePeriodSet {
		gp := gracePeriod.(int)
		request.GracePeriod = &gp
	}

	resp, err := client.PatchAPIClient(ctx, name, request)
	if err != nil {
		return fmt.Errorf("error creating API client: %w", err)
	}
	d.SetId(resp.APIClientID)
	if err := d.Set("api_key", resp.APIKey); err != nil {
		return err
	}
	if err := d.Set("api_client_id", resp.APIClientID); err != nil {
		return err
	}
	if err := d.Set("enabled", resp.Enabled); err != nil {
		return err
	}
	if err := d.Set("expiration_date", resp.ExpirationDate); err != nil {
		return err
	}
	if err := d.Set("last_used_at", resp.LastUsedAt); err != nil {
		return err
	}
	if err := d.Set("grace_period", resp.GracePeriod); err != nil {
		return err
	}
	if err := d.Set("name", name); err != nil {
		return err
	}
	return nil
}

func resourceApiClientUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	ctx := context.Background()
	id := d.Id()

	request := &APIClientUpdateRequest{}
	changed := false

	if d.HasChange("expiration_period") {
		str := d.Get("expiration_period").(string)
		request.ExpirationPeriod = &str
		changed = true
	}
	if d.HasChange("enabled") {
		b := d.Get("enabled").(bool)
		request.Enabled = &b
		changed = true
	}
	if d.HasChange("grace_period") {
		gp := d.Get("grace_period").(int)
		request.GracePeriod = &gp
		changed = true
	}
	if d.HasChange("regenerate_version") {
		regenerate := true
		request.Regenerate = &regenerate
		changed = true
	}

	if !changed {
		return resourceApiClientRead(d, meta)
	}

	resp, err := client.PatchAPIClient(ctx, id, request)
	if err != nil {
		return fmt.Errorf("error updating API client: %w", err)
	}
	if err := d.Set("api_key", resp.APIKey); err != nil {
		return err
	}
	if err := d.Set("api_client_id", resp.APIClientID); err != nil {
		return err
	}
	if err := d.Set("enabled", resp.Enabled); err != nil {
		return err
	}
	if err := d.Set("expiration_date", resp.ExpirationDate); err != nil {
		return err
	}
	if err := d.Set("last_used_at", resp.LastUsedAt); err != nil {
		return err
	}
	if err := d.Set("grace_period", resp.GracePeriod); err != nil {
		return err
	}
	if err := d.Set("name", d.Get("name")); err != nil {
		return err
	}
	return nil
}

// After: Implemented resourceApiClientRead function
func resourceApiClientRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	ctx := context.Background()
	id := d.Id()
	if id == "" {
		return nil
	}
	resp, err := client.GetAPIClient(ctx, id)
	if err != nil {
		d.SetId("")
		return nil
	}
	if err := d.Set("name", d.Get("name")); err != nil {
		return err
	}
	if err := d.Set("api_key", resp.APIKey); err != nil {
		return err
	}
	if err := d.Set("api_client_id", resp.APIClientID); err != nil {
		return err
	}
	if err := d.Set("enabled", resp.Enabled); err != nil {
		return err
	}
	if err := d.Set("expiration_date", resp.ExpirationDate); err != nil {
		return err
	}
	if err := d.Set("last_used_at", resp.LastUsedAt); err != nil {
		return err
	}
	if err := d.Set("grace_period", resp.GracePeriod); err != nil {
		return err
	}
	return nil
}

// DELETE /v3/api-client/{client_id}
func resourceApiClientDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	ctx := context.Background()
	id := d.Id()
	if err := client.DeleteAPIClient(ctx, id); err != nil {
		return fmt.Errorf("error deleting API client: %w", err)
	}
	d.SetId("")
	return nil
}
