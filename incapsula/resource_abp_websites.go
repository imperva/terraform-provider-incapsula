package incapsula

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type AbpTerraformAccount struct {
	AutoPublish   bool                       `json:"auto_publish"`
	LastPublish   string                     `json:"last_publish"`
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
	IncapsulaSiteId  int     `json:"incapsula_site_id"`
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

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				accountIdStr := d.Id()
				accountId, err := strconv.Atoi(accountIdStr)

				if err != nil {
					return nil, fmt.Errorf("Expected an account Id which must be an integer: %s", err.Error())
				}

				d.Set("account_id", accountId)
				d.SetId(accountIdStr)

				client := meta.(*Client)
				abpWebsites, diags := client.ReadAbpWebsites(accountId)

				for _, diagnostic := range diags {
					if diagnostic.Severity == diag.Error {
						return nil, fmt.Errorf("%s", diagnostic.Detail)
					}
				}

				setUniqueNameIds(abpWebsites)

				serializeAccount(d, *abpWebsites)

				log.Printf("[DEBUG] Import ABP websites for account id %d", accountId)

				return []*schema.ResourceData{d}, nil
			},
		},

		Description: "Provides an Incapsula ABP (Advanced Bot Protection) website resource. Allows you to enable and configure ABP for given websites.\n" +
			"\n" +
			"NOTE: Due to limitations in ABP, the API key/id used to deploy this resource must match the `account_id` used in the resource (API key/id for a parent account do not work). All Incapsula sites associated with the resource must also be defined in that account.",

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "The account these websites belongs to.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"auto_publish": {
				Description: "Whether to publish the changes automatically. Changes don't take take effect until they have been published.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"last_publish": {
				Description: "When the last publish was done for this terraform resource. Changes are published when `auto_publish` is true and the terraform config is applied.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"website_group": {
				Description: "List of website groups which are associated to ABP. Website groups are matched in a top-down fashion. If a more specific website group should take precedence over over a wild card entry then that should be higher in the list.",
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
							Description: "Unique user-defined identifier used to differentiate website groups whose `name` is identical",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name for the website group. Must be unique unless `name_id` is specified.",
						},
						"website": {
							Description: "List of websites within the website group. Websites are matched in a top-down fashion. If a more specific website should take precedence over over a wild card entry then that should be higher in the list",
							Type:        schema.TypeList,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"incapsula_site_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Which `incapsula_site` this website refers to",
									},
									"enable_mitigation": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Enables the ABP conditions for this website. Defaults to true.",
									},
								},
							},
							Optional: true,
						},
					},
				},
				Optional: true,
			},
		},
	}
}

