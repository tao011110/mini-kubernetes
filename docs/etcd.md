# etcd中存储的结构及其作用

````
.etcd ______ {pod meta}(k = podname, v = json(struct Pod))
		|___ {pod instance meta}(k = replica name(podname-`uuid`), v = json(struct PodInstance))
		|___ pod_name_list(optional)
		|___ deployment_name_list [name1, name2]//(apiserver->conntroller)
		|___ horizontalPodAutoScaler_name_list//(apiserver->conntroller)
		|___ pod_instance_ID_list //(controller->scheduler)
		|___ {node_podInstance_name_list(per node)}//(scheduler->kubelet)
		|___ {node_resource_record}
		|___ node_name_list
		|___ {podInstance_resource_record}
		|___ {replica_name_list(per pod)}
		|___ function_name_list
		|___ {function(struct)}
````
