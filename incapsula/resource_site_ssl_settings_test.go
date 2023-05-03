package incapsula

// TODO Uncomment those tests when we have the possibility to approve site ssl certificate
// TODO SSL settings endpoint will not work if the site does not have configured certificate.
//
//const sslSettingsResourceName = "incapsula_site_ssl_settings"
//
//func TestAccSiteSSLSettings_Basic(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck:     func() { testAccPreCheck(t) },
//		Providers:    testAccProviders,
//		CheckDestroy: testAccSiteCheckIncapsulaSiteSSLSettingsDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testSiteCheckIncapsulaIncapSiteSSLSettingsConfigBasic(t),
//				Check: resource.ComposeTestCheckFunc(
//					testCheckIncapsulaIncapSiteSSLSettingsExists(),
//					resource.TestCheckResourceAttr(sslSettingsResourceName, "site_id", "true"),
//					resource.TestCheckResourceAttr(sslSettingsResourceName, "hsts", "true"),
//				),
//			},
//		},
//	})
//}
//
//func testAccSiteCheckIncapsulaSiteSSLSettingsDestroy(state *terraform.State) error {
//	client := testAccProvider.Meta().(*Client)
//
//	for _, res := range state.RootModule().Resources {
//		if res.Type != "incapsula_site_ssl_settings" {
//			continue
//		}
//
//		ruleID, err := strconv.Atoi(res.Primary.ID)
//		if err != nil {
//			return fmt.Errorf("error parsing ID %v to int", res.Primary.ID)
//		}
//
//		siteID, ok := res.Primary.Attributes["site_id"]
//		if !ok {
//			return fmt.Errorf("incapsula Site ID does not exist for SSL settings")
//		}
//
//		var siteIDToInt, _ = strconv.Atoi(siteID)
//
//		_, statusCode, err := client.ReadSiteSSLSettings(siteIDToInt)
//		if statusCode != 404 {
//			return fmt.Errorf("incapsula Incap Site ssl settings %d (site id: %s) should have received 404 status code", ruleID, siteID)
//		}
//		if err == nil {
//			return fmt.Errorf("incapsula Incap Site ssl settings still exists for Site ID %s", siteID)
//		}
//	}
//
//	return nil
//}
//
//func testSiteCheckIncapsulaIncapSiteSSLSettingsConfigBasic(t *testing.T) string {
//	return fmt.Sprintf(`
//		resource "incapsula_site" "testacc-terraform-site" {
//			domain = "%s"
//			force_ssl = "true"
//           domain_validation = "dns"
//		}
//
//		resource "incapsula_site_ssl_settings" "incapsula_site_ssl_settings" {
//		  site_id = "${incapsula_site.testacc-terraform-site.id}"
//		  hsts {
//			is_enabled               = false
//			max_age                  = 31536000
//			sub_domains_included     = false
//			pre_loaded               = false
//		  }
//		}
//`, GenerateTestDomain(t))
//}
//
//func testCheckIncapsulaIncapSiteSSLSettingsExists() resource.TestCheckFunc {
//	return func(state *terraform.State) error {
//		res, ok := state.RootModule().Resources[sslSettingsResourceName]
//		if !ok {
//			return fmt.Errorf("incapsula Site SSL settings resource not found")
//		}
//
//		ruleID, err := strconv.Atoi(res.Primary.ID)
//		if err != nil {
//			return fmt.Errorf("error parsing ID %v to int", res.Primary.ID)
//		}
//
//		siteID, ok := res.Primary.Attributes["site_id"]
//		if !ok || siteID == "" {
//			return fmt.Errorf("incapsula Site ID does not exist for ssl settings %d", ruleID)
//		}
//
//		var siteIDToInt, _ = strconv.Atoi(siteID);
//
//		client := testAccProvider.Meta().(*Client)
//		_, statusCode, err := client.ReadSiteSSLSettings(siteIDToInt)
//		if statusCode != 200 {
//			return fmt.Errorf("incapsula site ssl settings (site id: %s) should have received 200 status code", siteID)
//		}
//		if err != nil {
//			return fmt.Errorf("incapsula site ssl settings (site id: %s) does not exist", siteID)
//		}
//
//		return nil
//	}
//}
