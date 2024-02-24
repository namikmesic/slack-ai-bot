package main

import (
	"context"
	"log"
	"os"

	"github.com/namikmesic/slack-ai-bot/internal/config"
	"github.com/namikmesic/slack-ai-bot/internal/slack"
	"github.com/namikmesic/slack-ai-bot/internal/slack/handlers"
)

func main() {
	cfg, err := config.LoadConfig([]string{".", "./config"}, "config")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	ctx := context.Background()

	apiToken := cfg.SlackBotToken
	appToken := cfg.SlackAppToken

	slackClient := slack.NewClient(ctx, apiToken, appToken, logger)

	// Register the AppMentionHandler
	appMentionHandler := handlers.NewAppMentionEventHandler(slackClient.ApiClient)
	slackClient.RegisterHandler("app_mention", appMentionHandler)

	// Start listening for events
	slackClient.StartListeningLoop()
}
