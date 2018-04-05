package main

import (
	"flag"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"fmt"
	"k8s.io/client-go/informers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const Port = ":8081"

func main() {
	fmt.Println("starting up...")

	sigs := make(chan os.Signal, 1) // Create channel to receive OS signals
	stop := make(chan struct{})

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT) // Register the sigs channel to receieve SIGTERM

	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")

	flag.Parse()

	glog.Infof("running on port %s", Port)
	clientSet := newClientSet(*kubeconfig)

	sharedInformers := informers.NewSharedInformerFactory(clientSet, 10*time.Minute)

	gardener := NewNaisGardener(clientSet, sharedInformers.Core().V1().Pods())

	sharedInformers.Start(stop)
	gardener.Run(stop)

	<-sigs
	fmt.Println("shutting  down...")
	close(stop)
}

// returns config using kubeconfig if provided, else from cluster context
func newClientSet(kubeconfig string) *kubernetes.Clientset {

	var config *rest.Config
	var err error

	if kubeconfig != "" {
		glog.Infof("using provided kubeconfig")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		glog.Infof("no kubeconfig provided, assuming we are running inside a cluster")
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	return clientset
}
