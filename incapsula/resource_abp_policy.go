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

// abpDirectiveActionRegexp matches a valid directive action: snake_case only.
var abpDirectiveActionRegexp = regexp.MustCompile(`^[_a-z][_a-z0-9]*$`)

type AbpPolicy struct {
	Id          string         `json:"id,omitempty"`
	AccountId   string         `json:"account_id,omitempty"`
	Name        string         `json:"name"`
	Description *string        `json:"description,omitempty"`
	Directives  []AbpDirective `json:"directives"`
	CreatedAt   string         `json:"created_at,omitempty"`
	ModifiedAt  string         `json:"modified_at,omitempty"`
}

type AbpDirective struct {
	Action                     string  `json:"action"`
	ConditionId                *string `json:"condition_id,omitempty"`
	SkipConditionId            *string `json:"skip_condition_id,omitempty"`
	ProofOfWorkConfigurationId *string `json:"proof_of_work_configuration_id,omitempty"`
}

func resourceAbpPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPolicyCreate,
		ReadContext:   resourceAbpPolicyRead,
		UpdateContext: resourceAbpPolicyUpdate,
		DeleteContext: resourceAbpPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAbpPolicyImport,
		},
		CustomizeDiff: resourceAbpPolicyCustomizeDiff,

		Description: `Provides an ABP Policy resource. A Policy is an ordered collection of
directives, each pairing an action with the conditions that trigger it, applied
to traffic that a Site's selector maps to this Policy.

All ` + "`directive`" + ` blocks must be specified explicitly. Conditions are attached to
a directive with ` + "`incapsula_abp_condition_list_entry`" + `, using the directive's
condition list as the entry's ` + "`parent_condition_list_id`" + `. Look that condition
list up with the ` + "`incapsula_abp_directive`" + ` data source, which addresses a
directive by its ` + "`action`" + ` (see the example below).

Do not read ` + "`directive[i].condition_list_id`" + ` off this resource instead. A list
index is not a stable handle on a directive, and on the apply that adds, removes
or reorders directives those nested ids are still reported at their old
positions, so conditions can be attached to the wrong directive or the apply can
fail outright. The data source re-reads them from the API and is unaffected.

A directive is identified by its ` + "`action`" + `, which is how the backend-owned
condition lists are kept attached to the right directive when directives are
added, removed or reordered. Two directives in the same policy therefore must
not share an ` + "`action`" + `.`,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "ABP account UUID this Policy belongs to. Changing it replaces the " +
					"Policy — Terraform deletes the existing one, together with its directives' " +
					"condition lists, and creates a new one — because the API cannot move a Policy " +
					"to another account.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Description:  "Policy name. 1..100 characters, without leading or trailing whitespace.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringLenBetween(1, 100), validateAbpName),
			},
			"description": {
				Description: "Description of the policy. Optional: the API stores no description " +
					"when it is omitted, which is reflected in state as the empty string.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"directive": {
				Description: "Ordered list of directives evaluated top-down for this policy. A policy must have at least one directive, and no two directives may share an `action`.",
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Description: "Action to take when this directive matches. Must be snake_case, 1..100 characters.",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: validation.All(
								validation.StringLenBetween(1, 100),
								validation.StringMatch(
									abpDirectiveActionRegexp,
									"action must be snake_case matching ^[_a-z][_a-z0-9]*$",
								),
							),
						},
						// Known limitation: these two are backend-owned ids
						// nested inside a TypeList, and SDKv2 keeps such a
						// computed value attached to its list index. Whenever
						// `directive` changes shape they are therefore reported
						// at their old position for the duration of that plan:
						// a newly added directive plans its id as null instead
						// of unknown, so a resource consuming it fails the
						// apply with "Provider produced inconsistent final
						// plan", and after a reorder a consumer silently
						// receives the id of whichever directive previously sat
						// at that index.
						//
						// There is no escape hatch: SetNew / SetNewComputed
						// (helper/schema/resource_diff.go) go through
						// checkKey(key, caller, false), which resolves
						// top-level schema keys only, so `directive.N.*` cannot
						// be marked unknown from CustomizeDiff. Fixing it means
						// publishing the ids as top-level action-keyed maps, or
						// moving the resource to terraform-plugin-framework.
						//
						// Until then the supported way to attach conditions to
						// a directive is the `incapsula_abp_directive` data
						// source, which re-reads the ids from the API and so is
						// immune to both problems. Prefer it over these
						// attributes in docs and examples.
						"condition_list_id": {
							Description: "Condition list containing conditions for this directive.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"skip_condition_list_id": {
							Description: "Condition list whose matches skip this directive. Only meaningful when `action` is `proof_of_work`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"proof_of_work_configuration_id": {
							Description: "ID of the proof-of-work configuration to apply. Only valid when `action` is `proof_of_work`.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// resourceAbpPolicyCustomizeDiff rejects two directives with the same action.
// Directive identity comes from the action (see extractAbpDirectives), so
// duplicates would make it impossible to tell which backend-owned condition
// list belongs to which directive. They are also unaddressable through the
// incapsula_abp_directive data source.
//
// Actions that are still unknown at plan time are skipped: Terraform re-plans
// every instance just before applying it, with the configuration fully
// resolved, so a duplicate hidden behind an unknown is caught then.
func resourceAbpPolicyCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, m any) error {
	raw, ok := diff.Get("directive").([]any)
	if !ok {
		return nil
	}
	seen := make(map[string]bool, len(raw))
	for i, item := range raw {
		directive, ok := item.(map[string]any)
		if !ok {
			continue
		}
		if !diff.NewValueKnown(fmt.Sprintf("directive.%d.action", i)) {
			continue
		}
		action, _ := directive["action"].(string)
		// Required and non-empty by validation; guards a nested object that
		// carries no action at all.
		if action == "" {
			continue
		}
		if seen[action] {
			return fmt.Errorf("policy has more than one directive with action %q; each directive in a policy must have a unique action", action)
		}
		seen[action] = true
	}
	return nil
}

