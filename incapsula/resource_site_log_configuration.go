package incapsula

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func resourceSiteLogConfiguration() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteLogConfigurationCreate,
		Read:   resourceSiteLogConfigurationRead,
		Update: resourceSiteLogConfigurationUpdate,
		Delete: resourceSiteLogConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				siteID := d.Id()
				if err := d.Set("site_id", siteID); err != nil {
					return nil, err
				}
				d.SetId(siteID)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"logs_account_id": {
				Description: "Available only for Enterprise Plan customers that purchased the Logs Integration SKU. Numeric identifier of the account that purchased the logs integration SKU and which collects the logs. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"log_level": {
				Description: "The log level. Options are `full`, `security`, and `none`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"data_storage_region": {
				Description: "The data region to use. Options are `APAC`, `AU`, `EU`, and `US`.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func resourceSiteLogConfigurationCreate(d *schema.ResourceData, m interface{}) error {
	// Log the creation process
	log.Printf("[INFO] Creating Incapsula site log configuration")

	// Call the update function to create the resource
	err := resourceSiteLogConfigurationUpdate(d, m)
	if err != nil {
		return err
	}

	// Set the resource ID
	siteID := d.Get("site_id").(string)
	d.SetId(siteID)

	// Log the resource ID
	log.Printf("[INFO] Incapsula site log configuration created with site_id: %s", siteID)

	return resourceSiteLogConfigurationRead(d, m)
}

func resourceSiteLogConfigurationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteId, err := strconv.Atoi(d.Get("site_id").(string))
	if err != nil {
		log.Printf("[ERROR] Could not convert site_id to int: %s\n", err)
		return err
	}

	// Get the log level for the site
	siteStatusResponse, err := client.SiteStatus("nil", siteId)
	if siteStatusResponse.LogLevel != "" {
		d.Set("log_level", siteStatusResponse.LogLevel)
	}

	// Get the data storage region for the site
	dataStorageRegionResponse, err := client.GetDataStorageRegion(d.Get("site_id").(string))
	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula site data storage region for site id: %s, %s\n", d.Get("site_id").(string), err)
		return err
	}
	d.Set("data_storage_region", dataStorageRegionResponse.Region)

	return nil
}

func resourceSiteLogConfigurationUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	// Update the log level and logs account id for the site
	err := updateSiteLogLevel(client, d)
	if err != nil {
		return err
	}

	// Update the data storage region for the site
	err = updateSiteDataStorageRegion(client, d)
	if err != nil {
		return err
	}

	return resourceSiteLogConfigurationRead(d, m)
}

func resourceSiteLogConfigurationDelete(d *schema.ResourceData, m interface{}) error {
	// Implement the delete logic here
	d.SetId("")
	return nil
}

func updateSiteLogLevel(client *Client, d *schema.ResourceData) error {
	if d.HasChange("log_level") ||
		d.HasChange("logs_account_id") {
		logLevel := d.Get("log_level").(string)
		logsAccountId := d.Get("logs_account_id").(string)
		siteID := d.Get("site_id").(string)
		err := client.UpdateLogLevel(siteID, logLevel, logsAccountId)
		if err != nil {
			log.Printf("[ERROR] Could not update Incapsula site log level: %s and logs account id: %s for site_id: %s %s\n", logLevel, logsAccountId, siteID, err)
			return err
		}
	}
	return nil
}

func updateSiteDataStorageRegion(client *Client, d *schema.ResourceData) error {
	if d.HasChange("data_storage_region") {
		dataStorageRegion := d.Get("data_storage_region").(string)
		siteID := d.Get("site_id").(string)
		_, err := client.UpdateDataStorageRegion(siteID, dataStorageRegion)
		if err != nil {
			log.Printf("[ERROR] Could not set Incapsula site data storage region with value (%s) for site_id: %s %s\n", dataStorageRegion, siteID, err)
			return err
		}
	}
	return nil
}
