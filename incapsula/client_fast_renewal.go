package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointFastRenewalConfiguration = "/certificates-ui/v3/sites/%s/certificate/fastRenewal"

type FastRenewalConfigurationDto struct {
	Errors []ApiErrorResponse         `json:"errors"`
	Data   []FastRenewalConfiguration `json:"data"`
}

type FastRenewalConfiguration struct {
	FastRenewal bool `json:"fastRenewal"`
}

func (c *Client) GetFastRenewalConfiguration(siteId string, caid string) (*FastRenewalConfigurationDto, error) {
	if siteId == "" {
		fmt.Errorf("[ERROR] site ID was not provided")
	}
	endpoint := fmt.Sprintf(endpointFastRenewalConfiguration, siteId)
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpoint)
	var params = map[string]string{}
	if caid == "" {
		params = nil
	} else {
		params["caid"] = caid
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadFastRenewalConfiguration)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when reading fast renewal configuration %s: %s", siteId, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when reading fast renewal configuration. site ID: %s", resp.StatusCode, siteId)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula Get fast renewal configuration response: %s\n", string(responseBody))

	var fastRenewalConfigurationResponse FastRenewalConfigurationDto
	err = json.Unmarshal([]byte(responseBody), &fastRenewalConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get fast renewal configuration response for site Id %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &fastRenewalConfigurationResponse, nil
}

func (c *Client) EnableFastRenewalConfiguration(siteId string, caid string) (*FastRenewalConfigurationDto, error) {
	if siteId == "" {
		fmt.Errorf("[ERROR] site ID was not provided")
	}
	endpoint := fmt.Sprintf(endpointFastRenewalConfiguration, siteId)
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpoint)
	var params = map[string]string{}
	if caid == "" {
		params = nil
	} else {
		params["caid"] = caid
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, nil, params, CreateFastRenewalConfiguration)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when enabling fast renewal configuration %s: %s", siteId, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when enabling fast renewal configuration. site ID: %s", resp.StatusCode, siteId)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula create fast renewal configuration response: %s\n", string(responseBody))

	var fastRenewalConfigurationResponse FastRenewalConfigurationDto
	err = json.Unmarshal([]byte(responseBody), &fastRenewalConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get fast renewal create configuration response for site Id %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &fastRenewalConfigurationResponse, nil
}

func (c *Client) DeleteFastRenewalConfiguration(siteId string, caid string) (*FastRenewalConfigurationDto, error) {
	if siteId == "" {
		fmt.Errorf("[ERROR] site ID was not provided")
	}
	endpoint := fmt.Sprintf(endpointFastRenewalConfiguration, siteId)
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpoint)
	var params = map[string]string{}
	if caid == "" {
		params = nil
	} else {
		params["caid"] = caid
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, params, DeleteFastRenewalConfiguration)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error from Incapsula service when deleting fast renewal configuration. site ID: %s\n: %s", siteId, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when deleting fast renewal configuration. site ID: %s", resp.StatusCode, siteId)
	}

	resp, err = c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadFastRenewalConfiguration)
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Incapsula delete fast renewal configuration response: %s\n", string(responseBody))

	var fastRenewalConfigurationResponse FastRenewalConfigurationDto
	err = json.Unmarshal([]byte(responseBody), &fastRenewalConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error parsing get fast renewal delete configuration response for site Id %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when parsing response on deleting fast renewal configuration. site ID: %s", resp.StatusCode, siteId)
	}

	return &fastRenewalConfigurationResponse, nil
}
