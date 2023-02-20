package incapsula

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const accountResourceUserType = "incapsula_account_user"
const accountResourceUserName = "test-terraform-account-user"
const accountResourceUserTypeName = accountResourceUserType + "." + accountResourceUserName
const accountUserEmail = "test-terraform@incaptest.com"
const accountUserEmailSpecialChar = "test1+terraform@incaptest.com"
const accountUserFirstName = "First"
const accountUserLastName = "Last"

func TestIncapsulaAccountUser_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountUserConfigBasic(t, accountUserEmail),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountUserExists(accountResourceUserTypeName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "email", accountUserEmail),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "first_name", accountUserFirstName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "last_name", accountUserLastName),
				),
			},
		},
	})
}

func TestIncapsulaAccountUser_BasicWithSpecialChar(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountUserConfigBasic(t, accountUserEmailSpecialChar),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountUserExists(accountResourceUserTypeName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "email", accountUserEmailSpecialChar),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "first_name", accountUserFirstName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "last_name", accountUserLastName),
				),
			},
		},
	})
}

func TestIncapsulaAccountUser_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountUserConfigBasic(t, accountUserEmail),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountUserExists(accountResourceUserTypeName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "email", accountUserEmail),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "role_ids.#", "0"),
				),
			},
			{
				Config: testCheckIncapsulaAccountUserConfigUpdate(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaAccountUserExists(accountResourceUserTypeName),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "email", accountUserEmail),
					resource.TestCheckResourceAttr(accountResourceUserTypeName, "role_ids.#", "1"),
				),
			},
		},
	})
}

func TestIncapsulaAccountUser_ImportBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaAccountUserConfigBasic(t, accountUserEmail),
			},
			{
				ResourceName:      accountResourceUserTypeName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateUserConfigID,
			},
		},
	})
}

func testACCStateUserConfigID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		fmt.Sprintf("Resource: %v", rs)
		if rs.Type != accountResourceUserType {
			continue
		}

		keyParts := strings.Split(rs.Primary.ID, "/")
		if len(keyParts) != 2 {
			return "", fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", keyParts)
		}
		keyAccountID, err := strconv.Atoi(keyParts[0])
		if err != nil {
			return "", fmt.Errorf("failed to convert account ID, actual value: %s, expected numeric id", keyParts[0])
		}
		keyEmail := keyParts[1]
		resourceID := fmt.Sprintf("%d/%s", keyAccountID, keyEmail)

		schemaAccountID, err := strconv.Atoi(rs.Primary.Attributes["account_id"])
		schemaEmail := rs.Primary.Attributes["email"]
		newID := fmt.Sprintf("%d/%s", schemaAccountID, schemaEmail)

		if strings.Compare(newID, resourceID) != 0 {
			// if newID != resourceID {
			return "", fmt.Errorf("Incapsula Account User does not exist")
		}
		return resourceID, nil
	}
	return "", fmt.Errorf("Error finding correct resource %s", accountResourceUserType)
}

func testCheckIncapsulaAccountUserDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != accountResourceUserType {
			continue
		}

		accountUserIDStr := res.Primary.ID
		if accountUserIDStr == "" {
			return fmt.Errorf("Incapsula account user ID does not exist")
		}

		keyParts := strings.Split(accountUserIDStr, "/")
		if len(keyParts) != 2 {
			return fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", accountUserIDStr)
		}
		keyAccountID, err := strconv.Atoi(keyParts[0])
		if err != nil {
			return fmt.Errorf("failed to convert account ID, actual value: %s, expected numeric id", keyParts[0])
		}
		keyEmail := keyParts[1]

		accountUserStatusResponse, err := client.GetAccountUser(keyAccountID, keyEmail)

		// Account object may have been deleted
		if accountUserStatusResponse != nil {
			return fmt.Errorf(
				"Incapsula account user with email: %s for account ID %d still exists",
				keyEmail,
				keyAccountID,
			)
		}
	}

	return nil
}

func testCheckIncapsulaAccountUserExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Incapsula account user resource not found: %s", name)
		}

		accountUserIDStr := res.Primary.ID
		if accountUserIDStr == "" {
			return fmt.Errorf("Incapsula account user ID does not exist")
		}

		keyParts := strings.Split(accountUserIDStr, "/")
		if len(keyParts) != 2 {
			return fmt.Errorf("Error parsing ID, actual value: %s, expected numeric id and string seperated by '/'\n", accountUserIDStr)
		}
		keyAccountID, err := strconv.Atoi(keyParts[0])
		if err != nil {
			return fmt.Errorf("failed to convert account ID, actual value: %s, expected numeric id", keyParts[0])
		}
		keyEmail := keyParts[1]

		client := testAccProvider.Meta().(*Client)
		accountUserStatusResponse, err := client.GetAccountUser(keyAccountID, keyEmail)
		if accountUserStatusResponse == nil {
			return fmt.Errorf(
				"Incapsula account user with email: %s for account ID %d does not exist",
				keyEmail,
				keyAccountID,
			)
		}

		return nil
	}
}

func testCheckIncapsulaAccountUserConfigBasic(t *testing.T, email string) string {
	return fmt.Sprintf(`
		data "incapsula_account_data" "account_data" {}

		resource "%s" "%s" {
			account_id = data.incapsula_account_data.account_data.current_account
			email = "%s"
			first_name = "%s"
			last_name = "%s"
		}`,
		accountResourceUserType, accountResourceUserName, email, accountUserFirstName, accountUserLastName,
	)
}

func testCheckIncapsulaAccountUserConfigUpdate(t *testing.T) string {
	return fmt.Sprintf(`
		data "incapsula_account_data" "account_data" {}

		data "incapsula_account_roles" "roles" {
		  account_id = data.incapsula_account_data.account_data.current_account
		}

		resource "%s" "%s" {
			account_id = data.incapsula_account_data.account_data.current_account
			email = "%s"
			first_name = "%s"
			last_name = "%s"
			role_ids = [
				data.incapsula_account_roles.roles.reader_role_id
			]
		}`,
		accountResourceUserType, accountResourceUserName, accountUserEmail, accountUserFirstName, accountUserLastName,
	)
}
