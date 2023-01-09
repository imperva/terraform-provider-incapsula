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

func resourceWaitingRoom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWaitingRoomCreate,
		ReadContext:   resourceWaitingRoomRead,
		UpdateContext: resourceWaitingRoomUpdate,
		DeleteContext: resourceWaitingRoomDelete,
		Importer: &schema.ResourceImporter{
			State: func(data *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(data.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/waiting_room_id", data.Id())
				}

				data.Set("site_id", idSlice[0])
				data.SetId(idSlice[1])

				return []*schema.ResourceData{data}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The waiting room name. Must be unique across all waiting room of the site.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The waiting room description.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "whether this waiting room is enabled or not.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:    true,
			},
			"html_template_base64": {
				Description: "The HTML template file path.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"filter": {
				Description: "The rule conditions that determine on which sessions this waiting room applies. (no filter means the waiting room applies for the whole site)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"bots_action_in_queuing_mode": {
				Description: "The waiting room bot handling action. Determines the waiting room behavior for legitimate bots trying to access your website during peak time",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"entrance_rate_threshold": {
				Description:      "The entrance rate activation threshold of the waiting room. The waiting room is activated when sessions per minute exceed the specified value.",
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
			"concurrent_sessions_threshold": {
				Description:      "The active users activation threshold of the waiting room. The waiting room is activated when number of active users reached specified value.",
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
			"inactivity_timeout": {
				Description:      "If waiting room conditions that limit the scope of the waiting room to a subset of the website have been defined, the user is considered active only when navigating the pages in scope of the conditions. A user who is inactive for a longer period of time is considered as having left the site.",
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          5,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},
			"queue_inactivity_timeout": {
				Description:      "Queue inactivity timeout. A user in the waiting room who is inactive for a longer period of time is considered as having left the queue. On returning to the site, the user moves to the end of the queue and needs to wait in line again if the waiting room is active.",
				Type:             schema.TypeInt,
				Optional:         true,
				Default:          1,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(1)),
			},

			//Computed values

			"account_id": {
				Description: "The account this waiting room belongs to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"created_at": {
				Description: "The waiting room creation date in milliseconds.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_at": {
				Description: "The last configuration change date in milliseconds.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_modified_by": {
				Description: "The waiting room mode. Indicates whether the waiting room is currently queuing or not.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mode": {
				Description: "The user who last modified the waiting room configuration",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceWaitingRoomCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteID := data.Get("site_id").(string)

	waitingRoom := WaitingRoomDTO{
		Name:                        data.Get("name").(string),
		Description:                 data.Get("description").(string),
		Filter:                      data.Get("filter").(string),
		HtmlTemplateBase64:          data.Get("html_template_base64").(string),
		Enabled:                     data.Get("enabled").(bool),
		BotsActionInQueuingMode:     data.Get("bots_action_in_queuing_mode").(string),
		EntranceRateThreshold:       data.Get("entrance_rate_threshold").(int),
		ConcurrentSessionsThreshold: data.Get("concurrent_sessions_threshold").(int),
		QueueInactivityTimeout:      data.Get("queue_inactivity_timeout").(int),
		InactivityTimeout:           data.Get("inactivity_timeout").(int),
	}

	if waitingRoom.EntranceRateThreshold != 0 {
		waitingRoom.EntranceRateEnabled = true
	}

	if waitingRoom.ConcurrentSessionsThreshold != 0 {
		waitingRoom.ConcurrentSessionsEnabled = true
	}

	waitingRoomDTOResponse, diags := client.CreateWaitingRoom(siteID, &waitingRoom)
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to create Waiting Room for Site ID %s", siteID)
		return diags
	} else if waitingRoomDTOResponse.Errors != nil {
		log.Printf("[ERROR] Failed to create Waiting Room for Site ID %s: %s", siteID, waitingRoomDTOResponse.Errors[0].Detail)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  waitingRoomDTOResponse.Errors[0].Title,
			Detail:   fmt.Sprintf("Failed to create Waiting Room for Site ID %s: %s", siteID, waitingRoomDTOResponse.Errors[0].Detail),
		}}
	} else if len(waitingRoomDTOResponse.Data) < 1 {
		log.Printf("[ERROR] Empty response received while creating Waiting Room for Site ID %s", siteID)
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Empty response",
			Detail:   fmt.Sprintf("[ERROR] Empty response received while creating Waiting Room for Site ID %s", siteID),
		}}
	}

	data.SetId(strconv.FormatInt(waitingRoomDTOResponse.Data[0].Id, 10))
	log.Printf("[INFO] Created Incapsula Waiting Room %d for Site: %s", waitingRoomDTOResponse.Data[0].Id, siteID)

	diags = append(diags, resourceWaitingRoomRead(ctx, data, m)[:]...)

	return diags
}

func resourceWaitingRoomRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)
	var diags diag.Diagnostics
	siteID := data.Get("site_id").(string)

	waitingRoomID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		log.Printf("[ERROR] The ID should be numeric. Currrent value: %s", data.Id())
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Waiting Room ID conversion error",
			Detail:   fmt.Sprintf("The ID should be numeric. Currrent value: %s", data.Id()),
		}}
	}

	waitingRoomDTOResponse, diags := client.ReadWaitingRoom(siteID, waitingRoomID)
	if waitingRoomDTOResponse.Errors != nil && waitingRoomDTOResponse.Errors[0].Status == 404 {
		data.SetId("")
	}
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to read Waiting Room %d for Site ID %s", waitingRoomID, siteID)
		return diags
	}

	if len(waitingRoomDTOResponse.Data) == 0 {
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Empty Response",
			Detail:   fmt.Sprintf("Error getting Waiting Room for site (%s): %d, empty result", siteID, waitingRoomID),
		}}
	}

	waitingRoom := waitingRoomDTOResponse.Data[0]

	data.Set("name", waitingRoom.Name)
	data.Set("description", waitingRoom.Description)
	data.Set("enabled", waitingRoom.Enabled)
	data.Set("html_template_base64", waitingRoom.HtmlTemplateBase64)
	data.Set("filter", waitingRoom.Filter)
	data.Set("bots_action_in_queuing_mode", waitingRoom.BotsActionInQueuingMode)
	if waitingRoom.EntranceRateEnabled {
		data.Set("entrance_rate_threshold", waitingRoom.EntranceRateThreshold)
	}
	if waitingRoom.ConcurrentSessionsEnabled {
		data.Set("concurrent_sessions_threshold", waitingRoom.ConcurrentSessionsThreshold)
	}
	data.Set("inactivity_timeout", waitingRoom.InactivityTimeout)
	data.Set("queue_inactivity_timeout", waitingRoom.QueueInactivityTimeout)
	data.Set("account_id", strconv.FormatInt(waitingRoom.AccountId, 10))
	data.Set("created_at", strconv.FormatInt(waitingRoom.CreatedAt, 10))
	data.Set("last_modified_at", strconv.FormatInt(waitingRoom.LastModifiedAt, 10))
	data.Set("last_modified_by", waitingRoom.LastModifiedBy)
	data.Set("mode", waitingRoom.Mode)

	return diags
}

func resourceWaitingRoomUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteID := data.Get("site_id").(string)

	waitingRoomID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		log.Printf("[ERROR] The ID should be numeric. Currrent value: %s", data.Id())
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Waiting Room ID conversion error",
			Detail:   fmt.Sprintf("The ID should be numeric. Currrent value: %s", data.Id()),
		}}
	}

	waitingRoom := WaitingRoomDTO{
		Id:                          waitingRoomID,
		Name:                        data.Get("name").(string),
		Description:                 data.Get("description").(string),
		Filter:                      data.Get("filter").(string),
		HtmlTemplateBase64:          data.Get("html_template_base64").(string),
		Enabled:                     data.Get("enabled").(bool),
		BotsActionInQueuingMode:     data.Get("bots_action_in_queuing_mode").(string),
		EntranceRateThreshold:       data.Get("entrance_rate_threshold").(int),
		ConcurrentSessionsThreshold: data.Get("concurrent_sessions_threshold").(int),
		QueueInactivityTimeout:      data.Get("queue_inactivity_timeout").(int),
		InactivityTimeout:           data.Get("inactivity_timeout").(int),
	}

	if waitingRoom.EntranceRateThreshold != 0 {
		waitingRoom.EntranceRateEnabled = true
	}

	if waitingRoom.ConcurrentSessionsThreshold != 0 {
		waitingRoom.ConcurrentSessionsEnabled = true
	}

	waitingRoomDTOResponse, diags := client.UpdateWaitingRoom(siteID, waitingRoomID, &waitingRoom)
	if waitingRoomDTOResponse.Errors != nil && waitingRoomDTOResponse.Errors[0].Status == 404 {
		data.SetId("")
	}
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to update Waiting Room %d for Site ID %s", waitingRoomID, siteID)
		return diags
	}

	diags = append(diags, resourceWaitingRoomRead(ctx, data, m)[:]...)

	return diags
}

func resourceWaitingRoomDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics
	siteID := data.Get("site_id").(string)

	waitingRoomID, err := strconv.ParseInt(data.Id(), 10, 64)
	if err != nil {
		log.Printf("[ERROR] The ID should be numeric. Currrent value: %s", data.Id())
		return []diag.Diagnostic{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Waiting Room ID conversion error",
			Detail:   fmt.Sprintf("The ID should be numeric. Currrent value: %s", data.Id()),
		}}
	}

	waitingRoomDTOResponse, diags := client.DeleteWaitingRoom(siteID, waitingRoomID)
	if waitingRoomDTOResponse.Errors != nil && waitingRoomDTOResponse.Errors[0].Status == 404 {
		data.SetId("")
	}
	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to update Waiting Room %d for Site ID %s", waitingRoomID, siteID)
		return diags
	}

	data.SetId("")

	return diags
}
