Security - api.threats.illegal_resource_access

Create
/api/prov/v1/sites/configure/security
site_id=int *
rule_id=api.threats.illegal_resource_access *
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
rule_id=api.threats.illegal_resource_access *
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
rule_id=api.threats.illegal_resource_access *
action_text=Block Request *
security_rule_action=api.threats.action.block_request * (reset to default)
