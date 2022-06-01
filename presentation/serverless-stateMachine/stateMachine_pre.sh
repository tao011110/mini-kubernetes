./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-stateMachine/state_machine.json
./kubectl get statemachine
./kubectl describe statemachine test_state_machine
./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-stateMachine/state1/state1.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-stateMachine/state2/state2.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-stateMachine/state3/state3.yaml
./kubectl create -f /home/go/src/mini-kubernetes/presentation/serverless-stateMachine/state4/state4.yaml
./kubectl get function
# curl -H "Content-type: application/json" -X GET -d '{"type": 1}' http://10.119.11.140:3306/state_machine/test_state_machine
# curl -H "Content-type: application/json" -X GET -d '{"type": 2}' http://10.119.11.140:3306/state_machine/test_state_machine

