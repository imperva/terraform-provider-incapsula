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

func resourceAbpProofOfWorkConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpProofOfWorkConfigurationCreate,
		ReadContext:   resourceAbpProofOfWorkConfigurationRead,
		UpdateContext: resourceAbpProofOfWorkConfigurationUpdate,
		DeleteContext: resourceAbpProofOfWorkConfigurationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpProofOfWorkConfigurationImport,
		},

		Description: `Provides an ABP Proof Of Work Configuration resource. A Proof Of Work
configuration tunes the cost of the proof-of-work challenge that clients must
complete to earn credits. It is referenced from policy directives whose action
is "proof_of_work".`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this Proof Of Work Configuration belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Human-readable name of the Proof Of Work Configuration.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},
			"difficulty": {
				Description: "Number of credits a client gets for completing a proof-of-work " +
					"challenge. The amount of work performed by the client is proportional to " +
					"this number. Must be >= 1.",
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"algorithm": {
				Description: "Proof Of Work algorithm. One of: bbs, sha1. Optional: when " +
					"omitted, the backend selects a default and that value is reflected " +
					"in state.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"bbs", "sha1"}, false),
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Proof Of Work Configuration was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Proof Of Work Configuration was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func extractAbpProofOfWorkConfiguration(data *schema.ResourceData) AbpProofOfWorkConfiguration {
	return AbpProofOfWorkConfiguration{
		Name:       data.Get("name").(string),
		Difficulty: int64(data.Get("difficulty").(int)),
		Algorithm:  data.Get("algorithm").(string),
	}
}

func serializeAbpProofOfWorkConfiguration(data *schema.ResourceData, config *AbpProofOfWorkConfiguration) error {
	if err := data.Set("account_id", config.AccountId); err != nil {
		return err
	}
	if err := data.Set("name", config.Name); err != nil {
		return err
	}
	if err := data.Set("difficulty", config.Difficulty); err != nil {
		return err
	}
	if err := data.Set("algorithm", config.Algorithm); err != nil {
		return err
	}
	if err := data.Set("created_at", config.CreatedAt); err != nil {
		return err
	}
	if err := data.Set("modified_at", config.ModifiedAt); err != nil {
		return err
	}
	return nil
}

func resourceAbpProofOfWorkConfigurationCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	created, err := client.CreateAbpProofOfWorkConfiguration(accountId, extractAbpProofOfWorkConfiguration(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Proof Of Work Configuration create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpProofOfWorkConfiguration(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Proof Of Work Configuration %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpProofOfWorkConfigurationRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	config, err := client.ReadAbpProofOfWorkConfiguration(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if config == nil {
		log.Printf("[INFO] ABP Proof Of Work Configuration %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpProofOfWorkConfiguration(data, config); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpProofOfWorkConfigurationUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	updated, err := client.UpdateAbpProofOfWorkConfiguration(id, extractAbpProofOfWorkConfiguration(data))
	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpProofOfWorkConfiguration(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpProofOfWorkConfigurationDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpProofOfWorkConfiguration(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpProofOfWorkConfigurationImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<proof_of_work_configuration_id>'")
	}

	client := m.(*Client)
	config, err := client.ReadAbpProofOfWorkConfiguration(id)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("ABP Proof Of Work Configuration %s not found", id)
	}

	data.SetId(id)
	if err := data.Set("account_id", config.AccountId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}
