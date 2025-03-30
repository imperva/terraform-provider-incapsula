package incapsula

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const abpWebsitesResourceName = "incapsula_abp_websites"
const accountConfigName = "testacc-terraform-abp-websites"
const abpWebsitesResource = abpWebsitesResourceName + "." + accountConfigName

func TestAccAbpWebsites_Basic(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesBasic(t, true, GenerateTestDomain(t)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesMultipleWebsites(t, GenerateTestDomain(t), GenerateTestDomain(t), GenerateTestDomain(t)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites-2"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.1.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.0.enable_mitigation", "false"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites-2"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.1.website.1.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.website.1.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesBasic2(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-2"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "false"),
				),
			},
		},
	})
}

func TestAccAbpWebsites_DuplicateWebsites(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		ErrorCheck: func(err error) error {
			// Normalize newlines as the error will have line breaks in it to limit its width
			msg := strings.ReplaceAll(err.Error(), "\n", " ")
			if strings.Contains(msg, "is referenced twice") {
				return nil
			}
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesDuplicate(t),
			},
		},
	})
}

func TestAccAbpWebsites_DuplicateWebsiteGroups(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		ErrorCheck: func(err error) error {
			// Normalize newlines as the error will have line breaks in it to limit its width
			msg := strings.ReplaceAll(err.Error(), "\n", " ")
			if strings.Contains(msg, "Found duplicate identifier (sites)") {
				return nil
			}
			return err
		},
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsiteGroupsDuplicate(t),
			},
		},
	})
}

func TestAccAbpWebsites_DuplicateNamesButDiscriminatedId(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesDiscriminatedId(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites"),
				),
			},
			{
				Config: testAccAbpWebsitesDiscriminatedId2(t),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name_id", ""),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name_id", "sites-2"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.2.name", "sites"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.2.name_id", "sites-3"),
				),
			},
		},
	})
}

func TestAccAbpWebsites_AutoPublish(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesAutoPublish(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "false"),
					// Can be set due to publishes from other tests, or earlier test runs
					// resource.TestCheckNoResourceAttr(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
					resource.TestCheckNoResourceAttr(abpWebsitesResource, "website_group.1"),
				),
			},
			{
				Config: testAccAbpWebsitesAutoPublish(t, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesAutoPublish(t, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "false"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
		},
	})
}

func TestAccAbpWebsites_ImportBasic(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesBasic(t, true, GenerateTestDomain(t)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites-1"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				ResourceName:            abpWebsitesResourceName + "." + accountConfigName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_publish", "last_publish"},
			},
		},
	})
}

