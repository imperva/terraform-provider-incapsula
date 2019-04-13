Security Exception - api.threats.cross_site_scripting

Create
/api/prov/v1/sites/configure/whitelists
site_id=int *
rule_id=api.threats.cross_site_scripting *
action=api.threats.action.block_request *
action_text=Block Request *ips=1.2.3.6 (optional)
urls=/myurl,/myurl2 (optional)
url_patterns=EQUALS,CONTAINS (optional)
countries=JM,US (optional)
client_apps=488,123 (optional)

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/whitelists
site_id=int *
rule_id=api.threats.cross_site_scripting *
action=api.threats.action.block_request *
whitelist_id=int *
action_text=Block Request *ips=1.2.3.6 (optional)
urls=/myurl,/myurl2 (optional)
url_patterns=EQUALS,CONTAINS (optional)
countries=JM,US (optional)
client_apps=488,123 (optional)

Delete
/api/prov/v1/sites/configure/whitelists
domain=string *
site_id=int *
rule_id=api.threats.cross_site_scripting *
whitelist_id=int *
delete_whitelist=true *
