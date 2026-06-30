resource "incapsula_abp_condition" "specific_visitor" {
  account_id  = var.account_id
  name        = "Specific visitor"
  description = "Match specific UA coming from selected IPs"
  code        = "(all headers.user_agent? (matches headers.user_agent re\"Mozilla\") (in visitor_ip 1.2.3.4 2.3.4.5))"
}

# When a template other than "custom" is set, `code` must follow that
# template's expected structure or the condition will fail to create.
resource "incapsula_abp_condition" "rate_limit" {
  account_id  = var.account_id
  name        = "Too many requests"
  description = "Trip when a session exceeds 120 requests per minute"
  template    = "rate_limiting"
  code        = "(all (not flags.no_token) (requests_per_minute > 120))"
}
