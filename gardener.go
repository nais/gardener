package main

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type gardener struct {
	podLister       listercorev1.PodLister
	podListerSynced cache.InformerSynced
	queue           workqueue.RateLimitingInterface
}

func NewNaisGardener(client *kubernetes.Clientset,
	podInformer informercorev1.PodInformer) *gardener {

	gardener := &gardener{
		podLister:       podInformer.Lister(),
		podListerSynced: podInformer.Informer().HasSynced,
	}

	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldPod, newPod interface{}) {
				triggered, name := FindPodsInCrashloopBackoff(client, newPod.(*v1.Pod))

				if triggered{
					glog.Infof("pod: %s is marked for weeding", name)
				}
			},
		},
	)
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
