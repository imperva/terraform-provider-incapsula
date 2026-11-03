package incapsula

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const siteV3ResourceName = "incapsula_site_v3.test-terraform-site-v3"

var siteName string

func GenerateTestSiteName(t *testing.T) string {
	if v := os.Getenv("INCAPSULA_API_ID"); v == "" && t != nil {
		t.Fatal("INCAPSULA_API_ID must be set for acceptance tests")
	}

	s3 := rand.NewSource(time.Now().UnixNano())
	r3 := rand.New(s3)
	siteName = "id" + os.Getenv("INCAPSULA_API_ID") + strconv.Itoa(r3.Intn(1000))
	return siteName
}

func TestIncapsulaSiteV3_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
				),
			},
			{
				ResourceName:      siteV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSiteV3Importer,
			},
		},
	})
}

func TestIncapsulaSiteV3_refId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", "ref_id = \"123456\""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
					resource.TestCheckResourceAttr(siteV3ResourceName, "ref_id", "123456"),
				),
			},
			{
				ResourceName:      siteV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSiteV3Importer,
			},
		},
	})
}

func TestIncapsulaSiteV3_isActive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", "active = false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
					resource.TestCheckResourceAttr(siteV3ResourceName, "active", "false"),
				),
			},
			{
				ResourceName:      siteV3ResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testSiteV3Importer,
			},
		},
	})
}

func TestIncapsulaSiteV3_AccountIdUpdateFails(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testCheckIncapsulaSiteV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testCheckIncapsulaSiteV3ConfigBasic(GenerateTestSiteName(nil), "CLOUD_WAF", ""),
				Check: resource.ComposeTestCheckFunc(
					testCheckIncapsulaSiteExists(siteV3ResourceName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "name", siteName),
					resource.TestCheckResourceAttr(siteV3ResourceName, "type", "CLOUD_WAF"),
				),
			},
			{
				Config:      testCheckIncapsulaSiteV3ConfigWithAccountId(siteName, "CLOUD_WAF", "999999"),
				ExpectError: regexp.MustCompile("account_id cannot be updated for an existing site"),
			},
		},
	})
}

func testCheckIncapsulaSiteV3Destroy(state *terraform.State) error {
	client := testAccProvider.Meta().(*Client)

	for _, res := range state.RootModule().Resources {
		if res.Type != "incapsula_site_v3" {
			continue
		}

		siteIDStr := res.Primary.ID
		if siteIDStr == "" {
			return fmt.Errorf("incapsula site v3 ID does not exist")
		}
		siteID, err := strconv.Atoi(siteIDStr)
		if err != nil {
			return fmt.Errorf("Site ID conversion error for %s: %s", siteIDStr, err)
		}

		siteV3Request := SiteV3Request{}
		siteV3Request.Name = siteName

		_, diags := client.GetV3Site(&siteV3Request, "123")

		if diags == nil {
			return fmt.Errorf("incapsula site for domain: %s (site id: %d) still exists", siteName, siteID)
		}
	}

	return nil
}

func testCheckIncapsulaSiteV3ConfigBasic(name string, siteType string, extraAttr string) string {
	return fmt.Sprintf(`
		resource "incapsula_site_v3" "test-terraform-site-v3" {
			name = "%s"
		    type = "%s"
			%s
		}`,
		name,
		siteType,
		extraAttr,
	)
}

func testCheckIncapsulaSiteV3ConfigWithAccountId(name string, siteType string, accountId string) string {
	return fmt.Sprintf(`
		resource "incapsula_site_v3" "test-terraform-site-v3" {
			name = "%s"
		    type = "%s"
			account_id = "%s"
		}`,
		name,
		siteType,
		accountId,
	)
}

func testSiteV3Importer(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {

		accountId1, err := strconv.Atoi(rs.Primary.Attributes["account_id"])
		if err != nil {
			return "", fmt.Errorf("Error parsing API ID %v to int", rs.Primary.Attributes["id"])
		}

		return fmt.Sprintf("%d/%s", accountId1, rs.Primary.ID), nil
	}
	return "", fmt.Errorf("Error finding an Site V3")
}

func siteV3CtyObject(values map[string]cty.Value) cty.Value {
	attrs := map[string]cty.Value{}
	for name, attrType := range resourceSiteV3().CoreConfigSchema().ImpliedType().AttributeTypes() {
		if v, ok := values[name]; ok {
			attrs[name] = v
		} else {
			attrs[name] = cty.NullVal(attrType)
		}
	}
	return cty.ObjectVal(attrs)
}

