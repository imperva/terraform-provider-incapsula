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
				siteId, err := strconv.Atoi(d.Get("site_id").(string))
				if err != nil {
					return nil, fmt.Errorf("site_id should be of type int")
				}
				d.SetId(strconv.Itoa(siteId))
				log.Printf("[DEBUG] Import ATO allowlist for site ID %d", siteId)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Site ID to get the allowlist for.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			// Computed attributes
			"allowlist": {
				Description: "The allowlist of IPs and IP ranges for the given site ID",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP excluded from mitigation",
						},
						"mask": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet excluded for the IP from mitigation",
						},
						"desc": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notes specified at the time of creating this allowlist entry",
						},
						"updated": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "IP excluded from mitigation",
						},
					},
				},
			},
		},
	}
}

func resourceATOSiteAllowlistRead(d *schema.ResourceData, m interface{}) error {

	// Fetch our http client
	client := m.(*Client)

	// Fetch the ATO allowlist of IPs and subnets
	siteId := d.Get("site_id").(int)
	atoAllowlistDTO, err := client.GetAtoSiteAllowlist(siteId)

	// Handle fetch error
	if err != nil {
		return fmt.Errorf("[Error] getting ATO allowlist: %s", err)
	}

	// Initialize the data map that terraform uses
	var atoAllowlistMap = make(map[string]interface{})

	/* Assign values from our received DTO to the map that terraform understands.
	This is defined in the schema at dataSourceATOAllowlist() */
	atoAllowlistMap["site_id"] = atoAllowlistDTO.siteId

	// Assign the allowlist if present to the terraform compatible map
	if atoAllowlistDTO.allowlist != nil {

		atoAllowlistMap["allowlist"] = make([]map[string]interface{}, len(atoAllowlistDTO.allowlist))

		for i, allowlistItem := range atoAllowlistDTO.allowlist {
			allowlistItemMap := make(map[string]interface{})
			allowlistItemMap["ip"] = allowlistItem.ip
			allowlistItemMap["mask"] = allowlistItem.mask
			allowlistItemMap["desc"] = allowlistItem.desc
			allowlistItemMap["updated"] = allowlistItem.updated
			atoAllowlistMap["allowlist"].([]map[string]interface{})[i] = allowlistItemMap

		}

	} else {
		atoAllowlistMap["allowlist"] = make([]interface{}, 0)
	}

	d.SetId(strconv.Itoa(siteId))
	err = d.Set("allowlist", atoAllowlistMap)
	if err != nil {
		e := fmt.Errorf("[Error] Error in reading allowlist values : %s", err)
		return e
	}

	return nil
}

func resourceATOSiteAllowlistUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)
	atoAllowlistMap := d.Get("map").(map[string]interface{})

	log.Printf("[DEBUG] Updating ATO site allowlist site ID %d \n", siteId)

	// Assign the allowlist if present to the terraform compatible map
	if atoAllowlistMap["allowlist"] != nil {
		atoAllowlistMap := d.Get("allowlist")
		atoAllowlistDTO, err := formAtoAllowlistDTOFromMap(atoAllowlistMap.(map[string]interface{}))
		if err != nil {
			e := fmt.Errorf("[Error] Error forming ATO allow list object for API call : %s", err)
			log.Printf(e.Error())
			return err
		}

		err = client.UpdateATOSiteAllowlist(siteId, atoAllowlistDTO)
		if err != nil {
			e := fmt.Errorf("[ERROR] Could not update ATO site allowlist for site ID : %d Error : %s \n", atoAllowlistDTO.siteId, err)
			return e
		}
	} else {
		// No update required as nil value implies that this resource is not managed by terraform
	}

	return resourceATOSiteAllowlistRead(d, m)
}

func resourceATOSiteAllowlistDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(int)

	log.Printf("[DEBUG] Deleting ATO site allowlist site ID %d \n", siteId)

	err := client.DeleteATOSiteAllowlist(siteId)
	if err != nil {
		e := fmt.Errorf("[ERROR] Could not delete ATO site allowlist for site ID : %d Error : %s \n", siteId, err)
		return e
	}

	return resourceATOSiteAllowlistRead(d, m)
}
