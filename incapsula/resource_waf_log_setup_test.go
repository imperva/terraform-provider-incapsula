package incapsula

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"math/rand"
	"os"
	"testing"
)

const wafSetupLogResourceType = "incapsula_waf_log_setup"
const wafSetupLogResourceName = "example_waf_log_setup"
const s3BucketName = "bucket_name/log_folder"
const s3AccessKey = "AKIAIOSFODNN7EXAMPLE"
const sftpUserName = "sampleuser"
const sftpHost = "dummyhost"
const sftpDestinationFolder = "/home/user_name/log_folder"
const accountID = 52065879

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestAccIncapsulaWAFLogSetupS3_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_waf_log_setup_test.go.TestAccIncapsulaWAFLogSetupS3_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccSecCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaWAFSetupLogS3ConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "s3_bucket_name", s3BucketName),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "s3_access_key", s3AccessKey),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "s3_secret_key", os.Getenv("S3_SECRET_KEY")),
				),
			},
		},
	})
}

func TestAccIncapsulaWAFLogSetupSFTP_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_waf_log_setup_test.go.TestAccIncapsulaWAFLogSetupSFTP_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccSecCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaWAFSetupLogSFTPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_host", sftpHost),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_user_name", sftpUserName),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_password", os.Getenv("SFTP_PASSWORD")),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_destination_folder", sftpDestinationFolder),
				),
			},
		},
	})
}

func TestAccIncapsulaWAFLogSetupSFTP_Disabled(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_waf_log_setup_test.go.TestAccIncapsulaWAFLogSetupSFTP_Disabled")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccSecCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaWAFSetupLogSFTPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_host", sftpHost),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_user_name", sftpUserName),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_password", os.Getenv("SFTP_PASSWORD")),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "sftp_destination_folder", sftpDestinationFolder),
				),
			},
		},
	})
}

func TestAccIncapsulaWAFLogSetupS3_Disabled(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_waf_log_setup_test.go.TestAccIncapsulaWAFLogSetupS3_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccSecCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaWAFSetupLogS3ConfigDisabled(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "s3_bucket_name", s3BucketName),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "s3_access_key", s3AccessKey),
					resource.TestCheckResourceAttr(wafSetupLogResourceType+"."+wafSetupLogResourceName, "s3_secret_key", os.Getenv("S3_SECRET_KEY")),
				),
			},
		},
	})
}

func testAccSecCheck(t *testing.T) {
	testAccProviderConfigure.Do(func() {
		if v := os.Getenv("S3_SECRET_KEY"); v == "" {
			os.Setenv("S3_SECRET_KEY", fmt.Sprint(rand.Intn(len(letterRunes))))
		}

		if v := os.Getenv("SFTP_PASSWORD"); v == "" {
			os.Setenv("SFTP_PASSWORD", fmt.Sprint(rand.Intn(len(letterRunes))))
		}
	})
	testAccPreCheck(t)
}

func getAccIncapsulaWAFSetupLogS3ConfigBasic() string {
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	return fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "%d"
			s3_bucket_name = "%s"
			s3_access_key = "%s"
			s3_secret_key = "%s"
		}`,
		wafSetupLogResourceType, wafSetupLogResourceName, accountID, s3BucketName, s3AccessKey, s3SecretKey,
	)
}

func getAccIncapsulaWAFSetupLogS3ConfigDisabled() string {
	s3SecretKey := os.Getenv("S3_SECRET_KEY")
	return fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "%d"
			enabled = false
			s3_bucket_name = "%s"
			s3_access_key = "%s"
			s3_secret_key = "%s"
		}`,
		wafSetupLogResourceType, wafSetupLogResourceName, accountID, s3BucketName, s3AccessKey, s3SecretKey,
	)
}

func getAccIncapsulaWAFSetupLogSFTPConfigBasic() string {
	sftpPassword := os.Getenv("SFTP_PASSWORD")
	return fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "%d"
			sftp_host = "%s"
			sftp_user_name = "%s"
			sftp_password = "%s"
			sftp_destination_folder = "%s"
		}`,
		wafSetupLogResourceType, wafSetupLogResourceName, accountID, sftpHost, sftpUserName, sftpPassword, sftpDestinationFolder,
	)
}

func getAccIncapsulaWAFSetupLogSFTPConfigDisabled() string {
	sftpPassword := os.Getenv("SFTP_PASSWORD")
	return fmt.Sprintf(`
		resource "%s" "%s" {
			account_id = "%d"
			enabled = false
			sftp_host = "%s"
			sftp_user_name = "%s"
			sftp_password = "%s"
			sftp_destination_folder = "%s"
		}`,
		wafSetupLogResourceType, wafSetupLogResourceName, accountID, sftpHost, sftpUserName, sftpPassword, sftpDestinationFolder,
	)
}
