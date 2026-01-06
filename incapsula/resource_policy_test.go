package incapsula

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const policyResourceType = "incapsula_policy"
const policyResourceName = "testacc-terraform-policy-"
const policyResourceTypeAndName = policyResourceType + "." + policyResourceName

const aclPolicyName = "acl-policy-test"
const fileUploadPolicyName = "file-upload-policy-test"

const wafPolicySettings = "[\n " +
	"   {\n" +
	"      \"settingsAction\": \"BLOCK\",\n" +
	"      \"policySettingType\": \"REMOTE_FILE_INCLUSION\"\n" +
	"    },\n" +
	"    {\n" +
	"      \"settingsAction\": \"BLOCK\",\n" +
	"      \"policySettingType\": \"ILLEGAL_RESOURCE_ACCESS\"\n" +
	"    },\n" +
	"    {\n" +
	"      \"settingsAction\": \"BLOCK\",\n" +
	"      \"policySettingType\": \"CROSS_SITE_SCRIPTING\"\n" +
	"    },\n" +
	"    {\n" +
	"      \"settingsAction\": \"BLOCK\",\n" +
	"      \"policySettingType\": \"SQL_INJECTION\"\n" +
	"    }\n" +
	"]"

const aclPolicySettingsUrlExceptions = "[\n" +
	"    {\n" +
	"        \"settingsAction\": \"BLOCK\",\n" +
	"        \"policySettingType\": \"IP\",\n" +
	"        \"data\": {\n" +
	"            \"ips\": [\n" +
	"                \"10.10.10.10\",\n" +
	"                \"10.10.10.11\",\n" +
	"                \"10.10.10.24\"\n" +
	"            ]\n" +
	"        },\n" +
	"        \"policyDataExceptions\": [\n" +
	"            {\n" +
	"                \"data\": [\n" +
	"                    {\n" +
	"                        \"exceptionType\": \"CLIENT_ID\",\n" +
	"                        \"values\": [\n" +
	"                            \"144\"\n" +
	"                        ]\n" +
	"                    },\n" +
	"                    {\n" +
	"                        \"exceptionType\": \"IP\",\n" +
	"                        \"values\": [\n" +
	"                            \"10.10.192.10\"\n" +
	"                        ]\n" +
	"                    }\n" +
	"                ],\n" +
	"                \"comment\": \"Adding first exception to policy settings\"\n" +
	"            }\n" +
	"        ]\n" +
	"    },\n" +
	"    {\n" +
	"        \"settingsAction\": \"BLOCK\",\n" +
	"        \"policySettingType\": \"GEO\",\n" +
	"        \"data\": {\n" +
	"            \"geo\": {\n" +
	"                \"countries\": [\n" +
	"                    \"WF\",\n" +
	"                    \"AD\"\n" +
	"                ],\n" +
	"                \"continents\": [\n" +
	"                    \"AS\",\n" +
	"                    \"EU\"\n" +
	"                ]\n" +
	"            }\n" +
	"        }\n" +
	"    }\n" +
	"]"

func TestAccIncapsulaPolicy_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaPolicyConfigBasic(t, aclPolicyName, true, "ACL", aclPolicySettingsUrlExceptions),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyExists(policyResourceTypeAndName+aclPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+aclPolicyName, "name", aclPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+aclPolicyName, "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+aclPolicyName, "policy_settings", aclPolicySettingsUrlExceptions),
				),
			},
			{
				ResourceName:      policyResourceTypeAndName + aclPolicyName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStatePolicyID,
			},
		},
	})
}

func TestAccIncapsulaFileUploadPolicy(t *testing.T) {
	const FileUploadPolicySettingsHashException = "[\n" +
		"    {\n" +
		"        \"settingsAction\": \"BLOCK\",\n" +
		"        \"policySettingType\": \"MALICIOUS_FILE_UPLOAD\",\n" +
		"        \"data\": {},\n" +
		"        \"policyDataExceptions\": [\n" +
		"            {\n" +
		"                \"data\": [\n" +
		"                    {\n" +
		"                        \"exceptionType\": \"FILE_HASH\",\n" +
		"                        \"values\": [\n" +
		"                            \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"\n" +
		"                        ]\n" +
		"                    }\n" +
		"                ],\n" +
		"                \"comment\": \"Adding first exception to policy settings\"\n" +
		"            }\n" +
		"        ]\n" +
		"    }\n" +
		"]"

	const FileUploadPolicySettingsHashAndIPException = "[\n" +
		"    {\n" +
		"        \"settingsAction\": \"BLOCK\",\n" +
		"        \"policySettingType\": \"MALICIOUS_FILE_UPLOAD\",\n" +
		"        \"data\": {},\n" +
		"        \"policyDataExceptions\": [\n" +
		"            {\n" +
		"                \"data\": [\n" +
		"                    {\n" +
		"                        \"exceptionType\": \"FILE_HASH\",\n" +
		"                        \"values\": [\n" +
		"                            \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"\n" +
		"                        ]\n" +
		"                    },\n" +
		"                    {\n" +
		"                        \"exceptionType\": \"IP\",\n" +
		"                        \"values\": [\n" +
		"                            \"10.10.192.10\"\n" +
		"                        ]\n" +
		"                    }\n" +
		"                ],\n" +
		"                \"comment\": \"Updating exception in policy settings\"\n" +
		"            }\n" +
		"        ]\n" +
		"    }\n" +
		"]"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaPolicyConfigBasic(t, fileUploadPolicyName, true, "FILE_UPLOAD", FileUploadPolicySettingsHashException),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyExists(policyResourceTypeAndName+fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "name", fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "policy_settings", FileUploadPolicySettingsHashException),
				),
			},
			{
				Config: testAccCheckIncapsulaPolicyConfigBasic(t, fileUploadPolicyName, true, "FILE_UPLOAD", FileUploadPolicySettingsHashAndIPException),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyExists(policyResourceTypeAndName+fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "name", fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "enabled", strconv.FormatBool(true)),
					resource.TestMatchResourceAttr(
						policyResourceTypeAndName+fileUploadPolicyName,
						"policy_settings",
						regexp.MustCompile(regexp.QuoteMeta(FileUploadPolicySettingsHashAndIPException)+`\n?`),
					),
				),
			},
			{
				ResourceName:      policyResourceTypeAndName + fileUploadPolicyName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStatePolicyID,
			},
		},
	})
}

