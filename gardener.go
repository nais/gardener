package main

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	informerappsv1 "k8s.io/client-go/informers/apps/v1"

	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	listerappsv1 "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	v12 "k8s.io/api/apps/v1"
)

type gardener struct {
	podLister       listercorev1.PodLister
	podListerSynced cache.InformerSynced

	deploymentLister       listerappsv1.DeploymentLister
	deploymentListerSynced cache.InformerSynced

	queue workqueue.RateLimitingInterface
}

func NewNaisGardener(client *kubernetes.Clientset,
	podInformer informercorev1.PodInformer,
	deploymentInformer informerappsv1.DeploymentInformer) *gardener {

	gardener := &gardener{
		podLister:              podInformer.Lister(),
		podListerSynced:        podInformer.Informer().HasSynced,
		deploymentLister:       deploymentInformer.Lister(),
		deploymentListerSynced: deploymentInformer.Informer().HasSynced,
	}

	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldPod, newPod interface{}) {
				triggered, err := FindPodsInCrashloopBackoff(client, newPod.(*v1.Pod))
				if err != nil {
					glog.Error("Error while looking for chrashloopbackoff on pod", newPod, err)
				}
				if triggered {
					glog.Infof("pod: %s is marked for weeding", newPod)
				}
			},
		},
	)
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldDeployment, newDeployment interface{}) {
			NotifyTeamsOfWeed(newDeployment.(*v12.Deployment))
		},
	})
	return gardener
}

func (gardener *gardener) Run(stop <-chan struct{}) {
	defer func() {
		// make sure the work queue is shut down which will trigger workers to end
		glog.Info("shutting down")
	}()

	glog.Info("waiting for cache sync")
	if !cache.WaitForCacheSync(
		stop,
		gardener.podListerSynced) {
		glog.Error("timed out waiting for cache sync")
		return
	}
}
