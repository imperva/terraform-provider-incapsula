package incapsula

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const abpDomainEncryptionKeyResourceName = "ABP Domain Encryption Key"

type AbpDomainEncryptionKey struct {
	Id               string  `json:"id,omitempty"`
	DomainId         string  `json:"domain_id,omitempty"`
	Key              string  `json:"key"`
	CreatedAt        string  `json:"created_at,omitempty"`
	FirstPublishedAt *string `json:"first_published_at,omitempty"`
}

type abpDomainEncryptionKeyList struct {
	Items []AbpDomainEncryptionKey `json:"items"`
}

func (c *Client) abpDomainEncryptionKeyDomainUrl(domainId string) string {
	return fmt.Sprintf("%s/v1/domain/%s/encryptionkey", c.config.BaseURLAPI, domainId)
}

func (c *Client) abpDomainEncryptionKeyUrl(keyId string) string {
	return fmt.Sprintf("%s/v1/encryptionkey/%s", c.config.BaseURLAPI, keyId)
}

func (c *Client) CreateAbpDomainEncryptionKey(domainId string, key AbpDomainEncryptionKey) (*AbpDomainEncryptionKey, error) {
	log.Printf("[INFO] Creating %s for domain %s", abpDomainEncryptionKeyResourceName, domainId)

	body, err := json.Marshal(map[string]string{"key": key.Key})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %s: %w", abpDomainEncryptionKeyResourceName, err)
	}

	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, c.abpDomainEncryptionKeyDomainUrl(domainId), body, CreateAbpDomainEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error creating %s for domain %s: %w", abpDomainEncryptionKeyResourceName, domainId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when creating %s: %w", abpDomainEncryptionKeyResourceName, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error status code %d from Incapsula service when creating %s for domain %s: %s", resp.StatusCode, abpDomainEncryptionKeyResourceName, domainId, string(responseBody))
	}

	var created AbpDomainEncryptionKey
	if err := json.Unmarshal(responseBody, &created); err != nil {
		return nil, fmt.Errorf("error parsing %s create response: %w; body: %s", abpDomainEncryptionKeyResourceName, err, string(responseBody))
	}
	return &created, nil
}

// ReadAbpDomainEncryptionKey looks up a key by listing all keys on the domain
// and matching by ID. The API does not expose a GET-by-id endpoint. Returns
// (nil, nil) if the domain or the key is gone.
func (c *Client) ReadAbpDomainEncryptionKey(domainId, keyId string) (*AbpDomainEncryptionKey, error) {
	log.Printf("[INFO] Reading %s %s on domain %s", abpDomainEncryptionKeyResourceName, keyId, domainId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, c.abpDomainEncryptionKeyDomainUrl(domainId), nil, ListAbpDomainEncryptionKeys)
	if err != nil {
		return nil, fmt.Errorf("error listing %ss on domain %s: %w", abpDomainEncryptionKeyResourceName, domainId, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body when listing %ss: %w", abpDomainEncryptionKeyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error status code %d when listing %ss on domain %s: %s", resp.StatusCode, abpDomainEncryptionKeyResourceName, domainId, string(responseBody))
	}

	var list abpDomainEncryptionKeyList
	if err := json.Unmarshal(responseBody, &list); err != nil {
		return nil, fmt.Errorf("error parsing %s list response: %w; body: %s", abpDomainEncryptionKeyResourceName, err, string(responseBody))
	}
	for i := range list.Items {
		if list.Items[i].Id == keyId {
			return &list.Items[i], nil
		}
	}
	return nil, nil
}

func (c *Client) DeleteAbpDomainEncryptionKey(keyId string) error {
	log.Printf("[INFO] Deleting %s %s", abpDomainEncryptionKeyResourceName, keyId)

	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, c.abpDomainEncryptionKeyUrl(keyId), nil, DeleteAbpDomainEncryptionKey)
	if err != nil {
		return fmt.Errorf("error deleting %s %s: %w", abpDomainEncryptionKeyResourceName, keyId, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body when deleting %s: %w", abpDomainEncryptionKeyResourceName, err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("error status code %d from Incapsula service when deleting %s %s: %s", resp.StatusCode, abpDomainEncryptionKeyResourceName, keyId, string(responseBody))
	}
	return nil
}
