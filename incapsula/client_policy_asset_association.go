package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type PolicyAssetAssociationStatus struct {
	Value   bool `json:"value"`
	IsError bool `json:"isError"`
}

// AddPolicyAssetAssociation adds a policy to be managed by Incapsula
func (c *Client) AddPolicyAssetAssociation(policyID, assetID, assetType string, currentAccountId *int) error {
	log.Printf("[INFO] Adding Incapsula Policy Asset Association: %s/%s/%s\n", policyID, assetID, assetType)

	// Post form to Incapsula
	reqURL := fmt.Sprintf("%s/policies/v2/assets/%s/%s/policies/%s", c.config.BaseURLAPI, assetType, assetID, policyID)
	if currentAccountId != nil && *currentAccountId != 0 {
		reqURL = fmt.Sprintf("%s?caid=%d", reqURL, *currentAccountId)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, CreatePolicyAssetAssociation)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when adding Policy Asset Association: %s", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Add Policy Asset Association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when adding Policy Asset Association: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

// DeletePolicyAssetAssociation deletes a policy asset association currently managed by Incapsula
func (c *Client) DeletePolicyAssetAssociation(policyID, assetID, assetType string, currentAccountId *int) error {
	log.Printf("[INFO] Deleting Incapsula Policy Asset Association: %s/%s/%s\n", policyID, assetID, assetType)

	// Delete request to Incapsula
	reqURL := fmt.Sprintf("%s/policies/v2/assets/%s/%s/policies/%s", c.config.BaseURLAPI, assetType, assetID, policyID)
	if currentAccountId != nil {
		reqURL = fmt.Sprintf("%s?caid=%d", reqURL, *currentAccountId)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil, DeletePolicyAssetAssociation)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when deleting Policy Asset Association (%s): %s", policyID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Delete Policy Asset Association JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting Policy Asset Association: %s", resp.StatusCode, string(responseBody))
	}

	return nil
}

func (c *Client) isPolicyAssetAssociated(policyID, assetID, assetType string, currentAccountId *int) (bool, error) {
	log.Printf("[INFO] Checking Policy Asset Association: %s/%s/%s\n", policyID, assetID, assetType)

	// Check with Policies if the association exist
	reqURL := fmt.Sprintf("%s/policies/v2/policies/%s/assets/%s/%s", c.config.BaseURLAPI, policyID, assetType, assetID)
	if currentAccountId != nil {
		reqURL = fmt.Sprintf("%s?caid=%d", reqURL, *currentAccountId)
	}
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil, ReadPolicyAssetAssociation)
	if err != nil {
		return false, fmt.Errorf("error from Incapsula service when checking if Policy Asset Association exist: %s/%s/%s, err: %s", policyID, assetID, assetType, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	log.Printf("[DEBUG] Incapsula isPolicyAssetAssociated for: %s/%s/%s , response is: %s\n", policyID, assetID, assetType, string(responseBody))

	// Check the response code
	if resp.StatusCode == 404 {
		// If policy asset is not associated 404 will be returned from policies
		return false, nil
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("Error status code %d from Incapsula service when checking the reading Policy Asset Association: %s/%s/%s, response is: %s", resp.StatusCode, policyID, assetID, assetType, string(responseBody))
	}

	// Parse the JSON
	var policyAssetAssociationStatus PolicyAssetAssociationStatus
	err = json.Unmarshal([]byte(responseBody), &policyAssetAssociationStatus)
	if err != nil {
		return false, fmt.Errorf("error parsing Policy Asset Association JSON response for Policy Asset Association: %d/%s/%s: %s\nresponse: %s, err: %s", resp.StatusCode, policyID, assetID, assetType, err, string(responseBody))
	}

	return true, nil
}
