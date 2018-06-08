package main

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

)

func FindPodsInCrashloopBackoff(client kubernetes.Interface, pod *v1.Pod) (bool, string) {

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			for _, set := range pod.OwnerReferences {
				replicaSet, _ := client.AppsV1().ReplicaSets(pod.Namespace).Get(set.Name, v12.GetOptions{})
				println("set: ", replicaSet.Name)

				for _, depl := range replicaSet.OwnerReferences {
					deployment, _ := client.AppsV1().Deployments(pod.Namespace).Get(depl.Name, v12.GetOptions{})
					println("deployment: ", deployment.Name)
					annotations := deployment.GetAnnotations()
					annotations["nais.io/gardener/status"] = "bad"
					client.AppsV1().Deployments(pod.Namespace).Update(deployment)
					return true, pod.Name
				}
			}
		}
	}
	return false, ""
}
