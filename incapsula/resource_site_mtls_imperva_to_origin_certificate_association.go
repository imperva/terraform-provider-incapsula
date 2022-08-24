package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceSiteMtlsCertificateAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteMtlsCertificateAssociationCreate,
		Read:   resourceSiteMtlsCertificateAssociationRead,
		Update: resourceSiteMtlsCertificateAssociationCreate,
		Delete: resourceSiteMtlsCertificateAssociationDelete,
		//todo
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of Incapsula Site to Imperva to Origin mutual TLS Certificate Association resource ID, expected site_id/certificate_id, got %s", d.Id())
				}

				d.Set("site_id", idSlice[0])
				d.Set("certificate_id", idSlice[1])

				log.Printf("[DEBUG] Importing Incapsula Site to Imperva to Origin mutual TLS Certificate Association for Site ID %s, mutual TLS Certificate Id %s", idSlice[0], idSlice[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"certificate_id": {
				Description: "The certificate file in base64 format.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceSiteMtlsCertificateAssociationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteID, certificateID, err := validateInput(d)
	if err != nil {
		return err
	}

	associationExists, err := client.GetSiteMtlsCertificateAssociation(certificateID, siteID)
	if err != nil {
		return err
	}

	if associationExists == false {
		return fmt.Errorf("Couldn't find the Incapsula Site - Imperva to Origin mutual TLS Certificate Association")
	}

	// Generate synthetic ID
	syntheticID := fmt.Sprintf("%d/%d", siteID, certificateID)
	d.SetId(syntheticID)
	log.Printf("[INFO] Created Incapsula Site to Imperva to Origin mutual TLS Certificate Association with ID: %s - site ID (%d) - certificate ID (%d)", syntheticID, siteID, certificateID)
	return nil
}

func resourceSiteMtlsCertificateAssociationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID, certificateID, err := validateInput(d)
	if err != nil {
		return err
	}

	err = client.CreateSiteMtlsCertificateAssociation(
		certificateID,
		siteID,
	)
	if err != nil {
		return err
		//todo -add error message
	}
	return resourceSiteMtlsCertificateAssociationRead(d, m)
}

func resourceSiteMtlsCertificateAssociationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID, certificateID, err := validateInput(d)

	err = client.DeleteSiteMtlsCertificateAssociation(
		certificateID,
		siteID,
	)
	if err != nil {
		//todo - check error
		return err
	}

	d.SetId("")
	return nil
}

func validateInput(d *schema.ResourceData) (int, int, error) {
	siteIDStr := d.Get("site_id").(string)
	certificateIDStr := d.Get("certificate_id").(string)

	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert Site Id for Incapsula Site to Imperva to Origin mutual TLS Certificate Association resource, actual value: %s, expected numeric id", siteIDStr)
	}

	certificateID, err := strconv.Atoi(certificateIDStr)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert certificate API Id for Incapsula Site to Imperva to Origin mutual TLS Certificate Association, actual value: %s, expected numeric id", certificateIDStr)
	}
	log.Printf("site_id %d\ncertificate_id - %d", siteID, certificateID)

	return siteID, certificateID, nil
}
