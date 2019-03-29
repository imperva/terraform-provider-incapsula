Rule - Block Ip (IncapRule)

Create
/api/prov/v1/sites/incapRules/add
site_id=int *
enabled=true *
priority=1 *
name=Sample Rule Block Ip IncapRule *
action=RULE_ACTION_BLOCK_IP *
filter= Full-URL == ""/someurl"""

Read
/api/prov/v1/sites/incapRules/list
include_ad_rules=No *
include_incap_rules=Yes *"

Update
/api/prov/v1/sites/incapRules/edit
site_id=int *
enabled=true *
priority=1 *
name=Sample Rule Block Ip IncapRule *
action=RULE_ACTION_BLOCK_IP *
filter= Full-URL == ""/someurl""
rule_id=int *"

Delete
/api/prov/v1/sites/incapRules/delete
rule_id=int *
