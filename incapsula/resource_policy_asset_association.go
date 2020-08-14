package incapsula

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePolicyAssetAssociation() *schema.Resource {
	return &schema.Resource{
		Create:   resourcePolicyAssetAssociationCreate,
		Read:     resourcePolicyAssetAssociationNil,
		Update:   resourcePolicyAssetAssociationNil,
		Delete:   resourcePolicyAssetAssociationDelete,
		Importer: nil,

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"policy_id": {
				Description: "The Policy ID for the asset association.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"asset_id": {
				Description: "The Asset ID for the asset association. Only type of asset supported at the moment is site.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"asset_type": {
				Description: "The Policy type for the asset association. Only value at the moment is `WEBSITE`.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourcePolicyAssetAssociationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := d.Get("policy_id").(string)
	assetID := d.Get("asset_id").(string)
	assetType := d.Get("asset_type").(string)

	err := client.AddPolicyAssetAssociation(policyID, assetID, assetType)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula policy asset association: policy ID (%s) - asset ID (%s) - asset type (%s) - %s\n", policyID, assetID, assetType, err)
		return err
	}

	// Generate synthetic ID
	syntheticID := fmt.Sprintf("%s-%s-%s", policyID, assetID, assetType)
	d.SetId(syntheticID)
	log.Printf("[INFO] Created Incapsula policy asset association with ID: %s - policy ID (%s) - asset ID (%s) - asset type (%s)\n", syntheticID, policyID, assetID, assetType)

	return nil
}

func resourcePolicyAssetAssociationNil(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePolicyAssetAssociationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := d.Get("policy_id").(string)
	assetID := d.Get("asset_id").(string)
	assetType := d.Get("asset_type").(string)

	err := client.DeletePolicyAssetAssociation(policyID, assetID, assetType)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
