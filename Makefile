run-cluster:
	go run ./cmd/cluster -addr 0.0.0.0:8081 

run-switch:
	go run ./cmd/switch -addr 0.0.0.0:8080 -cluster http://192.168.1.104:8081


run-worker1:
	go run ./cmd/worker -id worker-1 -addr 0.0.0.0:8091 -cluster http://192.168.1.104:8081

run-worker2:
	go run ./cmd/worker -id worker-2 -addr 0.0.0.0:8092 -cluster http://192.168.1.104:8081

root:root@tcp(127.0.0.1:3306)/sirocco


go go go
go run ./cmd/worker \
  -id=worker-2 \
  -addr=:8092 \
  -cluster=http://localhost:8081 \
  -mysql="sirocco:sirocco@tcp(127.0.0.1:3306)/sirocco"
