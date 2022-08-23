package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
	"strings"
)

func resourceSiteMtlsClientToImpervaCertificateAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteMtlsClientToImpervaCertificateAssociationCreate,
		Read:   resourceSiteMtlsClientToImpervaCertificateAssociationRead,
		Update: resourceSiteMtlsClientToImpervaCertificateAssociationCreate,
		Delete: resourceSiteMtlsClientToImpervaCertificateAssociationDelete,
		//todo
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				//todo!!!!! KATRIN change all error messages
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 3 || idSlice[0] == "" || idSlice[1] == "" || idSlice[2] == "" {
					return nil, fmt.Errorf("unexpected format of Incapsula Site to Client CA to Imperva Certificate Association resource ID, expected site_id/certificate_id, got %s", d.Id())
				}

				d.Set("account_id", idSlice[0])
				d.Set("site_id", idSlice[1])
				d.Set("certificate_id", idSlice[2])

				log.Printf("[DEBUG] Importing Incapsula Site to Imperva to Origin mutual TLS Certificate Association for Site ID %s, mutual TLS Certificate Id %s", idSlice[0], idSlice[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Account ID", //todo add
				Type:        schema.TypeString,
				Required:    true,
			},
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

func resourceSiteMtlsClientToImpervaCertificateAssociationRead(d *schema.ResourceData, m interface{}) error {
	//// Implement by reading the ListCertificatesResponse for the data center
	client := m.(*Client)

	accountIDStr := d.Get("account_id").(string)

	accountID, err := strconv.Atoi(accountIDStr)
	if err != nil {
		return fmt.Errorf("failed to convert Account Id for Incapsula  Site to mutual TLS Client to Imperva Certificate association resource, actual value: %s, expected numeric id", accountIDStr)
	}

	siteID, certificateID, err := validateInput(d)
	if err != nil {
		return err
	}

	mTLSCertificateData, err := client.GetSiteMtlsClientToImpervaCertificateAssociation(accountID, siteID, certificateID)
	if err != nil {
		return err
	}

	if mTLSCertificateData == nil {
		//todo KATRIN - change error message
		return fmt.Errorf("Couldn't find the Incapsula Site to Imperva to Origin mutual TLS Certificate Association")
	}

	// Generate synthetic ID
	syntheticID := fmt.Sprintf("%d/%d/%d", accountID, siteID, certificateID)
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
	//todo KATRIN - do we want to call read?
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
