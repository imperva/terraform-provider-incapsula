package incapsula

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const siemLogConfigurationResourceType = "incapsula_siem_log_configuration"
const siemLogConfigurationResourceName = "test_acc"
const siemLogConfigurationResource = siemLogConfigurationResourceType + "." + siemLogConfigurationResourceName

var siemLogConfigurationName = "SIEMLOGCONFIGURATION" + RandomLetterAndNumberString(10)

var siemLogConfigurationNameUpdated = "UPDATE" + siemLogConfigurationName

func TestSiemLogConfiguration_Basic(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_log_configuration.go.TestAccSiemLogConfiguration_Basic")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\"", "\"GOOGLE_ANALYTICS_IDS\", \"SIGNIFICANT_DOMAIN_DISCOVERY\", \"SIGNIFICANT_SCRIPT_DISCOVERY\", \"SIGNIFICANT_DATA_TRANSFER_DISCOVERY\", \"DOMAIN_DISCOVERY_ENFORCE_MODE\", \"CSP_HEADER_HEALTH\"", "\"CLOUD_WAF_ACCESS\", \"WAF_RAW_LOGS\"", "\"WAF_ANALYTICS_LOGS\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_csp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_cloudwaf"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationName+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "producer", "ABP"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationName+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "producer", "NETSEC"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationName+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "producer", "ATO"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationName+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "producer", "AUDIT"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "configuration_name", siemLogConfigurationName+"csp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "producer", "CSP"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "configuration_name", siemLogConfigurationName+"cloudwaf"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "producer", "CLOUD_WAF"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "format", "CEF"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "configuration_name", siemLogConfigurationName+"attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "producer", "ATTACK_ANALYTICS"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "format", "CEF"),
				),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_abp",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_netsec",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_ato",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_audit",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_csp",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_cloudwaf",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
			{
				ResourceName:      siemLogConfigurationResource + "_attackanalytics",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testACCStateSiemLogConfigurationID(siemLogConfigurationResourceType),
			},
		},
	})
}

func TestSiemLogConfiguration_Update(t *testing.T) {
	log.Printf("========================BEGIN TEST========================")
	log.Printf("[DEBUG]Running test resource_siem_log_configuration.go.TestAccSiemLogConfiguration_Update")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccIncapsulaSiemLogConfigurationDestroy(siemLogConfigurationResourceType),
		Steps: []resource.TestStep{
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\"", "\"GOOGLE_ANALYTICS_IDS\", \"SIGNIFICANT_DOMAIN_DISCOVERY\", \"SIGNIFICANT_SCRIPT_DISCOVERY\", \"SIGNIFICANT_DATA_TRANSFER_DISCOVERY\", \"DOMAIN_DISCOVERY_ENFORCE_MODE\", \"CSP_HEADER_HEALTH\"", "\"CLOUD_WAF_ACCESS\", \"WAF_RAW_LOGS\"", "\"WAF_ANALYTICS_LOGS\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_csp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_cloudwaf"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationName+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationName+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationName+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationName+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "configuration_name", siemLogConfigurationName+"csp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "configuration_name", siemLogConfigurationName+"cloudwaf"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "configuration_name", siemLogConfigurationName+"attackanalytics"),
				),
			},
			{
				Config: getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationNameUpdated, "\"ABP\"", "\"CONNECTION\", \"NETFLOW\"", "\"ATO\"", "\"AUDIT_TRAIL\"", "\"GOOGLE_ANALYTICS_IDS\", \"SIGNIFICANT_DOMAIN_DISCOVERY\", \"SIGNIFICANT_SCRIPT_DISCOVERY\", \"SIGNIFICANT_DATA_TRANSFER_DISCOVERY\", \"DOMAIN_DISCOVERY_ENFORCE_MODE\", \"CSP_HEADER_HEALTH\"", "\"CLOUD_WAF_ACCESS\", \"WAF_RAW_LOGS\"", "\"WAF_ANALYTICS_LOGS\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_abp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_netsec"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_ato"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_audit"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_csp"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_cloudwaf"),
					testCheckIncapsulaSiemLogConfigurationExists(siemLogConfigurationResource+"_attackanalytics"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_abp", "configuration_name", siemLogConfigurationNameUpdated+"abp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_netsec", "configuration_name", siemLogConfigurationNameUpdated+"netsec"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_ato", "configuration_name", siemLogConfigurationNameUpdated+"ato"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_audit", "configuration_name", siemLogConfigurationNameUpdated+"audit"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_csp", "configuration_name", siemLogConfigurationNameUpdated+"csp"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_cloudwaf", "configuration_name", siemLogConfigurationNameUpdated+"cloudwaf"),
					resource.TestCheckResourceAttr(siemLogConfigurationResource+"_attackanalytics", "configuration_name", siemLogConfigurationNameUpdated+"attackanalytics"),
				),
			},
		},
	})
}

