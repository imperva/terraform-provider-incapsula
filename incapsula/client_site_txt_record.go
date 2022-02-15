package incapsula

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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

// ReadTXTRecords gets the site TXT Records
func (c *Client) ReadTXTRecords(siteID int) (*TXTRecordResponse, error) {
	log.Printf("[INFO] Getting Incapsula TXT record(s) for siteID %d\n", siteID)

	// GET records from Incapsula
	reqURL := fmt.Sprintf("%s/sites/%d/settings/general/additionalTxtRecords", c.config.BaseURLRev2, siteID)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Error from Incapsula service when reading TXT record(s) for siteID: %d\n %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula GetTXTRecords JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when reading TXT record(s) for siteID: %d\n%s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var txtRecords TXTRecordResponse
	err = json.Unmarshal([]byte(responseBody), &txtRecords)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap TXT record(s) JSON response for siteID: %d\n%s\nresponse: %s", siteID, err, string(responseBody))
	}

	response := []byte(responseBody)
	if strings.Contains(string(response), "no TXT records") {
		return nil, fmt.Errorf("[ERROR] The Text Records for Site ID %d does not exist", siteID)
	}

	return &txtRecords, nil
}

// CreateTXTRecord creates TXT record(s) on the siteID.
func (c *Client) CreateTXTRecord(siteID int, txtRecordValueOne string, txtRecordValueTwo string, txtRecordValueThree string, txtRecordValueFour string, txtRecordValueFive string) (*TXTResponse, error) {
	log.Printf("[INFO] Create Incapsula TXT record(s) for siteID %d\n txt_record_value_one=%s, txt_record_value_two=%s, txt_record_value_three=%s, txt_record_value_four=%s, txt_record_value_five=%s", siteID, txtRecordValueOne, txtRecordValueTwo, txtRecordValueThree, txtRecordValueFour, txtRecordValueFive)

	// Post request to Incapsula
	values := url.Values{
		"txt_record_value_one":   {txtRecordValueOne},
		"txt_record_value_two":   {txtRecordValueTwo},
		"txt_record_value_three": {txtRecordValueThree},
		"txt_record_value_four":  {txtRecordValueFour},
		"txt_record_value_five":  {txtRecordValueFive},
	}

	// Post request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%d/settings/general/additionalTxtRecords", c.config.BaseURLRev2, siteID)
	resp, err := c.PostFormWithHeaders(reqURL, values)
	if err != nil {
		return nil, fmt.Errorf("Error creating TXT record(s) for siteID %d: %s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula CreateTXTRecord JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating TXT record(s) for siteID: %d\n%s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var txtResponse TXTResponse
	err = json.Unmarshal([]byte(responseBody), &txtResponse)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap TXT response JSON response for siteID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &txtResponse, nil
}

// UpdateTXTRecord updates TXT record(s) on the siteID.
func (c *Client) UpdateTXTRecord(siteID int, txtRecordValueOne string, txtRecordValueTwo string, txtRecordValueThree string, txtRecordValueFour string, txtRecordValueFive string) (*TXTResponse, error) {
	log.Printf("[INFO] Update Incapsula TXT record(s) for siteID %d\n txt_record_value_one=%s, txt_record_value_two=%s, txt_record_value_three=%s, txt_record_value_four=%s, txt_record_value_five=%s",
		siteID, txtRecordValueOne, txtRecordValueTwo, txtRecordValueThree, txtRecordValueFour, txtRecordValueFive)

	// Post request to Incapsula
	values := url.Values{
		"txt_record_value_one":   {txtRecordValueOne},
		"txt_record_value_two":   {txtRecordValueTwo},
		"txt_record_value_three": {txtRecordValueThree},
		"txt_record_value_four":  {txtRecordValueFour},
		"txt_record_value_five":  {txtRecordValueFive},
	}

	// Post request to Incapsula
	reqURL := fmt.Sprintf("%s/sites/%d/settings/general/additionalTxtRecords", c.config.BaseURLRev2, siteID)
	resp, err := c.PostFormWithHeaders(reqURL, values)
	if err != nil {
		return nil, fmt.Errorf("Error updating TXT record(s) for siteID: %d\n%s", siteID, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula UpadteTXTRecord JSON response: %s\n", string(responseBody))

	// Check the response code
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error status code %d from Incapsula service when updating TXT record(s) for siteID: %d\n%s", resp.StatusCode, siteID, string(responseBody))
	}

	// Parse the JSON
	var txtResponse TXTResponse
	err = json.Unmarshal([]byte(responseBody), &txtResponse)

	if err != nil {
		return nil, fmt.Errorf("Error parsing Incap TXT response JSON response for siteID %d: %s\nresponse: %s", siteID, err, string(responseBody))
	}

	return &txtResponse, nil
}

// DeleteTXTRecord deletes TXT record(s) on the siteID.
func (c *Client) DeleteTXTRecord(siteID int, recordNumber string) error {
	log.Printf("[INFO] Delete Incapsula TXT record number %s for siteID %d\n ", recordNumber, siteID)

	reqURL := fmt.Sprintf("%s/sites/%d/settings/general/additionalTxtRecords?record_number=%s", c.config.BaseURLRev2, siteID, recordNumber)
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil)
	if err != nil {
		log.Printf("[DEBUG] Error deleting TXT record for siteID %d: %s", siteID, err)
		return fmt.Errorf("Error deleting TXT record for siteID %d", siteID)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula DeleteTXTRecord JSON response: %s\n", string(responseBody))

	response := []byte(responseBody)
	// Check the response code
	// The response code of successful request is 400
	if resp.StatusCode != 400 && !strings.Contains(string(response), "OK") {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting TXT record for siteID %d: %s", resp.StatusCode, siteID, string(responseBody))
	}

	return nil
}

func (c *Client) DeleteTXTRecordAll(siteID int) error {
	log.Printf("[INFO] Delete Incapsula All TXT records for siteID %s\n ", siteID)
	reqURL := fmt.Sprintf("%s/sites/%d/settings/general/additionalTxtRecords/delete-all", c.config.BaseURLRev2, siteID)
	
	resp, err := c.DoJsonRequestWithHeaders(http.MethodDelete, reqURL, nil)
	if err != nil {
		return fmt.Errorf("Error deleting TXT records for siteID %d: %v", siteID, err)
	}
	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Incapsula DeleteTXTRecords JSON response: %s\n", string(responseBody))

	response := []byte(responseBody)
	// Check the response code
	if resp.StatusCode != 400 && !strings.Contains(string(response), "OK") {
		return fmt.Errorf("Error status code %d from Incapsula service when deleting all "+
			"TXT records for siteID %d: %s", resp.StatusCode, siteID, string(responseBody))
	}

	return nil
}
