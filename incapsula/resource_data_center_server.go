package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDataCenterServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataCenterServerCreate,
		Read:   resourceDataCenterServerRead,
		Update: resourceDataCenterServerUpdate,
		Delete: resourceDataCenterServerDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected dc_id/server_id", d.Id())
				}

				dcID, err := strconv.Atoi(idSlice[0])
				if err != nil {
					return nil, err
				}
				serverID := idSlice[1]

				d.Set("dc_id", dcID)
				d.SetId(serverID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"dc_id": {
				Description: "Numeric identifier of the data center server to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"server_address": {
				Description: "The server's address.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_standby": {
				Description: "Set the server as Active (P0) or Standby (P1).",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_enabled": {
				Description: "Enables the data center server.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "yes",
			},
		},
	}
}

func resourceDataCenterServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	dataCenterServerAddResponse, err := client.AddDataCenterServer(
		d.Get("dc_id").(string),
		d.Get("server_address").(string),
		d.Get("is_standby").(string),
	)

	if err != nil {
		return err
	}

	if d.Get("is_enabled") != "" {
		log.Printf("[INFO] Updating data center server server_id (%s) with is_enabled (%s)\n", dataCenterServerAddResponse.ServerID, d.Get("is_enabled").(string))
		_, err := client.EditDataCenterServer(dataCenterServerAddResponse.ServerID, d.Get("server_address").(string), d.Get("is_standby").(string), d.Get("is_enabled").(string))
		if err != nil {
			log.Printf("[ERROR] Could not update data center server server_id (%s) with is_enabled (%s) %s\n", dataCenterServerAddResponse.ServerID, d.Get("is_enabled").(string), err)
			return err
		}
	}

	// Set the server ID
	d.SetId(dataCenterServerAddResponse.ServerID)

	return resourceDataCenterServerRead(d, m)
}

func resourceDataCenterServerRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListDataCentersResponse for the data centers
	client := m.(*Client)

	listDataCentersResponse, err := client.ListDataCenters(d.Get("site_id").(string))
	if err != nil {
		return err
	}

	for _, dataCenter := range listDataCentersResponse.DCs {
		if dataCenter.ID == d.Get("dc_id").(string) {
			for _, server := range dataCenter.Servers {
				if server.Address == d.Get("server_address").(string) {
					d.Set("enabled", server.Enabled)
					d.Set("server_address", server.Address)
					d.Set("is_standby", server.IsStandBy)
				}
			}
		}
	}

	return nil
}

func resourceDataCenterServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.EditDataCenterServer(
		d.Id(),
		d.Get("server_address").(string),
		d.Get("is_standby").(string),
		d.Get("is_enabled").(string),
	)

	if err != nil {
		return err
	}

	return nil
}

func resourceDataCenterServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	serverID := d.Id()
	err := client.DeleteDataCenterServer(serverID)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
