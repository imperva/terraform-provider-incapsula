package incapsula

import (
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTXTRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceTXTRecordCreate,
		Read:   resourceTXTRecordRead,
		Update: resourceTXTRecordUpdate,
		Delete: resourceTXTRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required Argument
			"site_id": {
				Description: "Numeric identifier of the site.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			// Optional Arguments
			"txt_record_value_one": {
				Description: "New value for txt record number one.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"txt_record_value_two": {
				Description: "New value for txt record number two.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"txt_record_value_three": {
				Description: "New value for txt record number three.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"txt_record_value_four": {
				Description: "New value for txt record number four.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"txt_record_value_five": {
				Description: "New value for txt record number five.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceTXTRecordCreate(d *schema.ResourceData, m interface{}) error {
	// Implement by create the TXT Records

	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	TXTRecordOne := d.Get("txt_record_value_one").(string)
	TXTRecordTwo := d.Get("txt_record_value_two").(string)
	TXTRecordThree := d.Get("txt_record_value_three").(string)
	TXTRecordFour := d.Get("txt_record_value_four").(string)
	TXTRecordFive := d.Get("txt_record_value_five").(string)

	log.Printf("[INFO] Setting Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID)

	_, err := client.CreateTXTRecord(siteID, TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive)

	if err != nil {
		log.Printf("[ERROR] Could not set Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n%s", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID, err)
		return err
	}

	// Set the ID
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.Printf("[INFO] Set Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID)

	return resourceTXTRecordRead(d, m)
}

func resourceTXTRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)

	if (d.HasChange("txt_record_value_one") && d.Get("txt_record_value_one") != "") ||
		(d.HasChange("txt_record_value_two") && d.Get("txt_record_value_two") != "") ||
		(d.HasChange("txt_record_value_three") && d.Get("txt_record_value_three") != "") ||
		(d.HasChange("txt_record_value_four") && d.Get("txt_record_value_four") != "") ||
		(d.HasChange("txt_record_value_five") && d.Get("txt_record_value_five") != "") {

		TXTRecordOne := d.Get("txt_record_value_one").(string)
		TXTRecordTwo := d.Get("txt_record_value_two").(string)
		TXTRecordThree := d.Get("txt_record_value_three").(string)
		TXTRecordFour := d.Get("txt_record_value_four").(string)
		TXTRecordFive := d.Get("txt_record_value_five").(string)

		log.Printf("[INFO] Setting Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID)

		_, err := client.UpdateTXTRecord(siteID, TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive)

		if err != nil {
			log.Printf("[ERROR] Could not set Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n%s", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID, err)
			return err
		}

		// Set the ID
		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.Printf("[INFO] Set Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID)
	}
	return resourceTXTRecordRead(d, m)
}

func resourceTXTRecordRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the TXTRecordResponse for the TXT Records
	client := m.(*Client)

	recordResponse, err := client.ReadTXTRecords(d.Get("site_id").(int))

	// Gte TXT response object
	if recordResponse != nil {
		// Res can oscillate between strings and ints
		if recordResponse.Res == 0 {
			d.Set("txt_record_value_one", recordResponse.TxtRecordValueOne)
			d.Set("txt_record_value_two", recordResponse.TxtRecordValueTwo)
			d.Set("txt_record_value_three", recordResponse.TxtRecordValueThree)
			d.Set("txt_record_value_four", recordResponse.TxtRecordValueFour)
			d.Set("txt_record_value_five", recordResponse.TxtRecordValueFive)
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func resourceTXTRecordDelete(d *schema.ResourceData, m interface{}) error {
	// Implement by deleting the a TXT Record

	client := m.(*Client)
	siteID := d.Get("site_id").(int)

	if d.HasChange("txt_record_value_one") && d.Get("txt_record_value_one") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 1, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "1")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 1, for siteID: %d\n%s", siteID, err)
		}
	} else if d.HasChange("txt_record_value_two") && d.Get("txt_record_value_two") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 2, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "2")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 2, for siteID: %d\n%s", siteID, err)
		}
	} else if d.HasChange("txt_record_value_three") && d.Get("txt_record_value_three") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 3, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "3")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 3, for siteID: %d\n%s", siteID, err)
		}
	} else if d.HasChange("txt_record_value_four") && d.Get("txt_record_value_four") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 4, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "4")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 4, for siteID: %d\n%s", siteID, err)
		}
	} else if d.HasChange("txt_record_value_five") && d.Get("txt_record_value_five") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 5, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "5")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 5, for siteID: %d\n%s", siteID, err)
		}
	}
	return resourceTXTRecordRead(d, m)
}
