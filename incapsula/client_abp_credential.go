package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpCredentialResourceName = "ABP Credential"

type AbpCredential struct {
	Id         string  `json:"id,omitempty"`
	AccountId  string  `json:"account_id,omitempty"`
	Secret     string  `json:"secret,omitempty"`
	CreatedAt  string  `json:"created_at,omitempty"`
	ModifiedAt *string `json:"modified_at,omitempty"`
}

func (c *Client) abpCredentialAccountUrl(accountId string) string {
	return fmt.Sprintf("%s/v1/account/%s/credential", c.config.BaseURLAPI, accountId)
}

func (c *Client) abpCredentialUrl(credentialId string) string {
	return fmt.Sprintf("%s/v1/credential/%s", c.config.BaseURLAPI, credentialId)
}

func (c *Client) CreateAbpCredential(accountId string) (*AbpCredential, error) {
	log.Printf("[INFO] Creating %s for account %s", abpCredentialResourceName, accountId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpCredentialAccountUrl(accountId), nil, CreateAbpCredential)
	if err != nil {
		return nil, fmt.Errorf("error creating %s for account %s: %w", abpCredentialResourceName, accountId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpCredentialResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s for account %s: %s", resp.StatusCode, abpCredentialResourceName, accountId, string(responseBody))
	}

	var created AbpCredential
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpCredentialResourceName, err, string(responseBody))
	}
	return &created, nil
}

func (c *Client) ReadAbpCredential(credentialId string) (*AbpCredential, error) {
	log.Printf("[INFO] Reading %s %s", abpCredentialResourceName, credentialId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpCredentialUrl(credentialId), nil, ReadAbpCredential)
	if err != nil {
		return nil, fmt.Errorf("error reading %s %s: %w", abpCredentialResourceName, credentialId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when reading %s: %w", abpCredentialResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when reading %s %s: %s", resp.StatusCode, abpCredentialResourceName, credentialId, string(responseBody))
	}

	var cred AbpCredential
	if err := json.Unmarshal(responseBody, &cred); err != nil {
		return nil, fmt.Errorf("error parsing %s read response: %w; body: %s", abpCredentialResourceName, err, string(responseBody))
	}
	return &cred, nil
}

func (c *Client) DeleteAbpCredential(credentialId string) error {
	log.Printf("[INFO] Deleting %s %s", abpCredentialResourceName, credentialId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpCredentialUrl(credentialId), nil, DeleteAbpCredential)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpCredentialResourceName, credentialId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpCredentialResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpCredentialResourceName, credentialId, string(responseBody))
	}
	return nil
}
