# mini-kubernetes serverless functional编程模型
By hjk 

### 函数定义(参考 aliyun serverless)
````python
def handler(env):
    print("param1 in url is" + env.uri.get("param1", ""))
    print(env.body["name1"])
    return '{"hjk": "hjk"}'
````
- 函数式服务写在一个.py文件内
- 函数的最顶层函数(http trigger后直接调用的函数)以handler命名, 参数为一个有`uri`(表示请求的params)和`body`(请求体)两个成员的结构体
- 以`http:127.0.0.1:8080/?user=test`为例:
	- 获取路径中的`user`参数:
   ````python
   env.uri.get("user", "")
   ````
	- 请求体为**json格式**, 已经预先解析, 获取请求体中的成员:
   ````python
   env.body["name1"]
   ````
- 返回值为字符串, 如果需要返回结构体或其他信息需要预先解析为字符串
### 提交一个函数式服务
- 用户需要通过apiserver提交如下文件:
	- requirements.txt: 执行所需要的依赖, 部分高频使用的依赖已经预装在基础镜像中
	- python函数文件
	- yaml定义文件, 见示例
	- 提交命令:
   ````shell
   kubectl apply def.yaml 
   ````
	- 在提交函数时服务后, 用户可以通过masterIP:def.ActiverPort/function/taskName发起http trigger

### 制作镜像
- 拉取基础镜像: jingkaihe/python-functional-template
- (三个启动后运行的指令)启动镜像, 读取.py和requirements.txt内容, 在启动时使用echo指令将其分别写入到`/home/functionalTemplate/handler.py`和`/requirements.txt`, 命令类似于：
````shell
echo $string > /home/functionalTemplate/handler.py
````
- 运行`./prepare.sh`
- docker commit & docker pull提交镜像, 将镜像名存储在etcd中
- 建立预定义的service(单pod, pod单容器), pod的启动后命令为`./start.sh`
- service不需要立即创建
### http trigger
- 当向对应的端口和路径发送http请求后,Activer需要创建service(如果没有正在运行的), 在已经存在的service中做负载均衡, 记录访问频率并根据访问频率对service中的pod replica做扩缩容(缩容最低到0, 扩容有上限), 目标数目为上个30s收到的trigger数目除以100
- Activer调用service并向用户返回
### state machine(参考aws step functions)
####定义方式
- 提供json定义, 见示例(语法规则: https://docs.aws.amazon.com/zh_cn/step-functions/latest/dg/input-output-resultpath.html)
####调用方式
- 由activer与service交互完成执行
#### 示例
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
表示的work flow为:
![](serverless%20template/state%20machine/stepfunctions_graph.png)
