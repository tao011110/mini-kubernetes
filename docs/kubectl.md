# kubectl 说明文档
By lsh

4种基本指令：create、describe、get、delete
create：创建资源
describe：查看资源信息
get：获取该资源全部信息
delete：删除资源

![](https://notes.sjtu.edu.cn/uploads/upload_84a05a5914a0eda9097f4bd633b22ef5.png)

## Pod
- kubectl create -f xxx.yaml
  根据配置文件创建pod资源，xxx.yaml为yaml文件位置，如果尝试创建的pod和已有pod重名会报错。
- kubectl describe pod xxx.yaml/yyy
  查看pod信息，可以通过pod的名称或pod配置文件来确定pod
- kubectl get pods
  查看所有pod的简要信息
  ![](https://notes.sjtu.edu.cn/uploads/upload_f380e69d133df6b2ba615c4372c6b288.png)

- kubectl delete pod yyy
  根据pod名称删除对应pod

## Service（包括ClusterIP和Nodeport）
- kubectl create -f xxx.yaml
  根据配置文件创建service资源，xxx.yaml为yaml文件位置。
- kubectl describe service xxx.yaml/yyy
  查看service信息，可以通过service的名称或service配置文件来确定service
- kubectl get services
  查看所有service的简要信息
  ![](https://notes.sjtu.edu.cn/uploads/upload_373575493162940ff3bc238f08a6c1d0.png)

- kubectl delete service yyy
  根据service名称删除对应service

## DNS
- kubectl create -f xxx.yaml
  根据配置文件创建dns资源，xxx.yaml为yaml文件位置。
- kubectl describe dns xxx.yaml/yyy
  查看dns信息，可以通过dns的名称或dns配置文件来确定dns
- kubectl get dns
  查看所有dns的简要信息
  ![](https://notes.sjtu.edu.cn/uploads/upload_bbb1ba2e97e3160a9892ed359abd2a16.png)

(没有delete)

## Deployment
- kubectl create -f xxx.yaml
  根据配置文件创建deployment资源，xxx.yaml为yaml文件位置。
- kubectl describe deployment xxx.yaml/yyy
  查看deployment信息，可以通过deployment的名称或deployment配置文件来确定deployment
- kubectl get deployment
  查看所有deployment的简要信息
  ![](https://notes.sjtu.edu.cn/uploads/upload_bcdb55bbf73a46200f01e02d3cf70870.png)

- kubectl delete deployment yyy
  根据deployment名称删除对应deployment

## Autoscaler
- kubectl create -f xxx.yaml
  根据配置文件创建autoscaler资源，xxx.yaml为yaml文件位置。
- kubectl describe autoscaler xxx.yaml/yyy
  查看autoscaler信息，可以通过autoscaler的名称或autoscaler配置文件来确定autoscaler
- kubectl get autoscaler
  查看所有autoscaler的简要信息
  ![](https://notes.sjtu.edu.cn/uploads/upload_7ba82efee8768df9bf66fc0547de4dd5.png)

- kubectl delete autoscaler yyy
  根据autoscaler名称删除对应autoscaler

## GPUJob
- kubectl create -f xxx.yaml
  根据配置文件创建gpujobr资源，xxx.yaml为yaml文件位置。
- kubectl describe gpujob xxx.yaml/yyy
  查看gpujob信息，可以通过gpujob的名称或gpujob配置文件来确定gpujob
- kubectl get gpujob
  查看所有gpujob的简要信息，字段包括：
  NAME  POD-NODE  POD-STATUS  POD-ID  POD-STIME

(没有delete)

## Function
- kubectl create -f xxx.yaml
  根据配置文件创建function资源，xxx.yaml为yaml文件位置。
- kubectl describe function xxx.yaml/yyy
  查看function信息，可以通过function的名称或function配置文件来确定function
- kubectl get function
  查看所有function的简要信息，字段包括：
  NAME  VERSION  URL

- kubectl delete function yyy
  根据function名称删除对应function

## statemachine
- kubectl create -f xxx.json
  根据配置文件创建statemachine资源，xxx.json为json文件位置。
- kubectl describe statemachine yyy
  查看statemachine信息，可以通过statemachine的名称来确定statemachine
- kubectl get statemachine
  查看所有function的简要信息，字段包括：
  NAME  STARTAT  URL

- kubectl delete statemachine yyy
  根据statemachine名称删除对应statemachine


## 其它
- kubectl hello
  测试kubectl，如果正常会返回：Hi! This is kubectl for minik8s.
