package incapsula

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.SetId("12345")
				d.Set("site_id", d.Get("site_id").(string))
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "Numeric identifier of the site to operate on.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"certificate": {
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Optional Arguments
			"private_key": {
				Description: "The private key of the certificate in base64 format. Optional in case of PFX certificate file format.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"passphrase": {
				Description: "The passphrase used to protect your SSL certificate.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.AddCertificate(
		d.Get("site_id").(string),
		d.Get("certificate").(string),
		d.Get("private_key").(string),
		d.Get("passphrase").(string),
	)

	if err != nil {
		return err
	}

	// TODO: Setting this to arbitrary value as there is only one cert for each site.
	d.SetId("12345")

	return resourceCertificateRead(d, m)
}

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	// Implement by reading the ListCertificatesResponse for the data center
	client := m.(*Client)

	site_id := d.Get("site_id").(string)

	_, err := client.ListCertificates(site_id)

	if err != nil {
		log.Printf("[ERROR] Could not read custom certificate from Incapsula site for site_id: %s, %s\n", site_id, err)
		return err
	}

	d.SetId("12345")

	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	_, err := client.EditCertificate(
		d.Get("site_id").(string),
		d.Get("certificate").(string),
		d.Get("private_key").(string),
		d.Get("passphrase").(string),
	)

	if err != nil {
		return err
	}

	d.SetId("12345")

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := client.DeleteCertificate(d.Get("site_id").(string))

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
