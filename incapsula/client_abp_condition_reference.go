package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpConditionReferenceResourceName = "ABP Condition Reference"

// AbpConditionKind classifies the three variants of ConditionV1.
type AbpConditionKind string

const (
	AbpConditionKindLiteral   AbpConditionKind = "literal"
	AbpConditionKindReference AbpConditionKind = "reference"
	AbpConditionKindList      AbpConditionKind = "list"
)

// AbpConditionReference is the reference variant of ConditionV1: it attaches
// a literal condition or a condition list to a parent condition list.
type AbpConditionReference struct {
	Id         string   `json:"id,omitempty"`
	AccountId  string   `json:"account_id,omitempty"`
	Reference  string   `json:"reference"`
	Parent     string   `json:"parent"`
	Tags       []string `json:"tags"`
	State      string   `json:"state"`
	CreatedAt  string   `json:"created_at,omitempty"`
	ModifiedAt string   `json:"modified_at,omitempty"`
}

func (c *Client) abpConditionReferenceAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/condition", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpConditionReferenceUrl(id string) string {
	return fmt.Sprintf("%s/v1/condition/%s", c.config.BaseURLAPI, id)
}

// ReadAbpConditionKind fetches a condition by id and reports which variant
// (literal / reference / list) it is. The variant is determined by which
// fields are present in the response: a literal always carries `code` (even
// for managed conditions, where it may be the empty string), a reference
// carries `reference`, and a list carries neither.
//
// Returns (nil, nil) when the condition does not exist.
func (c *Client) ReadAbpConditionKind(conditionId string) (*AbpConditionKind, error) {
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpConditionReferenceUrl(conditionId), nil, ReadAbpCondition)
	if err != nil {
		return nil, fmt.Errorf("error reading ABP condition %s: %w", conditionId, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for ABP condition %s: %w", conditionId, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading ABP condition %s: %s", resp.StatusCode, conditionId, string(body))
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, fmt.Errorf("error parsing ABP condition %s: %w; body: %s", conditionId, err, string(body))
	}

	var kind AbpConditionKind
	switch {
	case hasJsonField(fields, "reference"):
		kind = AbpConditionKindReference
	case hasJsonField(fields, "code"):
		kind = AbpConditionKindLiteral
	default:
		kind = AbpConditionKindList
	}
	return &kind, nil
}

func hasJsonField(fields map[string]json.RawMessage, key string) bool {
	_, ok := fields[key]
	return ok
}

func (c *Client) CreateAbpConditionReference(accountId string, ref AbpConditionReference) (*AbpConditionReference, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpConditionReferenceResourceName, accountId)

	if ref.Tags == nil {
		ref.Tags = []string{}
	}

	body, err := json.Marshal(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpConditionReferenceResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpConditionReferenceAccountUrl(accountId), body, CreateAbpConditionReference)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpConditionReferenceResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpConditionReferenceResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpConditionReferenceResourceName, accountId, string(responseBody))
	}

	var created AbpConditionReference
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpConditionReferenceResourceName, err, string(responseBody))
	}
	return &created, nil
}

func (c *Client) ReadAbpConditionReference(id string) (*AbpConditionReference, error) {
	log.Printf("[INFO] Reading %s with id %s", abpConditionReferenceResourceName, id)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpConditionReferenceUrl(id), nil, ReadAbpConditionReference)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpConditionReferenceResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpConditionReferenceResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpConditionReferenceResourceName, id, string(responseBody))
	}

	var ref AbpConditionReference
	if err := json.Unmarshal(responseBody, &ref); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpConditionReferenceResourceName, err, string(responseBody))
	}
	if ref.Reference == "" {
		return nil, fmt.Errorf("ABP Condition %s is not a reference variant", id)
	}
	return &ref, nil
}

func (c *Client) UpdateAbpConditionReference(id string, ref AbpConditionReference) (*AbpConditionReference, error) {
	log.Printf("[INFO] Updating %s with id %s", abpConditionReferenceResourceName, id)

	if ref.Tags == nil {
		ref.Tags = []string{}
	}

	body, err := json.Marshal(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpConditionReferenceResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpConditionReferenceUrl(id), body, UpdateAbpConditionReference)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpConditionReferenceResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpConditionReferenceResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpConditionReferenceResourceName, id, string(responseBody))
	}

	var updated AbpConditionReference
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpConditionReferenceResourceName, err, string(responseBody))
	}
	return &updated, nil
}

func (c *Client) DeleteAbpConditionReference(id string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpConditionReferenceResourceName, id)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpConditionReferenceUrl(id), nil, DeleteAbpConditionReference)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpConditionReferenceResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpConditionReferenceResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpConditionReferenceResourceName, id, string(responseBody))
	}
	return nil
}
