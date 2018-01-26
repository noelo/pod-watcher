package main

import (
	"fmt"

	v1 "github.com/openshift/api/build/v1"
)

type gitlabWebhook struct{}

func (t gitlabWebhook) Publish(hook v1.BuildTriggerPolicy, uri *v1.GitBuildSource) {

	if hook.GitLabWebHook != nil {
		fmt.Println("GitLab trigger found with secret " + hook.GitLabWebHook.Secret)
		fmt.Println("GitLab URL " + uri.URI)
	}
}
