package incapsula

import (
	"regexp"
	"testing"

	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const hsmCertificateResourceName = "incapsula_custom_hsm_certificate"
const hsmCertificateName = "hsm-custom-certificate"

//Since the data to create HSM certificate is sensitive (Fortanix api key & id) we can only do negative testing
func TestAccImpervaCustomHsmCertificateWithWrongFortanixApiKey_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckImpervaCustomHsmCertificateWrongFortanixApiKeyConfig(t),
				ExpectError: regexp.MustCompile("Get ESK session for site 0 got error response from api.amer.smartkey.io API server connection_status: CC_OK http_status: 401 message: Unauthorized access"),
			},
		},
	})
}

func TestAccImpervaCustomHsmCertificateWithWrongCertificate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckImpervaCustomHsmCertificateWrongCertificate(t),
				ExpectError: regexp.MustCompile("invalid certificate"),
			},
		},
	})
}

func testAccCheckImpervaCustomHsmCertificateWrongFortanixApiKeyConfig(t *testing.T) string {
	certificate := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURuRENDQW9TZ0F3SUJBZ0lKQUk3SEpML3hwZXp5TUEwR0NTcUdTSWIzRFFFQkN3VUFNR2t4Q3pBSkJnTlYKQkFZVEFrbE1NUTh3RFFZRFZRUUlEQVpKYzNKaFpXd3hFREFPQmdOVkJBY01CM0psYUc5MmIzUXhFakFRQmdOVgpCQW9NQ1UxNVEyOXRjR0Z1ZVRFTU1Bb0dBMVVFQ3d3RFpHVjJNUlV3RXdZRFZRUUREQXhwYm1OaGNIUmxjM1F1ClkyOHdIaGNOTWpFd09EQTFNRGt5TmpJeldoY05Nak13T0RBMU1Ea3lOakl6V2pCcE1Rc3dDUVlEVlFRR0V3SkoKVERFUE1BMEdBMVVFQ0F3R1NYTnlZV1ZzTVJBd0RnWURWUVFIREFkeVpXaHZkbTkwTVJJd0VBWURWUVFLREFsTgplVU52YlhCaGJua3hEREFLQmdOVkJBc01BMlJsZGpFVk1CTUdBMVVFQXd3TWFXNWpZWEIwWlhOMExtTnZNSUlCCklqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FROEFNSUlCQ2dLQ0FRRUEyVCtJdEJ6aXdYSzllb3RhTzl6YXZVVTkKY2I5eDc4OGxpaWYyR1VtaUd0V0x6KzNURzVOb2lGRFpYY2tPMEFuUE1GTGdLMXRENERBbFdZcUdlcm1BNldMYQpCZDRXSy80RnhhWWFhM1AxWC9UcEVaMVBvSlVzeGRnYk9tSnUvQmVDdEpTbXZ0WlZrTC9jOEpKWHZ4WVBSRkxKCmt2NGU5bE5pcURQWDIwT2dUQlUxajN2amxhSUo0dEsyMjJyZ0pKQlhjOTdGUjUrY1pSejRDaS9LOFFuV0dHOW8KL08xbHdZUldrc2plMG9QZDRCR1hzNFg2eUNRQXRGMGI2Q1F1NDZ0OW5hemNieFlkaEZuY1g2YStQd0VyeDlNKwpVWlQvTlVNSTRlM01IYWtMU25uZi9TTkw0dDc0SnNpenVJZll4bGwzWExPWnVRQkhuYUppMVFzZHFOZ2dsd0lECkFRQUJvMGN3UlRBTEJnTlZIUThFQkFNQ0JEQXdFd1lEVlIwbEJBd3dDZ1lJS3dZQkJRVUhBd0V3SVFZRFZSMFIKQkJvd0dJSVdLaTV6YzJ4MFpXRnRMbWx1WTJGd2RHVnpkQzVqYnpBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQQpqYitMMkd0Yyt2UUR6b25OWlJpYmFEbEhpdVlUQVJ2VHgrTVNiaFMvMkN1NEhOdHBodnZ5S0dYMnl2aU54Q3RlClJEMEpKYzl4ZWs0U05KMllKdWcvR3RTa01remx0OENxdFJvSG9VWnU4VGFPTnpwMVg2YzRzVXhaN01DeG80RksKQ0hOeFU2MDArcENOZXFVMVNBUnZSUUVVNEpNb1RxditROElzc2dvUVB2dUhXdEFTUGV3eDdiR3FUdkxJWnR2UwpPRXduN20rZlJ1emFTbTcxRzVKTDFyek5aeit1TmJjWjZBTmEzZDNNc1BaSXIzWkp1bUNXQTIwUW5IL0xjUDIyCm9aUFpiYWVuODVUTVBOYU55VzErZHBmZHFnTXdoMS94UFZGeXFNMzg1UE5WeCtpekVVeGpuRnRFb1ZWMGtBV0wKY2dEUHlkSzYxeURLaWhlVDRrUEw5QT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0="
	result := testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
  		site_id = incapsula_site.testacc-terraform-site.id
  		certificate = "%s"
		api_detail {
			api_id = "fakeKey"
			api_key = "fakeId"
			hostname = "api.amer.smartkey.io"
		}
	}`, hsmCertificateResourceName, hsmCertificateName, certificate)
	return result
}

func testAccCheckImpervaCustomHsmCertificateWrongCertificate(t *testing.T) string {
	certificate := "fakeCertificate"
	result := testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
	resource "%s" "%s" {
  		site_id = incapsula_site.testacc-terraform-site.id
  		certificate = "%s"
		api_detail {
			api_id = "fakeKey"
			api_key = "fakeId"
			hostname = "api.amer.smartkey.io"
		}
	}`, hsmCertificateResourceName, hsmCertificateName, certificate)
	return result
}
