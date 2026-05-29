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

func resourceAbpConditionList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpConditionListCreate,
		ReadContext:   resourceAbpConditionListRead,
		UpdateContext: resourceAbpConditionListUpdate,
		DeleteContext: resourceAbpConditionListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpConditionListImport,
		},

		Description: `Provides an ABP Condition List resource. A condition list is a named
container that groups conditions; entries are managed via
` + "`incapsula_abp_condition_list_entry`" + `. Condition lists can be referenced
from policies and from other condition lists.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this Condition List belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Human-readable name of the condition list. 1..100 characters.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description: "Description of the condition list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Condition List was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Condition List was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func extractAbpConditionList(data *schema.ResourceData) AbpConditionList {
	return AbpConditionList{
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}
}

func serializeAbpConditionList(data *schema.ResourceData, list *AbpConditionList) error {
	if list.AccountId == "" {
		return fmt.Errorf("Managed condition lists are not supported: account_id of condition list %s is empty", list.Id)
	}
	if err := data.Set("account_id", list.AccountId); err != nil {
		return err
	}
	if err := data.Set("name", list.Name); err != nil {
		return err
	}
	if err := data.Set("description", list.Description); err != nil {
		return err
	}
	if err := data.Set("created_at", list.CreatedAt); err != nil {
		return err
	}
	if err := data.Set("modified_at", list.ModifiedAt); err != nil {
		return err
	}
	return nil
}

func resourceAbpConditionListCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	created, err := client.CreateAbpConditionList(accountId, extractAbpConditionList(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Condition List create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpConditionList(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Condition List %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpConditionListRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	list, err := client.ReadAbpConditionList(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if list == nil {
		log.Printf("[INFO] ABP Condition List %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpConditionList(data, list); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpConditionListUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	updated, err := client.UpdateAbpConditionList(id, extractAbpConditionList(data))
	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpConditionList(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpConditionListDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpConditionList(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpConditionListImport(ctx context.Context, data *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<condition_list_id>'")
	}

	client := m.(*Client)
	list, err := client.ReadAbpConditionList(id)
	if err != nil {
		return nil, err
	}
	if list == nil {
		return nil, fmt.Errorf("ABP Condition List %s not found", id)
	}
	if list.AccountId == "" {
		return nil, fmt.Errorf("ABP Condition List %s is a managed condition list and cannot be imported; only account-owned condition lists are supported", id)
	}

	data.SetId(id)
	if err := data.Set("account_id", list.AccountId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}
