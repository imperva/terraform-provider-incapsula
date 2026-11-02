package incapsula

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const genericHttpSiemConnectionResourceName = "test_acc_generic_http"
const genericHttpSiemConnectionResource = siemGenericHttpConnectionResourceType + "." + genericHttpSiemConnectionResourceName

var genericHttpSiemConnectionName = "SIEMCONNECTIONGENERICHTTP" + RandomLetterAndNumberString(10)

var genericHttpSiemConnectionNameUpdated = "UPDATED" + genericHttpSiemConnectionName

var genericHttpEnvVars = []string{
	"SIEM_CONNECTION_GENERIC_HTTP_URL",
	"SIEM_CONNECTION_GENERIC_HTTP_HEADER_NAME",
	"SIEM_CONNECTION_GENERIC_HTTP_SECRET",
	"SIEM_CONNECTION_GENERIC_HTTP_USERNAME",
	"SIEM_CONNECTION_GENERIC_HTTP_PASSWORD",
}

func isGenericHttpEnvVarExist() (bool, error) {
	skipTest := false
	for _, envVar := range genericHttpEnvVars {
		if os.Getenv(envVar) == "" {
			skipTest = true
			log.Printf("[ERROR] %s environment variable does not exist, if you want to test SIEM connection you must provide it", envVar)
		}
	}
	return skipTest, nil
}

func TestAccSiemGenericHttpConnection(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connnection_generic_http.go.TestAccSiemGenericHttpConnection")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemGenericHttpConnectionDestroy(siemGenericHttpConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isGenericHttpEnvVarExist,
				Config:   getAccIncapsulaGenericHttpSiemConnectionHeadersConfig(genericHttpSiemConnectionName, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemGenericHttpConnectionExists(genericHttpSiemConnectionResource),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "connection_name", genericHttpSiemConnectionName),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_type", HttpAuthTypeHeaders),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_header.#", "1"),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "disable_cert_verification", "true"),
				),
			},
			{
				SkipFunc:                isGenericHttpEnvVarExist,
				ResourceName:            genericHttpSiemConnectionResource,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_header.0.secret", "password"},
				ImportStateIdFunc:       testACCStateSiemGenericHttpConnectionID(siemGenericHttpConnectionResourceType),
			},
		},
	})
}

func TestAccSiemGenericHttpConnection_TwoHeaders(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connnection_generic_http.go.TestAccSiemGenericHttpConnection_TwoHeaders")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemGenericHttpConnectionDestroy(siemGenericHttpConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isGenericHttpEnvVarExist,
				Config:   getAccIncapsulaGenericHttpSiemConnectionHeadersConfig(genericHttpSiemConnectionName, 2),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemGenericHttpConnectionExists(genericHttpSiemConnectionResource),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_header.#", "2"),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_header.0.header_name", genericHttpAuthHeaderName(0)),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_header.1.header_name", genericHttpAuthHeaderName(1)),
				),
			},
		},
	})
}

func TestAccSiemGenericHttpConnection_BasicAuth(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connnection_generic_http.go.TestAccSiemGenericHttpConnection_BasicAuth")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemGenericHttpConnectionDestroy(siemGenericHttpConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isGenericHttpEnvVarExist,
				Config:   getAccIncapsulaGenericHttpSiemConnectionBasicAuthConfig(genericHttpSiemConnectionName),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemGenericHttpConnectionExists(genericHttpSiemConnectionResource),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_type", HttpAuthTypeBasicAuth),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "username", os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_USERNAME")),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "auth_header.#", "0"),
				),
			},
		},
	})
}

func TestAccSiemGenericHttpConnection_Update(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_connnection_generic_http.go.TestAccSiemGenericHttpConnection_Update")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemGenericHttpConnectionDestroy(siemGenericHttpConnectionResourceType),
		Steps: []resource.TestStep{
			{
				SkipFunc: isGenericHttpEnvVarExist,
				Config:   getAccIncapsulaGenericHttpSiemConnectionHeadersConfig(genericHttpSiemConnectionName, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemGenericHttpConnectionExists(genericHttpSiemConnectionResource),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "connection_name", genericHttpSiemConnectionName),
				),
			},
			{
				SkipFunc: isGenericHttpEnvVarExist,
				Config:   getAccIncapsulaGenericHttpSiemConnectionHeadersConfig(genericHttpSiemConnectionNameUpdated, 1),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemGenericHttpConnectionExists(genericHttpSiemConnectionResource),
					resource.TestCheckResourceAttr(genericHttpSiemConnectionResource, "connection_name", genericHttpSiemConnectionNameUpdated),
				),
			},
		},
	})
}

// Header names must be unique, so every header after the first gets an index suffix.
func genericHttpAuthHeaderName(index int) string {
	headerName := os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_HEADER_NAME")
	if index == 0 {
		return headerName
	}
	return fmt.Sprintf("%s-%d", headerName, index)
}

func getAccIncapsulaGenericHttpSiemConnectionHeadersConfig(connectionName string, authHeaderCount int) string {
	var authHeaders strings.Builder
	for index := 0; index < authHeaderCount; index++ {
		authHeaders.WriteString(fmt.Sprintf(`
			auth_header {
				header_name   = "%s"
				secret_format = "{secret}"
				secret        = "%s"
			}`,
			genericHttpAuthHeaderName(index),
			os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_SECRET"),
		))
	}

	return fmt.Sprintf(`
		resource "%s" "%s" {
			connection_name = "%s"
			url             = "%s"
			auth_type       = "HEADERS"
			disable_cert_verification = true
			headers = {
				"Content-Type" = "application/json"
			}
			%s
		}`,
		siemGenericHttpConnectionResourceType,
		genericHttpSiemConnectionResourceName,
		connectionName,
		os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_URL"),
		authHeaders.String(),
	)
}

