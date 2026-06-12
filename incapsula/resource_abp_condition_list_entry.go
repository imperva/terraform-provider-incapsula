package incapsula

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const abpConditionListEntryResourceName = "ABP Condition List Entry"

var abpConditionListEntryTagRegexp = regexp.MustCompile(`^[_a-z][_a-z0-9]*$`)

const (
	abpConditionStateInactive = "inactive"
	abpConditionStateMonitor  = "monitor"
	abpConditionStateActive   = "active"
)

func resourceAbpConditionListEntry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpConditionListEntryCreate,
		ReadContext:   resourceAbpConditionListEntryRead,
		UpdateContext: resourceAbpConditionListEntryUpdate,
		DeleteContext: resourceAbpConditionListEntryDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpConditionListEntryImport,
		},

		Description: `Provides an ABP Condition List Entry resource. A condition list entry
attaches a literal condition (` + "`condition_id`" + `) or a nested condition list
(` + "`condition_list_id`" + `) to a parent condition list, with tags and an active
state. The lifecycle of the entry is independent of the referenced condition,
which is owned by its own resource.

This is also how conditions are added to a policy: a directive exposes a
` + "`condition_list_id`" + ` (and, for ` + "`proof_of_work`" + ` directives, a
` + "`skip_condition_list_id`" + `), which can be used as the ` + "`parent_condition_list_id`" + `
of an entry to attach a condition to that directive. These IDs are available on
the ` + "`incapsula_abp_policy`" + ` resource and the ` + "`incapsula_abp_directive`" + ` data source.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description:  "ABP account UUID this Condition List Entry belongs to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"parent_condition_list_id": {
				Description:  "UUID of the parent condition list that this entry is attached to.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"condition_id": {
				Description:  "UUID of a literal condition referenced by this entry. Mutually exclusive with `condition_list_id`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"condition_id", "condition_list_id"},
			},
			"condition_list_id": {
				Description:  "UUID of a condition list referenced by this entry. Mutually exclusive with `condition_id`.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
				ExactlyOneOf: []string{"condition_id", "condition_list_id"},
			},
			"tags": {
				Description: "Snake_case tags applied to the entry. Each tag matches `^[_a-z][_a-z0-9]*$`.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						abpConditionListEntryTagRegexp,
						"tag must be snake_case matching ^[_a-z][_a-z0-9]*$",
					),
				},
			},
			"state": {
				Description: "Whether the entry is evaluated. One of: " +
					"`inactive` (skipped), `monitor` (evaluated, logged only), " +
					"`active` (evaluated and may trigger an action).",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					abpConditionStateInactive,
					abpConditionStateMonitor,
					abpConditionStateActive,
				}, false),
			},
			"created_at": {
				Description: "RFC3339 timestamp at which the Condition List Entry was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"modified_at": {
				Description: "RFC3339 timestamp at which the Condition List Entry was last modified.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func extractAbpConditionListEntry(data *schema.ResourceData) AbpCondition {
	reference := data.Get("condition_id").(string)
	if reference == "" {
		reference = data.Get("condition_list_id").(string)
	}

	tagsSet := data.Get("tags").(*schema.Set)
	tags := make([]string, 0, tagsSet.Len())
	for _, t := range tagsSet.List() {
		tags = append(tags, t.(string))
	}

	return AbpCondition{
		Kind:      AbpConditionKindReference,
		Reference: reference,
		Parent:    data.Get("parent_condition_list_id").(string),
		Tags:      tags,
		State:     data.Get("state").(string),
	}
}

// serializeAbpConditionListEntry writes the entry to state. The
// referencedKind argument determines whether the referenced condition is
// exposed as `condition_id` (a literal) or `condition_list_id` (a list);
// the unselected field is cleared.
func serializeAbpConditionListEntry(data *schema.ResourceData, ref *AbpCondition, referencedKind AbpConditionKind) error {
	if ref.Kind != AbpConditionKindReference {
		return fmt.Errorf("%s %s is not a reference variant (it is a %s)", abpConditionListEntryResourceName, ref.Id, ref.Kind)
	}
	if ref.AccountId == "" {
		return fmt.Errorf("Managed condition list entries are not supported: account_id of entry %s is empty", ref.Id)
	}
	if err := data.Set("account_id", ref.AccountId); err != nil {
		return err
	}
	if err := data.Set("parent_condition_list_id", ref.Parent); err != nil {
		return err
	}

	switch referencedKind {
	case AbpConditionKindLiteral:
		if err := data.Set("condition_id", ref.Reference); err != nil {
			return err
		}
		if err := data.Set("condition_list_id", ""); err != nil {
			return err
		}
	case AbpConditionKindList:
		if err := data.Set("condition_list_id", ref.Reference); err != nil {
			return err
		}
		if err := data.Set("condition_id", ""); err != nil {
			return err
		}
	default:
		return fmt.Errorf("referenced condition %s has unexpected kind %q; expected literal or list", ref.Reference, referencedKind)
	}

	tags := make([]any, 0, len(ref.Tags))
	for _, t := range ref.Tags {
		tags = append(tags, t)
	}
	if err := data.Set("tags", schema.NewSet(schema.HashString, tags)); err != nil {
		return err
	}

	if err := data.Set("state", ref.State); err != nil {
		return err
	}
	if err := data.Set("created_at", ref.CreatedAt); err != nil {
		return err
	}
	if err := data.Set("modified_at", ref.ModifiedAt); err != nil {
		return err
	}
	return nil
}

// validateReferencedConditionKind ensures `condition_id` points at a literal
// and `condition_list_id` points at a list. The backend does not enforce
// this, so we check before sending the request.
func validateReferencedConditionKind(client *Client, data *schema.ResourceData) (AbpConditionKind, diag.Diagnostics) {
	conditionId := data.Get("condition_id").(string)
	conditionListId := data.Get("condition_list_id").(string)

	var (
		referenced string
		expected   AbpConditionKind
		field      string
	)
	switch {
	case conditionId != "":
		referenced = conditionId
		expected = AbpConditionKindLiteral
		field = "condition_id"
	case conditionListId != "":
		referenced = conditionListId
		expected = AbpConditionKindList
		field = "condition_list_id"
	default:
		return "", diag.Errorf("exactly one of condition_id or condition_list_id must be set")
	}

	condition, err := client.ReadAbpCondition(referenced)
	if err != nil {
		return "", diag.FromErr(err)
	}
	if condition == nil {
		return "", diag.Errorf("%s %q does not exist", field, referenced)
	}
	if condition.Kind != expected {
		return "", diag.Errorf("%s %q must reference a %s condition, but it is a %s", field, referenced, expected, condition.Kind)
	}
	return expected, nil
}

func resourceAbpConditionListEntryCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	kind, diags := validateReferencedConditionKind(client, data)
	if diags.HasError() {
		return diags
	}

	created, err := client.CreateAbpCondition(accountId, extractAbpConditionListEntry(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("%s create response did not contain an id", abpConditionListEntryResourceName)
	}

	data.SetId(created.Id)
	if err := serializeAbpConditionListEntry(data, created, kind); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created %s %s in account %s", abpConditionListEntryResourceName, created.Id, accountId)
	return nil
}

func resourceAbpConditionListEntryRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	ref, err := client.ReadAbpCondition(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if ref == nil {
		log.Printf("[INFO] %s %s not found, removing from state", abpConditionListEntryResourceName, id)
		data.SetId("")
		return nil
	}

	referenced, err := client.ReadAbpCondition(ref.Reference)
	if err != nil {
		return diag.FromErr(err)
	}
	if referenced == nil {
		return diag.Errorf("referenced condition %s no longer exists", ref.Reference)
	}

	if err := serializeAbpConditionListEntry(data, ref, referenced.Kind); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpConditionListEntryUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	// condition_id / condition_list_id / parent_condition_list_id are all
	// ForceNew, so the kind is whatever the current config says it is.
	kind, diags := validateReferencedConditionKind(client, data)
	if diags.HasError() {
		return diags
	}

	updated, err := client.UpdateAbpCondition(id, extractAbpConditionListEntry(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpConditionListEntry(data, updated, kind); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpConditionListEntryDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpCondition(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpConditionListEntryImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<condition_list_entry_id>'")
	}

	client := m.(*Client)
	ref, err := client.ReadAbpCondition(id)
	if err != nil {
		return nil, err
	}
	if ref == nil {
		return nil, fmt.Errorf("%s %s not found", abpConditionListEntryResourceName, id)
	}
	if ref.Kind != AbpConditionKindReference {
		return nil, fmt.Errorf("ABP Condition %s is not a reference variant (it is a %s)", id, ref.Kind)
	}
	if ref.AccountId == "" {
		return nil, fmt.Errorf("%s %s is a managed entry and cannot be imported; only account-owned entries are supported", abpConditionListEntryResourceName, id)
	}

	referenced, err := client.ReadAbpCondition(ref.Reference)
	if err != nil {
		return nil, err
	}
	if referenced == nil {
		return nil, fmt.Errorf("referenced condition %s no longer exists", ref.Reference)
	}

	data.SetId(id)
	if err := serializeAbpConditionListEntry(data, ref, referenced.Kind); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{data}, nil
}
