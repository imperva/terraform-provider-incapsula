
variable "id" {
  default = ""
}
variable "key" {
  default = ""
}
provider "incapsula" {
  api_id = "${var.id}"
  api_key = "${var.key}"
}

resource "incapsula_site" "example-site" {
  domain = "mojo.beer.center" //Good for adding site and in response object
  account_id = "2398" //Good for adding site and in response object
  ref_id = "666" //Good for adding and updating site and in response object
  send_site_setup_emails = "true" //Good for adding site, not availble in site onject
  site_ip = "75.80.36.61" //Good for adding and updating site and in response object
  force_ssl = "true" //Good for adding site and in response object
  //naked_domain_san = "false" //TODO: Not available in TF site resource
  //wildcard_san = "false" //TODO: Not available in TF site resource
  //log_level = "full" //TODO: Fails with "reason": "Logs are not supported due to feature restrictions"
  //logs_account_id = "1034421" //TODO: Fails with "reason": "The provided logs_account_id invalid"

  ####################################################################
  # For updating the site ONLY!!!!
  # Available to modify after the site is created...
  #TODO: This entire site update part needs to have logic applied resource; one failure will return out of the resource_incapsula_site.go.resourceSiteUpdate function...
  ####################################################################
  # active = "bypass" //Good for updating site and in response object
  domain_validation = "dns" //TODO: This call fails after being set by Imperva and cannnot be changed. Default is DNS; call will error on HTML and pass on EMAIL but no change.
    //Error from Incapsula service when updating site for siteID 12608430:
    //{"res":1,"res_message":"Unexpected error","debug_info":{"id-info":"13007","domain_html":"An error occured while trying to modify site\u0027s domain validation"}}
  //approver = "joesph.moore@gmail.com" //TOD: This call fails when trying to apply on update.
    //Error from Incapsula service when updating site for siteID 12608430:
    //{"res":4202,"res_message":"Domain_email invalid","debug_info":{"approver":"approver email value is not a valid email value from the site\u0027s domain_emails","id-info":"13007"}}
  # ignore_ssl = "true" ////Good for updating site and in response object statusEnum key from pending-select-approver to pending_ssl_approval TODO: How to handle this in state
  # acceleration_level = "none" //Good for updating site and in response object
  # seal_location = "api.seal_location.bottom_right" //Good for updating site and in response object
  # domain_redirect_to_full = "true" //TODO: Good for updating site but not shown in reponse or in UI
  remove_ssl = "false" //Good for updating site and in response object
}

