ACL - api.acl.whitelisted_ips (no exceptions for this ACL)

Create
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.whitelisted_ips *
ips=192.168.1.1,1.2.3.5 (comma separated list) *

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.whitelisted_ips *
ips=192.168.1.1,1.2.3.5 (comma separated list) *

Delete
/api/prov/v1/sites/configure/acl
site_id=int *
rule_id=api.acl.whitelisted_ips *
ips="""" (empty string to clear rule) *
