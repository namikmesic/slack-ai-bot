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

		_, err := h.APIClient.PostEphemeral(appMentionEvent.Channel, appMentionEvent.User, slack.MsgOptionText("Please note, as an AI, there's a possibility I might not always be accurate!", false))

		if err != nil {
			return fmt.Errorf("failed to post warning message: %w", err)
		}

		_, _, err = h.APIClient.PostMessage(appMentionEvent.Channel, slack.MsgOptionText("Hello, thanks for mentioning me!", false), slack.MsgOptionTS(appMentionEvent.TimeStamp))

		if err != nil {
			return fmt.Errorf("failed to post response message: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to assert AppMentionEvent: %v", event.InnerEvent.Data)
}
