package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"strings"
	"testing"
)

func TestAccIncapsulaATOSiteAllowlistConfig_basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_ato_site_allowlist_configuration_test.TestAccIncapsulaATOSiteAllowlistConfig_basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateATOSiteAllowlistConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckATOSiteAllowlistConfigBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckATOSiteAllowlistConfigExists(atoSiteAllowlistConfigResource),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.0.ip", "192.10.20.0"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.0.mask", "24"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.0.desc", "Test IP 1"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.0.updated", "1632530998076"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.1.ip", "192.10.20.1"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.1.mask", "8"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.1.desc", "Test IP 2"),
					resource.TestCheckResourceAttr(atoSiteAllowlistConfigResource, "allowlist.1.updated", "1632530998077"),
				),
			},
			{
				Config:            testAccCheckATOSiteAllowlistConfigBasic(t),
				ResourceName:      atoSiteAllowlistConfigResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateATOSiteAllowlistID,
			},
		},
	})
}

func testCheckATOSiteAllowlistConfigExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		// The ATO site configuration takes upto a minute to deploy worldwide. To keep the tests consistent we wait
		//time.Sleep(60 * time.Second)

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula ATO Site allowlist resource not found: %s", name)
		}

		// Extract accountId and siteId from teh terraform state
		var accountIdString = res.Primary.Attributes["account_id"]
		var accountId int
		if accountIdString != "" {
			accountIdInt, err := strconv.Atoi(accountIdString)
			if err != nil {
				return fmt.Errorf("failed to convert account ID from import command, actual value: %s, expected numeric ID", accountIdString)
			}
			accountId = accountIdInt
		}
		var siteIdString = res.Primary.Attributes["site_id"]
		siteId, err := strconv.Atoi(siteIdString)
		if err != nil {
			fmt.Errorf("failed to convert site ID from import command, actual value: %s, expected numeric ID", siteIdString)
		}

		client := testAccProvider.Meta().(*Client)
		aTOAllowlistDTO, err := client.GetAtoSiteAllowlistWithRetries(accountId, siteId)
		if err != nil {
			return fmt.Errorf("Error in fetching ATO allowlistg for site ID %d, Error : %s", siteId, err)
		}
		if aTOAllowlistDTO == nil || aTOAllowlistDTO.Allowlist == nil {
			return fmt.Errorf("ATO site allowlist is not present for site ID %d", siteId)
		}

		return nil
	}
}

func testACCStateATOSiteAllowlistID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != atoSiteAllowlistResourceType {
			continue
		}

		resourceId := rs.Primary.ID
		schemaId := rs.Primary.Attributes["site_id"]

		if strings.Compare(schemaId, resourceId) != 0 {
			// if newID != resourceID {
			return "", fmt.Errorf("Incapsula ATO Site allowlist Config does not exist")
		}

		return schemaId, nil
	}
	return "", fmt.Errorf("Error finding correct resource %s", atoSiteAllowlistConfigResource)
}

func testACCStateATOSiteAllowlistConfigDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != atoSiteAllowlistResourceType {
			continue
		}

		siteId := rs.Primary.Attributes["site_id"]
		// If site id is not explicitly specified, attempt to extract it from the resource ID
		if siteId == "" {
			siteId = rs.Primary.ID

		}
		siteIdInt, err := strconv.Atoi(siteId)
		if err != nil {
			return fmt.Errorf("failed to convert site ID from import command, actual value : %s, expected numeric id", siteId)
		}
		accountId := rs.Primary.Attributes["account_id"]
		accountIdInt, _ := strconv.Atoi(accountId)

		atoAllowlistDTO, _ := client.GetAtoSiteAllowlist(accountIdInt, siteIdInt)
		if err != nil && atoAllowlistDTO != nil && atoAllowlistDTO.Allowlist != nil && len(atoAllowlistDTO.Allowlist) > 0 {
			return fmt.Errorf("resource %s for ATO site allowlist : Api Id %s, site ID %d still exists", atoSiteAllowlistResourceType, rs.Primary.ID, siteIdInt)
		}
		return nil
	}
	return nil
}

func testAccCheckATOSiteAllowlistConfigBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
		site_id             = %s.id
		allowlist			= [ { "ip": "192.10.20.0", "mask": "24", "desc": "Test IP 1", "updated": 1632530998076 }, { "ip": "192.10.20.1", "mask": "8", "desc": "Test IP 2", "updated": 1632530998077 } ]
		depends_on 			= ["%s"]
	}`,
		atoSiteAllowlistResourceType, atoSiteAllowlistResourceName, siteResourceName, siteResourceName,
	)
}