// extractAbpPolicy builds the request body from the configuration. An omitted
// description reaches here as the empty string, which is left out of the
// request so the API stores no description at all (`description` is
// `Option<String>` in both CreatePolicyV1 and UpdatePolicyV1, and an update
// that omits it clears it). The value read back is normalised to the empty
// string again by serializeAbpPolicy, so this round-trips without drift.
func extractAbpPolicy(data *schema.ResourceData) AbpPolicy {
	policy := AbpPolicy{
		Name:       data.Get("name").(string),
		Directives: extractAbpDirectives(data),
	}
	if description := data.Get("description").(string); description != "" {
		policy.Description = &description
	}
	return policy
}

// abpConditionLists holds the condition lists the backend has allocated for a
// single directive.
type abpConditionLists struct {
	conditionId     string
	skipConditionId string
}

// priorAbpConditionLists indexes the condition lists recorded in state by
// directive action. The backend owns these ids and an update must send each of
// them back on the directive it belongs to: an id that is not sent is deleted
// together with all of its entries, and an id sent on a different directive
// silently re-targets every condition attached to it onto another action.
//
// State is read rather than the post-diff value because SDKv2 keeps a computed
// nested attribute attached to its list index, which does not survive
// directives being inserted, removed or reordered.
func priorAbpConditionLists(data *schema.ResourceData) map[string]abpConditionLists {
	prior, _ := data.GetChange("directive")
	raw, ok := prior.([]any)
	if !ok {
		return nil
	}
	lists := make(map[string]abpConditionLists, len(raw))
	for _, item := range raw {
		directive, ok := item.(map[string]any)
		if !ok {
			continue
		}
		action, _ := directive["action"].(string)
		if action == "" {
			continue
		}
		// Actions are unique within a policy (see
		// resourceAbpPolicyCustomizeDiff), but a policy changed out of band
		// may still hold duplicates; keep the first one.
		if _, exists := lists[action]; exists {
			continue
		}
		conditionId, _ := directive["condition_list_id"].(string)
		skipConditionId, _ := directive["skip_condition_list_id"].(string)
		lists[action] = abpConditionLists{conditionId: conditionId, skipConditionId: skipConditionId}
	}
	return lists
}

