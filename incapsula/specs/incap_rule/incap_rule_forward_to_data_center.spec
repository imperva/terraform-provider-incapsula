Rule - Forward to Data Center (ADR)

Create
/api/prov/v1/sites/incapRules/add
site_id=int *
enabled=true *
priority=1 *
name=Sample Rule Forward to Data Center ADR *
dc_id=int *
action=RULE_ACTION_FORWARD_TO_DC *
allow_caching=false | true
filter= Full-URL == ""/someurl"""

Read
/api/prov/v1/sites/incapRules/list
include_ad_rules=Yes *
include_incap_rules=No *"

Update
/api/prov/v1/sites/incapRules/edit
site_id=int *
enabled=true *
priority=1 *
name=Sample Rule Forward to Data Center ADR *
dc_id=int *
action=RULE_ACTION_FORWARD_TO_DC *
allow_caching=false }| true
filter= Full-URL == ""/someurl""
rule_id=int *"

Delete
/api/prov/v1/sites/incapRules/delete
rule_id=int *
