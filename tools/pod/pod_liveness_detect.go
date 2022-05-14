package pod

//import (
//	"fmt"
//	"mini-kubernetes/tools/def"
//	"mini-kubernetes/tools/httpget"
//)

//func DetectLiveness(podInstance *def.PodInstance) {
//	path := fmt.Sprintf("http://%s:%d%s",
//		podInstance.ClusterIP,
//		podInstance.Spec.LivenessProbe.HttpGetRequest.Port,
//		podInstance.Spec.LivenessProbe.HttpGetRequest.Path)
//	client := httpget.Post(path)
//	for _, kv := range podInstance.Spec.LivenessProbe.HttpGetRequest.HttpHeaders {
//		client = client.AddHeader(kv.Name, kv.Value)
//	}
//	_, s := client.ContentType("application/json").Execute()
//	if s != "200 OK" {
//		if podInstance.PodInstanceStatus.LastDetectSuccess {
//			podInstance.PodInstanceStatus.LastDetectSuccess = false
//			podInstance.PodInstanceStatus.ConsecutiveFailures = 1
//		} else {
//			podInstance.PodInstanceStatus.ConsecutiveFailures++
//		}
//	} else {
//		if podInstance.PodInstanceStatus.LastDetectSuccess {
//			podInstance.PodInstanceStatus.ConsecutiveSuccesses++
//		} else {
//			podInstance.PodInstanceStatus.LastDetectSuccess = true
//			podInstance.PodInstanceStatus.ConsecutiveSuccesses = 1
//		}
//	}
//}
