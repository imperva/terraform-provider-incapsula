package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"
	"testing"
)

const accountSSLSettingsResourceName = "incapsula_account_ssl_settings"
const accountSSLSettingsResource = accountSSLSettingsResourceName + "." + accountSSLSettingsConfigName
const accountSSLSettingsConfigName = "testacc-terraform-account-ssl-settings"

func TestAccAccountSSLSettings_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_account_ssl_settings_test.TestAccAccountSSLSettings_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccAccountSSLSettingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountSSLSettingsFullUpdateConfig(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAccountSSLSettingsAfterFullUpdate(),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "use_wild_card_san_instead_of_fqdn", "false"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "add_naked_domain_san_for_www_sites", "true"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "allow_cname_validation", "true"),
					resource.TestCheckTypeSetElemAttr(accountSSLSettingsResource, "allowed_domains_for_cname_validation.*", "example.com"),
				),
			},
			{
				Config: testAccAccountSSLSettingsPartialUpdate1Config(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAccountSSLSettingsAfterPartialUpdate1(),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "use_wild_card_san_instead_of_fqdn", "true"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "add_naked_domain_san_for_www_sites", "false"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "allow_cname_validation", "false"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "allowed_domains_for_cname_validation.#", "0"),
				),
			},
			{
				Config: testAccAccountSSLSettingsPartialUpdate2Config(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAccountSSLSettingsAfterPartialUpdate2(),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "use_wild_card_san_instead_of_fqdn", "true"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "add_naked_domain_san_for_www_sites", "true"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "allow_cname_validation", "true"),
					resource.TestCheckTypeSetElemAttr(accountSSLSettingsResource, "allowed_domains_for_cname_validation.*", "example.com"),
				),
			},
			{
				Config: testAccAccountSSLSettingsPartialUpdate3Config(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckAccAccountSSLSettingsAfterPartialUpdate3(),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "use_wild_card_san_instead_of_fqdn", "false"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "add_naked_domain_san_for_www_sites", "true"),
					resource.TestCheckResourceAttr(accountSSLSettingsResource, "allow_cname_validation", "false"),
					resource.TestCheckTypeSetElemAttr(accountSSLSettingsResource, "allowed_domains_for_cname_validation.*", "example2.com"),
				),
			},
			{
				ResourceName:      accountSSLSettingsResource,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateAccountSSLSettingsID,
			},
		},
	})
}

func testValueForCnameValidationNotNull(val string) error {
	if val == "" {
		return fmt.Errorf("value for cname validation is missing on resource %s", accountSSLSettingsResourceName)
	}
	return nil
}

func testAccAccountSSLSettingsDestroy(s *terraform.State) error {
	return nil
}

func testCheckAccAccountSSLSettingsAfterFullUpdate() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, ok := state.RootModule().Resources[accountSSLSettingsResource]
		if !ok {
			return fmt.Errorf("incapsula Account SSL settings resource not found : %s", accountSSLSettingsResourceName)
		}

		client := testAccProvider.Meta().(*Client)
		settings, diagnostics := client.GetAccountSSLSettings("")
		if diagnostics != nil && diagnostics.HasError() {
			return fmt.Errorf("failed to get account ssl settings after full update resource %s", accountSSLSettingsResourceName)
		}
		if settings.Errors == nil && settings.Data != nil && !*settings.Data[0].ImpervaCertificate.UseWildCardSanInsteadOfFQDN && *settings.Data[0].ImpervaCertificate.AddNakedDomainSanForWWWSites &&
			settings.Data[0].ImpervaCertificate.Delegation.ValueForCNAMEValidation != "" && len(settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation) == 1 &&
			settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation[0] == "example.com" && *settings.Data[0].ImpervaCertificate.Delegation.AllowCNAMEValidation {
			return nil
		}
		return fmt.Errorf("resource %s was not updated correctly after full update", accountSSLSettingsResourceName)
	}
}

func testCheckAccAccountSSLSettingsAfterPartialUpdate1() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, ok := state.RootModule().Resources[accountSSLSettingsResource]
		if !ok {
			return fmt.Errorf("incapsula Account SSL settings resource not found : %s", accountSSLSettingsResourceName)
		}

		client := testAccProvider.Meta().(*Client)
		settings, diagnostics := client.GetAccountSSLSettings("")
		if diagnostics != nil && diagnostics.HasError() {
			return fmt.Errorf("failed to get account ssl settings after partial update1 resource %s", accountSSLSettingsResourceName)
		}
		if settings.Errors == nil && settings.Data != nil && *settings.Data[0].ImpervaCertificate.UseWildCardSanInsteadOfFQDN && !*settings.Data[0].ImpervaCertificate.AddNakedDomainSanForWWWSites &&
			settings.Data[0].ImpervaCertificate.Delegation.ValueForCNAMEValidation != "" && len(settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation) == 0 &&
			!*settings.Data[0].ImpervaCertificate.Delegation.AllowCNAMEValidation {
			return nil
		}
		return fmt.Errorf("resource %s was not updated correctly after full update", accountSSLSettingsResourceName)
	}
}

