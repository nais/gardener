package main

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"
)

func FindPodsInCrashloopBackoff(client kubernetes.Interface, pod *v1.Pod) (bool, string) {

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			for _, set := range pod.OwnerReferences {
				replicaSet, err := client.AppsV1().ReplicaSets(pod.Namespace).Get(set.Name, v12.GetOptions{})
				if err != nil {
					glog.Error("cannot fetch replicaset ", set, err)
					return false, ""
				}
				println("set: ", replicaSet.Name)

				for _, depl := range replicaSet.OwnerReferences {
					deployment, err := client.AppsV1().Deployments(pod.Namespace).Get(depl.Name, v12.GetOptions{})
					if err != nil {
						glog.Error("cannot fetch deployment ", deployment, err)
						return false, ""
					}
					println("deployment: ", deployment.Name)
					annotations := deployment.GetAnnotations()
					annotations["nais.io/gardener.status"] = "bad"
					upDeployment, err := client.AppsV1().Deployments(pod.Namespace).Update(deployment)
					if err != nil {
						glog.Error("cannot update deployment ", upDeployment, err)
						return false, ""
					}
					return true, upDeployment.Name
				}
			}
		}
	}
	return false, ""
}
