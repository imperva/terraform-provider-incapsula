package incapsula

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type OriginServerStruct struct {
	Address    string `json:"address"`
	IsEnabled  bool   `json:"isEnabled"`
	ServerMode string `json:"serverMode"`
	Weight     *int   `json:"weight"`
}

type DataCenterStruct struct {
	Name                string               `json:"name"`
	ID                  *int                 `json:"id"`
	IpMode              string               `json:"ipMode"`
	WebServersPerServer *int                 `json:"webServersPerServer"`
	DcLbAlgorithm       string               `json:"lbAlgorithm"`
	Weight              *int                 `json:"weight"`
	IsEnabled           bool                 `json:"isEnabled"`
	IsActive            bool                 `json:"isActive"`
	IsContent           bool                 `json:"isContent"`
	IsRestOfTheWorld    bool                 `json:"isRestOfTheWorld"`
	GeoLocations        []string             `json:"geoLocations"`
	OriginPoP           string               `json:"originPop"`
	OriginServers       []OriginServerStruct `json:"servers"`
}

type DataCentersStruct struct {
	SiteLbAlgorithm                    string             `json:"lbAlgorithm"`
	FailOverRequiredMonitors           string             `json:"failOverRequiredMonitors"`
	DataCenterMode                     string             `json:"dataCenterMode"`
	MinAvailableServersForDataCenterUp int                `json:"minAvailableServersForDataCenterUp"`
	KickStartURL                       string             `json:"kickStartURL"`
	KickStartUser                      string             `json:"kickStartUser"`
	KickStartPass                      string             `json:"kickStartPass"`
	IsPersistent                       bool               `json:"isPersistent"`
	DataCenters                        []DataCenterStruct `json:"dataCenters"`
}

type ApiErrorSource struct {
	Pointer   string `json:"pointer"`
	Parameter string `json:"parameter"`
}

type ApiError struct {
	ID      string         `json:"id"`
	Status  string         `json:"status"`
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Source  ApiErrorSource `json:"source"`
}

// Same DTO for: GET response, PUT request, and PUT response
type DataCentersConfigurationDTO struct {
	Errors []ApiError          `json:"errors"`
	Data   []DataCentersStruct `json:"data"`
}

// AddDataCenter adds an incap rule to be managed by Incapsula
func (c *Client) PutDataCentersConfiguration(siteID string, requestDTO DataCentersConfigurationDTO) (*DataCentersConfigurationDTO, error) {
	log.Printf("[INFO] Updating Incapsula data centers configuration for siteID: %s\n", siteID)

	baseURLv3 := c.config.BaseURL[:len(c.config.BaseURL)-3] + "/v3"
	dcsJSON, err := json.Marshal(requestDTO)
	if err != nil {
		return nil, fmt.Errorf("Failed to JSON marshal Data Centers configuration: %s", err)
	}

	// Prepare HTTP client PUT request to incapsula
	log.Printf("[DEBUG] Incapsula Update Data Centers configuration request: %s\n", string(dcsJSON))
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/sites/%s/data-centers-configuration", baseURLv3, siteID),
		bytes.NewReader(dcsJSON))
	if err != nil {
		return nil, fmt.Errorf("Error preparing HTTP PUT for updating Data Centers configuration with ID %s: %s", siteID, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-api-id", c.config.APIID)
	req.Header.Set("x-api-key", c.config.APIKey)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error executing update Data Centers configuration request for siteID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Update Data Centers configuration JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO DataCentersConfigurationDTO
	err = json.Unmarshal([]byte(responseBody), &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing update Data Centers configuration JSON response for siteID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &responseDTO, nil
}

// ListDataCenters gets the Incapsula list of data centers
func (c *Client) GetDataCentersConfiguration(siteID string) (*DataCentersConfigurationDTO, error) {
	log.Printf("[INFO] Getting Data Centers configuration (site_id: %s)\n", siteID)

	// Get request to Incapsula
	baseURLv3 := c.config.BaseURL[:len(c.config.BaseURL)-3] + "/v3"
	// Prepare HTTP client PUT request to incapsula
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/sites/%s/data-centers-configuration", baseURLv3, siteID), nil)
	if err != nil {
		return nil, fmt.Errorf("Error preparing HTTP GET for getting Data Centers configuration with ID %s: %s", siteID, err)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-api-id", c.config.APIID)
	req.Header.Set("x-api-key", c.config.APIKey)
	resp, err := c.httpClient.Do(req)

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula data centers JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var responseDTO DataCentersConfigurationDTO
	err = json.Unmarshal([]byte(responseBody), &responseDTO)
	if err != nil {
		return nil, fmt.Errorf("Error parsing data centers list JSON response for siteID: %s %s\nresponse: %s", siteID, err, string(responseBody))
	}

	// Look at the response status code from Incapsula
	if responseDTO.Errors != nil && len(responseDTO.Errors) > 0 {
		out, err := json.Marshal(responseDTO.Errors)
		if err != nil {
			panic(err)
		}
		return &responseDTO, fmt.Errorf("Error from Incapsula service when getting data centers configuration (site_id: %s): %s", siteID, string(out))
	}

	return &responseDTO, nil
}
