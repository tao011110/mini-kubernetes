# Kubernetes搭建指南
## 安装K8S
 - 在每台机器上执行(需要sudo权限)
````sh
# install dependency
apt-get update
apt-get install -y apt-transport-https ca-certificates curl

# add key
curl https://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg | apt-key add -

# add apt source
cat <<EOF >/etc/apt/sources.list.d/kubernetes.list
deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main
EOF

# install kube* and enable kubelet
apt-get update
apt-get install -y kubelet kubeadm kubectl
systemctl enable kubelet
````
## 创建K8S集群
在主节点(val00)上执行
````sh
# apiserver-advertise-address: val00 address，pod-network-cidr: default CIDR in the config of flannel
sudo kubeadm init --pod-network-cidr=10.244.0.0/16 --apiserver-advertise-address=192.168.15.100 --image-repository registry.aliyuncs.com/google_containers --kubernetes-version 1.20.5

# copy config file
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# install flannel in the k8s cluster
sudo kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
````
 - 此时集群已经创建好，kubectl get pod --all-namespaces可以查看组件是否运行成功
 - 在需要加入集群的机器上执行
````sh
# token和sha256会在主节点创建时打印到命令行
sudo kubeadm join 192.168.15.100:6443 --token eyxlsr.2fgt96lmo03hx6nb \
    --discovery-token-ca-cert-hash sha256:c178213ffb5807c73500050a02d57401889710de55a06369c32e02a7e33830d1
````
## 清理集群
 - 在每台子节点执行
````sh
sudo kubeadm reset
````
 - 在主节点(val00)上执行
````
# for each <node name> to be deleted
kubectl delete node <node name>

sudo kubeadm reset
rm $HOME/.kube/config
````
## 安装Helm
 - helm是K8S生态中的一个包管理工具，类似于python中的pip，可以方便的将应用打包成helm charts进行安装
````sh
curl https://baltocdn.com/helm/signing.asc | sudo apt-key add -
echo "deb https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list
sudo apt-get update
sudo apt-get install helm
````
## 安装Metrics-Server
 - 默认K8S是没有安装Metrics-Server的，安装Metrics-Server后可以使用kubectl top命令查看pod资源使用
 - K8S生态中还有更加高级复杂的集群监控组件(Prometheus)
````sh
wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# install metrics-server
kubectl apply -f components.yaml
````
 - 发现image pull失败了，网上查了一下发现都遇到过这个问题
````
# 手动拉取镜像并重新打tag (注意，任务被分配到哪个机器就去哪个机器pull)
sudo docker search metrics-server:v0.4.1
sudo docker pull phperall/metrics-server:v0.4.1
sudo docker tag phperall/metrics-server:v0.4.1 k8s.gcr.io/metrics-server/metrics-server:v0.4.2

# 此时依然不可以运行，describe pod 发现liveness报错
Warning  Unhealthy  19s (x5 over 99s)   kubelet            Liveness probe failed: HTTP probe failed with statuscode: 500
Warning  Unhealthy  16s (x6 over 116s)  kubelet            Readiness probe failed: HTTP probe failed with statuscode: 500

# 解决方案：添加--kubelet-insecure-tls参数到yaml文件里的container

# 重新部署,运行kubectl top nodes发现metrics-server可以正常工作
````
## val集群上使用K8S
 - 现在每台机器上都已经安装了K8S，使用K8S环境需要在主节点上创建集群，并将子节点加入到集群中