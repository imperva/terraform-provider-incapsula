package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const apiSecEndpointConfigResourceName = "incapsula_api_security_endpoint_config"
const apiSecEndpointConfigResource = apiSecEndpointConfigResourceName + "." + apiSecEndpointConfigName
const apiSecEndpointConfigName = "testacc-terraform-api-security-endpoint-config"

func TestAccIncapsulaApiSecurityEndpointConfig_Basic(t *testing.T) {
	log.Printf("========================BEGINTEST========================")
	log.Printf("[DEBUG]Running test resource_api_security_endpoint_config_test.TestAccIncapsulaApiSecurityEndpoint_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaApiSecurityEndpointConfigBasic(t), //toda ad path + method
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaApiSecurityEndpointExists(apiSecEndpointConfigResource),
					resource.TestCheckResourceAttr(apiSecEndpointConfigResource, "invalid_param_name_violation_action", "IGNORE"),
					resource.TestCheckResourceAttr(apiSecEndpointConfigResource, "invalid_param_value_violation_action", "DEFAULT"),
					resource.TestCheckResourceAttr(apiSecEndpointConfigResource, "path", "/users"),
					resource.TestCheckResourceAttr(apiSecEndpointConfigResource, "method", "GET"),
					resource.TestCheckResourceAttr(apiSecEndpointConfigResource, "missing_param_violation_action", "BLOCK_IP"),
				),
			},
			{
				ResourceName:      apiSecEndpointConfigResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateApiSecurityEndpointID,
			},
		},
	})
}

func testAccStateApiSecurityEndpointID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != apiSecEndpointConfigResourceName {
			continue
		}
		apiId, err := strconv.Atoi(rs.Primary.Attributes["api_id"])
		if err != nil {
			fmt.Errorf("Failed to convert API Id,actual value:%s, expected numeric id", rs.Primary.Attributes["api_id"])
		}

		endpointId, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			fmt.Errorf("Failed to convert API Id,actual value:%s, expected numeric id", rs.Primary.Attributes["api_id"])
		}

		return fmt.Sprintf("%d/%d", apiId, endpointId), nil
	}

	return "", fmt.Errorf("Error finding API Security Endpoint ID")
}

func testCheckIncapsulaApiSecurityEndpointExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("incapsula api security endpoint resource not found : %s", name)
		}

		endpointId := res.Primary.ID
		if endpointId == "" {
			return fmt.Errorf("incapsula api security endpoint ID does not exist")
		}

		apiId := res.Primary.Attributes["api_id"]
		if apiId == "" {
			return fmt.Errorf("incapsula api security endpoint ID does not exist")
		}
		apiIdInt, err := strconv.Atoi(apiId)
		if err != nil {
			return fmt.Errorf("failed to convert api security API ID is not numeric")
		}

		client := testAccProvider.Meta().(*Client)
		endpointListResponse, err := client.GetApiSecurityEndpointConfig(apiIdInt, endpointId)
		if err != nil {
			return fmt.Errorf("Incapsula Api Security Endpoint doesn't exist")
		}

		if endpointListResponse.Value.Id != 0 {
			return nil
		}
		return nil
	}
}

func testAccCheckIncapsulaApiSecurityEndpointConfigBasic(t *testing.T) string {
	return testAccCheckApiConfigBasic(t) + fmt.Sprintf(`
resource"%s""%s"{
api_id=%s.id
path = "/users"
method = "GET"
invalid_param_name_violation_action="IGNORE"
invalid_param_value_violation_action="DEFAULT"
missing_param_violation_action="BLOCK_IP"
depends_on=[%s]
}
`, apiSecEndpointConfigResourceName, apiSecEndpointConfigName, apiSecApiConfigResource, apiSecApiConfigResource,
	)
}
