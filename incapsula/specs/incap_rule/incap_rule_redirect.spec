Rule - Redirect (ADR)

Create
/api/prov/v1/sites/incapRules/add
site_id=int *
enabled=true *
priority=1
response_code=302 *
name=Sample Rule Redirect ADR *
action=RULE_ACTION_REDIRECT *
from=https://site1.com/url1
to=https://site2.com/url2 *
filter= AnyHeaderValue == ""testval""

Read
/api/prov/v1/sites/incapRules/list
include_ad_rules=Yes *
include_incap_rules=No *

Update
/api/prov/v1/sites/incapRules/edit
site_id=int *
enabled=true *
priority=1
response_code=302 *
name=Sample Rule Redirect ADR
action=RULE_ACTION_REDIRECT
from=https://site1.com/url1
to=https://site2.com/url2 *
filter= AnyHeaderValue == ""testval""
rule_id=int *

Delete
/api/prov/v1/sites/incapRules/delete
rule_id=int *
