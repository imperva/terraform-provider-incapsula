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
	cert := ""
	pkey := ""
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
