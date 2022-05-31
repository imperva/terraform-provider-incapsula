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
		Create: resourceCertificateHsmCreateAndUpdate,
		//HSM & custom certificate using same read from resource_certificate.go
		//but with different operation header in the rest request
		Read:   resourceCertificateRead,
		Update: resourceCertificateHsmCreateAndUpdate,
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

			//input hash will be created by terraform and saved on server, so we can identify any change in
			//the certificate without getting it.
			//we can't get the certificate since its sensitive data- we can only run over the current certificate.
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

// resourceCertificateHsmCreateAndUpdate Create and update are the same rest end point & logic, so we don't need
// additional function. This is the api behaviour, we can't just update part of the certificate
func resourceCertificateHsmCreateAndUpdate(d *schema.ResourceData, m interface{}) error {
	siteId := d.Get("site_id").(string)
	log.Printf("[DEBUG] Start createing HSM custome certificate for site id %s", siteId)
	client := m.(*Client)
	hSMDataDTO := HSMDataDTO{
		Certificate:    d.Get("certificate").(string),
		HsmDetailsList: getHsmDetailsFromResource(d),
	}

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
	log.Printf("[DEBUG] Done createing HSM custome certificate for site id %s, now reding the data", siteId)

	return resourceCertificateRead(d, m)
}

func getHsmDetailsFromResource(d *schema.ResourceData) []HSMDetailsDTO {
	siteId := d.Get("site_id").(string)
	log.Printf("[DEBUG] starting getHsmDetailsFromResource for site id  %s", siteId)
	var hsmDetailList []HSMDetailsDTO
	hsmDetails := d.Get("api_detail").(*schema.Set)
	for _, hsmDetail := range hsmDetails.List() {
		hsmDetailResource := hsmDetail.(map[string]interface{})
		assetDto := HSMDetailsDTO{
			KeyId:    hsmDetailResource["api_id"].(string),
			ApiKey:   hsmDetailResource["api_key"].(string),
			HostName: hsmDetailResource["hostname"].(string),
		}

		log.Printf("[DEBUG] getHsmDetailsFromResource hostname %s add for site id  %s", assetDto.HostName, siteId)
		hsmDetailList = append(hsmDetailList, assetDto)
	}

	log.Printf("[DEBUG] Done with getHsmDetailsFromResource for site id  %s", siteId)

	return hsmDetailList
}

func resourceCertificateHsmDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteId := d.Get("site_id").(string)
	log.Printf("[DEBUG] Strt removing HSM certificate for site id: %s with resourceCertificateHsmDelete", siteId)
	err := client.DeleteHsmCertificate(siteId)

	if err != nil {
		log.Printf("[ERROR] Removing HSM certificate for site id: %s faild with error: %s", siteId, err)
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
	hashFromString := calculateHashFromString(stringHash)
	log.Printf("[DEBUG] Hash for hsm custom certificate created: %s", hashFromString)

	return hashFromString
}

func calculateHashFromString(stringForHash string) string {
	h := sha1.New()
	h.Write([]byte(stringForHash))
	byteString := h.Sum(nil)
	result := hex.EncodeToString(byteString)

	return result
}
