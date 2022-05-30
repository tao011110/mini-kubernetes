./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/serviceForDNS.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/serviceForDNS2.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/podForDNS.yaml
#查看iptables 所有规则: sudo iptables -vL -t nat
./kubectl create -f /home/go/src/mini-kubernetes/presentation/dns/dns_pre.yaml
./kubectl get dns
./kubectl describe dns dns-test2
#在宿主机尝试curl www.minik8s.com:80/route2, curl www.minik8s.com:80/route3
#然后进入nginx容器，尝试curl www.minik8s.com:80/route2, curl www.minik8s.com:80/route3
#apt-get update
#apt-get install curl

