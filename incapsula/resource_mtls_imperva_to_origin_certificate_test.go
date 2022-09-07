package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strconv"

	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	b64 "encoding/base64"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"
)

const mtlsCrtificateResourceName = "incapsula_mtls_imperva_to_origin_certificate"
const mtlsCrtificateResource = mtlsCrtificateResourceName + "." + mtlsCrtificateName
const mtlsCrtificateName = "testacc-terraform-mtls-certificate"

var calculatedHashCACert = ""

func TestAccIncapsulaMtlsImpervaToOriginCertificate_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_mtls_imperva_to_origin_certificate.TestAccIncapsulaMtlsImpervaToOriginCertificate_Basic")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testACCStateMtlsImpervaToOriginCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMtlsImpervaToOriginCertificateBasic(t),
				Check: resource.ComposeTestCheckFunc(
					testCheckMtlsImpervaToOriginCertificateExists(),
					resource.TestCheckResourceAttr(mtlsCrtificateResource, "input_hash", calculatedHashCACert),
					resource.TestCheckResourceAttr(mtlsCrtificateResource, "certificate_name", "acceptance test certificate"),
				),
			},
		},
	})
}

func testACCStateMtlsImpervaToOriginCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != mtlsCrtificateResourceName {
			continue
		}
		return nil

		certificateID := rs.Primary.ID
		if certificateID == "" {
			fmt.Errorf("Parameter id was not found in resource %s", mtlsCrtificateResourceName)
		}

		_, err := client.GetMTLSCertificate(certificateID, "")
		if err == nil {
			return fmt.Errorf("Resource %s with cerificate ID %s still exists", mtlsCrtificateResourceName, certificateID)
		}
	}
	return fmt.Errorf("Error finding site_id")
}

func testCheckMtlsImpervaToOriginCertificateExists() resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[mtlsCrtificateResource]
		if !ok {
			return fmt.Errorf("Incapsula mTLS Imperva to Origin Certificate resource not found : %s", mtlsCrtificateResource)
		}
		certificateID := res.Primary.ID
		client := testAccProvider.Meta().(*Client)
		_, err := client.GetMTLSCertificate(certificateID, "")
		if err != nil {
			return fmt.Errorf("Incapsula mTLS Imperva to Origin Certificate with ID %s does not exist", certificateID)
		}
		return nil
	}
}

func testACCStateMtlsImpervaToOriginCertificateID(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != mtlsCrtificateResourceName {
			continue
		}
		_, err := strconv.Atoi(rs.Primary.Attributes["id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing mTLS Imperva to Origin Certificate ID %v to int", rs.Primary.Attributes["id"])
		}
		return rs.Primary.Attributes["id"], nil
	}
	return "", fmt.Errorf("Error finding an mTLS Imperva to Origin Certificate\"")
}

func testAccCheckMtlsImpervaToOriginCertificateBasic(t *testing.T) string {
	cert, pkey, err := certsetup()
	log.Printf("pkey mtls\n%v", pkey)
	log.Printf("\n\ncertStr mtls\n%v", cert)
	if err != nil {
		panic(err)
	}
	//clientTLSConf := string(clientTLSConf.Certificates[0].Certificate)
	return fmt.Sprintf(`
	resource"%s""%s"{
		certificate_name = "acceptance test certificate"
  		certificate = %s
  		private_key = %s
	}`,
		mtlsCrtificateResourceName, mtlsCrtificateName, cert, pkey,
	)
}

func certsetup() (cert, pkey string, err error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// create our private and public key
	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	// create the CA
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return "", "", err
	}

	// pem encode
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	// set up our server certificate
	certificate := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, certificate, ca, &certPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return "", "", err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	// encode PEM-encoded certificate with base64 algorith
	certificateBase64 := b64.StdEncoding.EncodeToString(certPEM.Bytes())
	// encode PEM-encoded certificate with base64 algorith
	privateKeyBase64 := b64.StdEncoding.EncodeToString(certPrivKeyPEM.Bytes())
	calculatedHashCACert = calculateHash(certificateBase64+"\n", "", privateKeyBase64+"\n")
	return fmt.Sprintf("<<EOT\n%s\nEOT", certificateBase64), fmt.Sprintf("<<EOT\n%s\nEOT", privateKeyBase64), nil
}
