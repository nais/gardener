package main

import (
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestWillTriggerOn50Restarts(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()
	pod := &v1.Pod{
		Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: 51}},
		},
	}
	pod.Annotations = make(map[string]string)
	triggered, _ := FindPodsInCrashloopBackoff(k8sclient, pod)
	if triggered == false {
		t.Fail()
	}
	if pod.GetAnnotations()["nais.io/gardener/status"] == "" || pod.GetAnnotations()["nais.io/gardener/status"] != "bad" {
		t.Fail()
	}
}

func TestWillNotTriggerOnLessThan50Restarts(t *testing.T) {
	clientset := fake.NewSimpleClientset()

	pod := &v1.Pod{
		Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: 49}},
		},
	}
	pod.Annotations = make(map[string]string)
	triggered, _ := FindPodsInCrashloopBackoff(clientset, pod)
	if triggered == true {
		t.Fail()
	}
	if pod.GetAnnotations()["nais.io/gardener/status"] != "" || pod.GetAnnotations()["nais.io/gardener/status"] == "bad" {
		t.Fail()
	}
}
