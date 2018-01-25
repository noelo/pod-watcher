package webhooks

import (
	"fmt"

	v1 "github.com/openshift/api/build/v1"
)

type githubWebhook struct{}

func (t githubWebhook) Publish(hook v1.BuildTriggerPolicy, uri v1.GitBuildSource) {

	if hook.GitHubWebHook != nil {
		fmt.Println("Github trigger found with secret " + hook.GitHubWebHook.Secret)
		fmt.Println("Github URL " + uri.URI)
	}
}
