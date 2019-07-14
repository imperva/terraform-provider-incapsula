package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strconv"
	"testing"
)

const certificateName = "Example custom certificate"
const certificateResourceName = "incapsula_custom_certificate.custom-certificate"

func testAccCheckCertificateUpload_goodConfig(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaCustomCertificateGoodConfig(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCertificateExists(certificateResourceName),
					resource.TestCheckResourceAttr(certificateResourceName, "name", certificateName),
				),
			},
			{
				ResourceName:      certificateResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateCertificateID,
			},
		},
	})
}

func testAccCheckCertificateUpload_badKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaCustomCertificateBadKey(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCertificateExists(certificateResourceName),
					resource.TestCheckResourceAttr(certificateResourceName, "name", certificateName),
				),
			},
			{
				ResourceName:      certificateResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateCertificateID,
			},
		},
	})
}

func testAccCheckCertificateUpload_badCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaCustomCertificateBadCertificate(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCertificateExists(certificateResourceName),
					resource.TestCheckResourceAttr(certificateResourceName, "name", certificateName),
				),
			},
			{
				ResourceName:      certificateResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateCertificateID,
			},
		},
	})
}

func testAccCheckCertificateUpload_badPassphrase(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaCustomCertificateBadPassphrase(),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCertificateExists(certificateResourceName),
					resource.TestCheckResourceAttr(certificateResourceName, "name", certificateName),
				),
			},
			{
				ResourceName:      certificateResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccStateCertificateID,
			},
		},
	})
}

func testAccStateCertificateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "incapsula_custom_certificate" {
			continue
		}

		certificateID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return "", fmt.Errorf("Error parsing ID %v to int", rs.Primary.ID)
		}
		siteID, err := strconv.Atoi(rs.Primary.Attributes["site_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing site_id %v to int", rs.Primary.Attributes["site_id"])
		}
		return fmt.Sprintf("%d/%d", siteID, certificateID), nil
	}

	return "", fmt.Errorf("Error finding site_id")
}

func testAccCheckIncapsulaCertificateDestroy(state *terraform.State) error {
	//client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site" {
			continue
		}

		siteID := res.Primary.ID
		if siteID == "" {
			return fmt.Errorf("Incapsula site ID does not exist")
		}

		err := "nil"
		// TODO: Update function to look for cert details on site object from ListCertificates when fix is in place in API
		//listCertificatesResponse, err := client.ListCertificates(siteID)
		//for _, dc := range listCertificatesResponse.Res {
		//	if dc.Name == certificateName {
		//		return fmt.Errorf("Incapsula custom certificate: %s (site_id: %s) still exists", certificateName, siteID)
		//	}
		//}
		if err == "nil" {
			return fmt.Errorf("Incapsula site for domain: %s (site id: %s) still exists", testAccDomain, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaCertificateExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		// TODO: Update function to look for cert details on site object from ListCertificates when fix is in place in API
		return nil
	}
}

