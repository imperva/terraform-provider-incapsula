package incapsula

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// aiApplicationSecurityGuardrailExceptionReasons are the only accepted exception_reason values (exact, case-sensitive).
var aiApplicationSecurityGuardrailExceptionReasons = []string{"BUSINESS_REQUIREMENT", "FALSE_POSITIVE", "OTHER"}

func resourceAiApplicationSecurityGuardrailException() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAiApplicationSecurityGuardrailExceptionCreate,
		ReadContext:   resourceAiApplicationSecurityGuardrailExceptionRead,
		UpdateContext: resourceAiApplicationSecurityGuardrailExceptionUpdate,
		DeleteContext: resourceAiApplicationSecurityGuardrailExceptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAiApplicationSecurityGuardrailExceptionImport,
		},
		CustomizeDiff: resourceAiApplicationSecurityGuardrailExceptionCustomizeDiff,

		Description: "Manages an AI Application Security guardrail exception: a pattern that one guardrail of a policy " +
			"lets through. The guardrail is chosen by its type and phase; the provider looks it up in the policy. " +
			"Every change except the account, application or policy is made in place (same id).",

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
				Description:  "The UUID of the AI Application Security application the policy belongs to.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"policy_id": {
				Type:         schema.TypeString,
				Description:  "The UUID of the AI Application Security policy whose guardrail this exception applies to.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"guardrail_type": {
				Type:         schema.TypeString,
				Description:  "Type of the policy's guardrail this exception applies to, e.g. PROMPT_INJECTION. RATE_LIMIT is not supported: it has no text to except.",
				Required:     true,
				ValidateFunc: validation.StringInSlice(aiApplicationSecurityGuardrailExceptionGuardrailTypes(), false),
			},
			"phase": {
				Type:         schema.TypeString,
				Description:  "Phase of that guardrail: PROMPT or RESPONSE.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{aiApplicationSecurityGuardrailPhasePrompt, aiApplicationSecurityGuardrailPhaseResponse}, false),
			},
			"exception_pattern": {
				Type: schema.TypeString,
				Description: "The text the guardrail lets through (by default 3 to 200 words). The service removes leading and " +
					"trailing whitespace when it saves it, so such whitespace is not a change.",
				Required:         true,
				ValidateFunc:     validation.StringIsNotWhiteSpace,
				DiffSuppressFunc: suppressOuterWhitespaceDiff,
			},
			"exception_reason": {
				Type:         schema.TypeString,
				Description:  "Why the exception exists: BUSINESS_REQUIREMENT, FALSE_POSITIVE or OTHER.",
				Required:     true,
				ValidateFunc: validation.StringInSlice(aiApplicationSecurityGuardrailExceptionReasons, false),
			},
			"note": {
				Type:             schema.TypeString,
				Description:      "Optional note for the people managing the exception; not sent to the gateway. Leading and trailing whitespace is removed when saved.",
				Optional:         true,
				DiffSuppressFunc: suppressOuterWhitespaceDiff,
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the exception is applied. Defaults to true.",
				Optional:    true,
				Default:     true,
			},
			"guardrail_id": {
				Type:        schema.TypeInt,
				Description: "Id of the guardrail the exception is attached to, looked up in the policy from guardrail_type and phase.",
				Computed:    true,
			},
			"author": {
				Type:        schema.TypeString,
				Description: "The user who created the exception.",
				Computed:    true,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Description: "The user who last changed the exception; empty if it was never changed.",
				Computed:    true,
			},
		},
	}
}

// aiApplicationSecurityGuardrailExceptionGuardrailTypes are the guardrail types an exception can be attached to:
// all but RATE_LIMIT, which has no text to except (the management UI leaves it out too).
func aiApplicationSecurityGuardrailExceptionGuardrailTypes() []string {
	var types []string
	for _, t := range aiApplicationSecurityGuardrailTypes() {
		if t != "RATE_LIMIT" {
			types = append(types, t)
		}
	}
	return types
}

// resourceAiApplicationSecurityGuardrailExceptionCustomizeDiff marks guardrail_id as unknown when the guardrail changes,
// so anything that reads it during the same apply gets the new id, not the old one.
func resourceAiApplicationSecurityGuardrailExceptionCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	if d.Id() != "" && (d.HasChange("guardrail_type") || d.HasChange("phase")) {
		return d.SetNewComputed("guardrail_id")
	}
	return nil
}

