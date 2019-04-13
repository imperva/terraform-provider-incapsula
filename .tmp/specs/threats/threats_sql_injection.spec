Security - api.threats.sql_injection

Create
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.sql_injection *
action_text=Block Request *
security_rule_action=api.threats.action.block_request * (default)
  ( api.threats.action.disabled |
    api.threats.action.alert |
    api.threats.action.block_request |
    api.threats.action.block_user |
    api.threats.action.block_ip )

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.sql_injection *
action_text=Block Request *
security_rule_action=api.threats.action.block_request * (default)
  ( api.threats.action.disabled |
    api.threats.action.alert |
    api.threats.action.block_request |
    api.threats.action.block_user |
    api.threats.action.block_ip )

Delete
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.sql_injection *
action_text=Block Request *
security_rule_action=api.threats.action.block_request * (reset to default)
