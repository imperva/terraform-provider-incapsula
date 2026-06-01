package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpConditionResourceName = "ABP Condition"

// AbpConditionKind classifies the three variants of ConditionV1.
type AbpConditionKind string

const (
	AbpConditionKindLiteral   AbpConditionKind = "literal"
	AbpConditionKindReference AbpConditionKind = "reference"
	AbpConditionKindList      AbpConditionKind = "list"
)

// AbpCondition is the Go view of ConditionV1, a tagged union with three
// variants:
//
//   - literal:   carries Name, Description, Code
//   - reference: carries Reference, Parent, Tags, State
//   - list:      carries Name, Description (acts as a container for
//     references via Parent on other AbpConditions)
//
// The variant is selected by Kind. Marshaling emits only the fields
// relevant to the active variant to satisfy the API's oneOf schema.
// Unmarshaling derives Kind from which discriminator fields are present
// in the JSON (not from their values), since managed literals may carry
// an empty `code`.
type AbpCondition struct {
	Kind AbpConditionKind

	Id        string
	AccountId string

	Name         string
	Description  string
	Code         string
	LastChangeBy string

	Reference string
	Parent    string
	Tags      []string
	State     string

	CreatedAt  string
	ModifiedAt string
}

// abpConditionWire is the on-the-wire shape of every variant. Pointer
// strings let UnmarshalJSON distinguish "field absent" from "field
// present and empty", which is how the variant is identified.
type abpConditionWire struct {
	Id           string   `json:"id,omitempty"`
	AccountId    string   `json:"account_id,omitempty"`
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Code         *string  `json:"code,omitempty"`
	LastChangeBy string   `json:"last_change_by,omitempty"`
	Reference    *string  `json:"reference,omitempty"`
	Parent       string   `json:"parent,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	State        string   `json:"state,omitempty"`
	CreatedAt    string   `json:"created_at,omitempty"`
	ModifiedAt   string   `json:"modified_at,omitempty"`
}

func (c *AbpCondition) UnmarshalJSON(data []byte) error {
	var w abpConditionWire
	if err := json.Unmarshal(data, &w); err != nil {
		return err
	}

	c.Id = w.Id
	c.AccountId = w.AccountId
	c.LastChangeBy = w.LastChangeBy
	c.Parent = w.Parent
	c.Tags = w.Tags
	c.State = w.State
	c.CreatedAt = w.CreatedAt
	c.ModifiedAt = w.ModifiedAt

	if w.Name != nil {
		c.Name = *w.Name
	}
	if w.Description != nil {
		c.Description = *w.Description
	}
	if w.Code != nil {
		c.Code = *w.Code
	}
	if w.Reference != nil {
		c.Reference = *w.Reference
	}

	switch {
	case w.Reference != nil:
		c.Kind = AbpConditionKindReference
	case w.Code != nil:
		c.Kind = AbpConditionKindLiteral
	default:
		c.Kind = AbpConditionKindList
	}
	return nil
}

// MarshalJSON emits only the fields that belong to the active variant so
// the request body matches one of the CreateConditionV1/UpdateConditionV1
// oneOf branches. The reference variant always emits a (possibly empty)
// tags array as required by UpdateConditionV1.
func (c AbpCondition) MarshalJSON() ([]byte, error) {
	switch c.Kind {
	case AbpConditionKindLiteral:
		return json.Marshal(struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Code        string `json:"code"`
		}{c.Name, c.Description, c.Code})
	case AbpConditionKindReference:
		tags := c.Tags
		if tags == nil {
			tags = []string{}
		}
		return json.Marshal(struct {
			Reference string   `json:"reference"`
			Parent    string   `json:"parent"`
			Tags      []string `json:"tags"`
			State     string   `json:"state"`
		}{c.Reference, c.Parent, tags, c.State})
	case AbpConditionKindList:
		return json.Marshal(struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}{c.Name, c.Description})
	default:
		return nil, fmt.Errorf("unknown ABP Condition kind %q", c.Kind)
	}
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

// ReadAbpCondition fetches a condition by id. Returns (nil, nil) if the
// condition does not exist. Callers should check the returned Kind to
// confirm the variant matches their resource.
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
