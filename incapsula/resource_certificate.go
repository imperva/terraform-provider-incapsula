package incapsula

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				ForceNew:    true,
			},
			"certificate": {
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			// Optional Arguments
			"private_key": {
				Description: "The private key of the certificate in base64 format. Optional in case of PFX certificate file format. This will be encoded in sha256 in terraform state.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase": {
				Description: "The passphrase used to protect your SSL certificate. This will be encoded in sha256 in terraform state.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"input_hash": {
				Description: "inputHash",
				Type:        schema.TypeString,
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					newHash := createHash(d)
					if newHash == old {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	inputHash := createHash(d)
	_, err := client.AddCertificate(
		d.Get("site_id").(string),
		d.Get("certificate").(string),
		d.Get("private_key").(string),
		d.Get("passphrase").(string),
		inputHash,
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

	siteID := d.Get("site_id").(string)
	operation := getOperation(d)
	listCertificatesResponse, err := client.ListCertificates(siteID, operation)

	// List data centers response object may indicate that the Site ID has been deleted (9413)
	if listCertificatesResponse != nil && listCertificatesResponse.Res == 9413 {
		log.Printf("[INFO] Incapsula Site with ID %s has already been deleted: %s\n", d.Get("site_id"), err)
		d.SetId("")
		return nil
	}

	if err != nil {
		log.Printf("[ERROR] Could not read custom certificate from Incapsula site for site_id: %s, %s\n", siteID, err)
		return err
	}

	d.Set("input_hash", listCertificatesResponse.SSL.CustomCertificate.InputHash)
	d.SetId("12345")

	return nil
}

func getOperation(d *schema.ResourceData) string {
	isCustomCertificate := d.Get("api_detail") != nil
	operation := ReadCustomCertificate
	if isCustomCertificate {
		operation = ReadHSMCustomCertificate
	}
	log.Printf("[DEBUG] Selected oprtaion type for rest request is: %s", operation)

	return operation
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	inputHash := createHash(d)

	_, err := client.EditCertificate(
		d.Get("site_id").(string),
		d.Get("certificate").(string),
		d.Get("private_key").(string),
		d.Get("passphrase").(string),
		inputHash,
	)

	if err != nil {
		return err
	}

	d.SetId("12345")
	return resourceCertificateRead(d, m)
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

func createHash(d *schema.ResourceData) string {
	certificate := d.Get("certificate").(string)
	passphrase := d.Get("passphrase").(string)
	privateKey := d.Get("private_key").(string)
	result := calculateHash(certificate, passphrase, privateKey)
	return result
}

func calculateHash(certificate, passphrase, privateKey string) string {
	h := sha1.New()
	stringForHash := certificate + privateKey + passphrase
	h.Write([]byte(stringForHash))
	byteString := h.Sum(nil)
	result := hex.EncodeToString(byteString)
	return result
}
