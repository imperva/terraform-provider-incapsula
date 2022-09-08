package incapsula

import (
	//"crypto/sha1"
	//"encoding/hex"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

func resourceMtlsImpervaToOriginCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceMTLSImpervaToOriginCertificateCreate,
		Read:   resourceMTLSImpervaToOriginCertificateRead,
		Update: resourceMTLSImpervaToOriginCertificateUpdate,
		Delete: resourceMTLSImpervaToOriginCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"certificate": {
				Description: "Your mTLS client certificate file in base64 format. Supported formats: PEM, DER and PFX. Only RSA certificates are currently supported. The certificate RSA key size must be 2048 bit or less. The certificate must be issued by a certificate authority (CA) and cannot be self-signed.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			// Optional Arguments
			"private_key": {
				Description: "Your private key file in base64 format. Supported formats: PEM, DER. If PFX certificate is used, then this field can remain empty.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase": {
				Description: "Your private key passphrase. Leave empty if the private key is not password protected.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"certificate_name": {
				Description: "A descriptive name for your mTLS client certificate.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"account_id": {
				Description: "Numeric identifier of the account to operate on. If not specified, operation will be performed on the account identified by the authentication parameters.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"input_hash": {
				Description: "Currently ignored. If terraform plan flags this field as changed, it means that any of: certificate, private_key, or passphrase has changed.",
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
	accountID := d.Get("account_id").(string)

	encodedCert := d.Get("certificate").(string)
	// Standard Base64 Decoding
	decodedCert, err := base64.StdEncoding.DecodeString(encodedCert)
	if err != nil {
		fmt.Printf("Error decoding Base64 encoded data %v", err)
	}

	encodedPKey := d.Get("private_key").(string)
	// Standard Base64 Decoding
	decodedPKey, err := base64.StdEncoding.DecodeString(encodedPKey)
	if err != nil {
		fmt.Printf("Error decoding Base64 encoded data %v", err)
	}

	_, err = client.AddMTLSCertificate(
		decodedCert,
		decodedPKey,
		d.Get("passphrase").(string),
		d.Get("certificate_name").(string),
		inputHash,
		accountID,
	)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Created mutual TLS Imperva to Origin Certificate with ID: %s\n", d.Id())
	return resourceMTLSImpervaToOriginCertificateRead(d, m)
}

func resourceMTLSImpervaToOriginCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(string)

	mTLSCertificateData, err := client.GetMTLSCertificate(d.Id(), accountID)
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
	accountID := d.Get("account_id").(string)

	encodedCert := d.Get("certificate").(string)
	// Standard Base64 Decoding
	decodedCert, err := base64.StdEncoding.DecodeString(encodedCert)
	if err != nil {
		log.Printf("Error decoding Base64 encoded data %v", err)
	}
	log.Println(string(decodedCert))

	encodedPKey := d.Get("private_key").(string)
	// Standard Base64 Decoding
	decodedPKey, err := base64.StdEncoding.DecodeString(encodedPKey)
	if err != nil {
		log.Printf("Error decoding Base64 encoded data %v", err)
	}
	log.Println(string(decodedPKey))

	mTLSCertificateData, err := client.UpdateMTLSCertificate(
		d.Id(),
		decodedCert,
		decodedPKey,
		d.Get("passphrase").(string),
		d.Get("certificate_name").(string),
		inputHash,
		accountID,
	)

	if err != nil {
		return err
	}

	//// TODO: Setting this to arbitrary value as there is only one cert for each site.
	d.SetId(strconv.Itoa(mTLSCertificateData.Id))
	d.Set("input_hash", mTLSCertificateData.Hash)

	log.Printf("[INFO] Updated mutual TLS Imperva to Origin Certificate with ID: %s\n", d.Id())
	return resourceMTLSImpervaToOriginCertificateRead(d, m)
}

func resourceMTLSImpervaToOriginCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	accountID := d.Get("account_id").(string)

	err := client.DeleteMTLSCertificate(d.Id(), accountID)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
