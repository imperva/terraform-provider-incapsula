package incapsula

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// Guardrail phase values.
const (
	aiApplicationSecurityGuardrailPhasePrompt   = "PROMPT"
	aiApplicationSecurityGuardrailPhaseResponse = "RESPONSE"
)

// aiApplicationSecurityGuardrailValidPhases maps each guardrail type to the phases it may run in,
// mirroring the backend's per-guardrail phase validation. Enforced at plan time by
// resourceAiApplicationSecurityPolicyCustomizeDiff so misconfigurations fail before any API call.
var aiApplicationSecurityGuardrailValidPhases = map[string][]string{
	"PROMPT_INJECTION":         {aiApplicationSecurityGuardrailPhasePrompt},
	"MODERATION":               {aiApplicationSecurityGuardrailPhaseResponse},
	"SYSTEM_PROMPT_LEAK":       {aiApplicationSecurityGuardrailPhaseResponse},
	"PII_STATIC":               {aiApplicationSecurityGuardrailPhasePrompt, aiApplicationSecurityGuardrailPhaseResponse},
	"RATE_LIMIT":               {aiApplicationSecurityGuardrailPhasePrompt, aiApplicationSecurityGuardrailPhaseResponse},
	"ZERO_SHOT_CLASSIFICATION": {aiApplicationSecurityGuardrailPhasePrompt, aiApplicationSecurityGuardrailPhaseResponse},
}

func aiApplicationSecurityGuardrailTypes() []string {
	types := make([]string, 0, len(aiApplicationSecurityGuardrailValidPhases))
	for t := range aiApplicationSecurityGuardrailValidPhases {
		types = append(types, t)
	}
	return types
}

func resourceAiApplicationSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAiApplicationSecurityPolicyCreate,
		ReadContext:   resourceAiApplicationSecurityPolicyRead,
		UpdateContext: resourceAiApplicationSecurityPolicyUpdate,
		DeleteContext: resourceAiApplicationSecurityPolicyDelete,
		CustomizeDiff: resourceAiApplicationSecurityPolicyCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAiApplicationSecurityPolicyImport,
		},

		Description: "Manages an AI Application Security policy for an application. A policy holds an ordered " +
			"set of guardrails, each attached to the PROMPT (request) or RESPONSE phase.",

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
				Description:  "The UUID of the AI Application Security application this policy belongs to.",
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
			"name": {
				Type:         schema.TypeString,
				Description:  "The policy name.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Type:         schema.TypeString,
				Description:  "An optional description of the policy.",
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 500),
			},
			"active": {
				Type:        schema.TypeBool,
				Description: "Whether the policy is active. Defaults to true.",
				Optional:    true,
				Default:     true,
			},
			"guardrail": {
				Type:        schema.TypeSet,
				Description: "The set of guardrails enforced by this policy. At least one is required.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Description:  "The guardrail type.",
							Required:     true,
							ValidateFunc: validation.StringInSlice(aiApplicationSecurityGuardrailTypes(), false),
						},
						"phase": {
							Type:         schema.TypeString,
							Description:  "The phase the guardrail runs in: PROMPT (request) or RESPONSE.",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{aiApplicationSecurityGuardrailPhasePrompt, aiApplicationSecurityGuardrailPhaseResponse}, false),
						},
						"mode": {
							Type:         schema.TypeString,
							Description:  "The enforcement mode: BLOCK, ALERT, or MASK.",
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"BLOCK", "ALERT", "MASK"}, false),
						},
						"active": {
							Type:        schema.TypeBool,
							Description: "Whether this guardrail is active. Defaults to true.",
							Optional:    true,
							Default:     true,
						},
						"config": {
							Type:         schema.TypeString,
							Description:  "Guardrail-specific configuration as a JSON object. Defaults to an empty object.",
							Optional:     true,
							Default:      "{}",
							ValidateFunc: validation.StringIsJSON,
							// guardrail is a TypeSet, whose element hash is computed from the raw
							// attribute values BEFORE DiffSuppressFunc runs — so a suppressor alone
							// cannot prevent spurious add/remove diffs when the backend returns
							// semantically-equal but reordered JSON. The set's Set hash function (see
							// below) normalizes this config before hashing so both sides hash
							// identically; DiffSuppressFunc then suppresses the intra-element string
							// diff for the (now hash-matched) element.
							DiffSuppressFunc: suppressEquivalentJSONStringDiffs,
						},
					},
				},
				// A guardrail's config is JSON whose key order is not significant, but a TypeSet
				// hashes the raw attribute values, so a config written with keys in one order
				// would hash differently than the same config read back sorted (see
				// flattenAiApplicationSecurityGuardrail) — producing a perpetual remove/add diff. Hash the
				// normalized config instead so semantically-equal guardrails land in the same
				// set element.
				Set: aiApplicationSecurityGuardrailHash,
			},
		},
	}
}

