go build ../tools/kubelet/kubelet.go
cp ../tools/kubelet/cadvisor .
go build ../tools/kubeproxy/proxy.go
./kubelet &
./cadvisor &
./proxy &
