package dispatcher

import (
	"log"
	"sync"

	"github.com/sashabaranov/go-openai"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

type MessageDispatcher struct {
	APIClient       *slack.Client
	WSClient        *socketmode.Client
	AIClient        *openai.Client
	botID           string
	Logger          *log.Logger
	TrackedThreads  sync.Map
	TrackedMessages sync.Map
}

func NewMessageDispatcher(apiClient *slack.Client, wsClient *socketmode.Client, botID string, logger *log.Logger, aiClient *openai.Client) *MessageDispatcher {
	return &MessageDispatcher{
		APIClient:       apiClient,
		WSClient:        wsClient,
		AIClient:        aiClient,
		botID:           botID,
		Logger:          logger,
		TrackedThreads:  sync.Map{},
		TrackedMessages: sync.Map{},
	}
}

func (d *MessageDispatcher) Dispatch(event socketmode.Event) {
	switch event.Type {
	case socketmode.EventTypeEventsAPI:
		eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
		if !ok {
			d.Logger.Printf("Error asserting EventsAPIEvent: %v", event.Data)
			return
		}

		// Handling the callback event
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			switch eventsAPIEvent.InnerEvent.Type {
			case "app_mention":
				// Handle app_mention events
				d.handleAppMention(eventsAPIEvent)
			case "message":
				// Handle message events
				d.handleMessage(eventsAPIEvent)
			default:
				// Handle any unexpected inner event types
				d.handleUnhandledEvent(eventsAPIEvent)
			}
		} else {
			// Handle any unexpected top-level event types
			d.handleUnhandledEvent(eventsAPIEvent)
		}

		d.WSClient.Ack(*event.Request)
	default:
		// Log or handle completely unexpected event types
		d.Logger.Printf("Received an unhandled event type: %s", event.Type)
	}
}

func (d *MessageDispatcher) handleAppMention(event slackevents.EventsAPIEvent) {
	appMentionEvent, ok := event.InnerEvent.Data.(*slackevents.AppMentionEvent)

	if !ok {
		d.Logger.Printf("Error asserting AppMentionEvent: %v", event.InnerEvent.Data)
		return
	}
	threadID := appMentionEvent.ThreadTimeStamp
	if threadID == "" {
		threadID = appMentionEvent.TimeStamp
	}

	if _, ok := d.TrackedThreads.LoadOrStore(threadID, true); ok {
		// check if the message is already tracked
		return
	}
	_, _, err := d.APIClient.PostMessage(appMentionEvent.Channel, slack.MsgOptionText("Hello, thanks for mentioning me!", false), slack.MsgOptionTS(appMentionEvent.TimeStamp))
	if err != nil {
		d.Logger.Printf("Error posting warning message: %v", err)
		return
	}
}

func (d *MessageDispatcher) handleMessage(event slackevents.EventsAPIEvent) {
	messageEvent, ok := event.InnerEvent.Data.(*slackevents.MessageEvent)

	if !ok {
		d.Logger.Printf("Error asserting MessageEvent: %v", event.InnerEvent.Data)
		return
	}

	if _, ok := d.TrackedThreads.Load(messageEvent.ThreadTimeStamp); ok {
		// check if the message is already tracked
		if _, ok := d.TrackedMessages.Load(messageEvent.TimeStamp); ok {
			return
		}

		// prevent the bot from replying to its own messages
		if messageEvent.User == d.botID {
			return
		}

		_, _, err := d.APIClient.PostMessage(messageEvent.Channel, slack.MsgOptionText("Thanks for writing in this thread!", false), slack.MsgOptionTS(messageEvent.TimeStamp))
		if err != nil {
			d.Logger.Printf("Error posting warning message: %v", err)
			return
		}
		// Extract message timestamp and store it to track the message
		d.TrackedMessages.Store(messageEvent.TimeStamp, true)
	}
}

func (d *MessageDispatcher) handleUnhandledEvent(event slackevents.EventsAPIEvent) {
	// Log the unhandled event for review or take other appropriate actions
	d.Logger.Printf("Received an unhandled callback event type: %s", event.InnerEvent.Type)
}
