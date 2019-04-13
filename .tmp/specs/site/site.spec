Site

Create
/api/prov/v1/sites/add
account_id: int *
domain: string

Read
/api/prov/v1/sites/status
site_id: int *

Update
/api/prov/v1/sites/configure
site_id: int *
param/value * (submit single param and value)
-------------
param=active
value=active | bypass
-------------
param=site_ip
value=1.2.3.4,1.2.3.5,some.cname.com (comma separated list)
-------------
param=domain_validation
value=email | html | dns
-------------
param=approver
value=my.approver@email.com (some approver email address)
-------------
param=ignore_ssl
value=true | ""
-------------
param=acceleration_level
value=none | standard | aggressive
-------------
param=seal_location
value=api.seal_location.bottom_left | api.seal_location.none | api.seal_location.right_bottom | api.seal_location.right | api.seal_location.left | api.seal_location.bottom_right | api.seal_location.bottom
-------------
param=domain_redirect_to_full
value= true | ""
-------------
param=remove_ssl
value=true | ""
-------------
param=ref_id
value=string (free-text field to add a unique identifier)

Delete
/api/prov/v1/sites/delete
site_id: int *
