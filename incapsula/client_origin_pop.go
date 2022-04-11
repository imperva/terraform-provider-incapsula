package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// SetOriginPOPResponse contains the relevant site information when setting an Incapsula Origin POP
type SetOriginPOPResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	DebugInfo  struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// SetOriginPOP sets the origin POP for given data center
func (c *Client) SetOriginPOP(dcID int, originPOP string) error {
	reqURL := fmt.Sprintf("%s/sites/datacenter/origin-pop/modify?dc_id=%d", c.config.BaseURL, dcID)
	if originPOP != "" {
		reqURL = fmt.Sprintf("%s&origin_pop=%s", reqURL, originPOP)
	}
	// Post request to Incapsula
	resp, err := c.DoJsonRequestWithHeaders(http.MethodPost, reqURL, nil, UpdateOriginPop)
	if err != nil {
		return fmt.Errorf("Error from Incapsula service when setting origin POP: %s for data center: %d: %s", originPOP, dcID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula SetOriginPOP JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var originPOPResponse SetOriginPOPResponse
	err = json.Unmarshal([]byte(responseBody), &originPOPResponse)
	if err != nil {
		return fmt.Errorf("Error parsing origin POP JSON response for origin POP: %s for data center: %d: %s", originPOP, dcID, err)
	}

	// Look at the response status code from Incapsula
	if originPOPResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when updating origin POP: %s for data center: %d: %s", originPOP, dcID, string(responseBody))
	}

	return nil
}
