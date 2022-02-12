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

const txtRecordResourceName = "incapsula_txt_record"
const txtRecordResource = txtRecordResourceName + "." + textRecordName
const textRecordName = "testacc-terraform-txt-record"

func TestAccIncapsulaTxtRecord_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_txt_record.TestAccIncapsulaApiSecurityApiConfig_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateTXTRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckTxtRecordBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckTxtRecordExists(txtRecordResource),
					resource.TestCheckResourceAttr(txtRecordResource, "txt_record_value_one", "test1"),
					resource.TestCheckResourceAttr(txtRecordResource, "txt_record_value_two", "test2"),
					resource.TestCheckResourceAttr(txtRecordResource, "txt_record_value_three", "test3"),
					resource.TestCheckResourceAttr(txtRecordResource, "txt_record_value_four", "test4"),
					resource.TestCheckResourceAttr(txtRecordResource, "txt_record_value_five", "test5"),
				),
			},
			{
				ResourceName:      txtRecordResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateTxtRecordID,
			},
		},
	})
}

func testACCStateTXTRecordDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != txtRecordResource {
			continue
		}
		//return nil
		siteID := rs.Primary.ID

		//siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			fmt.Errorf("Parameter site_id was not found in resource %s", txtRecordResourceName)
		}
		siteIDInt, err := strconv.Atoi(siteID)
		if err != nil {
			fmt.Errorf("failed to convert Site Id from import command, actual value : %s, expected numeric id", siteID)
		}

		recordResponse, err := client.ReadTXTRecords(siteIDInt)
		if err == nil || strings.Contains(recordResponse.ResMessage, "no TXT records") {
			return fmt.Errorf("Resource %s for Incapsula TXT Record : site ID %d still exists", txtRecordResourceName, siteIDInt)
		}
	}
	return nil
}

func testACCStateTxtRecordID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != txtRecordResourceName {
			continue
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site ID %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d", siteID), nil
	}
	return "", fmt.Errorf("Error finding a TXT Record for SiteID, state")
}

func testCheckTxtRecordExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, _ := state.RootModule().Resources[name]
		//if !ok {
		//	return fmt.Errorf("Incapsula TXT Record resource not found: %s", name)
		//}
		siteId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing ID %v to int", res.Primary.ID)
		}

		client := testAccProvider.Meta().(*Client)
		recordResponse, err := client.ReadTXTRecords(siteId)
		if err != nil || strings.Contains(recordResponse.ResMessage, "no TXT records") {
			fmt.Errorf("Incapsula TXT Record doesn't exist")
		}

		return nil
	}
}

func testAccCheckTxtRecordBasic(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "incapsula_txt_record" "%s" {
		site_id = "${incapsula_site.testacc-terraform-site.id}"
        txt_record_value_one = "test1"
        txt_record_value_two = "test2"
        txt_record_value_three = "test3"
        txt_record_value_four = "test4"
        txt_record_value_five = "test5"
  		depends_on = ["%s"]
	}`,
		textRecordName, siteResourceName,
	)
}
