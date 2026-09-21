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

func resourceAiApplicationSecurityApplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAiApplicationSecurityApplicationCreate,
		ReadContext:   resourceAiApplicationSecurityApplicationRead,
		UpdateContext: resourceAiApplicationSecurityApplicationUpdate,
		DeleteContext: resourceAiApplicationSecurityApplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAiApplicationSecurityApplicationImport,
		},
		CustomizeDiff: resourceAiApplicationSecurityApplicationCustomizeDiff,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Numeric identifier of the account to operate on. Defaults to the account of the API credentials when omitted.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the AI Application Security application.",
			},
			"application_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"SDK", "EDGE", "API"}, false),
				Description:  "Deployment type of the application. One of SDK, EDGE, API. Immutable after creation.",
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"US", "EU", "AU", "APAC"}, false),
				Description: "Data region of the application. Omit to inherit the account's data region on create " +
					"(the backend resolves it from account-management, falling back to US if that value is absent " +
					"or unrecognised). The resolved value is stored as Computed. Not ForceNew: changing the value " +
					"updates the application in place. Account-region inheritance is create-only: removing this " +
					"attribute from configuration after it has been set does NOT revert to the account default - " +
					"the last applied value is kept in state and the backend leaves the region unchanged when no " +
					"region is sent. Set the attribute explicitly to move an application between regions.",
			},
			"configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Application configuration. Required when application_type = EDGE; not supported for SDK or API application types.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"site_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"path": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "application/json",
						},
						"prompt_location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"blocked_response_structure": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_streaming": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"request": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"content_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"role_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"response": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"role_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"content_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"finish_reason_path": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"finish_reason_value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"end_of_stream_marker": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceAiApplicationSecurityApplicationCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	applicationType := d.Get("application_type").(string)
	_, hasConfig := d.GetOk("configuration")

	switch applicationType {
	case "EDGE":
		if !hasConfig {
			return fmt.Errorf("configuration is required for EDGE application type")
		}
	case "SDK", "API":
		if hasConfig {
			return fmt.Errorf("configuration is not supported for %s application type", applicationType)
		}
	}

	return nil
}

func expandAiApplicationSecurityApplicationConfig(configList []interface{}) *AiApplicationSecurityApplicationConfig {
	if len(configList) == 0 || configList[0] == nil {
		return nil
	}
	raw := configList[0].(map[string]interface{})

	config := &AiApplicationSecurityApplicationConfig{
		SiteId:                   int64(raw["site_id"].(int)),
		Path:                     raw["path"].(string),
		ContentType:              raw["content_type"].(string),
		PromptLocation:           raw["prompt_location"].(string),
		BlockedResponseStructure: raw["blocked_response_structure"].(string),
		IsStreaming:              raw["is_streaming"].(bool),
	}

	if reqList, ok := raw["request"].([]interface{}); ok && len(reqList) > 0 && reqList[0] != nil {
		reqRaw := reqList[0].(map[string]interface{})
		config.Request = &AiApplicationSecurityStreamingRequest{
			MessagePath: reqRaw["message_path"].(string),
			ContentPath: reqRaw["content_path"].(string),
			RolePath:    reqRaw["role_path"].(string),
		}
	}

	if respList, ok := raw["response"].([]interface{}); ok && len(respList) > 0 && respList[0] != nil {
		respRaw := respList[0].(map[string]interface{})
		config.Response = &AiApplicationSecurityStreamingResponse{
			RolePath:          respRaw["role_path"].(string),
			ContentPath:       respRaw["content_path"].(string),
			FinishReasonPath:  respRaw["finish_reason_path"].(string),
			FinishReasonValue: respRaw["finish_reason_value"].(string),
			EndOfStreamMarker: respRaw["end_of_stream_marker"].(string),
		}
	}

	return config
}

