# etcd中存储的结构及其作用
By hjk & tyc

````
.etcd __ //meta data:
		|___ {pod meta}(k = /pod/podName, v = json(struct Pod))
		|___ {pod instance meta}(k = podInstance/podName, v = json(struct PodInstance))
		|___ {replicas meta}(k = podInstance/replicaName(podName-`uuid`), v = json(struct PodInstance))
		|___ {service meta}(k = /service/serviceName, v = json(struct Service))
		|___ {function meta}(k = /function/functionName, v = json(struct Function))
		|___ {statemachine meta}(k = /statemachine/statemachineName, v = json(struct StateMachine))
		|___ {deployment meta}(k = /deployment/deploymentName, v = json(struct ParsedDeployment))
		|___ {autoscaler meta}(k = /autoscaler/autoscalerName, v = json(struct ParsedAutoscaler))
		|___ {gpujob meta}(k = /gpujob/gpujobName, v = json(struct GPUJob))
		|___ {node meta}(k = /node/nodeName, v = json(struct Node))
		|___ {dns meta}(k = /DNS/DNSName, v = json(struct DNSDetail))
		|___ {node_resource_record meta}(k = /resource_usage/nodeID)
		|___ {podInstance_resource_record meta}(k = resource_usage/podInstanceID)
		// message queue:
		|___ {node_name_list} [nodeName1, nodeName2] //(apiserver->controller)
		|___ {function_name_list} [functionName1, functionName2] //(apiserver->activer)
		|___ {deployment_list_name} [deploymentName1, deploymentName2]//(apiserver->controller)
		|___ {parsed_horizontal_pod_autoscaler_list_name} [autoscalerName1, autoscalerName2]//(apiserver->controller)
		|___ {node_podInstance_ID_list(per node)} [podInstanceID1, podInstanceID2]//(scheduler->kubelet)
		|___ {pod_instance_list_id} [podInstanceID1, podInstanceID2]//(controller->scheduler)
		|___ {replicas_name_list(per pod)} [podName1, podName2] //(apiserver->controller)
````