func TestAccAbpWebsites_ImportDuplicateNames(t *testing.T) {
	var websitesResponse AbpTerraformAccount

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "incapsula_account_data" "account_data" {
					}
					resource "%s" "%s" {
						account_id = data.incapsula_account_data.account_data.current_account
						auto_publish = true
						website_group {
							name = "sites"
						}
						website_group {
							name = "sites"
							name_id = "sites-1"
						}
						website_group {
							name = "sites-other"
						}
						website_group {
							name = "sites"
							name_id = "sites-2"
						}
					}`, abpWebsitesResourceName, accountConfigName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.1.name", "sites"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.2.name", "sites-other"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.3.name", "sites"),
				),
			},
			{
				ResourceName:      abpWebsitesResourceName + "." + accountConfigName,
				ImportState:       true,
				ImportStateVerify: true,
				// name_id is checked in `ImportStateCheck`. They differ because we must generate them during import
				ImportStateVerifyIgnore: []string{"auto_publish", "last_publish", "website_group.1.name_id", "website_group.3.name_id"},
				ImportStateCheck: func(is []*terraform.InstanceState) error {
					attrs := is[0].Attributes

					if err := checkEqual(attrs, "website_group.1.name_id", "sites-0"); err != nil {
						return err
					}
					if err := checkEqual(attrs, "website_group.3.name_id", "sites-1"); err != nil {
						return err
					}

					return nil
				},
			},
		},
	})
}

func TestAccAbpWebsites_ChangeSiteId(t *testing.T) {
	var websitesResponse AbpTerraformAccount
	var siteId1 = GenerateTestDomain(t)
	var siteId2 = GenerateTestDomain(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAbpWebsitesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAbpWebsitesChangeSiteId(t, siteId1, siteId2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
			{
				Config: testAccAbpWebsitesChangeSiteId(t, siteId1, siteId2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAbpWebsitesExists(&websitesResponse),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "account_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "auto_publish", "true"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "last_publish"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.name", "sites"),
					resource.TestCheckResourceAttrSet(abpWebsitesResource, "website_group.0.website.0.incapsula_site_id"),
					resource.TestCheckResourceAttr(abpWebsitesResource, "website_group.0.website.0.enable_mitigation", "true"),
				),
			},
		},
	})
}

func checkEqual(attrs map[string]string, l, r string) error {
	if attrs[l] != r {
		return fmt.Errorf("Key `%s` does not match `%s` != `%s` in %+v", l, attrs[l], r, attrs)
	}
	return nil
}

func testAccCheckAbpWebsitesExists(websitesresponse *AbpTerraformAccount) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[abpWebsitesResource]
		if !ok {
			return fmt.Errorf("Not found: %s", abpWebsitesResource)
		}

		accountID, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("Error parsing ID %s to int", rs.Primary.ID)
		}

		accountId := rs.Primary.Attributes["account_id"]

		client := testAccProvider.Meta().(*Client)

		id, err := strconv.Atoi(accountId)
		if err != nil {
			return err
		}

		response, _ := client.ReadAbpWebsites(id)
		if response == nil {
			return fmt.Errorf("Failed to retrieve ABP Websites (id=%d)", accountID)
		}

		*websitesresponse = *response
		return nil
	}
}

func testAccCheckIncapsulaSiteConfig(t *testing.T, name string) string {
	domain := createTestDomain(t, name)
	return fmt.Sprintf(`
		resource "incapsula_site" "%s" {
			domain = "%s"
		}`,
		name,
		domain,
	)
}

func createTestDomain(t *testing.T, id string) string {
	// Using examplewebsite.com like the other tests gives an error
	// "Error from Incapsula service when adding site for domain id1165506sites-1.examplewebsite.com:
	// {"res":1,"res_message":"Unexpected error","debug_info":{"problem":"[Load Balancing mode is disabled for this site, you need to reduce number of servers to 1]","id-info":"999999"}}"
	generatedDomain = strings.ToLower(strings.ReplaceAll(t.Name(), "_", "-")) + "-" + id + ".distil.ninja"
	return generatedDomain
}

func testAccAbpWebsitesMultipleWebsites(t *testing.T, domain1 string, domain2 string, domain3 string) string {
	return testAccCheckIncapsulaSiteConfig(t, domain1) +
		testAccCheckIncapsulaSiteConfig(t, domain2) +
		testAccCheckIncapsulaSiteConfig(t, domain3) +
		fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				incapsula_site_id = incapsula_site.sites-1.id
				enable_mitigation = true
			}
		}
		website_group {
			name = "sites-2"
			website {
				incapsula_site_id = incapsula_site.%s.id
				enable_mitigation = false
			}
			website {
				incapsula_site_id = incapsula_site.%s.id
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, domain1, domain2)
}

func testAccAbpWebsitesBasic(t *testing.T, mitigationEnabled bool, domainName string) string {
	return testAccCheckIncapsulaSiteConfigBasic(domainName) + fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites-1"
			website {
				incapsula_site_id = incapsula_site.testacc-terraform-site.id
				enable_mitigation = %t
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, mitigationEnabled)
}

func testAccAbpWebsitesBasic2(t *testing.T, mitigationEnabled bool) string {
	return testAccCheckIncapsulaSiteConfig(t, "sites-2") + fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites-2"
			website {
				incapsula_site_id = incapsula_site.sites-2.id
				enable_mitigation = %t
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, mitigationEnabled)
}

func testAccAbpWebsitesDuplicate(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfig(t, "sites-2") + fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites-2"
			website {
				incapsula_site_id = incapsula_site.sites-2.id
				enable_mitigation = false
			}
			website {
				incapsula_site_id = incapsula_site.sites-2.id
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccAbpWebsiteGroupsDuplicate(t *testing.T) string {
	return testAccCheckIncapsulaSiteConfig(t, "sites-1") + testAccCheckIncapsulaSiteConfig(t, "sites-2") + fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites"
			website {
				incapsula_site_id = incapsula_site.sites-1.id
				enable_mitigation = false
			}
		}
		website_group {
			name = "sites"
			website {
				incapsula_site_id = incapsula_site.sites-2.id
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccAbpWebsitesDiscriminatedId(t *testing.T) string {
	return fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites"
		}
		website_group {
			name = "sites"
			name_id = "sites-2"
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccAbpWebsitesDiscriminatedId2(t *testing.T) string {
	return fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites"
		}
		website_group {
			name = "sites"
			name_id = "sites-2"
		}
		website_group {
			name = "sites"
			name_id = "sites-3"
		}
	}`, abpWebsitesResourceName, accountConfigName)
}

func testAccAbpWebsitesAutoPublish(t *testing.T, autoPublish bool) string {
	return testAccCheckIncapsulaSiteConfig(t, "sites-1") + fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = %t
		website_group {
			name = "sites-1"
			website {
				incapsula_site_id = incapsula_site.sites-1.id
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, autoPublish)
}

func testAccAbpWebsitesChangeSiteId(t *testing.T, siteId string, siteId2 string) string {
	return testAccCheckIncapsulaSiteConfigBasic(siteId) + testAccCheckIncapsulaSiteConfig(t, siteId2) +
		fmt.Sprintf(`
	data "incapsula_account_data" "account_data" {
    }
	resource "%s" "%s" {
		account_id = data.incapsula_account_data.account_data.current_account
		auto_publish = true
		website_group {
			name = "sites"
			website {
				incapsula_site_id = incapsula_site.%s.id
				enable_mitigation = true
			}
		}
	}`, abpWebsitesResourceName, accountConfigName, siteId)
}

func testAccCheckAbpWebsitesDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != abpWebsitesResourceName {
			continue
		}

		accountID := res.Primary.ID
		if accountID == "" {
			return fmt.Errorf("Incapsula ABP Websites ID does not exist")
		}

		accountId := res.Primary.Attributes["account_id"]

		id, err := strconv.Atoi(accountId)
		if err != nil {
			return err
		}

		websitesResponse, _ := client.ReadAbpWebsites(id)
		if websitesResponse == nil {
			return fmt.Errorf("Failed to check ABP Websites status (id=%s)", accountID)
		}

		if len(websitesResponse.WebsiteGroups) != 0 {
			return fmt.Errorf("Found some website groups remaining after delete: %+v", websitesResponse)
		}
		// if websitesResponse.Errors[0].Status != 404 {
		// 	return fmt.Errorf("Incapsula ABP Websites with id %s still exists", accountID)
		// }
	}

	return nil
}
