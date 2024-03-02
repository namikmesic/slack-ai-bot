package utils

import (
	"context"
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/dispatcher"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient struct {
	APIClient  *slack.Client
	WsClient   *socketmode.Client
	Dispatcher *dispatcher.EventDispatcher
	Logger     *log.Logger
	Context    context.Context
}

func NewSlackClient(ctx context.Context, apiToken, appToken string, isDev bool, logger *log.Logger) *SlackClient {
	api := slack.New(apiToken, slack.OptionAppLevelToken(appToken), slack.OptionLog(logger), slack.OptionDebug(isDev))
	wsclient := socketmode.New(api, socketmode.OptionLog(logger), socketmode.OptionDebug(isDev))

	return &SlackClient{
		APIClient:  api,
		WsClient:   wsclient,
		Dispatcher: nil,
		Logger:     logger,
		Context:    ctx,
	}
}

func (c *SlackClient) RegisterNewEventDispatcher(dp *dispatcher.EventDispatcher) {
	c.Dispatcher = dp
}

func (c *SlackClient) Listen() {
	go func() {
		for {
			select {
			case <-c.Context.Done():
				c.Logger.Println("context cancelled, stopping listener loop")
				return
			case evt := <-c.WsClient.Events:
				go c.Dispatcher.Dispatch(evt)
			}
		}
	}()

	if err := c.WsClient.Run(); err != nil {
		c.Logger.Fatalf("error running websocket client: %v", err)
	}
}