####################################################################
# Custom Certificates
####################################################################
variable "certificate" {
  default = "-----BEGIN CERTIFICATE-----\nMIIDgjCCAmoCCQCk3MsAS5x+UjANBgkqhkiG9w0BAQsFADCBgjELMAkGA1UEBhMC\nVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlTYW4gRGllZ28xCzAJBgNVBAoMAlNF\nMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFzaC5iZWVyLmNlbnRlcjEdMBsGCSqG\nSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wHhcNMTkwNzA4MTU0MjQ0WhcNMjAwNzA3\nMTU0MjQ0WjCBgjELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRIwEAYDVQQHDAlT\nYW4gRGllZ28xCzAJBgNVBAoMAlNFMQswCQYDVQQLDAJTRTEZMBcGA1UEAwwQZGFz\naC5iZWVyLmNlbnRlcjEdMBsGCSqGSIb3DQEJARYOYmFAaW1wZXJ2YS5jb20wggEi\nMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCj0rYKUhVtNKQ/oKZdCfxvLKhQ\nLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1DhPRSDGgONdHe\n2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkPQ8dsPpE945lW\nsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607DqdZUS/O4XSx+d0\n5kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1wqRdwGJsUdq6\nkC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2uWCBSY/OLAgMB\nAAEwDQYJKoZIhvcNAQELBQADggEBABfNZcItHdsSpfp8h+1EP5BnRuoKj+l42EI5\nE9dVlqdOZ25+V5Ee899sn2Nj8h+/zVU3+IDO2abUPrDd2xZHaHdf0p69htSwFTHs\nEwUdPUUsKRSys7fVP1clHcKWswTcoWIzQiPZsDMoOQw/pzN05cXSzdo8wSWuEeBK\ncqRNd5BKPeeXbFa4i5TFzT/+pl8V075k16tzHSbT7QDk5fuZWYv/2jImw/lgS/nx\nDWtlprrgG6AX1FzovDs/NnNq/e7vZtn8sdOoO2pCSVymNvctNLV2tFcS8sPQDl5M\nIpnZa3kktAegjsCln1JvD0AFigXrF8wjK+FKGI8SPJfbTQ149+A=\n-----END CERTIFICATE-----"
}
variable "private_key" {
  default = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCj0rYKUhVtNKQ/\noKZdCfxvLKhQLbCNsOt94afUZCbM93/TYj7kHQaapJ6s9snPjN6dvRKo/0h+qx1D\nhPRSDGgONdHe2plv6M7h2gNhBF2853/GZLdNzO9GBHDI6VB9bFJpQvqBl+Cy7nkP\nQ8dsPpE945lWsQ7KMakikp1oJrFHmfalNMo+VQgOKPNc3jUlgmSNEwk3Cf607Dqd\nZUS/O4XSx+d05kRg3hmrjDxDyTwG2gQDJBGkdZ87HUqd5NC7KlrY5xuLkloq4Rt1\nwqRdwGJsUdq6kC8lPmikw2i3peTUu03T3OiZxBpKK6gNMcKe3uA3zSPdoY/mDY2u\nWCBSY/OLAgMBAAECggEAfDPprkNzWTta95594vYKO+OYrEknnRUwRV0LF/ziae2P\nLR1EX0AeXKpIXwwwDzpXITJS7mu1c8uZwTIQ5g/f6D4nopULVYwlJZhbjXd49hpx\nhmGfk8227te5BqnVS3IPvRx5vjz+r8obYFZb4JZDGa/v9okAlI04FS0hR/Bl4ckD\naIsztf4R+AO2dP6BxYZGIwcq3jkbf0BdyQpkw4Ds7pdKbSa+PsobseyI2NqR2ryX\n4HH4b89HZj8lfiniIN3tPV6uIvpPS6jJklLKy6zdkIFOng/OGwxXomGkrk9ZjBHm\nJx5yA5YfwPidyt80wO9/26wClXYidfKQC8mDN21owQKBgQDPQbNr/sGiI2QzTOpb\nYTx0FWzWMnn9N2XiQm5rcr9kM5WsXh+anlqP54MeXDGZ2f6L8+aGrghZ/78WbG9J\nDbtEc7qTSRw5LFRglqn32a3ppHToEzOVxsA3g/OBJT5lJJwGMTdeKEXtLMmkm/sz\n1ClFnYJ1I8rNcueI9936odDWKwKBgQDKWgGwWTbqVa3wVIOFvluxolQzo6TEBFbf\nQTJo7byO2iRZvhrZUUk8539Uz2px0Ilzxx61CszhNWDVNwgqsN7FtuzXuCwz9GzU\nyBWkzPKGzvK12aFMYoj/cPbcRfMpYWNoK/YfEKfTRkJJfrJSbWP2XlyEr69te8s7\nB/zxOtUIIQKBgEjoJcOhtF/i70aUkgRfKjLzrnuS+hK3QCHdmJY3oVgQRWCDI77y\nYY0ptZgielhStRZqT/eklM+EBaZPsr4SFIQ56bISD9mU3IG1vkivzFvaPD2/M3BG\noCtnQWt2vII75J7RBVcb9609ChnbvPw4b+RLSi8GzjqDZytpdi7KaXpNAoGAS2Ym\nYvObRs4ONhMHvvojaJk4DtXXO0Lyq9W7VuXe8MvP57CyiG+FfrAz/gIbg7VUwlNb\n2dHgbbpaDpim7mFhYQK8VdVGg0V8l/zGM9Y6OIk8Xw5sz+2XZrdNBN77sFudkt9u\nojyujEcNxBz1jUk9iju29aoREBakr6ZWVfy6DIECgYEAtXxrOsbMsbHhVGqgeGXy\nhLXIltR+7NIUaxpLHhYCMzK9SbyZvx/Hd6m34oTw9ws+tHFpeCyiVU+wQgmx0ARD\ncDLKOPIHTGYhq/H8Oc6/Dzfxs1L/hH34mw5u7hVtAaA+q8iaRGVZ797dTVSxw4U0\nRm+BCDRhDcvaG7qpvFj8T6k=\n-----END PRIVATE KEY-----"
}

