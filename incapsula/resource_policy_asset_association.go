package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
)

func resourcePolicyAssetAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicyAssetAssociationCreate,
		Read:   resourcePolicyAssetAssociationRead,
		Update: nil,
		Delete: resourcePolicyAssetAssociationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required Arguments
			"policy_id": {
				Description: "The Policy ID for the asset association.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"asset_id": {
				Description: "The Asset ID for the asset association. Only type of asset supported at the moment is site.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"asset_type": {
				Description: "The Policy type for the asset association. Only value at the moment is `WEBSITE`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			// Optional Arguments
			"account_id": {
				Description: "The Asset's Account ID",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourcePolicyAssetAssociationCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := d.Get("policy_id").(string)
	assetID := d.Get("asset_id").(string)
	assetType := d.Get("asset_type").(string)
	currentAccountId := d.Get("account_id").(int)

	err := client.AddPolicyAssetAssociation(policyID, assetID, assetType, &currentAccountId)

	if err != nil {
		log.Printf("[ERROR] Could not create Incapsula policy asset association: policy ID (%s) - asset ID (%s) - asset type (%s) - %s\n", policyID, assetID, assetType, err)
		return err
	}

	// Generate synthetic ID
	syntheticID := fmt.Sprintf("%s/%s/%s", policyID, assetID, assetType)
	d.SetId(syntheticID)
	log.Printf("[INFO] Created Incapsula policy asset association with ID: %s - policy ID (%s) - asset ID (%s) - asset type (%s)\n", syntheticID, policyID, assetID, assetType)

	return resourcePolicyAssetAssociationRead(d, m)
}

func resourcePolicyAssetAssociationRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := strings.Split(d.Id(), "/")[0]
	assetID := strings.Split(d.Id(), "/")[1]
	assetType := strings.Split(d.Id(), "/")[2]
	currentAccountId := getCurrentAccountId(d, client.accountStatus)
	if currentAccountId != nil {
		log.Printf("[INFO] Trying to read Incapsula Policy Asset Association: %s-%s-%s for account %d\n", policyID, assetID, assetType, *currentAccountId)
	} else {
		log.Printf("[INFO] Trying to read Incapsula Policy Asset Association: %s-%s-%s\n", policyID, assetID, assetType)
	}
	var isAssociated, err = client.isPolicyAssetAssociated(policyID, assetID, assetType, currentAccountId)

	if err != nil {
		log.Printf("[ERROR] Could not read Incapsula Policy Asset Association: %s-%s-%s, err: %s\n", policyID, assetID, assetType, err)
		return err
	}

	if !isAssociated {
		log.Printf("[ERROR] Could not find Incapsula Policy Asset Association: %s-%s-%s\n", policyID, assetID, assetType)
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] Successfully read Policy Asset Association exist: %s-%s-%s\n", policyID, assetID, assetType)
	syntheticID := fmt.Sprintf("%s/%s/%s", policyID, assetID, assetType)

	d.Set("asset_id", assetID)
	d.Set("asset_type", assetType)
	d.Set("policy_id", policyID)
	if currentAccountId != nil {
		d.Set("account_id", *currentAccountId)
	}
	d.SetId(syntheticID)

	return nil
}

func resourcePolicyAssetAssociationDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*Client)

	policyID := d.Get("policy_id").(string)
	assetID := d.Get("asset_id").(string)
	assetType := d.Get("asset_type").(string)
	currentAccountId := getCurrentAccountId(d, client.accountStatus)
	if currentAccountId != nil {
		log.Printf("[INFO] Trying to delete Incapsula Policy Asset Association: %s-%s-%s for account %d\n", policyID, assetID, assetType, *currentAccountId)
	} else {
		log.Printf("[INFO] Trying to delete Incapsula Policy Asset Association: %s-%s-%s\n", policyID, assetID, assetType)
	}
	err := client.DeletePolicyAssetAssociation(policyID, assetID, assetType, currentAccountId)

	if err != nil {
		return err
	}

	// Set the ID to empty
	// Implicitly clears the resource
	d.SetId("")

	return nil
}
