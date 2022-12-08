package incapsula

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"bytes"
	b64 "encoding/base64"
	"fmt"
)

const certificateResourceName = "incapsula_custom_certificate"
const certificateName = "custom-certificate"
const certificateResource = certificateResourceName + "." + certificateName

var calculatedHash = ""
var calculatedHashBase64 = "60f46b532df4f8dc794a2a151c7c6cc3a3b48fc3"

func TestAccIncapsulaCustomCertificate_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIncapsulaCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIncapsulaCustomCertificateGoodConfig(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaCertificateExists(certificateResourceName),
					resource.TestCheckResourceAttr(certificateResource, "input_hash", calculatedHashBase64),
				),
			},
		},
	})
}

func testCheckIncapsulaCertificateExists(name string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[certificateResource]
		if !ok {
			return fmt.Errorf("Incapsula Custom Certificate resource not found : %s", certificateResource)
		}

		siteID, ok := res.Primary.Attributes["site_id"]
		if !ok {
			return fmt.Errorf("Incapsula Custom Certificate Site ID %s does not exist", siteID)
		}

		client := testAccProvider.Meta().(*Client)
		listCertificatesResponse, _ := client.ListCertificates(siteID, ReadHSMCustomCertificate)
		if listCertificatesResponse == nil && listCertificatesResponse.Res == 9413 {
			return fmt.Errorf("Incapsula Custom Certificate : %s (SiteId : %s) does not exist", certificateResource, siteID)
		}
		return nil
	}
}

func testAccCheckIncapsulaCertificateDestroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)
	for _, rs := range state.RootModule().Resources {
		if rs.Type != certificateResourceName {
			continue
		}

		siteID := rs.Primary.Attributes["site_id"]
		if siteID == "" {
			fmt.Errorf("Parameter site_id was not found in resource %s", certificateResourceName)
		}

		// List certificates response object may indicate that the certificate has been deleted (9413)
		listCertificatesResponse, _ := client.ListCertificates(siteID, ReadHSMCustomCertificate)
		if listCertificatesResponse != nil && listCertificatesResponse.Res != 9413 {
			return fmt.Errorf("Resource %s for Incapsula Custom Certificate: site ID %s still exists", certificateResourceName, siteID)
		}
		return nil
	}
	return fmt.Errorf("Error finding site_id in destroy custom certificate test")
}

