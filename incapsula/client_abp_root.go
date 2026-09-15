package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpRootResourceName = "ABP Account ID"

type AbpRoot struct {
	AccountId string `json:"account_id"`
}

func (c *Client) abpRootUrl() string {
	return fmt.Sprintf("%s/botmanagement/v1/", c.config.BaseURLAPI)
}

// ReadAbpAccountId looks up the ABP account ID belonging to the configured API credentials.
func (c *Client) ReadAbpAccountId() (string, error) {
	log.Printf("[INFO] Reading %s", abpRootResourceName)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpRootUrl(), nil, ReadAbpAccountId)
	if err != nil {
		return "", fmt.Errorf("error reading %s: %w", abpRootResourceName, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body when reading %s: %w", abpRootResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error status code %d from Incapsula service when reading %s: %s", resp.StatusCode, abpRootResourceName, string(responseBody))
	}

	var root AbpRoot
	if err := json.Unmarshal(responseBody, &root); err != nil {
		return "", fmt.Errorf("error parsing %s read response: %w; body: %s", abpRootResourceName, err, string(responseBody))
	}
	return root.AccountId, nil
}
