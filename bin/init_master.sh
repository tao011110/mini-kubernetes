go build ../tools/master/master.go
go build ../tools/controller/controller.go
go build ../tools/scheduler/scheduler.go
go build ../tools/activer/activer.go
go build ../tools/kubectl/kubectl.go
./master > ../log/master_log &
./controller > ../log/controller_log &
./scheduler > ../log/scheduler_log &
./activer > ../log/activer_log &
