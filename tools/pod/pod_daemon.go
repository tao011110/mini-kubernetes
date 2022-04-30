package pod

import (
	"time"
)

func (podInstance *PodInstance) PodDaemon() {
	time.Sleep(time.Duration(podInstance.Spec.LivenessProbe.InitialDelaySeconds) * time.Second)
	for {
		if podInstance.Status == RUNNING {
			podInstance.DetectLiveness()
			if podInstance.PodInstanceStatus.LastDetectSuccess == false &&
				podInstance.PodInstanceStatus.ConsecutiveFailures >= uint(podInstance.Spec.LivenessProbe.FailureThreshold) {
				podInstance.RestartPod()
			}
		}
		if podInstance.Status == SUCCEEDED || podInstance.Status == FAILED {
			break
		}
		time.Sleep(time.Duration(podInstance.Spec.LivenessProbe.PeriodSeconds) * time.Second)
	}
	return
}
