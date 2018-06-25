package main

import (
	"testing"
	"k8s.io/client-go/kubernetes/fake"
	v13 "k8s.io/api/apps/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAnnotationsIsNotDoneWhenStatusIsNotBad(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()

	deployment := &v13.Deployment{}
	deployment.Name = deploymentname
	deployment.Namespace = namespace
	deployment.UID = "1234"
	deployment.Annotations = make(map[string]string)
	k8sclient.AppsV1().Deployments(namespace).Create(deployment)

	stub := func(string) error { return nil }
	err := NotifyTeamsOfWeed(stub, k8sclient, "clustername", deployment)
	if err != nil {
		t.Fail()
	}

	depl, _ := k8sclient.AppsV1().Deployments(namespace).Get("deployment", v12.GetOptions{})
	if depl.GetAnnotations()[annotationNotify] == "done" {
		t.Fail()
	}
}

func TestAnnotationsIsDoneWhenStatusIsBad(t *testing.T) {
	k8sclient := fake.NewSimpleClientset()

	deployment := &v13.Deployment{}
	deployment.Name = deploymentname
	deployment.Namespace = namespace
	deployment.UID = "1234"
	deployment.Annotations = make(map[string]string)
	deployment.Annotations[annotationStatus] = "bad"
	k8sclient.AppsV1().Deployments(namespace).Create(deployment)

	stub := func(string) error { return nil }
	err := NotifyTeamsOfWeed(stub, k8sclient, "clustername", deployment)
	if err != nil {
		t.Fail()
	}
	depl, _ := k8sclient.AppsV1().Deployments(namespace).Get("deployment", v12.GetOptions{})
	if depl.GetAnnotations()[annotationNotify] != "done" {
		t.Fail()
	}
}
