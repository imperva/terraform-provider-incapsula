
resource "incapsula_abp_websites" "abp_websites" {
    account_id = data.incapsula_account_data.account_data.current_account
    auto_publish = true
    website_group {
        name = "sites-1"
        website {
            site_id = incapsula_site.sites-1.id
            enable_mitigation = false
        }
    }
    website_group {
        name = "sites-2"
        website {
            site_id = incapsula_site.sites-2.id
            enable_mitigation = true
        }
    }
}


resource "incapsula_abp_websites" "abp_websites" {
    account_id = data.incapsula_account_data.account_data.current_account
    auto_publish = true
    website_group {
        name = "sites"
        website {
            site_id = incapsula_site.sites-1.id
            enable_mitigation = false
        }
        website {
            site_id = incapsula_site.sites-2.id
            enable_mitigation = true
        }
    }
    website_group {
        name = "sites" # Duplicate name
        name_id = "sites-2" # name_id can be used to disambiguate names in case of duplicates
        website {
            site_id = incapsula_site.sites-3.id
            enable_mitigation = true
        }
    }
}
