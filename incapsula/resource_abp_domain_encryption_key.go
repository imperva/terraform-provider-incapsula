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

func resourceAbpDomainEncryptionKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpDomainEncryptionKeyCreate,
		ReadContext:   resourceAbpDomainEncryptionKeyRead,
		DeleteContext: resourceAbpDomainEncryptionKeyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpDomainEncryptionKeyImport,
		},

		Description: `Provides an ABP Domain encryption key. The key material is base64-encoded
between 24 and 64 raw bytes (32-88 base64 characters).

NOTE: encryption keys cannot be modified in place. Changing ` + "`key`" + ` or
` + "`domain_id`" + ` replaces the resource, and by default Terraform deletes the
old key *before* creating the new one — which is exactly what you must avoid,
because the Account configuration has to be published before the analysis host
accepts a new key. Add ` + "`lifecycle { create_before_destroy = true }`" + ` to
the resource so the new key is created first and the previous one stays alive
until the publish completes (see the API documentation for the recommended
rotation procedure).

NOTE: the API returns key material in plaintext, so the ` + "`key`" + ` value is
stored in Terraform state. Protect the state as you would the key itself.`,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "Domain UUID this key belongs to. Changing it replaces the key: " +
					"Terraform deletes this key and creates a new one on the new Domain.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"key": {
				Description: "Base64-encoded key material (32..88 characters). Must be valid " +
					"base64; an invalid value is only rejected when the change is applied. " +
					"Changing it replaces the key (see the rotation note above).",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringLenBetween(32, 88),
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the key was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"first_published_at": {
				Description: "RFC3339 timestamp at which the key was first published to the analysis host. Empty until the first publish completes.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func serializeAbpDomainEncryptionKey(data *schema.ResourceData, k *AbpDomainEncryptionKey) error {
	if err := data.Set("domain_id", k.DomainId); err != nil {
		return err
	}
	if err := data.Set("key", k.Key); err != nil {
		return err
	}
	if err := data.Set("created_at", k.CreatedAt); err != nil {
		return err
	}
	first := ""
	if k.FirstPublishedAt != nil {
		first = *k.FirstPublishedAt
	}
	if err := data.Set("first_published_at", first); err != nil {
		return err
	}
	return nil
}

func resourceAbpDomainEncryptionKeyCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	domainId := data.Get("domain_id").(string)

	created, err := client.CreateAbpDomainEncryptionKey(domainId, AbpDomainEncryptionKey{
		Key: data.Get("key").(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Domain encryption key create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpDomainEncryptionKey(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Domain encryption key %s on domain %s", created.Id, domainId)
	return nil
}

func resourceAbpDomainEncryptionKeyRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	domainId := data.Get("domain_id").(string)
	id := data.Id()

	key, err := client.ReadAbpDomainEncryptionKey(domainId, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if key == nil {
		log.Printf("[INFO] ABP Domain encryption key %s on domain %s not found, removing from state", id, domainId)
		data.SetId("")
		return nil
	}

	if err := serializeAbpDomainEncryptionKey(data, key); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpDomainEncryptionKeyDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	if err := client.DeleteAbpDomainEncryptionKey(data.Id()); err != nil {
		return diag.FromErr(err)
	}
	data.SetId("")
	return nil
}

// resourceAbpDomainEncryptionKeyImport accepts "<domain_id>/<key_id>" because
// the read endpoint requires the parent domain id, which would otherwise be
// missing from state on import.
func resourceAbpDomainEncryptionKeyImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(strings.TrimSpace(data.Id()), "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("expected import ID to be '<domain_id>/<key_id>'")
	}
	domainId, keyId := parts[0], parts[1]

	client := m.(*Client)
	key, err := client.ReadAbpDomainEncryptionKey(domainId, keyId)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, fmt.Errorf("ABP Domain encryption key %s on domain %s not found", keyId, domainId)
	}

	data.SetId(keyId)
	if err := data.Set("domain_id", domainId); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
