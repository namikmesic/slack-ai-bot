package main

import (
	"context"
	"log"

	"github.com/namikmesic/slack-ai-bot/internal/ai"
	"github.com/namikmesic/slack-ai-bot/internal/config"
	"github.com/namikmesic/slack-ai-bot/internal/dispatcher"
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

	// Create OpenAI client
	openaiClient := ai.NewOpenAIClient(ctx, cfg, applogger)

	// Register the events handler
	messageDispatcher := dispatcher.NewMessageDispatcher(slackClient.APIClient, slackClient.WSClient, slackClient.BotID, slackClient.Logger, openaiClient)
	slackClient.RegisterNewMessageDispatcher(messageDispatcher)

	// Start listening for events
	slackClient.Listen()
}
