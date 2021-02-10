package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOriginPOP() *schema.Resource {
	return &schema.Resource{
		Create:   resourceOriginPOPUpdate,
		Read:     resourceOriginPOPRead,
		Update:   resourceOriginPOPUpdate,
		Delete:   resourceOriginPOPDelete,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"dc_id": {
				Description: "Numeric identifier of the data center.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"origin_pop": {
				Description: "The Origin POP code (must be lowercase), e.g: iad. Note, this field is create/update only. Reads are not supported as the API doesn't exist yet. Note that drift may happen.",
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					// Check if valid JSON
					d := val.(string)
					if strings.ToLower(d) != d {
						errs = append(errs, fmt.Errorf("%q must be lowercase, please check your origin POP code, got: %s", key, d))
					}
					return
				},
			},
		},
	}
}

func resourceOriginPOPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	dcID := d.Get("dc_id").(int)
	originPOP := d.Get("origin_pop").(string)

	log.Printf("[INFO] Setting Incapsula origin POP: %s for data center: %d\n", originPOP, dcID)

	err := client.SetOriginPOP(dcID, originPOP)

	if err != nil {
		log.Printf("[ERROR] Could not set Incapsula origin POP: %s for data center: %d: %s\n", originPOP, dcID, err)
		return err
	}

	// Set the ID
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	log.Printf("[INFO] Set Incapsula origin POP: %s for data center: %d\n", originPOP, dcID)

	return nil
}

func resourceOriginPOPRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceOriginPOPDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
