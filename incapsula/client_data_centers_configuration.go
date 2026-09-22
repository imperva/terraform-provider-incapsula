package incapsula

import (
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
// accountID is optional - pass 0 to operate on the account the API credentials belong to.
// It must be provided when the site belongs to a different account, e.g. when a parent
// account's API key manages a site that resides in one of its sub accounts.
func (c *Client) PutDataCentersConfiguration(siteID string, accountID int, requestDTO DataCentersConfigurationDTO) (*DataCentersConfigurationDTO, error) {
	log.Printf("[INFO] Updating Incapsula data centers configuration for siteID: %s, accountID: %d\n", siteID, accountID)

	baseURLv3 := c.config.BaseURL[:len(c.config.BaseURL)-3] + "/v3"
	dcsJSON, err := json.Marshal(requestDTO)
	reqURL := fmt.Sprintf("%s/sites/%s/data-centers-configuration", baseURLv3, siteID)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodPut, reqURL, dcsJSON, GetRequestParamsWithCaid(accountID), CreateDataCenterConfiguration)
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
// accountID is optional - pass 0 to operate on the account the API credentials belong to.
// See PutDataCentersConfiguration.
func (c *Client) GetDataCentersConfiguration(siteID string, accountID int) (*DataCentersConfigurationDTO, error) {
	log.Printf("[INFO] Getting Data Centers configuration (site_id: %s, account_id: %d)\n", siteID, accountID)

	// Get request to Incapsula
	baseURLv3 := c.config.BaseURL[:len(c.config.BaseURL)-3] + "/v3"
	reqURL := fmt.Sprintf("%s/sites/%s/data-centers-configuration", baseURLv3, siteID)
	resp, err := c.DoJsonAndQueryParamsRequestWithHeaders(http.MethodGet, reqURL, nil, GetRequestParamsWithCaid(accountID), ReadDataCenterConfiguration)
	if err != nil {
		return nil, fmt.Errorf("Error executing get Data Centers configuration request for siteID %s: %s", siteID, err)
	}

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

	return &responseDTO, nil
}