func extractAccount(data *schema.ResourceData) (AbpTerraformAccount, diag.Diagnostics) {

	oldWebsiteGroup, newWebsiteGroup := data.GetChange("website_group")

	usedNames := make(map[string]bool)
	var diags diag.Diagnostics
	for _, websiteGroup := range newWebsiteGroup.([]interface{}) {
		websiteGroup := websiteGroup.(map[string]interface{})

		nameId := websiteGroup["name_id"].(string)
		if nameId == "" {
			nameId = websiteGroup["name"].(string)
		}
		_, ok := usedNames[nameId]
		if ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Found duplicate identifier (%s) for website group", nameId),
				Detail:   fmt.Sprintf("Identifiers must be unique per website group. If you need duplicate `name`s you may specify `name_id` with an unique identifier"),
			})
		}
		usedNames[nameId] = true
	}

	nameToId := make(map[string]string)
	siteIdToId := make(map[int]string)
	for _, websiteGroup := range oldWebsiteGroup.([]interface{}) {
		websiteGroup := websiteGroup.(map[string]interface{})

		nameId := websiteGroup["name_id"].(string)
		if nameId == "" {
			nameId = websiteGroup["name"].(string)
		}
		nameToId[nameId] = websiteGroup["id"].(string)

		for _, website := range websiteGroup["website"].([]interface{}) {
			website := website.(map[string]interface{})

			incapsulaSiteId := website["incapsula_site_id"].(int)

			if _, ok := siteIdToId[incapsulaSiteId]; ok {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Found duplicate incapsula_site_id (%d) for website", incapsulaSiteId),
					Detail:   fmt.Sprintf("Each incapsula_site can only be used by a single website"),
				})
			} else {
				siteIdToId[incapsulaSiteId] = website["id"].(string)
			}
		}
	}

	if len(diags) > 0 {
		return AbpTerraformAccount{}, diags
	}

	autoPublish := data.Get("auto_publish").(bool)

	websiteGroups := make([]AbpTerraformWebsiteGroup, 0, 0)
	for _, websiteGroup := range newWebsiteGroup.([]interface{}) {
		websiteGroup := websiteGroup.(map[string]interface{})

		websites := make([]AbpTerraformWebsite, 0)
		for _, website := range websiteGroup["website"].([]interface{}) {
			website := website.(map[string]interface{})

			incapsulaSiteId := website["incapsula_site_id"].(int)

			var idOpt *string
			if id, ok := siteIdToId[incapsulaSiteId]; ok {
				idOpt = &id
			}

			websites = append(websites,
				AbpTerraformWebsite{
					Id:               idOpt,
					IncapsulaSiteId:  incapsulaSiteId,
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
			if id, ok := nameToId[name]; ok {
				idOpt = &id
			}
		} else {
			if id, ok := nameToId[*nameIdOpt]; ok {
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
	}, nil
}

func setUniqueNameIds(account *AbpTerraformAccount) {
	nameIds := make(map[string]bool, 0)
	for _, websiteGroup := range account.WebsiteGroups {
		nameIds[websiteGroup.Name] = true
	}

	names := make(map[string]bool, 0)
	for i, websiteGroup := range account.WebsiteGroups {
		if _, ok := names[websiteGroup.Name]; ok {
			nameId := websiteGroup.Name + "-0"
			for i := 0; i < 1000000; i++ {
				if _, ok := nameIds[nameId]; ok {
					nameId = websiteGroup.Name + "-" + strconv.Itoa(i)
				} else {
					nameIds[nameId] = true
					break
				}
			}

			account.WebsiteGroups[i].NameId = &nameId
		}
		names[websiteGroup.Name] = true
	}
}

func serializeAccount(data *schema.ResourceData, account AbpTerraformAccount) {

	websiteGroupsData := make([]interface{}, len(account.WebsiteGroups), len(account.WebsiteGroups))
	oldWebsiteGroups := data.Get("website_group").([]interface{})
	for i, websiteGroup := range account.WebsiteGroups {
		websiteGroupData := make(map[string]interface{})

		websitesData := make([]interface{}, len(websiteGroup.Websites), len(websiteGroup.Websites))
		for j, website := range websiteGroup.Websites {
			websiteData := make(map[string]interface{})

			if website.Id != nil {
				websiteData["id"] = *website.Id
			}
			websiteData["incapsula_site_id"] = website.IncapsulaSiteId
			websiteData["enable_mitigation"] = website.EnableMitigation

			websitesData[j] = websiteData
		}

		if websiteGroup.Id != nil {
			websiteGroupData["id"] = *websiteGroup.Id
		}
		websiteGroupData["name"] = websiteGroup.Name
		websiteGroupData["website"] = websitesData

		if i < len(oldWebsiteGroups) {
			oldWebsiteGroup := oldWebsiteGroups[i].(map[string]interface{})
			// Don't lose the name_id that might have been set by the user
			websiteGroupData["name_id"] = oldWebsiteGroup["name_id"]
		} else if websiteGroup.NameId != nil {
			websiteGroupData["name_id"] = *websiteGroup.NameId
		}

		websiteGroupsData[i] = websiteGroupData
	}

	if account.LastPublish != "" {
		data.Set("last_publish", account.LastPublish)
	}

	data.Set("website_group", websiteGroupsData)
}

func resourceAbpWebsitesCreate(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	var diags diag.Diagnostics

	accountId := data.Get("account_id").(int)
	account, diags := extractAccount(data)
	if diags != nil && diags.HasError() {
		return diags
	}
	var abpWebsites *AbpTerraformAccount

	abpWebsites, diags = client.CreateAbpWebsites(accountId, account)

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
	abpWebsites, diags = client.ReadAbpWebsites(accountId)

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
	account, diags := extractAccount(data)
	if diags != nil && diags.HasError() {
		return diags
	}

	var abpWebsites *AbpTerraformAccount
	abpWebsites, diags = client.UpdateAbpWebsites(accountId, account)

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

	_, diags = client.DeleteAbpWebsites(accountId, autoPublish)

	if diags != nil && diags.HasError() {
		log.Printf("[ERROR] Failed to delete ABP websites for Account ID %d", accountId)
		return diags
	}

	data.SetId("")

	return diags
}