func siteV3Diff(t *testing.T, state *terraform.InstanceState, config map[string]interface{}) (*terraform.InstanceDiff, error) {
	rawConfig := map[string]cty.Value{}
	for k, v := range config {
		switch val := v.(type) {
		case string:
			rawConfig[k] = cty.StringVal(val)
		case bool:
			rawConfig[k] = cty.BoolVal(val)
		}
	}
	if state == nil {
		state = &terraform.InstanceState{RawState: cty.NullVal(resourceSiteV3().CoreConfigSchema().ImpliedType())}
	}
	state.RawConfig = siteV3CtyObject(rawConfig)
	return resourceSiteV3().Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), nil)
}

func existingGcpSiteState(isLb *bool) *terraform.InstanceState {
	attributes := map[string]string{
		"id":         "123",
		"account_id": "1",
		"name":       "lb.example.com",
		"type":       "PUBLIC_CLOUD",
		"cloud_type": "GCP",
		"active":     "true",
	}
	rawState := map[string]cty.Value{
		"id":         cty.StringVal("123"),
		"account_id": cty.StringVal("1"),
		"name":       cty.StringVal("lb.example.com"),
		"type":       cty.StringVal("PUBLIC_CLOUD"),
		"cloud_type": cty.StringVal("GCP"),
		"active":     cty.True,
	}
	if isLb != nil {
		attributes["is_load_balancer_site"] = strconv.FormatBool(*isLb)
		rawState["is_load_balancer_site"] = cty.BoolVal(*isLb)
	}
	return &terraform.InstanceState{ID: "123", Attributes: attributes, RawState: siteV3CtyObject(rawState)}
}

func gcpSiteConfig(isLb bool) map[string]interface{} {
	return map[string]interface{}{
		"account_id":            "1",
		"name":                  "lb.example.com",
		"type":                  "PUBLIC_CLOUD",
		"cloud_type":            "GCP",
		"is_load_balancer_site": isLb,
	}
}

func TestSiteV3DiffAdoptsLoadBalancerFlagWhenStateHasNoValue(t *testing.T) {
	diff, err := siteV3Diff(t, existingGcpSiteState(nil), gcpSiteConfig(true))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if diff != nil {
		if _, ok := diff.Attributes["is_load_balancer_site"]; ok {
			t.Errorf("Expected no diff for is_load_balancer_site. Got: %+v", diff.Attributes["is_load_balancer_site"])
		}
		if diff.RequiresNew() {
			t.Errorf("Expected no replacement. Got: %+v", diff)
		}
	}
}

func TestSiteV3DiffPlansLoadBalancerFlagChangeInPlaceOnExistingSite(t *testing.T) {
	current := true
	diff, err := siteV3Diff(t, existingGcpSiteState(&current), gcpSiteConfig(false))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if diff == nil || diff.Attributes["is_load_balancer_site"] == nil || diff.Attributes["is_load_balancer_site"].New != "false" {
		t.Fatalf("Expected in-place diff for is_load_balancer_site. Got: %+v", diff)
	}
	if diff.RequiresNew() {
		t.Errorf("Expected no replacement. Got: %+v", diff)
	}
}

func TestSiteV3UpdateSendsChangedLoadBalancerFlag(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSiteV3().Schema, map[string]interface{}{"is_load_balancer_site": true})
	isLb := changedLoadBalancerSiteFlag(d)
	if isLb == nil || !*isLb {
		t.Fatalf("Expected is_load_balancer_site=true to be sent. Got: %v", isLb)
	}
}

func TestSiteV3UpdateOmitsUnchangedLoadBalancerFlag(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceSiteV3().Schema, map[string]interface{}{"name": "lb.example.com"})
	if isLb := changedLoadBalancerSiteFlag(d); isLb != nil {
		t.Fatalf("Expected is_load_balancer_site to be omitted. Got: %v", *isLb)
	}
}

func TestSiteV3DiffNoChangeWhenLoadBalancerFlagMatchesState(t *testing.T) {
	current := true
	diff, err := siteV3Diff(t, existingGcpSiteState(&current), gcpSiteConfig(true))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if diff != nil && diff.RequiresNew() {
		t.Errorf("Expected no replacement. Got: %+v", diff)
	}
}

func TestSiteV3DiffLeavesLoadBalancerFlagValidationToApi(t *testing.T) {
	config := map[string]interface{}{
		"name":                  "waf.example.com",
		"type":                  "CLOUD_WAF",
		"is_load_balancer_site": false,
	}
	if _, err := siteV3Diff(t, nil, config); err != nil {
		t.Fatalf("Expected no plan-time error. Got: %v", err)
	}
}

