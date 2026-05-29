package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpConditionListResourceName = "ABP Condition List"

// AbpConditionList is the list variant of ConditionV1: a named container that
// groups condition references via incapsula_abp_condition_list_entry.
type AbpConditionList struct {
	Id          string `json:"id,omitempty"`
	AccountId   string `json:"account_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at,omitempty"`
	ModifiedAt  string `json:"modified_at,omitempty"`
}

func (c *Client) abpConditionListAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/condition", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpConditionListUrl(id string) string {
	return fmt.Sprintf("%s/v1/condition/%s", c.config.BaseURLAPI, id)
}

func (c *Client) CreateAbpConditionList(accountId string, list AbpConditionList) (*AbpConditionList, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpConditionListResourceName, accountId)

	body, err := json.Marshal(list)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpConditionListResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpConditionListAccountUrl(accountId), body, CreateAbpConditionList)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpConditionListResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpConditionListResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpConditionListResourceName, accountId, string(responseBody))
	}

	var created AbpConditionList
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpConditionListResourceName, err, string(responseBody))
	}
	return &created, nil
}

func (c *Client) ReadAbpConditionList(id string) (*AbpConditionList, error) {
	log.Printf("[INFO] Reading %s with id %s", abpConditionListResourceName, id)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpConditionListUrl(id), nil, ReadAbpConditionList)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpConditionListResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpConditionListResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpConditionListResourceName, id, string(responseBody))
	}

	// ConditionV1 is a tagged union; only the list variant lacks both `code`
	// and `reference`. Reject the other variants so callers don't silently
	// treat a literal/reference as a list.
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(responseBody, &fields); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpConditionListResourceName, err, string(responseBody))
	}
	if _, ok := fields["code"]; ok {
		return nil, fmt.Errorf("ABP Condition %s is not a list variant (it is a literal)", id)
	}
	if _, ok := fields["reference"]; ok {
		return nil, fmt.Errorf("ABP Condition %s is not a list variant (it is a reference)", id)
	}

	var list AbpConditionList
	if err := json.Unmarshal(responseBody, &list); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpConditionListResourceName, err, string(responseBody))
	}
	return &list, nil
}

func (c *Client) UpdateAbpConditionList(id string, list AbpConditionList) (*AbpConditionList, error) {
	log.Printf("[INFO] Updating %s with id %s", abpConditionListResourceName, id)

	body, err := json.Marshal(list)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpConditionListResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpConditionListUrl(id), body, UpdateAbpConditionList)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpConditionListResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpConditionListResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpConditionListResourceName, id, string(responseBody))
	}

	var updated AbpConditionList
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpConditionListResourceName, err, string(responseBody))
	}
	return &updated, nil
}

func (c *Client) DeleteAbpConditionList(id string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpConditionListResourceName, id)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpConditionListUrl(id), nil, DeleteAbpConditionList)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpConditionListResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpConditionListResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpConditionListResourceName, id, string(responseBody))
	}
	return nil
}
