package incapsula

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const apiClientResourceType = "incapsula_api_client"
const apiClientResourceName = "example_api_client"
const apiClientName = "acceptance-api-client-test-1"

func TestAccIncapsulaAPIClient_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaAPIClientDestroy,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaAPIClientConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAPIClientExists(),
					resource.TestCheckResourceAttr(apiClientResourceType+"."+apiClientResourceName, "name", apiClientName),
				),
			},
			{
				ResourceName:      apiClientResourceType + "." + apiClientResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateAPIClientID,
			},
		},
	})
}

func getAccIncapsulaAPIClientConfigBasic() string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			name = "%s"
			enabled = true
			grace_period = 2000
			regenerate_version = 1.0
		}
	`,
		apiClientResourceType, apiClientResourceName, apiClientName,
	)
}

func testAccIncapsulaAPIClientDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	var apiClientID string
	for _, res := range state.RootModule().Resources {
		if res.Type != apiClientResourceType {
			continue
		}
		apiClientID = res.Primary.ID
		_, err := client.GetAPIClient(nil, apiClientID)
		if err == nil {
			return fmt.Errorf("incapsula API Client %s still exists", apiClientID)
		}
	}
	return nil
}

func testACCStateAPIClientID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != apiClientResourceType {
			continue
		}
		return rs.Primary.ID, nil
	}
	return "", fmt.Errorf("error finding API Client ID")
}

func testCheckIncapsulaAPIClientExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resource1 := apiClientResourceType + "." + apiClientResourceName
		res, ok := state.RootModule().Resources[resource1]
		if !ok {
			return fmt.Errorf("incapsula API Client resource not found : %s", apiClientResourceType)
		}
		apiID := res.Primary.ID
		if !ok || apiID == "" {
			return fmt.Errorf("incapsula API Client ID does not exist for API Client")
		}
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetAPIClient(nil, apiID)
		if err != nil {
			return err
		}
		return nil
	}
}
