package main

import (
	"k8s.io/api/core/v1"
	"testing"
)

func TestWillTriggerOn50Restarts(t *testing.T) {

	pod := &v1.Pod{
		Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: 51}},
		},
	}
	pod.Annotations = make(map[string]string)
	triggered, _ := FindPodsInCrashloopBackoff(pod)
	if triggered == false {
		t.Fail()
	}
	if pod.GetAnnotations()["nais.io/gardener/status"] == "" || pod.GetAnnotations()["nais.io/gardener/status"] != "bad" {
		t.Fail()
	}
}

func TestWillNotTriggerOnLessThan50Restarts(t *testing.T) {
	pod := &v1.Pod{
		Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: 49}},
		},
	}
	pod.Annotations = make(map[string]string)
	triggered, _ := FindPodsInCrashloopBackoff(pod)
	if triggered == true {
		t.Fail()
	}
	if pod.GetAnnotations()["nais.io/gardener/status"] != "" || pod.GetAnnotations()["nais.io/gardener/status"] == "bad" {
		t.Fail()
	}
}