func getAccIncapsulaGenericHttpSiemConnectionBasicAuthConfig(connectionName string) string {
	return fmt.Sprintf(`
		resource "%s" "%s" {
			connection_name = "%s"
			url             = "%s"
			auth_type       = "BASIC_AUTH"
			disable_cert_verification = true
			username        = "%s"
			password        = "%s"
		}`,
		siemGenericHttpConnectionResourceType,
		genericHttpSiemConnectionResourceName,
		connectionName,
		os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_URL"),
		os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_USERNAME"),
		os.Getenv("SIEM_CONNECTION_GENERIC_HTTP_PASSWORD"),
	)
}

func testAccIncapsulaSiemGenericHttpConnectionDestroy(resourceType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)
		for _, res := range state.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}
			// Reuses the Splunk read helper: it only checks the connection exists, regardless of storage type.
			if err := testAccReadSiemSplunkConnection(client, res.Primary.ID, res.Primary.Attributes["account_id"]); err == nil {
				return fmt.Errorf("[ERROR] Resource with ID=%s was not destroyed", res.Primary.ID)
			}
		}
		return nil
	}
}

func testCheckIncapsulaSiemGenericHttpConnectionExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection resource not found : %s", resourceName)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemConnection does not exist")
		}
		client := testAccProvider.Meta().(*Client)
		return testAccReadSiemSplunkConnection(client, res.Primary.ID, res.Primary.Attributes["account_id"])
	}
}

func testACCStateSiemGenericHttpConnectionID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			return fmt.Sprintf("%s/%s", rs.Primary.Attributes["account_id"], rs.Primary.ID), nil
		}
		return "", fmt.Errorf("[ERROR] Cannot find SiemConnection ID")
	}
}

func TestSiemGenericHttpConnectionResourceValidation(t *testing.T) {
	authHeader := map[string]interface{}{"header_name": "X-Key", "secret_format": "{secret}", "secret": "value"}
	testCases := []struct {
		name          string
		raw           map[string]interface{}
		expectedError string
	}{
		{
			name: "headers auth with one header is valid",
			raw:  map[string]interface{}{"auth_type": HttpAuthTypeHeaders, "auth_header": []interface{}{authHeader}},
		},
		{
			name:          "headers auth without auth_header fails",
			raw:           map[string]interface{}{"auth_type": HttpAuthTypeHeaders},
			expectedError: "auth_header",
		},
		{
			name:          "headers auth with username fails",
			raw:           map[string]interface{}{"auth_type": HttpAuthTypeHeaders, "auth_header": []interface{}{authHeader}, "username": "admin"},
			expectedError: "mutually exclusive",
		},
		{
			name: "basic auth with username and password is valid",
			raw:  map[string]interface{}{"auth_type": HttpAuthTypeBasicAuth, "username": "admin", "password": "pass"},
		},
		{
			name:          "basic auth with auth_header fails",
			raw:           map[string]interface{}{"auth_type": HttpAuthTypeBasicAuth, "username": "admin", "password": "pass", "auth_header": []interface{}{authHeader}},
			expectedError: "mutually exclusive",
		},
		{
			name:          "basic auth without password fails",
			raw:           map[string]interface{}{"auth_type": HttpAuthTypeBasicAuth, "username": "admin"},
			expectedError: "username and password must be provided",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSiemGenericHttpConnection().Schema, testCase.raw)

			err := siemGenericHttpConnectionResourceValidation(d)

			if testCase.expectedError == "" && err != nil {
				t.Errorf("Should not have received an error, got: %s", err)
			}
			if testCase.expectedError != "" && (err == nil || !strings.Contains(err.Error(), testCase.expectedError)) {
				t.Errorf("Should have received an error containing %q, got: %v", testCase.expectedError, err)
			}
		})
	}
}

func TestFlattenAuthHeadersKeepsSecretsFromStateByIndex(t *testing.T) {
	apiAuthHeaders := []AuthHeaderInfo{
		{HeaderName: "X-First", SecretFormat: "{secret}", Secret: "ciphertext-1"},
		{HeaderName: "X-Second", SecretFormat: "Bearer {secret}", Secret: "ciphertext-2"},
	}
	stateAuthHeaders := []interface{}{map[string]interface{}{"header_name": "X-First", "secret_format": "{secret}", "secret": "plain-1"}}

	flattened := flattenAuthHeadersKeepingSecrets(apiAuthHeaders, stateAuthHeaders)

	first := flattened[0].(map[string]interface{})
	second := flattened[1].(map[string]interface{})
	if first["secret"] != "plain-1" || second["secret"] != "" {
		t.Errorf("Should keep state secrets and never use API ciphertext, got: %v", flattened)
	}
	if second["header_name"] != "X-Second" || second["secret_format"] != "Bearer {secret}" {
		t.Errorf("Should take non-secret fields from the API, got: %v", second)
	}
}

func TestBuildHttpAuthenticationInfoForHeadersKeepsDeclarationOrder(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSiemGenericHttpConnection().Schema, map[string]interface{}{
		"auth_type": HttpAuthTypeHeaders,
		"auth_header": []interface{}{
			map[string]interface{}{"header_name": "X-First", "secret_format": "{secret}", "secret": "one"},
			map[string]interface{}{"header_name": "X-Second", "secret_format": "{secret}", "secret": "two"},
		},
	})

	authentication := buildHttpAuthenticationInfo(d)

	if len(authentication.AuthHeaders) != 2 || authentication.AuthHeaders[0].HeaderName != "X-First" || authentication.AuthHeaders[1].HeaderName != "X-Second" {
		t.Errorf("Should keep auth headers in declaration order, got: %+v", authentication.AuthHeaders)
	}
}