// aiApplicationSecurityGuardrailHash computes the set hash for a guardrail element from its identifying
// fields plus its normalized (canonical-JSON) config, so that guardrails differing only in
// config key order hash identically.
func aiApplicationSecurityGuardrailHash(v interface{}) int {
	m := v.(map[string]interface{})
	var b strings.Builder
	b.WriteString(m["type"].(string))
	b.WriteString("|")
	b.WriteString(m["phase"].(string))
	b.WriteString("|")
	b.WriteString(m["mode"].(string))
	b.WriteString("|")
	b.WriteString(strconv.FormatBool(m["active"].(bool)))
	b.WriteString("|")
	b.WriteString(normalizeAiApplicationSecurityGuardrailConfig(m["config"]))
	return schema.HashString(b.String())
}

// resourceAiApplicationSecurityPolicyCustomizeDiff validates each guardrail's type/phase combination
// against aiApplicationSecurityGuardrailValidPhases at plan time.
func resourceAiApplicationSecurityPolicyCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	guardrails, ok := d.Get("guardrail").(*schema.Set)
	if !ok {
		return nil
	}
	for _, raw := range guardrails.List() {
		g := raw.(map[string]interface{})
		guardrailType := g["type"].(string)
		phase := g["phase"].(string)
		validPhases, known := aiApplicationSecurityGuardrailValidPhases[guardrailType]
		if !known {
			// Unknown type is caught by the schema ValidateFunc; skip here.
			continue
		}
		if !stringInSlice(validPhases, phase) {
			return fmt.Errorf("guardrail type %q is not valid in phase %q; allowed phases: %s",
				guardrailType, phase, strings.Join(validPhases, ", "))
		}
	}
	return nil
}

func stringInSlice(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func resourceAiApplicationSecurityPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Get("application_id").(string)

	req, err := buildAiApplicationSecurityPolicyRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	policy, err := client.CreateAiApplicationSecurityPolicy(accountID, applicationID, req)
	if err != nil {
		return diag.Errorf("Failed to create AI Application Security policy for application %s: %s", applicationID, err)
	}

	d.SetId(policy.Id)

	return resourceAiApplicationSecurityPolicyRead(ctx, d, m)
}

func resourceAiApplicationSecurityPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Get("application_id").(string)
	if applicationID == "" {
		// On policyId-only import the application id is unknown; the backend looks the
		// policy up by id regardless of the path segment.
		applicationID = aiApplicationSecurityPolicyImportApplicationPlaceholder
	}
	policyID := d.Id()

	policy, err := client.GetAiApplicationSecurityPolicy(accountID, applicationID, policyID)
	if err != nil {
		return diag.Errorf("Failed to read AI Application Security policy %s: %s", policyID, err)
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	if policy.AccountId != 0 {
		d.Set("account_id", int(policy.AccountId))
	}
	d.Set("application_id", policy.ApplicationId)
	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("active", policy.Active)

	guardrails, err := flattenAiApplicationSecurityGuardrails(policy)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("guardrail", guardrails); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAiApplicationSecurityPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Get("application_id").(string)
	policyID := d.Id()

	req, err := buildAiApplicationSecurityPolicyRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.UpdateAiApplicationSecurityPolicy(accountID, applicationID, policyID, req); err != nil {
		return diag.Errorf("Failed to update AI Application Security policy %s: %s", policyID, err)
	}

	return resourceAiApplicationSecurityPolicyRead(ctx, d, m)
}

func resourceAiApplicationSecurityPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Get("application_id").(string)
	policyID := d.Id()

	if err := client.DeleteAiApplicationSecurityPolicy(accountID, applicationID, policyID); err != nil {
		return diag.Errorf("Failed to delete AI Application Security policy %s: %s", policyID, err)
	}

	d.SetId("")
	return nil
}

// buildAiApplicationSecurityPolicyRequest assembles the full desired-state write payload from the
// resource data, splitting guardrails into their request/response phase buckets.
func buildAiApplicationSecurityPolicyRequest(d *schema.ResourceData) (*AiApplicationSecurityPolicyRequest, error) {
	requestGuardrails, responseGuardrails, err := expandAiApplicationSecurityGuardrails(d.Get("guardrail").(*schema.Set))
	if err != nil {
		return nil, err
	}

	return &AiApplicationSecurityPolicyRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Active:      d.Get("active").(bool),
		Request:     requestGuardrails,
		Response:    responseGuardrails,
	}, nil
}

// expandAiApplicationSecurityGuardrails converts the guardrail set into the backend's phase-split
// request/response slices, injecting the "type" discriminator into each guardrail's config.
func expandAiApplicationSecurityGuardrails(set *schema.Set) (request []AiApplicationSecurityGuardrail, response []AiApplicationSecurityGuardrail, err error) {
	request = []AiApplicationSecurityGuardrail{}
	response = []AiApplicationSecurityGuardrail{}

	for _, raw := range set.List() {
		g := raw.(map[string]interface{})
		guardrail, err := buildAiApplicationSecurityGuardrail(g)
		if err != nil {
			return nil, nil, err
		}
		if guardrail.GuardrailPhase == aiApplicationSecurityGuardrailPhasePrompt {
			request = append(request, guardrail)
		} else {
			response = append(response, guardrail)
		}
	}

	return request, response, nil
}

