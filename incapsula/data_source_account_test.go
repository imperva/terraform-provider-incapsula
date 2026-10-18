package incapsula

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestFindAccountsByName(t *testing.T) {
	accounts := []SubAccountItem{
		{SubAccountID: 1, SubAccountName: "Prod"},
		{SubAccountID: 2, SubAccountName: "prod"},
		{SubAccountID: 3, SubAccountName: "Staging"},
	}

	if got := findAccountsByName(accounts, "staging"); len(got) != 1 || got[0].SubAccountID != 3 {
		t.Errorf("Expected case-insensitive match of ID 3, got %+v", got)
	}
	if got := findAccountsByName(accounts, "PROD"); len(got) != 2 {
		t.Errorf("Expected 2 matches, got %+v", got)
	}
	if got := findAccountsByName(accounts, "Staging "); len(got) != 0 {
		t.Errorf("Input must not be trimmed, got %+v", got)
	}
	if got := findAccountsByName(accounts, "missing"); len(got) != 0 {
		t.Errorf("Expected no matches, got %+v", got)
	}
}

const accountNamedDataSourceType = "incapsula_account"
const accountNamedDataSourceName = "account_by_name"
const accountNamedAccountName = "acceptance-account-data-source-test-1"

func TestAccIncapsulaDataSourceAccountByName_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test data_source_account_test.go.TestAccIncapsulaDataSourceAccountByName_Basic")

	dataSource := "data." + accountNamedDataSourceType + "." + accountNamedDataSourceName
	account := "incapsula_account.test-terraform-account"
	email := GenerateTestEmail(t)

	resource.Test(t, resource.TestCase{
		ErrorCheck:   SkipIfAccountTypeIsResellerEndUser(t),
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaAccountDestroy,
		Steps: []resource.TestStep{
			{
				// Look the account up with a different casing than the one it was created with
				Config: testAccCheckIncapsulaDataSourceAccountByNameConfig(email, strings.ToUpper(accountNamedAccountName)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSource, "id", account, "id"),
					resource.TestCheckResourceAttrPair(dataSource, "account_id", account, "id"),
					resource.TestCheckResourceAttr(dataSource, "account_name", accountNamedAccountName),
					resource.TestCheckResourceAttrSet(dataSource, "parent_id"),
				),
			},
		},
	})
}

func TestAccIncapsulaDataSourceAccountByName_NotFound(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test data_source_account_test.go.TestAccIncapsulaDataSourceAccountByName_NotFound")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "%s" "%s" {
						account_name = "acceptance-account-that-does-not-exist"
					}`, accountNamedDataSourceType, accountNamedDataSourceName),
				ExpectError: regexp.MustCompile("No account named"),
			},
		},
	})
}

func testAccCheckIncapsulaDataSourceAccountByNameConfig(email string, lookupName string) string {
	return fmt.Sprintf(`
		resource "incapsula_account" "test-terraform-account" {
			email        = "%s"
			account_name = "%s"
			plan_id      = "entTrial"
			support_all_tls_versions           = "false"
			naked_domain_san_for_new_www_sites = "true"
		}

		data "%s" "%s" {
			account_name = "%s"
			depends_on   = [incapsula_account.test-terraform-account]
		}`,
		email, accountNamedAccountName,
		accountNamedDataSourceType, accountNamedDataSourceName, lookupName,
	)
}
