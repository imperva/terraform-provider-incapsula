package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
)

// endpoints
const advancedCacheEndpoint = "sites/performance/response-headers"

// struct
type AddCacheHeaderResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	DebugInfo  struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// ConfigureAdvanceCache adds an WAF Sec rule
func (c *Client) ConfigureAdvanceCache(siteID int, cacheHeaders, cacheAllHeaders string) (*AddCacheHeaderResponse, error) {
	log.Printf("[INFO] Configuring Incapsula Advance Cacheing site id: %d\n", siteID)

	// Base URL values
	values := url.Values{
		"api_id":            {c.config.APIID},
		"api_key":           {c.config.APIKey},
		"site_id":           {strconv.Itoa(siteID)},
		"cache_headers":     {cacheHeaders},
		"cache_all_headers": {cacheAllHeaders},
	}

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, advancedCacheEndpoint), values)
	if err != nil {
		return nil, fmt.Errorf("Error adding  Advance WAF Caching for site id %d", siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula Advance Caching JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var addCacheHeaderResponse AddCacheHeaderResponse
	err = json.Unmarshal([]byte(responseBody), &addCacheHeaderResponse)
	if err != nil {
		return nil, fmt.Errorf("Error parsing add Advances Cache JSON response for site id %d", siteID)
	}

	// Look at the response status code from Incapsula
	if addCacheHeaderResponse.Res != 0 {
		return nil, fmt.Errorf("Error from Incapsula service when adding Advance Caching for site id %d: %s", siteID, string(responseBody))
	}

	return &addCacheHeaderResponse, nil
}
