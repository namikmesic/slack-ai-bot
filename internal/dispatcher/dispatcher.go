package dispatcher

import (
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/handlers"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type EventDispatcher struct {
	ApiClient        *slack.Client
	WsClient         *socketmode.Client
	APIEventHandlers map[string]handlers.EventsAPIEventHandler
	Logger           *log.Logger
}

func NewEventDispatcher(apiClient *slack.Client, wsClient *socketmode.Client, logger *log.Logger) *EventDispatcher {
	return &EventDispatcher{
		ApiClient:        apiClient,
		WsClient:         wsClient,
		APIEventHandlers: make(map[string]handlers.EventsAPIEventHandler),
		Logger:           logger,
	}
}

func (h *EventDispatcher) RegisterNeweventsAPIEventHandler(eventType string, handler handlers.EventsAPIEventHandler) {
	h.APIEventHandlers[eventType] = handler
}

func (h *EventDispatcher) Dispatch(event socketmode.Event) {
	switch event.Type {
	case socketmode.EventTypeConnecting:
		h.Logger.Printf("Attempting to establish connection, current status: %v", event.Type)
	case socketmode.EventTypeConnected:
		h.Logger.Printf("Sucesfully established connection, current status: %v", event.Type)
	case socketmode.EventTypeHello:
		h.Logger.Printf("Received hello event from Slack%v", event.Type)
	case socketmode.EventTypeEventsAPI:
		eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
		if !ok {
			h.Logger.Printf("Error asserting EventsAPIEvent: %v", event.Data)
			return
		}
		handler, ok := h.APIEventHandlers[eventsAPIEvent.InnerEvent.Type]
		if !ok {
			h.Logger.Printf("No handler registered for event type: %s", eventsAPIEvent.Type)
			return
		}

		if err := handler.Handle(eventsAPIEvent); err != nil {
			h.Logger.Printf("Error handling event: %v", err)
			return
		}

		h.WsClient.Ack(*event.Request)
	}
}
