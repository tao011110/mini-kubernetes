go build ../tools/master/master.go
go build ../tools/controller/controller.go
go build ../tools/scheduler/scheduler.go
go build ../tools/activer/activer.go
go build ../tools/kubectl/kubectl.go
./master & \
./controller & \
./scheduler & \
./activer &
