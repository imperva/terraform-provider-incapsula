package incapsula

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

var apiDetailsResource = schema.Resource{
	Schema: map[string]*schema.Schema{
		"api_id": {
			Description: "The api id in Fortanix",
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
		},
		"api_key": {
			Description: "the api key in Fortanix",
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
		},

		"hostname": {
			Description: "The Fortanix hostname ",
			Type:        schema.TypeString,
			Required:    true,
		},
	},
}

func resourceCustomCertificateHsm() *schema.Resource {
	return &schema.Resource{
		Create: resourceCertificateHsmCreate,
		Read:   resourceCertificateRead, //using same read from resource_certificate.go
		//Update: resourceCertificateHsmUpdate,
		Update: resourceCertificateHsmCreate,
		Delete: resourceCertificateHsmDelete,
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
			"certificate": { //The public certificate
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"api_detail": {
				Description: "The details of the API in Fortanix",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &apiDetailsResource,
				Set:         schema.HashResource(&apiDetailsResource),
			},

			//-----------
			"input_hash": {
				Description: "inputHash",
				Type:        schema.TypeString,
				Optional:    true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					hSMDataDTO := HSMDataDTO{
						Certificate:    d.Get("certificate").(string),
						HsmDetailsList: getHsmDetailsFromResource(d),
					}
					siteId := d.Get("site_id").(string)
					newHash := createHashFromHSMDataDTO(&hSMDataDTO, siteId)
					if newHash == old {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceCertificateHsmCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] Start createing HSM custome certificate")
	client := m.(*Client)
	hSMDataDTO := HSMDataDTO{
		Certificate:    d.Get("certificate").(string),
		HsmDetailsList: getHsmDetailsFromResource(d),
	}

	siteId := d.Get("site_id").(string)
	inputHash := createHashFromHSMDataDTO(&hSMDataDTO, siteId)
	_, err := client.AddHsmCertificate(
		siteId,
		inputHash,
		&hSMDataDTO,
	)

	if err != nil {
		log.Printf("[ERROR] Uploading HSM custom certificate to Site ID %s got error: %s\n", d.Get("site_id"), err)
		return err
	}

	d.SetId("12345")

	return resourceCertificateRead(d, m)
}

func getHsmDetailsFromResource(d *schema.ResourceData) []HSMDetailsDTO {
	var hsmDetailList []HSMDetailsDTO
	hsmDetails := d.Get("api_detail").(*schema.Set)
	for _, hsmDetail := range hsmDetails.List() {
		hsmDetailResource := hsmDetail.(map[string]interface{})
		assetDto := HSMDetailsDTO{
			KeyId:    hsmDetailResource["api_id"].(string),
			ApiKey:   hsmDetailResource["api_key"].(string),
			HostName: hsmDetailResource["hostname"].(string),
		}
		hsmDetailList = append(hsmDetailList, assetDto)
	}
	return hsmDetailList
}

//func resourceCertificateHsmRead(d *schema.ResourceData, m interface{}) error {
//	client := m.(*Client)
//	siteID := d.Get("site_id").(string)
//	listCertificatesResponse, err := client.ListCertificates(siteID)
//	log.Printf("[INFO] Reading HSM custome certificate for site id %s ", siteID)
//
//	// List data centers response object may indicate that the Site ID has been deleted (9413)
//	if listCertificatesResponse != nil && listCertificatesResponse.Res == 9413 {
//		log.Printf("[INFO] Incapsula Site ID %s has already been deleted: %s\n", d.Get("site_id"), err)
//		d.SetId("")
//		return nil
//	}
//
//	if err != nil {
//		log.Printf("[ERROR] Could not read custom certificate from Incapsula site for site_id: %s, %s\n", siteID, err)
//		return err
//	}
//
//	d.Set("input_hash", listCertificatesResponse.SSL.CustomCertificate.InputHash)
//	d.SetId("12345")
//
//	return nil
//}

func resourceCertificateHsmUpdate(d *schema.ResourceData, m interface{}) error {
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

func resourceCertificateHsmDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(string)
	log.Printf("[DEBUG] Strt removing HSM certificate for site id: %d with resourceCertificateHsmDelete", siteId)
	err := client.DeleteHsmCertificate(siteId)

	if err != nil {
		log.Printf("[ERROR] Removing HSM certificate for site id: %d faild with error: %s", siteId, err)
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}

func createHashFromHSMDataDTO(hSMDataDTO *HSMDataDTO, siteId string) string {
	log.Printf("[DEBUG] Starting create hash for hsm custom certificate")
	var hsmString string
	for _, hSMDetailsDTO := range hSMDataDTO.HsmDetailsList {
		hsmString += hSMDetailsDTO.HostName + hSMDetailsDTO.KeyId + hSMDetailsDTO.ApiKey
	}

	stringHash := siteId + hSMDataDTO.Certificate + hsmString
	result := calculateHashFromString(stringHash)

	return result
}

func calculateHashFromString(stringForHash string) string {
	h := sha1.New()
	h.Write([]byte(stringForHash))
	byteString := h.Sum(nil)
	result := hex.EncodeToString(byteString)

	return result
}
