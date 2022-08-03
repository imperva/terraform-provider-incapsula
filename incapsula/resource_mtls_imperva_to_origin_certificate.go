package incapsula

import (
	//"crypto/sha1"
	//"encoding/hex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func resourceMTLSImpervaToOriginCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceMTLSImpervaToOriginCertificateCreate,
		Read:   resourceMTLSImpervaToOriginCertificateRead,
		Update: resourceMTLSImpervaToOriginCertificateUpdate,
		Delete: resourceMTLSImpervaToOriginCertificateDelete,
		//todo
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
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
			"certificate_name": {
				Description: "The certificate name.",
				Type:        schema.TypeString,
				Optional:    true,
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

func resourceMTLSImpervaToOriginCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	inputHash := createHash(d)
	mTLSCertificateData, err := client.AddMTLSCertificate(
		d.Get("certificate").(string),
		d.Get("private_key").(string),
		d.Get("passphrase").(string),
		d.Get("certificate_name").(string),
		inputHash,
	)
	if err != nil {
		return err
	}

	//// TODO: Setting this to arbitrary value as there is only one cert for each site.
	d.SetId(strconv.Itoa(mTLSCertificateData.Id))
	d.Set("input_hash", mTLSCertificateData.Hash)

	return resourceMTLSImpervaToOriginCertificateRead(d, m)
}

func resourceMTLSImpervaToOriginCertificateRead(d *schema.ResourceData, m interface{}) error {
	//// Implement by reading the ListCertificatesResponse for the data center
	client := m.(*Client)
	mTLSCertificateData, err := client.GetTLSCertificate(d.Id())
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(mTLSCertificateData.Id))
	d.Set("input_hash", mTLSCertificateData.Hash)
	d.Set("certificate_name", mTLSCertificateData.Name)

	return nil
}

func resourceMTLSImpervaToOriginCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	inputHash := createHash(d)
	log.Printf("hash:\n%s", inputHash)
	mTLSCertificateData, err := client.UpdateMTLSCertificate(
		d.Id(),
		d.Get("certificate").(string),
		d.Get("private_key").(string),
		d.Get("passphrase").(string),
		d.Get("certificate_name").(string),
		inputHash,
	)
	if err != nil {
		return err
	}

	//// TODO: Setting this to arbitrary value as there is only one cert for each site.
	d.SetId(strconv.Itoa(mTLSCertificateData.Id))
	d.Set("input_hash", mTLSCertificateData.Hash)
	return resourceMTLSImpervaToOriginCertificateRead(d, m)
}

func resourceMTLSImpervaToOriginCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	err := client.DeleteMTLSCertificate(d.Id())

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}

//
//func createHash(d *schema.ResourceData) string {
//	certificate := d.Get("certificate").(string)
//	passphrase := d.Get("passphrase").(string)
//	privateKey := d.Get("private_key").(string)
//	result := calculateHash(certificate, passphrase, privateKey)
//	return result
//}
//
//func calculateHash(certificate, passphrase, privateKey string) string {
//	h := sha1.New()
//	stringForHash := certificate + privateKey + passphrase
//	h.Write([]byte(stringForHash))
//	byteString := h.Sum(nil)
//	result := hex.EncodeToString(byteString)
//	return result
//}

//todo ????
//func getOperation(d *schema.ResourceData) string {
//	isCustomCertificate := d.Get("api_detail") != nil
//	operation := ReadCustomCertificate
//	if isCustomCertificate {
//		operation = ReadHSMCustomCertificate
//	}
//	log.Printf("[DEBUG] Selected oprtaion type for rest request is: %s", operation)
//
//	return operation
//}
