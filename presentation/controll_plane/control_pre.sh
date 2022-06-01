./kubectl create -f /home/go/src/mini-kubernetes/presentation/controll_plane/control_pod.yaml
./kubectl get pods
./kubectl create -f /home/go/src/mini-kubernetes/presentation/controll_plane/control_service.yaml
./kubectl get services
#sudo iptables -vL -t nat 查看规则
#关掉控制面后，docker ps -a发现pod仍然存在
#重启控制面
./kubectl get pods
./kubectl get services
# curl 192.168.40.10:80
