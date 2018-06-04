package main

import (
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestWillTriggerOn50Restarts(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()

	namespace := "ns"
	podname := "name"
	pod := createPod(k8sclient, namespace, podname, 51)

	triggered, _ := FindPodsInCrashloopBackoff(k8sclient, pod)

	if triggered == false {
		t.Fail()
	}

	pod1, _ := k8sclient.CoreV1().Pods(namespace).Get(podname, v12.GetOptions{})

	if pod1.GetAnnotations()["nais.io/gardener/status"] == "" || pod1.GetAnnotations()["nais.io/gardener/status"] != "bad" {
		t.Fail()
	}
}

func TestWillNotTriggerOnLessThan50Restarts(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()

	namespace := "ns"
	podname := "name"
	pod := createPod(k8sclient, namespace, podname, 49)
	triggered, _ := FindPodsInCrashloopBackoff(k8sclient, pod)
	if triggered == true {
		t.Fail()
	}
	pod1, _ := k8sclient.CoreV1().Pods(namespace).Get(podname, v12.GetOptions{})

	if pod1.GetAnnotations()["nais.io/gardener/status"] != "" || pod1.GetAnnotations()["nais.io/gardener/status"] == "bad" {
		t.Fail()
	}
}

func createPod(k8sclient *fake.Clientset, namespace string, podName string, restarts int32) (*v1.Pod) {
	pod := &v1.Pod{
		Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: restarts}},
		},
	}

	pod.Name = podName
	pod.Namespace = namespace
	pod.Annotations = make(map[string]string)
	k8sclient.CoreV1().Pods(namespace).Create(pod)
	return pod
}