func flattenAiApplicationSecurityApplicationConfig(config *AiApplicationSecurityApplicationConfig) []interface{} {
	if config == nil {
		return nil
	}

	// content_type has a non-empty schema Default ("application/json"); the backend tags
	// contentType omitempty. Fall back to the default when it comes back empty so the Default
	// doesn't re-apply on the next plan and cause a perpetual in-place update.
	contentType := config.ContentType
	if contentType == "" {
		contentType = "application/json"
	}

	result := map[string]interface{}{
		"site_id":                    int(config.SiteId),
		"path":                       config.Path,
		"content_type":               contentType,
		"prompt_location":            config.PromptLocation,
		"blocked_response_structure": config.BlockedResponseStructure,
		"is_streaming":               config.IsStreaming,
	}

	if config.Request != nil {
		result["request"] = []interface{}{
			map[string]interface{}{
				"message_path": config.Request.MessagePath,
				"content_path": config.Request.ContentPath,
				"role_path":    config.Request.RolePath,
			},
		}
	}

	if config.Response != nil {
		result["response"] = []interface{}{
			map[string]interface{}{
				"role_path":            config.Response.RolePath,
				"content_path":         config.Response.ContentPath,
				"finish_reason_path":   config.Response.FinishReasonPath,
				"finish_reason_value":  config.Response.FinishReasonValue,
				"end_of_stream_marker": config.Response.EndOfStreamMarker,
			},
		}
	}

	return []interface{}{result}
}

func resourceAiApplicationSecurityApplicationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)

	req := &AiApplicationSecurityApplicationRequest{
		Name:            d.Get("name").(string),
		ApplicationType: d.Get("application_type").(string),
		Region:          d.Get("region").(string),
		Configuration:   expandAiApplicationSecurityApplicationConfig(d.Get("configuration").([]interface{})),
	}

	app, err := client.CreateAiApplicationSecurityApplication(accountID, req)
	if err != nil {
		return diag.Errorf("Could not create AI Application Security application: %s", err)
	}

	d.SetId(app.ApplicationId)
	log.Printf("[INFO] Created AI Application Security application with ID: %s", app.ApplicationId)

	return resourceAiApplicationSecurityApplicationRead(ctx, d, m)
}

func resourceAiApplicationSecurityApplicationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Id()

	app, err := client.GetAiApplicationSecurityApplication(accountID, applicationID)
	if err != nil {
		return diag.Errorf("Could not read AI Application Security application with ID %s: %s", applicationID, err)
	}

	if app == nil {
		log.Printf("[INFO] AI Application Security application with ID %s not found, removing from state", applicationID)
		d.SetId("")
		return nil
	}

	if app.AccountId != 0 {
		d.Set("account_id", int(app.AccountId))
	}
	d.Set("name", app.Name)
	d.Set("application_type", app.ApplicationType)
	d.Set("region", app.Region)
	d.Set("configuration", flattenAiApplicationSecurityApplicationConfig(app.Configuration))

	return nil
}

func resourceAiApplicationSecurityApplicationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Id()

	req := &AiApplicationSecurityApplicationRequest{
		Name:          d.Get("name").(string),
		Configuration: expandAiApplicationSecurityApplicationConfig(d.Get("configuration").([]interface{})),
	}

	if d.HasChange("region") {
		req.Region = d.Get("region").(string)
	}

	_, err := client.UpdateAiApplicationSecurityApplication(accountID, applicationID, req)
	if err != nil {
		return diag.Errorf("Could not update AI Application Security application with ID %s: %s", applicationID, err)
	}

	return resourceAiApplicationSecurityApplicationRead(ctx, d, m)
}

func resourceAiApplicationSecurityApplicationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	accountID := d.Get("account_id").(int)
	applicationID := d.Id()

	if err := client.DeleteAiApplicationSecurityApplication(accountID, applicationID); err != nil {
		return diag.Errorf("Could not delete AI Application Security application with ID %s: %s", applicationID, err)
	}

	d.SetId("")
	return nil
}

func resourceAiApplicationSecurityApplicationImport(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	raw := strings.TrimSpace(d.Id())
	if raw == "" {
		return nil, fmt.Errorf("expected import ID to be '<application_id>' or '<account_id>/<application_id>'")
	}

	var accountID string
	resourceID := raw

	if strings.Contains(raw, "/") {
		parts := strings.SplitN(raw, "/", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("invalid import ID %q: want '<application_id>' or '<account_id>/<application_id>'", raw)
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
