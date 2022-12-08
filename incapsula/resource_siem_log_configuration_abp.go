package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSiemLogConfigurationAbp() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiemLogConfigurationCreate,
		Read:   resourceSiemLogConfigurationRead,
		Update: resourceSiemLogConfigurationUpdate,
		Delete: resourceSiemLogConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Description: "Client account id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"configuration_name": {
				Description: "Name of the configuration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"producer": {
				Description:  "Type of the producer.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"ABP"}, false),
			},
			"datasets": {
				Description: "All datasets for the supported producers.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{"ABP"}, false),
				},
				Required: true,
			},
			"enabled": {
				Description: "True if the connection is enabled, false otherwise.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"connection_id": {
				Description: "The id of the connection for this log configuration.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "Version of the log configuration.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