func TestAccIncapsulaFileUploadPolicyWithoutExceptionsAndUpdate(t *testing.T) {
	const FileUploadPolicySettings = "[\n" +
		"    {\n" +
		"        \"settingsAction\": \"BLOCK\",\n" +
		"        \"policySettingType\": \"MALICIOUS_FILE_UPLOAD\",\n" +
		"        \"data\": {}\n" +
		"    }\n" +
		"]"

	const FileUploadPolicySettingsHashException = "[\n" +
		"    {\n" +
		"        \"settingsAction\": \"BLOCK\",\n" +
		"        \"policySettingType\": \"MALICIOUS_FILE_UPLOAD\",\n" +
		"        \"data\": {},\n" +
		"        \"policyDataExceptions\": [\n" +
		"            {\n" +
		"                \"data\": [\n" +
		"                    {\n" +
		"                        \"exceptionType\": \"FILE_HASH\",\n" +
		"                        \"values\": [\n" +
		"                            \"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\"\n" +
		"                        ]\n" +
		"                    }\n" +
		"                ],\n" +
		"                \"comment\": \"Adding first exception to policy settings\"\n" +
		"            }\n" +
		"        ]\n" +
		"    }\n" +
		"]"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaPolicyConfigBasic(t, fileUploadPolicyName, true, "FILE_UPLOAD", FileUploadPolicySettings),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyExists(policyResourceTypeAndName+fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "name", fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "enabled", strconv.FormatBool(true)),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "policy_settings", FileUploadPolicySettings),
				),
			},
			{
				Config: testAccCheckIncapsulaPolicyConfigBasic(t, fileUploadPolicyName, true, "FILE_UPLOAD", FileUploadPolicySettingsHashException),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaPolicyExists(policyResourceTypeAndName+fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "name", fileUploadPolicyName),
					resource.TestCheckResourceAttr(policyResourceTypeAndName+fileUploadPolicyName, "enabled", strconv.FormatBool(true)),
					resource.TestMatchResourceAttr(
						policyResourceTypeAndName+fileUploadPolicyName,
						"policy_settings",
						regexp.MustCompile(regexp.QuoteMeta(FileUploadPolicySettingsHashException)+`\n?`),
					),
				),
			},
			{
				ResourceName:      policyResourceTypeAndName + fileUploadPolicyName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStatePolicyID,
			},
		},
	})
}

func testAccStatePolicyID(state *terraform.State) (string, error) {
	for _, rs := range state.RootModule().Resources {
		if rs.Type != policyResourceType {
			continue
		}

		policyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		return fmt.Sprintf("%d", policyID), nil
	}

	return "", fmt.Errorf("Error finding Policy ID")
}

func testCheckIncapsulaPolicyExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula Policy resource not found: %s", name)
		}
		policyIDStr := res.Primary.ID
		if policyIDStr == "" {
			return fmt.Errorf("Incapsula policy ID does not exist")
		}

		_, err := strconv.Atoi(policyIDStr)
		if err != nil {
			return fmt.Errorf("Policy ID conversion error for %s: %s", policyIDStr, err)
		}
		client := testAccProvider.Meta().(*Client)
		_, err = client.GetPolicy(policyIDStr, nil)
		if err != nil {
			return fmt.Errorf("Get Incapsula Policy return error %s", err)
		}
		return nil
	}
}

func testAccCheckIncapsulaPolicyConfigBasic(t *testing.T, policyName string, enabled bool, policyType string, policySettings string) string {
	return testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + createPolicyResourceString(policyName, enabled, policyType, policySettings)
}

func createPolicyResourceString(policyName string, enabled bool, policyType string, policySettings string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
    name        = "%s"
    enabled     = %s
    policy_type = "%s"
    policy_settings = <<POLICY
%s 
POLICY
}`, policyResourceType, policyResourceName+policyName, policyName, strconv.FormatBool(enabled), policyType, policySettings,
	)
}

func testAccCheckIncapsulaPolicyDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_policy" {
			continue
		}

		policyId := res.Primary.ID
		if policyId == "" {
			// There is a bug in Terraform: https://github.com/hashicorp/terraform/issues/23635
			// Specifically, upgrades/destroys are happening simultaneously and not honoring
			// dependencies. In this case, it's possible that the site has already been deleted,
			// which means that all the sub resources will have been cleared out.
			// Ordinarily, this should return an error, but until this gets addressed, we're
			// going to simply return nil.
			// return fmt.Errorf("Incapsula Policy ID does not exist")
			return nil
		}

		getPolicyResponse, _ := client.GetPolicy(policyId, nil)

		if getPolicyResponse != nil {
			return fmt.Errorf("Incapsula Policy with name: %s (id: %d) still exists", getPolicyResponse.Value.Name, getPolicyResponse.Value.ID)
		}
	}
	return nil
}
