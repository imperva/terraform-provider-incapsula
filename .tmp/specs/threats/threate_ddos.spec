Security - api.threats.ddos

Create
/api/prov/v1/sites/configure/security
site_id=int *
activation_mode_text=Auto *
ddos_traffic_threshold=1000 *
rule_id=api.threats.ddos *
activation_mode=api.threats.ddos.activation_mode.auto * (default)
  ( ddos.activation_mode.auto |
    ddos.activation_mode.off |
    ddos.activation_mode.on )

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure/security
site_id=int *
activation_mode_text=Auto *
ddos_traffic_threshold=1000 *
rule_id=api.threats.ddos *
activation_mode=api.threats.ddos.activation_mode.auto * (default)
  ( ddos.activation_mode.auto |
    ddos.activation_mode.off |
    ddos.activation_mode.on )

Delete
/api/prov/v1/sites/configure/security
site_id=int *
activation_mode_text=Auto *
ddos_traffic_threshold=1000 *
rule_id=api.threats.ddos *
activation_mode=api.threats.ddos.activation_mode.auto * (reset to default)
