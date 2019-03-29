ACL - api.acl.blacklisted_countries

Create
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_countries' *
continents=AF,AS (comma separated list)
countries=AI,AN (comma separated list)"

Read
/api/prov/v1/sites/status	site_id: int *

Update
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id = ""api.acl.blacklisted_countries"" *
countries=""AI,AN"" (comma seperated string)  *"

Delete
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id = "api.acl.blacklisted_countries" *
countries= "" (empty string to clear rule) *