// suppressOuterWhitespaceDiff ignores a difference only in leading and trailing whitespace - the service trims it.
func suppressOuterWhitespaceDiff(k, old, new string, d *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

// guardrailsOfPhase returns the policy's guardrails of a phase: PROMPT ones are listed under request, RESPONSE ones under response.
func guardrailsOfPhase(policy *AiApplicationSecurityPolicyResponse, phase string) []AiApplicationSecurityGuardrail {
	if phase == aiApplicationSecurityGuardrailPhaseResponse {
		return policy.Response
	}
	return policy.Request
}

// guardrailIdOf returns the id of the policy's only guardrail with this type in this phase.
// None, or more than one, is an error - the provider never guesses.
func guardrailIdOf(policy *AiApplicationSecurityPolicyResponse, guardrailType string, phase string) (int64, error) {
	var ids []int64
	for _, g := range guardrailsOfPhase(policy, phase) {
		if g.GuardrailType == guardrailType {
			ids = append(ids, g.Id)
		}
	}
	switch len(ids) {
	case 1:
		return ids[0], nil
	case 0:
		return 0, fmt.Errorf("policy %s has no %s guardrail in phase %s", policy.Id, guardrailType, phase)
	default:
		return 0, fmt.Errorf("policy %s has %d %s guardrails in phase %s - can't tell which one to use", policy.Id, len(ids), guardrailType, phase)
	}
}

// guardrailTypeAndPhaseOf maps a guardrail id back to its type and phase in the policy.
func guardrailTypeAndPhaseOf(policy *AiApplicationSecurityPolicyResponse, guardrailID int64) (string, string, bool) {
	for _, phase := range []string{aiApplicationSecurityGuardrailPhasePrompt, aiApplicationSecurityGuardrailPhaseResponse} {
		for _, g := range guardrailsOfPhase(policy, phase) {
			if g.Id == guardrailID {
				return g.GuardrailType, phase, true
			}
		}
	}
	return "", "", false
}

func readAiApplicationSecurityGuardrailExceptionPolicy(client *Client, d *schema.ResourceData) (*AiApplicationSecurityPolicyResponse, error) {
	policyID := d.Get("policy_id").(string)
	policy, err := client.GetAiApplicationSecurityPolicy(d.Get("account_id").(int), d.Get("application_id").(string), policyID)
	if err != nil {
		return nil, err
	}
	if policy == nil {
		return nil, fmt.Errorf("policy %s not found", policyID)
	}
	return policy, nil
}

// buildAiApplicationSecurityGuardrailExceptionRequest looks up the chosen guardrail in the policy and builds the write body.
func buildAiApplicationSecurityGuardrailExceptionRequest(client *Client, d *schema.ResourceData) (*AiApplicationSecurityGuardrailExceptionRequest, error) {
	policy, err := readAiApplicationSecurityGuardrailExceptionPolicy(client, d)
	if err != nil {
		return nil, err
	}
	guardrailID, err := guardrailIdOf(policy, d.Get("guardrail_type").(string), d.Get("phase").(string))
	if err != nil {
		return nil, err
	}
	return &AiApplicationSecurityGuardrailExceptionRequest{
		GuardrailId:      guardrailID,
		ExceptionPattern: d.Get("exception_pattern").(string),
		ExceptionReason:  d.Get("exception_reason").(string),
		Note:             d.Get("note").(string),
		Active:           d.Get("active").(bool),
	}, nil
}

func resourceAiApplicationSecurityGuardrailExceptionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	req, err := buildAiApplicationSecurityGuardrailExceptionRequest(client, d)
	if err != nil {
		return diag.Errorf("Failed to create AI Application Security guardrail exception: %s", err)
	}

	created, err := client.CreateAiApplicationSecurityGuardrailException(d.Get("account_id").(int), d.Get("application_id").(string), d.Get("policy_id").(string), req)
	if err != nil {
		return diag.Errorf("Failed to create AI Application Security guardrail exception: %s", err)
	}

	d.SetId(created.Id)
	log.Printf("[INFO] Created AI Application Security guardrail exception with ID: %s", created.Id)
	return resourceAiApplicationSecurityGuardrailExceptionRead(ctx, d, m)
}

