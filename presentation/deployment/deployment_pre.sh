./kubectl create -f /home/go/src/mini-kubernetes/presentation/deployment/deployment_pre.yaml
#docker ps -q | wc -l统计正在运行容器数量
./kubectl get deployment
./kubectl describe deployment test-deployment
./kubectl create -f /home/go/src/mini-kubernetes/presentation/deployment/deployment_pre2.yaml
./kubectl get deployment
#删除一个pod,然后docker ps -q | wc -l看有没有恢复过来
#删除一个deployment
./kubectl delete deployment test-deployment
./kubectl delete deployment test-deployment2
#容器和deployment都删除了docker ps -q | wc -l
./kubectl get deployment

