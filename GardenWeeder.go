package main

import (
	"k8s.io/api/apps/v1"
	"net/http"
)

func NotifyTeamsOfWeed(deployment *v1.Deployment) {

	annotations := deployment.GetAnnotations()
	lables := deployment.GetLabels()

	app := lables["app"]
	status := annotations["nais.io/gardener.status"]

	if status == "bad" {
		slack := Client{"https://hooks.slack.com/services/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", &http.Client{}}
		//slack.Send(Message{Text:"", Channel:"", })
		slack.Simple("The application " + app + " has restarted more the 50 times. The deployment will be deleted")
	}
}