func testAccCheckIncapsulaCustomCertificateGoodConfig(t *testing.T) string {
	cert := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURyakNDQXBhZ0F3SUJBZ0lFQWVDNjFUQU5CZ2txaGtpRzl3MEJBUVVGQURDQmxqRUxNQWtHQTFVRUJoTUMKYVd3eER6QU5CZ05WQkFnTUJtMWxjbXRoZWpFUU1BNEdBMVVFQnd3SGNtVm9iM1p2ZERFVE1CRUdBMVVFQ2d3SwphVzF3WlhKMllTQmpZVEVYTUJVR0ExVUVDd3dPYVcxd1pYSjJZU0JqWVNCa1pYWXhHREFXQmdOVkJBTU1EMmx0CmNHVnlkbUVnWTJFZ1kyVnlkREVjTUJvR0NTcUdTSWIzRFFFSkFSWU5hVzF3ZGtCcGJYQjJMbU52YlRBZ0Z3MHkKTWpFeE1qa3dOekEzTkROYUdBOHlNRFV3TURReE5UQTNNRGMwTTFvd2RERUxNQWtHQTFVRUJoTUNhV3d4RHpBTgpCZ05WQkFnTUJtbHpjbUZsYkRFUk1BOEdBMVVFQnd3SWRHVnNJR0YyYVhZeEVEQU9CZ05WQkFvTUIwbHRjR1Z5CmRtRXhFakFRQmdOVkJBc01DVU5zYjNWa0lGZEJSakViTUJrR0ExVUVBd3dTWlhoaGJYQnNaWGRsWW5OcGRHVXUKWTI5dE1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBb1FXcVJuVVltSFZSSVAyTAp6M3NpTFN4SFhKYS9kRHh4eFo4RHFaUjA3TnBhM3lYR0duQlVqZFRsWkFleWRtU2w1a0UvNklNcHF1U1gzL243CmM4ZjRsM3VYTXIrYlV4WTNoYVNxK003clVUeHVSclBjS08rWGRXWm5GR0pjTERYbXNQUVpZOHZ1QWkyY3ExRVYKcTJSNlVlODNleXdqZ21RU0RLUEd6d0IrbGFiWnRnc0ZwV200bXJqaTVtVDNzd0xpaFNRV1E4Sko4TUVFUVlsOApESVk5QVNDejBHT3VRdnVHQlFYOUJheDFhK3lvaFVteGs0ZWl0WHdTVXNJT1d2cTk1cHBqVGpRWGNtV2hBeEFVCkVjZk5wRUpSRWV1em9OWmpiRHlYYVRMZ1JOUnhSUzlxRHZuelhSWHg5ZTdWZ0ovQWVpNENuNnJla0ZabjVaRWMKVWtLbFJRSURBUUFCb3lNd0lUQWZCZ05WSFJFRUdEQVdnaFFxTG1WNFlXMXdiR1YzWldKemFYUmxMbU52YlRBTgpCZ2txaGtpRzl3MEJBUVVGQUFPQ0FRRUFiTDNXeFN5TDlkRlF4a3lTVzg0OCtGamMyNTF0M0g0cVlnNzl2cEs1CnQ1WE9tV3BIcmZvUFJpWUNkd1plYnZ1VENVQURLN3VUZXc0eDB5R21VWksyZjNteVNZdjV0Mnc5UGtadTE4U0YKSEtJRS9WUkxPQWZnQ3ZHcVBsbTNkcHVTWmpxVFlvQ0dGRHowZ2gyb3V4Z1pjVUFHTnVvckovOHlMMTRIYmk1VApMcy9IMkZCaDVDVlFBRlNiWEtobWgyRU12bzROcXMvSCtQbVZXdUh1WCtlRGhiZUVHK1ZOdjB2NG5nSzFUTnp4CkJRa2MrL0phRys2RnllV0FBZXFGcDBRanB4ZUh2V0gwaDNMTjlja2pYZ1lnSnNPY0VlWUJvcVN4WWlmWlBtcCsKYkFuUWRPaWxSWXBnNE1zV0wyLzR0YnV3Rmt1cVFKWTZKWjZZK3hOVkoxUUtXdz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0="
	pkey := "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFb3dJQkFBS0NBUUVBb1FXcVJuVVltSFZSSVAyTHozc2lMU3hIWEphL2REeHh4WjhEcVpSMDdOcGEzeVhHCkduQlVqZFRsWkFleWRtU2w1a0UvNklNcHF1U1gzL243YzhmNGwzdVhNcitiVXhZM2hhU3ErTTdyVVR4dVJyUGMKS08rWGRXWm5GR0pjTERYbXNQUVpZOHZ1QWkyY3ExRVZxMlI2VWU4M2V5d2pnbVFTREtQR3p3QitsYWJadGdzRgpwV200bXJqaTVtVDNzd0xpaFNRV1E4Sko4TUVFUVlsOERJWTlBU0N6MEdPdVF2dUdCUVg5QmF4MWEreW9oVW14Cms0ZWl0WHdTVXNJT1d2cTk1cHBqVGpRWGNtV2hBeEFVRWNmTnBFSlJFZXV6b05aamJEeVhhVExnUk5SeFJTOXEKRHZuelhSWHg5ZTdWZ0ovQWVpNENuNnJla0ZabjVaRWNVa0tsUlFJREFRQUJBb0lCQUN5cTZhVko3bHk4anBqYwpQT0I5ZyttTUV1KzRVYnZvMkphOW1jSjlFRUowQUNsT3pUbWdWNVJRcnFHbEVQaU95d2FvcXhYUTdNb2ZSNUkyClFtN0gxa29QV3M0VklQMVhlR2QyV29kU3Z5eDEyeEY1NjJUZnNlQTdXL1RucERJUGNjTThzNTVmZjlMUzNGY1oKMHkwTVhuSkVMZHZaVHJCcEdpaXZkZ01PWEE4ZkhSZUNId09GeTlaRklKdUI2WmJtQ052UFNvR1I2OEF0ZEZLVApMSVc4UmZhN1o5c1I0Wi92YUtzQWZRbEdJUjVqWnZ2cnJWUTJVTWl3NE1BaFh0UENTSGFQbkMwS01ZZmQ2WGRECjVHOUxtR2NhdlQ2YUVFeGFpRnpZeTk1dTRhcngra29BMnY2SEowY1RYSnVXc3orTXNCU0pkV2hyQzZnc1BPOXIKUThFVm1BRUNnWUVBMVE0MWc4dmtVS0lDMHd2d0ttSjNmaUk1eWZ4b01DNDN3ek00WldtMTZTM3VrVXBxSXI3LwpWcFZydjQ4TzR3Mms5dFptenp4bUJERmpPdjY3bHhTb0Z0MWk4V1ZIRmJFUE1OR1FqcllxM1NvaENvRkZFRlJoClMwWFF3R3VpSlpxeWRRdXh1TkZ5cWhha3JpNGVNWjZCaWRKdWVFK25qWmljZWt5Zzd6OXZvWUVDZ1lFQXdYcUEKd2dtckNtbnBnTWI5Q2N3RC9SMlRoaHp1eUIxUTYrUWl6U3l3RHAxaXRScW1aZ2lQb3lhemFzZlpUSytweTBPMwpxaTl0UUNEVDcwQ2VqMDlSdjRhWXJrZ2ptMFVYVmJRbmkyVnB6TWp4RGxWTEg4enozNDB5LzBZdnlwUTZJYTV2Clh5QlNVSkRza09lMGZtVWg0VCtVbDc5NW1WdHNrOUcyTEJ4djNjVUNnWUVBbWFiRnNXYzZJV3kxM0w2ZlZmSHQKZTJuemcxZ2xTNW9KWFIxemJxL3VJVnllME9sNTRkVWRFTFJ5SUpScmlCUXZCRlZiajlsZk9XYmt5WWNzZ3FqRApFTHBZd1A3cFpSdHNlU2lwdUVKb1oxZ2F2QmkrVmlpRWdtUzNTQTVYd2didTdMcWlVVWU4Q2k2S1ZaT3M4dHY5ClVBZ1M0M0dPeE85cTZraVpSL0hYOEFFQ2dZQXVJQTVpTS92YTE3VWJSbFU4Nks1cXdZcFNCc1BHWVhiUlJlb20KRCtsSkVxeGRrS1RxM2srZ0RiSG9Xd3lyQTVYdko0MjV2T1RHelF5NWxTWTM5Q2tCQ0EyT1B4UitCOUt3VStxNQppTXZZVG05cGcxd05rTWJ6SEs1enZUL1hnODc0Q0tYMGY3Z2dET3pZL3VSQTNjdGQ3OUowK3VqNmJwbE1CRXJ4CjZUV2lJUUtCZ0N2OFFIYTNXekxkQm95a3B3ay8rVGZUUmFaK2pzTmhzRU9zUFpDdmRsemxYcVREa05TUHdwZjAKb2o0aVNrdWo4NHVQZStqK2tOdlAzbTB6dGRkRDBLNUtJbmw4aHBDV09CclRaZXVFK2dhajVVeXQ0OHVZN2hOTApIMHg3a3FmaUhyYlh6aHZKTTNsMVBoeEs5cVY3VVc0Zmo2VmkrMnliTUtLbEJPbnRqMTlUCi0tLS0tRU5EIFJTQSBQUklWQVRFIEtFWS0tLS0t"
	//cert, pkey := generateKeyPairBase64()
	certRes := fmt.Sprintf("<<EOT\n%s\nEOT", cert)
	pkeyRes := fmt.Sprintf("<<EOT\n%s\nEOT", pkey)
	//cert, privateKey := generateKeyPair()
	//generateKeyPairBase64
	result := testAccCheckIncapsulaSiteConfigBasic(GenerateTestDomain(t)) + fmt.Sprintf(`
resource "%s" "%s" {
  site_id = incapsula_site.testacc-terraform-site.id
  certificate = %s
  private_key = %s
depends_on = ["%s"]
}`, certificateResourceName, certificateName, certRes, pkeyRes, siteResourceName)
	return result
}

