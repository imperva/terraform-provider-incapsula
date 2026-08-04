package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpAccountSitePriorityResourceName = "ABP Account Site Priority"

type AbpAccountSitePriority struct {
	SiteIds []string `json:"site_ids"`
}

func (c *Client) abpAccountSitePriorityUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/site_priority", c.config.BaseURLAPI, accountId)
}

func (c *Client) ReadAbpAccountSitePriority(accountId string) (*AbpAccountSitePriority, error) {
	log.Printf("[INFO] Reading %s for account %s", abpAccountSitePriorityResourceName, accountId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpAccountSitePriorityUrl(accountId), nil, ReadAbpAccountSitePriority)
	if err != nil {
		return nil, fmt.Errorf("error reading %s for account %s: %w", abpAccountSitePriorityResourceName, accountId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpAccountSitePriorityResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s for account %s: %s", resp.StatusCode, abpAccountSitePriorityResourceName, accountId, string(responseBody))
	}

	var sp AbpAccountSitePriority
	if err := json.Unmarshal(responseBody, &sp); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpAccountSitePriorityResourceName, err, string(responseBody))
	}
	return &sp, nil
}

func (c *Client) UpdateAbpAccountSitePriority(accountId string, sp AbpAccountSitePriority) (*AbpAccountSitePriority, error) {
	log.Printf("[INFO] Updating %s for account %s", abpAccountSitePriorityResourceName, accountId)

	body, err := json.Marshal(sp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpAccountSitePriorityResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpAccountSitePriorityUrl(accountId), body, UpdateAbpAccountSitePriority)
	if err != nil {
		return nil, fmt.Errorf("error updating %s for account %s: %w", abpAccountSitePriorityResourceName, accountId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpAccountSitePriorityResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s for account %s: %s", resp.StatusCode, abpAccountSitePriorityResourceName, accountId, string(responseBody))
	}

	var updated AbpAccountSitePriority
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpAccountSitePriorityResourceName, err, string(responseBody))
	}
	return &updated, nil
}
