package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"strconv"
	"strings"
)

func resourceDataCenterServers() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataCenterServersCreate,
		Read:   resourceDataCenterServersRead,
		Update: resourceDataCenterServersUpdate,
		Delete: resourceDataCenterServersDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected dc_id/server_id", d.Id())
				}

				dcID, err := strconv.Atoi(idSlice[0])
				serverID := idSlice[1]
				if err != nil {
					return nil, err
				}

				d.Set("dc_id", dcID)
				d.SetId(serverID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"dc_id": {
				Description: "Numeric identifier of the data center server to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			// Optional Arguments
			"server_address": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_standby": {
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceDataCenterServersCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	dataCenterServerAddResponse, err := client.AddDataCenterServers(
		d.Get("dc_id").(int),
		d.Get("server_address").(string),
		d.Get("is_standby").(string),
	)

	if err != nil {
		return err
	}

	// Set the dc ID
	d.SetId(d.Get("dc_id").(string))
	d.SetId(strconv.Itoa(dataCenterServerAddResponse.ServerID))

	return resourceDataCenterServersRead(d, m)
}

func resourceDataCenterServersRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	siteStatusResponse, err := client.ListDataCenterServers(d.Get("site_id").(int))
	d.Set("todo", siteStatusResponse)

	if err != nil {
		return err
	}

	// todo: review response

	return nil
}

func resourceDataCenterServersUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.EditDataCenterServers(
		d.Id(),
		d.Get("server_address").(string),
		d.Get("is_standby").(string),
		d.Get("is_content").(string),
	)

	if err != nil {
		return err
	}

	return nil
}

func resourceDataCenterServersDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	serverID, _ := strconv.Atoi(d.Id())
	err := client.DeleteDataCenterServers(serverID)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
