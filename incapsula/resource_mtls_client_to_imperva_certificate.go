package incapsula

import (
	//"crypto/sha1"
	//"encoding/hex"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceMtlsClientToImpervaCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceClientCaCertificateCreate,
		Read:   resourceClientCaCertificateRead,
		Update: resourceClientCaCertificateUpdate,
		Delete: resourceClientCaCertificateDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.MarkNewResource()
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of Incapsula Client CA to Imperva Certificate resource ID, expected account_id/certificate_id, got %s", d.Id())
				}

				_, err := strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert Account Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				_, err = strconv.Atoi(idSlice[1])
				if err != nil {
					fmt.Errorf("failed to convert Certificate Id from import command, actual value: %s, expected numeric id", idSlice[1])
				}

				d.Set("account_id", idSlice[0])
				d.SetId(idSlice[1])

				log.Printf("[DEBUG] Importing Incapsula Site to Client to Imperva mutual TLS Certificate for Account ID %s, Certificate Id %s", idSlice[0], idSlice[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			// Required Arguments
			"certificate": {
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeString,
				Required:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if new == ignoreSensitiveVariableString {
						return true
					}
					return false
				},
			},
			"account_id": {
				Description: "Account ID to operte on",
				Type:        schema.TypeString,
				Required:    true,
			},
			// Optional Arguments
			"certificate_name": {
				Description: "The certificate name",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceClientCaCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Update action is not supported for current resource. Please create a new Mutual TLS Client to Imperva CA Certificate resource and only then, delete the old one.")
}
func resourceClientCaCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	encodedCert := d.Get("certificate").(string)
	// Standard Base64 Decoding
	var decodedCert []byte
	var err error

	decodedCert, err = base64.StdEncoding.DecodeString(encodedCert)
	if err != nil {
		fmt.Printf("Error decoding Base64 encoded data from certificate field of incapsula_site_tls_settings resource %v", err)
	}

	mTLSCertificate, err := client.AddClientCaCertificate(
		decodedCert,
		d.Get("account_id").(string),
		d.Get("certificate_name").(string),
	)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(mTLSCertificate.Id))
	return resourceClientCaCertificateRead(d, m)
}

func resourceClientCaCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	accountIDStr := d.Get("account_id").(string)
	certificateIDStr := d.Id()

	clientToImpervaCertificateData, certificateExits, err := client.GetClientCaCertificate(accountIDStr, certificateIDStr)
	if !certificateExits && !d.IsNewResource() {
		d.SetId("")
		return nil
	}

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(clientToImpervaCertificateData.Id))
	d.Set("certificate_name", clientToImpervaCertificateData.Name)
	return nil
}

func resourceClientCaCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	accountID := d.Get("account_id").(string)
	certificateID := d.Id()
	err := client.DeleteClientCaCertificate(accountID, certificateID)
	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")
	return nil
}
