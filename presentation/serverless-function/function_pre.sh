./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-function/function.yaml
./kubectl get function
./kubectl describe function function_test
# curl -H "Content-type: application/json" -X GET -d '{"userType":1}' http://10.119.11.140:3306/function/function_test
