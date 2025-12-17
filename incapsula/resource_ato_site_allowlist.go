package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func resourceATOSiteAllowlist() *schema.Resource {
	return &schema.Resource{
		Create: resourceATOSiteAllowlistUpdate,
		Read:   resourceATOSiteAllowlistRead,
		Update: resourceATOSiteAllowlistUpdate,
		Delete: resourceATOSiteAllowlistDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {

				// If this id is of form <account_id>/<site_id> extract the sub account Id as well
				if strings.Contains(d.Id(), "/") {
					keyParts := strings.Split(d.Id(), "/")
					print("id is %s", d.Id())
					if len(keyParts) != 2 {
						return nil, fmt.Errorf("Error parsing ID, actual value: %s, expected 2 numeric IDs seperated by '/'\n", d.Id())
					}
					accountId, err := strconv.Atoi(keyParts[0])
					if err != nil {
						return nil, fmt.Errorf("[ERROR] failed to convert account ID from import command, actual value: %s, expected numeric id", keyParts[0])
					}
					siteId, err := strconv.Atoi(keyParts[1])
					if err != nil {
						return nil, fmt.Errorf("[ERROR] failed to convert site ID from import command, actual value: %s, expected numeric id", keyParts[1])
					}

					d.Set("account_id", accountId)
					d.Set("site_id", siteId)
					d.Set("id", d.Id())
					log.Printf("[DEBUG] To Import ATO allowlsit configuration for account ID %d , site ID %d", accountId, siteId)
					return []*schema.ResourceData{d}, nil
				}

				siteId, err := strconv.Atoi(d.Id())
				err = d.Set("site_id", siteId)
				if err != nil {
					return nil, fmt.Errorf("[ERROR] failed to extract site ID from import command, actual value: %s, error : %s", d.Id(), err)
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
	atoAllowlistDTO, status, err := client.GetAtoSiteAllowlistWithRetries(accountId, siteId)

	// Handle fetch error
	if err != nil {
		return fmt.Errorf("[Error] getting ATO allowlist: %s", err)
	}

	// Initialize the data map that terraform uses
	var atoAllowlistEntry = make(map[string]interface{})

	// Assign values from our received DTO to the map that terraform understands.
	atoAllowlistEntry["site_id"] = atoAllowlistDTO.SiteId

	// Assign the allowlist if present to the terraform compatible map
	if atoAllowlistDTO.Allowlist != nil {

		atoAllowlistEntry["allowlist"] = make([]map[string]interface{}, len(atoAllowlistDTO.Allowlist))

		for i, allowlistItem := range atoAllowlistDTO.Allowlist {
			allowlistItemMap := make(map[string]interface{})
			allowlistItemMap["ip"] = allowlistItem.Ip
			allowlistItemMap["mask"] = allowlistItem.Mask
			allowlistItemMap["desc"] = allowlistItem.Desc
			atoAllowlistEntry["allowlist"].([]map[string]interface{})[i] = allowlistItemMap

		}

	} else {
		atoAllowlistEntry["allowlist"] = make([]interface{}, 0)
	}

	// Handle site does not exist in ATO. If this is a permissions issue, then let this be addressed in the update phase
	if status == http.StatusUnauthorized {
		// Remove this resource from the state file by setting empty ID as it does not exist. Terraform will remove it
		d.Set("id", "")
	} else {
		// Set ID for all other cases as siteId
		if accountId == 0 {
			d.SetId(strconv.Itoa(siteId))
		} else {
			d.SetId(fmt.Sprintf("%d/%d", accountId, siteId))
		}
	}

	err = d.Set("allowlist", atoAllowlistEntry["allowlist"])
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
		log.Print(e.Error())
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
