Security - api.threats.bot_access_control

Create
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.bot_access_control *
block_bad_bots=true | false (optional, default: true)
challenge_suspected_bots=true | false (optional, default: false)

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.bot_access_control *
block_bad_bots=true | false (optional, default: true)
challenge_suspected_bots=true | false (optional, default: false)

Delete
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.bot_access_control *
block_bad_bots=true (reset to defauilt) *
challenge_suspected_bots=false (reset to defauilt) *
