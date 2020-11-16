package incapsula

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// AddPolicyAssetAssociation adds a policy to be managed by Incapsula
func (c *Client) AddPolicyAssetAssociation(policyID, assetID, assetType string) error {
	log.Printf("[INFO] Adding Incapsula Policy Asset Association: %s-%s-%s\n", policyID, assetID, assetType)

	// Post form to Incapsula
	resp, err := c.httpClient.Post(
		fmt.Sprintf("%s/policies/v2/assets/%s/%s/policies/%s?api_id=%s&api_key=%s", c.config.BaseURLAPI, assetType, assetID, policyID, c.config.APIID, c.config.APIKey),
		"application/json",
		nil)
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
func (c *Client) DeletePolicyAssetAssociation(policyID, assetID, assetType string) error {
	log.Printf("[INFO] Deleting Incapsula Policy Asset Association: %s-%s-%s\n", policyID, assetID, assetType)

	// Delete request to Incapsula
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/policies/v2/assets/%s/%s/policies/%s?api_id=%s&api_key=%s", c.config.BaseURLAPI, assetType, assetID, policyID, c.config.APIID, c.config.APIKey),
		nil)
	if err != nil {
		return fmt.Errorf("Error preparing HTTP DELETE for deleting Policy Asset Association: %s", err)
	}
	resp, err := c.httpClient.Do(req)
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
