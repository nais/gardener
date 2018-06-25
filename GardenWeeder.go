package main

import (
	"k8s.io/api/apps/v1"
	"net/http"
	"fmt"
	"k8s.io/client-go/kubernetes"
)

const slackUrl = "https://hooks.slack.com/services/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"

var httpClient = &http.Client{}

func NotifyTeamsOfWeed(send func(string) error, client kubernetes.Interface, clustername string, deployment *v1.Deployment) error {

	annotations := deployment.GetAnnotations()

	status := annotations[annotationStatus]
	notify := annotations[annotationNotify]

	if status == "bad" && notify == "" {
		message := "The application " + deployment.Namespace + "." + deployment.Name + " has restarted more the 50 times in the cluster: " + clustername + ". The deployment will be deleted"

		err := send(message)
		if err != nil {
			return fmt.Errorf("Error when posting to slack %s ", err)
		}
		annotations[annotationNotify] = "done"
		deployment, err = client.AppsV1().Deployments(deployment.Namespace).Update(deployment)
		if err != nil {
			return fmt.Errorf("Error while updating deployment %s %s", deployment.Name, err)
		}

	}
	return nil
}

func SendMessage(message string) error {
	slack := Client{slackUrl, httpClient}
	return slack.Simple(message)
}
