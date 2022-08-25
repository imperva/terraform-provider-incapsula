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
				idSlice := strings.Split(d.Id(), "/")
				if len(idSlice) != 3 || idSlice[0] == "" || idSlice[1] == "" || idSlice[2] == "" {
					return nil, fmt.Errorf("unexpected format of Incapsula Client to Imperva CA Certificate Site Association resource ID, expected account_idsite_id/certificate_id, got %s", d.Id())
				}

				_, err := strconv.Atoi(idSlice[1])
				if err != nil {
					fmt.Errorf("failed to convert Site Id from import command, actual value: %s, expected numeric id", idSlice[1])
				}

				_, err = strconv.Atoi(idSlice[0])
				if err != nil {
					fmt.Errorf("failed to convert Account Id from import command, actual value: %s, expected numeric id", idSlice[0])
				}

				_, err = strconv.Atoi(idSlice[2])
				if err != nil {
					fmt.Errorf("failed to convert Certificate Id from import command, actual value: %s, expected numeric id", idSlice[2])
				}

				d.Set("account_id", idSlice[0])
				d.Set("site_id", idSlice[1])
				d.Set("certificate_id", idSlice[2])

				log.Printf("[DEBUG] Importing Incapsula Client to Imperva CA Certificate Site Association for Site ID %s, mutual TLS Certificate Id %s, account ID %s", idSlice[1], idSlice[2], idSlice[0])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"account_id": {
				Description: "Account ID to operate on",
				Type:        schema.TypeString,
				Required:    true,
			},
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
		return fmt.Errorf("Couldn't find the Incapsula Client to Imperva CA Certificate Site Association. Site Id %d, accountID %s, certificate ID %d", siteID)
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
	//todo KATRIN - do we want to call read here?
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
