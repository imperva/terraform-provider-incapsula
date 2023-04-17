package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func resourceATOSiteAllowlist() *schema.Resource {
	return &schema.Resource{
		Create: resourceATOSiteAllowlistUpdate,
		Read:   resourceATOSiteAllowlistRead,
		Update: resourceATOSiteAllowlistUpdate,
		Delete: resourceATOSiteAllowlistDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteId, err := strconv.Atoi(d.Id())
				err = d.Set("site_id", siteId)
				if err != nil {
					return nil, fmt.Errorf("failed to extract site ID from import command, actual value: %s, error : %s", d.Id(), err)
				}
				log.Printf("[DEBUG] Import ATO allowlist for site ID %d", siteId)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Account ID that the site belongs to.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"site_id": {
				Description: "Site ID to get the allowlist for.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"allowlist": {
				Description: "The allowlist of IPs and IP ranges for the given site ID",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					// Terraform does not allow us to granular type define the map
					Type: schema.TypeMap,
				},
			},
		},
	}
}

func resourceATOSiteAllowlistRead(d *schema.ResourceData, m interface{}) error {

	// Fetch our http client
	client := m.(*Client)

	// Fetch the ATO allowlist of IPs and subnets
	siteId, ok := d.Get("site_id").(int)

	if !ok {
		return fmt.Errorf("site_id should be of type int. Received : %s", d.Get("site_id"))
	}

	if siteId == 0 {
		siteIdFromResourceId, conversionError := strconv.Atoi(d.Id())
		if conversionError != nil {
			return fmt.Errorf("atleast one of id or site_id should be set for incapsula_ato_site_allowlist")
		}
		siteId = siteIdFromResourceId
	}

	accountId := d.Get("account_id").(int)
	atoAllowlistDTO, err := client.GetAtoSiteAllowlistWithRetries(accountId, siteId)

	// Handle fetch error
	if err != nil {
		return fmt.Errorf("[Error] getting ATO allowlist: %s", err)
	}

	// Initialize the data map that terraform uses
	var atoAllowlistMap = make(map[string]interface{})

	// Assign values from our received DTO to the map that terraform understands.
	atoAllowlistMap["site_id"] = atoAllowlistDTO.SiteId

	// Assign the allowlist if present to the terraform compatible map
	if atoAllowlistDTO.Allowlist != nil {

		atoAllowlistMap["allowlist"] = make([]map[string]interface{}, len(atoAllowlistDTO.Allowlist))

		for i, allowlistItem := range atoAllowlistDTO.Allowlist {
			allowlistItemMap := make(map[string]interface{})
			allowlistItemMap["ip"] = allowlistItem.Ip
			allowlistItemMap["mask"] = allowlistItem.Mask
			allowlistItemMap["desc"] = allowlistItem.Desc
			atoAllowlistMap["allowlist"].([]map[string]interface{})[i] = allowlistItemMap

		}

	} else {
		atoAllowlistMap["allowlist"] = make([]interface{}, 0)
	}

	d.SetId(strconv.Itoa(siteId))
	err = d.Set("allowlist", atoAllowlistMap["allowlist"])
	if err != nil {
		e := fmt.Errorf("[Error] Error in reading allowlist values : %s", err)
		return e
	}

	return nil
}

func resourceATOSiteAllowlistUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	accountId := d.Get("account_id").(int)

	atoAllowlistMap := make(map[string]interface{})
	atoAllowlistMap["account_id"] = accountId
	atoAllowlistMap["site_id"] = siteId
	atoAllowlistMap["allowlist"] = d.Get("allowlist").([]interface{})

	log.Printf("[DEBUG] Updating ATO site allowlist site ID %d \n", siteId)

	// convert terraform compatible map to ATOAllowlistDTO
	atoAllowlistDTO, err := formAtoAllowlistDTOFromMap(atoAllowlistMap)
	if err != nil {
		e := fmt.Errorf("[Error] Error forming ATO allow list object for API call : %s", err)
		log.Printf(e.Error())
		return err
	}

	err = client.UpdateATOSiteAllowlistWithRetries(atoAllowlistDTO)
	if err != nil {
		e := fmt.Errorf("[ERROR] Could not update ATO site allowlist for site ID : %d Error : %s \n", atoAllowlistDTO.SiteId, err)
		return e
	}

	return resourceATOSiteAllowlistRead(d, m)
}

func resourceATOSiteAllowlistDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	accountId := d.Get("account_id").(int)

	log.Printf("[DEBUG] Deleting ATO site allowlist for site ID %d \n", siteId)

	err := client.DeleteATOSiteAllowlist(accountId, siteId)
	if err != nil {
		e := fmt.Errorf("[ERROR] Could not delete ATO site allowlist for site ID : %d Error : %s \n", siteId, err)
		return e
	}

	return nil
}
