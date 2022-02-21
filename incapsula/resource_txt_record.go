package incapsula

import (
	"log"
	"strconv"
	"strings"

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
				ForceNew:    true,
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
	d.SetId(strconv.Itoa(siteID))
	log.Printf("[INFO] Create Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID)

	return resourceTXTRecordRead(d, m)
}

func resourceTXTRecordUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID := d.Get("site_id").(int)
	errDelete := deleteSpecificTXTRecordIfNeeded(d, siteID, client)
	if errDelete != nil {
		return errDelete
	}
	errUpdate := updateSpecificTXTRecordIfNeeded(d, siteID, client)
	if errUpdate != nil {
		return errUpdate
	}

	return resourceTXTRecordRead(d, m)
}

func updateSpecificTXTRecordIfNeeded(d *schema.ResourceData, siteID int, client *Client) error {
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
		log.Printf("[INFO] Update Incapsula TXT Records: %s, %s, %s, %s, %s, for siteID: %d\n", TXTRecordOne, TXTRecordTwo, TXTRecordThree, TXTRecordFour, TXTRecordFive, siteID)
	}
	return nil
}

func deleteSpecificTXTRecordIfNeeded(d *schema.ResourceData, siteID int, client *Client) error {
	if d.HasChange("txt_record_value_one") && d.Get("txt_record_value_one") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 1, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "1")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 1, for siteID: %d\n%s", siteID, err)
			return err
		}
	} else if d.HasChange("txt_record_value_two") && d.Get("txt_record_value_two") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 2, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "2")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 2, for siteID: %d\n%s", siteID, err)
			return err
		}
	} else if d.HasChange("txt_record_value_three") && d.Get("txt_record_value_three") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 3, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "3")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 3, for siteID: %d\n%s", siteID, err)
			return err
		}
	} else if d.HasChange("txt_record_value_four") && d.Get("txt_record_value_four") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 4, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "4")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 4, for siteID: %d\n%s", siteID, err)
			return err
		}
	} else if d.HasChange("txt_record_value_five") && d.Get("txt_record_value_five") == "" {
		log.Printf("[INFO] Delete Incapsula TXT Record 5, for siteID: %d", siteID)
		err := client.DeleteTXTRecord(siteID, "5")
		if err != nil {
			log.Printf("[ERROR] Could not delete Incapsula TXT Records 5, for siteID: %d\n%s", siteID, err)
			return err
		}
	}
	return nil
}

func resourceTXTRecordRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the TXTRecordResponse for the TXT Records
	client := m.(*Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] The ID should be numeric. Currrent value: %s", d.Id())
		return err
	}

	recordResponse, err := client.ReadTXTRecords(id)
	d.Set("site_id", id)

	// Gte TXT response object
	if recordResponse != nil {
		// Res can oscillate between strings and ints
		if recordResponse.Res == 0 && !strings.Contains(recordResponse.ResMessage, "no TXT records") {
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
	err := client.DeleteTXTRecordAll(siteID)
	if err != nil {
		log.Printf("[ERROR] Could not delete all Incapsula TXT Records, for siteID: %d\n%s", siteID, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")
	return nil
}
