./kubectl create -f /home/go/src/mini-kubernetes/presentation/autoscaler/new_version/autoscaler_mem.yaml
#docker ps -q | wc -l统计正在运行容器数量
./kubectl get autoscaler
./kubectl get pods
./kubectl describe autoscaler autoscaler_mem
./kubectl create -f /home/go/src/mini-kubernetes/presentation/autoscaler/new_version/service_autoscaler.yaml
#iptables查看所有规则 sudo iptables -vL -t nat
#在主机curl 192.168.40.90:80
#进入一个容器，boom，再docker ps -q | wc -l + ./kubectl get pods

#删除一个pod,然后docker ps -q | wc -l看有没有恢复过来
#删除一个autoscaler
./kubectl delete autoscaler autoscaler_mem
#容器和autoscaler都删除了docker ps -q | wc -l
./kubectl get autoscaler
./kubectl create -f /home/go/src/mini-kubernetes/presentation/autoscaler/new_version/autoscaler_cpu.yaml
./kubectl describe autoscaler autoscaler_cpu
./kubectl delete autoscaler autoscaler_cpu
./kubectl get autoscaler
