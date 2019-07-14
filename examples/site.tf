
provider "incapsula" {
  api_id = "your_api_id"
  api_key = "your_api_key"
}

resource "incapsula_site" "example-site" {
    domain = "your.fulldomain.com"
    account_id = "1014181"
    ref_id = "12345"
    send_site_setup_emails = "true"
    site_ip = "1.2.3.4"
    force_ssl = "true"
}

resource "incapsula_custom_certificate" "custom-certificate" {
    site_id = "${incapsula_site.example-site.id}"
    certificate = "${file("path/to/your/cert.crt")}"
    private_key = "${file("path/to/your/private_key.key")}"
    passphrase = "yourpassphrase"
}

