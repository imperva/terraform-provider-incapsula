package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

const endpointSiteLogLevel = "sites/setlog"

// UpdateLogLevel will update the site log level
func (c *Client) UpdateLogLevel(siteID, logLevel, logsAccountId string) error {
	type LogLevelResponse struct {
		Res        int    `json:"res"`
		ResMessage string `json:"res_message"`
		DebugInfo  struct {
			LogLevel      string `json:"log_level,omitempty"`
			LogsAccountId string `json:"logs_account_id,omitempty"`
		} `json:"debug_info"`
	}

	log.Printf("[INFO] Updating Incapsula log level (%s) for siteID: %s\n", logLevel, siteID)

	// Post form to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/%s", c.config.BaseURL, endpointSiteLogLevel), url.Values{
		"api_id":          {c.config.APIID},
		"api_key":         {c.config.APIKey},
		"site_id":         {siteID},
		"log_level":       {logLevel},
		"logs_account_id": {logsAccountId},
	})
	if err != nil {
		return fmt.Errorf("Error updating log level (%s) on site_id: %s: %s", logLevel, siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula update log level JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var logLevelResponse LogLevelResponse
	err = json.Unmarshal([]byte(responseBody), &logLevelResponse)
	if err != nil {
		return fmt.Errorf("Error parsing update log level JSON response for siteID %s: %s", siteID, err)
	}

	// Look at the response status code from Incapsula
	if logLevelResponse.Res != 0 {
		return fmt.Errorf("Error from Incapsula service when updating log level for siteID %s: %s", siteID, string(responseBody))
	}

	return nil
}
