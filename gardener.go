package gardener

import (
	"github.com/golang/glog"
	informercorev1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type NaisJanitor struct {
	podLister       listercorev1.PodLister
	podListerSynced cache.InformerSynced
	queue           workqueue.RateLimitingInterface
}

func NewNaisJanitor(client *kubernetes.Clientset,
	podInformer informercorev1.PodInformer) *NaisJanitor {

	janitor := &NaisJanitor{
		podLister:       podInformer.Lister(),
		podListerSynced: podInformer.Informer().HasSynced,
	}

	podInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(pod interface{}) {
				//podObj := pod.(*v14.Pod)
				//glog.Info("pod added: ", podObj.Name)
			},
			UpdateFunc: func(oldPod, newPod interface{}) {
				//oldPodObj := oldPod.(*v14.Pod)
				newPodObj := newPod.(*api.Pod)

				//glog.Info("pod changed.old: ", oldPodObj.Name, oldPodObj.UID)
				//glog.Info("pod changed.new: ", newPodObj.Name, newPodObj.UID)
				janitor.findPodsInCrashloopBackoff(newPodObj)

			},
			DeleteFunc: func(pod interface{}) {
				//podObj := pod.(*v14.Pod)
				//glog.Info("pod deleted: ", podObj.Name)
			},
		},
	)
	return janitor
}

func (janitor *NaisJanitor) findPodsInCrashloopBackoff(pod *api.Pod) {
	for _, containerStatus := range pod.Status.ContainerStatuses {
		restartCount := containerStatus.RestartCount
		//state := containerStatus.State
		//lastState := containerStatus.LastTerminationState

		glog.Info("name: ", pod.Name)
		glog.Info("restartcount: ", restartCount)
		//glog.Info("state: ", state)
		//glog.Info("lastState: ", lastState)
	}

}

func (janitor *NaisJanitor) Run(stop <-chan struct{}) {
	defer func() {
		// make sure the work queue is shut down which will trigger workers to end
		glog.Info("shutting down")
	}()

	glog.Info("waiting for cache sync")
	if !cache.WaitForCacheSync(
		stop,
		janitor.podListerSynced) {
		glog.Error("timed out waiting for cache sync")
		return
	}
}
