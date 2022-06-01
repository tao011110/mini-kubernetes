./kubectl create -f /home/go/src/mini-kubernetes/presentation/pod/pod_pre.yaml
./kubectl get pods
./kubectl describe pod pod_pre1
#然后进入ghost容器，尝试curl localhost 80
#apt-get update
#apt-get install curl
#curl localhost 80
#在挂载的/vpath目录下创建文件
#进入nginx容器，进入/vpath目录，发现确实创建好了;再去主机/home/mount_test目录下，发现宿主机中也创建好了
#展示一下yaml文件中的with和notWith
./kubectl create -f /home/go/src/mini-kubernetes/presentation/pod/pod_pre2.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/pod/pod_pre3.yaml
./kubectl get pods
./kubectl delete pod pod_pre1
./kubectl get pods
#docker ps -a + ./kubectl get pods 发现确实删除了