variable "passphrase" {
  default = "webco123"
}

resource "incapsula_custom_certificate" "custom-certificate" {
  site_id = "${incapsula_site.example-site.id}"
  certificate = "${var.certificate}"
  private_key = "${var.private_key}"
  passphrase = "${var.passphrase}"
}

####################################################################
# Security Rules
####################################################################

# Security Rule: Country
# resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule" {
#   site_id = "${incapsula_site.example-site.id}"
#   rule_id = "api.acl.blacklisted_countries"
#   countries = "AI,AN"
#   //continents = "SA"
#   depends_on = ["incapsula_site.example-site"]
# }
#
# # Security Rule: Country IP Exception
# //resource "incapsula_acl_security_rule" "example-global-blacklist-country-rule_exception" {
# //  rule_id = "api.acl.blacklisted_countries"
# //  site_id = "${incapsula_site.example-site.id}"
# //  ips = "192.168.1.1,192.168.1.2"
# //  urls = "/myurl,/myurl2"
# //  url_patterns = "EQUALS,CONTAINS"
# //  countries = "JM,US"
# //  client_apps= "488,123"
# //  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-country-rule"]
# //}
#
# # Security Rule: Blacklist IP
# resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule" {
#   site_id = "${incapsula_site.example-site.id}"
#   rule_id = "api.acl.blacklisted_ips"
#   ips = "192.168.1.0/24"
#   depends_on = ["incapsula_site.example-site"]
# }
#
# # Security Rule: Blacklist IP Exception
# //resource "incapsula_acl_security_rule" "example-global-blacklist-ip-rule_exception" {
# //  rule_id = "api.acl.blacklisted_ips"
# //  site_id = "${incapsula_site.example-site.id}"
# //  ips = "192.168.1.1,192.168.1.2"
# //  urls = "/myurl,/myurl2"
# //  url_patterns = "EQUALS,CONTAINS"
# //  countries = "JM,US"
# //  client_apps= "488,123"
# //  depends_on = ["incapsula_site.example-site", "incapsula_acl_security_rule.example-global-blacklist-ip-rule"]
# //}
#
# # Security Rule: URL
# resource "incapsula_acl_security_rule" "example-global-blacklist-url-rule" {
#   rule_id = "api.acl.blacklisted_urls"
#   site_id = "${incapsula_site.example-site.id}"
#   url_patterns = "CONTAINS,EQUALS"
#   urls = "/alpha,/bravo"
#   depends_on = ["incapsula_site.example-site"]
# }
#
# # Security Rule: Whitelist IP
# resource "incapsula_acl_security_rule" "example-global-whitelist-ip-rule" {
#   rule_id = "api.acl.whitelisted_ips"
#   site_id = "${incapsula_site.example-site.id}"
#   ips = "192.168.1.3,192.168.1.4"
#   depends_on = ["incapsula_site.example-site"]
# }
#
# ####################################################################
# # Incap Rules
# ####################################################################
#
# # Incap Rule: Alert
# resource "incapsula_incap_rule" "example-incap-rule-alert" {
#   priority = "1"
#   name = "Example incap rule alert"
#   site_id = "${incapsula_site.example-site.id}"
#   action = "RULE_ACTION_ALERT"
#   filter = "Full-URL == \"/someurl\""
#   depends_on = ["incapsula_site.example-site"]
# }

//resource "aws_route53_record" "dash_record" {
//  depends_on = ["incapsula_site.dash_api"]
//  name = "dash.${data.aws_route53_zone.zone.name}"
//  type = "CNAME"
//  zone_id = "${data.aws_route53_zone.zone.zone_id}"
//  ttl = "180"
//  records = ["${incapsula_site.dash_api.dns_cname_record_value}"]
//}
//
//output "incap_ald_url" {
//  value = "${aws_route53_record.dash_record.fqdn}"
//}

output "incap_siteID" {
  value = "${incapsula_site.example-site.id}"
}
