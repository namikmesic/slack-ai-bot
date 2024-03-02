package main

import (
	"context"
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/config"
	"github.com/namikmesic/slack-ai-bot/internal/dispatcher"
	"github.com/namikmesic/slack-ai-bot/internal/handlers"
	"github.com/namikmesic/slack-ai-bot/internal/logger"
	"github.com/namikmesic/slack-ai-bot/internal/utils"
)

func main() {
	cfg, err := config.LoadConfig([]string{".", "./config"}, "config")

	if err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	// Create a new logger
	applogger := logger.New()

	// Create a new context
	ctx := context.Background()

	// Create a new Slack client
	slackClient := utils.NewSlackClient(ctx, cfg, applogger)

	// Register the events handler
	eventsDispatcher := dispatcher.NewEventDispatcher(slackClient.APIClient, slackClient.WsClient, slackClient.Logger)

	// Register the AppMention handler
	appMentionHandler := handlers.NewAppMentionEventHandler(slackClient.APIClient)
	eventsDispatcher.RegisterNeweventsAPIEventHandler("app_mention", appMentionHandler)

	slackClient.RegisterNewEventDispatcher(eventsDispatcher)

	// Start listening for events
	slackClient.Listen()
}
