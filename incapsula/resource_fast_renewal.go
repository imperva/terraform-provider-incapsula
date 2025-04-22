package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceFastRenewal() *schema.Resource {
	return &schema.Resource{
		Read:   resourceFastRenewalConfigurationRead,
		Create: resourceFastRenewalConfigurationCreate,
		Delete: resourceFastRenewalConfigurationReadDelete,
		Update: resourceFastRenewalConfigurationReadUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.Set("site_id", d.Id())
				d.SetId(d.Id())
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"account_id": {
				Description: "Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"fast_renewal": {
				Type:        schema.TypeBool,
				Description: "The fast renewal configuration. If true, then fast renewal is enabled. If false, then fast renewal is disabled.",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The fast renewal configuration id",
				Computed:    true,
			},
		},
	}
}

func resourceFastRenewalConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	fastRenewalConfigurationDto, err := client.GetFastRenewalConfiguration(d.Get("site_id").(string), d.Get("account_id").(string))
	if err != nil {
		return err
	}

	if fastRenewalConfigurationDto.Errors != nil && len(fastRenewalConfigurationDto.Errors) > 0 {
		if fastRenewalConfigurationDto.Errors[0].Status == 404 || fastRenewalConfigurationDto.Errors[0].Status == 401 {
			log.Printf("[INFO] Operation not allowed: %s\n", fastRenewalConfigurationDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(fastRenewalConfigurationDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting fast renewal configuration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId(d.Get("site_id").(string))
	d.Set("fast_renewal", fastRenewalConfigurationDto.Data[0].FastRenewal)
	return nil
}

func resourceFastRenewalConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	fastRenewal := d.Get("fast_renewal").(bool)
	siteId := d.Get("site_id").(string)
	accountId := d.Get("account_id").(string)
	client := m.(*Client)
	var fastRenewalConfigurationDto *FastRenewalConfigurationDto
	var err error
	if fastRenewal {
		log.Printf("[DEBUG] going to enable fast renewal for site %s\n", siteId)
		fastRenewalConfigurationDto, err = client.EnableFastRenewalConfiguration(siteId, accountId)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] going to disable fast renewal for site %s\n", siteId)
		fastRenewalConfigurationDto, err = client.DeleteFastRenewalConfiguration(siteId, accountId)
		if err != nil {
			return err
		}
	}

	if fastRenewalConfigurationDto.Errors != nil && len(fastRenewalConfigurationDto.Errors) > 0 {
		if fastRenewalConfigurationDto.Errors[0].Status == 404 || fastRenewalConfigurationDto.Errors[0].Status == 401 {
			log.Printf("[INFO] Operation not allowed: %s\n", fastRenewalConfigurationDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(fastRenewalConfigurationDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting create fast renewal configuration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId(siteId)
	d.Set("fast_renewal", fastRenewalConfigurationDto.Data[0].FastRenewal)
	return nil
}

func resourceFastRenewalConfigurationReadDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(string)
	accountId := d.Get("account_id").(string)
	fastRenewalConfigurationDto, err := client.DeleteFastRenewalConfiguration(siteId, accountId)
	if err != nil {
		return err
	}

	if fastRenewalConfigurationDto.Errors != nil && len(fastRenewalConfigurationDto.Errors) > 0 {
		if fastRenewalConfigurationDto.Errors[0].Status == 404 || fastRenewalConfigurationDto.Errors[0].Status == 401 {
			log.Printf("[INFO] Operation not allowed: %s\n", fastRenewalConfigurationDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(fastRenewalConfigurationDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting delete fast renewal configuration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId(d.Get("site_id").(string))
	d.Set("fast_renewal", fastRenewalConfigurationDto.Data[0].FastRenewal)
	return nil
}

func resourceFastRenewalConfigurationReadUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceFastRenewalConfigurationCreate(d, m)
}
