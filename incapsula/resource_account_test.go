package incapsula

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testEmail = "example@imperva.com"
const accountResourceName = "incapsula_account.test-terraform-account"

func TestIncapsulaAccount_Basic(t *testing.T) {
	email := GenerateTestEmail(t)
	resource.Test(t, resource.TestCase{
		ErrorCheck:   SkipIfAccountTypeIsResellerEndUser(t),
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountConfigBasic(email),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountExists(accountResourceName),
					resource.TestCheckResourceAttr(accountResourceName, "email", email),
				),
			},
		},
	})
}

func TestIncapsulaAccount_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ErrorCheck:   SkipIfAccountTypeIsResellerEndUser(t),
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountConfigBasic(GenerateTestEmail(t)),
			},
			{
				ResourceName:      "incapsula_account.test-terraform-account",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateAccountID,
			},
		},
	})
}

func TestIncapsulaAccount_Http2Defaults(t *testing.T) {
	testIncapsulaAccountHttp2Client(t, true, false)
}

func TestIncapsulaAccount_Http2ClientAndOriginEnabled(t *testing.T) {
	testIncapsulaAccountHttp2Client(t, true, true)
}

func TestIncapsulaAccount_Http2ClientAndOriginDisabled(t *testing.T) {
	testIncapsulaAccountHttp2Client(t, false, false)
}

func testIncapsulaAccountHttp2Client(t *testing.T, enableHttp2ForNewSites bool, enableHttp2ToOriginForNewSites bool) {
	email := GenerateTestEmail(t)
	resource.Test(t, resource.TestCase{
		ErrorCheck:   SkipIfAccountTypeIsResellerEndUser(t),
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testHttp2AccountConfig(email, enableHttp2ForNewSites, enableHttp2ToOriginForNewSites),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountExists(accountResourceName),
					resource.TestCheckResourceAttr(accountResourceName, "email", email),
					resource.TestCheckResourceAttr(accountResourceName, "enable_http2_for_new_sites", strconv.FormatBool(enableHttp2ForNewSites)),
					resource.TestCheckResourceAttr(accountResourceName, "enable_http2_to_origin_for_new_sites", strconv.FormatBool(enableHttp2ToOriginForNewSites)),
				),
			},
		},
	})
}

func testCheckIncapsulaAccountDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_account" {
			continue
		}

		accountIDStr := res.Primary.ID
		if accountIDStr == "" {
			return fmt.Errorf("Incapsula account ID does not exist")
		}

		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return fmt.Errorf("Account ID conversion error for %s: %s", accountIDStr, err)
		}

		_, err = client.AccountStatus(accountID, ReadAccount)

		if err == nil {
			return fmt.Errorf("Incapsula account id: %d still exists", accountID)
		}
	}

	return nil
}

func testCheckIncapsulaAccountExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula account resource not found: %s", name)
		}

		accountIDStr := res.Primary.ID
		if accountIDStr == "" {
			return fmt.Errorf("Incapsula account ID does not exist")
		}

		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return fmt.Errorf("Account ID conversion error for %s: %s", accountIDStr, err)
		}

		client := testAccProvider.Meta().(*Client)
		accountStatusResponse, err := client.AccountStatus(accountID, ReadAccount)
		if accountStatusResponse == nil {
			return fmt.Errorf("Incapsula account id: %d does not exist", accountID)
		}

		return nil
	}
}

func testCheckIncapsulaAccountConfigBasic(email string) string {
	return fmt.Sprintf(`
			resource "incapsula_account" "test-terraform-account" {
			email = "%s"
            account_name = "testTerraform"
            plan_id = "entTrial"
            support_all_tls_versions = "false"
			naked_domain_san_for_new_www_sites = "true"
		}`,
		email,
	)
}

func testHttp2AccountConfig(email string, enableHttp2ForNewSites bool, enableHttp2ToOriginForNewSites bool) string {
	return fmt.Sprintf(`
		resource "incapsula_account" "test-terraform-account" {
			email = "%s"
			enable_http2_for_new_sites = "%t"
			enable_http2_to_origin_for_new_sites = "%t"
            account_name = "testTerraform"
            plan_id = "entTrial"
            support_all_tls_versions = "false"
		}`,
		email,
		enableHttp2ForNewSites,
		enableHttp2ToOriginForNewSites,
	)
}

func SkipIfAccountTypeIsResellerEndUser(t *testing.T) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}
		if strings.Contains(err.Error(), "Operation not allowed") {
			t.Skipf("skipping test since account type is RESELLER_END_USER. Error: %s", err.Error())
		}

		return err
	}
}

func GenerateTestEmail(t *testing.T) string {
	if v := os.Getenv("INCAPSULA_API_ID"); v == "" {
		t.Fatal("INCAPSULA_API_ID must be set for acceptance tests")
	}

	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	generatedDomain = "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000)) + testEmail

	return generatedDomain
}

func testACCStateAccountID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_account" {
			continue
		}
		accountID := rs.Primary.ID

		return accountID, nil
	}
	return "", fmt.Errorf("Error finding an Account\"")
}
