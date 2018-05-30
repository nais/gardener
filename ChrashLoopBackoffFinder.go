package main

import (
	"k8s.io/api/core/v1"
)

func FindPodsInCrashloopBackoff(pod *v1.Pod) (bool, string) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			labels := pod.GetLabels()
			labels["nais.io/status"]="bad"
			return true, pod.Name
		}
	}
	return false, ""
}
