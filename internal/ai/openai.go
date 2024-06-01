package ai

import (
	"context"
	"log"
	"sync"

	"github.com/namikmesic/slack-ai-bot/internal/config"
	"github.com/sashabaranov/go-openai"
)

type OpenAiClient struct {
	OpenaiClient    *openai.Client // Ensuring pointer type for the client
	Context         context.Context
	Logger          *log.Logger // Including Logger for logging purposes
	TrackedMessages sync.Map
}

// NewOpenAIClient initializes a new client for interacting with OpenAI's API.
func NewOpenAIClient(ctx context.Context, cfg *config.AppConfig, logger *log.Logger) *OpenAiClient {
	// Extracting the OpenAI token from the configuration.
	aiClient := openai.NewClient(cfg.OpenAiToken) // Using configuration to get the OpenAI token

	return &OpenAiClient{
		OpenaiClient: aiClient,
		Context:      ctx,
		Logger:       logger, // Assigning the passed logger to the struct
	}
}

// RespondToMessage sends a message to OpenAI's API and returns the AI's response.
func (c *OpenAiClient) RespondToMessage(threadID string, role, message string) (string, bool) {
	// Load existing conversation or initialize it
	val, loaded := c.TrackedMessages.Load(threadID)
	if !loaded {
		val = []openai.ChatCompletionMessage{}
	}

	// Convert the loaded value to the expected slice type
	messages := val.([]openai.ChatCompletionMessage)

	// Append the new message to the conversation
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    role, // typically "user" for incoming messages
		Content: message,
	})

	// Always update the map with the new slice since the append might have created a new slice
	c.TrackedMessages.Store(threadID, messages)

	// Making the completion request to OpenAI
	resp, err := c.OpenaiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		},
	)
	if err != nil {
		c.Logger.Printf("error creating completion: %v", err)
		return "", false
	}

	// Returning the completion text
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		return resp.Choices[0].Message.Content, true
	}

	// In case there is no content or choices are empty
	return "", false
}

// AddMessage adds a new message to the tracked conversation identified by conversationID.
func (c *OpenAiClient) AddMessage(conversationID string, role, content string) {
	val, _ := c.TrackedMessages.LoadOrStore(conversationID, []openai.ChatCompletionMessage{})
	messages := append(val.([]openai.ChatCompletionMessage), openai.ChatCompletionMessage{
		Role:    role,    // Define the role of the message (e.g., "user", "assistant")
		Content: content, // The text content of the message
	})
	c.TrackedMessages.Store(conversationID, messages)
}

// GetMessages retrieves messages for a given conversation ID.
func (c *OpenAiClient) GetMessages(conversationID string) ([]openai.ChatCompletionMessage, bool) {
	val, ok := c.TrackedMessages.Load(conversationID)
	if !ok {
		return nil, false
	}
	return val.([]openai.ChatCompletionMessage), true
}