// All variations of good and bad certificate configs
func testAccCheckIncapsulaCustomCertificateGoodConfig() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_custom_certificate" "custom-certificate" {
  site_id = "${incapsula_site.example-site.id}"
  certificate = "-----BEGIN CERTIFICATE-----\nMIIDgjCCAmoCCQCk3MsAS5x+UjANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlTYW4gRGllZ28xCzAJBgNVBAoMAlNF\nMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFzaC5iZWVyLmNlbnRlcjEdMBsGCSqG\nSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wHhcNMTkwNzA4MTU0MjQ0WhcNMjAwNzA3\nMTU0MjQ0WjCBgjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlT\nYW4gRGllZ28xCzAJBgNVBAoMAlNFMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFz\naC5iZWVyLmNlbnRlcjEdMBsGCSqGSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCj0rYKUhVtNKQ/oKZdCfxvLKhQ\nLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1DhPRSDGgONdHe\n2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkPQ8dsPpE945lW\nsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607DqdZUS/O4XSx+d0\n5kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1wqRdwGJsUdq6\nkC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2uWCBSY/OLAgMB\nAAEwDQYJKoZIhvcNAQELBQADggEBABfNZcItHdsSpfp8h+1EP5BnRuoKj+l42EI5\nE9dVlqdOZ25+V5Ee899sn2Nj8h+/zVU3+IDO2abUPrDd2xZHaHdf0p69htSwFTHs\nEwUdPUUsKRSys7fVP1clHcKWswTcoWIzQiPZsDMoOQw/pzN05cXSzdo8wSWuEeBK\ncqRNd5BKPeeXbFa4i5TFzT/+pl8V075k16tzHSbT7QDk5fuZWYv/2jImw/lgS/nx\nDWtlprrgG6AX1FzovDs/NnNq/e7vZtn8sdOoO2pCSVymNvctNLV2tFcS8sPQDl5M\nIpnZa3kktAegjsCln1JvD0AFigXrF8wjK+FKGI8SPJfbTQ149+A=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCj0rYKUhVtNKQ/\noKZdCfxvLKhQLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1D\nhPRSDGgONdHe2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkP\nQ8dsPpE945lWsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607Dqd\nZUS/O4XSx+d05kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1\nwqRdwGJsUdq6kC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2u\nWCBSY/OLAgMBAAECggEAfDPprkNzWTta95594vYKO+OYrEknnRUwRV0LF/ziae2P\nLR1EX0AeXKpIXwwwDzpXITJS7mu1c8uZwTIQ5g/f6D4nopULVYwlJZhbjXd49hpx\nhmGfk8227te5BqnVS3IPvRx5vjz+r8obYFZb4JZDGa/v9okAlI04FS0hR/Bl4ckD\naIsztf4R+AO2dP6BxYZGIwcq3jkbf0BdyQpkw4Ds7pdKbSa+PsobseyI2NqR2ryX\n4HH4b89HZj8lfiniIN3tPV6uIvpPS6jJklLKy6zdkIFOng/OGwxXomGkrk9ZjBHm\nJx5yA5YfwPidyt80wO9/26wClXYidfKQC8mDN21owQKBgQDPQbNr/sGiI2QzTOpb\nYTx0FWzWMnn9N2XiQm5rcr9kM5WsXh+anlqP54MeXDGZ2f6L8+aGrghZ/78WbG9J\nDbtEc7qTSRw5LFRglqn32a3ppHToEzOVxsA3g/OBJT5lJJwGMTdeKEXtLMmkm/sz\n1ClFnYJ1I8rNcueI9936odDWKwKBgQDKWgGwWTbqVa3wVIOFvluxolQzo6TEBFbf\nQTJo7byO2iRZvhrZUUk8539Uz2px0Ilzxx61CszhNWDVNwgqsN7FtuzXuCwz9GzU\nyBWkzPKGzvK12aFMYoj/cPbcRfMpYWNoK/YfEKfTRkJJfrJSbWP2XlyEr69te8s7\nB/zxOtUIIQKBgEjoJcOhtF/i70aUkgRfKjLzrnuS+hK3QCHdmJY3oVgQRWCDI77y\nYY0ptZgielhStRZqT/eklM+EBaZPsr4SFIQ56bISD9mU3IG1vkivzFvaPD2/M3BG\noCtnQWt2vII75J7RBVcb9609ChnbvPw4b+RLSi8GzjqDZytpdi7KaXpNAoGAS2Ym\nYvObRs4ONhMHvvojaJk4DtXXO0Lyq9W7VuXe8MvP57CyiG+FfrAz/gIbg7VUwlNb\n2dHgbbpaDpim7mFhYQK8VdVGg0V8l/zGM9Y6OIk8Xw5sz+2XZrdNBN77sFudkt9u\nojyujEcNxBz1jUk9iju29aoREBakr6ZWVfy6DIECgYEAtXxrOsbMsbHhVGqgeGXy\nhLXIltR+7NIUaxpLHhYCMzK9SbyZvx/Hd6m34oTw9ws+tHFpeCyiVU+wQgmx0ARD\ncDLKOPIHTGYhq/H8Oc6/Dzfxs1L/hH34mw5u7hVtAaA+q8iaRGVZ797dTVSxw4U0\nRm+BCDRhDcvaG7qpvFj8T6k=\n-----END PRIVATE KEY-----"
  passphrase = "webco123"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckIncapsulaCustomCertificateBadKey() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_custom_certificate" "custom-certificate" {
  site_id = "${incapsula_site.example-site.id}"
  certificate = "-----BEGIN CERTIFICATE-----\nMIIDgjCCAmoCCQCk3MsAS5x+UjANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlTYW4gRGllZ28xCzAJBgNVBAoMAlNF\nMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFzaC5iZWVyLmNlbnRlcjEdMBsGCSqG\nSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wHhcNMTkwNzA4MTU0MjQ0WhcNMjAwNzA3\nMTU0MjQ0WjCBgjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlT\nYW4gRGllZ28xCzAJBgNVBAoMAlNFMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFz\naC5iZWVyLmNlbnRlcjEdMBsGCSqGSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCj0rYKUhVtNKQ/oKZdCfxvLKhQ\nLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1DhPRSDGgONdHe\n2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkPQ8dsPpE945lW\nsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607DqdZUS/O4XSx+d0\n5kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1wqRdwGJsUdq6\nkC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2uWCBSY/OLAgMB\nAAEwDQYJKoZIhvcNAQELBQADggEBABfNZcItHdsSpfp8h+1EP5BnRuoKj+l42EI5\nE9dVlqdOZ25+V5Ee899sn2Nj8h+/zVU3+IDO2abUPrDd2xZHaHdf0p69htSwFTHs\nEwUdPUUsKRSys7fVP1clHcKWswTcoWIzQiPZsDMoOQw/pzN05cXSzdo8wSWuEeBK\ncqRNd5BKPeeXbFa4i5TFzT/+pl8V075k16tzHSbT7QDk5fuZWYv/2jImw/lgS/nx\nDWtlprrgG6AX1FzovDs/NnNq/e7vZtn8sdOoO2pCSVymNvctNLV2tFcS8sPQDl5M\nIpnZa3kktAegjsCln1JvD0AFigXrF8wjK+FKGI8SPJfbTQ149+A=\n-----END CERTIFICATE-----"
  private_key = "some other bad value"
  passphrase = "webco123"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckIncapsulaCustomCertificateBadCertificate() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_custom_certificate" "custom-certificate" {
  site_id = "${incapsula_site.example-site.id}"
  certificate = "some bad value"
  private_key = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCj0rYKUhVtNKQ/\noKZdCfxvLKhQLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1D\nhPRSDGgONdHe2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkP\nQ8dsPpE945lWsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607Dqd\nZUS/O4XSx+d05kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1\nwqRdwGJsUdq6kC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2u\nWCBSY/OLAgMBAAECggEAfDPprkNzWTta95594vYKO+OYrEknnRUwRV0LF/ziae2P\nLR1EX0AeXKpIXwwwDzpXITJS7mu1c8uZwTIQ5g/f6D4nopULVYwlJZhbjXd49hpx\nhmGfk8227te5BqnVS3IPvRx5vjz+r8obYFZb4JZDGa/v9okAlI04FS0hR/Bl4ckD\naIsztf4R+AO2dP6BxYZGIwcq3jkbf0BdyQpkw4Ds7pdKbSa+PsobseyI2NqR2ryX\n4HH4b89HZj8lfiniIN3tPV6uIvpPS6jJklLKy6zdkIFOng/OGwxXomGkrk9ZjBHm\nJx5yA5YfwPidyt80wO9/26wClXYidfKQC8mDN21owQKBgQDPQbNr/sGiI2QzTOpb\nYTx0FWzWMnn9N2XiQm5rcr9kM5WsXh+anlqP54MeXDGZ2f6L8+aGrghZ/78WbG9J\nDbtEc7qTSRw5LFRglqn32a3ppHToEzOVxsA3g/OBJT5lJJwGMTdeKEXtLMmkm/sz\n1ClFnYJ1I8rNcueI9936odDWKwKBgQDKWgGwWTbqVa3wVIOFvluxolQzo6TEBFbf\nQTJo7byO2iRZvhrZUUk8539Uz2px0Ilzxx61CszhNWDVNwgqsN7FtuzXuCwz9GzU\nyBWkzPKGzvK12aFMYoj/cPbcRfMpYWNoK/YfEKfTRkJJfrJSbWP2XlyEr69te8s7\nB/zxOtUIIQKBgEjoJcOhtF/i70aUkgRfKjLzrnuS+hK3QCHdmJY3oVgQRWCDI77y\nYY0ptZgielhStRZqT/eklM+EBaZPsr4SFIQ56bISD9mU3IG1vkivzFvaPD2/M3BG\noCtnQWt2vII75J7RBVcb9609ChnbvPw4b+RLSi8GzjqDZytpdi7KaXpNAoGAS2Ym\nYvObRs4ONhMHvvojaJk4DtXXO0Lyq9W7VuXe8MvP57CyiG+FfrAz/gIbg7VUwlNb\n2dHgbbpaDpim7mFhYQK8VdVGg0V8l/zGM9Y6OIk8Xw5sz+2XZrdNBN77sFudkt9u\nojyujEcNxBz1jUk9iju29aoREBakr6ZWVfy6DIECgYEAtXxrOsbMsbHhVGqgeGXy\nhLXIltR+7NIUaxpLHhYCMzK9SbyZvx/Hd6m34oTw9ws+tHFpeCyiVU+wQgmx0ARD\ncDLKOPIHTGYhq/H8Oc6/Dzfxs1L/hH34mw5u7hVtAaA+q8iaRGVZ797dTVSxw4U0\nRm+BCDRhDcvaG7qpvFj8T6k=\n-----END PRIVATE KEY-----"
  passphrase = "webco123"
}`, certificateName, siteResourceName,
	)
}

func testAccCheckIncapsulaCustomCertificateBadPassphrase() string {
	return testAccCheckIncapsulaSiteConfig_basic(testAccDomain) + fmt.Sprintf("%s%s%s", `
