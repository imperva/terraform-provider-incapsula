package utils

import "fmt"

type RuleType int
type RuleAction int

const (
	SIMPLIFIED_REDIRECT RuleType = iota
	REDIRECT
	REWRITE
	REWRITE_RESPONSE
	FORWARD
)
const (
	RULE_ACTION_SIMPLIFIED_REDIRECT RuleAction = iota
	RULE_ACTION_REDIRECT
	RULE_ACTION_REWRITE_HEADER
	RULE_ACTION_REWRITE_COOKIE
	RULE_ACTION_REWRITE_URL
	RULE_ACTION_DELETE_HEADER
	RULE_ACTION_DELETE_COOKIE
	RULE_ACTION_RESPONSE_REWRITE_HEADER
	RULE_ACTION_RESPONSE_DELETE_HEADER
	RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE
	RULE_ACTION_CUSTOM_ERROR_RESPONSE
	RULE_ACTION_FORWARD_TO_DC
	RULE_ACTION_FORWARD_TO_PORT
)

func (ruleType RuleType) ValidRule() bool {
	switch ruleType {

	case SIMPLIFIED_REDIRECT:
		return true

	case REDIRECT:
		return true

	case REWRITE:
		return true

	case REWRITE_RESPONSE:
		return true

	case FORWARD:
		return true

	default:
		return false
	}
}
func (ruleType RuleType) String() string {
	switch ruleType {

	case SIMPLIFIED_REDIRECT:
		return "SIMPLIFIED_REDIRECT"

	case REDIRECT:
		return "REDIRECT"

	case REWRITE:
		return "REWRITE"

	case REWRITE_RESPONSE:
		return "REWRITE_RESPONSE"

	case FORWARD:
		return "FORWARD"

	default:
		return "Unknown Rule"
	}
}
func RuleFromString(rule string) (RuleType, error) {
	switch rule {
	case "SIMPLIFIED_REDIRECT":
		return SIMPLIFIED_REDIRECT, nil
	case "REDIRECT":
		return REDIRECT, nil
	case "REWRITE":
		return REWRITE, nil
	case "REWRITE_RESPONSE":
		return REWRITE_RESPONSE, nil
	case "FORWARD":
		return FORWARD, nil
	default:
		return -1, fmt.Errorf("Invalid Rule: %s", rule)
	}
}
func (ruleAction RuleAction) String() string {
	switch ruleAction {
	case RULE_ACTION_SIMPLIFIED_REDIRECT:
		return "RULE_ACTION_SIMPLIFIED_REDIRECT"
	case RULE_ACTION_DELETE_COOKIE:
		return "RULE_ACTION_DELETE_COOKIE"
	case RULE_ACTION_CUSTOM_ERROR_RESPONSE:
		return "RULE_ACTION_CUSTOM_ERROR_RESPONSE"
	case RULE_ACTION_FORWARD_TO_PORT:
		return "RULE_ACTION_FORWARD_TO_PORT"
	case RULE_ACTION_REDIRECT:
		return "RULE_ACTION_REDIRECT"
	case RULE_ACTION_REWRITE_HEADER:
		return "RULE_ACTION_REWRITE_HEADER"
	case RULE_ACTION_REWRITE_COOKIE:
		return "RULE_ACTION_REWRITE_COOKIE"
	case RULE_ACTION_REWRITE_URL:
		return "RULE_ACTION_REWRITE_URL"
	case RULE_ACTION_DELETE_HEADER:
		return "RULE_ACTION_DELETE_HEADER"
	case RULE_ACTION_RESPONSE_REWRITE_HEADER:
		return "RULE_ACTION_RESPONSE_REWRITE_HEADER"
	case RULE_ACTION_RESPONSE_DELETE_HEADER:
		return "RULE_ACTION_RESPONSE_DELETE_HEADER"
	case RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE:
		return "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE"
	case RULE_ACTION_FORWARD_TO_DC:
		return "RULE_ACTION_FORWARD_TO_DC"
	default:
		return "Unknown Rule Action"
	}
}
func RuleActionFromString(ruleAction string) (RuleAction, error) {
	switch ruleAction {
	case "RULE_ACTION_SIMPLIFIED_REDIRECT":
		return RULE_ACTION_SIMPLIFIED_REDIRECT, nil
	case "RULE_ACTION_DELETE_COOKIE":
		return RULE_ACTION_DELETE_COOKIE, nil
	case "RULE_ACTION_CUSTOM_ERROR_RESPONSE":
		return RULE_ACTION_CUSTOM_ERROR_RESPONSE, nil
	case "RULE_ACTION_FORWARD_TO_PORT":
		return RULE_ACTION_FORWARD_TO_PORT, nil
	case "RULE_ACTION_REDIRECT":
		return RULE_ACTION_REDIRECT, nil
	case "RULE_ACTION_REWRITE_HEADER":
		return RULE_ACTION_REWRITE_HEADER, nil
	case "RULE_ACTION_REWRITE_COOKIE":
		return RULE_ACTION_REWRITE_COOKIE, nil
	case "RULE_ACTION_REWRITE_URL":
		return RULE_ACTION_REWRITE_URL, nil
	case "RULE_ACTION_DELETE_HEADER":
		return RULE_ACTION_DELETE_HEADER, nil
	case "RULE_ACTION_RESPONSE_REWRITE_HEADER":
		return RULE_ACTION_RESPONSE_REWRITE_HEADER, nil
	case "RULE_ACTION_RESPONSE_DELETE_HEADER":
		return RULE_ACTION_RESPONSE_DELETE_HEADER, nil
	case "RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE":
		return RULE_ACTION_RESPONSE_REWRITE_RESPONSE_CODE, nil
	case "RULE_ACTION_FORWARD_TO_DC":
		return RULE_ACTION_FORWARD_TO_DC, nil
	default:
		return -1, fmt.Errorf("Inavlid Rule Action %s\n", ruleAction)
	}
}

type Response struct {
	RuleDetails []RuleDetails `json:"data"`
}

type RuleDetails struct {
	Name                  string     `json:"name"`
	Action                RuleAction `json:"action"`
	Filter                string     `json:"filter,omitempty"`
	AddMissing            bool       `json:"add_missing,omitempty"`
	From                  string     `json:"from,omitempty"`
	To                    string     `json:"to,omitempty"`
	ResponseCode          int        `json:"response_code,omitempty"`
	RewriteExisting       *bool      `json:"rewrite_existing,omitempty"`
	RewriteName           string     `json:"rewrite_name,omitempty"`
	DCID                  *int       `json:"dc_id,omitempty"`
	PortForwardingContext string     `json:"port_forwarding_context,omitempty"`
	PortForwardingValue   string     `json:"port_forwarding_value,omitempty"`
	RateContext           string     `json:"rate_context,omitempty"`
	RateInterval          *int       `json:"rate_interval,omitempty"`
	ErrorType             string     `json:"error_type,omitempty"`
	ErrorResponseFormat   string     `json:"error_response_format,omitempty"`
	ErrorResponseData     string     `json:"error_response_data,omitempty"`
	MultipleDeletions     *bool      `json:"multiple_deletions,omitempty"`
	OverrideWafRule       string     `json:"overrideWafRule,omitempty"`
	OverrideWafAction     string     `json:"overrideWafAction,omitempty"`
	Enabled               *bool      `json:"enabled"`
}

// define big json that contain all the value and use pointers.
