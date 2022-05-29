./kubectl create -f /home/go/src/mini-kubernetes/presentation/service/service_pre.yaml
#查看所有规则iptables sudo iptables -vL -t nat
./kubectl get services
./kubectl describe service test-service-for-deployment
#⽤⼾能够通过虚拟ip访问Service，由minik8s将请求具体转发⾄对应的pod。 curl 192.168.40.30 80
./kubectl create -f /home/go/src/mini-kubernetes/presentation/service/service_pre2.yaml
./kubectl get services
# Service内的pod也可以通过虚拟ip访问其他Service：进入上一个service的pod，去curl 192.168.40.40 80
./kubectl delete service test-service-for-deployment
./kubectl delete service test-service-for-deployment2
#展示service删干净了
./kubectl get services
# iptables干净：sudo iptables -vL -t nat
