Provides an Incapsula ABP (Advanced Bot Protection) websites resource. Allows for ABP to enabled and configured for given websites.


```hcl
resource "incapsula_abp_websites" "abp_websites" {
    account_id = data.incapsula_account_data.account_data.current_account
    auto_publish = true
    website_group {
        name = "sites-1"
        website {
            website_id = incapsula_site.sites-1.id
            enable_mitigation = false
        }
    }
    website_group {
        name = "sites-2"
        website {
            website_id = incapsula_site.sites-2.id
            enable_mitigation = true
        }
    }
}
```
