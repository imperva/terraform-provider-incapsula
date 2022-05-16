package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceOriginPOP() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource is deprecated. It will be removed in a future version. Please use resource incapsula_data_centers_configuration instead.",
		Create:             resourceOriginPOPUpdate,
		Read:               resourceOriginPOPRead,
		Update:             resourceOriginPOPUpdate,
		Delete:             resourceOriginPOPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"dc_id": {
				Description: "Numeric identifier of the data center.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"site_id": {
				Description: "Numeric identifier of the site.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"origin_pop": {
				Description: "The Origin POP code (must be lowercase), e.g: iad. Note, this field is create/update only. Reads are not supported as the API doesn't exist yet. Note that drift may happen.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					// Check if valid JSON
					d := val.(string)
					if strings.ToLower(d) != d {
						errs = append(errs, fmt.Errorf("%q must be lowercase, please check your origin POP code, got: %s", key, d))
					}
					return
				},
			},
		},
	}
}

func resourceOriginPOPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	dcID := d.Get("dc_id").(int)
	siteID := d.Get("site_id").(int)
	originPOP := d.Get("origin_pop").(string)

	log.Printf("[INFO] Setting Incapsula origin POP: %s for data center: %d\n", originPOP, dcID)

	err := client.SetOriginPOP(dcID, originPOP)

	if err != nil {
		log.Printf("[ERROR] Could not set Incapsula origin POP: %s for data center: %d: %s\n", originPOP, dcID, err)
		return err
	}

	log.Printf("[INFO] Set Incapsula origin POP: %s for data center: %d\n", originPOP, dcID)

	syntheticID := fmt.Sprintf("%d/%d", siteID, dcID)
	d.SetId(syntheticID)

	return resourceOriginPOPRead(d, m)
}

func resourceOriginPOPRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListDataCentersResponse for the data centers
	client := m.(*Client)
	if !strings.Contains(d.Id(), "/") {
		log.Printf("[ERROR] The origin_pop resource in your state file is in the old, unsupported format. /n" +
			"We recommend to use the new resource of data_center_configuration which replaced this resource./n" +
			"If you choose to continue using this resource, you must update the resource in your configuration files according to the new format. /n " +
			"The old resource will now be removed from state file, and will be updated back on the next terraform plan.")
		d.SetId("")
		return nil
	}
	siteID := strings.Split(d.Id(), "/")[0]
	dcID := strings.Split(d.Id(), "/")[1]

	listDataCentersResponse, err := client.ListDataCenters(siteID)

	if err != nil {
		log.Printf("[ERROR] Could not read origin POP for data center: %s, site: %s %s\n", dcID, siteID, err)
		return err
	}

	found := false

	for _, dataCenter := range listDataCentersResponse.DCs {
		if dataCenter.ID == dcID {
			originPop := dataCenter.OriginPop
			if originPop != "" {
				siteIDInteger, _ := strconv.Atoi(siteID)
				dcIDInteger, _ := strconv.Atoi(dcID)
				d.Set("site_id", siteIDInteger)
				d.Set("dc_id", dcIDInteger)
				d.Set("origin_pop", originPop)
				syntheticID := fmt.Sprintf("%s/%s", siteID, dcID)
				d.SetId(syntheticID)
				found = true
				break
			}
		}
	}

	if !found {
		d.SetId("")
	}

	return nil
}

func resourceOriginPOPDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	dcID := d.Get("dc_id").(int)
	err := client.SetOriginPOP(dcID, "")

	if err != nil {
		log.Printf("[ERROR] Could not delete Incapsula origin POP for data center: %d: %s\n", dcID, err)
		return err
	}

	log.Printf("[INFO] Deleted Incapsula origin POP for data center: %d\n", dcID)

	return nil
}
