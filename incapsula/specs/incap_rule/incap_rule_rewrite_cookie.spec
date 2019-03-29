Rule - Rewrite Cookie (ADR)

Create
/api/prov/v1/sites/incapRules/add
site_id=int *
enabled=true *
priority=18 *
name=Sample Rule Rewrite Cookie ADR *
action=RULE_ACTION_REWRITE_COOKIE *
add_missing=true | false
from=some_optional_value
to=some_new_value *
allow_caching=false | true
filter= Full-URL == "/someurl"
rewrite_name=my_cookie_name *

Read
/api/prov/v1/sites/incapRules/list
include_ad_rules=Yes *
include_incap_rules=No *

Update
/api/prov/v1/sites/incapRules/edit
site_id=int *
enabled=true *
priority=18 *
name=Sample Rule Rewrite Cookie ADR *
action=RULE_ACTION_REWRITE_COOKIE *
add_missing=true | false
from=some_optional_value
to=some_new_value *
allow_caching=false | true
filter= Full-URL == "/someurl"
rewrite_name=my_cookie_name *
rule_id=int *

Delete
/api/prov/v1/sites/incapRules/delete
rule_id=int *
