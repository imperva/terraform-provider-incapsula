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
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/rule_id", d.Id())
				}

				siteID, err := strconv.Atoi(idSlice[0])
				ruleID := idSlice[1]
				if err != nil {
					return nil, err
				}

				d.Set("site_id", siteID)
				d.Set("rule_id", ruleID)
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": &schema.Schema{
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"name": &schema.Schema{
				Description: "todo",
				Type:        schema.TypeString,
				Required:    true,
			},
			"server_address": &schema.Schema{
				Description: "todo",
				Type:        schema.TypeString,
				Required:    true,
			},

			// Optional Arguments
			"is_standby": &schema.Schema{
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"is_content": &schema.Schema{
				Description: "todo",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceDataCenterCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.AddDataCenter(
		d.Get("site_id").(int),
		d.Get("name").(string),
		d.Get("server_address").(string),
		d.Get("is_standby").(string),
		d.Get("is_content").(string),
	)

	if err != nil {
		return err
	}

	// Set the rule ID
	d.SetId(d.Get("dc_id").(string))

	return resourceDataCenterRead(d, m)
}

func resourceDataCenterRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the SiteResponse for the site
	client := m.(*Client)

	siteStatusResponse, err := client.ListDataCenters(d.Get("site_id").(int))
	d.Set("todo", siteStatusResponse)

	if err != nil {
		return err
	}

	// todo: review response

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

	err := client.DeleteDataCenter(d.Get("dc_id").(int))

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
