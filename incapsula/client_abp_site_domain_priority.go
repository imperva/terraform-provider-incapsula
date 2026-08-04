package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpSiteDomainPriorityResourceName = "ABP Site Domain Priority"

type AbpSiteDomainPriority struct {
	DomainIds []string `json:"domain_ids"`
}

func (c *Client) abpSiteDomainPriorityUrl(siteId string) string {
	return fmt.Sprintf("%s/v1/site/%s/domain_priority", c.config.BaseURLAPI, siteId)
}

func (c *Client) ReadAbpSiteDomainPriority(siteId string) (*AbpSiteDomainPriority, error) {
	log.Printf("[INFO] Reading %s for site %s", abpSiteDomainPriorityResourceName, siteId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpSiteDomainPriorityUrl(siteId), nil, ReadAbpSiteDomainPriority)
	if err != nil {
		return nil, fmt.Errorf("error reading %s for site %s: %w", abpSiteDomainPriorityResourceName, siteId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpSiteDomainPriorityResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s for site %s: %s", resp.StatusCode, abpSiteDomainPriorityResourceName, siteId, string(responseBody))
	}

	var dp AbpSiteDomainPriority
	if err := json.Unmarshal(responseBody, &dp); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpSiteDomainPriorityResourceName, err, string(responseBody))
	}
	return &dp, nil
}

func (c *Client) UpdateAbpSiteDomainPriority(siteId string, dp AbpSiteDomainPriority) (*AbpSiteDomainPriority, error) {
	log.Printf("[INFO] Updating %s for site %s", abpSiteDomainPriorityResourceName, siteId)

	body, err := json.Marshal(dp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpSiteDomainPriorityResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpSiteDomainPriorityUrl(siteId), body, UpdateAbpSiteDomainPriority)
	if err != nil {
		return nil, fmt.Errorf("error updating %s for site %s: %w", abpSiteDomainPriorityResourceName, siteId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpSiteDomainPriorityResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s for site %s: %s", resp.StatusCode, abpSiteDomainPriorityResourceName, siteId, string(responseBody))
	}

	var updated AbpSiteDomainPriority
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpSiteDomainPriorityResourceName, err, string(responseBody))
	}
	return &updated, nil
}