func extractAbpDirectives(data *schema.ResourceData) []AbpDirective {
	prior := priorAbpConditionLists(data)
	raw := data.Get("directive").([]any)
	directives := make([]AbpDirective, len(raw))
	for i, item := range raw {
		m := item.(map[string]any)
		action := m["action"].(string)
		d := AbpDirective{Action: action}
		// A directive whose action is not in state is new to this policy: its
		// condition lists are left out so the backend allocates them. On
		// create nothing is in state, which also keeps condition_id out of the
		// request, where the API does not accept it.
		if lists, found := prior[action]; found {
			if lists.conditionId != "" {
				conditionId := lists.conditionId
				d.ConditionId = &conditionId
			}
			if lists.skipConditionId != "" {
				skipConditionId := lists.skipConditionId
				d.SkipConditionId = &skipConditionId
			}
			delete(prior, action)
		}
		if powId, ok := m["proof_of_work_configuration_id"].(string); ok && powId != "" {
			d.ProofOfWorkConfigurationId = &powId
		}
		directives[i] = d
	}
	return directives
}

func flattenAbpDirectives(directives []AbpDirective) []any {
	out := make([]any, len(directives))
	for i, d := range directives {
		m := map[string]any{"action": d.Action}
		if d.ConditionId != nil {
			m["condition_list_id"] = *d.ConditionId
		}
		if d.SkipConditionId != nil {
			m["skip_condition_list_id"] = *d.SkipConditionId
		}
		if d.ProofOfWorkConfigurationId != nil {
			m["proof_of_work_configuration_id"] = *d.ProofOfWorkConfigurationId
		}
		out[i] = m
	}
	return out
}

// serializeAbpPolicy writes the server's view of a policy back into state.
func serializeAbpPolicy(data *schema.ResourceData, policy *AbpPolicy) error {
	if policy.AccountId != "" {
		if err := data.Set("account_id", policy.AccountId); err != nil {
			return err
		}
	}
	if err := data.Set("name", policy.Name); err != nil {
		return err
	}
	description := ""
	if policy.Description != nil {
		description = *policy.Description
	}
	if err := data.Set("description", description); err != nil {
		return err
	}
	if err := data.Set("directive", flattenAbpDirectives(policy.Directives)); err != nil {
		return err
	}
	return nil
}

func resourceAbpPolicyCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	accountId := data.Get("account_id").(string)

	created, err := client.CreateAbpPolicy(accountId, extractAbpPolicy(data))
	if err != nil {
		return diag.FromErr(err)
	}
	if created.Id == "" {
		return diag.Errorf("ABP Policy create response did not contain an id")
	}

	data.SetId(created.Id)
	if err := serializeAbpPolicy(data, created); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Created ABP Policy %s in account %s", created.Id, accountId)
	return nil
}

func resourceAbpPolicyRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	policy, err := client.ReadAbpPolicy(id)
	if err != nil {
		return diag.FromErr(err)
	}

	if policy == nil {
		log.Printf("[INFO] ABP Policy %s not found, removing from state", id)
		data.SetId("")
		return nil
	}

	if err := serializeAbpPolicy(data, policy); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpPolicyUpdate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	updated, err := client.UpdateAbpPolicy(id, extractAbpPolicy(data))
	if err != nil {
		return diag.FromErr(err)
	}

	if updated == nil {
		data.SetId("")
		return nil
	}

	if err := serializeAbpPolicy(data, updated); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceAbpPolicyDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	id := data.Id()

	if err := client.DeleteAbpPolicy(id); err != nil {
		return diag.FromErr(err)
	}

	data.SetId("")
	return nil
}

func resourceAbpPolicyImport(ctx context.Context, data *schema.ResourceData, m any) ([]*schema.ResourceData, error) {
	id := strings.TrimSpace(data.Id())
	if id == "" {
		return nil, fmt.Errorf("expected import ID to be '<policy_id>'")
	}

	client := m.(*Client)
	policy, err := client.ReadAbpPolicy(id)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("ABP Policy %s not found", id)
	}

	data.SetId(id)
	if err := serializeAbpPolicy(data, policy); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}
