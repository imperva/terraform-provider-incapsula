# Provider information, like api_id and api_key can be specified here
# or as environment variables: INCAPSULA_API_ID and INCAPSULA_API_KEY
provider "incapsula" {
  api_id = "foo"
  api_key = "bar"
}

# Site information
resource "incapsula_site" "example-terraform-site" {
  domain = "examplesite.com"
}