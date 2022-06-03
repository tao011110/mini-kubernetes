# mini-kubernetes资源定义格式
参见`repo_url/docs/yaml_template`和`repo_url/docs/serverless template`

## node
node上的kubelet启动时输入node的name和master的url, kubelet使用自己的ip和端口向master注册
````shell
./kubelet [-name node1] -master 10.119.11.140:8000
````

## pod
仿照kubernetes格式定义
````yaml
apiVersion: v1
kind: Pod
metadata:
  name: string
  label: string
nodeName: <string>
nodeSelector:
    with: pod-name
    notwith: pod-name

spec:
  containers:
    - name: string
      image: string
      command: [string]
      args: [string]
      workingDir: string
      volumeMounts:
        - name: string
          mountPath: string

      ports:
        - name: string
          containerPort: int
          hostPort: int
          protocol: string

      resources:
        limits:
          cpu: string
          memory: string
        requests:
          cpu: string
          memory: string

  volumes:
    - name: string
      hostPath: string
````
## service
仿照kubernetes定义, 提供nodeIP和clusterIP两种可选项

- clusterIP:
````yaml
apiVersion: v1
kind: Service
metadata:
  name: engine
spec:
  type: ClusterIP
  clusterIP: string       
  ports:                                  
    - port: 8080                          
      targetPort: 8080                    
      protocol: TCP
  selector:
    name: enginehttpmanage
````
- nodeIP:
````yaml
apiVersion: v1
kind: Service
metadata:
  name: servicename
spec:
  type: NodePort
  clusterIP: string
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 8080
      protocol: tcp
  selector:
    name: servicename
````

## deployment
仿照kubernetes定义
````yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deployment-hjk
spec:
  replicas: 3
  template:
    metadata:
      labels:
        name: pod-hjk
    spec:
      containers:
        - name: fileserver
          image: dplsming/nginx-fileserver:1.0
          volumeMounts:
            - name: fileserver-volume
              mountPath: /usr/share/nginx/html/files
          ports:
            - name: fileserver80
              containerPort: 80
              protocol: TCP

        - name: downloader
          image: dplsming/aria2ng-downloader:1.0
          volumeMounts:
            - name: downloader-volume
              mountPath: /data
          ports:
            - name: downloader6800
              containerPort: 6800
              protocol: TCP
            - name: downloader6880
              containerPort: 6880
              protocol: TCP

      volumes:
        - name: fileserver-volume
          hostPath:
            path: /home/hjk

        - name: downloader-volume
          hostPath:
            path: /home/hjk
````

## auto-scaler
仿照kubernetes定义, 做细微改动, auto-scaler中定义pod参数, 使auto-scaler定义不再依赖deployment

可以定义最小/最大平均CPU(单位为m, 毫核), 最小/最大平均内存占用(单位为K/M/G)
````yaml
apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: mynginx
spec:
  minReplicas: 2
  maxReplicas: 10
  metrics:
    CPU:
      targetMinValue : num
      targetMaxValue : num
    memory:
      targetMinValue: num
      targetMaxValue: num
  template:
    metadata:
      labels:
        name: pod-hjk
    spec:
      containers:
        - name: fileserver
          image: dplsming/nginx-fileserver:1.0
          volumeMounts:
            - name: fileserver-volume
              mountPath: /usr/share/nginx/html/files
          ports:
            - name: fileserver80
              containerPort: 80
              protocol: TCP

        - name: downloader
          image: dplsming/aria2ng-downloader:1.0
          volumeMounts:
            - name: downloader-volume
              mountPath: /data
          ports:
            - name: downloader6800
              containerPort: 6800
              protocol: TCP
            - name: downloader6880
              containerPort: 6880
              protocol: TCP

      volumes:
        - name: fileserver-volume
          hostPath:
            path: /home/hjk

        - name: downloader-volume
          hostPath:
            path: /home/hjk
````

## dns
````yaml
kind: DNS
name: dns-name
host: example.com
paths:
  - path: /
    service: service-name
    port: 80
  - path: /temp
    service: service-name
    port: 80
````

## GPU job
````yaml
name: <name>
kind: GPUJob
sourceCodePath: <path>
MakefilePath: <path>
ResultPath: <dirPath, end with '/'>
slurm:
    jobName: <string>
    partition: <string, e.g.a100>
    cpusPerTask: <int>
    ntasksPerNode: <int>
    Node: <int>
    GPU: <int>
    time: <dd-hh:mm:ss, string>
    targetExecutableFileName: <string>
````

## function/activity
````yaml
kind: activity
name: <string>
function: <.py文件路径>
requirements: <requirements.txt文件路径>
version: <int>
````

## workflow/state machine
在AWS state machine定义格式基础上适当简化
````json
{
    "Name": "test_state_machine1",
    "StartAt": "Check Stock Price",
    "States": {
        "Check Stock Price": {
            "Type": "Task",
            "Resource": "$activityName$",
            "Next": "Generate Buy/Sell recommendation"
        },
        "Generate Buy/Sell recommendation": {
            "Type": "Task",
            "Resource": "$activityName$",
            "Next": "Request Human Approval"
        },
        "Request Human Approval": {
            "Type": "Task",
            "Resource": "$activityName$",
            "Next": "Buy or Sell?"
        },
        "Buy or Sell?": {
            "Type": "Choice",
            "Choices": [
                {
                    "Variable": "$.recommended_type",
                    "StringEquals": "buy",
                    "Next": "Buy Stock"
                },
                {
                    "Variable": "$.recommended_type",
                    "StringEquals": "sell",
                    "Next": "Sell Stock"
                }
            ]
        },
        "Buy Stock": {
            "Type": "Task",
            "Resource": "$activityName$",
            "Next": "Report Result"
        },
        "Sell Stock": {
            "Type": "Task",
            "Resource": "$activityName$",
            "Next": "Report Result"
        },
        "Report Result": {
            "Type": "Task",
            "Resource": "$activityName$",
            "End": true
        }
    }
}
````
