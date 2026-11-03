package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	neturl "net/url"
)

type GcpRegionsResponse struct {
	Data   []string    `json:"data"`
	Errors []APIErrors `json:"errors"`
}

func getGcpRegionsUrl(baseURL string, accountID string) string {
	url := fmt.Sprintf("%s/anywhere-provisioner/v1/igc/regions", baseURL)
	if accountID != "" {
		url = fmt.Sprintf("%s?caid=%s", url, neturl.QueryEscape(accountID))
	}
	return url
}

func (c *Client) GetGcpRegions(accountID string) (*GcpRegionsResponse, error) {
	log.Printf("[INFO] Getting Incapsula GCP supported regions\n")

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet,
		getGcpRegionsUrl(c.config.BaseURLAPI, accountID),
		nil,
		ReadGcpRegions)

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service while reading GCP regions: %s", err)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	log.Printf("[DEBUG] Incapsula get GCP regions JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading GCP regions: %s", resp.StatusCode, string(responseBody))
	}

	var response GcpRegionsResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("Error parsing GCP regions JSON response: %s\nresponse: %s", err, string(responseBody))
	}

	if len(response.Errors) > 0 {
		return nil, fmt.Errorf("Error from Incapsula service when reading GCP regions: %s", response.Errors[0].Detail)
	}

	return &response, nil
}