func getAccIncapsulaSiemLogConfigurationConfigBasic(siemLogConfigurationName string, abpDatasets string, netsecDatasets string, atoDatasets string, auditDatasets string, cspDatasets string, cloudWafDatasets string, attackAnalyticsDatasets string) string {
	return getAccIncapsulaS3ArnSiemConnectionConfigBasic(s3ArnSiemConnectionName, s3ArnSiemConnectionTestPath) + fmt.Sprintf(`
		resource "%s" "%s" {
			configuration_name = "%s"
  			producer = "ABP"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_abp", siemLogConfigurationName+"abp",
		abpDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {	
			configuration_name = "%s"
			producer = "NETSEC"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_netsec", siemLogConfigurationName+"netsec",
		netsecDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {	
			configuration_name = "%s"
			producer = "ATO"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_ato", siemLogConfigurationName+"ato",
		atoDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {	
			configuration_name = "%s"
			producer = "AUDIT"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_audit", siemLogConfigurationName+"audit",
		auditDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
		resource "%s" "%s" {
			configuration_name = "%s"
			producer = "CSP"
			datasets = [%s]
			enabled = true
			connection_id = %s.%s.id
		}`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_csp", siemLogConfigurationName+"csp",
		cspDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
	    resource "%s" "%s" {
	       configuration_name = "%s"
	       producer = "CLOUD_WAF"
	       datasets = [%s]
	       enabled = true
	       connection_id = %s.%s.id
           format = "CEF"
        }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_cloudwaf", siemLogConfigurationName+"cloudwaf",
		cloudWafDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	) + fmt.Sprintf(`
      	 resource "%s" "%s" {
      	    configuration_name = "%s"
      	    producer = "ATTACK_ANALYTICS"
      	    datasets = [%s]
      	    enabled = true
      	    connection_id = %s.%s.id
			format = "CEF"
         }`,
		siemLogConfigurationResourceType, siemLogConfigurationResourceName+"_attackanalytics", siemLogConfigurationName+"attackanalytics",
		attackAnalyticsDatasets, siemConnectionResourceType, s3ArnSiemConnectionResourceName,
	)
}

func testAccReadSiemLogConfiguration(client *Client, ID string, accountId string) error {
	log.Printf("[INFO] SiemLogConfiguration ID: %s", ID)
	siemLogConfiguration, statusCode, err := client.ReadSiemLogConfiguration(ID, accountId)
	if err != nil {
		return err
	}

	if (*statusCode == 200) && (siemLogConfiguration != nil) && (len(siemLogConfiguration.Data) == 1) && (siemLogConfiguration.Data[0].ID == ID) {
		log.Printf("[INFO] SiemLogConfiguration : %v\n", siemLogConfiguration)
		return nil
	} else if *statusCode == 400 {
		return fmt.Errorf("[ERROR] SiemLogConfiguration with id %s not found", ID)
	} else {
		return fmt.Errorf("[ERROR] Unsupported operation. Response status code: %d", *statusCode)
	}
}

func testAccIncapsulaSiemLogConfigurationDestroy(resourceType string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client := testAccProvider.Meta().(*Client)

		for _, res := range state.RootModule().Resources {
			if res.Type != resourceType {
				continue
			}

			err := testAccReadSiemLogConfiguration(client, res.Primary.ID, res.Primary.Attributes["account_id"])
			if err != nil {
				return nil
			} else {
				return fmt.Errorf("[ERROR] Resource with ID=%s was not destroyed", res.Primary.ID)
			}

		}

		return nil
	}
}

func testCheckIncapsulaSiemLogConfigurationExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		res, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("[ERROR] Incapsula SiemLogConfiguration resource not found : %s", resource)
		}

		if res.Primary.ID == "" {
			return fmt.Errorf("[ERROR] Incapsula SiemLogConfiguration does not exist")
		} else {
			client := testAccProvider.Meta().(*Client)
			return testAccReadSiemLogConfiguration(client, res.Primary.ID, res.Primary.Attributes["account_id"])
		}
	}
}

