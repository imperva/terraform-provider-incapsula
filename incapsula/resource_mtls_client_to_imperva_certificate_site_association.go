package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceMtlsClientToImpervaCertificateSiteAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteMtlsClientToImpervaCertificateAssociationCreate,
		Read:   resourceSiteMtlsClientToImpervaCertificateAssociationRead,
		Update: resourceSiteMtlsClientToImpervaCertificateAssociationCreate,
		Delete: resourceSiteMtlsClientToImpervaCertificateAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				d.MarkNewResource()
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 2 || idSlice[0] == "" || idSlice[1] == "" {
					return nil, fmt.Errorf("unexpected format of Incapsula Client to Imperva CA Certificate Site Association resource ID, expected site_id/certificate_id, got %s", d.Id())
				}

				_, err := strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				_, err = strconv.Atoi(idSlice[1])
				if err != nil {
					fmt.Errorf("failed to convert Certificate Id from import command, actual value: %s, expected numeric id", idSlice[1])
				}

				d.Set("site_id", idSlice[0])
				d.Set("certificate_id", idSlice[1])

				log.Printf("[DEBUG] Importing Incapsula Client to Imperva CA Certificate Site Association for Site ID %s, mutual TLS Certificate Id %s,", idSlice[0], idSlice[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"site_id": {
				Description: "The certificate file in base64 format.", //sopported formats
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

func resourceSiteMtlsClientToImpervaCertificateAssociationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	siteID, certificateID, err := validateInput(d)
	if err != nil {
		return err
	}

	mTLSCertificateData, associationExists, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(siteID, certificateID)
	if err != nil {
		return err
	}
	if !associationExists && !d.IsNewResource() {
		log.Printf("Site to mutual TLS Imperva to Origin Certificate association with Site ID %d, Certificate ID %d doesn't exist any more. The resource will be deleted from terraform state.", siteID, certificateID)
		d.SetId("")
		return nil
	}
	if mTLSCertificateData == nil {
		return fmt.Errorf("Couldn't find the Incapsula Client to Imperva CA Certificate Site Association. Site Id %d, certificate ID %d", siteID, certificateID)
	}

	// Generate synthetic ID
	syntheticID := fmt.Sprintf("%d/%d", siteID, certificateID)
	d.SetId(syntheticID)
	log.Printf("[INFO] Created Incapsula Site to Imperva to Origin mutual TLS Certificate Association with ID: %s - site ID (%d) - certificate ID (%d)", syntheticID, siteID, certificateID)
	return nil
}

func resourceSiteMtlsClientToImpervaCertificateAssociationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID, certificateID, err := validateInput(d)
	if err != nil {
		return err
	}

	err = client.CreateSiteMtlsClientToImpervaCertificateAssociation(
		certificateID,
		siteID,
	)
	if err != nil {
		return err
		//todo -add error message
	}
	return resourceSiteMtlsClientToImpervaCertificateAssociationRead(d, m)
}

func resourceSiteMtlsClientToImpervaCertificateAssociationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)
	siteID, certificateID, err := validateInput(d)

	err = client.DeleteSiteMtlsClientToImpervaCertificateAssociation(
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
