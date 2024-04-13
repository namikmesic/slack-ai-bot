package utils

import (
	"context"
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/config"
	"github.com/namikmesic/slack-ai-bot/internal/dispatcher"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type SlackClient struct {
	APIClient  *slack.Client
	WSClient   *socketmode.Client
	BotID      string
	Dispatcher *dispatcher.MessageDispatcher
	Logger     *log.Logger
	Context    context.Context
}

func getBotID(api *slack.Client) (string, error) {
	// The AuthTest method checks authentication and returns basic information about the bot.
	authTest, err := api.AuthTest()
	if err != nil {
		// Return an empty string and the error if the AuthTest fails.
		return "", err
	}
	// The UserID from AuthTest is the bot's user ID.
	return authTest.UserID, nil
}

func NewSlackClient(ctx context.Context, cfg *config.AppConfig, logger *log.Logger) *SlackClient {
	api := slack.New(cfg.SlackBotToken, slack.OptionAppLevelToken(cfg.SlackAppToken), slack.OptionLog(logger), slack.OptionDebug(cfg.IsDevelopment))
	wsclient := socketmode.New(api, socketmode.OptionLog(logger), socketmode.OptionDebug(cfg.IsDevelopment))
	// Get the bot ID

	botID, err := getBotID(api)
	if err != nil {
		logger.Fatalf("error getting bot ID: %v", err)
	}

	return &SlackClient{
		APIClient:  api,
		WSClient:   wsclient,
		BotID:      botID,
		Dispatcher: nil,
		Logger:     logger,
		Context:    ctx,
	}
}

func (c *SlackClient) RegisterNewMessageDispatcher(dp *dispatcher.MessageDispatcher) {
	c.Dispatcher = dp
}

func (c *SlackClient) Listen() {
	go func() {
		for {
			select {
			case <-c.Context.Done():
				c.Logger.Println("context cancelled, stopping listener loop")
				return
			case evt := <-c.WSClient.Events:
				go c.Dispatcher.Dispatch(evt)
			}
		}
	}()

	if err := c.WSClient.Run(); err != nil {
		c.Logger.Fatalf("error running websocket client: %v", err)
	}
}
