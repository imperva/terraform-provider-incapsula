package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const endpointShortRenewalCycleConfiguration = "/certificates-ui/v3/sites/%s/certificate/shortRenewalCycle"

type ShortRenewalCycleConfigurationDto struct {
	Errors []ApiErrorResponse               `json:"errors"`
	Data   []ShortRenewalCycleConfiguration `json:"data"`
}

type ShortRenewalCycleConfiguration struct {
	ShortRenewalCycle bool `json:"shortRenewalCycle"`
}

func (c *Client) GetShortRenewalCycleConfiguration(siteId string, caid string) (*ShortRenewalCycleConfigurationDto, error) {
	if siteId == "" {
		fmt.Errorf("[ERROR] site id was not provided")
	}
	endpoint := fmt.Sprintf(endpointShortRenewalCycleConfiguration, siteId)
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpoint)
	var params = map[string]string{}
	if caid == "" {
		params = nil
	} else {
		params["caid"] = caid
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadShortRenewalCycleConfiguration)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error from Incapsula service when reading short renewal cycle configuration %s: %s", siteId, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when reading short renewal cycle configuration. site id: %s", resp.StatusCode, siteId)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] incapsula Get short renewal cycle configuration response: %s\n", string(responseBody))

	var shortRenewalCycleConfigurationResponse ShortRenewalCycleConfigurationDto
	err = json.Unmarshal([]byte(responseBody), &shortRenewalCycleConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error parsing get short renewal cycle configuration response for site id %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &shortRenewalCycleConfigurationResponse, nil
}

func (c *Client) EnableShortRenewalCycleConfiguration(siteId string, caid string) (*ShortRenewalCycleConfigurationDto, error) {
	if siteId == "" {
		fmt.Errorf("[ERROR] site id was not provided")
	}
	endpoint := fmt.Sprintf(endpointShortRenewalCycleConfiguration, siteId)
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpoint)
	var params = map[string]string{}
	if caid == "" {
		params = nil
	} else {
		params["caid"] = caid
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPost, reqURL, nil, params, CreateShortRenewalCycleConfiguration)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error from Incapsula service when enabling short renewal cycle configuration %s: %s", siteId, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when enabling short renewal cycle configuration. site id: %s", resp.StatusCode, siteId)
	}

	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] incapsula create short renewal cycle configuration response: %s\n", string(responseBody))

	var shortRenewalCycleConfigurationResponse ShortRenewalCycleConfigurationDto
	err = json.Unmarshal([]byte(responseBody), &shortRenewalCycleConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error parsing get short renewal cycle create configuration response for site id %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	return &shortRenewalCycleConfigurationResponse, nil
}

func (c *Client) DeleteShortRenewalCycleConfiguration(siteId string, caid string) (*ShortRenewalCycleConfigurationDto, error) {
	if siteId == "" {
		fmt.Errorf("[ERROR] site id was not provided")
	}
	endpoint := fmt.Sprintf(endpointShortRenewalCycleConfiguration, siteId)
	reqURL := fmt.Sprintf("%s%s", c.config.BaseURLAPI, endpoint)
	var params = map[string]string{}
	if caid == "" {
		params = nil
	} else {
		params["caid"] = caid
	}
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodDelete, reqURL, nil, params, DeleteShortRenewalCycleConfiguration)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error from Incapsula service when deleting short renewal cycle configuration. site id: %s\n: %s", siteId, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when deleting short renewal cycle configuration. site id: %s", resp.StatusCode, siteId)
	}

	resp, err = c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, params, ReadShortRenewalCycleConfiguration)
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] incapsula delete short renewal cycle configuration response: %s\n", string(responseBody))

	var shortRenewalCycleConfigurationResponse ShortRenewalCycleConfigurationDto
	err = json.Unmarshal([]byte(responseBody), &shortRenewalCycleConfigurationResponse)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] error parsing get short renewal cycle delete configuration response for site id %s: %s\nresponse: %s", siteId, err, string(responseBody))
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[ERROR] bad status code %d from Incapsula service when parsing response on deleting short renewal cycle configuration. site id: %s", resp.StatusCode, siteId)
	}

	return &shortRenewalCycleConfigurationResponse, nil
}
