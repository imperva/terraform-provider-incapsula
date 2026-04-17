resource "incapsula_site" "example-site" {
  domain = "examplesite.incaptest.co"
}

# Manual test for issue #626: verify rule deletion works (disable before delete)
# To test the fix:
# 1. terraform apply - creates the rule
# 2. terraform destroy - should disable rule first, then delete (fix in resourceIncapRuleDelete)
resource "incapsula_incap_rule" "example_rule" {
  site_id = incapsula_site.example-site.id
  name    = "Example Alert Rule - Test Disable Before Delete"
  action  = "RULE_ACTION_ALERT"
  filter  = "Full-URL == \"/admin\""
  enabled = true
}

# Additional test: verify disable works by updating the rule
# Uncomment to test disabling without deletion
# resource "incapsula_incap_rule" "test_disable_only" {
#   site_id = incapsula_site.example-site.id
#   name    = "Test Disable Only"
#   action  = "RULE_ACTION_ALERT"
#   filter  = "Full-URL == \"/test\""
#   enabled = true
# }
# Then: terraform apply -var="enable_disable_test=false" to set enabled=false and verify UpdateIncapRule works