// buildAiApplicationSecurityGuardrail maps one guardrail set element to an AiApplicationSecurityGuardrail, injecting
// the guardrail type into the config JSON as the "type" discriminator required by the
// backend's polymorphic BaseGuardrailConfigDto.
func buildAiApplicationSecurityGuardrail(g map[string]interface{}) (AiApplicationSecurityGuardrail, error) {
	guardrailType := g["type"].(string)

	configStr, _ := g["config"].(string)
	if strings.TrimSpace(configStr) == "" {
		configStr = "{}"
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal([]byte(configStr), &configMap); err != nil {
		return AiApplicationSecurityGuardrail{}, fmt.Errorf("guardrail config for type %s is not valid JSON: %s", guardrailType, err)
	}
	if configMap == nil {
		configMap = map[string]interface{}{}
	}
	if _, present := configMap["type"]; !present {
		configMap["type"] = guardrailType
	}

	configRaw, err := json.Marshal(configMap)
	if err != nil {
		return AiApplicationSecurityGuardrail{}, fmt.Errorf("failed to encode guardrail config for type %s: %s", guardrailType, err)
	}

	return AiApplicationSecurityGuardrail{
		GuardrailType:  guardrailType,
		GuardrailMode:  g["mode"].(string),
		GuardrailPhase: g["phase"].(string),
		Active:         g["active"].(bool),
		Config:         configRaw,
	}, nil
}

// flattenAiApplicationSecurityGuardrails merges the phase-split response slices back into a single set,
// stripping the injected "type" discriminator so it round-trips against the user's config.
func flattenAiApplicationSecurityGuardrails(policy *AiApplicationSecurityPolicyResponse) ([]interface{}, error) {
	all := make([]interface{}, 0, len(policy.Request)+len(policy.Response))

	for _, g := range policy.Request {
		f, err := flattenAiApplicationSecurityGuardrail(g, aiApplicationSecurityGuardrailPhasePrompt)
		if err != nil {
			return nil, err
		}
		all = append(all, f)
	}
	for _, g := range policy.Response {
		f, err := flattenAiApplicationSecurityGuardrail(g, aiApplicationSecurityGuardrailPhaseResponse)
		if err != nil {
			return nil, err
		}
		all = append(all, f)
	}

	return all, nil
}

// normalizeAiApplicationSecurityGuardrailConfig canonicalizes a guardrail config JSON string (sorted keys,
// no insignificant whitespace) so it matches flattenAiApplicationSecurityGuardrail's json.Marshal output.
// Because config lives inside a TypeSet, the stored value feeds the set-element hash; storing a
// canonical form on both the config and refresh sides keeps the hashes stable. Invalid JSON is
// returned unchanged so ValidateFunc surfaces the error instead of this masking it.
func normalizeAiApplicationSecurityGuardrailConfig(v interface{}) string {
	s, ok := v.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return "{}"
	}
	var parsed interface{}
	if err := json.Unmarshal([]byte(s), &parsed); err != nil {
		return s
	}
	// The "type" discriminator is injected on write and stripped on read (see
	// flattenAiApplicationSecurityGuardrail). Drop it here too so a user who embeds "type" in
	// their config JSON hashes identically to the read-back form, avoiding a perpetual diff.
	if m, ok := parsed.(map[string]interface{}); ok {
		delete(m, "type")
	}
	canonical, err := json.Marshal(parsed)
	if err != nil {
		return s
	}
	return string(canonical)
}

func flattenAiApplicationSecurityGuardrail(g AiApplicationSecurityGuardrail, phaseFallback string) (map[string]interface{}, error) {
	phase := g.GuardrailPhase
	if phase == "" {
		phase = phaseFallback
	}

	configStr := "{}"
	if len(g.Config) > 0 {
		var configMap map[string]interface{}
		if err := json.Unmarshal(g.Config, &configMap); err != nil {
			return nil, fmt.Errorf("guardrail config from backend for type %s is not valid JSON: %s", g.GuardrailType, err)
		}
		// The "type" discriminator is injected on write; strip it so state matches config.
		delete(configMap, "type")
		configRaw, err := json.Marshal(configMap)
		if err != nil {
			return nil, fmt.Errorf("failed to encode guardrail config for type %s: %s", g.GuardrailType, err)
		}
		configStr = string(configRaw)
	}

	return map[string]interface{}{
		"type":   g.GuardrailType,
		"mode":   g.GuardrailMode,
		"phase":  phase,
		"active": g.Active,
		"config": configStr,
	}, nil
}

func resourceAiApplicationSecurityPolicyImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	raw := strings.TrimSpace(d.Id())
	if raw == "" {
		return nil, fmt.Errorf("expected import ID to be '<policy_id>' or '<account_id>/<policy_id>'")
	}

	var accountID string
	resourceID := raw

	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid import ID %q: want '<policy_id>' or '<account_id>/<policy_id>'", raw)
		}
		accountID = strings.TrimSpace(parts[0])
		resourceID = strings.TrimSpace(parts[1])
	}

	// Set the canonical Terraform ID to just the resource ID.
	d.SetId(resourceID)

	// Only set account_id if it was provided.
	if accountID != "" {
		caid, err := strconv.Atoi(accountID)
		if err != nil {
			return nil, fmt.Errorf("invalid account_id %q: must be numeric", accountID)
		}
		if err := d.Set("account_id", caid); err != nil {
			return nil, fmt.Errorf("setting account_id: %w", err)
		}
	}

	return []*schema.ResourceData{d}, nil
}
