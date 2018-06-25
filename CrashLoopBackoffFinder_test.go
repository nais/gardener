package main

import (
	"testing"

	v13 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

const namespace = "ns"
const podname = "podname"
const deploymentname = "deployment"
const rsname = "replicaSet"

func TestWillTriggerOn50Restarts(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()

	pod := createResources(k8sclient, 51)

	triggered, _ := FindPodsInCrashloopBackoff(k8sclient, pod)

	if triggered == false {
		t.Fail()
	}

	depl, _ := k8sclient.AppsV1().Deployments(namespace).Get("deployment", v12.GetOptions{})

	if depl.GetAnnotations()[annotationStatus] == "" || depl.GetAnnotations()[annotationStatus] != "bad" {
		t.Fail()
	}
}

func TestWillNotTriggerOnLessThan50Restarts(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()
	pod := createResources(k8sclient, 49)
	triggered, _ := FindPodsInCrashloopBackoff(k8sclient, pod)
	if triggered == true {
		t.Fail()
	}
	depl, _ := k8sclient.AppsV1().Deployments(namespace).Get("deployment", v12.GetOptions{})

	if depl.GetAnnotations()[annotationStatus] != "" || depl.GetAnnotations()[annotationStatus] == "bad" {
		t.Fail()
	}
}

func createResources(k8sclient *fake.Clientset, restarts int32) *v1.Pod {

	deployment := &v13.Deployment{}
	deployment.Name = deploymentname
	deployment.Namespace = namespace
	deployment.UID = "1234"
	deployment.Annotations = make(map[string]string)

	replicaSet := &v13.ReplicaSet{}
	replicaSet.Name = rsname
	replicaSet.Namespace = namespace
	replicaSet.UID = "12345"
	replicaSet.OwnerReferences = []v12.OwnerReference{{UID: deployment.GetUID()}}

	pod := &v1.Pod{
		Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: restarts}}},
	}
	pod.OwnerReferences = []v12.OwnerReference{{UID: replicaSet.GetUID()}}

	pod.Name = podname
	pod.Namespace = namespace
	pod.Annotations = make(map[string]string)

	k8sclient.AppsV1().Deployments(namespace).Create(deployment)
	k8sclient.AppsV1().ReplicaSets(namespace).Create(replicaSet)
	k8sclient.CoreV1().Pods(namespace).Create(pod)

	return pod
}
