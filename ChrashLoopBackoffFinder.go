package main

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"fmt"
)

func FindPodsInCrashloopBackoff(client kubernetes.Interface, pod *v1.Pod) (bool, error) {

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			for _, set := range pod.OwnerReferences {
				replicaSet, err := client.AppsV1().ReplicaSets(pod.Namespace).Get(set.Name, metav1.GetOptions{})
				if err != nil {
					return false, fmt.Errorf("cannot fetch replicaset %s, %s", replicaSet, err)
				}
				for _, depl := range replicaSet.OwnerReferences {
					deployment, err := client.AppsV1().Deployments(pod.Namespace).Get(depl.Name, metav1.GetOptions{})
					if err != nil {
						return false, fmt.Errorf("cannot fetch  deployment %s, %s", deployment, err)
					}
					annotations := deployment.GetAnnotations()
					annotations[annotation_status] = "bad"

					upDeployment, err := client.AppsV1().Deployments(pod.Namespace).Update(deployment)
					if err != nil {
						return false, fmt.Errorf("cannot update deployment %s, %s", upDeployment, err)
					}
					return true, nil
				}
			}
		}
	}
	return false, nil
}