func generateKeyPair() (string, string) {
	template := getCertificateTemplate()

	// generate private key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println(err)
	}

	publickey := &privatekey.PublicKey

	// create a self-signed certificate. template = parent
	var parent = template
	certificate, err := x509.CreateCertificate(rand.Reader, template, parent, publickey, privatekey)
	if err != nil {
		fmt.Println(err)
	}

	// the final version of private key that will be encoded by base64
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privatekey),
		},
	)

	// Encode certificate using PEM algorythm
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	certificateRes := out.String()
	pkeyRes := string(privateKeyPEM)
	calculatedHash = calculateHash(certificateRes+"\n", "", pkeyRes+"\n")

	return fmt.Sprintf("<<EOT\n%s\nEOT", certificateRes), fmt.Sprintf("<<EOT\n%s\nEOT", pkeyRes)
}

func generateKeyPairBase64() (string, string) {
	cert, pkey := generateKeyPair()
	// encode PEM-encoded certificate with base64 algorith
	certificateBase64 := b64.StdEncoding.EncodeToString([]byte(cert))
	// encode PEM-encoded certificate with base64 algorith
	privateKeyBase64 := b64.StdEncoding.EncodeToString([]byte(pkey))

	//save calculated hash for it's verification in step 1 of the test(verify create)
	calculatedHashBase64 = calculateHash(certificateBase64+"\n", "", privateKeyBase64+"\n")
	return certificateBase64, privateKeyBase64
}

func getCertificateTemplate() *x509.Certificate {
	template := &x509.Certificate{
		IsCA:                  true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"Earth"},
			Organization: []string{"Mother Nature"},
		},
		Issuer: pkix.Name{
			CommonName:   "dash.beer.center",
			Country:      []string{"IL"},
			Locality:     []string{"Rehovot"},
			Organization: []string{"MyCompany1"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(0, 0, 1),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	return template
}