resource "incapsula_custom_certificate" "custom-certificate" {
  site_id = "${incapsula_site.example-site.id}"
  certificate = "-----BEGIN CERTIFICATE-----\nMIIDgjCCAmoCCQCk3MsAS5x+UjANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlTYW4gRGllZ28xCzAJBgNVBAoMAlNF\nMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFzaC5iZWVyLmNlbnRlcjEdMBsGCSqG\nSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wHhcNMTkwNzA4MTU0MjQ0WhcNMjAwNzA3\nMTU0MjQ0WjCBgjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlT\nYW4gRGllZ28xCzAJBgNVBAoMAlNFMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFz\naC5iZWVyLmNlbnRlcjEdMBsGCSqGSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCj0rYKUhVtNKQ/oKZdCfxvLKhQ\nLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1DhPRSDGgONdHe\n2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkPQ8dsPpE945lW\nsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607DqdZUS/O4XSx+d0\n5kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1wqRdwGJsUdq6\nkC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2uWCBSY/OLAgMB\nAAEwDQYJKoZIhvcNAQELBQADggEBABfNZcItHdsSpfp8h+1EP5BnRuoKj+l42EI5\nE9dVlqdOZ25+V5Ee899sn2Nj8h+/zVU3+IDO2abUPrDd2xZHaHdf0p69htSwFTHs\nEwUdPUUsKRSys7fVP1clHcKWswTcoWIzQiPZsDMoOQw/pzN05cXSzdo8wSWuEeBK\ncqRNd5BKPeeXbFa4i5TFzT/+pl8V075k16tzHSbT7QDk5fuZWYv/2jImw/lgS/nx\nDWtlprrgG6AX1FzovDs/NnNq/e7vZtn8sdOoO2pCSVymNvctNLV2tFcS8sPQDl5M\nIpnZa3kktAegjsCln1JvD0AFigXrF8wjK+FKGI8SPJfbTQ149+A=\n-----END CERTIFICATE-----"
  private_key = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCj0rYKUhVtNKQ/\noKZdCfxvLKhQLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1D\nhPRSDGgONdHe2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkP\nQ8dsPpE945lWsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607Dqd\nZUS/O4XSx+d05kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1\nwqRdwGJsUdq6kC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2u\nWCBSY/OLAgMBAAECggEAfDPprkNzWTta95594vYKO+OYrEknnRUwRV0LF/ziae2P\nLR1EX0AeXKpIXwwwDzpXITJS7mu1c8uZwTIQ5g/f6D4nopULVYwlJZhbjXd49hpx\nhmGfk8227te5BqnVS3IPvRx5vjz+r8obYFZb4JZDGa/v9okAlI04FS0hR/Bl4ckD\naIsztf4R+AO2dP6BxYZGIwcq3jkbf0BdyQpkw4Ds7pdKbSa+PsobseyI2NqR2ryX\n4HH4b89HZj8lfiniIN3tPV6uIvpPS6jJklLKy6zdkIFOng/OGwxXomGkrk9ZjBHm\nJx5yA5YfwPidyt80wO9/26wClXYidfKQC8mDN21owQKBgQDPQbNr/sGiI2QzTOpb\nYTx0FWzWMnn9N2XiQm5rcr9kM5WsXh+anlqP54MeXDGZ2f6L8+aGrghZ/78WbG9J\nDbtEc7qTSRw5LFRglqn32a3ppHToEzOVxsA3g/OBJT5lJJwGMTdeKEXtLMmkm/sz\n1ClFnYJ1I8rNcueI9936odDWKwKBgQDKWgGwWTbqVa3wVIOFvluxolQzo6TEBFbf\nQTJo7byO2iRZvhrZUUk8539Uz2px0Ilzxx61CszhNWDVNwgqsN7FtuzXuCwz9GzU\nyBWkzPKGzvK12aFMYoj/cPbcRfMpYWNoK/YfEKfTRkJJfrJSbWP2XlyEr69te8s7\nB/zxOtUIIQKBgEjoJcOhtF/i70aUkgRfKjLzrnuS+hK3QCHdmJY3oVgQRWCDI77y\nYY0ptZgielhStRZqT/eklM+EBaZPsr4SFIQ56bISD9mU3IG1vkivzFvaPD2/M3BG\noCtnQWt2vII75J7RBVcb9609ChnbvPw4b+RLSi8GzjqDZytpdi7KaXpNAoGAS2Ym\nYvObRs4ONhMHvvojaJk4DtXXO0Lyq9W7VuXe8MvP57CyiG+FfrAz/gIbg7VUwlNb\n2dHgbbpaDpim7mFhYQK8VdVGg0V8l/zGM9Y6OIk8Xw5sz+2XZrdNBN77sFudkt9u\nojyujEcNxBz1jUk9iju29aoREBakr6ZWVfy6DIECgYEAtXxrOsbMsbHhVGqgeGXy\nhLXIltR+7NIUaxpLHhYCMzK9SbyZvx/Hd6m34oTw9ws+tHFpeCyiVU+wQgmx0ARD\ncDLKOPIHTGYhq/H8Oc6/Dzfxs1L/hH34mw5u7hVtAaA+q8iaRGVZ797dTVSxw4U0\nRm+BCDRhDcvaG7qpvFj8T6k=\n-----END PRIVATE KEY-----"
  passphrase = "some bad passphrase"
}`, certificateName, siteResourceName,
	)
}
