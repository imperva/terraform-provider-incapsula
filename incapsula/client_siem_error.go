package incapsula

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type siemApiErrorEntry struct {
	ID     string `json:"id"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

type siemApiErrorResponse struct {
	Errors []siemApiErrorEntry `json:"errors"`
}

// formatSiemApiError turns a siem-config-service errors[] body into a readable message.
// Falls back to the raw body when it is not in the standard format.
func formatSiemApiError(operation string, entity string, statusCode int, body string) error {
	var response siemApiErrorResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil || len(response.Errors) == 0 {
		return fmt.Errorf("received failure response for operation: %s on %s\nstatus code: %d\nbody: %s",
			operation, entity, statusCode, body)
	}

	var message strings.Builder
	fmt.Fprintf(&message, "received failure response for operation: %s on %s\nstatus code: %d", operation, entity, statusCode)
	for _, apiError := range response.Errors {
		if apiError.Code != "" {
			fmt.Fprintf(&message, "\ncode: %s", apiError.Code)
		}
		// 5xx responses carry no detail by design, so the title is the only description available.
		if apiError.Detail != "" {
			fmt.Fprintf(&message, "\ndetail: %s", apiError.Detail)
		} else if apiError.Title != "" {
			fmt.Fprintf(&message, "\ntitle: %s", apiError.Title)
		}
		if apiError.ID != "" {
			fmt.Fprintf(&message, "\ntrace id: %s", apiError.ID)
		}
	}
	return errors.New(message.String())
}
