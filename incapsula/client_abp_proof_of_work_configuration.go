package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpProofOfWorkConfigurationResourceName = "ABP Proof Of Work Configuration"

type AbpProofOfWorkConfiguration struct {
	Id         string `json:"id,omitempty"`
	AccountId  string `json:"account_id,omitempty"`
	Name       string `json:"name"`
	Difficulty int64  `json:"difficulty"`
	Algorithm  string `json:"algorithm,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	ModifiedAt string `json:"modified_at,omitempty"`
}

func (c *Client) abpProofOfWorkConfigurationAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/proof_of_work_configuration", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpProofOfWorkConfigurationUrl(id string) string {
	return fmt.Sprintf("%s/v1/proof_of_work_configuration/%s", c.config.BaseURLAPI, id)
}

func (c *Client) CreateAbpProofOfWorkConfiguration(accountId string, config AbpProofOfWorkConfiguration) (*AbpProofOfWorkConfiguration, error) {
	log.Printf("[INFO] Creating %s in ABP account %s", abpProofOfWorkConfigurationResourceName, accountId)

	body, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpProofOfWorkConfigurationResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpProofOfWorkConfigurationAccountUrl(accountId), body, CreateAbpProofOfWorkConfiguration)
	if err != nil {
		return nil, fmt.Errorf("error creating %s in ABP account %s: %w", abpProofOfWorkConfigurationResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpProofOfWorkConfigurationResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s in ABP account %s: %s", resp.StatusCode, abpProofOfWorkConfigurationResourceName, accountId, string(responseBody))
	}

	var created AbpProofOfWorkConfiguration
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpProofOfWorkConfigurationResourceName, err, string(responseBody))
	}
	return &created, nil
}

func (c *Client) ReadAbpProofOfWorkConfiguration(id string) (*AbpProofOfWorkConfiguration, error) {
	log.Printf("[INFO] Reading %s with id %s", abpProofOfWorkConfigurationResourceName, id)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpProofOfWorkConfigurationUrl(id), nil, ReadAbpProofOfWorkConfiguration)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpProofOfWorkConfigurationResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpProofOfWorkConfigurationResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpProofOfWorkConfigurationResourceName, id, string(responseBody))
	}

	var config AbpProofOfWorkConfiguration
	if err := json.Unmarshal(responseBody, &config); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpProofOfWorkConfigurationResourceName, err, string(responseBody))
	}
	return &config, nil
}

func (c *Client) UpdateAbpProofOfWorkConfiguration(id string, config AbpProofOfWorkConfiguration) (*AbpProofOfWorkConfiguration, error) {
	log.Printf("[INFO] Updating %s with id %s", abpProofOfWorkConfigurationResourceName, id)

	body, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpProofOfWorkConfigurationResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPut, c.abpProofOfWorkConfigurationUrl(id), body, UpdateAbpProofOfWorkConfiguration)
	if err != nil {
		return nil, fmt.Errorf("error updating %s %s: %w", abpProofOfWorkConfigurationResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when updating %s: %w", abpProofOfWorkConfigurationResourceName, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d from Incapsula service when updating %s %s: %s", resp.StatusCode, abpProofOfWorkConfigurationResourceName, id, string(responseBody))
	}

	var updated AbpProofOfWorkConfiguration
	if err := json.Unmarshal(responseBody, &updated); err != nil {
		return nil, fmt.Errorf("error parsing %s update response: %w; body: %s", abpProofOfWorkConfigurationResourceName, err, string(responseBody))
	}
	return &updated, nil
}

func (c *Client) DeleteAbpProofOfWorkConfiguration(id string) error {
	log.Printf("[INFO] Deleting %s with id %s", abpProofOfWorkConfigurationResourceName, id)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpProofOfWorkConfigurationUrl(id), nil, DeleteAbpProofOfWorkConfiguration)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpProofOfWorkConfigurationResourceName, id, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpProofOfWorkConfigurationResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpProofOfWorkConfigurationResourceName, id, string(responseBody))
	}
	return nil
}
