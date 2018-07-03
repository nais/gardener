package main

import (
	"k8s.io/api/apps/v1"
	"net/http"
	"k8s.io/client-go/kubernetes"
	"github.com/golang/glog"
)

const slackUrl = "https://hooks.slack.com/services/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"
const channel = "nais_gardener"

var httpClient = &http.Client{}

func NotifyTeamsOfWeed(send func(string, string) error, client kubernetes.Interface, clustername string, deployment *v1.Deployment) {

	annotations := deployment.GetAnnotations()

	status := annotations[annotationStatus]
	notify := annotations[annotationNotify]

	if status == "bad" && notify == "" {
		message := "The application " + deployment.Namespace + "." + deployment.Name + " has restarted more the 50 times in the cluster: " + clustername + ". The deployment will be deleted"

		err := send(message, channel)
		if err != nil {
			glog.Errorf("Error when posting to slack %s ", err)
			return
		}
		annotations[annotationNotify] = "done"
		deployment, err = client.AppsV1().Deployments(deployment.Namespace).Update(deployment)
		if err != nil {
			glog.Errorf("Error while updating deployment %s %s", deployment.Name, err)
		}

	}
}

func SendMessage(message string, channel string) error {
	slack := Client{slackUrl, httpClient}
	return slack.SimpleToChannel(message, channel)
}
