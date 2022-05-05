package pod

import (
	"mini-kubernetes/tools/def"
	"time"
)

func PodDaemon(podInstance *def.PodInstance) {
	time.Sleep(time.Duration(podInstance.Spec.LivenessProbe.InitialDelaySeconds) * time.Second)
	for {
		if podInstance.Status == def.RUNNING {
			DetectLiveness(podInstance)
			if podInstance.PodInstanceStatus.LastDetectSuccess == false &&
				podInstance.PodInstanceStatus.ConsecutiveFailures >= uint(podInstance.Spec.LivenessProbe.FailureThreshold) {
				RestartPod(podInstance)
			}
		}
		if podInstance.Status == def.SUCCEEDED || podInstance.Status == def.FAILED {
			break
		}
		time.Sleep(time.Duration(podInstance.Spec.LivenessProbe.PeriodSeconds) * time.Second)
	}
	return
}
