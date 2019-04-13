ACL - api.acl.blacklisted_urls

Create
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_urls *
urls=/admin,wp-admin (comma separated list)
url_patterns=PREFIX,CONTAINS (comma separated list)

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_urls *
urls=/admin,wp-admin (comma separated list)
url_patterns=PREFIX,CONTAINS (comma separated list)

Delete
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_urls *
urls="" (empty string to clear rule)
url_patterns="" (empty string to clear rule)
