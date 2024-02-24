package handlers

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// AppMentionHandler handles events of type AppMention.
type AppMentionHandler struct {
	ApiClient *slack.Client
}

// NewAppMentionEventHandler creates a new instance of AppMentionHandler.
func NewAppMentionEventHandler(apiClient *slack.Client) *AppMentionHandler {
	return &AppMentionHandler{
		ApiClient: apiClient,
	}
}

// Handle checks if the event is an AppMentionEvent and processes it.
func (h *AppMentionHandler) Handle(event slackevents.EventsAPIEvent) error {
	// Assert the event type to *slackevents.AppMentionEvent
	if appMentionEvent, ok := event.InnerEvent.Data.(*slackevents.AppMentionEvent); ok {
		// Now that we've asserted the type, proceed with handling the app mention.
		_, _, err := h.ApiClient.PostMessage(appMentionEvent.Channel, slack.MsgOptionText("Hello, thanks for mentioning me!", false))
		if err != nil {
			return fmt.Errorf("failed to post message: %w", err)
		}
		return nil
	}
	// If the event is not an AppMentionEvent, this handler does nothing.
	// Depending on your design, you might want to log this or handle it differently.
	return fmt.Errorf("event passed to AppMentionHandler is not an AppMentionEvent")
}
