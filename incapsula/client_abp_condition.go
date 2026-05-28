package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpConditionResourceName = "ABP Condition"

type AbpCondition struct {
	Id           string `json:"id,omitempty"`
	AccountId    string `json:"account_id,omitempty"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Code         string `json:"code"`
	LastChangeBy string `json:"last_change_by,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	ModifiedAt   string `json:"modified_at,omitempty"`
}

func (c *Client) abpConditionAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/condition", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpConditionUrl(conditionId string) string {
	return fmt.Sprintf("%s/v1/condition/%s", c.config.BaseURLAPI, conditionId)
}

func (c *Client) CreateAbpCondition(accountId string, condition AbpCondition) (*AbpCondition, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpConditionResourceName, accountId)

	body, err := json.Marshal(condition)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpConditionResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpConditionAccountUrl(accountId), body, CreateAbpCondition)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpConditionResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpConditionResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpConditionResourceName, accountId, string(responseBody))
	}

	var created AbpCondition
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpConditionResourceName, err, string(responseBody))
	}
	return &created, nil
}

func (c *Client) ReadAbpCondition(conditionId string) (*AbpCondition, error) {
	log.Printf("[INFO] Reading %s with id %s", abpConditionResourceName, conditionId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpConditionUrl(conditionId), nil, ReadAbpCondition)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpConditionResourceName, conditionId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpConditionResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpConditionResourceName, conditionId, string(responseBody))
	}

	var condition AbpCondition
	if err := json.Unmarshal(responseBody, &condition); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpConditionResourceName, err, string(responseBody))
	}
	return &condition, nil
}

func (c *Client) UpdateAbpCondition(conditionId string, condition AbpCondition) (*AbpCondition, error) {
	log.Printf("[INFO] Updating %s with id %s", abpConditionResourceName, conditionId)

	body, err := json.Marshal(condition)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpConditionResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpConditionUrl(conditionId), body, UpdateAbpCondition)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpConditionResourceName, conditionId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpConditionResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpConditionResourceName, conditionId, string(responseBody))
	}

	var updated AbpCondition
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpConditionResourceName, err, string(responseBody))
	}
	return &updated, nil
}

func (c *Client) DeleteAbpCondition(conditionId string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpConditionResourceName, conditionId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpConditionUrl(conditionId), nil, DeleteAbpCondition)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpConditionResourceName, conditionId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpConditionResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpConditionResourceName, conditionId, string(responseBody))
	}
	return nil
}
