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
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of Incapsula Client CA to Imperva Certificate resource ID, expected account_id/certificate_id, got %s", d.Id())
				}

				_, err := strconv.Atoi(idSlice[1])
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", idSlice[1])
				}

				_, err = strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert Account Id from import command, actual value: %s, expected numeric id", idSlice[0])
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
			},
			"account_id": {
				Description: "Account ID", //todo add
				Type:        schema.TypeString,
				Required:    true,
			},
			// Optional Arguments
			"certificate_name": {
				Description: "The private key of the certificate in base64 format. Optional in case of PFX certificate file format. This will be encoded in sha256 in terraform state.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			//
			//"input_hash": {
			//	Description: "inputHash",
			//	Type:        schema.TypeString,
			//	Optional:    true,
			//	DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			//		newHash := createHash(d)
			//		if newHash == old {
			//			return true
			//		}
			//		return false
			//	},
			//},
		},
	}
}

func resourceClientCaCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	return fmt.Errorf("Update action is not supported fore resource incapsula_mtls_client_to_imperva_certificate. Please delete resource and create new instead")
}
func resourceClientCaCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	//inputHash := createHash(d)
	encodedCert := d.Get("certificate").(string)
	// Standard Base64 Decoding
	decodedCert, err := base64.StdEncoding.DecodeString(encodedCert)
	if err != nil {
		//todo KATRIN add info to error msg
		fmt.Printf("Error decoding Base64 encoded data %v", err)
	}
	fmt.Println(string(decodedCert))

	mTLSCertificate, err := client.AddClientCaCertificate(
		decodedCert,
		d.Get("account_id").(string),
		d.Get("certificate_name").(string),
	)

	if err != nil {
		return err
	}

	// TODO: Setting this to arbitrary value as there is only one cert for each site.
	d.SetId(strconv.Itoa(mTLSCertificate.Id))
	return resourceClientCaCertificateRead(d, m)
}

func resourceClientCaCertificateRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	accountIDStr := d.Get("account_id").(string)
	certificateIDStr := d.Id()

	clientToImpervaCertificateData, err := client.GetClientCaCertificate(accountIDStr, certificateIDStr)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(clientToImpervaCertificateData.Id))
	//d.Set("input_hash", clientToImpervaCertificateData.Hash)
	d.Set("certificate_name", clientToImpervaCertificateData.Name)
	//todo KATRIN   we don't get accountID in response! what do we do in this case?

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
