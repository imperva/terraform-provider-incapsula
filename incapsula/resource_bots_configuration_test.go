package incapsula

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const botsConfigurationResource = "incapsula_bots_configuration"
const botsConfigurationName = "testacc-terraform-bots-configuration"
const botsConfigurationResourceName = botsConfigurationResource + "." + botsConfigurationName

func TestAccIncapsulaBotsConfiguration_Basic(t *testing.T) {
	log.Printf("======================== BEGIN TEST ========================")
	log.Printf("[DEBUG] Running test resource_bots_configuration_test.TestAccIncapsulaBotsConfiguration_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaBotsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaBotsConfigurationBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaBotsConfigurationExists(botsConfigurationResourceName),
					resource.TestCheckResourceAttr(botsConfigurationResourceName, "canceledGoodBots", "[6, 2, 1]"),
					resource.TestCheckResourceAttr(botsConfigurationResourceName, "badBots", "[530, 20, 537]"),
				),
			},
			{
				ResourceName:      botsConfigurationResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateBotsConfigurationID,
			},
		},
	})
}

func testAccStateBotsConfigurationID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != botsConfigurationResource {
			continue
		}

		dccID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d", dccID), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func testAccCheckIncapsulaBotsConfigurationDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteID := res.Primary.ID
		if siteID == "" {
			// There is a bug in Terraform: https://github.com/hashicorp/terraform/issues/23635
			// Specifically, upgrades/destroys are happening simulatneously and not honoroing
			// dependencies. In this case, it's possible that the site has already been deleted,
			// which means that all of the subresources will have been cleared out.
			// Ordinarily, this should return an error, but until this gets addressed, we're
			// going to simply return nil.
			// return fmt.Errorf("Incapsula site ID does not exist")
			return nil
		}

		listBotResponse, _ := client.GetBotAccessControlConfiguration(siteID)

		// See comment above - the bot may have already been deleted
		// This workaround will be removed in the future
		if listBotResponse == nil || listBotResponse.Data == nil || len(listBotResponse.Data) == 0 {
			return nil
		}

		// Nothing to check here.
		// Destroying incapsula_bots_configuration should not have any backend effect.

	}

	return nil
}

func testCheckIncapsulaBotsConfigurationExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		siteRes, siteOk := state.RootModule().Resources[siteResourceName]
		if !siteOk {
			return fmt.Errorf("Incapsula site resource not found: %s", siteResourceName)
		}

		siteIDString := siteRes.Primary.ID
		if siteIDString == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		siteID, err := strconv.Atoi(siteIDString)
		if err != nil {
			return fmt.Errorf("Error parsing Site ID %v to int", siteIDString)
		}

		domain := siteRes.Primary.Attributes["domain"]

		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula bots configuration resource not found: %s", name)
		}

		dccID := res.Primary.ID
		if dccID == "" {
			return fmt.Errorf("Incapsula bots configuration ID does not exist")
		}

		client := testAccProvider.Meta().(*Client)

		// If the site has already been deleted then return nil
		// Otherwise check the bot list
		_, err = client.SiteStatus(domain, siteID)
		if err != nil {
			return nil
		}

		responseDTO, err := client.GetBotAccessControlConfiguration(siteIDString)
		if responseDTO == nil || responseDTO.Data == nil || len(responseDTO.Data) == 0 {
			return fmt.Errorf("Incapsula bots configuration: %s (Site ID: %d) does not exist\n%s", name, siteID, err)
		}

		if siteRes.Primary.ID != res.Primary.ID {
			return fmt.Errorf("The ID of Incapsula bots configuration: %s is invalid (ID: %s). "+
				"It should be identical to site_id (%s)", name, dccID, siteRes.Primary.ID)
		}

		return nil
	}
}

func testAccCheckIncapsulaBotsConfigurationBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = %s.id
  canceledGoodBots = [6, 2, 1]
  badBots = [530, 20, 537]
}`, botsConfigurationResource, botsConfigurationName, siteResourceName,
	)
}
