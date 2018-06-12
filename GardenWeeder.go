package main

import (
	"net/http"
	"k8s.io/api/apps/v1"
)

func NotifyTeamsOfWeed(deployment *v1.Deployment) {

	annotations := deployment.GetAnnotations()
	status := annotations["nais.io/gardener.status"]
	if status == "bad" {
		slack := Client{"https://hooks.slack.com/services/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", &http.Client{}}
		slack.Simple("hei audun")
	}
}
