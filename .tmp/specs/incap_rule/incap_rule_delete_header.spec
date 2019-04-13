Rule - Delete Header (ADR)

Create
/api/prov/v1/sites/incapRules/add
enabled=true | false *
name=string (no special characters) *
action=RULE_ACTION_DELETE_HEADER *
filter=string *
rule_id=int *

Read
/api/prov/v1/sites/incapRules/list
include_ad_rules=Yes *
include_incap_rules=No *

Update
/api/prov/v1/sites/incapRules/edit

Delete
/api/prov/v1/sites/incapRules/delete
rule_id=int *
