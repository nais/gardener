package main

import (
	"k8s.io/api/core/v1"
	"github.com/golang/glog"
)

func FindPodsInCrashloopBackoff(pod *v1.Pod) (bool, string) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		glog.Infof("restartcount: %s: %d ", pod.Name, containerStatus.RestartCount)
		if containerStatus.RestartCount > 50 {
			return true, pod.Name
		}
	}
	return false, ""
}
