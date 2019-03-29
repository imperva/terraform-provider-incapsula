Data Centers

Create
/api/prov/v1/sites/dataCenters/add
site_id: int *
name: string *
server_address: IP or CNAME *
is_standby: yes | ""
is_content: yes | ""

Read
/api/prov/v1/sites/dataCenters/list
site_id: int

Update
/api/prov/v1/sites/dataCenters/edit
dc_id: int *
name: string *
is_standby: yes | ""
is_content: yes | ""

Delete
/api/prov/v1/sites/dataCenters/delete
dc_id: int *
