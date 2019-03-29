ACL Exception - api.acl.blacklisted_ips

Create
/api/prov/v1/sites/configure/whitelists
site_id=int *
rule_id=api.acl.blacklisted_ips *
ips=1.2.3.6 (comma separated list)
urls=/myurl,/myurl2 (comma separated list)
url_patterns=EQUALS,CONTAINS (comma separated list)
countries=JM,US (comma separated list)
client_apps=488,123 (comma separated list)

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/whitelists
site_id=int *
rule_id=api.acl.blacklisted_ips *
whitelist_id=int *
ips=1.2.3.6 (comma separated list)
urls=/myurl,/myurl2 (comma separated list)
url_patterns=EQUALS,CONTAINS (comma separated list)
countries=JM,US (comma separated list)
client_apps=488,123 (comma separated list)

Delete
/api/prov/v1/sites/configure/whitelists
domain=string *
site_id=int *
rule_id=api.acl.blacklisted_ips *
whitelist_id=int *
delete_whitelist=true *