func testCheckAccAccountSSLSettingsAfterPartialUpdate2() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, ok := state.RootModule().Resources[accountSSLSettingsResource]
		if !ok {
			return fmt.Errorf("incapsula Account SSL settings resource not found : %s", accountSSLSettingsResourceName)
		}

		client := testAccProvider.Meta().(*Client)
		settings, diagnostics := client.GetAccountSSLSettings("")
		if diagnostics != nil && diagnostics.HasError() {
			return fmt.Errorf("failed to get account ssl settings after partia2 update1 resource %s", accountSSLSettingsResourceName)
		}
		if settings.Errors == nil && settings.Data != nil && *settings.Data[0].ImpervaCertificate.UseWildCardSanInsteadOfFQDN && *settings.Data[0].ImpervaCertificate.AddNakedDomainSanForWWWSites &&
			settings.Data[0].ImpervaCertificate.Delegation.ValueForCNAMEValidation != "" && len(settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation) == 1 &&
			settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation[0] == "example.com" && *settings.Data[0].ImpervaCertificate.Delegation.AllowCNAMEValidation {
			return nil
		}
		return fmt.Errorf("resource %s was not updated correctly after full update", accountSSLSettingsResourceName)
	}
}

func testCheckAccAccountSSLSettingsAfterPartialUpdate3() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		_, ok := state.RootModule().Resources[accountSSLSettingsResource]
		if !ok {
			return fmt.Errorf("incapsula Account SSL settings resource not found : %s", accountSSLSettingsResourceName)
		}

		client := testAccProvider.Meta().(*Client)
		settings, diagnostics := client.GetAccountSSLSettings("")
		if diagnostics != nil && diagnostics.HasError() {
			return fmt.Errorf("failed to get account ssl settings after partia3 update1 resource %s", accountSSLSettingsResourceName)
		}
		if settings.Errors == nil && settings.Data != nil && !*settings.Data[0].ImpervaCertificate.UseWildCardSanInsteadOfFQDN && *settings.Data[0].ImpervaCertificate.AddNakedDomainSanForWWWSites &&
			settings.Data[0].ImpervaCertificate.Delegation.ValueForCNAMEValidation != "" && len(settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation) == 1 &&
			settings.Data[0].ImpervaCertificate.Delegation.AllowedDomainsForCNAMEValidation[0] == "example2.com" && !*settings.Data[0].ImpervaCertificate.Delegation.AllowCNAMEValidation {
			return nil
		}
		return fmt.Errorf("resource %s was not updated correctly after full update", accountSSLSettingsResourceName)
	}
}

func testAccAccountSSLSettingsFullUpdateConfig(t *testing.T) string {
	return fmt.Sprintf(`
	 data "incapsula_account_data" "account_data" {
    }
	resource"%s""%s"{
    account_id = data.incapsula_account_data.account_data.current_account
    use_wild_card_san_instead_of_fqdn = false
    add_naked_domain_san_for_www_sites = true
    allow_cname_validation = true
    allowed_domains_for_cname_validation = ["example.com"]
	}`,
		accountSSLSettingsResourceName, accountSSLSettingsConfigName,
	)
}

func testAccAccountSSLSettingsPartialUpdate1Config(t *testing.T) string {
	return fmt.Sprintf(`
	 data "incapsula_account_data" "account_data" {
    }
	resource"%s""%s"{
    account_id = data.incapsula_account_data.account_data.current_account
    use_wild_card_san_instead_of_fqdn = true
    add_naked_domain_san_for_www_sites = false
	}`,
		accountSSLSettingsResourceName, accountSSLSettingsConfigName,
	)
}

func testAccAccountSSLSettingsPartialUpdate2Config(t *testing.T) string {
	return fmt.Sprintf(`
	 data "incapsula_account_data" "account_data" {
    }
	resource"%s""%s"{
    account_id = data.incapsula_account_data.account_data.current_account
    allow_cname_validation = true
    allowed_domains_for_cname_validation = ["example.com"]
	}`,
		accountSSLSettingsResourceName, accountSSLSettingsConfigName,
	)
}

func testAccAccountSSLSettingsPartialUpdate3Config(t *testing.T) string {
	return fmt.Sprintf(`
    data "incapsula_account_data" "account_data" {
    }
	resource"%s""%s"{
    account_id = data.incapsula_account_data.account_data.current_account
    use_wild_card_san_instead_of_fqdn = false
    allowed_domains_for_cname_validation = ["example2.com"]
	}`,
		accountSSLSettingsResourceName, accountSSLSettingsConfigName,
	)
}

func testACCStateAccountSSLSettingsID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != accountSSLSettingsResourceName {
			continue
		}
		accountID, err := strconv.Atoi(rs.Primary.Attributes["account_id"])
		if err != nil {
			return "", fmt.Errorf("error parsing account ID for import Account Account SSL settings. Value %s", rs.Primary.Attributes["account_id"])
		}

		return fmt.Sprintf("%d", accountID), nil
	}
	return "", fmt.Errorf("error finding an Account SSL settings\"")
}
