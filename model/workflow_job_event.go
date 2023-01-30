package model

import (
	"encoding/json"
	"io"

	"github.com/google/go-github/v47/github"
)

type WorkflowJobEvent struct {
	WorkflowJob *github.WorkflowJob `json:"workflow_job,omitempty"`

	Action *string `json:"action,omitempty"`

	// The following fields are only populated by Webhook events.

	// Org is not nil when the webhook is configured for an organization or the event
	// occurs from activity in a repository owned by an organization.
	Org          *github.Organization `json:"organization,omitempty"`
	Repo         *github.Repository   `json:"repository,omitempty"`
	Sender       *github.User         `json:"sender,omitempty"`
	Installation *github.Installation `json:"installation,omitempty"`

	// Not present in google/go-github in v50.0.0 yet
	Deployment *github.Deployment `json:"deployment,omitempty"`
}

func WorkflowJobEventFromJSON(data io.Reader) *WorkflowJobEvent {
	decoder := json.NewDecoder(data)
	var event WorkflowJobEvent
	if err := decoder.Decode(&event); err != nil {
		return nil
	}
	return &event
}
