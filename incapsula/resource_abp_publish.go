package incapsula

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAbpPublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpPublishCreate,
		ReadContext:   resourceAbpPublishRead,
		DeleteContext: resourceAbpPublishDelete,

		Description: `Publishes a specific ABP preflight, making the configuration captured by
that preflight active on the Analysis Host.

NOTE: publishing activates the **entire account** configuration as it was when
the preflight was taken — not just the objects this Terraform configuration
manages. Objects created in the UI or by another workspace are published along
with yours.

NOTE: publishing is asynchronous. The API accepts the request and this resource
succeeds as soon as it has been accepted; it does not wait for the publish to
finish and does not track its outcome, so a publish that fails afterwards is not
reported here and is not retried. Check the publish status in the UI or the API
if it matters.

NOTE: if the preflight has become unusable by the time the publish runs — it was
deleted, or the account was published from elsewhere in the meantime — the apply
fails with an explanation. Running ` + "`terraform apply`" + ` again recovers:
the snapshot is taken anew and that one is published.

NOTE: the resource's ` + "`id`" + ` is the ID of the published *preflight*, not
of the publish itself. Destroying the resource only drops it from Terraform
state: nothing is unpublished or rolled back, and a ` + "`terraform destroy`" + `
of a whole configuration leaves the account's live configuration as it was at the
last publish. To go back to an earlier configuration, publish an earlier
preflight.`,

		Schema: map[string]*schema.Schema{
			"preflight_id": {
				Description: "The preflight to publish. Changing it replaces the resource, which " +
					"is how another publish is triggered: Terraform forgets the previous publish " +
					"(it is not undone) and publishes the new preflight.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},
		},
	}
}

func resourceAbpPublishCreate(ctx context.Context, data *schema.ResourceData, m any) diag.Diagnostics {
	client := m.(*Client)
	preflightId := data.Get("preflight_id").(string)

	preflight, diags := client.GetAbpPreflight(preflightId)
	if diags.HasError() {
		return diags
	}

	if preflight == nil {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ABP preflight no longer exists",
			Detail: fmt.Sprintf("Preflight %s was not found, so there is nothing to publish. "+
				"Run `terraform apply` again: the snapshot is taken anew and then published.", preflightId),
		})
	}

	if !preflight.CanPublish() {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "ABP preflight is not publishable",
			Detail: fmt.Sprintf("Preflight %s is older than the configuration last published to the "+
				"account, so it can no longer be published. Run `terraform apply` again: the snapshot "+
				"is taken anew and then published.", preflightId),
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
