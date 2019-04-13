ACL Exception - api.acl.blacklisted_countries

Create
/api/prov/v1/sites/configure/whitelists
site_id=int *
rule_id=api.acl.blacklisted_countries *
ips=1.2.3.6 (comma separated list)
urls=/myurl,/myurl2 (comma separated list)
url_patterns=EQUALS,CONTAINS (comma separated list)
client_app_types=DataScraper (optional)

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/whitelists
site_id=int *
rule_id=api.acl.blacklisted_countries *
whitelist_id=int *
ips=1.2.3.6 (comma separated list)
urls=/myurl,/myurl2 (comma separated list)
url_patterns=EQUALS,CONTAINS (comma separated list)
client_app_types=DataScraper (optional)

Delete
/api/prov/v1/sites/configure/whitelists	"domain=string *
site_id=int *
rule_id=api.acl.blacklisted_countries *
whitelist_id=int *
delete_whitelist=true *
