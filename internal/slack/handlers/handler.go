package handlers

import (
	"github.com/slack-go/slack/slackevents"
)

// EventHandler defines the interface for handling Slack events.
type EventHandler interface {
	Handle(event slackevents.EventsAPIEvent) error
}
