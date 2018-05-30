package main

import (
	"k8s.io/api/core/v1"
	"testing"
)

func TestWillTriggerOn50Restarts(t *testing.T) {
	triggered, _ := FindPodsInCrashloopBackoff(
		&v1.Pod{
			Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: 51}},
			}})
	if triggered == false {
		t.Fail()
	}
}

func TestWillNotTriggerOnLessThan50Restarts(t *testing.T) {
	triggered, _ := FindPodsInCrashloopBackoff(
		&v1.Pod{
			Status: v1.PodStatus{ContainerStatuses: []v1.ContainerStatus{{RestartCount: 49}},
			}})
	if triggered == true{
			t.Fail()
	}
}