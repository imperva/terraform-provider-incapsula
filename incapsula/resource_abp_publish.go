package incapsula

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAbpPublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPublishCreate,
		ReadContext:   resourceAbpPublishRead,
		DeleteContext: resourceAbpPublishDelete,

		Description: "Publishes a specific ABP preflight. Replacing the resource by setting a preflight_id fires another publish.\n",

		Schema: map[string]*schema.Schema{
			"preflight_id": {
				Description: "The preflight to publish. Change of it always triggers new publish to be created",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceAbpPublishCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	preflightId := data.Get("preflight_id").(string)

	preflightStatus, diags := client.GetAbpPreflightStatus(preflightId)
	if diags.HasError() {
		return diags
	}

	if !preflightStatus.CanPublish {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ABP preflight is not publishable",
			Detail:   fmt.Sprintf("Preflight %s s older than the latest published preflight: not allowed for publishing", preflightId),
		})
	}

	publishDiags := client.PublishAbpPreflight(preflightId)
	diags = append(diags, publishDiags...)
	if publishDiags.HasError() {
		return diags
	}

	log.Printf("[INFO] Published ABP preflight %s", preflightId)
	data.SetId(preflightId)
	return diags
}

func resourceAbpPublishRead(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	return nil
}

func resourceAbpPublishDelete(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	data.SetId("")
	return nil
}