func testACCStateSiemLogConfigurationID(resourceType string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		for _, rs := range state.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			return fmt.Sprintf("%s/%s", rs.Primary.Attributes["account_id"], rs.Primary.ID), nil
		}

		return "", fmt.Errorf("[ERROR] Cannot find SiemLogConfiguration ID")
	}
}

func TestBuildSiemLogConfigurationPayloadOmitsPayloadWhenBatchPackagingTypeUnset(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSiemLogConfiguration().Schema, map[string]interface{}{})

	payload := buildSiemLogConfigurationPayload(d)

	if payload != nil {
		t.Errorf("Should not build a payload when batch_packaging_type is unset, got: %+v", payload)
	}
}

func TestBuildSiemLogConfigurationPayloadSetsBatchPackagingType(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSiemLogConfiguration().Schema, map[string]interface{}{"batch_packaging_type": "IMPERVA_METADATA_PREFIX"})

	payload := buildSiemLogConfigurationPayload(d)

	if payload == nil || payload.BatchPackaging.Type != "IMPERVA_METADATA_PREFIX" {
		t.Errorf("Should send batch_packaging_type as payload.batchPackaging.type, got: %+v", payload)
	}
}

func TestSiemLogConfigurationDataSerializesPayloadOnlyWhenSet(t *testing.T) {
	withoutPayload, _ := json.Marshal(SiemLogConfigurationData{})
	withPayload, _ := json.Marshal(SiemLogConfigurationData{Payload: &Payload{BatchPackaging: BatchPackaging{Type: "NONE"}}})

	if strings.Contains(string(withoutPayload), "payload") {
		t.Errorf("Should omit payload when unset, got: %s", withoutPayload)
	}
	if !strings.Contains(string(withPayload), `"payload":{"batchPackaging":{"type":"NONE"}}`) {
		t.Errorf("Should serialize payload.batchPackaging.type, got: %s", withPayload)
	}
}

func TestBatchPackagingTypeFromPayloadReadsBackType(t *testing.T) {
	if batchPackagingTypeFromPayload(nil) != "" {
		t.Errorf("Should map a missing payload to an empty batch_packaging_type")
	}
	if batchPackagingTypeFromPayload(&Payload{BatchPackaging: BatchPackaging{Type: "NONE"}}) != "NONE" {
		t.Errorf("Should map payload.batchPackaging.type to batch_packaging_type")
	}
}

func TestSiemLogConfigurationSchemaAcceptsAnyProducerAndDataset(t *testing.T) {
	testCases := []struct {
		name     string
		producer string
		datasets []interface{}
	}{
		{name: "dataset of another producer", producer: "ATTACK_ANALYTICS", datasets: []interface{}{"CLOUD_WAF_ACCESS"}},
		{name: "dataset unknown to the provider", producer: "CLOUD_WAF", datasets: []interface{}{"NEW_DATASET"}},
		{name: "producer unknown to the provider", producer: "NEW_PRODUCER", datasets: []interface{}{"NEW_DATASET"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"configuration_name": "name",
				"producer":           testCase.producer,
				"datasets":           testCase.datasets,
				"enabled":            true,
				"connection_id":      "connection",
			})

			diagnostics := resourceSiemLogConfiguration().Validate(config)

			if diagnostics.HasError() {
				t.Errorf("Should leave producer and dataset validation to the API, got: %v", diagnostics)
			}
		})
	}
}

func TestSiemLogConfigurationCreateSurfacesApiDatasetValidationError(t *testing.T) {
	detail := "Dataset CLOUD_WAF_ACCESS is not valid for provider ATTACK_ANALYTICS."
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(`{"errors":[{"status":400,"id":"abc","code":"1050","title":"Bad Request","detail":"` + detail + `"}]}`))
	}))
	defer server.Close()
	client := &Client{config: &Config{BaseURLAPI: server.URL}, httpClient: &http.Client{}}
	d := schema.TestResourceDataRaw(t, resourceSiemLogConfiguration().Schema, map[string]interface{}{
		"configuration_name": "name",
		"producer":           "ATTACK_ANALYTICS",
		"datasets":           []interface{}{"CLOUD_WAF_ACCESS"},
		"enabled":            true,
		"connection_id":      "connection",
	})

	err := resourceSiemLogConfigurationCreate(d, client)

	if err == nil || !strings.Contains(err.Error(), detail) {
		t.Errorf("Should return the API detail, got: %v", err)
	}
}
