curl -X GET http://:33333@localhost:33333/v2/catalog


#########################Mysql######################

# provision
curl -i -X PUT http://user:pass@localhost:33333/v2/service_instances/mysql-shared-instance-abc123 -d '{
  "service_id":"f614fcc2-3cb2-4400-aa93-87714417f2cf",
  "plan_id":"726aa087-5a18-4a04-b449-5d455cecb29b",
  "organization_guid": "default",
  "space_guid":"space-guid",
  "accepts_incomplete":true,
  "parameters": {"ami_id":"ami-ecb68a84"}
}' -H "Content-Type: application/json"

sleep 30

curl -i -X GET 'http://user:pass@localhost:33333/v2/service_instances/mysql-shared-instance-abc123/last_operation' 

sleep 20

# bind
curl -i -X PUT http://user:pass@localhost:33333/v2/service_instances/mysql-shared-instance-abc123/service_bindings/d9d0c79e-ffbc-457e-93b6-14377a75e139 -d '{
  "plan_id":        "726aa087-5a18-4a04-b449-5d455cecb29b",
  "service_id":     "f614fcc2-3cb2-4400-aa93-87714417f2cf",
  "app_guid":       "app-guid"
}' -H "Content-Type: application/json"

sleep 3

# unbind
curl -i -X DELETE -L 'http://user:pass@localhost:33333/v2/service_instances/mysql-shared-instance-abc123/service_bindings/d9d0c79e-ffbc-457e-93b6-14377a75e139?service_id=f614fcc2-3cb2-4400-aa93-87714417f2cf&plan_id=726aa087-5a18-4a04-b449-5d455cecb29b' 

sleep 3

# deprovision
curl -i -X DELETE -L 'http://user:pass@localhost:33333/v2/service_instances/mysql-shared-instance-abc123?service_id=f614fcc2-3cb2-4400-aa93-87714417f2cf&plan_id=726aa087-5a18-4a04-b449-5d455cecb29b'


