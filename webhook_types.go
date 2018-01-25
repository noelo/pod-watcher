package webhooks

import (
	v1 "github.com/openshift/api/build/v1"
)

// WebhookProcessor used a base interface for different webhook publisher implementations
type WebhookProcessor interface {
	Publish(hook v1.BuildTriggerPolicy, uri v1.GitBuildSource)
}
