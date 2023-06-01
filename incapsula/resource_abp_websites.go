package incapsula

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AbpTerraformAccount struct {
	AutoPublish   bool                       `json:"auto_publish"`
	WebsiteGroups []AbpTerraformWebsiteGroup `json:"website_groups"`
}

type AbpTerraformWebsiteGroup struct {
	Id       *string `json:"id"`
	NameId   *string
	Name     string                `json:"name"`
	Websites []AbpTerraformWebsite `json:"websites"`
}

type AbpTerraformWebsite struct {
	Id               *string `json:"id"`
	WebsiteId        int     `json:"website_id"`
	EnableMitigation bool    `json:"enable_mitigation"`
}

func (a *AbpTerraformWebsiteGroup) UniqueId() string {
	if a.NameId == nil {
		return a.Name
	}
	return *a.NameId
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
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unique user-defined identifier used to differentiate websites groups whose `name` are identical",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the website group. Must be unique unless `name_id` is specified.",
						},
						"website": {
							Type: schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"website_id": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"enable_mitigation": {
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

	// Map the unique names to
	nameToId := make(map[string]string)
	oldWebsiteGroup, newWebsiteGroup := data.GetChange("website_group")
	for _, websiteGroup := range oldWebsiteGroup.([]interface{}) {
		websiteGroup := websiteGroup.(map[string]interface{})

		nameId := websiteGroup["name_id"].(string)
		if nameId == "" {
			nameId = websiteGroup["name"].(string)
		}
		nameToId[nameId] = websiteGroup["id"].(string)
	}

	autoPublish := data.Get("auto_publish").(bool)

	var websiteGroups []AbpTerraformWebsiteGroup
	for _, websiteGroup := range newWebsiteGroup.([]interface{}) {
		websiteGroup := websiteGroup.(map[string]interface{})

		var websites []AbpTerraformWebsite
		for _, website := range websiteGroup["website"].([]interface{}) {
			website := website.(map[string]interface{})

			id := website["id"].(string)
			var idOpt *string
			if id != "" {
				idOpt = &id
			}
			websites = append(websites,
				AbpTerraformWebsite{
					Id:               idOpt,
					WebsiteId:        website["website_id"].(int),
					EnableMitigation: website["enable_mitigation"].(bool),
				})
		}

		name := websiteGroup["name"].(string)

		nameId := websiteGroup["name_id"].(string)

		var nameIdOpt *string
		if nameId != "" {
			nameIdOpt = &nameId
		}

		// The items in the website group list may have shifted around making the ids not match name/name_id. Thus we use name/name_id
		// to lookup the id of the old configuration.
		var idOpt *string
		if nameIdOpt == nil {
			id, ok := nameToId[name]
			if ok {
				idOpt = &id
			}
		} else {
			id, ok := nameToId[*nameIdOpt]
			if ok {
				idOpt = &id
			}
		}

		websiteGroups = append(websiteGroups, AbpTerraformWebsiteGroup{
			Id:       idOpt,
			NameId:   nameIdOpt,
			Name:     name,
			Websites: websites,
		})
	}

	return AbpTerraformAccount{
		AutoPublish:   autoPublish,
		WebsiteGroups: websiteGroups,
	}
}

func serializeAccount(data *schema.ResourceData, account AbpTerraformAccount) {

	// We never store this on the server side, just in the terraform state so ignore what the server sends
	// data.Set("auto_publish", account.AutoPublish)

	websiteGroupsData := make([]interface{}, len(account.WebsiteGroups), len(account.WebsiteGroups))
	for i, websiteGroup := range account.WebsiteGroups {
		websiteGroupData := make(map[string]interface{})

		websitesData := make([]interface{}, len(websiteGroup.Websites), len(websiteGroup.Websites))
		for j, website := range websiteGroup.Websites {
			websiteData := make(map[string]interface{})

			if website.Id != nil {
				websiteData["id"] = *website.Id
			}
			websiteData["website_id"] = website.WebsiteId
			websiteData["enable_mitigation"] = website.EnableMitigation

			websitesData[j] = websiteData
		}

		if websiteGroup.Id != nil {
			websiteGroupData["id"] = *websiteGroup.Id
		}
		websiteGroupData["name"] = websiteGroup.Name
		websiteGroupData["website"] = websitesData

		websiteGroupsData[i] = websiteGroupData
	}

	data.Set("website_group", websiteGroupsData)
}

func resourceAbpWebsitesCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)
	account := extractAccount(data)
	var abpWebsites *AbpTerraformAccount

	abpWebsites, diags = client.CreateAbpWebsites(strconv.Itoa(accountId), account)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to create ABP websites for Account ID %d", accountId)
		return diags
	}

	serializeAccount(data, *abpWebsites)

	data.SetId(strconv.Itoa(accountId))

	return diags
}

func resourceAbpWebsitesRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)

	var abpWebsites *AbpTerraformAccount
	abpWebsites, diags = client.ReadAbpWebsites(strconv.Itoa(accountId))

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to read ABP websites for Account ID %d", accountId)
		return diags
	}

	serializeAccount(data, *abpWebsites)

	return diags
}

func resourceAbpWebsitesUpdate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)
	account := extractAccount(data)

	var abpWebsites *AbpTerraformAccount
	abpWebsites, diags = client.UpdateAbpWebsites(strconv.Itoa(accountId), account)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to update ABP websites for Account ID %d", accountId)
		return diags
	}

	serializeAccount(data, *abpWebsites)

	return diags
}

func resourceAbpWebsitesDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)
	autoPublish := data.Get("auto_publish").(bool)

	_, diags = client.DeleteAbpWebsites(strconv.Itoa(accountId), autoPublish)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to delete ABP websites for Account ID %d", accountId)
		return diags
	}

	data.SetId("")

	return diags
}
