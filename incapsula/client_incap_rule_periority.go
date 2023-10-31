package incapsula

import (
	"encoding/json"
	"fmt"
	"github.com/terraform-providers/terraform-provider-incapsula/utils"
	"io/ioutil"
	"log"
	"net/http"
)

func (c *Client) ReadIncapRulePriorities(siteID string, ruleType utils.RuleType) (*utils.Response, int, error) {
	log.Printf("[INFO] Getting Incapsula Incap ruleType Rule %s for Site ID %s\n", ruleType.String(), siteID)

	reqURL := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", c.config.BaseURLRev3, siteID, ruleType.String())
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadIncapRule)

	if err != nil {
		return nil, 0, fmt.Errorf("error from Incapsula service when reading Incap Rule Catagories %s for Site ID %s: %s", ruleType.String(), siteID, err)
	}
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, resp.StatusCode, fmt.Errorf("error status code %d from Incapsula service when reading Incap Rule Catagorie %s for Site ID %s: %s", resp.StatusCode, ruleType.String(), siteID, string(responseBody))
	}
	var rulesPriorities utils.Response
	err = json.Unmarshal(responseBody, &rulesPriorities)
	if err != nil {
		return nil, 0, fmt.Errorf("Error parsing Incap Rule Catagorie %s JSON response for Site ID %s: %s\nresponse: %s", ruleType.String(), siteID, err, string(responseBody))
	}
	log.Printf("[INFO] Getting Incapsula Incap ruleType Rule %s for Site ID %s\n - finished", ruleType.String(), siteID)

	return &rulesPriorities, resp.StatusCode, nil
}

func (c *Client) UpdateIncapRulePriorities(siteID string, ruleType utils.RuleType, rule []utils.RuleDetails) (*utils.Response, error) {
	log.Printf("[INFO] Updating Incapsula Incap Rule ruleType %s for Site ID %s\n", ruleType.String(), siteID)
	ruleJSON, err := json.Marshal(rule)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal IncapRule: %s", err)
	}
	log.Printf("[DEBUG] Update rule DTO request: %v\n", string(ruleJSON[:]))

	// Put request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%s/delivery-rules-configuration?category=%s", c.config.BaseURLRev3, siteID, ruleType.String())
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, reqURL, ruleJSON, UpdateIncapRule)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when updating Incap Rule  ruleType %s for Site ID %s: %s", ruleType.String(), siteID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Incap Rule JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating Incap Rule ruleType %s for Site ID %s: %s", resp.StatusCode, ruleType.String(), siteID, string(responseBody))
	}
	var rulesPriorities utils.Response
	err = json.Unmarshal(responseBody, &rulesPriorities)
	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap Rule ruleType %s JSON response for Site ID %s: %s\nresponse: %s", ruleType.String(), siteID, err, string(responseBody))
	}

	return &rulesPriorities, nil
}
