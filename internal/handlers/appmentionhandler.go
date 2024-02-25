package handlers

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type AppMentionHandler struct {
	APIClient *slack.Client
}

func NewAppMentionEventHandler(apiClient *slack.Client) *AppMentionHandler {
	return &AppMentionHandler{
		APIClient: apiClient,
	}
}

func (h *AppMentionHandler) Handle(event slackevents.EventsAPIEvent) error {
	if appMentionEvent, ok := event.InnerEvent.Data.(*slackevents.AppMentionEvent); ok {
		_, _, err := h.APIClient.PostMessage(appMentionEvent.Channel, slack.MsgOptionText("Hello, thanks for mentioning me!", false))

		if err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to assert AppMentionEvent: %v", event.InnerEvent.Data)
}