func resourceAiApplicationSecurityGuardrailExceptionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	guardrailException, err := client.GetAiApplicationSecurityGuardrailException(d.Get("account_id").(int), d.Get("application_id").(string), d.Get("policy_id").(string), d.Id())
	if err != nil {
		return diag.Errorf("Failed to read AI Application Security guardrail exception %s: %s", d.Id(), err)
	}
	if guardrailException == nil {
		log.Printf("[INFO] AI Application Security guardrail exception %s not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	// the response names the guardrail by id only - type and phase come from the policy
	policy, err := readAiApplicationSecurityGuardrailExceptionPolicy(client, d)
	if err != nil {
		return diag.Errorf("Failed to read AI Application Security guardrail exception %s: %s", d.Id(), err)
	}
	guardrailType, phase, found := guardrailTypeAndPhaseOf(policy, guardrailException.GuardrailId)
	if !found {
		// e.g. the policy's guardrails were recreated outside Terraform - the next apply attaches the exception again
		log.Printf("[WARN] Guardrail %d of AI Application Security guardrail exception %s is no longer in its policy", guardrailException.GuardrailId, d.Id())
	}

	if d.Get("account_id").(int) == 0 && policy.AccountId != 0 {
		d.Set("account_id", int(policy.AccountId))
	}
	d.Set("guardrail_type", guardrailType)
	d.Set("phase", phase)
	d.Set("guardrail_id", int(guardrailException.GuardrailId))
	d.Set("exception_pattern", guardrailException.ExceptionPattern)
	d.Set("exception_reason", guardrailException.ExceptionReason)
	d.Set("note", stringValueOrEmpty(guardrailException.Note))
	d.Set("active", guardrailException.Active)
	d.Set("author", guardrailException.Author)
	d.Set("updated_by", stringValueOrEmpty(guardrailException.UpdatedBy))
	return nil
}

func resourceAiApplicationSecurityGuardrailExceptionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	req, err := buildAiApplicationSecurityGuardrailExceptionRequest(client, d)
	if err != nil {
		// keep the prior state: nothing was changed, and the planned values must not be saved as if they were
		d.Partial(true)
		return diag.Errorf("Failed to update AI Application Security guardrail exception %s: %s", d.Id(), err)
	}

	if _, err := client.UpdateAiApplicationSecurityGuardrailException(d.Get("account_id").(int), d.Get("application_id").(string), d.Get("policy_id").(string), d.Id(), req); err != nil {
		d.Partial(true)
		return diag.Errorf("Failed to update AI Application Security guardrail exception %s: %s", d.Id(), err)
	}
	return resourceAiApplicationSecurityGuardrailExceptionRead(ctx, d, m)
}

func resourceAiApplicationSecurityGuardrailExceptionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	if err := client.DeleteAiApplicationSecurityGuardrailException(d.Get("account_id").(int), d.Get("application_id").(string), d.Get("policy_id").(string), d.Id()); err != nil {
		return diag.Errorf("Failed to delete AI Application Security guardrail exception %s: %s", d.Id(), err)
	}
	d.SetId("")
	return nil
}

// resourceAiApplicationSecurityGuardrailExceptionImport takes '<application_id>/<policy_id>/<id>', optionally prefixed
// with '<account_id>/' - the service finds a guardrail exception only within its policy.
func resourceAiApplicationSecurityGuardrailExceptionImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(strings.TrimSpace(d.Id()), "/")
	if len(parts) == 4 {
		caid, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid account_id %q: must be numeric", parts[0])
		}
		d.Set("account_id", caid)
		parts = parts[1:]
	}
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, fmt.Errorf("invalid import ID %q: want '<application_id>/<policy_id>/<id>' or '<account_id>/<application_id>/<policy_id>/<id>'", d.Id())
	}

	d.Set("application_id", parts[0])
	d.Set("policy_id", parts[1])
	d.SetId(parts[2])
	return []*schema.ResourceData{d}, nil
}

// stringValueOrEmpty returns the pointed-to string, or "" for nil (the service returns null for an unset note).
func stringValueOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
