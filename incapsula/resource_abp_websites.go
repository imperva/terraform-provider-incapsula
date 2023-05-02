package incapsula

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AbpTerraformAccount struct {
	AccountId     int                        `json:"account_id"`
	AutoPublish   bool                       `json:"auto_publish"`
	WebsiteGroups []AbpTerraformWebsiteGroup `json:"website_groups"`
}

type AbpTerraformWebsiteGroup struct {
	Id       *string               `json:"id"`
	Name     string                `json:"name"`
	Websites []AbpTerraformWebsite `json:"websites"`
}

type AbpTerraformWebsite struct {
	Id                *string `json:"id"`
	WebsiteId         int     `json:"website_id"`
	MitigationEnabled bool    `json:"mitigation_enabled"`
}

func resourceAbpWebsites() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAbpWebsitesCreate,
		ReadContext:   resourceAbpWebsitesRead,
		UpdateContext: resourceAbpWebsitesUpdate,
		DeleteContext: resourceAbpWebsitesDelete,

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account these websites belongs to.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"auto_publish": {
				Description: "Whether to publish the changes automatically.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"website_group": {
				Description: "Whether to publish the changes automatically.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"website": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"website_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"mitigation_enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
								},
							},
							Required: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func extractAccount(data *schema.ResourceData) AbpTerraformAccount {

	accountId := data.Get("account_id").(int)
	autoPublish := data.Get("auto_publish").(bool)

	var websiteGroups []AbpTerraformWebsiteGroup
	for _, websiteGroup := range data.Get("website_group").([]interface{}) {
		websiteGroup := websiteGroup.(map[string]interface{})

		var websites []AbpTerraformWebsite
		for _, website := range websiteGroup["website"].([]interface{}) {
			website := website.(map[string]interface{})

			websites = append(websites,
				AbpTerraformWebsite{
					Id:                nil,
					WebsiteId:         website["website_id"].(int),
					MitigationEnabled: website["mitigation_enabled"].(bool),
				})
		}

		websiteGroups = append(websiteGroups, AbpTerraformWebsiteGroup{
			Id:       nil,
			Name:     websiteGroup["name"].(string),
			Websites: websites,
		})
	}

	return AbpTerraformAccount{
		AccountId:     accountId,
		AutoPublish:   autoPublish,
		WebsiteGroups: websiteGroups,
	}
}

func resourceAbpWebsitesCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	account := extractAccount(data)
	_, diags = client.CreateAbpWebsites(strconv.Itoa(account.AccountId), account)

	return diags
}

func resourceAbpWebsitesRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)

	_, diags = client.ReadAbpWebsites(strconv.Itoa(accountId))

	return diags
}

func resourceAbpWebsitesUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	account := extractAccount(data)

	_, diags = client.UpdateAbpWebsites(strconv.Itoa(account.AccountId), account)
	return diags
}

func resourceAbpWebsitesDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)

	_, diags = client.DeleteAbpWebsites(strconv.Itoa(accountId))

	data.SetId("")
	return diags
}
