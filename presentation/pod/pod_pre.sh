./kubectl create -f /home/go/src/mini-kubernetes/presentation/pod/pod_pre.yaml
./kubectl get pods
./kubectl describe pod pod1
#然后进入ghost容器，尝试curl localhost 80
#apt-get update
#apt-get install curl
#curl localhost 80
./kubectl delete pod pod1
./kubectl get pods
#docker ps -a 发现确实删除了
