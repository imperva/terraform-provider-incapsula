ACL - api.acl.blacklisted_ips

Create
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_ips' *
ips=192.168.1.1,1.2.3.5 (comma separated list) *

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_ips' *
ips=192.168.1.1,1.2.3.5 (comma separated list) *

Delete
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.blacklisted_ips' *
ips="" (empty string to clear rule) *
