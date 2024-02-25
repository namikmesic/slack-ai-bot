package handlers

import "github.com/slack-go/slack/slackevents"

type EventsAPIEventHandler interface {
	Handle(event slackevents.EventsAPIEvent) error
}
