package utils

type Response struct {
	Category    string        `json:"category"`
	RuleDetails []RuleDetails `json:"data"`
	Errors      []APIErrors   `json:"errors"`
}

type RuleDetails struct {
	Id                      *int   `json:"id"`
	Name                    string `json:"rule_name"`
	Action                  string `json:"action"`
	Filter                  string `json:"filter,omitempty"`
	AddMissing              *bool  `json:"add_if_missing,omitempty"`
	From                    string `json:"from,omitempty"`
	To                      string `json:"to,omitempty"`
	ResponseCode            *int   `json:"response_code,omitempty"`
	RewriteExisting         *bool  `json:"rewrite_existing,omitempty"`
	RewriteName             string `json:"rewrite_name,omitempty"`
	CookieName              string `json:"cookie_name"`
	HeaderName              string `json:"header_name"`
	DCID                    *int   `json:"dc_id,omitempty"`
	PortForwardingContext   string `json:"port_forwarding_context,omitempty"`
	PortForwardingValue     string `json:"port_forwarding_value,omitempty"`
	ErrorType               string `json:"error_type,omitempty"`
	ErrorResponseFormat     string `json:"error_response_format,omitempty"`
	ErrorResponseData       string `json:"error_response_data,omitempty"`
	MultipleHeaderDeletions *bool  `json:"multiple_headers_deletion"`
	Enabled                 *bool  `json:"enabled"`
}

type APIErrors struct {
	Status int               `json:"status"`
	Id     string            `json:"id"`
	Code   int               `json:"code"`
	Source map[string]string `json:"source"`
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
}
