package main

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/golang/glog"
)

func FindPodsInCrashloopBackoff(client kubernetes.Interface, pod *v1.Pod) {

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.RestartCount > 50 {
			for _, set := range pod.OwnerReferences {
				replicaSet, err := client.AppsV1().ReplicaSets(pod.Namespace).Get(set.Name, metav1.GetOptions{})
				if err != nil {
					glog.Errorf("cannot fetch replicaset %s.%s: %s", pod.Namespace, set.Name, err)
					return

				}
				for _, depl := range replicaSet.OwnerReferences {
					deployment, err := client.AppsV1().Deployments(pod.Namespace).Get(depl.Name, metav1.GetOptions{})
					if err != nil {
						glog.Errorf("cannot fetch  deployment %s.%s: %s", pod.Namespace, depl.Name, err)
						return
					}
					annotations := deployment.GetAnnotations()
					annotations[annotationStatus] = "bad"

					_, err2 := client.AppsV1().Deployments(pod.Namespace).Update(deployment)
					if err2 != nil {
						glog.Errorf("cannot update deployment %s.%s: %s", pod.Namespace, deployment.Name, err2)
					}

				}
			}
		}
	}
}
