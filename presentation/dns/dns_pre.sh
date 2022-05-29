./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/podForDNS.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/serviceForDNS.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/serviceForDNS2.yaml
#查看iptables 所有规则: sudo iptables -vL -t nat
./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/dns_pre.yaml
./kubectl get dns
./kubectl describe dns dns-test
#然后进入ghost容器，尝试curl localhost 80
#apt-get update
#apt-get install curl
#curl localhost 80

