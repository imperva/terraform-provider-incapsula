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

func resourceAbpCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpCredentialCreate,
		ReadContext:   resourceAbpCredentialRead,
		UpdateContext: resourceAbpCredentialUpdate,
		DeleteContext: resourceAbpCredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpCredentialImport,
		},
		CustomizeDiff: resourceAbpCredentialCustomizeDiff,

		Description: `Provides an ABP Account credential used to authenticate with the Analysis
Host. The secret is generated server-side on creation.
The secret is encrypted by ` + "`rsa_key`" + ` using OAEP with SHA256 and the label 'abp_credential', and stored
only as ` + "`encrypted_secret`" + `; the plaintext ` + "`secret`" + ` is never written to
state.
Decrypt locally with, for example:
` + "`terraform output -raw encrypted_secret | base64 -d | openssl pkeyutl -decrypt -inkey <your-private-key-file> -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 -pkeyopt rsa_oaep_label:$(echo -n 'abp_credential' | xxd -p)`" + `.

` + "`rsa_key`" + ` only decides which key the secret in state is encrypted to, so
changing it is an in-place re-encryption of the same secret: the credential on
the server is left untouched. For the same reason ` + "`terraform import`" + ` is
safe — the imported credential simply has no ` + "`encrypted_secret`" + ` until the
first apply, which fills it in place without replacing anything.

NOTE: the secret itself cannot be rotated in place. Changing
` + "`account_id`" + ` replaces the credential, which generates a new secret. The
Account configuration must be published for the analysis host to begin accepting
a new credential, and you should keep the previous credential alive until that
publish completes (see the API documentation for the recommended rotation
procedure).`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Account UUID this credential belongs to. Changing it replaces the " +
					"credential — Terraform deletes the existing one and creates a new one with a " +
					"new secret — because the API cannot move a credential to another account.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"rsa_key": {
				Description: "A full PEM-encoded RSA public key used for encrypting the secret in state. Must be a `BEGIN PUBLIC KEY` (SPKI) PEM block of at least 3072 bits; other formats and smaller keys are rejected. The secret is encrypted to this key and stored only in `encrypted_secret`, keeping the plaintext out of state. Changing it re-encrypts the same secret in place.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"encrypted_secret": {
				Description: "The RSA-encrypted, base64-encoded credential secret",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the credential was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the credential was last modified. Always empty: the API does not track modifications of a credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

const abpCredentialEncryptionLabel = "abp_credential"

func serializeAbpCredential(data *schema.ResourceData, c *AbpCredential) error {
	if err := data.Set("account_id", c.AccountId); err != nil {
		return err
	}
	if err := data.Set("created_at", c.CreatedAt); err != nil {
		return err
	}
	modified := ""
	if c.ModifiedAt != nil {
		modified = *c.ModifiedAt
	}
	if err := data.Set("modified_at", modified); err != nil {
		return err
	}
	return nil
}

func resourceAbpCredentialCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	created, err := client.CreateAbpCredential(accountId)
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Credential create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpCredential(data, created); err != nil {
		return diag.FromErr(err)
	}
	rsaKey := data.Get("rsa_key").(string)
	encrypted, err := encryptRsa([]byte(rsaKey), []byte(created.Secret), abpCredentialEncryptionLabel)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("encrypted_secret", encrypted); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Created ABP Credential %s on account %s", created.Id, accountId)
	return nil
}

func resourceAbpCredentialRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	cred, err := client.ReadAbpCredential(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if cred == nil {
		log.Printf("[INFO] ABP Credential %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpCredential(data, cred); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// resourceAbpCredentialUpdate re-encrypts the existing secret. rsa_key is the
// only updatable attribute and is a purely local concern — it decides which key
// the secret in state is encrypted to — so the credential on the server is read
// but never modified.
func resourceAbpCredentialUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	cred, err := client.ReadAbpCredential(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if cred == nil {
		return diag.Errorf("ABP Credential %s no longer exists, so its secret cannot be re-encrypted. Run 'terraform apply -refresh-only' to drop it from state and create a replacement.", id)
	}
	if cred.Secret == "" {
		return diag.Errorf("ABP Credential %s read response did not contain a secret", id)
	}

	if err := serializeAbpCredential(data, cred); err != nil {
		return diag.FromErr(err)
	}

	rsaKey := data.Get("rsa_key").(string)
	encrypted, err := encryptRsa([]byte(rsaKey), []byte(cred.Secret), abpCredentialEncryptionLabel)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := data.Set("encrypted_secret", encrypted); err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[INFO] Re-encrypted the secret of ABP Credential %s", id)
	return nil
}

// resourceAbpCredentialCustomizeDiff marks encrypted_secret as unknown whenever
// rsa_key changes, so that the re-encryption is visible in the plan.
func resourceAbpCredentialCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, m any) error {
	if diff.Id() != "" && diff.HasChange("rsa_key") {
		return diff.SetNewComputed("encrypted_secret")
	}
	return nil
}

func resourceAbpCredentialDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	if err := client.DeleteAbpCredential(data.Id()); err != nil {
		return diag.FromErr(err)
	}
	data.SetId("")
	return nil
}

// resourceAbpCredentialImport imports an existing credential by its ID.
// rsa_key is not part of the remote object, so it is left empty; the first
// apply after the import re-encrypts the secret in place and fills
// encrypted_secret without replacing the credential.
func resourceAbpCredentialImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<credential_id>'")
	}

	client := m.(*Client)
	cred, err := client.ReadAbpCredential(id)
	if err != nil {
		return nil, err
	}
	if cred == nil {
		return nil, fmt.Errorf("ABP Credential %s not found", id)
	}

	data.SetId(id)
	if err := serializeAbpCredential(data, cred); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
