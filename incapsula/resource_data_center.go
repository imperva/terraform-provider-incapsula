package incapsula

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDataCenter() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataCenterCreate,
		Read:   resourceDataCenterRead,
		Update: resourceDataCenterUpdate,
		Delete: resourceDataCenterDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/dc_id", d.Id())
				}

				siteID := idSlice[0]
				dcID := idSlice[1]

				d.Set("site_id", siteID)
				d.SetId(dcID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The new data center's name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"server_address": {
				Description: "The server's address. Possible values: IP, CNAME.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"is_enabled": {
				Description: "Enables the data center.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "true",
			},
			"is_content": {
				Description: "The data center will be available for specific resources (Forward Delivery Rules).",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceDataCenterCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	dataCenterAddResponse, err := client.AddDataCenter(
		d.Get("site_id").(string),
		d.Get("name").(string),
		d.Get("server_address").(string),
		d.Get("is_content").(string),
	)

	if err != nil {
		return err
	}

	if d.Get("is_enabled") != "" {
		if d.Get("is_enabled") != "" {
			log.Printf("[INFO] Updating data center datacenter_id (%s) with is_enabled (%s)\n", dataCenterAddResponse.DataCenterID, d.Get("is_enabled").(string))
		}
		_, err := client.EditDataCenter(dataCenterAddResponse.DataCenterID, d.Get("name").(string), d.Get("is_content").(string), d.Get("is_enabled").(string))
		if err != nil {
			log.Printf("[ERROR] Could not update data center datacenter_id (%s) with is_enabled (%s) %s\n", dataCenterAddResponse.DataCenterID, d.Get("is_enabled").(string), err)
			return err
		}
	}

	// Set the dc ID
	d.SetId(dataCenterAddResponse.DataCenterID)

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	return resourceDataCenterRead(d, m)
}

func resourceDataCenterRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListDataCentersResponse for the data center
	client := m.(*Client)

	listDataCentersResponse, err := client.ListDataCenters(d.Get("site_id").(string))

	// List data centers response object may indicate that the Site ID has been deleted (9413)
	if listDataCentersResponse != nil {
		// Res can oscillate between strings and ints
		var resString string
		if resNumber, ok := listDataCentersResponse.Res.(float64); ok {
			resString = fmt.Sprintf("%d", int(resNumber))
		} else {
			resString = listDataCentersResponse.Res.(string)
		}
		if resString == "9413" {
			log.Printf("[INFO] Incapsula Site ID %s has already been deleted: %s\n", d.Get("site_id"), err)
			d.SetId("")
			return nil
		}
	}

	if err != nil {
		return err
	}

	found := false

	for _, dataCenter := range listDataCentersResponse.DCs {
		if dataCenter.ID == d.Id() {
			d.Set("name", dataCenter.Name)
			d.Set("is_enabled", dataCenter.Enabled)
			d.Set("is_content", dataCenter.ContentOnly)
			// Server address is the first value in the nested servers object
			d.Set("server_address", dataCenter.Servers[0].Address)
			found = true
		}
	}

	if !found {
		log.Printf("[INFO] Incapsula Data Center ID %s for Site ID %s has already been deleted: %s\n", d.Id(), d.Get("site_id"), err)
		d.SetId("")
		return nil
	}

	return nil
}

func resourceDataCenterUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err := client.EditDataCenter(
			d.Id(),
			d.Get("name").(string),
			d.Get("is_content").(string),
			d.Get("is_enabled").(string),
		)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error updating datacenter: %s", err))
		}

		return resource.NonRetryableError(resourceDataCenterRead(d, m))
	})
}

func resourceDataCenterDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := client.DeleteDataCenter(d.Id())

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	// There may be a timing/race condition here
	// Set an arbitrary period to sleep
	time.Sleep(3 * time.Second)

	return nil
}
