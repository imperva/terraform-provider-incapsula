package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
)

// TXTResponse contains the relevant information when creating/updating/deleting a TXT record
type TXTResponse struct {
	Res        int    `json:"res"`
	ResMessage string `json:"res_message"`
	DebugInfo  struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// TXTRecordResponse contains the relevant information when getting a TXT record
type TXTRecordResponse struct {
	SiteID              int    `json:"site_id"`
	TxtRecordValueOne   string `json:"txt_record_value_one"`
	TxtRecordValueTwo   string `json:"txt_record_value_two"`
	TxtRecordValueThree string `json:"txt_record_value_three"`
	TxtRecordValueFour  string `json:"txt_record_value_four"`
	TxtRecordValueFive  string `json:"txt_record_value_five"`
	Res                 int    `json:"res"`
	ResMessage          string `json:"res_message"`
	DebugInfo           struct {
		IDInfo string `json:"id-info"`
	} `json:"debug_info"`
}

// GetTXTRecords gets the site TXT Records
func (c *Client) GetTXTRecords(siteID string) (*TXTRecordResponse, error) {
	log.Printf("[INFO] Getting Incapsula TXT record(s) for Site ID %s\n", siteID)

	// GET records from Incapsula
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/sites/%s/settings/general/additionalTxtRecords?api_id=%s&api_key=%s", c.config.BaseURLRev2, siteID, c.config.APIID, c.config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading TXT record(s) for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula GetTXTRecords JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading TXT record(s) for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var txtRecords TXTRecordResponse
	err = json.Unmarshal([]byte(responseBody), &txtRecords)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap TXT record(s) JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &txtRecords, nil
}

// CreateTXTRecord creates TXT record(s) on the site ID.
func (c *Client) CreateTXTRecord(siteID string, txtRecordValueOne string, txtRecordValueTwo string, txtRecordValueThree string, txtRecordValueFour string, txtRecordValueFive string) error {
	log.Printf("[INFO] Create Incapsula TXT record(s) for Site ID %s\n txt_record_value_one=%s, txt_record_value_two=%s, txt_record_value_three=%s, txt_record_value_four=%s, txt_record_value_five=%s", siteID, txtRecordValueOne, txtRecordValueTwo, txtRecordValueThree, txtRecordValueFour, txtRecordValueFive)

	// Post request to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("#{c.config.BaseURLRev2}/sites/#{siteID}/settings/general/additionalTxtRecords"), url.Values{
		"api_id":                 {c.config.APIID},
		"api_key":                {c.config.APIKey},
		"txt_record_value_one":   {txtRecordValueOne},
		"txt_record_value_two":   {txtRecordValueTwo},
		"txt_record_value_three": {txtRecordValueThree},
		"txt_record_value_four":  {txtRecordValueFour},
		"txt_record_value_five":  {txtRecordValueFive},
	})

	if err != nil {
		return fmt.Errorf("Error creating TXT record(s) for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula CreateTXTRecord JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when creating TXT record(s) for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	return nil
}

// UpdateTXTRecord updates TXT record(s) on the site ID.
func (c *Client) UpdateTXTRecord(siteID string, txtRecordValueOne string, txtRecordValueTwo string, txtRecordValueThree string, txtRecordValueFour string, txtRecordValueFive string) (*TXTResponse, error) {
	log.Printf("[INFO] Update Incapsula TXT record(s) for Site ID %s\n "+
		"txt_record_value_one=%s, txt_record_value_two=%s, txt_record_value_three=%s, txt_record_value_four=%s, txt_record_value_five=%s",
		siteID, txtRecordValueOne, txtRecordValueTwo, txtRecordValueThree, txtRecordValueFour, txtRecordValueFive)

	// Post request to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("%s/sites/%s/settings/general/additionalTxtRecords", c.config.BaseURLRev2, siteID), url.Values{
		"api_id":                 {c.config.APIID},
		"api_key":                {c.config.APIKey},
		"txt_record_value_one":   {txtRecordValueOne},
		"txt_record_value_two":   {txtRecordValueTwo},
		"txt_record_value_three": {txtRecordValueThree},
		"txt_record_value_four":  {txtRecordValueFour},
		"txt_record_value_five":  {txtRecordValueFive},
	})

	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when updating TXT record(s) for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula UpdateTXTRecord JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating TXT record(s) for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var txtResponse TXTResponse
	err = json.Unmarshal([]byte(responseBody), &txtResponse)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap TXT response JSON response for Site ID %s: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &txtResponse, nil
}

// DeleteTXTRecord deletes TXT record(s) on the site ID.
func (c *Client) DeleteTXTRecord(siteID string, txtRecordValueOne string, txtRecordValueTwo string, txtRecordValueThree string, txtRecordValueFour string, txtRecordValueFive string) error {
	log.Printf("[INFO] Update Incapsula TXT record(s) for Site ID %s\n "+
		"txt_record_value_one=%s, txt_record_value_two=%s, txt_record_value_three=%s, txt_record_value_four=%s, txt_record_value_five=%s",
		siteID, txtRecordValueOne, txtRecordValueTwo, txtRecordValueThree, txtRecordValueFour, txtRecordValueFive)

	// Post request to Incapsula
	resp, err := c.httpClient.PostForm(fmt.Sprintf("#{c.config.BaseURLRev2}/sites/#{siteID}/settings/general/additionalTxtRecords"), url.Values{
		"api_id":                 {c.config.APIID},
		"api_key":                {c.config.APIKey},
		"txt_record_value_one":   {txtRecordValueOne},
		"txt_record_value_two":   {txtRecordValueTwo},
		"txt_record_value_three": {txtRecordValueThree},
		"txt_record_value_four":  {txtRecordValueFour},
		"txt_record_value_five":  {txtRecordValueFive},
	})

	if err != nil {
		return fmt.Errorf("Error deleting TXT record(s) for Site ID %s: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula DeleteTXTRecord JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting TXT record(s) for Site ID %s: %s", resp.StatusCode, siteID, string(responseBody))
	}

	return nil
}
