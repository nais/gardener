package main

import (
	"github.com/golang/glog"
	"k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

const channel = "nais_gardener"

var httpClient = &http.Client{}

func NotifyTeamsOfWeed(send func(string, string, string) error, client kubernetes.Interface, clustername string, slackUrl string, deployment *v1.Deployment) {

	annotations := deployment.GetAnnotations()

	status := annotations[annotationStatus]
	notify := annotations[annotationNotify]

	if status == "bad" && notify == "" {
		message := "The application " + deployment.Namespace + "." + deployment.Name + " has restarted more the 50 times in the cluster: " + clustername + ". The deployment will be deleted"

		annotations[annotationNotify] = "done"
		deployment, err := client.AppsV1().Deployments(deployment.Namespace).Update(deployment)
		if err != nil {
			glog.Errorf("Error while updating deployment %s %s", deployment.Name, err)
			return
		}
		err = send(message, channel, slackUrl)
		if err != nil {
			glog.Errorf("Error when posting to slack %s ", err)
			return
		}
	}
}

func SendMessage(message string, channel string, slackUrl string) error {
	slack := Client{slackUrl, httpClient}
	return slack.SimpleToChannel(message, channel)
}
