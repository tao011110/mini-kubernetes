./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-function/function.yaml
./kubectl get function
./kubectl get pods
./kubectl describe function function_test
# curl -H "Content-type: application/json" -X GET -d '{"userType":1}' http://10.119.11.140:3306/function/function_test
# 多次访问，使得动态扩容
./kubectl get pods
# 停止访问，scale-to-0
./kubectl get pods

./kubectl update hard -f /home/go/src/mini-kubernetes/presentation/serverless-function/function_update.yaml
./kubectl get function
./kubectl describe function function_test
# curl -H "Content-type: application/json" -X GET -d '{"userType":1}' http://10.119.11.140:3306/function/function_test
./kubectl delete function function_test
./kubectl get function