func TestSiteV3DiffAllowsLoadBalancerFlagOnNewGcpSite(t *testing.T) {
	config := gcpSiteConfig(true)
	delete(config, "account_id")
	diff, err := siteV3Diff(t, nil, config)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if diff == nil || diff.Attributes["is_load_balancer_site"] == nil || diff.Attributes["is_load_balancer_site"].New != "true" {
		t.Errorf("Expected is_load_balancer_site=true in create diff. Got: %+v", diff)
	}
}

func TestSiteV3DiffOmittedLoadBalancerFlagKeepsStateOnUnrelatedChange(t *testing.T) {
	current := true
	config := gcpSiteConfig(true)
	delete(config, "is_load_balancer_site")
	config["name"] = "renamed.example.com"
	diff, err := siteV3Diff(t, existingGcpSiteState(&current), config)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if diff == nil || diff.Attributes["name"] == nil {
		t.Fatalf("Expected a name diff. Got: %+v", diff)
	}
	if _, ok := diff.Attributes["is_load_balancer_site"]; ok {
		t.Errorf("Expected no diff for an omitted is_load_balancer_site. Got: %+v", diff.Attributes["is_load_balancer_site"])
	}
	if diff.RequiresNew() {
		t.Errorf("Expected no replacement. Got: %+v", diff)
	}
}

func applySiteV3Update(t *testing.T, state *terraform.InstanceState, config map[string]interface{}) string {
	var patchBody string
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPatch {
			body, _ := io.ReadAll(req.Body)
			patchBody = string(body)
		}
		rw.WriteHeader(200)
		rw.Write([]byte(`{"data":[{"id":123,"name":"renamed.example.com","type":"PUBLIC_CLOUD","cloud":"GCP","accountId":1,"active":true}]}`))
	}))
	defer server.Close()
	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	diff, err := siteV3Diff(t, state, config)
	if err != nil {
		t.Fatalf("unexpected diff error: %s", err)
	}
	if diff == nil {
		t.Fatalf("Expected a diff")
	}
	if _, diags := resourceSiteV3().Apply(context.Background(), state, diff, client); diags.HasError() {
		t.Fatalf("unexpected apply error: %v", diags)
	}
	if patchBody == "" {
		t.Fatalf("Expected a PATCH request")
	}
	return patchBody
}

func TestSiteV3ApplySendsChangedLoadBalancerFlagInPatch(t *testing.T) {
	current := true
	config := gcpSiteConfig(false)
	config["name"] = "renamed.example.com"
	body := applySiteV3Update(t, existingGcpSiteState(&current), config)
	if !strings.Contains(body, `"isLoadBalancerSite":false`) {
		t.Errorf("Expected isLoadBalancerSite=false in PATCH body. Got: %s", body)
	}
}

func TestSiteV3ApplyOmitsLoadBalancerFlagWhenStateHasNoValue(t *testing.T) {
	config := gcpSiteConfig(true)
	config["name"] = "renamed.example.com"
	body := applySiteV3Update(t, existingGcpSiteState(nil), config)
	if strings.Contains(body, "isLoadBalancerSite") {
		t.Errorf("Expected no isLoadBalancerSite in PATCH body. Got: %s", body)
	}
}

func TestSiteV3ApplyOmitsUnchangedLoadBalancerFlag(t *testing.T) {
	current := true
	config := gcpSiteConfig(true)
	delete(config, "is_load_balancer_site")
	config["name"] = "renamed.example.com"
	body := applySiteV3Update(t, existingGcpSiteState(&current), config)
	if strings.Contains(body, "isLoadBalancerSite") {
		t.Errorf("Expected no isLoadBalancerSite in PATCH body. Got: %s", body)
	}
}

func TestSiteV3ApplyKeepsPriorStateWhenPatchIsRejected(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(400)
		rw.Write([]byte(`{"errors":[{"status":400,"detail":"isLoadBalancerSite cannot be set"}]}`))
	}))
	defer server.Close()
	client := &Client{config: &Config{APIID: "foo", APIKey: "bar", BaseURLAPI: server.URL}, httpClient: &http.Client{}}

	current := true
	state := existingGcpSiteState(&current)
	config := gcpSiteConfig(false)
	config["name"] = "renamed.example.com"
	diff, err := siteV3Diff(t, state, config)
	if err != nil {
		t.Fatalf("unexpected diff error: %s", err)
	}
	newState, diags := resourceSiteV3().Apply(context.Background(), state, diff, client)
	if !diags.HasError() {
		t.Fatalf("Expected an apply error")
	}
	if newState.Attributes["name"] != "lb.example.com" {
		t.Errorf("Expected name to stay lb.example.com. Got: %s", newState.Attributes["name"])
	}
	if newState.Attributes["is_load_balancer_site"] != "true" {
		t.Errorf("Expected is_load_balancer_site to stay true. Got: %s", newState.Attributes["is_load_balancer_site"])
	}
}
