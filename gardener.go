package main

import (
	"github.com/golang/glog"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/api/core/v1"
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
				gardener.findPodsInCrashloopBackoff(newPod.(*v1.Pod))

			},
		},
	)
	return gardener
}

func (gardener *gardener) findPodsInCrashloopBackoff(pod *v1.Pod) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		glog.Infof("restartcount: %s: %d ", pod.Name, containerStatus.RestartCount)
	}

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
