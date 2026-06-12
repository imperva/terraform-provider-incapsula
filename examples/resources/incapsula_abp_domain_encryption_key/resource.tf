resource "incapsula_abp_domain_encryption_key" "test_com" {
  domain_id = incapsula_abp_domain.test_com.id
  key       = "U2VjcmV0IGtleSB1c2luZyBzdGF0ZS1vZi10aGUtYXJ0IGJhc2U2NCBlbmNyeXB0aW9u"
}
