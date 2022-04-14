# k8s中涉及的网络

## k8s中有3种不同类型的ip
 - nodeIP: 各node的物理网卡ip(对外可见)
 - podIP: 为各Pod分配的虚拟ip(对外不可见)
 - clusterIP(serviceIP): 各service (对外可见)

## 集群
