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
	originPopModifyUrl := ""
	if originPOP != "" {
		originPopModifyUrl = fmt.Sprintf("%s/sites/datacenter/origin-pop/modify?api_id=%s&api_key=%s&origin_pop=%s&dc_id=%d", c.config.BaseURL, c.config.APIID, c.config.APIKey, originPOP, dcID)
	} else { // setting the origin pop to NONE is done by not sending the origin_pop query param
		originPopModifyUrl = fmt.Sprintf("%s/sites/datacenter/origin-pop/modify?api_id=%s&api_key=%s&dc_id=%d", c.config.BaseURL, c.config.APIID, c.config.APIKey, dcID)
	}

	// Post request to Incapsula
	req, err := http.NewRequest(
		http.MethodPost,
		originPopModifyUrl,
		nil)
	if err != nil {
		return fmt.Errorf("Error preparing HTTP POST for setting Incapsula origin POP: %s for data center: %d: %s", originPOP, dcID, err)
	}
	resp, err := c.httpClient.Do(req)
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
