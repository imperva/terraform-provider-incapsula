Security - api.threats.backdoor

Create
/api/prov/v1/sites/configure/security
site_id=int *
action_text=Auto-Quarantine *
rule_id=api.threats.backdoor *
security_rule_action=api.threats.action.quarantine_url (default) *
  ( api.threats.action.alert |
    api.threats.action.disabled |
    api.threats.action.quarantine_url )

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/security
site_id=int *
action_text=Auto-Quarantine *
rule_id=api.threats.backdoor *
security_rule_action=api.threats.action.quarantine_url (default) *
  ( api.threats.action.alert |
    api.threats.action.disabled |
    api.threats.action.quarantine_url )

Delete
/api/prov/v1/sites/configure/security
site_id=int *
action_text=Auto-Quarantine *
rule_id=api.threats.backdoor *
security_rule_action=api.threats.action.quarantine_url (reset to default) *
