./kubectl create -f /home/go/src/mini-kubernetes/presentation/autoscaler/autoscaler_pre.yaml
#docker ps -q | wc -l统计正在运行容器数量
./kubectl get autoscaler
./kubectl describe autoscaler test-autoscaler
#在etcd中查看资源:etcdctl get --prefix "resource_usage/podInstance/test-auto"
#删除一个pod,然后docker ps -q | wc -l看有没有恢复过来
#删除一个autoscaler
./kubectl delete autoscaler test-autoscaler
#容器和autoscaler都删除了docker ps -q | wc -l
./kubectl get autoscaler

