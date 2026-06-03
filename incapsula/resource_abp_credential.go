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
		DeleteContext: resourceAbpCredentialDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpCredentialImport,
		},

		Description: `Provides an ABP Account credential used to authenticate with the Analysis
Host. The secret is generated server-side on creation and cannot be retrieved
again later, so it is exposed at creation time only.

If ` + "`rsa_key`" + ` is supplied, the secret is encrypted to that key using OAEP with SHA256 and the label 'abp_credential', and stored
only as ` + "`encrypted_secret`" + `; the plaintext ` + "`secret`" + ` is never written to
state. This is strongly recommended, since ` + "`Sensitive`" + ` attributes are
redacted from CLI output but are still stored in plaintext in the state file.
Decrypt locally with, for example:
` + "`terraform output -raw encrypted_secret | base64 -d | openssl pkeyutl -decrypt -inkey <your-private-key-file> -pkeyopt rsa_padding_mode:oaep -pkeyopt rsa_oaep_md:sha256 -pkeyopt rsa_oaep_label:$(echo -n 'abp_credential' | xxd -p)`" + `.

NOTE: credentials cannot be modified in place. Changing ` + "`account_id`" + ` or
` + "`rsa_key`" + ` forces resource replacement. The Account configuration must be
published for the analysis host to begin accepting a new credential, and you
should keep the previous credential alive until that publish completes (see the
API documentation for the recommended rotation procedure).`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "Account UUID this credential belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"rsa_key": {
				Description: "A full PEM-encoded RSA public key to use for encrypting the secret in state. If set, the secret is encrypted to this key and stored only in `encrypted_secret`, keeping the plaintext out of state.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
			},
			"secret": {
				Description: "Base64-encoded credential secret. Generated server-side at creation time and never returned again afterwards. Empty when `rsa_key` is set.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
			"encrypted_secret": {
				Description: "The RSA-encrypted, base64-encoded credential secret. Only populated when `rsa_key` is set.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the credential was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the credential was last modified. Empty if it has never been modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

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

	if rsaKey, ok := data.Get("rsa_key").(string); ok && rsaKey != "" {
		encrypted, err := encryptRsa([]byte(rsaKey), []byte(created.Secret), "abp_credential")
		if err != nil {
			return diag.FromErr(err)
		}
		if err := data.Set("encrypted_secret", encrypted); err != nil {
			return diag.FromErr(err)
		}
		// Keep the plaintext secret out of state when a rsa_key is provided.
		if err := data.Set("secret", ""); err != nil {
			return diag.FromErr(err)
		}
	} else if err := data.Set("secret", created.Secret); err != nil {
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

func resourceAbpCredentialDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	if err := client.DeleteAbpCredential(data.Id()); err != nil {
		return diag.FromErr(err)
	}
	data.SetId("")
	return nil
}

// resourceAbpCredentialImport imports an existing credential by its ID. The
// secret is not returned by the read endpoint and will remain empty in state.
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
