package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
)

type PscTargetResponse struct {
	Address string      `json:"data"`
	Errors  []APIErrors `json:"errors"`
}

func getPscTargetUrl(baseURL string, region string, accountID string) string {
	url := fmt.Sprintf("%s/anywhere-provisioner/v1/igc/psc-target?region=%s", baseURL, neturl.QueryEscape(region))
	if accountID != "" {
		url = fmt.Sprintf("%s&caid=%s", url, neturl.QueryEscape(accountID))
	}
	return url
}

func (c *Client) GetPscTarget(region string, accountID string) (*PscTargetResponse, error) {
	log.Printf("[INFO] Getting Incapsula GCP PSC target for region: %s\n", region)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		getPscTargetUrl(c.config.BaseURLAPI, region, accountID),
		nil,
		ReadGcpPscTarget)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while reading GCP PSC target for region %s: %s", region, err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read GCP PSC target response for region %s: %s", region, err)
	}

	log.Printf("[DEBUG] Incapsula get GCP PSC target JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		message := string(responseBody)
		var errorResponse PscTargetResponse
		if json.Unmarshal(responseBody, &errorResponse) == nil && len(errorResponse.Errors) > 0 && errorResponse.Errors[0].Detail != "" {
			message = errorResponse.Errors[0].Detail
		}
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading GCP PSC target for region %s: %s", resp.StatusCode, region, message)
	}

	var response PscTargetResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing GCP PSC target JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error from Incapsula service when reading GCP PSC target for region %s: %s", region, response.Errors[0].Detail)
	}

	return &response, nil
}
