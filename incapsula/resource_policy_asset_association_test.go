package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

const policyAssetAssociationResourceType = "incapsula_policy_asset_association"
const policyAssetAssociationResourceName = "testacc-terraform-asset-policy-association"
const policyAssetAssociationResourceTypeAndName = policyAssetAssociationResourceType + "." + policyAssetAssociationResourceName
const assetTypeWebsite = "WEBSITE"

func TestAccIncapsulaPolicyAssetAssociation_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaPolicyAssetAssociationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaPolicyAssetAssociationConfigBasic(t, "${incapsula_site.testacc-terraform-site.id}", assetTypeWebsite),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyAssetAssociationExists(policyAssetAssociationResourceTypeAndName),
					resource.TestCheckResourceAttr(policyAssetAssociationResourceTypeAndName, "asset_type", assetTypeWebsite),
				),
			},
			{
				ResourceName:      policyAssetAssociationResourceTypeAndName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStatePolicyAssetAssociationID,
			},
		},
	})
}

func testAccStatePolicyAssetAssociationID(state *terraform.State) (string, error) {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != policyAssetAssociationResourceType {
			continue
		}

		return fmt.Sprintf("%s", rs.Primary.ID), nil
	}

	return "", fmt.Errorf("Error finding Policy Asset Association ID")
}

func testCheckIncapsulaPolicyAssetAssociationExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Policy Asset Association resource not found: %s", name)
		}
		syntheticID := res.Primary.ID
		if syntheticID == "" {
			return fmt.Errorf("Incapsula policy Asset Association ID does not exist")
		}

		splittedIds := strings.Split(syntheticID, "/")
		policyId := splittedIds[0]
		assetID := splittedIds[1]
		assetType := splittedIds[2]
		client := testAccProvider.Meta().(*Client)
		isAssociated, err := client.isPolicyAssetAssociated(policyId, assetID, assetType, nil)
		if err != nil {
			return fmt.Errorf("Get Incapsula Policy Asset Association return error %s", err)
		}
		if !isAssociated {
			return fmt.Errorf("Asset is not associated with Incapsula Policy %s", policyId)
		}
		return nil
	}
}

func testAccCheckIncapsulaPolicyAssetAssociationConfigBasic(t *testing.T, asset_id string, assetTypeWebsite string) string {
	return testAccCheckIncapsulaPolicyConfigBasic(t, policyAssetAssociationResourceName, true, "ACL", aclPolicySettingsUrlExceptions) + fmt.Sprintf(`
		resource "%s" "%s" {
			policy_id    = "${%s.id}"
			asset_id     = "%s"
			asset_type   = "%s"
		}`, policyAssetAssociationResourceType, policyAssetAssociationResourceName, policyResourceTypeAndName+policyAssetAssociationResourceName, asset_id, assetTypeWebsite)
}

func testAccCheckIncapsulaPolicyAssetAssociationDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	var siteId string
	var policyId string
	for _, res := range state.RootModule().Resources {
		if res.Type == policyResourceType {
			policyId = res.Primary.ID
			if policyId == "" {
				// There is a bug in Terraform: https://github.com/hashicorp/terraform/issues/23635
				// Specifically, upgrades/destroys are happening simultaneously and not honoring
				// dependencies. In this case, it's possible that the site has already been deleted,
				// which means that all the sub resources will have been cleared out.
				// Ordinarily, this should return an error, but until this gets addressed, we're
				// going to simply return nil.
				// return fmt.Errorf("Incapsula policy ID does not exist")
				return nil
			}
		}
		if res.Type == "incapsula_site" {
			siteId = res.Primary.ID
			if siteId == "" {
				// There is a bug in Terraform: https://github.com/hashicorp/terraform/issues/23635
				// Specifically, upgrades/destroys are happening simultaneously and not honoring
				// dependencies. In this case, it's possible that the site has already been deleted,
				// which means that all the sub resources will have been cleared out.
				// Ordinarily, this should return an error, but until this gets addressed, we're
				// going to simply return nil.
				// return fmt.Errorf("Incapsula site ID does not exist")
				return nil
			}
		}
	}
	if siteId == "" {
		return fmt.Errorf("Failed to find incapsula Site")
	}
	if policyId == "" {
		return fmt.Errorf("Failed to find incapsula Policy")
	}

	isAssetAssociated, err := client.isPolicyAssetAssociated(policyId, siteId, assetTypeWebsite, nil)
	if err != nil {
		return fmt.Errorf("Get Incapsula Policy Asset Association return error %s", err)
	}
	if isAssetAssociated {
		return fmt.Errorf("Incapsula Site %s association with Policy (id: %s) still exists", siteId, policyId)
	}
	return nil
}
