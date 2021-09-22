package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const apiConfigUrl = "/api-security/api/"

type ApiConfigurationResponse struct {
	Id                           int              `json:"id"`
	SiteId                       int              `json:"siteId"`
	HostName                     string           `json:"hostName"`
	BasePath                     string           `json:"basePath"`
	Description                  string           `json:"description"`
	LastModified                 int              `json:"lastModified"`
	ViolationActions             ViolationActions `json:"violationActions"`
	SpecificationViolationAction string           `json:"specificationViolationAction"`
}
type ApiSecurityApiConfigGetResponse struct {
	Value struct {
		//ApiConfig ApiConfigurationResponse
		Id                           int              `json:"id"`
		SiteId                       int              `json:"siteId"`
		HostName                     string           `json:"hostName"`
		BasePath                     string           `json:"basePath"`
		Description                  string           `json:"description"`
		LastModified                 int              `json:"lastModified"`
		ViolationActions             ViolationActions `json:"violationActions"`
		SpecificationViolationAction string           `json:"specificationViolationAction"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

type ApiSecurityApiConfigGetAllResponse struct {
	Value struct {
		Endpoints []EndpointResponse `json:"endpoints"`
		ApiConfig ApiConfigurationResponse
	} `json:"value"`
	IsError bool `json:"isError"`
}

type ApiSecurityApiConfigPostResponse struct {
	Value struct {
		ApiId int `json:"apiId"`
	} `json:"value"`
	IsError bool `json:"isError"`
}

type ApiSecurityApiConfigGetFileResponse struct {
	Value   string `json:"value"`
	IsError bool   `json:"isError"`
}

type ApiSecurityApiConfigPostPayload struct {
	ValidateHost     bool
	Description      string
	ApiSpecification string
	BasePath         string
	ViolationActions ViolationActions `json:"violation_actions"`
}

type ApiSecurityApiConfigDeleteResponse struct {
	Value   string `json:"value"`
	IsError bool   `json:"isError"`
}

func (c *Client) CreateApiSecurityApiConfig(siteId int, apiConfigPayload *ApiSecurityApiConfigPostPayload) (*ApiSecurityApiConfigPostResponse, error) {
	log.Printf("[INFO] Creating Incapsula API Security API Configuration for Site ID %d\\n", siteId)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//In current implementation validateHost value is always set as "false". Will be changed in next releases
	fw, err := writer.CreateFormField("validateHost")
	if err != nil {
		log.Printf("failed to create %s formdata field", "validateHost")
	}
	_, err = io.Copy(fw, strings.NewReader("false"))
	if err != nil {
		log.Printf("failed to write %s formdata field", "validateHost")
	}

	if apiConfigPayload.Description != "" {
		fw, err := writer.CreateFormField("description")
		if err != nil {
			log.Printf("failed to create %s formdata field", "description")
		}
		_, err = io.Copy(fw, strings.NewReader(apiConfigPayload.Description))
		if err != nil {
			log.Printf("failed to write %s formdata field", "description")
		}
	}
	violationActionsStr, err := json.Marshal(apiConfigPayload.ViolationActions)
	if err != nil {
		fmt.Println(err)
	}

	if apiConfigPayload.ViolationActions != (ViolationActions{}) {
		fw, err := writer.CreateFormField("violationActions")
		if err != nil {
			log.Printf("failed to create %s formdata field", "violationActions")
		}
		_, err = io.Copy(fw, strings.NewReader(string(violationActionsStr)))
		if err != nil {
			log.Printf("failed to write %s formdata field", "violationActions")
		}
	}

	if apiConfigPayload.BasePath != "" {
		fw, err := writer.CreateFormField("basePath")
		if err != nil {
			log.Printf("failed to create %s formdata field", "basePath")
		}
		_, err = io.Copy(fw, strings.NewReader(apiConfigPayload.BasePath))
		if err != nil {
			log.Printf("failed to write %s formdata field", "basePath")
		}
	}

	fw, err = writer.CreateFormFile("apiSpecification", filepath.Base("swagger"))
	if err != nil {
		log.Printf("failed to create %s formdata field", "apiSpecification")
	}
	fw.Write([]byte(apiConfigPayload.ApiSpecification))

	writer.Close()

	reqURL := fmt.Sprintf("%s%s%d", c.config.BaseURLAPI, apiConfigUrl, siteId)
	contentType := writer.FormDataContentType()
	resp, err := c.DoJsonRequestWithHeadersForm(http.MethodPost, reqURL, body.Bytes(), contentType)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error adding API Security API Config for site %d: %s", siteId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Create Api-Security API Config JSON response: %s\n", string(responseBody))

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service while creating API Security API Config for Site ID %d: %v", resp.StatusCode, siteId, string(responseBody))
	}
	// Dump JSON
	var apiAddResponse ApiSecurityApiConfigPostResponse
	err = json.Unmarshal([]byte(responseBody), &apiAddResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing add API Security API Config JSON response for site id %d: %s", siteId, err)
	}

	return &apiAddResponse, nil
}

// UpdateApiSecurityApiConfig updates the Api-Security Api Config
func (c *Client) UpdateApiSecurityApiConfig(siteId int, apiId string, apiConfigPayload *ApiSecurityApiConfigPostPayload) (*ApiSecurityApiConfigPostResponse, error) {
	log.Printf("[INFO] Updating Incapsula API Security API Configuration for Site ID %d, API Config ID %s\n", siteId, apiId)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//In current implementation validateHost value is always set as "false". Will be changed in next releases
	fw, err := writer.CreateFormField("validateHost")
	if err != nil {
		log.Printf("failed to create %s formdata field", "validateHost")
	}
	_, err = io.Copy(fw, strings.NewReader("false"))
	if err != nil {
		log.Printf("failed to write %s formdata field", "validateHost")
	}
	if apiConfigPayload.Description != "" {
		fw, err := writer.CreateFormField("description")
		if err != nil {
			log.Printf("failed to create %s formdata field", "description")
		}
		_, err = io.Copy(fw, strings.NewReader(apiConfigPayload.Description))
		if err != nil {
			log.Printf("failed to write %s formdata field", "description")
		}
	}
	if apiConfigPayload.BasePath != "" {
		fw, err := writer.CreateFormField("basePath")
		if err != nil {
			log.Printf("failed to create %s formdata field", "basePath")
		}
		_, err = io.Copy(fw, strings.NewReader(apiConfigPayload.BasePath))
		if err != nil {
			log.Printf("failed to write %s formdata field", "basePath")
		}
	}
	//init violation actions JSON
	violationActionsStr, err := json.Marshal(apiConfigPayload.ViolationActions)
	if err != nil {
		fmt.Println(err)
	}

	if apiConfigPayload.ViolationActions != (ViolationActions{}) {
		fw, err := writer.CreateFormField("violationActions")
		if err != nil {
			log.Printf("failed to create %s formdata field", "violationActions")
		}
		_, err = io.Copy(fw, strings.NewReader(string(violationActionsStr)))
		if err != nil {
			log.Printf("failed to write %s formdata field", "violationActions")
		}
	}

	fw, err = writer.CreateFormFile("apiSpecification", filepath.Base("swagger"))
	if err != nil {
		log.Printf("failed to create %s formdata field", "apiSpecification")
	}
	fw.Write([]byte(apiConfigPayload.ApiSpecification))

	writer.Close()

	reqURL := fmt.Sprintf("%s%s%d/%s", c.config.BaseURLAPI, apiConfigUrl, siteId, apiId)
	contentType := writer.FormDataContentType()
	resp, err := c.DoJsonRequestWithHeadersForm(http.MethodPost, reqURL, body.Bytes(), contentType)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error updating API Security API Config for site id %d, API id %s :%s", siteId, apiId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Update Api-Security API Config JSON response: %s\n", string(responseBody))

	// Look at the response status code from Incapsula
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error from Incapsula service while updating API Security API for siteId %d, API id %s : %s", siteId, apiId, string(responseBody))
	}
	// Dump JSON
	var apiAddResponse ApiSecurityApiConfigPostResponse
	err = json.Unmarshal([]byte(responseBody), &apiAddResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing update API Security API Config JSON response for Site ID %d, API id %s", siteId, apiId)
	}
	return &apiAddResponse, nil
}

// GetApiSecurityApiConfig gets the Api-Security Api Config
func (c *Client) GetApiSecurityApiConfig(siteId int, apiId int) (*ApiSecurityApiConfigGetResponse, error) {
	log.Printf("[INFO] Getting Incapsula Api-Security API Config for Site ID %d, API Config ID %d\n", siteId, apiId)

	url := fmt.Sprintf("%s%s%d/%d", c.config.BaseURLAPI, apiConfigUrl, siteId, apiId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading Api-Security Api Config for Api ID %d: %s", apiId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Read Api-Security API Config JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading Api-Security Api Config for Api ID %d: %s", resp.StatusCode, apiId, string(responseBody))
	}

	// Parse the JSON
	var apiConfigGetResponse ApiSecurityApiConfigGetResponse
	err = json.Unmarshal([]byte(responseBody), &apiConfigGetResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing GET Api-Security Api Config JSON response for API ID %d: %s\nresponse: %s", apiId, err, string(responseBody))
	}
	return &apiConfigGetResponse, nil
}

// GetApiSecurityApiSwaggerConfig gets the Api-Security  API Config Swagger file content
func (c *Client) GetApiSecurityApiSwaggerConfig(siteId int, apiId int) (*ApiSecurityApiConfigGetFileResponse, error) {
	log.Printf("[INFO] Getting Incapsula Api-Security API Swagger Config for Site ID %d, API Config ID %d\n", siteId, apiId)

	url := fmt.Sprintf("%s%sfile/%d/%d", c.config.BaseURLAPI, apiConfigUrl, siteId, apiId)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading Api-Security Api Config for Api ID %d: %s", apiId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Read Api-Security API Config Swagger JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading Api-Security Api Config for Api ID %d: %s", resp.StatusCode, apiId, string(responseBody))
	}

	// Dump JSON
	var apiSecurityApiConfigGetFileResponse ApiSecurityApiConfigGetFileResponse
	err = json.Unmarshal([]byte(responseBody), &apiSecurityApiConfigGetFileResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing GET Api-Security Api Config JSON response for API ID %d: %s\nresponse: %s", apiId, err, string(responseBody))
	}
	return &apiSecurityApiConfigGetFileResponse, nil
}

// DeleteApiSecurityApiConfig deletes the Api-Security Api + endpoints Config
func (c *Client) DeleteApiSecurityApiConfig(siteID int, apiID string) error {
	log.Printf("[INFO] Deleting Incapsula API Security API for ID %s\n", apiID)

	// Delete request to Incapsula
	reqURL := fmt.Sprintf("%s%s%d/%s", c.config.BaseURLAPI, apiConfigUrl, siteID, apiID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Error from Incapsula service when deleting API Secirity API Config with Site ID %d, API ID %s, : %s", siteID, apiID, err)
	}

	// Read the body
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("[ERROR] Error status code %d from Incapsula service when deleting API Security API Config for Site ID %d, API Config ID %s: %s", resp.StatusCode, siteID, apiID, string(responseBody))
	}
	// Dump JSON
	var apiSecurityApiConfigDeleteResponse ApiSecurityApiConfigDeleteResponse
	err = json.Unmarshal([]byte(responseBody), &apiSecurityApiConfigDeleteResponse)
	if err != nil {
		return fmt.Errorf("[ERROR] Error parsing delete API Secirity API Config JSON response for Site ID %d, API Config ID %s: %s", siteID, apiID, err)
	}

	return nil
}
