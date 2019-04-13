Data Center Servers

Create
/api/prov/v1/sites/dataCenters/servers/add
dc_id: 	int
server_address: IP, CNAME
is_standby: yes (optional)

Read
/api/prov/v1/sites/dataCenters/list
site_id: int *

Update
/api/prov/v1/sites/dataCenters/servers/edit
server_id: int *
server_address: IP or CNAME *
is_standby: yes | ""
is_content: yes | ""

Delete
/api/prov/v1/sites/dataCenters/servers/delete
server_id: int *
