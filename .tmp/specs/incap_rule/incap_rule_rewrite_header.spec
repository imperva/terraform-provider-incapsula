Rule - Rewrite Header (ADR)

Create
/api/prov/v1/sites/incapRules/add
site_id=int *
enabled=true *
priority=17 *
name=Sample Rule Rewrite Header ADR *
action=RULE_ACTION_REWRITE_HEADER *
add_missing=true | false
from=some_optional_value
to=some_new_value *
allow_caching=false | true
filter= Full-URL == "/testurl"
rewrite_name=mytestheader

Read
/api/prov/v1/sites/incapRules/list
include_ad_rules=Yes *
include_incap_rules=No *"

Update
/api/prov/v1/sites/incapRules/edit
site_id=int *
enabled=true *
priority=17 *
name=Sample Rule Rewrite Header ADR *
action=RULE_ACTION_REWRITE_HEADER *
add_missing=true | false
from=some_optional_value
to=some_new_value *
allow_caching=false | true
filter= Full-URL == "/testurl"
rewrite_name=mytestheader
rule_id=int *

Delete
/api/prov/v1/sites/incapRules/delete
rule_id=int *
