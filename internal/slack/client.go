package slack

import (
	"context"
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/slack/handlers"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient struct {
	ApiClient *slack.Client
	Client    *socketmode.Client
	Handlers  map[string]handlers.EventHandler
	Logger    *log.Logger
	Context   context.Context
}

func NewClient(ctx context.Context, apiToken, appToken string, logger *log.Logger) *SlackClient {
	api := slack.New(apiToken, slack.OptionAppLevelToken(appToken))
	client := socketmode.New(api)

	return &SlackClient{
		ApiClient: api,
		Client:    client,
		Handlers:  make(map[string]handlers.EventHandler),
		Logger:    logger,
		Context:   ctx,
	}
}

func (sc *SlackClient) RegisterHandler(eventType string, handler handlers.EventHandler) {
	sc.Handlers[eventType] = handler
}

func (sc *SlackClient) StartListeningLoop() {
	go func() {
		for {
			select {
			case <-sc.Context.Done():
				sc.Logger.Println("Context cancelled, stopping listener loop.")
				return
			case evt := <-sc.Client.Events:
				sc.handleEvent(evt)
			}
		}
	}()
	sc.Client.Run()
}

func (sc *SlackClient) handleEvent(evt socketmode.Event) {
	switch evt.Type {
	case socketmode.EventTypeConnecting:
		sc.Logger.Println("Attempting to connect to Slack via Socket Mode.")
	case socketmode.EventTypeConnected:
		sc.Logger.Println("Connected to Slack via Socket Mode.")
	case socketmode.EventTypeHello:
		sc.Logger.Println("Received hello event from Slack.")
	case socketmode.EventTypeEventsAPI:
		eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
		if !ok {
			sc.Logger.Printf("Error asserting EventsAPIEvent: %v", evt.Data)
			return
		}
		sc.dispatchEvent(eventsAPIEvent)
		sc.Client.Ack(*evt.Request)
	default:
		sc.Logger.Printf("Unhandled event type received: %v", evt.Type)
	}
}

func (sc *SlackClient) dispatchEvent(eventsAPIEvent slackevents.EventsAPIEvent) {
	if handler, ok := sc.Handlers[eventsAPIEvent.InnerEvent.Type]; ok {
		err := handler.Handle(eventsAPIEvent)
		if err != nil {
			sc.Logger.Printf("Error handling event: %v", err)
		}
	} else {
		sc.Logger.Printf("No handler registered for event type: %s", eventsAPIEvent.Type)
	}
}
