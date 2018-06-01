package main

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func FindPodsInCrashloopBackoff(client kubernetes.Interface, pod *v1.Pod) (bool, string) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			annotations := pod.GetAnnotations()
			annotations["nais.io/gardener/status"] = "bad"
			client.CoreV1().Pods(pod.Namespace).Update(pod)
			return true, pod.Name
		}
	}
	return false, ""
}
