package main

import "k8s.io/api/core/v1"

func FindPodsInCrashloopBackoff(pod *v1.Pod) (bool, string) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			//glog.Infof("restartcount: %s: %d ", pod.Name, containerStatus.RestartCount)
			return true, pod.Name
		}
	}
	return false, ""
}
