package utils

import (
	"context"
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/dispatcher"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient struct {
	ApiClient  *slack.Client
	WsClient   *socketmode.Client
	Dispatcher *dispatcher.EventDispatcher
	Logger     *log.Logger
	Context    context.Context
}

func NewSlackClient(ctx context.Context, apiToken, appToken string, logger *log.Logger) *SlackClient {
	api := slack.New(apiToken, slack.OptionAppLevelToken(appToken), slack.OptionLog(logger))
	wsclient := socketmode.New(api, socketmode.OptionLog(logger))

	return &SlackClient{
		ApiClient:  api,
		WsClient:   wsclient,
		Dispatcher: nil,
		Logger:     logger,
		Context:    ctx,
	}
}

func (c *SlackClient) RegisterNewEventDispatcher(dispatcher *dispatcher.EventDispatcher) {
	c.Dispatcher = dispatcher
}

func (c *SlackClient) Listen() {
	go func() {
		for {
			select {
			case <-c.Context.Done():
				c.Logger.Println("Context cancelled, stopping listener loop.")
				return
			case evt := <-c.WsClient.Events:
				go c.Dispatcher.Dispatch(evt)
			}
		}
	}()
	c.WsClient.Run()
}
