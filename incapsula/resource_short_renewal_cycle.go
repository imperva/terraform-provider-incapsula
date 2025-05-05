package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceShortRenewalCycle() *schema.Resource {
	return &schema.Resource{
		Read:   resourceShortRenewalCycleConfigurationRead,
		Create: resourceShortRenewalCycleConfigurationCreate,
		Delete: resourceShortRenewalCycleConfigurationReadDelete,
		Update: resourceShortRenewalCycleConfigurationReadUpdate,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				log.Printf("[DEBUG] Starting to import short renewal cycle. Parameters: %s\n", d.Id())
				if len(idSlice) < 2 || len(idSlice) > 3 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected site_id/managed_certificate_settings_id/account_id", d.Id())
				}

				err := d.Set("site_id", idSlice[0])
				if err != nil {
					return nil, err
				}

				_, err = strconv.Atoi(idSlice[0])
				if err != nil || idSlice[0] == "" {
					return nil, fmt.Errorf("failed to convert site Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				d.SetId(idSlice[0])

				err = d.Set("managed_certificate_settings_id", idSlice[1])
				if err != nil {
					return nil, err
				}

				_, err = strconv.Atoi(idSlice[1])
				if err != nil || idSlice[1] == "" {
					return nil, fmt.Errorf("failed to convert managed certificate settings id from import command, actual value: %s, expected numeric id", idSlice[1])
				}

				if len(idSlice) == 3 {
					_, err = strconv.Atoi(idSlice[2])
					if err != nil {
						return nil, fmt.Errorf("failed to convert account Id from import command, actual value: %s, expected numeric id", idSlice[2])
					}

					err = d.Set("account_id", idSlice[2])
					if err != nil {
						return nil, err
					}
				}

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
			"managed_certificate_settings_id": {
				Description: "Numeric identifier of the site ssl settings to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"short_renewal_cycle": {
				Type:        schema.TypeBool,
				Description: "The short renewal cycle configuration. If true, then short renewal cycle is enabled. If false, then short renewal cycle is disabled.",
				Required:    true,
			},
		},
	}
}

func resourceShortRenewalCycleConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	shortRenewalCycleConfigurationDto, err := client.GetShortRenewalCycleConfiguration(d.Get("site_id").(string), d.Get("account_id").(string))
	if err != nil {
		return err
	}

	if shortRenewalCycleConfigurationDto.Errors != nil && len(shortRenewalCycleConfigurationDto.Errors) > 0 {
		if shortRenewalCycleConfigurationDto.Errors[0].Status == 404 || shortRenewalCycleConfigurationDto.Errors[0].Status == 401 {
			log.Printf("[INFO] operation not allowed: %s\n", shortRenewalCycleConfigurationDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(shortRenewalCycleConfigurationDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting short renewal cycleconfiguration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId(d.Get("site_id").(string))
	d.Set("short_renewal_cycle", shortRenewalCycleConfigurationDto.Data[0].ShortRenewalCycle)
	return nil
}

func resourceShortRenewalCycleConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	shortRenewalCycle := d.Get("short_renewal_cycle").(bool)
	siteId := d.Get("site_id").(string)
	accountId := d.Get("account_id").(string)
	client := m.(*Client)
	var shortRenewalCycleConfigurationDto *ShortRenewalCycleConfigurationDto
	var err error
	if shortRenewalCycle {
		log.Printf("[DEBUG] going to enable short renewal cyclefor site %s\n", siteId)
		shortRenewalCycleConfigurationDto, err = client.EnableShortRenewalCycleConfiguration(siteId, accountId)
		if err != nil {
			return err
		}
	} else {
		log.Printf("[DEBUG] going to disable short renewal cyclefor site %s\n", siteId)
		shortRenewalCycleConfigurationDto, err = client.DeleteShortRenewalCycleConfiguration(siteId, accountId)
		if err != nil {
			return err
		}
	}

	if shortRenewalCycleConfigurationDto.Errors != nil && len(shortRenewalCycleConfigurationDto.Errors) > 0 {
		if shortRenewalCycleConfigurationDto.Errors[0].Status == 404 || shortRenewalCycleConfigurationDto.Errors[0].Status == 401 {
			log.Printf("[INFO] operation not allowed: %s\n", shortRenewalCycleConfigurationDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(shortRenewalCycleConfigurationDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting create short renewal cycleconfiguration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId(siteId)
	d.Set("short_renewal_cycle", shortRenewalCycleConfigurationDto.Data[0].ShortRenewalCycle)
	return nil
}

func resourceShortRenewalCycleConfigurationReadDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(string)
	accountId := d.Get("account_id").(string)
	shortRenewalCycleConfigurationDto, err := client.DeleteShortRenewalCycleConfiguration(siteId, accountId)
	if err != nil {
		return err
	}

	if shortRenewalCycleConfigurationDto.Errors != nil && len(shortRenewalCycleConfigurationDto.Errors) > 0 {
		if shortRenewalCycleConfigurationDto.Errors[0].Status == 404 || shortRenewalCycleConfigurationDto.Errors[0].Status == 401 {
			log.Printf("[INFO] operation not allowed: %s\n", shortRenewalCycleConfigurationDto.Errors[0].Detail)
			d.SetId("")
			return nil
		}

		out, err := json.Marshal(shortRenewalCycleConfigurationDto.Errors)
		if err != nil {
			return err
		}
		return fmt.Errorf("error getting delete short renewal cycleconfiguration for site (%s): %s", d.Get("site_id"), string(out))
	}

	d.SetId(d.Get("site_id").(string))
	d.Set("short_renewal_cycle", shortRenewalCycleConfigurationDto.Data[0].ShortRenewalCycle)
	return nil
}

func resourceShortRenewalCycleConfigurationReadUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceShortRenewalCycleConfigurationCreate(d, m)
}
