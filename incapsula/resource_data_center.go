package incapsula

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
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

				siteID, err := strconv.Atoi(idSlice[0])
				if err != nil {
					return nil, err
				}
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
				Description: "Is enabled",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "yes",
			},
			"is_standby": {
				Description: "Is standby",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_content": {
				Description: "Is content",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceDataCenterCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	dataCenterAddResponse, err := client.AddDataCenter(
		d.Get("site_id").(string),
		d.Get("name").(string),
		d.Get("server_address").(string),
		d.Get("is_standby").(string),
		d.Get("is_content").(string),
	)

	if err != nil {
		return err
	}

	// Set the dc ID
	d.SetId(dataCenterAddResponse.DataCenterID)

	return resourceDataCenterRead(d, m)
}

func resourceDataCenterRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListDataCentersResponse for the data center
	client := m.(*Client)

	listDataCentersResponse, err := client.ListDataCenters(d.Get("site_id").(string))
	if err != nil {
		return err
	}

	for _, dataCenter := range listDataCentersResponse.DCs {
		if dataCenter.Name == d.Get("name").(string) {
			d.Set("enabled", dataCenter.Enabled)
			d.Set("is_content", dataCenter.ContentOnly)
		}
	}

	return nil
}

func resourceDataCenterUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.EditDataCenter(
		d.Get("dc_id").(int),
		d.Get("name").(string),
		d.Get("is_standby").(string),
		d.Get("is_content").(string),
	)

	if err != nil {
		return err
	}

	return nil
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

	return nil
}
