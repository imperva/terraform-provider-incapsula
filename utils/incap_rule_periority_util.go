package utils

type RuleType int

const (
	SIMPLIFIED_REDIRECT RuleType = iota
	REDIRECT
	REWRITE
	REWRITE_RESPONSE
	FORWARD
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
		return "Unknown"
	}
}

type RuleDetails struct {
	Filter          string `json:"filter"`
	From            string `json:"from"`
	To              string `json:"to"`
	HeaderName      string `json:"header_name"`
	RewriteExisting bool   `json:"rewrite_existing"`
	AddIfMissing    bool   `json:"add_if_missing"`
	Name            string `json:"name"`
	Action          string `json:"action"`
	Enabled         bool   `json:"enabled"`
}
